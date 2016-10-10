package main

import (
	"log"
	"os"
	"strings"

	"github.com/akfork/app"
	"github.com/daaku/go.httpgzip"
	"github.com/gohttp/logger"
	"github.com/gohttp/serve"
	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Service string `long:"service" default:":3000" description:"ip:port to bind to"`
	Gzip    bool   `long:"gzip" description:"whether to enable gzip encoding or not"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			log.Printf("error: %v\n", err.Error())
			os.Exit(1)
		} else {
			log.Printf("%v\n", err.Error())
			os.Exit(0)
		}
	}

	a := app.New()
	a.Use(logger.New())
	a.Use(serve.New("./"))

	if opts.Gzip {
		a.Use(httpgzip.NewHandler)
	}

	log.Printf("HTTP listening at: %v", opts.Service)
	a.Listen(opts.Service)
}
