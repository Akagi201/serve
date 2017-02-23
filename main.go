package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/ossrs/go-oryx-lib/https"
)

func main() {
	_, httpPort, err := net.SplitHostPort(opts.HTTPListenAddr)
	if err != nil {
		log.Fatalf("http port parse error: %v", err)
	}
	_, httpsPort, err := net.SplitHostPort(opts.HTTPSListenAddr)
	if err != nil {
		log.Fatalf("https port parse error: %v", err)
	}

	if httpsPort != "0" && httpsPort != "443" {
		log.Fatalln("https port must be 0(disabled) or 443(enabled)")
	}

	if httpPort == "0" && httpsPort == "0" {
		log.Fatalln("http and https are both disabled")
	}

	cacheFile := opts.CacheFile
	if !path.IsAbs(opts.CacheFile) && path.IsAbs(os.Args[0]) {
		cacheFile = path.Join(path.Dir(os.Args[0]), opts.CacheFile)
	}

	html := opts.HTML
	if !path.IsAbs(opts.HTML) && path.IsAbs(os.Args[0]) {
		html = path.Join(path.Dir(os.Args[0]), opts.HTML)
	}

	fh := http.FileServer(http.Dir(html))
	http.Handle("/", fh)

	var protos []string
	if httpPort != "0" {
		protos = append(protos, fmt.Sprintf("http(:%v)", httpPort))
	}

	if httpsPort != "0" {
		var domains string
		if len(opts.HTTPSDomains) == 0 {
			domains = "all domains"
		}
		domains = strings.Join(opts.HTTPSDomains, ",")
		protos = append(protos, fmt.Sprintf("https(:%v, %v, %v)", httpsPort, domains, cacheFile))

		if opts.UseLetsEncrypt {
			protos = append(protos, "letsencrypt")
		} else {
			protos = append(protos, fmt.Sprintf("self-sign(%v, %v)", opts.SelfSignKey, opts.SelfSignCert))
		}
	}

	log.Printf("%v html root at %v", strings.Join(protos, ", "), html)

	wg := sync.WaitGroup{}
	go func() {
		defer wg.Done()

		if httpPort == "0" {
			return
		}

		if err := http.ListenAndServe(fmt.Sprintf(":%v", httpPort), nil); err != nil {
			panic(err)
		}
	}()
	wg.Add(1)

	go func() {
		defer wg.Done()

		if httpsPort == "0" {
			return
		}

		var err error
		var m https.Manager

		if opts.UseLetsEncrypt {
			if m, err = https.NewLetsencryptManager("", opts.HTTPSDomains, cacheFile); err != nil {
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
	wg.Add(1)

	wg.Wait()
}
