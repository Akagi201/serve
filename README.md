# serve

Simple http server for localhost development

## Features
- [x] Use gohttp as http framework.
- [x] Static http file server.
- [x] Support https.
- [ ] Support http2.
- [ ] Support WebSocket.
- [ ] Support browser-sync like features.

## Install
* `go get github.com/Akagi201/serve`

## Run
* `./serve -h`

### HTTPS serve
* `openssl genrsa -out server.key 2048`
* `openssl req -new -x509 -key server.key -out server.crt -days 365`
* `sudo ./serve --http=0.0.0.0:8888 --domains=akagi201.org`

### HTTP serve
* `serve --http=:8888 --https=:0`
