package utils

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Date wraps time.Time to handle date-only and datetime formats
type Date struct {
	time.Time
}

// UnmarshalJSON parses a date from JSON (supports date and datetime)
func (d *Date) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "null" || str == "" {
		return nil
	}

	// Try parsing as date first (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", str); err == nil {
		d.Time = t
		return nil
	}

	// Try parsing as datetime (RFC3339)
	if t, err := time.Parse(time.RFC3339, str); err == nil {
		d.Time = t
		return nil
	}

	// Try parsing as datetime without timezone
	if t, err := time.Parse("2006-01-02T15:04:05", str); err == nil {
		d.Time = t
		return nil
	}

	return fmt.Errorf("cannot parse date: %s", str)
}

// MarshalJSON outputs the date as YYYY-MM-DD for JSON
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Time.Format("2006-01-02") + `"`), nil
}

// Value returns the date for database storage (driver.Valuer)
func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

// Scan reads a date value from the database (sql.Scanner)
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	}

	return fmt.Errorf("cannot scan %T into Date", value)
}
