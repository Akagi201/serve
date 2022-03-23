FROM golang:latest as builder

RUN mkdir /tmp/building

COPY . /tmp/building/

RUN go version

RUN cd /tmp/building && \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o /serve ./cmd/serve/

FROM alpine:latest

COPY --from=builder /serve /home

WORKDIR /home

CMD [ "./serve", "--addr", "0.0.0.0:8080" ]
