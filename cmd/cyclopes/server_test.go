package cyclopes

import (
	"net"
	"net/url"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	t.Run("Server the ./ path", func(t *testing.T) {
		Server("./")
		url, err := url.Parse(DEFAULT_URL)
		if err != nil {
			t.Error(err)
		}

		host, port, _ := net.SplitHostPort(url.Host)
		_, err = net.DialTimeout("tcp", net.JoinHostPort(host, port), 2000*time.Millisecond)
		if err != nil {
			t.Error(err)
		}
	})
}
