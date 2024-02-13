ARG GOLANG_VERSION=1
ARG TARGET_BASE_IMAGE=debian:slim
FROM golang:${GOLANG_VERSION} AS builder

ARG RELEASE_VERSION
ARG CGO_ENABLED=0

COPY . .

RUN go build -o /dasel -ldflags="-X 'github.com/tomwright/dasel/v2/internal.Version=${RELEASE_VERSION}'" ./cmd/dasel

FROM ${TARGET_BASE_IMAGE}

COPY --from=builder --chmod=777 /dasel /usr/local/bin/dasel

ENTRYPOINT ["/usr/local/bin/dasel"]
CMD []
