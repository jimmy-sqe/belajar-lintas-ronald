package httphelper

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type Logger func(ctx context.Context, reqBody, respBody []byte, statusCode int, status, message string)

func LogRequestResponse(ctx context.Context, req *http.Request, resp *http.Response, msg string, loggerFn Logger) {
	var reqBody, respBody []byte
	var status string
	var statusCode int

	if req != nil && req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
	}
	if resp != nil {
		statusCode = resp.StatusCode
		status = resp.Status

		if resp.Body != nil {
			respBody, _ = io.ReadAll(resp.Body)
		}
	}

	loggerFn(ctx, reqBody, respBody, statusCode, status, msg)
}

func CopyHTTPResponseBody(resp *http.Response) io.ReadCloser {
	if resp.Body == nil {
		return nil
	}

	var respBody []byte
	respBody, _ = io.ReadAll(resp.Body)
	defer func() {
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	}()

	return io.NopCloser(bytes.NewBuffer(respBody))
}
