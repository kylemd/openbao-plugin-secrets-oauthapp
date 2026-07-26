package clientctx_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/openbao/openbao-plugin-secrets-oauthapp/v3/pkg/oauth2ext/clientctx"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestWithUpdatedRequest(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, err := io.Copy(w, r.Body)
		require.NoError(t, err)
	}
	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	ctx := context.Background()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, srv.Client())
	ctx = clientctx.WithUpdatedRequestBody(ctx, func(body []byte) ([]byte, error) {
		return bytes.Join([][]byte{body, []byte("amended")}, []byte("+")), nil
	})

	client := oauth2.NewClient(ctx, nil)
	resp, err := client.Post(srv.URL, "text/plain", strings.NewReader("original"))
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, []byte("original+amended"), body)
}

func TestWithDefaultUserAgent(t *testing.T) {
	tests := []struct {
		name     string
		supplied string
		expected string
	}{
		{
			name:     "adds default",
			expected: clientctx.DefaultUserAgent,
		},
		{
			name:     "preserves supplied value",
			supplied: "Provider-specific client",
			expected: "Provider-specific client",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var received string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				received = r.Header.Get("User-Agent")
				w.WriteHeader(http.StatusNoContent)
			}))
			defer srv.Close()

			ctx := context.Background()
			ctx = context.WithValue(ctx, oauth2.HTTPClient, srv.Client())
			ctx = clientctx.WithDefaultUserAgent(ctx)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, http.NoBody)
			require.NoError(t, err)
			if tt.supplied != "" {
				req.Header.Set("User-Agent", tt.supplied)
			}

			resp, err := oauth2.NewClient(ctx, nil).Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, tt.expected, received)
		})
	}
}
