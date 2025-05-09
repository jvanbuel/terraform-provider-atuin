FROM golang:1.24.3-alpine
ENV ATUIN_HOST=""

RUN set -ex; \
    apk update; \
    apk add --no-cache git

WORKDIR /app
COPY . .

CMD CGO_ENABLED=0 TF_ACC=1 ATUIN_HOST=$ATUIN_HOST go test ./...
