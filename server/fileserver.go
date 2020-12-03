package server

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
)

/*
To allow deep-linking into angular router child paths on go chi server
https://github.com/go-chi/chi/issues/403#issuecomment-469152247

Usage examples
public -> public url path relative to server root.
static -> static files directory relative to project root

		StaticFileServer(router, "/", "www")
		StaticFileServer(router, "/site", "static/")
		StaticFileServer(router, "/admin", "dist/")
		StaticFileServer(router, "/", "../admin2/")

*/
func StaticFileServer(r chi.Router, public string, static string) {

	// everything up to the r.Get call is executed the first time the function is called
	if strings.ContainsAny(public, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	root, _ := filepath.Abs(static)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		panic("Static Documents Directory Not Found")
	}

	fs := http.StripPrefix(public, http.FileServer(http.Dir(root)))

	if public != "/" && public[len(public)-1] != '/' {
		r.Get(public, http.RedirectHandler(public+"/", 301).ServeHTTP)
		public += "/"
	}

	log.Printf("Serving spa index.html from: %s", http.Dir(root))

	// Register the Get request for the specified path, most likely /*
	r.Get(public+"*", func(w http.ResponseWriter, r *http.Request) {
		file := strings.Replace(r.RequestURI, public, "/", 1)
		// if the requested resource was not found, pass the request to the client
		if _, err := os.Stat(root + file); os.IsNotExist(err) {
			http.ServeFile(w, r, path.Join(root, "index.html"))
			return
		}
		// if the requested resource was found, serve it
		fs.ServeHTTP(w, r)
	})
}
