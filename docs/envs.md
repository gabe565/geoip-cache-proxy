# Environment Variables

| Name | Usage | Default |
| --- | --- | --- |
| `GEOIP_ACCOUNT_ID` | MaxMind account ID | `0` |
| `GEOIP_CACHE_DURATION` | Length of time to cache MaxMind response | `12h0m0s` |
| `GEOIP_DEBUG_ADDR` | Debug pprof listen address | `:6060` |
| `GEOIP_DOWNLOAD_ADDR` | Listen address | `:8081` |
| `GEOIP_DOWNLOAD_HOST` | MaxMind download host | `download.maxmind.com` |
| `GEOIP_HTTP_TIMEOUT` | HTTP request timeout | `30s` |
| `GEOIP_LICENSE_KEY` | MaxMind license key | ` ` |
| `GEOIP_LOG_FORMAT` | Log format (one of auto, color, plain, json) | `auto` |
| `GEOIP_LOG_LEVEL` | Log level (one of trace, debug, info, warn, error) | `info` |
| `GEOIP_REDIS_DB` | Redis database | `0` |
| `GEOIP_REDIS_HOST` | Redis host | `localhost` |
| `GEOIP_REDIS_PASSWORD` | Redis password | ` ` |
| `GEOIP_REDIS_PORT` | Redis port | `6379` |
| `GEOIP_TRANSLATE_INGRESS_NGINX_URLS` | Automatically translate ingress-nginx's expected file names to Maxmind paths. | `true` |
| `GEOIP_UPDATES_ADDR` | Listen address | `:8080` |
| `GEOIP_UPDATES_HOST` | MaxMind updates host | `updates.maxmind.com` |