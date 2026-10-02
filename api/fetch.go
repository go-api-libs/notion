package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	eg := errgroup.Group{}

	eg.Go(func() error { return fetchLLMs(ctx) })
	eg.Go(func() error { return fetchOpenAPI(ctx) })

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}

func fetchLLMs(ctx context.Context) error {
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://developers.notion.com/openapi.json", nil)
	if err != nil {
		return err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	f, err := os.Create("api/openapi.json")
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, rsp.Body); err != nil {
		return err
	}

	return nil
}
