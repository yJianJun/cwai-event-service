package utils

import (
	"encoding/json"
	"time"
)

type DateTime int64

func NewDateTime(t time.Time) DateTime {
	return DateTime(t.Unix())
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	t := d.AsTime()
	formatted := t.Local().Format(time.RFC3339) // 2020-06-18T16:19:27+08:00
	return json.Marshal(formatted)
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	layout := `"` + time.RFC3339 + `"`
	if parsed, err := time.Parse(layout, string(data)); err != nil {
		return err
	} else {
		*t = DateTime(parsed.Unix())
		return nil
	}
}

func (d DateTime) String() string {
	return d.AsTime().Local().Format(time.RFC3339)
}

func (d DateTime) AsTime() time.Time {
	return time.Unix(int64(d), 0)
}

func Date() DateTime {
	return NewDateTime(time.Now())
}
