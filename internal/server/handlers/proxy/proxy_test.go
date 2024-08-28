package proxy

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestUpstreamUrl(t *testing.T) {
	inputs := []struct {
		name      string
		url       string
		translate bool
		expected  string
	}{
		{
			"maxmind download untranslated",
			"https://download.maxmind.com/geoip/databases/GeoLite2-Country.tar.gz",
			false,
			"https://download.maxmind.com/geoip/databases/GeoLite2-Country.tar.gz",
		},
		{
			"maxmind download translated",
			"https://download.maxmind.com/geoip/databases/GeoLite2-Country.tar.gz",
			true,
			"https://download.maxmind.com/geoip/databases/GeoLite2-Country/download?suffix=tar.gz",
		},
	}

	for _, tc := range inputs {
		t.Run(tc.name, func(t *testing.T) {
			u, err := url.Parse(tc.url)
			if err != nil {
				t.Errorf("Unparseable url: %q: %s", tc.url, err)
			}

			host := u.Host
			req, err := http.NewRequest(http.MethodGet, tc.url, nil)
			if err != nil {
				t.Errorf("failed creating http request to %q: %s", tc.url, err)
			}

			actual := upstreamURL(host, req, tc.translate)

			assert.Equal(t, tc.expected, actual.String(), "upstreamUrl does not meet expectations")

		})
	}

}
