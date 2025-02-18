FROM gcr.io/distroless/static:nonroot
LABEL org.opencontainers.image.source="https://github.com/gabe565/geoip-cache-proxy"
WORKDIR /
COPY geoip-cache-proxy /
ENTRYPOINT ["/geoip-cache-proxy"]
