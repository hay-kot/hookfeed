package handlers

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/hay-kot/httpkit/server"
)

type StaticController struct {
	fs     embed.FS
	prefix string
}

func NewStaticController(fs embed.FS, prefix string) *StaticController {
	return &StaticController{
		fs:     fs,
		prefix: prefix,
	}
}

// HandleStatic serves static files from the embedded filesystem
func (ctrl *StaticController) HandleStatic(w http.ResponseWriter, r *http.Request) error {
	// If path is prefixed with /api, we don't serve static files.
	if strings.HasPrefix(r.URL.Path, "/api") {
		return server.Error().
			Msg("not found").
			Status(http.StatusNotFound).
			Write(r.Context(), w)
	}

	// Strip the /app prefix if present
	urlPath := r.URL.Path
	urlPath = strings.TrimPrefix(urlPath, "/app")

	// If the path is empty or just "/", serve index.html
	if urlPath == "" || urlPath == "/" {
		urlPath = "/index.html"
	}

	err := ctrl.serveFile(w, r, urlPath)
	if err != nil {
		// Fallback to the index.html file for SPA routing
		err = ctrl.serveFile(w, r, "/index.html")
		if err != nil {
			return server.Error().
				Msg("file not found").
				Status(http.StatusNotFound).
				Write(r.Context(), w)
		}
	}
	return nil
}

// serveFile serves a single file from the embedded filesystem
func (ctrl *StaticController) serveFile(w http.ResponseWriter, r *http.Request, urlPath string) error {
	// Clean the path
	urlPath = path.Clean(urlPath)

	// Remove leading slash
	urlPath = strings.TrimPrefix(urlPath, "/")

	// Construct the full path
	fullPath := path.Join(ctrl.prefix, urlPath)

	// Try to open the file
	file, err := ctrl.fs.Open(fullPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	// Get file info to check if it's a directory
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// If it's a directory, return error
	if stat.IsDir() {
		return fs.ErrNotExist
	}

	// Use http.ServeContent to serve the file with proper headers
	http.ServeContent(w, r, stat.Name(), stat.ModTime(), file.(io.ReadSeeker))
	return nil
}
