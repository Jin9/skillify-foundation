package order

import (
	"testing"
)

// TestValidateTransition_Allowed verifies the 10 entries in allowedTransitions
// that are explicitly permitted. Each maps to a row in td.json §state_machine.
func TestValidateTransition_Allowed(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		from  OrderStatus
		to    OrderStatus
		actor ActorRole
	}{
		// 1. PENDING_PAYMENT → PAID  (SYSTEM, events.payment.completed)
		{
			name:  "PENDING_PAYMENT→PAID/SYSTEM",
			from:  StatusPendingPayment,
			to:    StatusPaid,
			actor: RoleSystem,
		},
		// 2. PENDING_PAYMENT → PAYMENT_FAILED  (SYSTEM, events.payment.failed)
		{
			name:  "PENDING_PAYMENT→PAYMENT_FAILED/SYSTEM",
			from:  StatusPendingPayment,
			to:    StatusPaymentFailed,
			actor: RoleSystem,
		},
		// 3. PENDING_PAYMENT → PAYMENT_EXPIRED  (SYSTEM, events.payment.expired / reservation.expired)
		{
			name:  "PENDING_PAYMENT→PAYMENT_EXPIRED/SYSTEM",
			from:  StatusPendingPayment,
			to:    StatusPaymentExpired,
			actor: RoleSystem,
		},
		// 4. PENDING_PAYMENT → CANCELLED  (CUSTOMER, order.cancel-mine)
		{
			name:  "PENDING_PAYMENT→CANCELLED/CUSTOMER",
			from:  StatusPendingPayment,
			to:    StatusCancelled,
			actor: RoleCustomer,
		},
		// 5. PENDING_PAYMENT → CANCELLED  (ADMIN, order.update-status-admin)
		{
			name:  "PENDING_PAYMENT→CANCELLED/ADMIN",
			from:  StatusPendingPayment,
			to:    StatusCancelled,
			actor: RoleAdmin,
		},
		// 6. PENDING_PAYMENT → CANCELLED  (SYSTEM, order.cancel-on-checkout-failure)
		{
			name:  "PENDING_PAYMENT→CANCELLED/SYSTEM",
			from:  StatusPendingPayment,
			to:    StatusCancelled,
			actor: RoleSystem,
		},
		// 7. PAID → PACKING  (ADMIN, order.update-status-admin)
		{
			name:  "PAID→PACKING/ADMIN",
			from:  StatusPaid,
			to:    StatusPacking,
			actor: RoleAdmin,
		},
		// 8. PACKING → SHIPPED  (ADMIN, order.update-status-admin; trackingNumber required)
		{
			name:  "PACKING→SHIPPED/ADMIN",
			from:  StatusPacking,
			to:    StatusShipped,
			actor: RoleAdmin,
		},
		// 9. SHIPPED → DELIVERED  (ADMIN, order.update-status-admin)
		{
			name:  "SHIPPED→DELIVERED/ADMIN",
			from:  StatusShipped,
			to:    StatusDelivered,
			actor: RoleAdmin,
		},
		// 10. PAID → CANCELLED  (ADMIN, order.update-status-admin; reason required)
		{
			name:  "PAID→CANCELLED/ADMIN",
			from:  StatusPaid,
			to:    StatusCancelled,
			actor: RoleAdmin,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := ValidateTransition(tc.from, tc.to, tc.actor); err != nil {
				t.Errorf("expected allowed transition to succeed, got: %v", err)
			}
		})
	}
}

// TestValidateTransition_Forbidden verifies 14 explicitly forbidden transitions.
// Categories:
//   - Terminal-state exits (PAYMENT_FAILED, PAYMENT_EXPIRED, CANCELLED, DELIVERED → anything)
//   - Wrong-actor on a valid (from,to) pair
//   - Non-existent (from,to) pair regardless of actor
//   - Self-transition
func TestValidateTransition_Forbidden(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		from  OrderStatus
		to    OrderStatus
		actor ActorRole
	}{
		// Terminal-state exits (4 terminal states → any forward move).
		// 1. PAYMENT_FAILED → PAID (SYSTEM) — terminal, no exit
		{
			name:  "PAYMENT_FAILED→PAID/SYSTEM terminal exit",
			from:  StatusPaymentFailed,
			to:    StatusPaid,
			actor: RoleSystem,
		},
		// 2. PAYMENT_EXPIRED → PAID (SYSTEM) — terminal, no exit
		{
			name:  "PAYMENT_EXPIRED→PAID/SYSTEM terminal exit",
			from:  StatusPaymentExpired,
			to:    StatusPaid,
			actor: RoleSystem,
		},
		// 3. CANCELLED → PENDING_PAYMENT (SYSTEM) — terminal, no exit
		{
			name:  "CANCELLED→PENDING_PAYMENT/SYSTEM terminal exit",
			from:  StatusCancelled,
			to:    StatusPendingPayment,
			actor: RoleSystem,
		},
		// 4. DELIVERED → CANCELLED (ADMIN) — terminal, no exit
		{
			name:  "DELIVERED→CANCELLED/ADMIN terminal exit",
			from:  StatusDelivered,
			to:    StatusCancelled,
			actor: RoleAdmin,
		},

		// Wrong-actor on a defined (from, to) pair.
		// 5. PENDING_PAYMENT → PAID — CUSTOMER cannot drive this (only SYSTEM via event)
		{
			name:  "PENDING_PAYMENT→PAID/CUSTOMER wrong actor",
			from:  StatusPendingPayment,
			to:    StatusPaid,
			actor: RoleCustomer,
		},
		// 6. PENDING_PAYMENT → PAID — ADMIN cannot drive this directly
		{
			name:  "PENDING_PAYMENT→PAID/ADMIN wrong actor",
			from:  StatusPendingPayment,
			to:    StatusPaid,
			actor: RoleAdmin,
		},
		// 7. PAID → PACKING — CUSTOMER cannot drive this (only ADMIN)
		{
			name:  "PAID→PACKING/CUSTOMER wrong actor",
			from:  StatusPaid,
			to:    StatusPacking,
			actor: RoleCustomer,
		},
		// 8. PACKING → SHIPPED — CUSTOMER cannot drive this (only ADMIN)
		{
			name:  "PACKING→SHIPPED/CUSTOMER wrong actor",
			from:  StatusPacking,
			to:    StatusShipped,
			actor: RoleCustomer,
		},

		// Non-existent (from, to) pair — not defined for any actor.
		// 9. PAID → PENDING_PAYMENT — no backward transition defined
		{
			name:  "PAID→PENDING_PAYMENT/ADMIN no backward transition",
			from:  StatusPaid,
			to:    StatusPendingPayment,
			actor: RoleAdmin,
		},
		// 10. PACKING → PAID — no backward transition defined
		{
			name:  "PACKING→PAID/ADMIN no backward transition",
			from:  StatusPacking,
			to:    StatusPaid,
			actor: RoleAdmin,
		},
		// 11. SHIPPED → PACKING — no backward transition defined
		{
			name:  "SHIPPED→PACKING/ADMIN no backward transition",
			from:  StatusShipped,
			to:    StatusPacking,
			actor: RoleAdmin,
		},
		// 12. PENDING_PAYMENT → DELIVERED — skip-ahead not defined
		{
			name:  "PENDING_PAYMENT→DELIVERED/ADMIN skip-ahead",
			from:  StatusPendingPayment,
			to:    StatusDelivered,
			actor: RoleAdmin,
		},

		// Self-transition — explicitly forbidden by ValidateTransition.
		// 13. PENDING_PAYMENT → PENDING_PAYMENT (SYSTEM)
		{
			name:  "PENDING_PAYMENT→PENDING_PAYMENT/SYSTEM self-transition",
			from:  StatusPendingPayment,
			to:    StatusPendingPayment,
			actor: RoleSystem,
		},
		// 14. PAID → PAID (ADMIN)
		{
			name:  "PAID→PAID/ADMIN self-transition",
			from:  StatusPaid,
			to:    StatusPaid,
			actor: RoleAdmin,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateTransition(tc.from, tc.to, tc.actor)
			if err == nil {
				t.Errorf("expected forbidden transition to return error, but got nil")
			}
			if _, ok := err.(*ErrInvalidTransition); !ok {
				t.Errorf("expected *ErrInvalidTransition, got %T: %v", err, err)
			}
		})
	}
}

// TestIsTerminal verifies the four terminal states and a non-terminal state.
func TestIsTerminal(t *testing.T) {
	t.Parallel()

	terminals := []OrderStatus{
		StatusPaymentFailed,
		StatusPaymentExpired,
		StatusCancelled,
		StatusDelivered,
	}
	for _, s := range terminals {
		s := s
		t.Run(string(s)+" is terminal", func(t *testing.T) {
			t.Parallel()
			if !IsTerminal(s) {
				t.Errorf("expected %s to be terminal", s)
			}
		})
	}

	nonTerminals := []OrderStatus{
		StatusPendingPayment,
		StatusPaid,
		StatusPacking,
		StatusShipped,
	}
	for _, s := range nonTerminals {
		s := s
		t.Run(string(s)+" is not terminal", func(t *testing.T) {
			t.Parallel()
			if IsTerminal(s) {
				t.Errorf("expected %s to be non-terminal", s)
			}
		})
	}
}

// TestRequiredFieldsFor verifies the RequiredFieldsFor helper returns the right fields.
func TestRequiredFieldsFor(t *testing.T) {
	t.Parallel()

	// PACKING→SHIPPED/ADMIN requires trackingNumber.
	fields := RequiredFieldsFor(StatusPacking, StatusShipped, RoleAdmin)
	if len(fields) != 1 || fields[0] != "trackingNumber" {
		t.Errorf("PACKING→SHIPPED/ADMIN: expected [trackingNumber], got %v", fields)
	}

	// Unknown transition returns nil.
	fields = RequiredFieldsFor(StatusDelivered, StatusPaid, RoleAdmin)
	if fields != nil {
		t.Errorf("unknown transition: expected nil, got %v", fields)
	}
}
