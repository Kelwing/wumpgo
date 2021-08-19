package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHTTPClient_Request(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string

		request *request

		wantResponse *DiscordResponse
		wantErr      string
	}{
		{},
	}

	for _, tc := range tt {
		handler := http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {

			writer.WriteHeader(200)
			_, err := writer.Write([]byte(`{"test":"data"}`))
			require.NoError(t, err)
		})
		srv := httptest.NewServer(handler)
		t.Cleanup(func() {
			srv.Close()
		})

		h := DefaultHTTPClient{
			doer:          http.DefaultClient,
			authorization: "uwu",
			userAgent:     "rawr",
		}

		resp, err := h.Request(tc.request)
		if tc.wantErr == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, tc.wantErr)
		}
		assert.Equal(t, tc.wantResponse, resp)
	}
}
