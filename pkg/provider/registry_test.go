package provider_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openbao/openbao-plugin-secrets-oauthapp/v3/pkg/oauth2ext/clientctx"
	"github.com/openbao/openbao-plugin-secrets-oauthapp/v3/pkg/provider"
	"github.com/openbao/openbao-plugin-secrets-oauthapp/v3/pkg/testutil"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestRegistryAddsDefaultUserAgent(t *testing.T) {
	var received string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	r := provider.NewRegistry()
	r.MustRegister("http-check", func(ctx context.Context, vsn int, opts map[string]string) (provider.Provider, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, http.NoBody)
		require.NoError(t, err)
		resp, err := oauth2.NewClient(ctx, nil).Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		return testutil.MockFactory()(ctx, vsn, opts)
	})

	_, err := r.New(t.Context(), "http-check", nil)
	require.NoError(t, err)
	require.Equal(t, clientctx.DefaultUserAgent, received)
}
