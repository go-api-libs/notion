// Command openapi-codegen generates Go client, server, and type code from an OpenAPI spec.
package main

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/MarkRosemaker/jsonutil"
	"github.com/MarkRosemaker/openapi"
	enrich "github.com/MarkRosemaker/openapi-enrich"
	"github.com/MarkRosemaker/openapi-enrich/cassette"
)

var jsonOpts = json.JoinOptions(
	jsontext.Multiline(true),
	json.WithMarshalers(json.MarshalToFunc(jsonutil.HTTPHeaderMarshal)),
	json.WithUnmarshalers(json.UnmarshalFromFunc(jsonutil.HTTPHeaderUnmarshal)),
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "openapi-enrich: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	var (
		specPath string
		iaPath   string
	)

	flag.StringVar(&specPath, "spec", "api/openapi.json", "path to OpenAPI spec file")
	flag.StringVar(&iaPath, "ia", "api/interactions.json", "path to interactions file")
	flag.Parse()

	doc, err := openapi.LoadFromFile(specPath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		doc = enrich.NewDocument()
	}

	ias, err := jsonutil.ReadFile[cassette.Interactions](iaPath, jsonOpts)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		ias = cassette.Interactions{}
	}

	// Call requests that don't have a response yet
	for i, ia := range ias {
		if ia.Response.StatusCode > 0 {
			continue
		}

		req, err := ia.Request.Create(ctx)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		ias[i].Response, err = cassette.NewResponse(resp)
		if err != nil {
			return err
		}
	}

	if err := enrich.Enrich(doc, ias); err != nil {
		return err
	}

	// Sort responses and components (but not paths to keep the order)
	for _, path := range doc.Paths {
		for _, op := range path.Operations {
			op.Responses.Sort()
		}
	}
	doc.Components.SortMaps()

	if err := doc.Validate(); err != nil {
		return fmt.Errorf("produced invalid doc: %w", err)
	}

	if err := doc.WriteToFile(specPath); err != nil {
		return err
	}

	if err := jsonutil.WriteFile(iaPath, ias, jsonOpts); err != nil {
		return err
	}

	return nil
}
