package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hay-kot/hookfeed/backend/pkgs/httpclient"
	"github.com/hay-kot/hookfeed/backend/pkgs/utils"
	"github.com/hay-kot/httpkit/server"
)

func processResp(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read the entire response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		resp.Body.Close() //nolint:errcheck

		// Replace the body with a new reader so it can be read again
		resp.Body = io.NopCloser(bytes.NewReader(body))

		if bytes.HasPrefix(body, []byte("{")) {
			var t server.ErrorResp

			err := json.Unmarshal(body, &t)
			if err != nil {
				return fmt.Errorf("failed to unmarshal error response: %w", err)
			}

			utils.Dump(t)

			// TODO embed error resp into specific type
		}

		// TODO: handle specific error codes and marshal into
		// a more specific error type
		return errors.New("unexpected status code: " + resp.Status)
	}
	return nil
}

func post[TResp any](c *httpclient.Client, ctx context.Context, path string, data any) (TResp, error) {
	resp, err := c.PostCtx(ctx, path, httpclient.Body(data))
	if err != nil {
		var zero TResp
		return zero, err
	}
	defer resp.Body.Close() //nolint:errcheck

	err = processResp(resp)
	if err != nil {
		var zero TResp
		return zero, err
	}

	return httpclient.DecodeJSON[TResp](resp)
}

func patch[TResp any](c *httpclient.Client, ctx context.Context, path string, data any) (TResp, error) {
	resp, err := c.PatchCtx(ctx, path, httpclient.Body(data))
	if err != nil {
		var zero TResp
		return zero, err
	}
	defer resp.Body.Close() //nolint:errcheck

	err = processResp(resp)
	if err != nil {
		var zero TResp
		return zero, err
	}

	return httpclient.DecodeJSON[TResp](resp)
}

func put[TResp any](c *httpclient.Client, ctx context.Context, path string, data any) (TResp, error) {
	resp, err := c.PutCtx(ctx, path, httpclient.Body(data))
	if err != nil {
		var zero TResp
		return zero, err
	}
	defer resp.Body.Close() //nolint:errcheck

	err = processResp(resp)
	if err != nil {
		var zero TResp
		return zero, err
	}

	return httpclient.DecodeJSON[TResp](resp)
}

func get[TResp any](c *httpclient.Client, ctx context.Context, path string) (TResp, error) {
	resp, err := c.GetCtx(ctx, path)
	if err != nil {
		var zero TResp
		return zero, err
	}
	defer resp.Body.Close() //nolint:errcheck

	err = processResp(resp)
	if err != nil {
		var zero TResp
		return zero, err
	}

	return httpclient.DecodeJSON[TResp](resp)
}

func delete(c *httpclient.Client, ctx context.Context, path string) error {
	resp, err := c.DeleteCtx(ctx, path)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	return processResp(resp)
}
