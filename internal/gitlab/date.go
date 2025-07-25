package gitlab

import (
	"strings"
	"time"
)

type Date struct {
	time.Time
}

const layout = "2006-01-02"

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return nil
	}

	t, err := time.Parse(layout, s)
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}
