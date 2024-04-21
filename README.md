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

For a full list of configuration options, see the command [docs](docs/geoip-cache-proxy.md).

Any flag can be provided as an env by capitaling it, changing `-` to `_`, and prefixing it with `GEOIP_`.  
For example `--cache-duration=12h` could also be configured with the env `GEOIP_CACHE_DURATION=12h`.

### Usage with geoipupdate

To configure [`maxmind/geoipupdate`](https://github.com/maxmind/geoipupdate) to use this proxy, make sure this proxy is available over HTTPS, then set `GEOIP_UPDATE_HOST` to the host of the HTTPS endpoint.
