package media

import (
	"context"
	"errors"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	defaultFetchMaxBytes = 256 << 10 // 256 KiB
	maxFetchBytesLimit   = 2 << 20   // 2 MiB
	maxFetchRedirects    = 5
	fetchTimeout         = 20 * time.Second
)

var (
	errSSRFBlocked     = errors.New("target address is blocked by SSRF policy")
	errInvalidFetchURL = errors.New("invalid fetch URL")
	reScriptBlock      = regexp.MustCompile(`(?is)<script\b[^>]*>.*?</script>`)
	reStyleBlock       = regexp.MustCompile(`(?is)<style\b[^>]*>.*?</style>`)
	reTitleTag         = regexp.MustCompile(`(?is)<title\b[^>]*>(.*?)</title>`)
	reHTMLTag          = regexp.MustCompile(`(?s)<[^>]+>`)
	reMultiWhitespace  = regexp.MustCompile(`\s+`)

	blockedPrefixes = []netip.Prefix{
		netip.MustParsePrefix("0.0.0.0/8"),       // "This" network (RFC 1122)
		netip.MustParsePrefix("100.64.0.0/10"),   // Shared Address Space / CGNAT (RFC 6598)
		netip.MustParsePrefix("192.0.0.0/24"),    // IETF Protocol Assignments
		netip.MustParsePrefix("192.0.2.0/24"),    // TEST-NET-1
		netip.MustParsePrefix("198.18.0.0/15"),   // Benchmarking
		netip.MustParsePrefix("198.51.100.0/24"), // TEST-NET-2
		netip.MustParsePrefix("203.0.113.0/24"),  // TEST-NET-3
		netip.MustParsePrefix("240.0.0.0/4"),     // Reserved / future use + broadcast
		netip.MustParsePrefix("64:ff9b::/96"),    // NAT64 well-known prefix
		netip.MustParsePrefix("64:ff9b:1::/48"),  // NAT64 local-use prefix
		netip.MustParsePrefix("::/96"),           // IPv4-compatible IPv6
	}
)

// WebFetchRequest describes a POST /v1/web/fetch request body.
type WebFetchRequest struct {
	URL      string `json:"url"`
	MaxBytes int    `json:"max_bytes,omitempty"`
	Provider string `json:"provider,omitempty"`
}

// WebFetchResponse describes the structured output of POST /v1/web/fetch.
type WebFetchResponse struct {
	URL         string `json:"url"`
	FinalURL    string `json:"final_url"`
	StatusCode  int    `json:"status_code"`
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	Truncated   bool   `json:"truncated"`
	ByteLength  int    `json:"byte_length"`
}

func validateFetchURL(raw string) (*url.URL, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, fmt.Errorf("%w: url is required", errInvalidFetchURL)
	}
	u, err := url.Parse(trimmed)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errInvalidFetchURL, err)
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return nil, fmt.Errorf("%w: only http and https schemes are allowed", errInvalidFetchURL)
	}
	host := strings.TrimSpace(u.Hostname())
	if host == "" {
		return nil, fmt.Errorf("%w: host is required", errInvalidFetchURL)
	}
	return u, nil
}

func isBlockedAddr(addr netip.Addr) bool {
	if !addr.IsValid() || addr.Zone() != "" {
		return true
	}
	addr = addr.Unmap()
	if addr.IsLoopback() || addr.IsPrivate() || addr.IsLinkLocalUnicast() ||
		addr.IsLinkLocalMulticast() || addr.IsMulticast() || addr.IsUnspecified() ||
		addr.IsInterfaceLocalMulticast() {
		return true
	}
	for _, pfx := range blockedPrefixes {
		if pfx.Contains(addr) {
			return true
		}
	}
	return false
}

// resolveAndCheckHostSSRF resolves host once, validates every returned IP
// against the SSRF policy, and returns the verified IPs so callers can dial
// a validated IP directly without a second DNS lookup (eliminating TOCTOU
// DNS rebinding).
func resolveAndCheckHostSSRF(ctx context.Context, host string) ([]netip.Addr, error) {
	lower := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(host), "."))
	if lower == "" || lower == "localhost" || strings.HasSuffix(lower, ".localhost") ||
		strings.HasSuffix(lower, ".local") || strings.HasSuffix(lower, ".internal") {
		return nil, fmt.Errorf("%w: host %q is not permitted", errSSRFBlocked, host)
	}
	if ip, err := netip.ParseAddr(lower); err == nil {
		if isBlockedAddr(ip) {
			return nil, fmt.Errorf("%w: ip %s is not permitted", errSSRFBlocked, ip.String())
		}
		return []netip.Addr{ip.Unmap()}, nil
	}
	ips, err := net.DefaultResolver.LookupNetIP(ctx, "ip", lower)
	if err != nil {
		return nil, fmt.Errorf("resolve host %q: %w", host, err)
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("%w: host %q resolved no addresses", errSSRFBlocked, host)
	}
	verified := make([]netip.Addr, 0, len(ips))
	for _, ip := range ips {
		if isBlockedAddr(ip) {
			return nil, fmt.Errorf("%w: host %q resolved to blocked address", errSSRFBlocked, host)
		}
		verified = append(verified, ip.Unmap())
	}
	return verified, nil
}

func validateRedirectTarget(ctx context.Context, rawURL string, allowPrivate bool) error {
	parsed, err := validateFetchURL(rawURL)
	if err != nil {
		return err
	}
	if !allowPrivate {
		if _, err := resolveAndCheckHostSSRF(ctx, parsed.Hostname()); err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) executeBuiltinFetch(ctx context.Context, req WebFetchRequest) (WebFetchResponse, error) {
	parsed, err := validateFetchURL(req.URL)
	if err != nil {
		return WebFetchResponse{}, err
	}
	if !m.allowPrivateFetch {
		if _, err := resolveAndCheckHostSSRF(ctx, parsed.Hostname()); err != nil {
			return WebFetchResponse{}, err
		}
	}

	maxBytes := req.MaxBytes
	if maxBytes <= 0 {
		maxBytes = defaultFetchMaxBytes
	}
	if maxBytes > maxFetchBytesLimit {
		maxBytes = maxFetchBytesLimit
	}

	dialer := &net.Dialer{Timeout: 10 * time.Second}
	transport := &http.Transport{
		DialContext: func(dctx context.Context, network, address string) (net.Conn, error) {
			host, port, splitErr := net.SplitHostPort(address)
			if splitErr != nil {
				return nil, splitErr
			}
			if m.allowPrivateFetch {
				return dialer.DialContext(dctx, network, net.JoinHostPort(host, port))
			}
			verifiedIPs, ssrfErr := resolveAndCheckHostSSRF(dctx, host)
			if ssrfErr != nil {
				return nil, ssrfErr
			}
			// Dial the already-validated IP directly to prevent DNS rebinding TOCTOU.
			// http.Transport preserves the original hostname for TLS SNI verification.
			var lastDialErr error
			for _, ip := range verifiedIPs {
				conn, dialErr := dialer.DialContext(dctx, network, net.JoinHostPort(ip.String(), port))
				if dialErr == nil {
					return conn, nil
				}
				lastDialErr = dialErr
			}
			return nil, lastDialErr
		},
	}
	defer transport.CloseIdleConnections()

	client := &http.Client{
		Transport: transport,
		Timeout:   fetchTimeout,
		CheckRedirect: func(nextReq *http.Request, via []*http.Request) error {
			if len(via) >= maxFetchRedirects {
				return errors.New("too many redirects")
			}
			return validateRedirectTarget(nextReq.Context(), nextReq.URL.String(), m.allowPrivateFetch)
		},
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), http.NoBody)
	if err != nil {
		return WebFetchResponse{}, err
	}
	httpReq.Header.Set("User-Agent", "omeg4mon-router/web-fetch")
	httpReq.Header.Set("Accept", "text/html,application/xhtml+xml,text/plain,application/json;q=0.9,*/*;q=0.5")

	resp, err := client.Do(httpReq)
	if err != nil {
		return WebFetchResponse{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	limited := io.LimitReader(resp.Body, int64(maxBytes)+1)
	rawBody, err := io.ReadAll(limited)
	if err != nil {
		return WebFetchResponse{}, err
	}
	truncated := len(rawBody) > maxBytes
	if truncated {
		rawBody = rawBody[:maxBytes]
	}

	contentType := resp.Header.Get("Content-Type")
	title, text := extractReadableContent(contentType, string(rawBody))

	finalURL := parsed.String()
	if resp.Request != nil && resp.Request.URL != nil {
		finalURL = resp.Request.URL.String()
	}

	return WebFetchResponse{
		URL:         parsed.String(),
		FinalURL:    finalURL,
		StatusCode:  resp.StatusCode,
		ContentType: contentType,
		Title:       title,
		Text:        text,
		Truncated:   truncated,
		ByteLength:  len(rawBody),
	}, nil
}

func extractReadableContent(contentType, body string) (title, text string) {
	lowerCT := strings.ToLower(contentType)
	isHTML := strings.Contains(lowerCT, "text/html") || strings.Contains(lowerCT, "application/xhtml")
	if !isHTML {
		return "", strings.TrimSpace(body)
	}

	if m := reTitleTag.FindStringSubmatch(body); len(m) > 1 {
		title = strings.TrimSpace(reMultiWhitespace.ReplaceAllString(html.UnescapeString(m[1]), " "))
	}
	cleaned := reScriptBlock.ReplaceAllString(body, " ")
	cleaned = reStyleBlock.ReplaceAllString(cleaned, " ")
	cleaned = reHTMLTag.ReplaceAllString(cleaned, " ")
	cleaned = html.UnescapeString(cleaned)
	text = strings.TrimSpace(reMultiWhitespace.ReplaceAllString(cleaned, " "))
	return title, text
}
