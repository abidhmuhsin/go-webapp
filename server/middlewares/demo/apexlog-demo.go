package main

import (
	"errors"
	"os"
	"time"

	"github.com/apex/log"
	// "github.com/apex/log/handlers/json"
	// "github.com/apex/log/handlers/multi"
	"github.com/apex/log/handlers/text"
)

func main() {
	log.SetHandler(text.New(os.Stderr))
	// log.SetHandler(cli.Default)
	// log.SetHandler(multi.New(
	// 	text.New(os.Stderr),
	// 	json.New(os.Stderr),
	// ))
	log.SetLevel(log.InfoLevel)

	ctx := log.WithFields(log.Fields{
		"file": "something.png",
		"type": "image/png",
		"user": "tobi",
	})

	for range time.Tick(time.Second * 5) {
		ctx.Info("upload")
		ctx.Info("upload complete")
		ctx.Warn("upload retry")
		ctx.WithError(errors.New("unauthorized")).Error("upload failed")
		ctx.Errorf("failed to upload %s", "img.png")
	}
}

/*
APEX Logger
Package log implements a simple structured logging API inspired by Logrus, designed with centralization in mind.

    WithField(key, value) for single fields
    WithFields(fields) for multiple fields
    Pluggable log handlers
    Package singleton for convenience
	Simplified handlers, no formatter/hook distinction
    log.Interface shares the Logger and Entry method set

	Built-in handlers (text, json, cli, es, discard, logfmt, memory, kinesis, multi)
    Tracing support
	Centralization - centralized logging

 -- Sample o/p
multi json+text

←[34m  INFO←[0m[0115] upload                    ←[34mfile←[0m=something.png ←[34mtype←[0m=image/png ←[34muser←[0m=tobi
{"fields":{"file":"something.png","type":"image/png","user":"tobi"},"level":"info","timestamp":"2021-02-11T21:03:14.3500871+05:30","message":"upload"}
←[34m  INFO←[0m[0115] upload complete           ←[34mfile←[0m=something.png ←[34mtype←[0m=image/png ←[34muser←[0m=tobi
{"fields":{"file":"something.png","type":"image/png","user":"tobi"},"level":"info","timestamp":"2021-02-11T21:03:14.3520919+05:30","message":"upload complete"}
←[33m  WARN←[0m[0115] upload retry              ←[33mfile←[0m=something.png ←[33mtype←[0m=image/png ←[33muser←[0m=tobi
{"fields":{"file":"something.png","type":"image/png","user":"tobi"},"level":"warn","timestamp":"2021-02-11T21:03:14.3560929+05:30","message":"upload retry"}
←[31m ERROR←[0m[0115] upload failed             ←[31merror←[0m=unauthorized ←[31mfile←[0m=something.png ←[31mtype←[0m=image/png ←[31muser←[0m=tobi
{"fields":{"error":"unauthorized","file":"something.png","type":"image/png","user":"tobi"},"level":"error","timestamp":"2021-02-11T21:03:14.3570961+05:30","message":"upload failed"}
←[31m ERROR←[0m[0115] failed to upload img.png  ←[31mfile←[0m=something.png ←[31mtype←[0m=image/png ←[31muser←[0m=tobi
{"fields":{"file":"something.png","type":"image/png","user":"tobi"},"level":"error","timestamp":"2021-02-11T21:03:14.3620975+05:30","message":"failed
to upload img.png"}

*/
