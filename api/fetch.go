package main

import (
	"cmp"
	"context"
	"encoding/json/jsontext"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MarkRosemaker/openapi"
	edit "github.com/MarkRosemaker/openapi-edit"
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

	if err := fixOpenAPI(); err != nil {
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

// consolidateErrors replaces the schemas whose names start with prefix, each an allOf of the common error and a code
// and status, with one schema called name that allows every code and status they did.
func consolidateErrors(doc *openapi.Document, prefix, name, common string) error {
	var names []string

	for n := range doc.Components.Schemas.ByIndex() {
		if strings.HasPrefix(n, prefix) {
			names = append(names, n)
		}
	}

	if len(names) == 0 {
		return fmt.Errorf("no schemas start with %q", prefix)
	}

	commonSchema, ok := doc.Components.Schemas[common]
	if !ok {
		return fmt.Errorf("schema %q not found", common)
	}

	codes, statuses := &jsonSet{}, &jsonSet{}

	for _, n := range names {
		if err := collectError(doc.Components.Schemas[n], common, codes, statuses); err != nil {
			return fmt.Errorf("%s: %w", n, err)
		}
	}

	first := names[0]
	*doc.Components.Schemas[first] = openapi.Schema{AllOf: openapi.SchemaList{
		{Ref: &openapi.SchemaRef{Identifier: schemaRefPrefix + common, Value: commonSchema}},
		{
			Type: openapi.TypeObject,
			Properties: openapi.Schemas{
				"code":   {Type: openapi.TypeString, Enum: codes.values},
				"status": {Type: openapi.TypeInteger, Enum: statuses.values},
			},
			Required: []string{"code", "status"},
		},
	}}

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
				return fmt.Errorf("oneOf[%d]: %w", i, err)
			}
		}

		return nil
	}

	if len(s.AllOf) != 2 || s.AllOf[0].Ref == nil || s.AllOf[0].Ref.Identifier != schemaRefPrefix+common {
		return fmt.Errorf("want allOf of %s and one schema", common)
	}

	for prop, p := range s.AllOf[1].Properties.ByIndex() {
		switch prop {
		case "code":
			codes.add(p)
		case "status":
			statuses.add(p)
		case "additional_data": // a narrower shape of what the common error already allows
		default:
			return fmt.Errorf("unexpected property %q", prop)
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

func fixOpenAPI() error {
	doc, err := openapi.LoadFromFile(pathOpenAPI)
	if err != nil {
		return err
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

	if err := consolidateErrors(doc, "error_api_", "error_api", "publicApiCommonErrorResponse"); err != nil {
		return err
	}

	// TODO: apply
	// - openapi-enrich
	// - openapi-flatten
	// - openapi-compress
	// - openapi-flatten
	// - openapi-codegen -client -debug

	if err := doc.Validate(); err != nil {
		return fmt.Errorf("validating schema: %w", err)
	}

	if err := doc.WriteToFile(pathOpenAPI); err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}

	return nil
}
