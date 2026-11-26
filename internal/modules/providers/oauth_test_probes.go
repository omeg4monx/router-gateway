package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

// oauthTestProbeSpec declaratively describes how to test an OAuth provider's
// live connectivity without wasting billable inference tokens (adopted from
// xrouter/src/app/api/providers/[id]/test/testUtils.js).
type oauthTestProbeSpec struct {
	LiveHostSuffix        string
	Method                string
	URL                   string
	BuildURL              func(accessToken string) string
	NoAuthHeader          bool
	ExtraHeaders          map[string]string
	Refreshable           bool
	SoftSuccessOnThrottle bool
}

var oauthTestProbes = map[string]oauthTestProbeSpec{
	"antigravity": {
		LiveHostSuffix: "googleapis.com",
		Method:         http.MethodGet,
		URL:            googleUserInfoV1URL + "?alt=json",
		ExtraHeaders:   map[string]string{"x-request-source": "local"},
		Refreshable:    true,
	},
	"gemini-oauth": {
		LiveHostSuffix: "googleapis.com",
		Method:         http.MethodGet,
		URL:            googleUserInfoV1URL + "?alt=json",
		Refreshable:    true,
	},
	"github-copilot": {
		LiveHostSuffix: "githubcopilot.com",
		Method:         http.MethodGet,
		URL:            githubUserInfoURL,
		ExtraHeaders: map[string]string{
			"User-Agent": "omeg4mon-router",
			"Accept":     "application/vnd.github+json",
		},
		Refreshable: false,
	},
	"claude": {
		LiveHostSuffix: "anthropic.com",
		Method:         http.MethodGet,
		URL:            "https://api.anthropic.com/v1/models",
		ExtraHeaders: map[string]string{
			"anthropic-version": defaultAnthropicVersion,
			"anthropic-beta":    "oauth-2025-04-20",
		},
		Refreshable:           true,
		SoftSuccessOnThrottle: true,
	},
	"iflow": {
		LiveHostSuffix: "iflow.cn",
		Method:         http.MethodGet,
		BuildURL: func(token string) string {
			return "https://iflow.cn/api/oauth/getUserInfo?accessToken=" + url.QueryEscape(token)
		},
		NoAuthHeader: true,
		Refreshable:  true,
	},
	"gitlab": {
		LiveHostSuffix: "gitlab.com",
		Method:         http.MethodGet,
		URL:            "https://gitlab.com/api/v4/user",
		Refreshable:    true,
	},
}

func (m *Module) executeOAuthTestProbe(
	ctx context.Context,
	connectionID string,
	creds contracts.UpstreamCredentials,
	spec oauthTestProbeSpec,
) (result testConnectionResponse, healthStatus string) {
	method := spec.Method
	if method == "" {
		method = http.MethodGet
	}

	doProbe := func(token string) (int, error) {
		targetURL := spec.URL
		if spec.BuildURL != nil {
			targetURL = spec.BuildURL(token)
		}
		req, err := http.NewRequestWithContext(ctx, method, targetURL, http.NoBody)
		if err != nil {
			return 0, err
		}
		req.Header.Set("Accept", "application/json")
		if !spec.NoAuthHeader && token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		for k, v := range spec.ExtraHeaders {
			req.Header.Set(k, v)
		}
		resp, err := m.httpClient.Do(req)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxTestErrorBodyBytes))
		return resp.StatusCode, nil
	}

	start := time.Now()
	statusCode, probeErr := doProbe(creds.AccessToken)
	if spec.Refreshable && connectionID != "" &&
		(probeErr != nil || statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden) {
		if refreshed, refreshErr := m.RefreshOAuthToken(ctx, connectionID); refreshErr == nil && refreshed.AccessToken != "" {
			statusCode, probeErr = doProbe(refreshed.AccessToken)
		}
	}
	latencyMS := time.Since(start).Milliseconds()

	if probeErr != nil {
		return testConnectionResponse{
			OK:        false,
			LatencyMS: latencyMS,
			Error:     truncateError(sanitizeUpstreamError(probeErr.Error(), creds)),
		}, "error"
	}
	if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
		return testConnectionResponse{
			OK:         true,
			StatusCode: statusCode,
			LatencyMS:  latencyMS,
		}, "ok"
	}
	if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		return testConnectionResponse{
			OK:         false,
			StatusCode: statusCode,
			LatencyMS:  latencyMS,
			Error:      "Token invalid or revoked",
		}, "error"
	}
	if spec.SoftSuccessOnThrottle &&
		(statusCode == http.StatusRequestTimeout || statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError) {
		return testConnectionResponse{
			OK:         true,
			StatusCode: statusCode,
			LatencyMS:  latencyMS,
			Error:      fmt.Sprintf("upstream returned HTTP %d — transient throttle", statusCode),
		}, "ok"
	}
	return testConnectionResponse{
		OK:         false,
		StatusCode: statusCode,
		LatencyMS:  latencyMS,
		Error:      strings.TrimSpace(fmt.Sprintf("upstream returned HTTP %d", statusCode)),
	}, "error"
}
