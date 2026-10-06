package notion

import (
	"encoding/json/v2"
	"testing"
	"time"
)

func TestDateOrDateTime(t *testing.T) {
	var d DateOrDateTime
	if err := json.Unmarshal([]byte(`"2026-10-05"`), &d); err != nil {
		t.Fatal(err)
	}

	if d.Date == nil || d.Time != nil || d.Date.String() != "2026-10-05" {
		t.Errorf("a date decoded as %+v", d)
	}

	d = DateOrDateTime{}
	if err := json.Unmarshal([]byte(`"2026-10-05T20:33:00.000+02:00"`), &d); err != nil {
		t.Fatal(err)
	}

	if d.Time == nil || d.Date != nil || d.Time.Format(time.RFC3339) != "2026-10-05T20:33:00+02:00" {
		t.Errorf("a date-time decoded as %+v", d)
	}

	start := time.Date(2026, 10, 5, 20, 33, 0, 0, time.FixedZone("", 2*60*60))

	b, err := json.Marshal(DateRequest{Start: DateOrDateTime{Time: &start}})
	if err != nil {
		t.Fatal(err)
	}

	if want := `{"start":"2026-10-05T20:33:00+02:00"}`; string(b) != want {
		t.Errorf("got %s, want %s", b, want)
	}
}
