# Compile the executable
FROM golang:1.11 AS builder

WORKDIR /go/src/plex-sync
COPY . .

RUN go build -v ./...

# Create the final container
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /go/src/plex-sync/Plex-Sync /app/Plex-Sync
COPY --from=builder /go/src/plex-sync/README.md /app/Plex-Sync/README.md
COPY --from=builder /go/src/plex-sync/config/config.toml.default /app/Plex-Sync/config/config.toml
COPY --from=builder /go/src/plex-sync/config/tvshows.txt.default /app/Plex-Sync/config/tvshows.txt
COPY --from=builder /go/src/plex-sync/web/static /app/Plex-Sync/web/static
COPY --from=builder /go/src/plex-sync/web/templates /app/Plex-Sync/web/templates

CMD ["/app/Plex-Sync"]
