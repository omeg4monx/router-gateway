package cloudsync

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strings"
)

func validateEndpointURL(ctx context.Context, rawURL string, allowPrivate bool) (*url.URL, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return nil, errors.New("endpoint_url is required")
	}
	u, err := url.Parse(trimmed)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint_url: %w", err)
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "https" && scheme != "http" {
		return nil, errors.New("endpoint_url must use http or https scheme")
	}
	if !allowPrivate && scheme != "https" {
		return nil, errors.New("endpoint_url must use https unless private endpoints are allowed")
	}
	host := strings.TrimSpace(u.Hostname())
	if host == "" {
		return nil, errors.New("endpoint_url host is required")
	}
	if allowPrivate {
		return u, nil
	}

	lowerHost := strings.ToLower(host)
	if lowerHost == "localhost" || strings.HasSuffix(lowerHost, ".localhost") || strings.HasSuffix(lowerHost, ".local") || strings.HasSuffix(lowerHost, ".internal") {
		return nil, errors.New("private or local endpoint host is not allowed")
	}

	if ip, err := netip.ParseAddr(host); err == nil {
		if isBlockedAddr(ip) {
			return nil, errors.New("private, loopback, or link-local IP address is not allowed")
		}
		return u, nil
	}

	ips, err := net.DefaultResolver.LookupNetIP(ctx, "ip", host)
	if err == nil {
		for _, ip := range ips {
			if isBlockedAddr(ip) {
				return nil, errors.New("endpoint host resolves to a private, loopback, or link-local IP address")
			}
		}
	}
	return u, nil
}

func isBlockedAddr(ip netip.Addr) bool {
	ip = ip.Unmap()
	return ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() ||
		ip.IsUnspecified()
}
