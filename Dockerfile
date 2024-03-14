FROM golang:1.21-alpine
ENV ATUIN_HOST=""

RUN set -ex; \
    apk update; \
    apk add --no-cache git

WORKDIR /app
COPY . .

CMD CGO_ENABLED=0 ATUIN_HOST=$ATUIN_HOST go test ./...
