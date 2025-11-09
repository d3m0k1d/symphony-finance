package iso8601

import (
	"encoding/json"
	"time"
)

type Time struct {
	time.Time
}

const layout = "2006-01-02T15:04:05.999999"

func (t *Time) UnmarshalJSON(bits []byte) error {
	var s string
	json.Unmarshal(bits, &s)
	tt, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}
func (t Time) MarshalJSON() ([]byte, error) {
	s := t.Format(layout)
	return json.Marshal(s)
}
