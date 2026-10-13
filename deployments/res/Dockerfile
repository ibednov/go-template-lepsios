FROM golang:1.27-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
ENV GOWORK=off
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/go-template-lepsios ./cmd

FROM alpine:3.22

RUN apk add --no-cache ca-certificates postgresql-client \
    && mkdir -p /app/var/files

WORKDIR /app
COPY --from=builder /app/bin/go-template-lepsios ./go-template-lepsios
COPY --from=builder /app/migrations ./migrations
COPY deployments/res/docker-entrypoint.sh ./docker-entrypoint.sh

RUN chmod +x ./docker-entrypoint.sh

ARG API_CONTAINER_PORT=8080
EXPOSE ${API_CONTAINER_PORT}

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:${PORT:-8080}/health/live || exit 1

ENTRYPOINT ["./docker-entrypoint.sh"]
