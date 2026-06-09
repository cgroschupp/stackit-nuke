FROM alpine:3.24 as base
RUN apk add --no-cache ca-certificates
RUN adduser -D stackit-nuke

FROM golang:1.26 AS build
COPY / /src
WORKDIR /src
ENV CGO_ENABLED=0
RUN \
  --mount=type=cache,target=/go/pkg \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -ldflags '-s -w -extldflags="-static"' -o bin/stackit-nuke main.go

FROM base AS goreleaser
ENTRYPOINT ["/usr/local/bin/stackit-nuke"]
COPY stackit-nuke /usr/local/bin/stackit-nuke
USER stackit-nuke

FROM base
ENTRYPOINT ["/usr/local/bin/stackit-nuke"]
COPY --from=build --chmod=755 /src/bin/stackit-nuke /usr/local/bin/stackit-nuke
RUN chmod +x /usr/local/bin/stackit-nuke
USER stackit-nuke