package app

import "time"

// SetTimestamps sets CreatedAt (if zero) and UpdatedAt to the current local time.
// Pass pointers to the struct's time fields.
func SetTimestamps(createdAt, updatedAt *time.Time) {
	now := time.Now().Local()
	if createdAt.IsZero() {
		*createdAt = now
	}
	*updatedAt = now
}
