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
	addr = addr.Unmap()
	if !addr.IsValid() {
		return true
	}
	if addr.IsLoopback() || addr.IsPrivate() || addr.IsLinkLocalUnicast() ||
		addr.IsLinkLocalMulticast() || addr.IsMulticast() || addr.IsUnspecified() ||
		addr.IsInterfaceLocalMulticast() {
		return true
	}
	// Block CGNAT 100.64.0.0/10 and benchmark/reserved ranges.
	cgnat := netip.MustParsePrefix("100.64.0.0/10")
	return cgnat.Contains(addr)
}

func checkHostSSRF(ctx context.Context, host string) error {
	lower := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(host), "."))
	if lower == "" || lower == "localhost" || strings.HasSuffix(lower, ".localhost") ||
		strings.HasSuffix(lower, ".local") || strings.HasSuffix(lower, ".internal") {
		return fmt.Errorf("%w: host %q is not permitted", errSSRFBlocked, host)
	}
	if ip, err := netip.ParseAddr(lower); err == nil {
		if isBlockedAddr(ip) {
			return fmt.Errorf("%w: ip %s is not permitted", errSSRFBlocked, ip.String())
		}
		return nil
	}
	ips, err := net.DefaultResolver.LookupNetIP(ctx, "ip", lower)
	if err != nil {
		return fmt.Errorf("resolve host %q: %w", host, err)
	}
	if len(ips) == 0 {
		return fmt.Errorf("%w: host %q resolved no addresses", errSSRFBlocked, host)
	}
	for _, ip := range ips {
		if isBlockedAddr(ip) {
			return fmt.Errorf("%w: host %q resolved to blocked address", errSSRFBlocked, host)
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
		if err := checkHostSSRF(ctx, parsed.Hostname()); err != nil {
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
			if !m.allowPrivateFetch {
				if err := checkHostSSRF(dctx, host); err != nil {
					return nil, err
				}
			}
			return dialer.DialContext(dctx, network, net.JoinHostPort(host, port))
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
			if _, err := validateFetchURL(nextReq.URL.String()); err != nil {
				return err
			}
			if !m.allowPrivateFetch {
				if err := checkHostSSRF(nextReq.Context(), nextReq.URL.Hostname()); err != nil {
					return err
				}
			}
			return nil
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
