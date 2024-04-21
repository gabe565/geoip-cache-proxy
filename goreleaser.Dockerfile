FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY geoip-cache-proxy /
ENTRYPOINT ["/geoip-cache-proxy"]
