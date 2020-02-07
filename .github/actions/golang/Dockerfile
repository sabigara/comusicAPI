FROM golang:1.13-alpine

ENV CGO_ENABLED 0
RUN apk add make
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
