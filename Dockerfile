FROM golang:1.24.6-alpine
ENV ATUIN_HOST=""

RUN set -ex; \
    apk update; \
    apk add --no-cache git

WORKDIR /app
COPY . .

CMD ["sh", "-c", "CGO_ENABLED=0 TF_ACC=1 ATUIN_HOST=$ATUIN_HOST go test ./..."]
