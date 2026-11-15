package notion

import "testing"

// TestUpdateADataSource_Documented reads each change of a data source's schema Notion documents, and writes it back
// as it was: renaming a property, the options of a select, a multi-select or a status, by name or by ID, a number's
// format, a formula, a relation, a rollup by the names or the IDs of its properties, a unique ID's prefix, and a
// property turned into another type. The client says the type of each property and relation itself.
func TestUpdateADataSource_Documented(t *testing.T) {
	t.Parallel()

	for _, properties := range []string{
		`{"Old name":{"name":"New name"}}`,
		`{"Stage":{"type":"select","select":{"options":[{"name":"Done","color":"green"},{"id":"abc","name":"Doing"}]}}}`,
		`{"Tags":{"type":"multi_select","multi_select":{"options":[{"name":"A","color":"red","description":"the first"}]}}}`,
		`{"State":{"type":"status","status":{"options":[{"name":"Blocked","color":"red"},{"id":"xyz"}]}}}`,
		`{"Cost":{"type":"number","number":{"format":"dollar"}}}`,
		`{"Total":{"type":"formula","formula":{"expression":"prop(\"Cost\") * 2"}}}`,
		`{"Tasks":{"type":"relation","relation":{"data_source_id":"dd456007-6c66-4bba-957e-ea501dcda3a6","type":"single_property","single_property":{}}}}`,
		`{"Spent":{"type":"rollup","rollup":{"function":"sum","relation_property_name":"Tasks","rollup_property_name":"Cost"}}}`,
		`{"Spent":{"type":"rollup","rollup":{"function":"average","relation_property_id":"a%3Fb","rollup_property_id":"c%3Dd"}}}`,
		`{"ID":{"type":"unique_id","unique_id":{"prefix":"TASK"}}}`,
		`{"Notes":{"name":"Remarks","type":"rich_text","rich_text":{}}}`,
		`{"Done":{"type":"checkbox","checkbox":{}},"Due":{"type":"date","date":{}},"Owner":{"type":"people","people":{}}}`,
	} {
		in := `{"properties":` + properties + `}`
		roundTrip[UpdateADataSource](t, in, in)
	}
}
