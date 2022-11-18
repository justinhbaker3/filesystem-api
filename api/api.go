package api

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

type Handler struct {
	cfg    Config
	router *gin.Engine
	fs     fs.FS
}

type Config struct {
	port              int
	startingDirectory string
}

type GetResponse struct {
	FileOrDirectory      string               `json:"type"`
	GetDirectoryResponse GetDirectoryResponse `json:"contents,omitempty"`
	GetFileResponse      string               `json:"file_body,omitempty"`
}

type GetDirectoryResponse struct {
	Content []DirectoryContent `json:"content"`
}

type DirectoryContent struct {
	FileOrDirectory string `json:"type"`
	Name            string `json:"name"`
}

func (h *Handler) Initialize(cfg *Config) {
	h.router = gin.Default()
	h.fs = os.DirFS(cfg.startingDirectory)
	h.registerRoutes()
}

func (c *Config) Flags() []cli.Flag {
	flags := []cli.Flag{
		&cli.IntFlag{
			Name:        "port",
			Usage:       "port to start the filesystem api on",
			Destination: &c.port,
			EnvVars:     []string{"PORT"},
			Value:       8080,
		},
		&cli.StringFlag{
			Name:        "directory",
			Usage:       "starting directory to serve the api from",
			Destination: &c.startingDirectory,
			EnvVars:     []string{"DIRECTORY"},
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
	fs.ReadDir()
	c.String(http.StatusOK, filePath)
}

func (h *Handler) Start() error {
	h.registerRoutes()
	return h.router.Run(fmt.Sprintf(":%d", h.port))
}
