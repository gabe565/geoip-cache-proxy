## geoip-cache-proxy

A GeoIP database caching proxy

```
geoip-cache-proxy [flags]
```

### Options

```
      --account-id int                 MaxMind account ID
      --cache-duration duration        Length of time to cache MaxMind response (default 12h0m0s)
      --debug-addr string              Debug pprof listen address (default ":6060")
      --download-addr string           Listen address (default ":8081")
      --download-host string           MaxMind download host (default "download.maxmind.com")
  -h, --help                           help for geoip-cache-proxy
      --http-timeout duration          HTTP request timeout (default 30s)
      --license-key string             MaxMind license key
      --log-format string              Log format (auto, color, plain, json) (default "auto")
  -l, --log-level string               Log level (trace, debug, info, warn, error, fatal, panic) (default "info")
      --redis-db int                   Redis database
      --redis-host string              Redis host (default "localhost")
      --redis-password string          Redis password
      --redis-port uint16              Redis port (default 6379)
      --translate-ingress-nginx-urls   Automatically translate ingress-nginx's expected file names to Maxmind paths. (default true)
      --updates-addr string            Listen address (default ":8080")
      --updates-host string            MaxMind updates host (default "updates.maxmind.com")
  -v, --version                        version for geoip-cache-proxy
```

