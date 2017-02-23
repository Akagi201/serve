package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	LogLevel        string   `long:"log_level" default:"info" description:"log level"`
	HTTPListenAddr  string   `long:"http" default:"0.0.0.0:80" description:"HTTP address to listen at, :0 to disable it"`
	HTTPSListenAddr string   `long:"https" default:"0.0.0.0:443" description:"HTTPS address to listen at, :0 to disable it"`
	HTTPSDomains    []string `long:"domains" description:"the allow domains, empty to allow all."`
	HTML            string   `long:"html" default:"./" description:"the root html path to serve at"`
	CacheFile       string   `long:"cache" default:"./letsencrypt.cache" description:"the cache for https."`
	UseLetsEncrypt  bool     `long:"lets" description:"whether to use letsencrypt CA. self sign if not."`
	SelfSignKey     string   `long:"ssk" default:"server.key" description:"self-sign key, user can build it: openssl genrsa -out server.key 2048"`
	SelfSignCert    string   `long:"ssc" default:"server.crt" description:"self-sign cert, user can build it: openssl req -new -x509 -key server.key -out server.crt -days 365"`
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func init() {
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash|flags.IgnoreUnknown)

	_, err := parser.Parse()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}
}

func init() {
	if level, err := log.ParseLevel(strings.ToLower(opts.LogLevel)); err != nil {
		log.SetLevel(level)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
