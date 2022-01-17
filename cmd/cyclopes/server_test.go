package cyclopes

import (
	"context"
	"net"
	"net/url"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	t.Run("Server the ./ path", func(t *testing.T) {
		server := Server("./")

		defer server.Shutdown(context.Background())

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

	t.Run("Server the ./non-existing path", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Did not panic")
			}
		}()
		Server("./non-existing")
	})
}
