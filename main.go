package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Akagi201/light"
	"github.com/daaku/go.httpgzip"
	"github.com/gohttp/logger"
	"github.com/gohttp/serve"
	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Host string `long:"host" default:"0.0.0.0" description:"ip to bind to"`
	Port uint16 `long:"port" default:"3000" description:"port to bind to"`
	Gzip bool   `long:"gzip" description:"whether to enable gzip encoding or not"`
	Path string `long:"path" default:"./" description:"the root path to serve at"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			log.Printf("error: %v\n", err.Error())
			os.Exit(1)
		} else {
			// log.Printf("%v\n", err.Error())
			os.Exit(0)
		}
	}

	app := light.New()
	app.Use(logger.New())
	app.Use(serve.New(opts.Path))

	if opts.Gzip {
		app.Use(httpgzip.NewHandler)
	}

	log.Printf("HTTP listening at: %v:%v", opts.Host, opts.Port)
	app.Listen(fmt.Sprintf("%v:%d", opts.Host, opts.Port))
}
