package checkout

// pricing.go is the SINGLE canonical location for the §11.2 shipping fee tier formula.
// DO NOT compute shipping fees anywhere else in this service.
//
// BA-PR-002: shippingFee = 60 THB if subtotal < 1500 THB (exclusive); 0 if subtotal >= 1500 THB (inclusive).
// All values are in THB minor units (1 THB = 100 minor units, i.e. satang).

const (
	// ShippingFeeBelowThreshold is the shipping charge when the subtotal is below the free-shipping threshold.
	// 60 THB expressed in minor units.
	ShippingFeeBelowThreshold = int64(60_00) // 6000 minor = 60 THB

	// ShippingFeeThreshold is the inclusive boundary above which shipping is free.
	// 1500 THB expressed in minor units.
	ShippingFeeThreshold = int64(1500_00) // 150000 minor = 1500 THB
)

// ComputeShippingFee returns the shipping fee for the given subtotal (both in minor units).
// The boundary at exactly 1500 THB is inclusive → shippingFee = 0 (BA-PR-002).
func ComputeShippingFee(subtotalMinor int64) int64 {
	if subtotalMinor < ShippingFeeThreshold {
		return ShippingFeeBelowThreshold
	}
	return 0
}

// ComputeGrandTotal returns subtotal + shippingFee - couponDiscount, clamped to zero.
// couponDiscountMinor is always 0 in MVP but included for future-proofing without
// inventing business rules.
func ComputeGrandTotal(subtotalMinor, shippingFeeMinor, couponDiscountMinor int64) int64 {
	total := subtotalMinor + shippingFeeMinor - couponDiscountMinor
	if total < 0 {
		total = 0
	}
	return total
}

// minorToTHB converts minor units to whole-THB float64 for the response envelope.
func minorToTHB(minor int64) float64 {
	return float64(minor) / 100.0
}
