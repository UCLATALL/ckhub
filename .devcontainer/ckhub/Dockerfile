# syntax=docker/dockerfile:1

ARG GOLANG_IMAGE=docker.io/library/golang:alpine

################################################################################
# Golang stage
################################################################################
# hadolint ignore=DL3006
FROM $GOLANG_IMAGE AS golang

SHELL ["/bin/ash", "-euxo", "pipefail", "-c"]

# hadolint ignore=DL3018
RUN apk add --no-cache --upgrade \
    bc \
    binutils-gold \
    curl \
    docker-cli \
    gcc \
    git \
    git-lfs \
    gnupg \
    make \
    musl-dev

################################################################################
# Goswagger stage
################################################################################
FROM golang AS goswagger

RUN go install -trimpath -ldflags "-extldflags -static" \
    github.com/go-swagger/go-swagger/cmd/swagger@latest

################################################################################
# Gotestsum stage
################################################################################
FROM golang AS golangci

RUN go install -trimpath -ldflags "-extldflags -static" \
    github.com/golangci/golangci-lint/cmd/golangci-lint@latest

################################################################################
# Gotestsum stage
################################################################################
FROM golang AS gotestsum

RUN go install -trimpath -ldflags "-extldflags -static" \
    gotest.tools/gotestsum@latest

################################################################################
# Tools stage
################################################################################
FROM golang AS tools

COPY --from=goswagger /go/bin /usr/local/bin/
COPY --from=golangci /go/bin /usr/local/bin/
COPY --from=gotestsum /go/bin /usr/local/bin/

COPY scripts /usr/local/bin/
RUN chmod +x /usr/local/bin/*

################################################################################
# Final stage
################################################################################
FROM golang

COPY --from=tools /usr/local/bin /usr/local/bin/

ENTRYPOINT [ "await" ]
