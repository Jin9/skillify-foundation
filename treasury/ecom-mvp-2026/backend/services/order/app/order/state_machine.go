package order

import "fmt"

// transitionKey is the composite key used in the 2-D allowed-transitions map.
type transitionKey struct {
	From  OrderStatus
	To    OrderStatus
	Actor ActorRole
}

// transitionMeta carries human-readable context for a valid transition.
type transitionMeta struct {
	Trigger       string
	RequiredFields []string // e.g. ["trackingNumber"], ["reason"]
}

// allowedTransitions encodes the full state-machine table from td.json §state_machine.
// Key: (from, to, actor) → meta. Only explicit entries are permitted; everything
// else is implicitly forbidden.
//
// Encoding choice: 2-D map keyed by (from, to, actor) struct.
// Rationale: O(1) lookup, no iteration, single source of truth, easily audited
// against the TD table. A slice/switch alternative would need O(n) scan or manual
// indexing — no benefit at this cardinality (~11 rows).
var allowedTransitions = map[transitionKey]transitionMeta{
	// PENDING_PAYMENT → PAID  (SYSTEM via events.payment.completed)
	{StatusPendingPayment, StatusPaid, RoleSystem}: {
		Trigger: "events.payment.completed",
	},
	// PENDING_PAYMENT → PAYMENT_FAILED  (SYSTEM via events.payment.failed)
	{StatusPendingPayment, StatusPaymentFailed, RoleSystem}: {
		Trigger: "events.payment.failed",
	},
	// PENDING_PAYMENT → PAYMENT_EXPIRED  (SYSTEM via events.payment.expired OR events.reservation.expired)
	{StatusPendingPayment, StatusPaymentExpired, RoleSystem}: {
		Trigger:       "events.payment.expired OR events.reservation.expired",
		RequiredFields: []string{},
	},
	// PENDING_PAYMENT → CANCELLED  (CUSTOMER via order.cancel-mine)
	{StatusPendingPayment, StatusCancelled, RoleCustomer}: {
		Trigger: "order.cancel-mine",
	},
	// PENDING_PAYMENT → CANCELLED  (ADMIN via order.update-status-admin)
	{StatusPendingPayment, StatusCancelled, RoleAdmin}: {
		Trigger:       "order.update-status-admin",
		RequiredFields: []string{"reason"},
	},
	// PENDING_PAYMENT → CANCELLED  (SYSTEM via order.cancel-on-checkout-failure compensation)
	{StatusPendingPayment, StatusCancelled, RoleSystem}: {
		Trigger:       "order.cancel-on-checkout-failure",
		RequiredFields: []string{"reason"},
	},
	// PAID → PACKING  (ADMIN via order.update-status-admin)
	{StatusPaid, StatusPacking, RoleAdmin}: {
		Trigger: "order.update-status-admin",
	},
	// PACKING → SHIPPED  (ADMIN via order.update-status-admin; trackingNumber required — ORD-006)
	{StatusPacking, StatusShipped, RoleAdmin}: {
		Trigger:       "order.update-status-admin",
		RequiredFields: []string{"trackingNumber"},
	},
	// SHIPPED → DELIVERED  (ADMIN via order.update-status-admin)
	{StatusShipped, StatusDelivered, RoleAdmin}: {
		Trigger: "order.update-status-admin",
	},
	// PAID → CANCELLED  (ADMIN via order.update-status-admin; reason required — ORD-005, PR-005/PR-006)
	{StatusPaid, StatusCancelled, RoleAdmin}: {
		Trigger:       "order.update-status-admin",
		RequiredFields: []string{"reason"},
	},
}

// ErrInvalidTransition is returned by ValidateTransition for forbidden moves.
type ErrInvalidTransition struct {
	From   OrderStatus
	To     OrderStatus
	Actor  ActorRole
	Reason string
}

func (e *ErrInvalidTransition) Error() string {
	return fmt.Sprintf("invalid transition %s→%s for actor %s: %s", e.From, e.To, e.Actor, e.Reason)
}

// ValidateTransition returns nil when (from, to, actor) is in the allowed-transitions
// table, or an *ErrInvalidTransition with a human-readable reason otherwise.
//
// This is a pure function with no I/O. It is the single gate through which ALL
// state-changing paths must pass (handlers, consumers, compensation endpoint).
func ValidateTransition(from, to OrderStatus, actor ActorRole) error {
	// Terminal-state guard: fast-path for the three largest categories of forbidden
	// transitions (all of §state_machine.forbidden_transitions reachable from terminal states).
	switch from {
	case StatusPaymentFailed, StatusPaymentExpired, StatusCancelled, StatusDelivered:
		return &ErrInvalidTransition{
			From:   from,
			To:     to,
			Actor:  actor,
			Reason: fmt.Sprintf("%s is a terminal state; no further transitions are allowed (§9.3)", from),
		}
	}

	// No-op self-transition is not in the allowed table; callers must handle the
	// idempotent-retry case (same status) BEFORE calling ValidateTransition.
	if from == to {
		return &ErrInvalidTransition{
			From:   from,
			To:     to,
			Actor:  actor,
			Reason: "self-transition not permitted via state machine (handle idempotent retry at caller)",
		}
	}

	key := transitionKey{From: from, To: to, Actor: actor}
	if _, ok := allowedTransitions[key]; ok {
		return nil
	}

	// Produce a meaningful denial reason by checking whether the (from, to) pair
	// exists at all (wrong actor) vs. the pair being inherently forbidden.
	for k := range allowedTransitions {
		if k.From == from && k.To == to {
			return &ErrInvalidTransition{
				From:   from,
				To:     to,
				Actor:  actor,
				Reason: fmt.Sprintf("actor %s does not have authority to drive %s→%s (409 INVALID_ORDER_STATE)", actor, from, to),
			}
		}
	}

	return &ErrInvalidTransition{
		From:   from,
		To:     to,
		Actor:  actor,
		Reason: fmt.Sprintf("transition %s→%s is not defined in the allowed-transitions table (409 INVALID_ORDER_STATE)", from, to),
	}
}

// RequiredFieldsFor returns the list of required fields for a valid (from, to, actor)
// transition. Returns nil if the transition is not in the table.
func RequiredFieldsFor(from, to OrderStatus, actor ActorRole) []string {
	key := transitionKey{From: from, To: to, Actor: actor}
	if meta, ok := allowedTransitions[key]; ok {
		return meta.RequiredFields
	}
	return nil
}

// IsTerminal returns true if the given status is a terminal state.
func IsTerminal(s OrderStatus) bool {
	switch s {
	case StatusPaymentFailed, StatusPaymentExpired, StatusDelivered, StatusCancelled:
		return true
	}
	return false
}
