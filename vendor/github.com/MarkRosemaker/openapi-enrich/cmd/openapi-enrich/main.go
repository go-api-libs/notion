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
	"github.com/MarkRosemaker/openapi-enrich/recorder"
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

	prevIas, err := jsonutil.ReadFile[cassette.Interactions](iaPath, jsonOpts)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	tr := recorder.NewTransport(http.DefaultTransport, prevIas)

	// Call requests that don't have a response yet
	for _, ia := range prevIas {
		if ia.Response.StatusCode > 0 {
			continue
		}

		req, err := ia.Request.Create(ctx)
		if err != nil {
			return err
		}

		if _, err := tr.RoundTrip(req); err != nil {
			return err
		}
	}

	if err := enrich.Enrich(doc, tr.Interactions); err != nil {
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

	if err := jsonutil.WriteFile(iaPath, tr.Interactions, jsonOpts); err != nil {
		return err
	}

	return nil
}
