package notion

import "testing"

// TestQuickFilter_Documented reads every condition Notion documents for a filter, as a view's quick filter for a
// property takes it, and writes it back as it was, as an empty quick filter, which only adds the property to the
// filter bar.
func TestQuickFilter_Documented(t *testing.T) {
	t.Parallel()

	roundTrip[QuickFilterCondition](t, `{}`, `{}`)

	for _, d := range documented {
		for _, typ := range d.types {
			for _, c := range d.conditions {
				in := `{"` + typ + `":{` + c + `}}`
				roundTrip[QuickFilterCondition](t, in, in)
			}
		}
	}
}

// TestGroupBy reads each way Notion groups a view, by the type of the property it groups by, and writes it back as
// it was.
func TestGroupBy(t *testing.T) {
	t.Parallel()

	for _, in := range []string{
		`{"type":"select","property_id":"a","sort":{"type":"manual"},"hide_empty_groups":true}`,
		`{"type":"multi_select","property_id":"a","sort":{"type":"ascending"}}`,
		`{"type":"status","property_id":"a","group_by":"group","sort":{"type":"manual"}}`,
		`{"type":"person","property_id":"a","sort":{"type":"descending"}}`,
		`{"type":"created_by","property_id":"a","sort":{"type":"ascending"}}`,
		`{"type":"relation","property_id":"a","sort":{"type":"manual"}}`,
		`{"type":"date","property_id":"a","group_by":"week","sort":{"type":"ascending"},"start_day_of_week":1}`,
		`{"type":"title","property_id":"title","group_by":"alphabet_prefix","sort":{"type":"ascending"}}`,
		`{"type":"number","property_id":"a","sort":{"type":"ascending"},"range_start":0,"range_end":100,"range_size":10}`,
		`{"type":"checkbox","property_id":"a","sort":{"type":"manual"}}`,
		`{"type":"formula","property_id":"a","group_by":{"type":"date","group_by":"month","sort":{"type":"descending"}}}`,
		`{"type":"formula","property_id":"a","group_by":{"type":"number","sort":{"type":"ascending"},"range_size":5}}`,
		`{"type":"formula","property_id":"a","group_by":{"type":"checkbox","sort":{"type":"manual"}}}`,
	} {
		roundTrip[GroupByConfig](t, in, in)
	}
}
