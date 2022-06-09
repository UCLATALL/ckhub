# syntax=docker/dockerfile:1

ARG ALPINE_IMAGE=docker.io/library/alpine:latest

################################################################################
# Playground stage
################################################################################
# hadolint ignore=DL3006
FROM $ALPINE_IMAGE AS ckhub

ARG TARGETOS
ARG TARGETARCH

COPY bin/ckhub-${TARGETOS}-${TARGETARCH} /usr/local/bin/ckhub

################################################################################
# Final stage
################################################################################
# hadolint ignore=DL3006
FROM $ALPINE_IMAGE

# hadolint ignore=DL3018
RUN set -eux; \
    apk add --no-cache --upgrade ca-certificates; \
    addgroup -g 1000 ckhub; \
    adduser -HDu 1000 -G ckhub ckhub

COPY --from=ckhub /usr/local/bin /usr/local/bin

USER 1000:1000

ENTRYPOINT [ "ckhub" ]
