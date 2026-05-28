package cart

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
)

// ClearIdempotencyRecord holds the result of a previously processed clear request.
// Defined here (not in access) so that CartStoragePort can reference it without
// creating an import cycle: cart ← access ← cart.
type ClearIdempotencyRecord struct {
	OrderID uuid.UUID
	Removed int
}

// CartStoragePort is the storage interface used by Service.
// Defined in the cart package so that access (which imports cart for domain types)
// can implement it without creating an import cycle.
// Methods that span a transaction accept a pgx.Tx; the port owns BeginTx so callers
// never need to import the access or pgxpool packages directly.
type CartStoragePort interface {
	// Cart lifecycle
	GetOrCreateByUserID(ctx context.Context, userID uuid.UUID) (Cart, error)
	GetCartByUserID(ctx context.Context, userID uuid.UUID) (Cart, error)
	GetCartWithItems(ctx context.Context, userID uuid.UUID) (Cart, []CartItem, error)
	TouchUpdatedAt(ctx context.Context, tx pgx.Tx, cartID uuid.UUID) error

	// Cart items
	GetItemByID(ctx context.Context, cartItemID, cartID uuid.UUID) (CartItem, error)
	UpsertItem(ctx context.Context, tx pgx.Tx, cartID, productID uuid.UUID, qty int) error
	UpdateItemQty(ctx context.Context, tx pgx.Tx, cartItemID, cartID uuid.UUID, qty int) (int64, error)
	DeleteItem(ctx context.Context, tx pgx.Tx, cartItemID, cartID uuid.UUID) (int64, error)
	DeleteItems(ctx context.Context, tx pgx.Tx, cartID uuid.UUID, itemIDs []uuid.UUID) (int64, error)

	// Idempotency (clear-on-checkout)
	GetClearIdempotency(ctx context.Context, orderID uuid.UUID) (ClearIdempotencyRecord, error)
	SaveClearIdempotency(ctx context.Context, tx pgx.Tx, orderID uuid.UUID, removed int) error

	// Transaction management
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

// Sentinel errors returned by service methods. Handlers map these to HTTP codes.
var (
	ErrProductNotFound   = errors.New("cart: product not found")
	ErrProductInactive   = errors.New("cart: product inactive")
	ErrProductDeleted    = errors.New("cart: product deleted")
	ErrCartItemNotFound  = errors.New("cart: cart item not found")
	ErrStockInsufficient = errors.New("cart: stock insufficient")
)

// StockInsufficientDetail carries the available qty when stock is exceeded.
type StockInsufficientDetail struct {
	AvailableQty int
}

// ErrWithDetail wraps a sentinel with additional context for the handler.
type ErrWithDetail struct {
	Sentinel error
	Detail   any
}

func (e *ErrWithDetail) Error() string { return e.Sentinel.Error() }
func (e *ErrWithDetail) Unwrap() error { return e.Sentinel }

// Service orchestrates all cart business logic.
type Service struct {
	storage   CartStoragePort
	catalog   *CatalogClient
	inventory *InventoryClient
	fanoutCfg FanoutConfig
}

// FanoutConfig controls parallel catalog+inventory enrichment.
type FanoutConfig struct {
	TimeoutMS   int
	MaxItems    int
	Concurrency int
}

func NewService(storage CartStoragePort, catalog *CatalogClient, inventory *InventoryClient, fanoutCfg FanoutConfig) *Service {
	return &Service{
		storage:   storage,
		catalog:   catalog,
		inventory: inventory,
		fanoutCfg: fanoutCfg,
	}
}

// AddItem validates, upserts the item into the cart (ON CONFLICT merge), and returns the enriched cart.
// No stock check — stock guard lives at cart.update-item and checkout (per INV-002 / AMB-001).
func (s *Service) AddItem(ctx context.Context, userID, productID uuid.UUID, qty int) (CartView, error) {
	// 1. Verify the product exists and is ACTIVE via catalog.
	detail, err := s.catalog.GetProductDetail(ctx, productID)
	if err != nil {
		if isUpstreamTimeout(err) {
			return CartView{}, fmt.Errorf("catalog: %w", err)
		}
		return CartView{}, fmt.Errorf("catalog.GetProductDetail: %w", err)
	}
	switch detail.Status {
	case ProductStatusDeleted:
		return CartView{}, ErrProductDeleted
	case ProductStatusInactive:
		return CartView{}, ErrProductInactive
	case ProductStatusActive:
		// OK
	default:
		return CartView{}, ErrProductInactive
	}

	// 2. Ensure the cart row exists for this user.
	c, err := s.storage.GetOrCreateByUserID(ctx, userID)
	if err != nil {
		return CartView{}, serror.Wrap(err)
	}

	// 3. Upsert item in a transaction (merge qty on conflict) and touch updated_at.
	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return CartView{}, serror.Wrap(err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	if err = s.storage.UpsertItem(ctx, tx, c.CartID, productID, qty); err != nil {
		return CartView{}, serror.Wrap(err)
	}
	if err = s.storage.TouchUpdatedAt(ctx, tx, c.CartID); err != nil {
		return CartView{}, serror.Wrap(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return CartView{}, serror.Wrap(err)
	}

	// 4. Return enriched cart view.
	return s.ReadCart(ctx, userID)
}

// UpdateItem validates qty, checks stock (if qty>0), updates the cart item, returns enriched cart.
func (s *Service) UpdateItem(ctx context.Context, userID, cartItemID uuid.UUID, qty int) (CartView, error) {
	// 1. Fetch cart for this user.
	c, err := s.storage.GetCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CartView{}, ErrCartItemNotFound
		}
		return CartView{}, serror.Wrap(err)
	}

	// 2. Check the item belongs to this cart.
	item, err := s.storage.GetItemByID(ctx, cartItemID, c.CartID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CartView{}, ErrCartItemNotFound
		}
		return CartView{}, serror.Wrap(err)
	}

	// 3. qty=0 → treat as remove (CART-003 / AMB-004).
	if qty == 0 {
		return s.removeItemInternal(ctx, userID, cartItemID, c.CartID)
	}

	// 4. Stock-check via inventory (CART-002 — UX guard, NOT a reservation).
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(s.fanoutCfg.TimeoutMS)*time.Millisecond)
	defer cancel()

	availableQty, err := s.inventory.SingleStock(timeoutCtx, item.ProductID)
	if err != nil {
		if isUpstreamTimeout(err) {
			return CartView{}, fmt.Errorf("inventory: %w", err)
		}
		return CartView{}, serror.Wrap(err)
	}

	if qty > availableQty {
		return CartView{}, &ErrWithDetail{
			Sentinel: ErrStockInsufficient,
			Detail:   StockInsufficientDetail{AvailableQty: availableQty},
		}
	}

	// 5. UPDATE qty inside a transaction.
	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return CartView{}, serror.Wrap(err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	rows, err := s.storage.UpdateItemQty(ctx, tx, cartItemID, c.CartID, qty)
	if err != nil {
		return CartView{}, serror.Wrap(err)
	}
	if rows == 0 {
		return CartView{}, ErrCartItemNotFound
	}
	if err = s.storage.TouchUpdatedAt(ctx, tx, c.CartID); err != nil {
		return CartView{}, serror.Wrap(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return CartView{}, serror.Wrap(err)
	}

	return s.ReadCart(ctx, userID)
}

// RemoveItem deletes a cart item and returns the enriched cart (idempotent — no-op if already gone).
func (s *Service) RemoveItem(ctx context.Context, userID, cartItemID uuid.UUID) (CartView, error) {
	c, err := s.storage.GetCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// No cart at all — return empty cart view.
			return emptyCartView(userID), nil
		}
		return CartView{}, serror.Wrap(err)
	}
	return s.removeItemInternal(ctx, userID, cartItemID, c.CartID)
}

// removeItemInternal is the shared delete-and-touch path used by both RemoveItem and UpdateItem(qty=0).
func (s *Service) removeItemInternal(ctx context.Context, userID, cartItemID, cartID uuid.UUID) (CartView, error) {
	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return CartView{}, serror.Wrap(err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	rows, err := s.storage.DeleteItem(ctx, tx, cartItemID, cartID)
	if err != nil {
		return CartView{}, serror.Wrap(err)
	}
	if rows > 0 {
		slog.Info("cart: item removed", slog.String("cart_item_id", cartItemID.String()), slog.String("cart_id", cartID.String()))
	}
	if err = s.storage.TouchUpdatedAt(ctx, tx, cartID); err != nil {
		return CartView{}, serror.Wrap(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return CartView{}, serror.Wrap(err)
	}

	return s.ReadCart(ctx, userID)
}

// ReadCart loads the cart, fans out to catalog+inventory in parallel (errgroup + semaphore),
// and returns the enriched CartView with server-computed subtotal.
func (s *Service) ReadCart(ctx context.Context, userID uuid.UUID) (CartView, error) {
	c, items, err := s.storage.GetCartWithItems(ctx, userID)
	if err != nil {
		return CartView{}, serror.Wrap(err)
	}

	// No cart or empty cart — return zero-value view.
	if c.CartID == uuid.Nil || len(items) == 0 {
		return CartView{
			CustomerUserID: userID,
			Items:          []CartItemView{},
			Subtotal:       0,
			Currency:       "THB",
			UpdatedAt:      c.UpdatedAt,
		}, nil
	}

	// Cap at FANOUT_MAX_ITEMS to bound fan-out (CART-004 bounded list).
	if len(items) > s.fanoutCfg.MaxItems {
		items = items[:s.fanoutCfg.MaxItems]
	}

	// Parallel fan-out: catalog per-item + single inventory bulk-read.
	// errgroup with a timeout context; semaphore limits concurrent catalog calls.
	fanoutTimeout := time.Duration(s.fanoutCfg.TimeoutMS) * time.Millisecond
	fanoutCtx, cancel := context.WithTimeout(ctx, fanoutTimeout)
	defer cancel()

	productIDs := make([]uuid.UUID, len(items))
	for i, item := range items {
		productIDs[i] = item.ProductID
	}

	// Catalog results — one slot per item.
	details := make([]ProductDetail, len(items))
	// Inventory results — filled by the bulk-read goroutine.
	var stockMap map[uuid.UUID]int
	var stockMu sync.Mutex

	g, gCtx := errgroup.WithContext(fanoutCtx)

	// Semaphore channel to cap concurrent catalog calls.
	sem := make(chan struct{}, s.fanoutCfg.Concurrency)

	for i, item := range items {
		i, item := i, item // capture
		g.Go(func() error {
			sem <- struct{}{}
			defer func() { <-sem }()

			d, err := s.catalog.GetProductDetail(gCtx, item.ProductID)
			if err != nil {
				if isUpstreamTimeout(err) {
					// Per spec: on catalog timeout → that line carries available:false, reason:"catalog_unreachable".
					// We do NOT fail the whole read; store a zero detail with a sentinel status.
					details[i] = ProductDetail{
						ProductID: item.ProductID,
						Status:    "CATALOG_UNREACHABLE",
					}
					return nil
				}
				return err
			}
			details[i] = d
			return nil
		})
	}

	// Single inventory bulk-read in parallel with catalog calls.
	g.Go(func() error {
		m, err := s.inventory.BulkReadStock(gCtx, productIDs)
		if err != nil {
			// Inventory failure degrades gracefully: all items get availableQty=0.
			slog.WarnContext(gCtx, "inventory bulk-read failed, degrading gracefully", slog.Any("error", err))
			stockMu.Lock()
			stockMap = make(map[uuid.UUID]int)
			stockMu.Unlock()
			return nil
		}
		stockMu.Lock()
		stockMap = m
		stockMu.Unlock()
		return nil
	})

	if err := g.Wait(); err != nil {
		return CartView{}, fmt.Errorf("cart.read fan-out: %w", err)
	}

	// Enrich items.
	views := make([]CartItemView, len(items))
	var subtotal float64

	for i, item := range items {
		d := details[i]
		availQty := stockMap[item.ProductID]

		checkoutable := d.Status == ProductStatusActive && availQty > 0
		lineSubtotal := d.CurrentPrice * float64(item.Qty)

		// Catalog-unreachable sentinel: flag as unavailable.
		productStatus := d.Status
		if productStatus == "CATALOG_UNREACHABLE" {
			checkoutable = false
		}

		views[i] = CartItemView{
			CartItemID:    item.CartItemID,
			ProductID:     item.ProductID,
			SKU:           d.SKU,
			Name:          d.Name,
			Image:         d.Image,
			CurrentPrice:  d.CurrentPrice,
			Qty:           item.Qty,
			LineSubtotal:  lineSubtotal,
			AvailableQty:  availQty,
			ProductStatus: productStatus,
			Checkoutable:  checkoutable,
		}

		if checkoutable {
			subtotal += lineSubtotal
		}
	}

	return CartView{
		CartID:         c.CartID,
		CustomerUserID: c.UserID,
		Items:          views,
		Subtotal:       subtotal,
		Currency:       "THB",
		UpdatedAt:      c.UpdatedAt,
	}, nil
}

// ClearOnCheckout removes specified cart items idempotently (keyed by orderId).
func (s *Service) ClearOnCheckout(ctx context.Context, userID uuid.UUID, cartItemIDs []uuid.UUID, orderID uuid.UUID) (int, error) {
	// 1. Check idempotency cache.
	rec, err := s.storage.GetClearIdempotency(ctx, orderID)
	if err == nil {
		// Already processed — return cached result.
		return rec.Removed, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, serror.Wrap(err)
	}

	// 2. Fetch cart_id for this user.
	c, err := s.storage.GetCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// No cart → nothing to clear; persist idempotency with removed=0.
			return 0, s.persistClearIdempotencyNoTx(ctx, orderID, 0)
		}
		return 0, serror.Wrap(err)
	}

	// 3. Delete specified items + touch updated_at + save idempotency in a single tx.
	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return 0, serror.Wrap(err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	removed, err := s.storage.DeleteItems(ctx, tx, c.CartID, cartItemIDs)
	if err != nil {
		return 0, serror.Wrap(err)
	}
	if err = s.storage.TouchUpdatedAt(ctx, tx, c.CartID); err != nil {
		return 0, serror.Wrap(err)
	}
	if err = s.storage.SaveClearIdempotency(ctx, tx, orderID, int(removed)); err != nil {
		return 0, serror.Wrap(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, serror.Wrap(err)
	}

	return int(removed), nil
}

// persistClearIdempotencyNoTx saves an idempotency record outside a transaction (no-cart case).
func (s *Service) persistClearIdempotencyNoTx(ctx context.Context, orderID uuid.UUID, removed int) error {
	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err = s.storage.SaveClearIdempotency(ctx, tx, orderID, removed); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// emptyCartView returns a zero-value CartView for a user with no cart.
func emptyCartView(userID uuid.UUID) CartView {
	return CartView{
		CustomerUserID: userID,
		Items:          []CartItemView{},
		Subtotal:       0,
		Currency:       "THB",
		UpdatedAt:      time.Time{},
	}
}

// isUpstreamTimeout detects context deadline or cancellation errors from HTTP client calls.
func isUpstreamTimeout(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}
