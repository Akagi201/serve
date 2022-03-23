# serve

A dead simple file server.

## Features

- [x] Serve static files
- [x] Github CI/CD
- [x] Deploy to digitalocean kubernetes cluster
- [ ] Auto https
- [ ] Safe upload file
- [ ] Multi stage docker build

### HTTPS certs

- `openssl genrsa -out server.key 2048`
- `openssl req -new -x509 -key server.key -out server.crt -days 365`
- Convert crt to pem: `openssl x509 -in server.crt -out server.pem -outform PEM`
