package middlewares

import (
	"net/http"
	"runtime/debug"
	"time"

	log "github.com/go-kit/kit/log"
)

/*
// https://blog.questionable.services/article/guide-logging-middleware-go/

func ExampleMiddleware(next http.Handler) http.Handler {
  // We wrap our anonymous function, and cast it to a http.HandlerFunc
  // Because our function signature matches ServeHTTP(w, r), this allows
  // our function (type) to implicitly satisify the http.Handler interface.
  return http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      // Logic before - reading request values, putting things into the
      // request context, performing authentication

      // Important that we call the 'next' handler in the chain. If we don't,
      // then request handling will stop here.
      next.ServeHTTP(w, r)
      // Logic after - useful for logging, metrics, etc.
      //
      // It's important that we don't use the ResponseWriter after we've called the
      // next handler: we may cause conflicts when trying to write the response
    }
  )
}
*/

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func LoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Log(
						"err", err,
						"trace", debug.Stack(),
					)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.Log(
				"status", wrapped.status,
				"ip", r.RemoteAddr,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		}

		return http.HandlerFunc(fn)
	}
}
