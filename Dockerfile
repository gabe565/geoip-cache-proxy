#syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.6.1 AS xx

FROM --platform=$BUILDPLATFORM golang:1.25.1-alpine AS build
WORKDIR /app

COPY --from=xx / /

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETPLATFORM
RUN --mount=type=cache,target=/root/.cache \
  CGO_ENABLED=0 xx-go build -ldflags='-w -s' -trimpath


FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /app/geoip-cache-proxy /
ENTRYPOINT ["/geoip-cache-proxy"]
