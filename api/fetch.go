package main

import (
	"cmp"
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/MarkRosemaker/openapi"
	"golang.org/x/sync/errgroup"
)

const path = "api/openapi.json"

var isClaudeCode = os.Getenv("CLAUDECODE") != ""

func main() {
	eg := errgroup.Group{}

	ctx := context.Background()
	eg.Go(func() error { return fetchLLMs(ctx) })
	eg.Go(func() error { return fetchOpenAPI(ctx) })

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

func fetchOpenAPI(ctx context.Context) error {
	var body io.ReadCloser
	if isClaudeCode {
		var err error

		body, err = os.Open("api/openapi-official.json")
		if err != nil {
			return err
		}
	} else {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://developers.notion.com/openapi.json", nil)
		if err != nil {
			return err
		}

		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		body = rsp.Body
	}
	defer body.Close()

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, body); err != nil {
		return err
	}

	return nil
}

func fixOpenAPI() error {
	doc, err := openapi.LoadFromFile(path)
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

	return doc.WriteToFile(path)
}
