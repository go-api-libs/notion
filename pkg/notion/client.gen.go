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
	"strconv"
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
	bearer string
}

// NewClient creates a new Client, setting the bearer token to [os.Getenv]("NOTION_API_KEY").
func NewClient() (*Client, error) {
	bearer := os.Getenv("NOTION_API_KEY")
	if bearer == "" {
		return nil, errors.New("bearer token NOTION_API_KEY not provided")
	}

	return &Client{
		bearer: "Bearer " + strings.TrimPrefix(bearer, "Bearer "),
		cli:    http.DefaultClient,
	}, nil
}

// Retrieves a Page object using the ID in the request path. This endpoint exposes page properties, not page content.
//
//	GET /pages/{id}
func (c *Client) GetPage(ctx context.Context, id uuid.UUID) (*Page, error) {
	return GetPage[Page](ctx, c, id)
}

// Retrieves a Page object using the ID in the request path. This endpoint exposes page properties, not page content.
// You can define a custom result to unmarshal the response into.
//
//	GET /pages/{id}
func GetPage[R any](ctx context.Context, c *Client, id uuid.UUID) (*R, error) {
	u := baseURL.JoinPath("pages", id.String())
	req := (&http.Request{
		Header: http.Header{
			"Authorization":  []string{c.bearer},
			"Notion-Version": []string{"2022-06-28"},
			"User-Agent":     []string{userAgent},
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
		// Returns the page that was requested or created.
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

// GetBlocks96245c8f178444a482ad1941127c3ec3Children defines an operation.
//
//	GET /blocks/96245c8f-1784-44a4-82ad-1941127c3ec3/children
func (c *Client) GetBlocks96245c8f178444a482ad1941127c3ec3Children(ctx context.Context, params *GetBlocks96245c8f178444a482ad1941127c3ec3ChildrenParams) (*GetBlocks96245c8f178444a482ad1941127c3ec3ChildrenOkJSONResponse, error) {
	return GetBlocks96245c8f178444a482ad1941127c3ec3Children[GetBlocks96245c8f178444a482ad1941127c3ec3ChildrenOkJSONResponse](ctx, c, params)
}

// GetBlocks96245c8f178444a482ad1941127c3ec3Children defines an operation.
// You can define a custom result to unmarshal the response into.
//
//	GET /blocks/96245c8f-1784-44a4-82ad-1941127c3ec3/children
func GetBlocks96245c8f178444a482ad1941127c3ec3Children[R any](ctx context.Context, c *Client, params *GetBlocks96245c8f178444a482ad1941127c3ec3ChildrenParams) (*R, error) {
	u := baseURL.JoinPath("/blocks/96245c8f-1784-44a4-82ad-1941127c3ec3/children")

	if params != nil && params.PageSize != 0 {
		u.RawQuery = url.Values{"page_size": []string{strconv.Itoa(params.PageSize)}}.Encode()
	}

	req := (&http.Request{
		Header: http.Header{
			"Authorization":  []string{c.bearer},
			"Notion-Version": []string{"2022-06-28"},
			"User-Agent":     []string{userAgent},
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
