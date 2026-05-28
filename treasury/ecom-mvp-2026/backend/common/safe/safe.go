// Package safe provides null-safe pointer dereferencing helpers.
//
// These utilities are useful when working with optional fields in
// DTOs, database results, or API responses that use *T pointer types.
package safe

// Deref safely dereferences a pointer, returning the zero value of T if nil.
//
//	var name *string = nil
//	safe.Deref(name) // returns ""
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// DerefOr safely dereferences a pointer, returning fallback if nil.
//
//	var count *int = nil
//	safe.DerefOr(count, 42) // returns 42
func DerefOr[T any](p *T, fallback T) T {
	if p == nil {
		return fallback
	}
	return *p
}

// Ptr returns a pointer to the given value. Useful for inline literals.
//
//	safe.Ptr("hello")  // returns *string
//	safe.Ptr(42)       // returns *int
func Ptr[T any](v T) *T {
	return &v
}

// IsNil returns true if the pointer is nil.
func IsNil[T any](p *T) bool {
	return p == nil
}
