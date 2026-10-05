package main

import (
	"bytes"
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

// mergeTaggedUnion turns the schema name, a union of objects told apart by the property tag, into one object: the
// properties the variants share, tag as an enum of their values, and each variant's own properties, optional. The
// union may nest unions, and may be one part of an allOf whose other parts are objects every variant shares. The
// schemas it merged are removed once nothing refers to them. Every reference to one of aliases, unions within it that
// stand for the variants without the allOf's other parts, is repointed at name, whose properties from those parts
// become optional. It fails if a variant is not a reference to an object, has no single value for tag, or disagrees
// with another on a property they share, so a change upstream is noticed.
func mergeTaggedUnion(doc *openapi.Document, name, tag string, aliases ...string) error {
	s, ok := doc.Components.Schemas[name]
	if !ok {
		return componentErr(name, errors.New("not found"))
	}

	common, union, merged, err := splitAllOf(s)
	if err != nil {
		return componentErr(name, err)
	}

	var names []string

	variants := []*openapi.Schema{}

	var collect func(openapi.SchemaList) error

	collect = func(alts openapi.SchemaList) error {
		for _, alt := range alts {
			if alt.Ref == nil || alt.Ref.Value == nil {
				return errors.New("a variant is not a reference")
			}

			n := strings.TrimPrefix(alt.Ref.Identifier, schemaRefPrefix)
			merged = append(merged, n)

			if sub := alternatives(alt.Ref.Value); len(sub) > 0 {
				if err := collect(sub); err != nil {
					return err
				}

				continue
			}

			if alt.Ref.Value.Type != openapi.TypeObject {
				return fmt.Errorf("variant %s is neither an object nor a union", n)
			}

			names = append(names, n)
			variants = append(variants, alt.Ref.Value)
		}

		return nil
	}

	if err := collect(alternatives(union)); err != nil {
		return componentErr(name, err)
	}

	if len(variants) == 0 {
		return componentErr(name, errors.New("not a union"))
	}

	tags := &jsonSet{}
	seen := map[string][]byte{}
	shared := map[string]int{tag: len(variants)}
	required := map[string]int{}

	commonProps := openapi.Schemas{}

	for _, c := range common {
		for prop, p := range c.Properties.ByIndex() {
			b, err := json.Marshal(p)
			if err != nil {
				return err
			}

			seen[prop] = b
			commonProps.Set(prop, p)
		}

		if len(aliases) == 0 {
			for _, r := range c.Required {
				required[r] = len(variants)
			}
		}
	}

	for i, v := range variants {
		if t, ok := v.Properties[tag]; !ok || len(t.Const) == 0 {
			return componentErr(names[i], &errpath.ErrField{Field: "properties", Err: &errpath.ErrKey{
				Key: tag, Err: errors.New("has no single value"),
			}})
		} else {
			tags.add(t)
		}

		for prop, p := range v.Properties.ByIndex() {
			if prop == tag {
				continue
			}

			b, err := json.Marshal(p)
			if err != nil {
				return err
			}

			if prev, ok := seen[prop]; ok && !bytes.Equal(prev, b) {
				return componentErr(names[i], &errpath.ErrField{Field: "properties", Err: &errpath.ErrKey{
					Key: prop, Err: errors.New("differs from another variant's"),
				}})
			}

			seen[prop] = b
			shared[prop]++
		}

		for _, r := range v.Required {
			required[r]++
		}
	}

	// the common parts' properties first, then what the variants share, in the first one's order, then each one's own
	props := openapi.Schemas{}
	for prop, p := range commonProps.ByIndex() {
		props.Set(prop, p)
	}

	for _, own := range []bool{false, true} {
		for _, v := range variants {
			for prop, p := range v.Properties.ByIndex() {
				switch _, done := props[prop]; {
				case done, own == (shared[prop] == len(variants)):
				case prop == tag:
					props.Set(tag, &openapi.Schema{Type: openapi.TypeString, Enum: tags.values})
				default:
					props.Set(prop, p)
				}
			}
		}
	}

	var req []string

	for prop := range props.ByIndex() {
		if required[prop] >= len(variants) {
			req = append(req, prop)
		}
	}

	s.Replace(&openapi.Schema{
		Title:       s.Title,
		Description: s.Description,
		Type:        openapi.TypeObject,
		Properties:  props,
		Required:    req,
	})

	redirect := map[string]string{}
	for _, a := range aliases {
		if !slices.Contains(merged, a) {
			return componentErr(name, fmt.Errorf("%s is not one of its unions", a))
		}

		redirect[a] = name
	}

	if err := edit.RedirectSchemas(doc, redirect); err != nil {
		return fmt.Errorf("redirecting schemas: %w", err)
	}

	return removeUnreferenced(doc, merged)
}

// splitAllOf returns the object parts of s's allOf and its one union, with the names of the parts it refers to, or s
// itself as the union when it has no allOf.
func splitAllOf(s *openapi.Schema) (common []*openapi.Schema, union *openapi.Schema, names []string, err error) {
	if len(s.AllOf) == 0 {
		return nil, s, nil, nil
	}

	for i, part := range s.AllOf {
		p := part
		if part.Ref != nil {
			p = part.Ref.Value
			names = append(names, strings.TrimPrefix(part.Ref.Identifier, schemaRefPrefix))
		}

		switch {
		case len(alternatives(p)) > 0 && union == nil:
			union = p
		case p.Type == openapi.TypeObject:
			common = append(common, p)
		default:
			return nil, nil, nil, &errpath.ErrField{Field: "allOf", Err: &errpath.ErrIndex{
				Index: i, Err: errors.New("is neither an object nor the one union"),
			}}
		}
	}

	if union == nil {
		return nil, nil, nil, errors.New("has no union in its allOf")
	}

	return common, union, names, nil
}

// alternatives are s's oneOf, or else its anyOf.
func alternatives(s *openapi.Schema) openapi.SchemaList {
	if len(s.OneOf) > 0 {
		return s.OneOf
	}

	return s.AnyOf
}

// removeUnreferenced removes those of the component schemas names nothing refers to, until each that remains is
// referred to, since removing one can leave another unreferenced.
func removeUnreferenced(doc *openapi.Document, names []string) error {
	for removed := true; removed; {
		removed = false

		spec := &bytes.Buffer{}
		if err := doc.WriteJSON(spec); err != nil {
			return err
		}

		for _, n := range names {
			if _, ok := doc.Components.Schemas[n]; ok &&
				!bytes.Contains(spec.Bytes(), []byte(`"`+schemaRefPrefix+n+`"`)) {
				delete(doc.Components.Schemas, n)

				removed = true
			}
		}
	}

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

	if err := addRequestID(doc, "Page", "Database"); err != nil {
		return nil, err
	}

	if err := allowDateTimes(doc, "dateResponse", "start", "end"); err != nil {
		return nil, err
	}

	for _, u := range []struct {
		name    string
		aliases []string
	}{
		{name: "Block"},
		// a rollup's array holds property values without their id
		{name: "PropertyValue", aliases: []string{"simpleOrArrayPropertyValueResponse"}},
		{name: "PropertyConfig"},
		{name: "RichText"},
	} {
		if err := mergeTaggedUnion(doc, u.name, "type", u.aliases...); err != nil {
			return nil, err
		}
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
