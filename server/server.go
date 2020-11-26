package server

import (
	"compress/flate"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	bookroutes "abidhmuhsin.com/gowebapp/mvc/routes"
	v1 "abidhmuhsin.com/gowebapp/server/api/v1"
	users "abidhmuhsin.com/gowebapp/server/crudjsonusers"
	validatedusers "abidhmuhsin.com/gowebapp/server/validatedusers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httplog"
)

// HelloWorld is a sample handler
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

// NewRouter returns a new HTTP handler that implements the main server routes
func NewRouter() http.Handler {
	router := chi.NewRouter()

	// Set up our middleware with sane defaults
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	compressor := middleware.NewCompressor(flate.DefaultCompression)
	router.Use(compressor.Handler)
	router.Use(middleware.Timeout(60 * time.Second))

	//router.Use(middleware.Logger)
	// Logger
	logger := httplog.NewLogger("httplog-example", httplog.Options{
		// JSON: true,
		Concise: true,
		// Tags: map[string]string{
		// 	"version": "v1.0-81aa4244d9fc8076a",
		// 	"env":     "dev",
		// },
	})
	router.Use(httplog.RequestLogger(logger))
	router.Use(middleware.Heartbeat("/ping"))

	// Set up our root handlers
	router.Get("/", HelloWorld)

	// Set up our API
	router.Mount("/api/v1/", v1.NewRouter())

	// Set up users API
	router.Mount("/api/users/", users.NewRouter())

	// Set up validated users API
	router.Mount("/api/users-v/", validatedusers.NewRouter())

	// Pass router to Books mvc and register /book/ - routes
	bookroutes.RegisterBookStoreRoutes(router)

	// Set up static file serving
	staticPath, _ := filepath.Abs("../../static/")
	fs := http.FileServer(http.Dir(staticPath))
	router.Handle("/*", fs)

	return router
}
