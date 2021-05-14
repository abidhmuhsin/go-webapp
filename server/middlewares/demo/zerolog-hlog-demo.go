package main

import (
	"errors"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	zlog "github.com/rs/zerolog/log"
)

func main() {

	// UNIX Time is faster and smaller than most timestamps
	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Initialize logger with default output to console with basic context fields
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("host", "local-pc").
		Int("id", 007).
		Logger()

	//Add additional key value to log context
	log = log.With().Str("user", "admin").Logger()

	// Log a line
	log.Log().Msg("Starting up")
	// Log an Info line with sublogger
	log.Info().Msg("Starting up with Info")
	// Log an Error line
	log.Error().Err(errors.New("unauthorized")).Msg("This is a sample error log")

	//Add file and line number to log -- line no. not useful for http middleware logging
	sublogger := log.With().Caller().Logger() // any logs logged with sublogger will have line number
	// Log a Debug line with sublogger
	sublogger.Debug().Msg("Starting up sublogger (with lineno) with Debug")

	//Set as standard logger outputs -- so any logs using standard package will also gets logged into zerolog output
	stdlog.SetFlags(stdlog.LstdFlags | stdlog.Lshortfile)
	stdlog.SetOutput(log)
	// Try one message using standard logger
	stdlog.Print("This message was printed with standard go log package")

	// Replace global zerolog/log with our custom logger with fields added
	zlog.Logger = log.With().Str("global-zlog", "logged from zerolog/log global logger").Logger()
	// import "github.com/rs/zerolog/log" anywhere in the app and will use the zerologs global logger with above extra fields
	zlog.Log().Msg("This message was printed with global logger from zerolog/log package")

	// loggedRouter := LogMiddleWare(handler) // this dint work..

	r := chi.NewRouter()

	// Install the logger handler to the router with default output on the console setup above
	r.Use(hlog.NewHandler(log)) // -- this line is mandatory since it installs hlog into the context.

	// Install some provided extra handler to set some request's context fields.
	// Thanks to that handler, all our logs will come with some prepopulated fields.
	r.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))
	r.Use(hlog.RemoteAddrHandler("ip"))
	//r.Use(hlog.UserAgentHandler("user_agent"))
	r.Use(hlog.RefererHandler("referer"))
	//r.Use(hlog.RequestIDHandler("req_id", "Request-Id"))

	// Finally add another middleware based custom handler logging with hlog
	r.Use(LoggerMiddleWare)
	// loggedRouter := LoggerMiddleWare(r) dint work

	// After all the middlewares are added - attach the routes
	r.Get("/", HelloWorld)

	// Start our HTTP server
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal().Err(err).Msg("Startup failed")
	}
}

// HelloWorld is a sample handler
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	//	time.Sleep(1 * time.Second)
	fmt.Fprintf(w, "Hello world!")
}

// NewRouter returns an HTTP handler that implements the routes for the API
func NewRouter() http.Handler {
	r := chi.NewRouter()

	//	r.Use(LogMiddleWare)

	// Register the API routes
	r.Get("/", HelloWorld)
	//	r.Get("/{name}", HelloName)

	return r
}

func LoggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the logger from the request's context. You can safely assume it
		// will be always there: if the handler is removed, hlog.FromRequest
		// will return a no-op logger.
		hlog.FromRequest(r).Info().
			Str("username", "Administrator").
			Str("status", "ok").
			Msg("Some request was performed")

		// That will do the loging, we can execute the handler
		next.ServeHTTP(w, r)
	})
}

/*
---Sample Output---
{"host":"local-pc","id":7,"user":"admin","time":"2021-02-11T20:45:26+05:30","message":"Starting up"}
{"level":"info","host":"local-pc","id":7,"user":"admin","time":"2021-02-11T20:45:26+05:30","message":"Starting up with Info"}
{"level":"error","host":"local-pc","id":7,"user":"admin","error":"unauthorized","time":"2021-02-11T20:54:40+05:30","message":"This is a sample error log"}
{"level":"debug","host":"local-pc","id":7,"user":"admin","time":"2021-02-11T20:45:26+05:30","caller":"C:/Coding/Go/AccuServe/go-webapp/server/middlewares/demo/zerolog-hlog-demo.go:40","message":"Starting up sublogger (with lineno) with Debug"}
{"host":"local-pc","id":7,"user":"admin","time":"2021-02-11T20:45:26+05:30","message":"This message was printed with standard go log package"}
// middleware logs
{"level":"info","host":"local-pc","id":7,"user":"admin","ip":"::1","username":"Administrator","status":"ok","time":"2021-02-11T20:45:33+05:30","message":"Some request was performed"}
{"level":"info","host":"local-pc","id":7,"user":"admin","ip":"::1","method":"GET","url":"/","status":200,"size":12,"duration":0,"time":"2021-02-11T20:45:33+05:30"}
*/
