// This file provides a client with methods as well as functions to interact with the HTTP API.
//
// Code generated by github.com/MarkRosemaker DO NOT EDIT.

package notion

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/MarkRosemaker/jsonutil"
	"github.com/go-api-libs/api"
	"github.com/go-json-experiment/json"
	"github.com/google/uuid"
)

const (
	userAgent = "NotionGoAPILibrary/2022-06-28 (https://github.com/go-api-libs/notion)"
)

var (
	baseURL = &url.URL{
		Host:   "api.notion.com",
		Path:   "/v1",
		Scheme: "https",
	}

	jsonOpts = json.JoinOptions(
		json.RejectUnknownMembers(true),
		json.WithMarshalers(json.JoinMarshalers(
			json.MarshalToFunc(jsonutil.URLMarshal))),
		json.WithUnmarshalers(json.JoinUnmarshalers(
			json.UnmarshalFromFunc(jsonutil.URLUnmarshal))))
)

// Client conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The HTTP client to use for requests.
	cli *http.Client
	// The bearer token to use.
	Bearer string
}

// NewClient creates a new Client, setting the bearer token to [os.Getenv]("NOTION_TOKEN").
func NewClient() (*Client, error) {
	bearer := os.Getenv("NOTION_TOKEN")
	if bearer == "" {
		return nil, errors.New("bearer token \"NOTION_TOKEN\" not provided")
	}
	return &Client{
		Bearer: "Bearer " + strings.TrimPrefix(bearer, "Bearer "),
		cli:    http.DefaultClient,
	}, nil
}

// GetPage defines an operation.
//
//	GET /pages/{id}
func (c *Client) GetPage(ctx context.Context, id uuid.UUID) (*Page, error) {
	return GetPage[Page](ctx, c, id)
}

// GetPage defines an operation.
// You can define a custom result to unmarshal the response into.
//
//	GET /pages/{id}
func GetPage[R any](ctx context.Context, c *Client, id uuid.UUID) (*R, error) {
	u := baseURL.JoinPath("pages", id.String())
	req := (&http.Request{
		Header: http.Header{
			"NotionVersion": []string{"2022-06-28"},
			"User-Agent":    []string{userAgent},
		},
		Host:       u.Host,
		Method:     http.MethodGet,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		URL:        u,
	}).WithContext(ctx)

	rsp, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	switch rsp.StatusCode {
	case http.StatusOK:
		// TODO
		switch mt, _, _ := strings.Cut(rsp.Header.Get("Content-Type"), ";"); mt {
		case "application/json":
			var out R
			if err := json.UnmarshalRead(rsp.Body, &out, jsonOpts); err != nil {
				return nil, api.WrapDecodingError(rsp, err)
			}

			return &out, nil
		default:
			return nil, api.NewErrUnknownContentType(rsp)
		}
	default:
		return nil, api.NewErrUnknownStatusCode(rsp)
	}
}
