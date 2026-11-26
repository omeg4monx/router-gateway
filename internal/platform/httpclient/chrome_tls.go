package httpclient

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
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
	dialer := &net.Dialer{
		Timeout:   defaultDialTimeout,
		KeepAlive: defaultKeepAlive,
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConns:          32,
		MaxIdleConnsPerHost:   16,
		IdleConnTimeout:       defaultIdleConnTimeout,
		TLSHandshakeTimeout:   defaultTLSHandshakeTimeout,
		ResponseHeaderTimeout: defaultResponseHeaderTimeout,
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			rawConn, err := dialer.DialContext(ctx, network, addr)
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
