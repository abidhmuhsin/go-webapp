package server

import (
	"compress/flate"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	bookroutes "abidhmuhsin.com/gowebapp/mvc/routes"
	v1 "abidhmuhsin.com/gowebapp/server/api/v1"
	users "abidhmuhsin.com/gowebapp/server/crudjsonusers"
	authcontroller "abidhmuhsin.com/gowebapp/server/jwtauth"
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

	//router.Use(middleware.Logger)		// plain, text logger
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
	router.Get("/h", HelloWorld)

	// Set up our API
	router.Mount("/api/v1/", v1.NewRouter())

	// Set up users API
	router.Mount("/api/users/", users.NewRouter())

	// Set up validated users API
	router.Mount("/api/users-v/", validatedusers.NewRouter())

	// Pass router to Books mvc and register /book/ - routes
	bookroutes.RegisterBookStoreRoutes(router)

	// JWT - Auth based calls
	router.Mount("/api/jwt", authcontroller.NewRouter()) // ending / on mount path is optional

	// Set up static file serving
	staticPath, _ := filepath.Abs("static")
	fs := http.FileServer(http.Dir(staticPath))
	log.Printf("Serving static web files from: %s", staticPath)
	router.Handle("/*", fs) // Tested working with only "/*". Dint work with /static/

	// Set up SPA - static file serving with Deep linking sub routes for SPAs
	StaticFileServer(router, "/spa", "static/spa/")
	StaticFileServer(router, "/spa2", "static/spa2/")

	return router
}
