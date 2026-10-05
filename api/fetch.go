package main

import (
	"cmp"
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/openapi"
	codegen "github.com/MarkRosemaker/openapi-codegen"
	compress "github.com/MarkRosemaker/openapi-compress"
	edit "github.com/MarkRosemaker/openapi-edit"
	enrich "github.com/MarkRosemaker/openapi-enrich"
	"github.com/MarkRosemaker/openapi-enrich/cassette"
	flatten "github.com/MarkRosemaker/openapi-flatten"
	"github.com/ettle/strcase"
	"golang.org/x/sync/errgroup"
)

const (
	pathOpenAPI     = "api/openapi.json"
	schemaRefPrefix = "#/components/schemas/"
)

var isClaudeCode = os.Getenv("CLAUDECODE") != ""

func main() {
	eg := errgroup.Group{}

	ctx := context.Background()
	eg.Go(func() error { return fetchLLMs(ctx) })
	eg.Go(func() error { return persistOpenAPI(ctx) })

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}

	doc, err := fixOpenAPI()
	if err != nil {
		log.Fatal(err)
	}

	if err := generateCode(doc); err != nil {
		log.Fatal(err)
	}
}

func fetchLLMs(ctx context.Context) error {
	if isClaudeCode {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://developers.notion.com/llms.txt", nil)
	if err != nil {
		return err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	f, err := os.Create("api/llms.txt")
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, rsp.Body); err != nil {
		return err
	}

	return nil
}

func readOfficial(ctx context.Context) (io.ReadCloser, error) {
	if isClaudeCode {
		return os.Open("api/openapi-official.json")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://developers.notion.com/openapi.json", nil)
	if err != nil {
		return nil, err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return rsp.Body, nil
}

func fetchOpenAPI(ctx context.Context) (*openapi.Document, error) {
	rc, err := readOfficial(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return openapi.LoadFromReader(rc)
}

func persistOpenAPI(ctx context.Context) error {
	doc, err := fetchOpenAPI(ctx)
	if err != nil {
		return err
	}

	var w io.Writer

	f, err := os.Create(pathOpenAPI)
	if err != nil {
		return err
	}
	defer f.Close()

	if isClaudeCode {
		w = f
	} else {
		f2, err := os.Create("api/openapi-official.json")
		if err != nil {
			return err
		}
		defer f2.Close()

		w = io.MultiWriter(f, f2)
	}

	return doc.WriteJSON(w)
}

// consolidateErrors replaces the schemas whose names start with any of prefixes, each an allOf of the common error and
// a code and status, with one schema called name that allows every code and status they did. Its code refers to a
// schema of its own, called codeName.
func consolidateErrors(doc *openapi.Document, prefixes []string, name, common, codeName string) error {
	var names []string

	for n := range doc.Components.Schemas.ByIndex() {
		if slices.ContainsFunc(prefixes, func(p string) bool { return strings.HasPrefix(n, p) }) {
			names = append(names, n)
		}
	}

	if len(names) == 0 {
		return fmt.Errorf("no schemas start with any of %q", prefixes)
	}

	commonSchema, ok := doc.Components.Schemas[common]
	if !ok {
		return fmt.Errorf("schema %q not found", common)
	}

	codes, statuses := &jsonSet{}, &jsonSet{}

	for _, n := range names {
		if err := collectError(doc.Components.Schemas[n], common, codes, statuses); err != nil {
			return componentErr(n, err)
		}
	}

	codeSchema := &openapi.Schema{Type: openapi.TypeString, Enum: codes.values}
	doc.Components.Schemas.Set(codeName, codeSchema)

	// set one by one, since a map literal leaves their order to chance
	props := openapi.Schemas{}
	props.Set("code", &openapi.Schema{Ref: &openapi.SchemaRef{Identifier: schemaRefPrefix + codeName, Value: codeSchema}})
	props.Set("status", &openapi.Schema{Type: openapi.TypeInteger, Enum: statuses.values})

	first := names[0]
	doc.Components.Schemas[first].Replace(&openapi.Schema{AllOf: openapi.SchemaList{
		{Ref: &openapi.SchemaRef{Identifier: schemaRefPrefix + common, Value: commonSchema}},
		{Type: openapi.TypeObject, Properties: props, Required: []string{"code", "status"}},
	}})

	redirect := map[string]string{}
	for _, n := range names[1:] {
		redirect[n] = first
	}

	if err := edit.RedirectSchemas(doc, redirect); err != nil {
		return fmt.Errorf("redirecting schemas: %w", err)
	}

	if err := edit.RenameSchema(doc, first, name); err != nil {
		return fmt.Errorf("renaming %s to %s: %w", first, name, err)
	}

	return nil
}

// collectError adds the codes and statuses s allows to codes and statuses, failing on anything it would not keep.
func collectError(s *openapi.Schema, common string, codes, statuses *jsonSet) error {
	if len(s.OneOf) > 0 {
		for i, alt := range s.OneOf {
			if err := collectError(alt, common, codes, statuses); err != nil {
				return &errpath.ErrField{Field: "oneOf", Err: &errpath.ErrIndex{Index: i, Err: err}}
			}
		}

		return nil
	}

	if len(s.AllOf) != 2 || s.AllOf[0].Ref == nil || s.AllOf[0].Ref.Identifier != schemaRefPrefix+common {
		return &errpath.ErrField{Field: "allOf", Err: fmt.Errorf("want %s and one schema", common)}
	}

	for prop, p := range s.AllOf[1].Properties.ByIndex() {
		switch prop {
		case "code":
			codes.add(p)
		case "status":
			statuses.add(p)
		case "additional_data": // a narrower shape of what the common error already allows
		default:
			return &errpath.ErrField{Field: "allOf", Err: &errpath.ErrIndex{Index: 1, Err: &errpath.ErrField{
				Field: "properties", Err: &errpath.ErrKey{Key: prop, Err: errors.New("unexpected property")},
			}}}
		}
	}

	return nil
}

// nameEnum gives s's enum values the Go names names, which must match them one for one. It fails if s already has
// extensions, so a change upstream that adds some is noticed instead of overwritten.
// See https://github.com/oapi-codegen/oapi-codegen/blob/main/docs/extensions.md#x-enum-varnames--x-enumnames.
func nameEnum(s *openapi.Schema, names ...string) error {
	if len(s.Enum) != len(names) {
		return &errpath.ErrField{Field: "enum", Err: fmt.Errorf("has %d values, want %d", len(s.Enum), len(names))}
	}

	if len(s.Extensions) > 0 {
		return fmt.Errorf("already has extensions: %s", s.Extensions)
	}

	ext, err := json.Marshal(map[string][]string{"x-enum-varnames": names})
	if err != nil {
		return err
	}

	s.Extensions = ext

	return nil
}

// propertyErr reports err as one of the property prop of the component schema name.
func propertyErr(name, prop string, err error) error {
	return componentErr(name, &errpath.ErrField{Field: "properties", Err: &errpath.ErrKey{Key: prop, Err: err}})
}

// componentErr reports err as one of the component schema name.
func componentErr(name string, err error) error {
	return &errpath.ErrField{Field: "components", Err: &errpath.ErrField{
		Field: "schemas", Err: &errpath.ErrKey{Key: name, Err: err},
	}}
}

// nameBranches moves each inline branch of a union into the component schemas, under its parent's
// name and what tells it apart from the other branches -- see [telling].
// flatten would name it after its title, which other unions' branches share, or its position, which says nothing.
func nameBranches(doc *openapi.Document) {
	for name, s := range doc.Components.Schemas.ByIndex() {
		nameBranchesIn(doc, s, name)
	}

	// the bodies, under the names flatten gives them
	inContent := func(c openapi.Content, name string) {
		for _, mt := range c {
			if mt.Schema != nil {
				nameBranchesIn(doc, mt.Schema, strcase.ToGoPascal(cmp.Or(mt.Schema.Title, name)))
			}
		}
	}

	for _, p := range doc.Paths.ByIndex() {
		for _, op := range p.Operations {
			if op.RequestBody != nil && op.RequestBody.Value != nil {
				inContent(op.RequestBody.Value.Content, op.OperationID)
			}

			for code, rsp := range op.Responses.ByIndex() {
				if rsp.Value != nil {
					inContent(rsp.Value.Content, op.OperationID+" "+cmp.Or(code.StatusText(), string(code)))
				}
			}
		}
	}
}

// nameBranchesIn is nameBranches for s, the schema flatten would name name, and what it holds.
func nameBranchesIn(doc *openapi.Document, s *openapi.Schema, name string) {
	if s == nil || s.Ref != nil {
		return
	}

	for _, union := range []struct {
		alts openapi.SchemaList
		kind string
	}{{s.OneOf, "OneOf"}, {s.AnyOf, "AnyOf"}} {
		for i, alt := range union.alts {
			branch := fmt.Sprintf("%s%s%d", name, union.kind, i)

			if v, ok := telling(union.alts, i); ok && alt.Ref == nil && alt.Type == openapi.TypeObject {
				if n := strcase.ToGoPascal(name + " " + v); doc.Components.Schemas[n] == nil {
					moved := new(openapi.Schema)
					moved.Replace(alt)
					doc.Components.Schemas.Set(n, moved)
					alt.Replace(&openapi.Schema{Ref: &openapi.SchemaRef{Identifier: schemaRefPrefix + n, Value: moved}})
					alt, branch = moved, n
				}
			}

			nameBranchesIn(doc, alt, branch)
		}
	}

	for i, part := range s.AllOf {
		nameBranchesIn(doc, part, fmt.Sprintf("%sAllOf%d", name, i))
	}

	for prop, p := range s.Properties.ByIndex() {
		nameBranchesIn(doc, p, strcase.ToGoPascal(name+" "+strings.ReplaceAll(prop, "/", " ")))
	}

	nameBranchesIn(doc, s.Items, name+"Item")

	if s.AdditionalProperties != nil {
		nameBranchesIn(doc, s.AdditionalProperties.Schema, name+"Value")
	}
}

// telling is what tells the branch alts[i] of a union apart from the others: a value it fixes a member to and no
// other branch does, the member type's first, or else the one member it requires.
func telling(alts openapi.SchemaList, i int) (string, bool) {
	s := alts[i]

	props := []string{"type"}
	for prop := range s.Properties.ByIndex() {
		props = append(props, prop)
	}

	for _, prop := range props {
		if v, ok := constString(s.Properties[prop]); ok && !slices.ContainsFunc(alts, func(other *openapi.Schema) bool {
			w, ok := constString(deref(other).Properties[prop])
			return other != s && ok && w == v
		}) {
			return v, true
		}
	}

	if len(s.Required) == 1 {
		return s.Required[0], true
	}

	return "", false
}

// constString is the string p fixes its value to, if it does.
func constString(p *openapi.Schema) (string, bool) {
	var v string
	return v, p != nil && json.Unmarshal(p.Const, &v) == nil
}

// deref is the schema s stands for: the one it refers to, if it is a reference.
func deref(s *openapi.Schema) *openapi.Schema {
	if s.Ref != nil && s.Ref.Value != nil {
		return s.Ref.Value
	}

	return s
}

// variants are the names of the component schemas the union name lists.
func variants(doc *openapi.Document, name string) ([]string, error) {
	s, ok := doc.Components.Schemas[name]
	if !ok {
		return nil, componentErr(name, errors.New("not found"))
	}

	var names []string

	for _, alt := range alternatives(s) {
		n, ok := strings.CutPrefix(alt.Ref.Identifier, schemaRefPrefix)
		if alt.Ref == nil || !ok {
			return nil, componentErr(name, errors.New("lists a variant that is not a component"))
		}

		names = append(names, n)
	}

	return names, nil
}

// flattenUnions lists, in each union among the schemas names or their allOf parts, the variants of the unions it lists
// instead of those unions, which tells the variants apart by their tag just the same.
func flattenUnions(doc *openapi.Document, names ...string) error {
	for _, n := range names {
		s, ok := doc.Components.Schemas[n]
		if !ok {
			return componentErr(n, errors.New("not found"))
		}

		for _, part := range append(openapi.SchemaList{s}, s.AllOf...) {
			if part.OneOf != nil {
				part.OneOf = leaves(part.OneOf)
			}

			if part.AnyOf != nil {
				part.AnyOf = leaves(part.AnyOf)
			}
		}
	}

	return nil
}

// leaves are the alternatives alts lists, with each that is a union replaced by its own leaves.
func leaves(alts openapi.SchemaList) openapi.SchemaList {
	var flat openapi.SchemaList

	for _, alt := range alts {
		v := alt
		if alt.Ref != nil {
			v = alt.Ref.Value
		}

		if sub := alternatives(v); len(sub) > 0 {
			flat = append(flat, leaves(sub)...)
		} else {
			flat = append(flat, alt)
		}
	}

	return flat
}

// alternatives are s's oneOf, or else its anyOf.
func alternatives(s *openapi.Schema) openapi.SchemaList {
	if len(s.OneOf) > 0 {
		return s.OneOf
	}

	return s.AnyOf
}

// addRequestID adds the request_id Notion returns on every top-level object, but the official spec leaves out, to the
// schemas names.
func addRequestID(doc *openapi.Document, names ...string) error {
	for _, n := range names {
		s, ok := doc.Components.Schemas[n]
		if !ok {
			return componentErr(n, errors.New("not found"))
		}

		if _, ok := s.Properties["request_id"]; ok {
			return propertyErr(n, "request_id", errors.New("already there"))
		}

		s.Properties.Set("request_id", &openapi.Schema{Type: openapi.TypeString, Format: openapi.FormatUUID})
	}

	return nil
}

// allowDateTimes drops the date format from the properties of the schema name, which Notion declares as dates but
// sends as a date or a date-time: "an ISO 8601 date, with optional time". It fails unless each of them is, or is one
// of, a string with the date format, so a change upstream is noticed.
func allowDateTimes(doc *openapi.Document, name string, props ...string) error {
	s, ok := doc.Components.Schemas[name]
	if !ok {
		return componentErr(name, errors.New("not found"))
	}

	for _, prop := range props {
		p, ok := s.Properties[prop]
		if !ok {
			return propertyErr(name, prop, errors.New("not found"))
		}

		dates := 0

		for _, alt := range append(openapi.SchemaList{p}, p.OneOf...) {
			if alt.Type == openapi.TypeString && alt.Format == openapi.FormatDate {
				alt.Format = ""
				dates++
			}
		}

		if dates == 0 {
			return propertyErr(name, prop, errors.New("is not a date"))
		}
	}

	return nil
}

// jsonSet holds the distinct values of a const or enum, in the order first seen.
type jsonSet struct {
	seen   map[string]bool
	values []jsontext.Value
}

func (s *jsonSet) add(p *openapi.Schema) {
	if s.seen == nil {
		s.seen = map[string]bool{}
	}

	vals := p.Enum
	if p.Const != nil {
		vals = append(vals, p.Const)
	}

	for _, v := range vals {
		if !s.seen[string(v)] {
			s.seen[string(v)] = true
			s.values = append(s.values, v)
		}
	}
}

func fixOpenAPI() (*openapi.Document, error) {
	doc, err := openapi.LoadFromFile(pathOpenAPI)
	if err != nil {
		return nil, err
	}

	for _, p := range doc.Paths.ByIndex() {
		for _, op := range p.Operations {
			for code, override := range map[openapi.StatusCode]string{ //nolint:exhaustive
				"200": "",
				"202": "",
				"400": "",
				"401": "",
				"403": "",
				"404": "",
				"406": "",
				"409": "",
				"429": "",
				"500": "",
				"503": "",
				"504": "",
				"529": "",
			} {
				rsp := op.Responses[code]
				if rsp == nil {
					continue
				}

				rsp.Value.Description = cmp.Or(
					rsp.Value.Description,
					override,
					code.StatusText(),
				)
			}
		}
	}

	if err := consolidateErrors(doc, []string{"error_api_", "error_oauth_"},
		"Error", "publicApiCommonErrorResponse", "ErrorCode"); err != nil {
		return nil, err
	}

	notionVersion, ok := doc.Components.Parameters["notionVersion"]
	if !ok {
		return nil, errors.New("parameter notionVersion not found")
	}

	if err := nameEnum(notionVersion.Value.Schema, "Current"); err != nil {
		return nil, &errpath.ErrField{Field: "components", Err: &errpath.ErrField{
			Field: "parameters", Err: &errpath.ErrKey{Key: "notionVersion", Err: &errpath.ErrField{Field: "schema", Err: err}},
		}}
	}

	if err := edit.RenameSchemas(doc, map[string]string{
		"blockObjectResponse":             "Block",
		"databaseObjectResponse":          "Database",
		"databasePropertyConfigResponse":  "PropertyConfig",
		"pageObjectResponse":              "Page",
		"pagePropertyValueWithIdResponse": "PropertyValue",
		"richTextItemResponse":            "RichText",
	}); err != nil {
		return nil, fmt.Errorf("renaming schemas: %w", err)
	}

	// the official spec spells out every array of RichText inline
	if err := edit.ExtractSchema(doc, "RichTexts", func(s *openapi.Schema) bool {
		return s.Type == openapi.TypeArray && s.Items != nil && s.Items.Ref != nil &&
			s.Items.Ref.Identifier == schemaRefPrefix+"RichText"
	}); err != nil {
		return nil, fmt.Errorf("naming the arrays of RichText: %w", err)
	}

	blocks, err := variants(doc, "Block")
	if err != nil {
		return nil, err
	}

	// each block variant, so that they keep differing only in their type and its member
	if err := addRequestID(doc, append([]string{"Page", "Database"}, blocks...)...); err != nil {
		return nil, err
	}

	if err := flattenUnions(doc, "PropertyValue"); err != nil {
		return nil, err
	}

	nameBranches(doc)

	if err := allowDateTimes(doc, "dateResponse", "start", "end"); err != nil {
		return nil, err
	}

	if err := applyPasses(doc); err != nil {
		return nil, err
	}

	if err := doc.Validate(); err != nil {
		return nil, fmt.Errorf("validating schema: %w", err)
	}

	if err := doc.WriteToFile(pathOpenAPI); err != nil {
		return nil, fmt.Errorf("writing to file: %w", err)
	}

	return doc, nil
}

// applyPasses runs, after the edits above, the passes the openapi-* commands would, in the order a spec goes through
// them, then sorts as they do.
func applyPasses(doc *openapi.Document) error {
	ias, err := cassette.InteractionsReadFile("api/interactions.json")
	if err != nil {
		return err
	}

	// the placeholder for the next interaction to record has no response
	ias = slices.DeleteFunc(ias, func(ia cassette.Interaction) bool { return ia.Response.StatusCode == 0 })

	flattenDoc := func(d *openapi.Document) error { return flatten.Document(d, flatten.Config{MarkOrigin: true}) }

	for _, pass := range []struct {
		name string
		run  func(*openapi.Document) error
	}{
		{"enrich", func(d *openapi.Document) error { return enrich.Enrich(d, ias) }},
		{"flatten", flattenDoc},
		{"compress", func(d *openapi.Document) error { return compress.Document(d, compress.Config{}) }},
		{"flatten again", flattenDoc},
	} {
		if err := pass.run(doc); err != nil {
			return fmt.Errorf("%s: %w", pass.name, err)
		}
	}

	for _, path := range doc.Paths {
		for _, op := range path.Operations {
			op.Responses.Sort()
		}
	}

	doc.Components.SortMaps()

	return nil
}

func generateCode(doc *openapi.Document) error {
	return codegen.Generate(codegen.Config{
		Spec:        doc,
		PackageName: "notion",
		OutputDir:   "pkg/notion",
		Types:       true,
		Client:      true,
		ClientTest:  true,
	})
}
