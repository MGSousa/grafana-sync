# Build App
FROM golang:1.20.0-alpine3.17 AS builder

WORKDIR ${GOPATH}/src/github.com/mpostument/grafana-sync

COPY . ./

RUN go build -o /go/bin/grafana-sync .


# Create small image with binary
FROM alpine:3.17

RUN apk --no-cache add ca-certificates

COPY --from=builder /go/bin/grafana-sync /usr/bin/grafana-sync

ENTRYPOINT ["/usr/bin/grafana-sync"]
