package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	plugin.ServeHTTP(nil, w, r)

	result := w.Result()
	assert.Equal(result.StatusCode, http.StatusNotFound)
}

func TestPreviewEndpoint(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{}

	t.Run("No X-Mattermost-ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/preview/", nil)

		plugin.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.Equal(http.StatusNotFound, result.StatusCode)
	})

	t.Run("No fileID", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/preview/", nil)
		r.Header.Add("Mattermost-User-ID", "test")

		plugin.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.Equal(http.StatusBadRequest, result.StatusCode)
	})
}
