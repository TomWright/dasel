ARG GOLANG_VERSION=1
ARG TARGET_BASE_IMAGE=debian:bookworm-slim
FROM golang:${GOLANG_VERSION} AS builder

ARG MAJOR_VERSION=v2
ARG RELEASE_VERSION=master
ARG CGO_ENABLED=0

COPY . .

RUN go build -o /dasel -ldflags="-w -s -X 'github.com/tomwright/dasel/${MAJOR_VERSION}/internal.Version=${RELEASE_VERSION}'" ./cmd/dasel

FROM ${TARGET_BASE_IMAGE}

COPY --from=builder --chmod=755 /dasel /usr/local/bin/dasel

ENTRYPOINT ["/usr/local/bin/dasel"]
CMD ["--help"]
