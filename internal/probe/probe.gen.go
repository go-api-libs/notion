package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func main() {
	start := time.Now()
	defer func() { log.Printf("Probe took %s", time.Since(start)) }()

	r, err := record()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		// Ensure recorder is saved once we're done
		if err := r.Stop(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := probe(); err != nil {
		log.Fatal(err)
	}
}

func record() (*recorder.Recorder, error) {
	r, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       fmt.Sprintf("pkg/notion/vcr/%s", time.Now().Format(time.DateOnly)),
		Mode:               recorder.ModeReplayWithNewEpisodes,
		SkipRequestLatency: true,
		RealTransport:      http.DefaultTransport,
	})
	if err != nil {
		return nil, err
	}

	r.SetMatcher(matcher)
	r.AddHook(maskAuthorization, recorder.BeforeSaveHook)
	r.AddHook(maskSecrets, recorder.BeforeSaveHook)

	// replace the default transport with the recorder
	http.DefaultClient.Transport = r

	return r, nil
}

func matcher(r *http.Request, i cassette.Request) bool {
	u, err := url.Parse(i.URL)
	if err != nil {
		panic(err)
	}

	return r.Method == i.Method &&
		r.URL.Scheme == u.Scheme &&
		r.URL.Opaque == u.Opaque &&
		r.URL.Host == r.Host &&
		r.URL.Path == u.Path &&
		r.URL.Fragment == u.Fragment &&
		maps.EqualFunc(r.URL.Query(), u.Query(), slices.Equal) &&
		getBody(r) == i.Body
}

func getBody(r *http.Request) string {
	if r.Body == nil {
		return ""
	}

	if r.GetBody == nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(b))
		return string(b)
	}

	body, err := r.GetBody()
	if err != nil {
		panic(err)
	}
	defer body.Close()

	b, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func maskAuthorization(i *cassette.Interaction) error {
	auth := i.Request.Headers.Get("Authorization")
	if auth == "" {
		return nil
	}

	i.Request.Headers.Set("Authorization", maskAuthString(auth))

	return nil
}

func maskAuthString(s string) string {
	if parts := strings.Split(s, " "); len(parts) == 2 {
		// mask the token except the type (Bearer, Basic, etc.)
		return fmt.Sprintf("%s %s", parts[0], strings.Repeat("*", len(parts[1])))
	}

	return strings.Repeat("*", len(s))
}
