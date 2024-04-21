# geoip-cache-proxy

A GeoIP database caching proxy.

Requests to `updates.maxmind.com` and `download.maxmind.com` can be passed through this proxy, and the response contents will be saved for 24 hours. Mainly intended to be used in Kubernetes. If an app downloads the database and crash-loops, this app will help avoid GeoIP rate limits.
