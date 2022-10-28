package rest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultHTTPClient_Request(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string

		request *request

		wantMethod string
		wantBody   []byte
		wantHeader http.Header

		wantResponse *DiscordResponse
		wantErr      string
	}{
		{
			name: "standard request",

			wantBody: []byte(`test request body`),
			wantHeader: http.Header{
				"Accept-Encoding": {"gzip"},
				"Authorization":   {"Bot not.a.token"},
				"Content-Length":  {"17"},
				"Content-Type":    {"application/json"},
				"User-Agent":      {"NCSA Mosaic/3.0 (Windows 95)"},
			},
			request: &request{
				method:      http.MethodGet,
				body:        []byte(`test request body`),
				contentType: JsonContentType,
				ctx:         context.Background(),
			},
			wantResponse: &DiscordResponse{
				StatusCode: 200,
				Body:       []byte(`{"test":"data"}`),
				Header: http.Header{
					"Content-Length": {"15"},
					"Content-Type":   {"text/plain; charset=utf-8"},
					"Date":           {"fixed"},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
				gotBody, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				assert.Equal(t, tc.wantBody, gotBody)
				assert.Equal(t, tc.wantHeader, r.Header)

				writer.Header().Set("Date", "fixed")
				writer.WriteHeader(200)
				_, err = writer.Write([]byte(`{"test":"data"}`))
				require.NoError(t, err)
			})
			srv := httptest.NewServer(handler)
			t.Cleanup(func() {
				srv.Close()
			})

			tc.request.path = srv.URL

			h := DefaultHTTPClient{
				doer:          http.DefaultClient,
				authorization: "Bot not.a.token",
				userAgent:     "NCSA Mosaic/3.0 (Windows 95)",
			}

			resp, err := h.Request(tc.request)
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.wantResponse, resp)
		})

	}
}
