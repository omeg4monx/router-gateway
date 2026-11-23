package proxypools

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	maxScanCandidates = 50
	maxScanWorkers    = 8
)

// PoolTestResponse reports the outcome of probing a configured proxy pool.
type PoolTestResponse struct {
	ID        string `json:"id"`
	Healthy   bool   `json:"healthy"`
	Status    string `json:"status"`
	LatencyMS int64  `json:"latency_ms"`
	Error     string `json:"error,omitempty"`
	TestedAt  string `json:"tested_at"`
}

// ScanCandidateResult reports the probe result for a single candidate proxy URL.
type ScanCandidateResult struct {
	DisplayURL     string `json:"display_url"`
	Scheme         string `json:"scheme,omitempty"`
	HasCredentials bool   `json:"has_credentials"`
	Healthy        bool   `json:"healthy"`
	LatencyMS      int64  `json:"latency_ms"`
	Error          string `json:"error,omitempty"`
	AddedPoolID    string `json:"added_pool_id,omitempty"`
}

// ScanResponse holds the batch probe outcomes returned by POST /api/proxypools/scan.
type ScanResponse struct {
	Results []ScanCandidateResult `json:"results"`
}

func probeProxy(ctx context.Context, proxyURL *url.URL, targetURL string, timeout time.Duration) (latencyMS int64, err error) {
	if proxyURL == nil {
		return 0, errors.New("proxy url is nil")
	}
	if timeout <= 0 {
		timeout = defaultProbeTimeout
	}

	probeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	addr := proxyDialAddress(proxyURL)
	dialer := &net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(probeCtx, "tcp", addr)
	if err != nil {
		return time.Since(start).Milliseconds(), fmt.Errorf("dial proxy %s: %w", addr, err)
	}
	_ = conn.Close()

	if targetURL != "" {
		parsedTarget, parseErr := url.Parse(targetURL)
		if parseErr != nil || (parsedTarget.Scheme != "http" && parsedTarget.Scheme != "https") || parsedTarget.Host == "" {
			return time.Since(start).Milliseconds(), errors.New("invalid target_url")
		}
		tr := &http.Transport{
			Proxy:                 http.ProxyURL(proxyURL),
			DialContext:           dialer.DialContext,
			ResponseHeaderTimeout: timeout,
		}
		defer tr.CloseIdleConnections()

		client := &http.Client{Transport: tr, Timeout: timeout}
		req, reqErr := http.NewRequestWithContext(probeCtx, http.MethodHead, parsedTarget.String(), http.NoBody)
		if reqErr != nil {
			return time.Since(start).Milliseconds(), reqErr
		}
		resp, doErr := client.Do(req)
		if doErr != nil {
			return time.Since(start).Milliseconds(), fmt.Errorf("proxy http probe: %w", doErr)
		}
		_ = resp.Body.Close()
		if resp.StatusCode >= 500 {
			return time.Since(start).Milliseconds(), fmt.Errorf("proxy target returned status %d", resp.StatusCode)
		}
	}

	return time.Since(start).Milliseconds(), nil
}

func proxyDialAddress(u *url.URL) string {
	if u.Port() != "" {
		return u.Host
	}
	switch u.Scheme {
	case "https":
		return net.JoinHostPort(u.Hostname(), "443")
	case "socks5", "socks5h":
		return net.JoinHostPort(u.Hostname(), "1080")
	default:
		return net.JoinHostPort(u.Hostname(), "80")
	}
}

func runConcurrentScan(
	ctx context.Context,
	candidates []string,
	timeout time.Duration,
	onHealthy func(idx int, parsed parsedProxyURL, latencyMS int64) (string, error),
) []ScanCandidateResult {
	results := make([]ScanCandidateResult, len(candidates))
	if len(candidates) == 0 {
		return results
	}

	workers := min(maxScanWorkers, len(candidates))
	jobsCh := make(chan int, len(candidates))
	for i := range candidates {
		jobsCh <- i
	}
	close(jobsCh)

	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for idx := range jobsCh {
				raw := candidates[idx]
				parsed, err := parseAndValidateProxyURL(raw)
				if err != nil {
					results[idx] = ScanCandidateResult{
						DisplayURL: redactRawCandidate(raw),
						Healthy:    false,
						Error:      err.Error(),
					}
					continue
				}
				lat, probeErr := probeProxy(ctx, parsed.Parsed, "", timeout)
				res := ScanCandidateResult{
					DisplayURL:     parsed.DisplayURL,
					Scheme:         parsed.Scheme,
					HasCredentials: parsed.HasCredentials,
					Healthy:        probeErr == nil,
					LatencyMS:      lat,
				}
				if probeErr != nil {
					res.Error = probeErr.Error()
				} else if onHealthy != nil {
					if poolID, addErr := onHealthy(idx, parsed, lat); addErr != nil {
						res.Error = addErr.Error()
					} else {
						res.AddedPoolID = poolID
					}
				}
				results[idx] = res
			}
		}()
	}
	wg.Wait()
	return results
}

func redactRawCandidate(raw string) string {
	if u, err := url.Parse(raw); err == nil && u.Host != "" {
		if u.User != nil {
			return u.Scheme + "://***@" + u.Host
		}
		return u.Scheme + "://" + u.Host
	}
	return "<invalid-url>"
}
