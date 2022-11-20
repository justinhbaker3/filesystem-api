package api

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerStarts(t *testing.T) {
}

func TestHandler(t *testing.T) {
	h := Handler{}
	h.Initialize(&Config{
		port:              8080,
		startingDirectory: "/",
	})

	t.Run("test retrieving file", func(t *testing.T) {
		fsMock := fstest.MapFS{
			"file_one": &fstest.MapFile{
				Data: []byte("file contents for file one"),
			},
		}
		h.fs = fsMock
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/file_one", nil)
		h.router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		expectedResponse := &GetResponse{
			IsDirectory:     false,
			GetFileResponse: "file contents for file one",
		}
		resp := &GetResponse{}
		err := json.Unmarshal(w.Body.Bytes(), resp)
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, resp)
	})

	t.Run("test retrieving directory", func(t *testing.T) {
		fsMock := fstest.MapFS{
			"parent/directory_one": &fstest.MapFile{
				Mode: fs.ModeDir,
			},
			"parent/file_one": &fstest.MapFile{},
		}
		h.fs = fsMock
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/parent", nil)
		h.router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		expectedResponse := &GetResponse{
			IsDirectory: true,
			GetDirectoryResponse: []DirectoryContent{
				{
					IsDirectory: true,
					Name:        "directory_one",
					Permissions: "----------",
				},
				{
					IsDirectory: false,
					Name:        "file_one",
					Permissions: "----------",
				},
			},
		}
		resp := &GetResponse{}
		err := json.Unmarshal(w.Body.Bytes(), resp)
		require.NoError(t, err)
		assert.Equal(t, expectedResponse, resp)
	})

	t.Run("test not found", func(t *testing.T) {
		fsMock := fstest.MapFS{}
		h.fs = fsMock
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/file_one", nil)
		h.router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("test base directory", func(t *testing.T) {
		fsMock := fstest.MapFS{
			"file_one": &fstest.MapFile{
				Data: []byte("file contents for file one"),
			},
		}
		h.fs = fsMock
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		h.router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		expectedResponse := &GetResponse{
			IsDirectory: true,
			GetDirectoryResponse: []DirectoryContent{
				{
					IsDirectory: false,
					Name:        "file_one",
					Size:        26,
					Permissions: "----------",
				},
			},
		}
		resp := &GetResponse{}
		err := json.Unmarshal(w.Body.Bytes(), resp)
		require.NoError(t, err)
		assert.Equal(t, expectedResponse.IsDirectory, resp.IsDirectory)
	})
}
