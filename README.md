# geoip-cache-proxy

A GeoIP database caching proxy.

Requests to `updates.maxmind.com` and `download.maxmind.com` can be passed through this proxy, and the response contents will be saved for 24 hours. Mainly intended to be used in Kubernetes. If an app downloads the database and crash-loops, this app will help avoid GeoIP rate limits.

## Installation

A Docker image is provided at [`ghcr.io/gabe565/geoip-cache-proxy`](https://ghcr.io/gabe565/geoip-cache-proxy).

## Usage

Redis is required, and can be configured with `GEOIP_REDIS_ADDR` and `GEOIP_REDIS_PASSWORD`.

When run, geoip-cache-proxy will start different servers for each MaxMind endpoint:
- `localhost:8080` will proxy requests to `updates.maxmind.com`
- `localhost:8081` will proxy requests to `download.maxmind.com`
- `localhost:6060` will serve health checks and a pprof endpoint.

For a full list of configuration options, see the [command-line docs](docs/geoip-cache-proxy.md) and [environment variable reference](docs/envs.md).

Any flag can be provided as an env by capitaling it, changing `-` to `_`, and prefixing it with `GEOIP_`.  
For example `--cache-duration=12h` could also be configured with the env `GEOIP_CACHE_DURATION=12h`.

### Usage with geoipupdate

To configure [`maxmind/geoipupdate`](https://github.com/maxmind/geoipupdate) to use this proxy, set `GEOIP_UPDATE_HOST` to the URL of the proxy's `updates` endpoint.

Since version [v7.1.0](https://github.com/maxmind/geoipupdate/releases/tag/v7.1.0), geoipupdate supports HTTP endpoints. To configure it to use this proxy without HTTPS, set `GEOIP_UPDATE_HOST` to the host and scheme of the proxy's endpoint (e.g., `http://localhost:8080`). Note that geoipupdate defaults to HTTPS, so you must explicitly specify `http://` if the proxy is not served over HTTPS.
