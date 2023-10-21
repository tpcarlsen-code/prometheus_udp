FROM alpine:latest

RUN apk --no-cache add curl

ARG TARGETARCH
COPY ./build/$TARGETARCH/prometheus_udp /usr/local/bin/prometheus_udp

EXPOSE ${HTTP_PORT:-9231}

HEALTHCHECK --interval=10s --timeout=3s \
  CMD curl -f http://localhost:${HTTP_PORT:-9231}/-/health || exit 1
CMD ["/usr/local/bin/prometheus_udp"]
