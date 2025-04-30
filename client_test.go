package chartmetric_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/musicx-fm/chartmetric-go-client"
	"github.com/musicx-fm/chartmetric-go-client/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_Client_GetAny(t *testing.T) {
	ts := chartmetricTestServer(map[string]http.HandlerFunc{
		"POST /token": func(w http.ResponseWriter, r *http.Request) {
			assert.NotNil(t, r.Body)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(testdata.TokenResponse))
		},
		"GET /genres": func(w http.ResponseWriter, r *http.Request) {
			assert.True(t, strings.HasPrefix(r.Header.Get("Authorization"), "Bearer "))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(testdata.GenresResponse))
		},
	})
	defer ts.Close()

	client := chartmetric.NewClient("test-refresh-token", chartmetric.WithBaseURL(ts.URL))

	responseData, err := client.GetAny(context.Background(), "/genres", nil)
	assert.NoError(t, err)
	assert.NotNil(t, responseData)
	assert.Equal(t, testdata.GenresResponse, string(responseData))
}

func chartmetricTestServer(handlers map[string]http.HandlerFunc) *httptest.Server {
	mux := http.NewServeMux()
	for path, handler := range handlers {
		mux.HandleFunc(path, handler)
	}

	return httptest.NewServer(mux)

}
