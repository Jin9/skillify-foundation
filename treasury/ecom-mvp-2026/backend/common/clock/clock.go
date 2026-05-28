// Package clock wraps time operations behind the Clock interface,
// enabling deterministic testing via mock implementations (Fowler: Clock Wrapper).
package clock

import (
	"fmt"
	"time"
)

const (
	// DefaultLayout is the default time format for parsing.
	DefaultLayout = "2006-01-02T15:04:05"

	// Bangkok timezone identifier.
	Bangkok = "Asia/Bangkok"
)

// Clock abstracts system-clock access so domain logic can be tested deterministically.
type Clock interface {
	Now() time.Time
	NowIn(location string) (time.Time, error)
	Parse(layout, value string) (time.Time, error)
	ParseIn(layout, value, location string) (time.Time, error)
}

type systemClock struct{}

var _ Clock = (*systemClock)(nil)

// New returns a Clock backed by the real system clock.
func New() Clock {
	return &systemClock{}
}

func (c *systemClock) Now() time.Time {
	return time.Now().UTC()
}

func (c *systemClock) NowIn(location string) (time.Time, error) {
	loc, err := loadLocation(location)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(loc), nil
}

func (c *systemClock) Parse(layout, value string) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, parseError(value, err)
	}
	return t.UTC(), nil
}

func (c *systemClock) ParseIn(layout, value, location string) (time.Time, error) {
	loc, err := loadLocation(location)
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return time.Time{}, parseError(value, err)
	}
	return t, nil
}

// NowUTC returns the current UTC time.
func NowUTC() time.Time {
	return time.Now().UTC()
}

// NowBangkok returns the current time in Asia/Bangkok.
func NowBangkok() (time.Time, error) {
	loc, err := loadLocation(Bangkok)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(loc), nil
}

// ParseUTC parses a time string using DefaultLayout in UTC.
func ParseUTC(value string) (time.Time, error) {
	t, err := time.Parse(DefaultLayout, value)
	if err != nil {
		return time.Time{}, parseError(value, err)
	}
	return t.UTC(), nil
}

// ParseBangkok parses a time string using DefaultLayout in Asia/Bangkok.
func ParseBangkok(value string) (time.Time, error) {
	loc, err := loadLocation(Bangkok)
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.ParseInLocation(DefaultLayout, value, loc)
	if err != nil {
		return time.Time{}, parseError(value, err)
	}
	return t, nil
}

// ToBangkok converts a time.Time to Asia/Bangkok timezone.
func ToBangkok(t time.Time) (time.Time, error) {
	loc, err := loadLocation(Bangkok)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

func loadLocation(name string) (*time.Location, error) {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, fmt.Errorf("clock: load location %q: %w", name, err)
	}
	return loc, nil
}

func parseError(value string, err error) error {
	return fmt.Errorf("clock: parse %q: %w", value, err)
}
