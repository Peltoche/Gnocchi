package assets

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssetsHTTPHandler(t *testing.T) {
	t.Run("Success without hot-reload", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/assets/site.webmanifest", nil)
		srv := chi.NewRouter()

		NewHTTPHandler(Config{HotReload: false}).Register(srv, nil)
		srv.ServeHTTP(w, r)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, http.StatusOK, res.StatusCode)

		// Check the responde bose
		webmanifestFile, err := os.ReadFile("./public/site.webmanifest")
		require.NoError(t, err)
		response, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, webmanifestFile, response)

		assert.Equal(t, "max-age=31536000", res.Header.Get("Cache-Control"))
		assert.NotEmpty(t, res.Header.Get("Cache-Control"))
	})
}
