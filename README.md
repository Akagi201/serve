# serve

Simple http server for localhost development

## Features
- [x] Use gohttp as http framework.
- [x] Static http file server.
- [x] Support gzip.
- [ ] Support http2.
- [ ] Support http proxy.
- [ ] Support https proxy to http.
- [ ] Support websocket.
- [ ] Support browser-sync like features.

## Build
* docker: `docker build -t serve .`
* `go build main.go -o serve`

## Run
* `--host`: default host is `0.0.0.0`
* `--port`: default port is `3000`
* `--gzip`: default is `false`
