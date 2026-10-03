package main

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/MarkRosemaker/openapi"
	edit "github.com/MarkRosemaker/openapi-edit"
	merge "github.com/MarkRosemaker/openapi-merge"
	"golang.org/x/sync/errgroup"
)

const pathOpenAPI = "api/openapi.json"

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

	// TODO(Claude): use merge and edit to transform the spec
	_ = merge.Schema
	// if err := = merge.Schema(doc.Components.Schemas["foo"], doc.Components.Schemas["bar"], false); err != nil {
	// return fmt.Errorf("merging bar and foo: %w", err)
	// }

	if err := edit.RedirectSchemas(doc, map[string]string{}); err != nil {
		return fmt.Errorf("redirecting schemas: %w", err)
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
