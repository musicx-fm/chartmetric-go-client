package chartmetric

import (
	"fmt"
	"strings"
	"time"
)

const DateFormat = "2006-01-02"

// Date is a custom Time type for handling date strings in YYYY-MM-DD format.
type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		return nil
	}

	parsed, err := time.Parse(DateFormat, s)
	if err != nil {
		return fmt.Errorf("parse date: %w", err)
	}

	d.Time = parsed

	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Date) String() string {
	return fmt.Sprintf("%q", d.Time.Format(DateFormat))
}
