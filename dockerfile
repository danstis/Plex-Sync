# Compile the executable
FROM golang:1.11 AS builder

WORKDIR /go/src/plex-sync
COPY . .

ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux go build -a -o plex-sync .

# Create the final container
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl

WORKDIR /root/
COPY --from=builder /go/src/plex-sync/plex-sync .
COPY --from=builder /go/src/plex-sync/README.md .
COPY --from=builder /go/src/plex-sync/config/config.toml.default ./config/config.toml
COPY --from=builder /go/src/plex-sync/config/tvshows.txt.default ./config/tvshows.txt
COPY --from=builder /go/src/plex-sync/web/static/ ./web/static/
COPY --from=builder /go/src/plex-sync/web/templates/ ./templates/

EXPOSE 8123/tcp

CMD ["./plex-sync"]
