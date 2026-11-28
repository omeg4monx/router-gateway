package httpclient

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestNewChromeTLSClientWithProxySendsHTTPConnect(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	gotConnect := make(chan string, 1)
	go func() {
		conn, acceptErr := ln.Accept()
		if acceptErr != nil {
			return
		}
		defer conn.Close()
		req, readErr := http.ReadRequest(bufio.NewReader(conn))
		if readErr != nil {
			return
		}
		gotConnect <- req.Method + " " + req.Host + " " + req.Header.Get("Proxy-Authorization")
		_, _ = conn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\nContent-Length: 0\r\n\r\n"))
	}()

	proxyURL, _ := url.Parse("http://alice:secret@" + ln.Addr().String())
	client := NewChromeTLSClientWithProxy(5*time.Second, func(*http.Request) (*url.URL, error) {
		return proxyURL, nil
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.invalid:443/ping", http.NoBody)
	resp, doErr := client.Do(req)
	if resp != nil {
		_ = resp.Body.Close()
	}
	if doErr == nil || !strings.Contains(doErr.Error(), "proxy CONNECT rejected with HTTP 502") {
		t.Fatalf("expected CONNECT rejection error, got: %v", doErr)
	}

	select {
	case line := <-gotConnect:
		if !strings.HasPrefix(line, "CONNECT example.invalid:443 Basic ") {
			t.Fatalf("unexpected CONNECT line: %q", line)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for HTTP CONNECT request on proxy")
	}
}
