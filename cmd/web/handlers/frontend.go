package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

func Frontend(fs http.Handler, indexPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("/ui/dist", r.URL.Path)

		_, err := os.Stat(path)
		if err == nil {
			fs.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, indexPath)
	}
}
