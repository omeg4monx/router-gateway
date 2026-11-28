package httpclient

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	utls "github.com/refraction-networking/utls"
)

// bunBoringSSLClientHelloSpec mirrors the exact BoringSSL/Bun TLS 1.3
// ClientHello (17 ciphers, 13 extensions, X25519MLKEM768, ALPN http/1.1,
// JA4 t13d1713h1_5b57614c22b0_7f527383e249) used by OpenCode CLI v2 so
// OpenCode Free's Cloudflare/gateway TLS gate accepts native Go requests.
func bunBoringSSLClientHelloSpec() utls.ClientHelloSpec {
	return utls.ClientHelloSpec{
		TLSVersMin: tls.VersionTLS12,
		TLSVersMax: tls.VersionTLS13,
		CipherSuites: []uint16{
			utls.TLS_AES_128_GCM_SHA256,
			utls.TLS_AES_256_GCM_SHA384,
			utls.TLS_CHACHA20_POLY1305_SHA256,
			utls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			utls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			utls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			utls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			utls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			utls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			utls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			utls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			utls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			utls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			utls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			utls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			utls.TLS_RSA_WITH_AES_128_CBC_SHA,
			utls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		CompressionMethods: []uint8{0},
		Extensions: []utls.TLSExtension{
			&utls.SNIExtension{},
			&utls.ExtendedMasterSecretExtension{},
			&utls.RenegotiationInfoExtension{Renegotiation: utls.RenegotiateOnceAsClient},
			&utls.SupportedCurvesExtension{
				Curves: []utls.CurveID{
					utls.X25519MLKEM768,
					utls.X25519,
					utls.CurveP256,
					utls.CurveP384,
				},
			},
			&utls.SupportedPointsExtension{SupportedPoints: []uint8{0}},
			&utls.SessionTicketExtension{},
			&utls.ALPNExtension{AlpnProtocols: []string{"http/1.1"}},
			&utls.StatusRequestExtension{},
			&utls.SignatureAlgorithmsExtension{
				SupportedSignatureAlgorithms: []utls.SignatureScheme{
					utls.ECDSAWithP256AndSHA256,
					utls.PSSWithSHA256,
					utls.PKCS1WithSHA256,
					utls.ECDSAWithP384AndSHA384,
					utls.PSSWithSHA384,
					utls.PKCS1WithSHA384,
					utls.PSSWithSHA512,
					utls.PKCS1WithSHA512,
					utls.Ed25519,
					utls.ECDSAWithP521AndSHA512,
					utls.PKCS1WithSHA1,
				},
			},
			&utls.SCTExtension{},
			&utls.KeyShareExtension{
				KeyShares: []utls.KeyShare{
					{Group: utls.X25519MLKEM768},
					{Group: utls.X25519},
				},
			},
			&utls.PSKKeyExchangeModesExtension{Modes: []uint8{utls.PskModeDHE}},
			&utls.SupportedVersionsExtension{Versions: []uint16{tls.VersionTLS13, tls.VersionTLS12}},
		},
	}
}

// NewChromeTLSClient returns a pooled *http.Client whose HTTPS connections use
// BoringSSL/Bun TLS 1.3 ClientHello impersonation over HTTP/1.1.
func NewChromeTLSClient(timeout time.Duration) *http.Client {
	return NewChromeTLSClientWithProxy(timeout, http.ProxyFromEnvironment)
}

func dialThroughHTTPProxy(ctx context.Context, dialer *net.Dialer, proxyURL *url.URL, targetAddr string) (net.Conn, error) {
	proxyHost := proxyURL.Host
	if _, _, err := net.SplitHostPort(proxyHost); err != nil {
		if strings.EqualFold(proxyURL.Scheme, "https") {
			proxyHost = net.JoinHostPort(proxyHost, "443")
		} else {
			proxyHost = net.JoinHostPort(proxyHost, "80")
		}
	}

	rawConn, err := dialer.DialContext(ctx, "tcp", proxyHost)
	if err != nil {
		return nil, fmt.Errorf("httpclient: dial proxy %s: %w", proxyHost, err)
	}

	if deadline, ok := ctx.Deadline(); ok {
		_ = rawConn.SetDeadline(deadline)
	}
	defer func() {
		_ = rawConn.SetDeadline(time.Time{})
	}()

	connectReq := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Opaque: targetAddr},
		Host:   targetAddr,
		Header: make(http.Header),
	}
	if proxyURL.User != nil {
		creds := proxyURL.User.String()
		connectReq.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(creds)))
	}

	if err := connectReq.Write(rawConn); err != nil {
		_ = rawConn.Close()
		return nil, fmt.Errorf("httpclient: write proxy CONNECT: %w", err)
	}

	br := bufio.NewReader(rawConn)
	resp, err := http.ReadResponse(br, connectReq)
	if err != nil {
		_ = rawConn.Close()
		return nil, fmt.Errorf("httpclient: read proxy CONNECT response: %w", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_ = rawConn.Close()
		return nil, fmt.Errorf("httpclient: proxy CONNECT rejected with HTTP %d", resp.StatusCode)
	}
	return rawConn, nil
}

func isSOCKS5Scheme(scheme string) bool {
	s := strings.ToLower(strings.TrimSpace(scheme))
	return s == "socks5" || s == "socks5h"
}

// dialThroughSOCKS5Proxy establishes a TCP tunnel through a SOCKS5/SOCKS5H proxy
// (RFC 1928 + RFC 1929 username/password auth) so the caller can perform a
// Chrome/Bun uTLS handshake directly over the tunneled connection.
func dialThroughSOCKS5Proxy(ctx context.Context, dialer *net.Dialer, proxyURL *url.URL, targetAddr string) (net.Conn, error) {
	proxyHost := proxyURL.Host
	if _, _, err := net.SplitHostPort(proxyHost); err != nil {
		proxyHost = net.JoinHostPort(proxyHost, "1080")
	}
	host, portStr, err := net.SplitHostPort(targetAddr)
	if err != nil {
		return nil, fmt.Errorf("httpclient: invalid target addr %q: %w", targetAddr, err)
	}
	portNum, err := strconv.Atoi(portStr)
	if err != nil || portNum < 1 || portNum > 65535 {
		return nil, fmt.Errorf("httpclient: invalid target port %q", portStr)
	}

	rawConn, err := dialer.DialContext(ctx, "tcp", proxyHost)
	if err != nil {
		return nil, fmt.Errorf("httpclient: dial socks5 proxy %s: %w", proxyHost, err)
	}
	if deadline, ok := ctx.Deadline(); ok {
		_ = rawConn.SetDeadline(deadline)
	}
	defer func() {
		_ = rawConn.SetDeadline(time.Time{})
	}()

	// 1. Greeting & auth method selection.
	var greeting []byte
	var user, pass string
	if proxyURL.User != nil {
		user = proxyURL.User.Username()
		pass, _ = proxyURL.User.Password()
		greeting = []byte{0x05, 0x02, 0x00, 0x02}
	} else {
		greeting = []byte{0x05, 0x01, 0x00}
	}
	if _, err := rawConn.Write(greeting); err != nil {
		_ = rawConn.Close()
		return nil, fmt.Errorf("httpclient: socks5 write greeting: %w", err)
	}
	methodResp := make([]byte, 2)
	if _, err := io.ReadFull(rawConn, methodResp); err != nil {
		_ = rawConn.Close()
		return nil, fmt.Errorf("httpclient: socks5 read method: %w", err)
	}
	if methodResp[0] != 0x05 {
		_ = rawConn.Close()
		return nil, fmt.Errorf("httpclient: socks5 unexpected version %d", methodResp[0])
	}
	switch methodResp[1] {
	case 0x00:
		// No authentication required.
	case 0x02:
		if len(user) > 255 || len(pass) > 255 {
			_ = rawConn.Close()
			return nil, fmt.Errorf("httpclient: socks5 credentials too long")
		}
		authBuf := make([]byte, 0, 3+len(user)+len(pass))
		authBuf = append(authBuf, 0x01, byte(len(user)&0xff))
		authBuf = append(authBuf, user...)
		authBuf = append(authBuf, byte(len(pass)&0xff))
		authBuf = append(authBuf, pass...)
		if _, err := rawConn.Write(authBuf); err != nil {
			_ = rawConn.Close()
			return nil, fmt.Errorf("httpclient: socks5 write auth: %w", err)
		}
		authResp := make([]byte, 2)
		if _, err := io.ReadFull(rawConn, authResp); err != nil {
			_ = rawConn.Close()
			return nil, fmt.Errorf("httpclient: socks5 read auth: %w", err)
		}
		if authResp[1] != 0x00 {
			_ = rawConn.Close()
			return nil, fmt.Errorf("httpclient: socks5 auth rejected with status %d", authResp[1])
		}
	default:
		_ = rawConn.Close()
		return nil, fmt.Errorf("httpclient: socks5 unsupported auth method %d", methodResp[1])
	}

	// 2. CONNECT command.
	reqBuf := make([]byte, 0, 7+len(host))
	reqBuf = append(reqBuf, 0x05, 0x01, 0x00)
	if ip := net.ParseIP(host); ip != nil {
		if ip4 := ip.To4(); ip4 != nil {
			reqBuf = append(reqBuf, 0x01)
			reqBuf = append(reqBuf, ip4...)
		} else {
			reqBuf = append(reqBuf, 0x04)
			reqBuf = append(reqBuf, ip.To16()...)
		}
	} else {
		if len(host) > 255 {
			_ = rawConn.Close()
			return nil, fmt.Errorf("httpclient: socks5 hostname too long")
		}
		reqBuf = append(reqBuf, 0x03, byte(len(host)&0xff))
		reqBuf = append(reqBuf, host...)
	}
	var portBytes [2]byte
	binary.BigEndian.PutUint16(portBytes[:], uint16(portNum))
	reqBuf = append(reqBuf, portBytes[:]...)

	if _, err := rawConn.Write(reqBuf); err != nil {
		_ = rawConn.Close()
		return nil, fmt.Errorf("httpclient: socks5 write CONNECT: %w", err)
	}

	replyHdr := make([]byte, 4)
	if _, err := io.ReadFull(rawConn, replyHdr); err != nil {
		_ = rawConn.Close()
		return nil, fmt.Errorf("httpclient: socks5 read CONNECT reply: %w", err)
	}
	if replyHdr[0] != 0x05 || replyHdr[1] != 0x00 {
		_ = rawConn.Close()
		return nil, fmt.Errorf("httpclient: socks5 CONNECT rejected with code %d", replyHdr[1])
	}

	// Drain bound address + port from reply.
	var addrLen int
	switch replyHdr[3] {
	case 0x01:
		addrLen = 4
	case 0x04:
		addrLen = 16
	case 0x03:
		lBuf := make([]byte, 1)
		if _, err := io.ReadFull(rawConn, lBuf); err != nil {
			_ = rawConn.Close()
			return nil, fmt.Errorf("httpclient: socks5 read domain length: %w", err)
		}
		addrLen = int(lBuf[0])
	default:
		_ = rawConn.Close()
		return nil, fmt.Errorf("httpclient: socks5 unsupported reply atyp %d", replyHdr[3])
	}
	if _, err := io.ReadFull(rawConn, make([]byte, addrLen+2)); err != nil {
		_ = rawConn.Close()
		return nil, fmt.Errorf("httpclient: socks5 read bound addr: %w", err)
	}
	return rawConn, nil
}

// NewChromeTLSClientWithProxy returns a pooled *http.Client using BoringSSL/Bun
// TLS 1.3 ClientHello impersonation and the supplied proxy resolver function.
// For HTTP and SOCKS5 proxies targeting HTTPS endpoints, the tunnel is established
// inside DialTLSContext so the inner TLS handshake still uses Chrome uTLS
// instead of Go's standard crypto/tls.
func NewChromeTLSClientWithProxy(timeout time.Duration, proxyFunc func(*http.Request) (*url.URL, error)) *http.Client {
	if proxyFunc == nil {
		proxyFunc = http.ProxyFromEnvironment
	}
	dialer := &net.Dialer{
		Timeout:   defaultDialTimeout,
		KeepAlive: defaultKeepAlive,
	}

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			u, err := proxyFunc(req)
			if err != nil || u == nil {
				return u, err
			}
			// Delegate HTTP and SOCKS5 HTTPS tunnels to DialTLSContext so uTLS is preserved.
			if req != nil && req.URL != nil && strings.EqualFold(req.URL.Scheme, "https") &&
				(strings.EqualFold(u.Scheme, "http") || isSOCKS5Scheme(u.Scheme)) {
				return nil, nil
			}
			return u, nil
		},
		DialContext:           dialer.DialContext,
		MaxIdleConns:          32,
		MaxIdleConnsPerHost:   16,
		IdleConnTimeout:       defaultIdleConnTimeout,
		TLSHandshakeTimeout:   defaultTLSHandshakeTimeout,
		ResponseHeaderTimeout: defaultResponseHeaderTimeout,
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			var (
				rawConn net.Conn
				err     error
			)
			probeReq := &http.Request{URL: &url.URL{Scheme: "https", Host: addr}}
			if proxyURL, proxyErr := proxyFunc(probeReq.WithContext(ctx)); proxyErr == nil && proxyURL != nil {
				switch {
				case strings.EqualFold(proxyURL.Scheme, "http"):
					rawConn, err = dialThroughHTTPProxy(ctx, dialer, proxyURL, addr)
				case isSOCKS5Scheme(proxyURL.Scheme):
					rawConn, err = dialThroughSOCKS5Proxy(ctx, dialer, proxyURL, addr)
				default:
					rawConn, err = dialer.DialContext(ctx, network, addr)
				}
			} else {
				rawConn, err = dialer.DialContext(ctx, network, addr)
			}
			if err != nil {
				return nil, err
			}
			host, _, splitErr := net.SplitHostPort(addr)
			if splitErr != nil {
				host = addr
			}
			uConn := utls.UClient(rawConn, &utls.Config{
				ServerName: host,
				MinVersion: tls.VersionTLS12,
			}, utls.HelloCustom)
			spec := bunBoringSSLClientHelloSpec()
			if err := uConn.ApplyPreset(&spec); err != nil {
				_ = rawConn.Close()
				return nil, err
			}
			if err := uConn.HandshakeContext(ctx); err != nil {
				_ = rawConn.Close()
				return nil, err
			}
			return uConn, nil
		},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}

var (
	sharedChromeTLSOnce   sync.Once
	sharedChromeTLSClient *http.Client
)

// SharedChromeTLSClient returns a process-wide pooled BoringSSL/Bun uTLS client
// suitable for OpenCode Free probes and requests.
func SharedChromeTLSClient() *http.Client {
	sharedChromeTLSOnce.Do(func() {
		sharedChromeTLSClient = NewChromeTLSClient(0)
	})
	return sharedChromeTLSClient
}
