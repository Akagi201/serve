package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	flags "github.com/jessevdk/go-flags"
	"github.com/ossrs/go-oryx-lib/https"
)

var opts struct {
	HTTPListenAddr  string   `long:"http" default:"0.0.0.0:80" description:"HTTP address to listen at"`
	HTTPSListenAddr string   `long:"https" default:"0.0.0.0:443" description:"HTTPS address to listen at"`
	HTTPSDomains    []string `long:"domains" description:"the allow domains, empty to allow all."`
	Path            string   `long:"path" default:"./" description:"the root path to serve at"`
	CacheFile       string   `long:"cache" default:"./letsencrypt.cache" description:"the cache for https."`
	UseLetsEncrypt  bool     `long:"lets" description:"whether to use letsencrypt CA. self sign if not."`
	SelfSignKey     string   `long:"ssk" default:"server.key" description:"self-sign key, user can build it: openssl genrsa -out server.key 2048"`
	SelfSignCert    string   `long:"ssc" default:"server.crt" description:"self-sign cert, user can build it: openssl req -new -x509 -key server.key -out server.crt -days 365"`
}

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Level = logrus.InfoLevel
	f := new(logrus.TextFormatter)
	f.TimestampFormat = "2006-01-02 15:04:05"
	f.FullTimestamp = true
	log.Formatter = f
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			log.Fatalf("cli parse error: %v", err)
		} else {
			return
		}
	}

	fh := http.FileServer(http.Dir(opts.Path))
	http.Handle("/", fh)

	var protos []string
	_, httpPort, err := net.SplitHostPort(opts.HTTPListenAddr)
	if err != nil {
		log.Fatalf("http port parse error: %v", err)
	}
	protos = append(protos, fmt.Sprintf("http(:%v)", httpPort))

	_, httpsPort, err := net.SplitHostPort(opts.HTTPSListenAddr)
	if err != nil {
		log.Fatalf("https port parse error: %v", err)
	}

	var domains string
	if len(opts.HTTPSDomains) == 0 {
		domains = "all domains"
	}
	domains = strings.Join(opts.HTTPSDomains, ",")
	protos = append(protos, fmt.Sprintf("https(:%v, %v, %v)", httpsPort, domains, opts.CacheFile))

	if opts.UseLetsEncrypt {
		protos = append(protos, "letsencrypt")
	} else {
		protos = append(protos, fmt.Sprintf("self-sign(%v, %v)", opts.SelfSignKey, opts.SelfSignCert))
	}
	log.Printf("%v html root at %v", strings.Join(protos, ", "), string(opts.Path))

	wg := sync.WaitGroup{}
	go func() {
		defer wg.Done()

		if err := http.ListenAndServe(fmt.Sprintf(":%v", httpPort), nil); err != nil {
			panic(err)
		}
	}()
	wg.Add(1)

	go func() {
		defer wg.Done()

		var err error
		var m https.Manager

		if opts.UseLetsEncrypt {
			if m, err = https.NewLetsencryptManager("", opts.HTTPSDomains, opts.CacheFile); err != nil {
				panic(err)
			}
		} else {
			if m, err = https.NewSelfSignManager(opts.SelfSignCert, opts.SelfSignKey); err != nil {
				panic(err)
			}
		}

		svr := &http.Server{
			Addr: fmt.Sprintf(":%v", httpsPort),
			TLSConfig: &tls.Config{
				GetCertificate: m.GetCertificate,
			},
		}

		if err := svr.ListenAndServeTLS("", ""); err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}
