# serve

Simple http server for localhost development

## Features
- [x] Use gohttp as http framework.
- [x] Static http file server.
- [ ] Support http2.
- [ ] Support gzip.
- [ ] Support http proxy.
- [ ] Support https proxy to http.
- [ ] Support websocket.
- [ ] Support browser-sync like features.

## Build
* docker: `docker build -t serve .`
* `go build main.go -o serve`

## Run
* `--service`: default service is `:3000`
