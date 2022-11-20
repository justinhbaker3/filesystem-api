package api

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

// Handler is used to instantiate a new instance of an API.
// Users should create an empty Handler (h := Handler{}), then call
// h.Initialize() and h.Start()
type Handler struct {
	cfg    *Config
	router *gin.Engine
	fs     fs.FS
}

// Config holds configuration values used to create a new Handler
type Config struct {
	port              int
	startingDirectory string
}

// GetResponse is used to build JSON responses from the API
type GetResponse struct {
	IsDirectory          bool               `json:"is_directory"`
	GetDirectoryResponse []DirectoryContent `json:"directory_content,omitempty"`
	GetFileResponse      string             `json:"file_body,omitempty"`
}

// DirectoryContent holds information about a directory
type DirectoryContent struct {
	IsDirectory bool   `json:"is_directory"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Permissions string `json:"permissions"`
}

// Initialize must be called before Start()
func (h *Handler) Initialize(cfg *Config) {
	h.router = gin.Default()
	h.fs = os.DirFS(cfg.startingDirectory)
	h.cfg = cfg
	h.registerRoutes()
}

// Start starts the API, listening on the configured port
func (h *Handler) Start() error {
	return h.router.Run(fmt.Sprintf(":%d", h.cfg.port))
}

// Flags exposes configuration values for the API as environment variables
func (c *Config) Flags() []cli.Flag {
	flags := []cli.Flag{
		&cli.IntFlag{
			Name:        "port",
			Usage:       "port to start the filesystem api on",
			Destination: &c.port,
			EnvVars:     []string{"FILESYSTEM_API_PORT"},
			Value:       8080,
		},
		&cli.StringFlag{
			Name:        "directory",
			Usage:       "starting directory to serve the api from",
			Destination: &c.startingDirectory,
			EnvVars:     []string{"FILESYSTEM_API_DIRECTORY"},
			Value:       "/",
		},
	}
	return flags
}

func (h *Handler) registerRoutes() {
	h.router.GET("/*filepath", h.get)
}

func (h *Handler) get(c *gin.Context) {
	filePath := c.Param("filepath")
	log.Print(filePath)
	// remove the leading slash unless we are at the base path
	if filePath != "/" {
		filePath = filePath[1:]
	} else {
		// if we are at the base path, use relative path syntax
		filePath = "."
	}
	fileInfo, err := fs.Stat(h.fs, filePath)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if fileInfo.IsDir() {
		dirContnet, err := h.getDirectoryContent(filePath)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		resp := &GetResponse{
			IsDirectory:          true,
			GetDirectoryResponse: dirContnet,
		}
		c.JSON(http.StatusOK, resp)
		return
	} else {
		fileContent, err := h.getFileContent(filePath)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		resp := &GetResponse{
			IsDirectory:     false,
			GetFileResponse: string(fileContent),
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

func (h *Handler) getFileContent(filePath string) (string, error) {
	fileConents, err := fs.ReadFile(h.fs, filePath)
	if err != nil {
		return "", err
	}
	return string(fileConents), nil
}

func (h *Handler) getDirectoryContent(filePath string) ([]DirectoryContent, error) {
	dirEntries, err := fs.ReadDir(h.fs, filePath)
	if err != nil {
		return nil, err
	}
	dirContent := []DirectoryContent{}
	for _, entry := range dirEntries {
		content := DirectoryContent{
			IsDirectory: entry.IsDir(),
			Name:        entry.Name(),
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		content.Size = info.Size()
		content.Permissions = info.Mode().Perm().String()

		dirContent = append(dirContent, content)
	}
	return dirContent, nil
}
