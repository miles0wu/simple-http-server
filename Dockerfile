FROM golang:1.17.1-alpine3.14 AS builder
COPY . /app/
WORKDIR /app/cmd/api-server
RUN go build \
  -ldflags "-s -w -X main.BuildAt=`date +%FT%T%z`" \
  -o simple-http-server

FROM alpine:3.14
EXPOSE 8080
COPY --from=builder /app/cmd/api-server/simple-http-server /usr/local/sbin/simple-http-server
ENTRYPOINT ["simple-http-server"]
