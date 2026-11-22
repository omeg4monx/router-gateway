package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
)

const (
	maxProvidersBodyBytes   = 64 << 10 // 64 KiB
	maxTestErrorBodyBytes   = 4096
	maxNameRunes            = 128
	maxIDRunes              = 64
	maxSliceItems           = 256
	maxSliceItemRunes       = 128
	maxHeadersCount         = 64
	maxHeaderKeyRunes       = 128
	maxHeaderValueRunes     = 1024
	maxTestErrorSummaryLen  = 256
	defaultConnectionPrefix = "conn_"
	defaultNodePrefix       = "node_"
)

type catalogOAuthDescriptor struct {
	AuthorizeURL       string   `json:"authorize_url"`
	TokenURL           string   `json:"token_url"`
	Scopes             []string `json:"scopes"`
	DefaultRedirectURI string   `json:"default_redirect_uri"`
	ClientIDConfigured bool     `json:"client_id_configured"`
}

type catalogAdapterDescriptor struct {
	ID             string                  `json:"id"`
	DisplayName    string                  `json:"display_name"`
	Category       string                  `json:"category"`
	Protocol       string                  `json:"protocol"`
	AuthType       string                  `json:"auth_type"`
	DefaultBaseURL string                  `json:"default_base_url"`
	DefaultModels  []string                `json:"default_models"`
	OAuth          *catalogOAuthDescriptor `json:"oauth"`
}

type catalogResponse struct {
	Adapters    []catalogAdapterDescriptor `json:"adapters"`
	CustomNodes []ProviderNode             `json:"custom_nodes"`
}

type listNodesResponse struct {
	Nodes []ProviderNode `json:"nodes"`
}

type createNodeRequest struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Protocol string            `json:"protocol"`
	BaseURL  string            `json:"base_url"`
	Models   []string          `json:"models"`
	Headers  map[string]string `json:"headers"`
}

type listConnectionsResponse struct {
	Connections []ConnectionSummary `json:"connections"`
}

type createConnectionRequest struct {
	ProviderID   string            `json:"provider_id"`
	Name         string            `json:"name"`
	Protocol     string            `json:"protocol"`
	BaseURL      string            `json:"base_url"`
	AuthType     string            `json:"auth_type"`
	APIKey       string            `json:"api_key"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	Priority     *int              `json:"priority"`
	Enabled      *bool             `json:"enabled"`
	Models       []string          `json:"models"`
	Headers      map[string]string `json:"headers"`
}

type patchConnectionRequest struct {
	Name         *string            `json:"name"`
	BaseURL      *string            `json:"base_url"`
	APIKey       *string            `json:"api_key"`
	AccessToken  *string            `json:"access_token"`
	RefreshToken *string            `json:"refresh_token"`
	Priority     *int               `json:"priority"`
	Enabled      *bool              `json:"enabled"`
	Models       *[]string          `json:"models"`
	Headers      *map[string]string `json:"headers"`
}

type testConnectionResponse struct {
	OK         bool   `json:"ok"`
	StatusCode int    `json:"status_code"`
	LatencyMS  int64  `json:"latency_ms"`
	Error      string `json:"error"`
}

type oauthStartRequest struct {
	ConnectionName string   `json:"connection_name"`
	RedirectURI    string   `json:"redirect_uri"`
	ClientID       string   `json:"client_id"`
	ClientSecret   string   `json:"client_secret"`
	BaseURL        string   `json:"base_url"`
	Models         []string `json:"models"`
}

type oauthStartResponse struct {
	State        string `json:"state"`
	AuthorizeURL string `json:"authorize_url"`
	RedirectURI  string `json:"redirect_uri"`
	ExpiresAt    string `json:"expires_at"`
}

type oauthExchangeRequest struct {
	State          string            `json:"state"`
	Code           string            `json:"code"`
	CallbackURL    string            `json:"callback_url"`
	ConnectionID   string            `json:"connection_id"`
	ConnectionName string            `json:"connection_name"`
	Priority       *int              `json:"priority"`
	Enabled        *bool             `json:"enabled"`
	Headers        map[string]string `json:"headers"`
}

func (m *Module) handleGetCatalog(w http.ResponseWriter, r *http.Request) {
	nodes, err := m.repo.ListNodes(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load custom provider nodes")
		return
	}

	adapters := m.ListAdapters()
	descriptors := make([]catalogAdapterDescriptor, 0, len(adapters))
	for _, adapter := range adapters {
		desc := catalogAdapterDescriptor{
			ID:             adapter.ID(),
			DisplayName:    adapter.DisplayName(),
			Category:       adapter.Category(),
			Protocol:       adapter.Protocol(),
			AuthType:       adapter.AuthType(),
			DefaultBaseURL: adapter.DefaultBaseURL(),
			DefaultModels:  cloneStringSlice(adapter.DefaultModels()),
		}
		if adapter.OAuthConfig() != nil {
			resolved, err := m.resolveOAuthConfig(adapter.ID(), "", "")
			if err == nil {
				desc.OAuth = &catalogOAuthDescriptor{
					AuthorizeURL:       resolved.AuthorizeURL,
					TokenURL:           resolved.TokenURL,
					Scopes:             cloneStringSlice(resolved.Scopes),
					DefaultRedirectURI: resolved.DefaultRedirectURI,
					ClientIDConfigured: resolved.ClientID != "",
				}
			}
		}
		descriptors = append(descriptors, desc)
	}

	httpserver.WriteJSON(w, http.StatusOK, catalogResponse{
		Adapters:    descriptors,
		CustomNodes: nodes,
	})
}

func (m *Module) handleListNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := m.repo.ListNodes(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to list custom provider nodes")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, listNodesResponse{Nodes: nodes})
}

func (m *Module) handleCreateNode(w http.ResponseWriter, r *http.Request) {
	var req createNodeRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	var (
		nodeID string
		err    error
	)
	if strings.TrimSpace(req.ID) != "" {
		nodeID, err = validateSlugID(req.ID)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		if _, exists := m.GetAdapter(nodeID); exists {
			httpserver.WriteError(w, r, http.StatusConflict, "node_conflict", "node id conflicts with a built-in provider adapter")
			return
		}
	} else {
		nodeID, err = generateID(defaultNodePrefix)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate node identifier")
			return
		}
	}

	name, err := validateName(req.Name)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	protocol := strings.TrimSpace(req.Protocol)
	if !isValidProtocol(protocol) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "protocol must be one of openai-chat, openai-responses, anthropic, gemini")
		return
	}
	baseURL, err := validateHTTPURL(req.BaseURL)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "base_url: "+err.Error())
		return
	}
	models, err := validateStringSlice("models", req.Models)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	headers, err := validateHeaders(req.Headers)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	created, err := m.repo.CreateNode(r.Context(), ProviderNode{
		ID:       nodeID,
		Name:     name,
		Protocol: protocol,
		BaseURL:  baseURL,
		Models:   models,
		Headers:  headers,
	})
	if err != nil {
		if errors.Is(err, ErrNodeAlreadyExists) {
			httpserver.WriteError(w, r, http.StatusConflict, "node_conflict", "custom provider node id already exists")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to create custom provider node")
		return
	}

	m.publishChanged("created", "node", created.ID, created.ID, actorFromContext(r.Context()), created.UpdatedAt)
	httpserver.WriteJSON(w, http.StatusCreated, created)
}

func (m *Module) handleDeleteNode(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "node id is required")
		return
	}

	deletedAt, err := m.repo.DeleteNode(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNodeNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "node_not_found", "custom provider node not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to delete custom provider node")
		return
	}

	m.publishChanged("deleted", "node", id, id, actorFromContext(r.Context()), deletedAt)
	w.WriteHeader(http.StatusNoContent)
}

func (m *Module) handleListConnections(w http.ResponseWriter, r *http.Request) {
	items, err := m.repo.ListConnections(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to list provider connections")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, listConnectionsResponse{Connections: items})
}

func (m *Module) handleCreateConnection(w http.ResponseWriter, r *http.Request) {
	var req createConnectionRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	providerID := strings.TrimSpace(req.ProviderID)
	if providerID == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "provider_id is required")
		return
	}

	name, err := validateName(req.Name)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	var (
		defaultProtocol string
		defaultBaseURL  string
		defaultAuthType string
		defaultModels   []string
		defaultHeaders  map[string]string
	)
	if adapter, ok := m.GetAdapter(providerID); ok {
		defaultProtocol = adapter.Protocol()
		defaultBaseURL = adapter.DefaultBaseURL()
		defaultAuthType = adapter.AuthType()
		defaultModels = adapter.DefaultModels()
		defaultHeaders = map[string]string{}
	} else {
		node, nodeErr := m.repo.GetNode(r.Context(), providerID)
		if nodeErr != nil {
			if errors.Is(nodeErr, ErrNodeNotFound) {
				httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_provider", "unknown provider_id")
				return
			}
			httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to resolve provider")
			return
		}
		defaultProtocol = node.Protocol
		defaultBaseURL = node.BaseURL
		defaultAuthType = authTypeAPIKey
		defaultModels = node.Models
		defaultHeaders = node.Headers
	}

	protocol := strings.TrimSpace(req.Protocol)
	if protocol == "" {
		protocol = defaultProtocol
	}
	if !isValidProtocol(protocol) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "protocol must be one of openai-chat, openai-responses, anthropic, gemini")
		return
	}

	rawBaseURL := strings.TrimSpace(req.BaseURL)
	if rawBaseURL == "" {
		rawBaseURL = defaultBaseURL
	}
	baseURL, err := validateHTTPURL(rawBaseURL)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "base_url: "+err.Error())
		return
	}

	authType := strings.TrimSpace(req.AuthType)
	if authType == "" {
		authType = defaultAuthType
	}
	if !isValidAuthType(authType) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "auth_type must be api_key or oauth")
		return
	}

	creds := storedCredentials{
		APIKey:       strings.TrimSpace(req.APIKey),
		AccessToken:  strings.TrimSpace(req.AccessToken),
		RefreshToken: strings.TrimSpace(req.RefreshToken),
	}
	switch authType {
	case authTypeAPIKey:
		if creds.APIKey == "" {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "api_key is required for api_key connections")
			return
		}
	case authTypeOAuth:
		if creds.AccessToken == "" && creds.RefreshToken == "" {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "access_token or refresh_token is required for oauth connections")
			return
		}
	}

	priority := defaultPriorityValue
	if req.Priority != nil {
		if *req.Priority < 0 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "priority must be non-negative")
			return
		}
		priority = *req.Priority
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	var models []string
	if req.Models != nil {
		models, err = validateStringSlice("models", req.Models)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	} else {
		models = cloneStringSlice(defaultModels)
	}

	var headers map[string]string
	if req.Headers != nil {
		headers, err = validateHeaders(req.Headers)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	} else {
		headers = cloneStringMap(defaultHeaders)
	}

	connID, err := generateID(defaultConnectionPrefix)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate connection identifier")
		return
	}

	enc, err := m.sealCredentials(creds)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to encrypt connection credentials")
		return
	}

	summary, err := m.repo.CreateConnection(r.Context(), createConnectionParams{
		ID:             connID,
		ProviderID:     providerID,
		Name:           name,
		Protocol:       protocol,
		BaseURL:        baseURL,
		AuthType:       authType,
		CredentialsEnc: enc,
		Priority:       priority,
		Enabled:        enabled,
		Models:         models,
		Headers:        headers,
	})
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to create provider connection")
		return
	}

	m.publishChanged("created", "connection", summary.ID, summary.ProviderID, actorFromContext(r.Context()), summary.UpdatedAt)
	httpserver.WriteJSON(w, http.StatusCreated, summary)
}

func (m *Module) handleGetConnection(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "connection id is required")
		return
	}

	summary, err := m.repo.GetConnectionSummary(r.Context(), id)
	if err != nil {
		if errors.Is(err, contracts.ErrConnectionNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "connection_not_found", "provider connection not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load provider connection")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, summary)
}

func (m *Module) handlePatchConnection(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "connection id is required")
		return
	}

	var req patchConnectionRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	if req.Name == nil &&
		req.BaseURL == nil &&
		req.APIKey == nil &&
		req.AccessToken == nil &&
		req.RefreshToken == nil &&
		req.Priority == nil &&
		req.Enabled == nil &&
		req.Models == nil &&
		req.Headers == nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "at least one field must be provided")
		return
	}

	var (
		cleanName    string
		cleanBaseURL string
		cleanModels  []string
		cleanHeaders map[string]string
		err          error
	)
	if req.Name != nil {
		cleanName, err = validateName(*req.Name)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	}
	if req.BaseURL != nil {
		cleanBaseURL, err = validateHTTPURL(*req.BaseURL)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "base_url: "+err.Error())
			return
		}
	}
	if req.Priority != nil && *req.Priority < 0 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "priority must be non-negative")
		return
	}
	if req.Models != nil {
		cleanModels, err = validateStringSlice("models", *req.Models)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	}
	if req.Headers != nil {
		cleanHeaders, err = validateHeaders(*req.Headers)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	}

	var validationErr error
	summary, err := m.repo.UpdateConnection(r.Context(), id, func(cur connectionRecord) (updateConnectionValues, error) {
		vals := updateConnectionValues{
			Name:           cur.Summary.Name,
			BaseURL:        cur.Summary.BaseURL,
			CredentialsEnc: cur.CredentialsEnc,
			Priority:       cur.Summary.Priority,
			Enabled:        cur.Summary.Enabled,
			Models:         cur.Summary.Models,
			Headers:        cur.Summary.Headers,
			TokenExpiresAt: cur.Summary.TokenExpiresAt,
		}
		if req.Name != nil {
			vals.Name = cleanName
		}
		if req.BaseURL != nil {
			vals.BaseURL = cleanBaseURL
		}
		if req.Priority != nil {
			vals.Priority = *req.Priority
		}
		if req.Enabled != nil {
			vals.Enabled = *req.Enabled
		}
		if req.Models != nil {
			vals.Models = cleanModels
		}
		if req.Headers != nil {
			vals.Headers = cleanHeaders
		}

		if req.APIKey != nil || req.AccessToken != nil || req.RefreshToken != nil {
			stored, decErr := m.unsealCredentials(cur.CredentialsEnc)
			if decErr != nil {
				return updateConnectionValues{}, decErr
			}
			if req.APIKey != nil {
				v := strings.TrimSpace(*req.APIKey)
				if cur.Summary.AuthType == authTypeAPIKey && v == "" {
					validationErr = errors.New("api_key must not be empty")
					return updateConnectionValues{}, validationErr
				}
				stored.APIKey = v
			}
			if req.AccessToken != nil {
				stored.AccessToken = strings.TrimSpace(*req.AccessToken)
			}
			if req.RefreshToken != nil {
				stored.RefreshToken = strings.TrimSpace(*req.RefreshToken)
			}
			if cur.Summary.AuthType == authTypeOAuth && stored.AccessToken == "" && stored.RefreshToken == "" {
				validationErr = errors.New("oauth connection requires access_token or refresh_token")
				return updateConnectionValues{}, validationErr
			}
			enc, encErr := m.sealCredentials(stored)
			if encErr != nil {
				return updateConnectionValues{}, encErr
			}
			vals.CredentialsEnc = enc
		}
		return vals, nil
	})
	if err != nil {
		if validationErr != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", validationErr.Error())
			return
		}
		if errors.Is(err, contracts.ErrConnectionNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "connection_not_found", "provider connection not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to update provider connection")
		return
	}

	m.publishChanged("updated", "connection", summary.ID, summary.ProviderID, actorFromContext(r.Context()), summary.UpdatedAt)
	httpserver.WriteJSON(w, http.StatusOK, summary)
}

func (m *Module) handleDeleteConnection(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "connection id is required")
		return
	}

	providerID, deletedAt, err := m.repo.DeleteConnection(r.Context(), id)
	if err != nil {
		if errors.Is(err, contracts.ErrConnectionNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "connection_not_found", "provider connection not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to delete provider connection")
		return
	}

	m.publishChanged("deleted", "connection", id, providerID, actorFromContext(r.Context()), deletedAt)
	w.WriteHeader(http.StatusNoContent)
}

func (m *Module) handleTestConnection(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "connection id is required")
		return
	}

	creds, err := m.ResolveCredentials(r.Context(), id)
	if err != nil {
		if errors.Is(err, contracts.ErrConnectionNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "connection_not_found", "provider connection not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to resolve connection credentials")
		return
	}

	var testReq *http.Request
	if adapter, ok := m.GetAdapter(creds.ProviderID); ok {
		testReq, err = adapter.BuildTestRequest(r.Context(), creds)
	} else {
		testReq, err = buildProtocolTestRequest(r.Context(), creds)
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "failed to build connectivity test request")
		return
	}

	start := time.Now()
	resp, doErr := m.httpClient.Do(testReq) //nolint:gosec // G704: upstream provider base_url is configured by operator and validated http/https
	latencyMS := time.Since(start).Milliseconds()

	result := testConnectionResponse{
		LatencyMS: latencyMS,
	}
	statusText := "error"

	if doErr != nil {
		result.OK = false
		result.Error = truncateError(sanitizeUpstreamError(doErr.Error(), creds))
	} else {
		defer resp.Body.Close()
		bodySnippet, _ := io.ReadAll(io.LimitReader(resp.Body, maxTestErrorBodyBytes))
		result.StatusCode = resp.StatusCode
		if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
			result.OK = true
			statusText = "ok"
		} else {
			result.OK = false
			msg := fmt.Sprintf("upstream returned HTTP %d", resp.StatusCode)
			if trimmed := strings.TrimSpace(string(bodySnippet)); trimmed != "" {
				msg = msg + ": " + sanitizeUpstreamError(trimmed, creds)
			}
			result.Error = truncateError(msg)
		}
	}

	testedAt := time.Now().UTC().Format(rfc3339NanoFixed)
	_ = m.repo.UpdateTestResult(r.Context(), id, testedAt, statusText, result.Error)

	httpserver.WriteJSON(w, http.StatusOK, result)
}

func (m *Module) handleOAuthStart(w http.ResponseWriter, r *http.Request) {
	providerID := strings.TrimSpace(chi.URLParam(r, "provider"))
	adapter, ok := m.GetAdapter(providerID)
	if !ok {
		httpserver.WriteError(w, r, http.StatusNotFound, "provider_not_found", "unknown provider")
		return
	}
	if adapter.OAuthConfig() == nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "oauth_not_supported", "provider does not support OAuth")
		return
	}

	var req oauthStartRequest
	if !decodeOptionalJSONBody(w, r, &req) {
		return
	}

	cfg, err := m.resolveOAuthConfig(providerID, req.ClientID, req.ClientSecret)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "oauth_not_supported", "provider does not support OAuth")
		return
	}
	if strings.TrimSpace(cfg.ClientID) == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request",
			fmt.Sprintf("oauth client_id is required; set %s_CLIENT_ID or provide client_id", oauthEnvPrefix(providerID)))
		return
	}

	redirectURI := strings.TrimSpace(req.RedirectURI)
	if redirectURI == "" {
		redirectURI = cfg.DefaultRedirectURI
	}
	redirectURI, err = validateHTTPURL(redirectURI)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "redirect_uri: "+err.Error())
		return
	}

	baseURL := strings.TrimSpace(req.BaseURL)
	if baseURL == "" {
		baseURL = adapter.DefaultBaseURL()
	}
	baseURL, err = validateHTTPURL(baseURL)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "base_url: "+err.Error())
		return
	}

	connName := strings.TrimSpace(req.ConnectionName)
	if connName == "" {
		connName = adapter.DisplayName()
	} else {
		connName, err = validateName(connName)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	}

	var models []string
	if req.Models != nil {
		models, err = validateStringSlice("models", req.Models)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	} else {
		models = adapter.DefaultModels()
	}

	verifier, challenge, err := generatePKCE()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate PKCE verifier")
		return
	}
	state, err := generateOAuthState()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate OAuth state")
		return
	}

	verifierEnc, err := m.vault.Encrypt([]byte(verifier))
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to encrypt PKCE verifier")
		return
	}
	var clientSecretEnc string
	if cfg.ClientSecret != "" {
		clientSecretEnc, err = m.vault.Encrypt([]byte(cfg.ClientSecret))
		if err != nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to encrypt OAuth client secret")
			return
		}
	}

	authorizeURL, err := buildAuthorizeURL(cfg, redirectURI, state, challenge)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to build authorization URL")
		return
	}

	now := time.Now().UTC()
	createdAt := now.Format(rfc3339NanoFixed)
	expiresAt := now.Add(oauthSessionTTL).Format(rfc3339NanoFixed)

	if err := m.repo.CreateOAuthSession(r.Context(), oauthSessionRecord{
		State:           state,
		ProviderID:      providerID,
		ConnectionName:  connName,
		CodeVerifierEnc: verifierEnc,
		RedirectURI:     redirectURI,
		ClientID:        cfg.ClientID,
		ClientSecretEnc: clientSecretEnc,
		BaseURL:         baseURL,
		Models:          models,
		ExpiresAt:       expiresAt,
		CreatedAt:       createdAt,
	}); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to persist OAuth session")
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, oauthStartResponse{
		State:        state,
		AuthorizeURL: authorizeURL,
		RedirectURI:  redirectURI,
		ExpiresAt:    expiresAt,
	})
}

func (m *Module) handleOAuthExchange(w http.ResponseWriter, r *http.Request) {
	providerID := strings.TrimSpace(chi.URLParam(r, "provider"))
	adapter, ok := m.GetAdapter(providerID)
	if !ok {
		httpserver.WriteError(w, r, http.StatusNotFound, "provider_not_found", "unknown provider")
		return
	}
	if adapter.OAuthConfig() == nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "oauth_not_supported", "provider does not support OAuth")
		return
	}

	var req oauthExchangeRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	code, state, oauthErr, oauthErrDesc, err := parseOAuthCallbackInput(req.CallbackURL, req.Code, req.State)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	if oauthErr != "" {
		msg := "oauth provider returned error: " + oauthErr
		if oauthErrDesc != "" {
			msg = msg + " (" + oauthErrDesc + ")"
		}
		httpserver.WriteError(w, r, http.StatusBadRequest, "oauth_error", truncateError(msg))
		return
	}

	var validatedConnName string
	if strings.TrimSpace(req.ConnectionName) != "" {
		validatedConnName, err = validateName(req.ConnectionName)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	}

	priority := defaultPriorityValue
	if req.Priority != nil {
		if *req.Priority < 0 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "priority must be non-negative")
			return
		}
		priority = *req.Priority
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	headers, err := validateHeaders(req.Headers)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	updateConnID := strings.TrimSpace(req.ConnectionID)
	if updateConnID != "" {
		existing, lookupErr := m.repo.GetConnectionSummary(r.Context(), updateConnID)
		if lookupErr != nil {
			if errors.Is(lookupErr, contracts.ErrConnectionNotFound) {
				httpserver.WriteError(w, r, http.StatusNotFound, "connection_not_found", "target provider connection not found")
				return
			}
			httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load target provider connection")
			return
		}
		if existing.ProviderID != providerID || existing.AuthType != authTypeOAuth {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_connection", "target connection must be an oauth connection for the same provider")
			return
		}
	}

	sess, err := m.repo.GetOAuthSession(r.Context(), state)
	if err != nil {
		if errors.Is(err, ErrOAuthSessionNotFound) {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_oauth_state", "oauth session not found or already consumed")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load oauth session")
		return
	}

	if sess.ProviderID != providerID {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_oauth_state", "oauth session provider mismatch")
		return
	}

	expTime, err := time.Parse(time.RFC3339Nano, sess.ExpiresAt)
	if err != nil || !time.Now().UTC().Before(expTime) {
		_ = m.repo.DeleteOAuthSession(r.Context(), state)
		httpserver.WriteError(w, r, http.StatusBadRequest, "oauth_session_expired", "oauth session has expired")
		return
	}

	verifierBytes, err := m.vault.Decrypt(sess.CodeVerifierEnc)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to decrypt PKCE verifier")
		return
	}
	var clientSecret string
	if sess.ClientSecretEnc != "" {
		secretBytes, decErr := m.vault.Decrypt(sess.ClientSecretEnc)
		if decErr != nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to decrypt OAuth client secret")
			return
		}
		clientSecret = string(secretBytes)
	}

	cfg, err := m.resolveOAuthConfig(providerID, sess.ClientID, clientSecret)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to resolve OAuth config")
		return
	}

	form := make(map[string][]string)
	form["grant_type"] = []string{"authorization_code"}
	form["code"] = []string{code}
	form["redirect_uri"] = []string{sess.RedirectURI}
	form["client_id"] = []string{sess.ClientID}
	form["code_verifier"] = []string{string(verifierBytes)}
	if clientSecret != "" {
		form["client_secret"] = []string{clientSecret}
	}
	for k, v := range cfg.ExtraTokenParams {
		if strings.TrimSpace(k) != "" && v != "" {
			form[k] = []string{v}
		}
	}

	tok, err := m.postTokenEndpoint(r.Context(), cfg.TokenURL, form)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadGateway, "oauth_exchange_failed", "upstream OAuth token exchange failed")
		return
	}

	stored := storedCredentials{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		TokenType:    tok.TokenType,
		Scope:        tok.Scope,
		ClientID:     sess.ClientID,
		ClientSecret: clientSecret,
		TokenURL:     cfg.TokenURL,
	}
	sealed, err := m.sealCredentials(stored)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to encrypt OAuth tokens")
		return
	}

	var tokenExpiresAt *string
	if tok.ExpiresIn > 0 {
		exp := time.Now().UTC().Add(time.Duration(tok.ExpiresIn) * time.Second).Format(rfc3339NanoFixed)
		tokenExpiresAt = &exp
	}

	connName := sess.ConnectionName
	if validatedConnName != "" {
		connName = validatedConnName
	}

	errOAuthConnectionMismatch := errors.New("providers: target connection mismatch")
	var createParams createConnectionParams
	var updateFn func(cur connectionRecord) (updateConnectionValues, error)

	if updateConnID == "" {
		newConnID, idErr := generateID(defaultConnectionPrefix)
		if idErr != nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate connection identifier")
			return
		}
		createParams = createConnectionParams{
			ID:             newConnID,
			ProviderID:     providerID,
			Name:           connName,
			Protocol:       adapter.Protocol(),
			BaseURL:        sess.BaseURL,
			AuthType:       authTypeOAuth,
			CredentialsEnc: sealed,
			Priority:       priority,
			Enabled:        enabled,
			Models:         sess.Models,
			Headers:        headers,
			TokenExpiresAt: tokenExpiresAt,
		}
	} else {
		updateFn = func(cur connectionRecord) (updateConnectionValues, error) {
			if cur.Summary.ProviderID != providerID || cur.Summary.AuthType != authTypeOAuth {
				return updateConnectionValues{}, errOAuthConnectionMismatch
			}
			nextName := cur.Summary.Name
			if validatedConnName != "" {
				nextName = validatedConnName
			}
			nextPriority := cur.Summary.Priority
			if req.Priority != nil {
				nextPriority = priority
			}
			nextEnabled := cur.Summary.Enabled
			if req.Enabled != nil {
				nextEnabled = enabled
			}
			nextHeaders := cur.Summary.Headers
			if req.Headers != nil {
				nextHeaders = headers
			}
			return updateConnectionValues{
				Name:           nextName,
				BaseURL:        cur.Summary.BaseURL,
				CredentialsEnc: sealed,
				Priority:       nextPriority,
				Enabled:        nextEnabled,
				Models:         cur.Summary.Models,
				Headers:        nextHeaders,
				TokenExpiresAt: tokenExpiresAt,
			}, nil
		}
	}

	summary, action, err := m.repo.CompleteOAuthExchange(r.Context(), state, updateConnID, createParams, updateFn)
	if err != nil {
		switch {
		case errors.Is(err, ErrOAuthSessionNotFound):
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_oauth_state", "oauth session not found or already consumed")
		case errors.Is(err, errOAuthConnectionMismatch):
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_connection", "target connection must be an oauth connection for the same provider")
		case errors.Is(err, contracts.ErrConnectionNotFound):
			httpserver.WriteError(w, r, http.StatusNotFound, "connection_not_found", "target provider connection not found")
		default:
			httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to persist OAuth connection")
		}
		return
	}

	m.publishChanged(action, "connection", summary.ID, summary.ProviderID, actorFromContext(r.Context()), summary.UpdatedAt)
	httpserver.WriteJSON(w, http.StatusCreated, summary)
}

func (m *Module) handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	code := q.Get("code")
	state := q.Get("state")
	oauthErr := q.Get("error")
	oauthErrDesc := q.Get("error_description")

	sessionStatus := "unverified"
	providerID := ""
	if state != "" {
		sess, err := m.repo.GetOAuthSession(r.Context(), state)
		if err == nil {
			providerID = sess.ProviderID
			if exp, parseErr := time.Parse(time.RFC3339Nano, sess.ExpiresAt); parseErr == nil && time.Now().UTC().Before(exp) {
				sessionStatus = "verified"
			} else {
				sessionStatus = "expired"
			}
		} else {
			sessionStatus = "unknown"
		}
	}
	if oauthErr != "" {
		sessionStatus = "error"
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	accept := strings.ToLower(r.Header.Get("Accept"))
	if strings.Contains(accept, "application/json") && !strings.Contains(accept, "text/html") {
		httpserver.WriteJSON(w, http.StatusOK, map[string]string{
			"status":            sessionStatus,
			"provider_id":       providerID,
			"state":             state,
			"code":              code,
			"error":             oauthErr,
			"error_description": oauthErrDesc,
		})
		return
	}

	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	if host == "" {
		host = "localhost:20128"
	}
	fullURL := fmt.Sprintf("%s://%s%s", scheme, host, r.URL.RequestURI())

	title := "OAuth Callback Received"
	if oauthErr != "" {
		title = "OAuth Authorization Error"
	}

	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"utf-8\"><title>")
	b.WriteString(html.EscapeString(title))
	b.WriteString("</title><style>body{font-family:system-ui,sans-serif;max-width:44rem;margin:2rem auto;padding:0 1rem;line-height:1.5}pre{background:#f4f4f5;padding:.75rem;overflow-x:auto;border-radius:4px}</style></head><body><h1>")
	b.WriteString(html.EscapeString(title))
	b.WriteString("</h1><p>Status: <strong>")
	b.WriteString(html.EscapeString(sessionStatus))
	b.WriteString("</strong>")
	if providerID != "" {
		b.WriteString(" (provider: <code>")
		b.WriteString(html.EscapeString(providerID))
		b.WriteString("</code>)")
	}
	b.WriteString("</p>")
	if oauthErr != "" {
		b.WriteString("<p>Error: <code>")
		b.WriteString(html.EscapeString(oauthErr))
		b.WriteString("</code></p>")
		if oauthErrDesc != "" {
			b.WriteString("<p>Description: ")
			b.WriteString(html.EscapeString(oauthErrDesc))
			b.WriteString("</p>")
		}
	} else {
		b.WriteString("<p>Copy the callback URL or parameters below into the router OAuth exchange step:</p>")
		b.WriteString("<h2>Callback URL</h2><pre>")
		b.WriteString(html.EscapeString(fullURL))
		b.WriteString("</pre><h2>State</h2><pre>")
		b.WriteString(html.EscapeString(state))
		b.WriteString("</pre><h2>Code</h2><pre>")
		b.WriteString(html.EscapeString(code))
		b.WriteString("</pre>")
	}
	b.WriteString("</body></html>")
	_, _ = io.WriteString(w, b.String())
}

func validateName(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("name must not be empty")
	}
	if utf8.RuneCountInString(trimmed) > maxNameRunes {
		return "", fmt.Errorf("name must not exceed %d characters", maxNameRunes)
	}
	return trimmed, nil
}

func validateSlugID(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("id must not be empty")
	}
	if utf8.RuneCountInString(trimmed) > maxIDRunes {
		return "", fmt.Errorf("id must not exceed %d characters", maxIDRunes)
	}
	for i := 0; i < len(trimmed); i++ {
		c := trimmed[i]
		isLower := c >= 'a' && c <= 'z'
		isDigit := c >= '0' && c <= '9'
		if !isLower && !isDigit && c != '-' && c != '_' {
			return "", errors.New("id must contain only lowercase letters, digits, hyphens, or underscores")
		}
	}
	return trimmed, nil
}

func validateStringSlice(field string, items []string) ([]string, error) {
	if len(items) == 0 {
		return []string{}, nil
	}
	if len(items) > maxSliceItems {
		return nil, fmt.Errorf("%s must not exceed %d items", field, maxSliceItems)
	}
	out := make([]string, 0, len(items))
	for _, raw := range items {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			return nil, fmt.Errorf("%s entries must not be empty", field)
		}
		if utf8.RuneCountInString(trimmed) > maxSliceItemRunes {
			return nil, fmt.Errorf("%s entry must not exceed %d characters", field, maxSliceItemRunes)
		}
		out = append(out, trimmed)
	}
	return out, nil
}

func validateHeaders(headers map[string]string) (map[string]string, error) {
	if len(headers) == 0 {
		return map[string]string{}, nil
	}
	if len(headers) > maxHeadersCount {
		return nil, fmt.Errorf("headers must not exceed %d entries", maxHeadersCount)
	}
	out := make(map[string]string, len(headers))
	for k, v := range headers {
		cleanKey := strings.TrimSpace(k)
		if cleanKey == "" {
			return nil, errors.New("header name must not be empty")
		}
		if utf8.RuneCountInString(cleanKey) > maxHeaderKeyRunes {
			return nil, fmt.Errorf("header name must not exceed %d characters", maxHeaderKeyRunes)
		}
		if strings.ContainsAny(cleanKey, "\r\n") || strings.ContainsAny(v, "\r\n") {
			return nil, errors.New("headers must not contain newline characters")
		}
		if utf8.RuneCountInString(v) > maxHeaderValueRunes {
			return nil, fmt.Errorf("header value must not exceed %d characters", maxHeaderValueRunes)
		}
		out[cleanKey] = v
	}
	return out, nil
}

func sanitizeUpstreamError(msg string, creds contracts.UpstreamCredentials) string {
	sanitized := msg
	if creds.APIKey != "" {
		sanitized = strings.ReplaceAll(sanitized, creds.APIKey, "[REDACTED]")
	}
	if creds.AccessToken != "" {
		sanitized = strings.ReplaceAll(sanitized, creds.AccessToken, "[REDACTED]")
	}
	return sanitized
}

func truncateError(msg string) string {
	trimmed := strings.TrimSpace(msg)
	if len(trimmed) <= maxTestErrorSummaryLen {
		return trimmed
	}
	return trimmed[:maxTestErrorSummaryLen]
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxProvidersBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			httpserver.WriteError(w, r, http.StatusRequestEntityTooLarge, "request_too_large", "request body exceeds maximum allowed size")
			return false
		}
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "invalid JSON request body")
		return false
	}
	if dec.More() {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "request body must contain a single JSON object")
		return false
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "request body must contain a single JSON object")
		return false
	}
	return true
}

func decodeOptionalJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxProvidersBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return true
		}
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			httpserver.WriteError(w, r, http.StatusRequestEntityTooLarge, "request_too_large", "request body exceeds maximum allowed size")
			return false
		}
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "invalid JSON request body")
		return false
	}
	if dec.More() {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "request body must contain a single JSON object")
		return false
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "request body must contain a single JSON object")
		return false
	}
	return true
}
