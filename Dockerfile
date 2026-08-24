FROM golang:1.26.1-alpine
ENV ATUIN_HOST=""

RUN set -ex; \
    apk update; \
    apk add --no-cache git

WORKDIR /app
COPY . .

RUN adduser -D -u 1000 appuser && chown -R 1000:1000 /app
USER 1000
CMD ["sh", "-c", "CGO_ENABLED=0 TF_ACC=1 ATUIN_HOST=$ATUIN_HOST go test ./..."]
