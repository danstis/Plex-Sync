# Compile the executable
FROM golang:1.12 AS builder

ARG VERSION="0.0.0-dev"

WORKDIR /plex-sync/
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o plex-sync -ldflags "-X github.com/danstis/Plex-Sync/plex.Version=$VERSION" .

# Create the final container
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl gettext

WORKDIR /plex-sync/
COPY --from=builder /plex-sync/plex-sync .
COPY ./README.md .
COPY ./config/config.toml.default ./config/config.toml
COPY ./config/tvshows.txt.default ./config/tvshows.txt
COPY ./web/static/ ./web/static/
COPY ./web/templates/ ./web/templates/

EXPOSE 8123/tcp

ENTRYPOINT ["./plex-sync"]


# docker build --rm -t danstis/plex-sync . && docker run -it --rm -p 8123:8123 -v "/tmp/plexsync/config:/root/config" -v "/tmp/plexsync/logs:/root/logs" danstis/plex-sync
