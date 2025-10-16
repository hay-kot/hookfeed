// Package httpclient provides a client for working with HTTP requests, using
// a method based API.
package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
)

type ClientMiddleware = func(*http.Request) (*http.Request, error)

type Client struct {
	client *http.Client
	base   string             // Base URL
	mw     []ClientMiddleware // Global Middleware

	mu     sync.RWMutex
	bearer string
}

func New(client *http.Client, base string) *Client {
	c := &Client{
		client: client,
		base:   strings.TrimRight(base, "/"),
		mw:     nil,
	}

	c.Use(func(req *http.Request) (*http.Request, error) {
		c.mu.RLock()
		defer c.mu.RUnlock()
		if c.bearer != "" {
			req.Header.Set("Authorization", "Bearer "+c.bearer)
		}

		return req, nil
	})

	return c
}

func (c *Client) Use(mws ...ClientMiddleware) {
	c.mw = append(c.mw, mws...)
}

func (c *Client) SetBearer(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bearer = token
}

func (c *Client) Get(url string, mw ...ClientMiddleware) (*http.Response, error) {
	url = c.path(url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) GetCtx(ctx context.Context, url string, mw ...ClientMiddleware) (*http.Response, error) {
	url = c.path(url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) Post(url string, payload io.Reader, mw ...ClientMiddleware) (*http.Response, error) {
	url = c.path(url)
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) PostCtx(ctx context.Context, url string, payload io.Reader, mw ...ClientMiddleware) (*http.Response, error) {
	url = c.path(url)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) Put(url string, payload io.Reader, mw ...ClientMiddleware) (*http.Response, error) {
	url = c.path(url)
	req, err := http.NewRequest(http.MethodPut, url, payload)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) PutCtx(ctx context.Context, url string, payload io.Reader, mw ...ClientMiddleware) (*http.Response, error) {
	url = c.path(url)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, payload)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) Patch(url string, payload io.Reader, mw ...ClientMiddleware) (*http.Response, error) {
	url = c.path(url)
	req, err := http.NewRequest(http.MethodPatch, url, payload)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) PatchCtx(ctx context.Context, url string, payload io.Reader, mw ...ClientMiddleware) (*http.Response, error) {
	url = c.path(url)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, payload)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) Delete(url string, mw ...ClientMiddleware) (*http.Response, error) {
	url = c.path(url)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) DeleteCtx(ctx context.Context, url string, mw ...ClientMiddleware) (*http.Response, error) {
	url = c.path(url)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

// Do will execute the provided request, applying any middleware that has been
// registered with the client. Any middleware provided to this function will
// be applied after the middleware registered with the client.
//
// Example:
//
//	Request -> ClientMiddleware(Request) -> FunctionMiddleware(ClientMiddleware(Request))
func (c *Client) Do(req *http.Request, rmw []ClientMiddleware) (*http.Response, error) {
	for _, mw := range c.mw {
		var err error
		req, err = mw(req)
		if err != nil {
			return nil, err
		}
	}

	for _, mw := range rmw {
		var err error
		req, err = mw(req)
		if err != nil {
			return nil, err
		}
	}

	return c.client.Do(req)
}

// path will safely join the base URL and the provided path and return a string
// that can be used in a request.
func (c *Client) path(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	}

	if url == "" {
		return c.base
	}

	return c.base + "/" + strings.TrimLeft(url, "/")
}

// DecodeJSON will decode the response body into the provided value.
func DecodeJSON[T any](r *http.Response) (T, error) {
	var zero T

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&zero); err != nil {
		return zero, err
	}
	return zero, nil
}

// Body will marshal the provided value into JSON and return a bytes.Reader
// that can be used as the body of a request. If the value cannot be
// marshaled, it will panic.
func Body(t any) *bytes.Reader {
	jsonData, err := json.Marshal(t)
	if err != nil {
		panic("failed to marshal body: " + err.Error())
	}

	return bytes.NewReader(jsonData)
}

func MwJSON(req *http.Request) (*http.Request, error) {
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}
