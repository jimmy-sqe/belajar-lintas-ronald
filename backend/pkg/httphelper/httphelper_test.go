package httphelper

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPHelper_LogRequestResponse(t *testing.T) {
	tt := []struct {
		name     string
		ctx      context.Context
		req      *http.Request
		resp     *http.Response
		msg      string
		loggerFn Logger
	}{
		{
			name: "should log request",
			ctx:  t.Context(),
			req: &http.Request{
				Method:        http.MethodPost,
				Body:          io.NopCloser(bytes.NewBufferString("{\"otp\":\"557070\",\"identifier\":\"6285612345678\",\"channel\":\"whatsapp\",\"metadata\":null}")),
				ContentLength: 107,
			},
			resp: &http.Response{
				StatusCode:    http.StatusUnauthorized,
				Status:        "401 Unauthorized",
				Body:          io.NopCloser(bytes.NewBufferString("{\"exp\":\"token expired\"}")),
				ContentLength: 30,
			},
			msg: "SendOtpCallback ReqRes Body",
			loggerFn: func(ctx context.Context, reqBody, respBody []byte, statusCode int, status, message string) {
				assert.Equal(t, "{\"otp\":\"557070\",\"identifier\":\"6285612345678\",\"channel\":\"whatsapp\",\"metadata\":null}", string(reqBody))
				assert.Equal(t, "{\"exp\":\"token expired\"}", string(respBody))
				assert.Equal(t, 401, statusCode)
				assert.Equal(t, "401 Unauthorized", status)
				assert.Equal(t, "SendOtpCallback ReqRes Body", message)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			LogRequestResponse(tc.ctx, tc.req, tc.resp, tc.msg, tc.loggerFn)
		})
	}
}

func TestHTTPHelper_CopyHTTPResponseBody(t *testing.T) {
	tt := []struct {
		name     string
		resp     *http.Response
		expected io.ReadCloser
	}{
		{
			name: "should copy response body (respBody != nil)",
			resp: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString("{ \"errors\": [ { \"code\": \"Bad Request\", \"field\": \"file\", \"message\": \"FAILED\", \"detail\": \"File Is Empty\" } ] }")),
			},
			expected: io.NopCloser(bytes.NewBufferString("{ \"errors\": [ { \"code\": \"Bad Request\", \"field\": \"file\", \"message\": \"FAILED\", \"detail\": \"File Is Empty\" } ] }")),
		},
		{
			name: "not return error although response body is nil",
			resp: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       nil,
			},
			expected: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, CopyHTTPResponseBody(tc.resp))
		})
	}
}
