FROM golang:1.24-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-s -w" -o vs .

FROM alpine:3.22.1

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S app \
 && adduser -S -G app -H -s /sbin/nologin app

COPY --from=builder --chown=app:app /app/vs /app/vs

USER app

EXPOSE 8090

ENTRYPOINT ["/app/vs", "run"]
