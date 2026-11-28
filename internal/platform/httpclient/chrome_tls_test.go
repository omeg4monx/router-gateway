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

func TestNewChromeTLSClientWithProxySendsSOCKS5AuthAndConnect(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	gotTarget := make(chan string, 1)
	go func() {
		conn, acceptErr := ln.Accept()
		if acceptErr != nil {
			return
		}
		defer conn.Close()

		// 1. Method greeting: [0x05, nmethods, methods...]
		hdr := make([]byte, 2)
		if _, err := conn.Read(hdr); err != nil || hdr[0] != 0x05 {
			return
		}
		methods := make([]byte, hdr[1])
		if _, err := conn.Read(methods); err != nil {
			return
		}
		// Select username/password auth (0x02).
		_, _ = conn.Write([]byte{0x05, 0x02})

		// 2. RFC 1929 username/password subnegotiation.
		authVer := make([]byte, 2)
		if _, err := conn.Read(authVer); err != nil || authVer[0] != 0x01 {
			return
		}
		user := make([]byte, authVer[1])
		_, _ = conn.Read(user)
		pLen := make([]byte, 1)
		_, _ = conn.Read(pLen)
		pass := make([]byte, pLen[0])
		_, _ = conn.Read(pass)
		_, _ = conn.Write([]byte{0x01, 0x00}) // auth success

		// 3. CONNECT request: [0x05, 0x01, 0x00, atyp, ...]
		reqHdr := make([]byte, 4)
		if _, err := conn.Read(reqHdr); err != nil || reqHdr[1] != 0x01 || reqHdr[3] != 0x03 {
			return
		}
		dLen := make([]byte, 1)
		_, _ = conn.Read(dLen)
		host := make([]byte, dLen[0])
		_, _ = conn.Read(host)
		port := make([]byte, 2)
		_, _ = conn.Read(port)
		gotTarget <- string(user) + ":" + string(pass) + "@" + string(host)
		// Reject with general failure (0x01) so the test asserts tunnel negotiation cleanly.
		_, _ = conn.Write([]byte{0x05, 0x01, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	}()

	proxyURL, _ := url.Parse("socks5://bob:pw123@" + ln.Addr().String())
	client := NewChromeTLSClientWithProxy(5*time.Second, func(*http.Request) (*url.URL, error) {
		return proxyURL, nil
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://opencode.invalid:443/ping", http.NoBody)
	resp, doErr := client.Do(req)
	if resp != nil {
		_ = resp.Body.Close()
	}
	if doErr == nil || !strings.Contains(doErr.Error(), "socks5 CONNECT rejected") {
		t.Fatalf("expected socks5 CONNECT rejection error, got: %v", doErr)
	}

	select {
	case target := <-gotTarget:
		if target != "bob:pw123@opencode.invalid" {
			t.Fatalf("unexpected SOCKS5 auth/target: %q", target)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for SOCKS5 handshake")
	}
}
