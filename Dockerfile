FROM golang:1.22-alpine AS builder

RUN apk add --no-cache ca-certificates git
WORKDIR /src

# Better layer caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a static-ish binary
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/matcha .

FROM alpine:3.20

# HTTPS calls (RSS feeds, OpenAI/LocalAI, etc.) need CA certs
RUN apk add --no-cache ca-certificates tzdata \
  && addgroup -S matcha && adduser -S matcha -G matcha

WORKDIR /app
COPY --from=builder /out/matcha /usr/local/bin/matcha

# Conventional mount points
RUN mkdir -p /config /data && chown -R matcha:matcha /config /data
USER matcha

# Default: run once using a mounted config
ENTRYPOINT ["matcha"]
CMD ["-c", "/config/config.yaml"]
