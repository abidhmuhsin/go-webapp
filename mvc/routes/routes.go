package routes

import (
	"abidhmuhsin.com/gowebapp/mvc/controllers"
	"github.com/go-chi/chi"
)

var RegisterBookStoreRoutes = func(router *chi.Mux) {

	router.Route("/book", func(r chi.Router) {
		// add sub routes
		r.Get("/", controllers.GetBook)
		r.Post("/", controllers.CreateBook)
		r.Get("/{bookId}", controllers.GetBookById)
		r.Put("/{bookId}", controllers.UpdateBook)
		r.Delete("/{bookId}", controllers.DeleteBook)
	})

}
