package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	h := NewHandler()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo/bar", nil)
	h.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/foo/bar", w.Body.String())
}
