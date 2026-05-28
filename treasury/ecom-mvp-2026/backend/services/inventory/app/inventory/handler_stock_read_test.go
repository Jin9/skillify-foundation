package inventory

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/inventory/config"
)

// ─── fakeStockStorage ────────────────────────────────────────────────────────

// fakeStockStorage is a minimal StockStorage stub for handler tests.
// Only GetBySKU is used by StockReadHandler; other methods return nil / no-op.
type fakeStockStorage struct {
	row *StockLevel
	err error
}

func (f *fakeStockStorage) GetBySKU(_ context.Context, _ string) (*StockLevel, error) {
	return f.row, f.err
}

func (f *fakeStockStorage) GetBySKUForUpdate(_ context.Context, _ pgx.Tx, _ string) (*StockLevel, error) {
	return nil, nil
}

func (f *fakeStockStorage) GetManyBySKUsForUpdate(_ context.Context, _ pgx.Tx, _ []string) (map[string]*StockLevel, error) {
	return nil, nil
}

func (f *fakeStockStorage) DecrementAvailableIncrementReserved(_ context.Context, _ pgx.Tx, _ string, _ int) error {
	return nil
}

func (f *fakeStockStorage) IncrementAvailableDecrementReserved(_ context.Context, _ pgx.Tx, _ string, _ int) error {
	return nil
}

func (f *fakeStockStorage) DecrementReservedIncrementSold(_ context.Context, _ pgx.Tx, _ string, _ int) error {
	return nil
}

func (f *fakeStockStorage) DecrementSoldIncrementAvailable(_ context.Context, _ pgx.Tx, _ string, _ int) error {
	return nil
}

func (f *fakeStockStorage) AdjustQuantities(_ context.Context, _ pgx.Tx, _ string, _ int) error {
	return nil
}

func (f *fakeStockStorage) Insert(_ context.Context, _ pgx.Tx, _ string) error {
	return nil
}

// ─── fakeReservationStorage ──────────────────────────────────────────────────

type fakeReservationStorage struct{}

func (f *fakeReservationStorage) Insert(_ context.Context, _ pgx.Tx, _ *Reservation) error {
	return nil
}

func (f *fakeReservationStorage) FindByOrderID(_ context.Context, _ pgx.Tx, _ uuid.UUID) ([]*Reservation, error) {
	return nil, nil
}

func (f *fakeReservationStorage) FindByOrderIDForUpdate(_ context.Context, _ pgx.Tx, _ uuid.UUID) ([]*Reservation, error) {
	return nil, nil
}

func (f *fakeReservationStorage) FindExpiredCandidates(_ context.Context, _ pgx.Tx, _ int) ([]*Reservation, error) {
	return nil, nil
}

func (f *fakeReservationStorage) GetByIDForUpdate(_ context.Context, _ pgx.Tx, _ uuid.UUID) (*Reservation, error) {
	return nil, nil
}

func (f *fakeReservationStorage) MarkExpired(_ context.Context, _ pgx.Tx, _ uuid.UUID) error { return nil }

func (f *fakeReservationStorage) MarkCommitted(_ context.Context, _ pgx.Tx, _ uuid.UUID) error {
	return nil
}

func (f *fakeReservationStorage) MarkReleased(_ context.Context, _ pgx.Tx, _ uuid.UUID, _ ReleaseReason) error {
	return nil
}

// ─── fakeOutboxStorage ───────────────────────────────────────────────────────

type fakeOutboxStorage struct{}

func (f *fakeOutboxStorage) Insert(_ context.Context, _ pgx.Tx, _ *OutboxEvent) error { return nil }

func (f *fakeOutboxStorage) FetchUnpublished(_ context.Context, _ pgx.Tx, _ int) ([]*OutboxEvent, error) {
	return nil, nil
}

func (f *fakeOutboxStorage) MarkPublished(_ context.Context, _ pgx.Tx, _ interface{}) error {
	return nil
}

// ─── fakeConsumedStorage ─────────────────────────────────────────────────────

type fakeConsumedStorage struct{}

func (f *fakeConsumedStorage) TryInsert(_ context.Context, _ pgx.Tx, _ uuid.UUID, _ string) (bool, error) {
	return true, nil
}

// ─── helpers ─────────────────────────────────────────────────────────────────

// newTestService builds a *Service with the given stock storage and no-op stubs for the rest.
func newTestService(stock StockStorage) *Service {
	return &Service{
		StockStorage:    stock,
		ResvStorage:     &fakeReservationStorage{},
		OutboxStorage:   &fakeOutboxStorage{},
		ConsumedStorage: &fakeConsumedStorage{},
		Cfg:             config.Config{},
	}
}

// ─── Tests ───────────────────────────────────────────────────────────────────

func TestStockReadHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           string
		storageRow     *StockLevel
		storageErr     error
		wantHTTPStatus int
		wantCode       string
		wantSKU        string
	}{
		{
			name: "ok — sku found returns stock levels",
			body: `{"sku":"SKU-001"}`,
			storageRow: &StockLevel{
				SKU:          "SKU-001",
				AvailableQty: 10,
				ReservedQty:  2,
				SoldQty:      5,
			},
			wantHTTPStatus: http.StatusOK,
			wantCode:       string(CodeSuccess),
			wantSKU:        "SKU-001",
		},
		{
			name:           "404 — sku not found",
			body:           `{"sku":"NO-SUCH-SKU"}`,
			storageErr:     pgx.ErrNoRows,
			wantHTTPStatus: http.StatusNotFound,
			wantCode:       string(CodeNotFound),
		},
		{
			name:           "400 — missing sku field",
			body:           `{}`,
			wantHTTPStatus: http.StatusBadRequest,
			wantCode:       string(CodeBadRequest),
		},
		{
			name:           "400 — malformed JSON",
			body:           `{not-json}`,
			wantHTTPStatus: http.StatusBadRequest,
			wantCode:       string(CodeBadRequest),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(&fakeStockStorage{row: tt.storageRow, err: tt.storageErr})

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(
				http.MethodPost,
				"/api/v1/inventory/stock/read",
				bytes.NewBufferString(tt.body),
			)
			c.Request.Header.Set("Content-Type", "application/json")

			svc.StockReadHandler(c)

			if w.Code != tt.wantHTTPStatus {
				t.Fatalf("http status: want %d got %d (body=%s)", tt.wantHTTPStatus, w.Code, w.Body.String())
			}

			var env struct {
				Code string `json:"code"`
				Data *struct {
					SKU          string `json:"sku"`
					AvailableQty int    `json:"availableQty"`
					ReservedQty  int    `json:"reservedQty"`
					SoldQty      int    `json:"soldQty"`
				} `json:"data"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
				t.Fatalf("unmarshal: %v (body=%s)", err, w.Body.String())
			}
			if env.Code != tt.wantCode {
				t.Errorf("code: want %q got %q", tt.wantCode, env.Code)
			}
			if tt.wantHTTPStatus == http.StatusOK {
				if env.Data == nil {
					t.Fatal("expected data payload on 200 response")
				}
				if env.Data.SKU != tt.wantSKU {
					t.Errorf("sku: want %q got %q", tt.wantSKU, env.Data.SKU)
				}
				if env.Data.AvailableQty != tt.storageRow.AvailableQty {
					t.Errorf("availableQty: want %d got %d", tt.storageRow.AvailableQty, env.Data.AvailableQty)
				}
			}
		})
	}
}
