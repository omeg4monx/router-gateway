package providers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const rfc3339NanoFixed = "2006-01-02T15:04:05.000000000Z"

var (
	// ErrNodeNotFound indicates no custom provider node exists with the requested ID.
	ErrNodeNotFound = errors.New("providers: custom provider node not found")
	// ErrNodeAlreadyExists indicates a custom provider node ID is already registered.
	ErrNodeAlreadyExists = errors.New("providers: custom provider node already exists")
	// ErrOAuthSessionNotFound indicates no pending OAuth session exists for the state token.
	ErrOAuthSessionNotFound = errors.New("providers: oauth session not found")
)

// ProviderNode represents a user-defined compatible upstream provider template.
type ProviderNode struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Protocol  string            `json:"protocol"`
	BaseURL   string            `json:"base_url"`
	Models    []string          `json:"models"`
	Headers   map[string]string `json:"headers"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

// ConnectionSummary represents public metadata for a configured upstream
// provider connection. Secret credentials and ciphertext are intentionally excluded.
type ConnectionSummary struct {
	ID             string            `json:"id"`
	ProviderID     string            `json:"provider_id"`
	Name           string            `json:"name"`
	Protocol       string            `json:"protocol"`
	BaseURL        string            `json:"base_url"`
	AuthType       string            `json:"auth_type"`
	Priority       int               `json:"priority"`
	Enabled        bool              `json:"enabled"`
	Models         []string          `json:"models"`
	Headers        map[string]string `json:"headers"`
	HasCredentials bool              `json:"has_credentials"`
	TokenExpiresAt *string           `json:"token_expires_at"`
	LastTestedAt   *string           `json:"last_tested_at"`
	LastTestStatus string            `json:"last_test_status"`
	LastTestError  string            `json:"last_test_error"`
	CreatedAt      string            `json:"created_at"`
	UpdatedAt      string            `json:"updated_at"`
}

type storedCredentials struct {
	APIKey       string         `json:"api_key,omitempty"`
	AccessToken  string         `json:"access_token,omitempty"`
	RefreshToken string         `json:"refresh_token,omitempty"`
	IDToken      string         `json:"id_token,omitempty"`
	TokenType    string         `json:"token_type,omitempty"`
	Scope        string         `json:"scope,omitempty"`
	ClientID     string         `json:"client_id,omitempty"`
	ClientSecret string         `json:"client_secret,omitempty"`
	TokenURL     string         `json:"token_url,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

type connectionRecord struct {
	Summary        ConnectionSummary
	CredentialsEnc string
}

type createConnectionParams struct {
	ID             string
	ProviderID     string
	Name           string
	Protocol       string
	BaseURL        string
	AuthType       string
	CredentialsEnc string
	Priority       int
	Enabled        bool
	Models         []string
	Headers        map[string]string
	TokenExpiresAt *string
}

type updateConnectionValues struct {
	Name           string
	BaseURL        string
	CredentialsEnc string
	Priority       int
	Enabled        bool
	Models         []string
	Headers        map[string]string
	TokenExpiresAt *string
}

type oauthSessionRecord struct {
	State           string
	ProviderID      string
	ConnectionName  string
	CodeVerifierEnc string
	RedirectURI     string
	ClientID        string
	ClientSecretEnc string
	BaseURL         string
	Models          []string
	ExpiresAt       string
	CreatedAt       string
}

type repository struct {
	db       *db.DB
	onChange func()
}

func newRepository(database *db.DB, onChange func()) *repository {
	return &repository{db: database, onChange: onChange}
}

func (r *repository) notifyChanged() {
	if r != nil && r.onChange != nil {
		r.onChange()
	}
}

type rowScanner interface {
	Scan(dest ...any) error
}

const selectNodeColumns = `id, name, protocol, base_url, models_json, headers_json, created_at, updated_at`

const selectConnectionColumns = `id, provider_id, name, protocol, base_url, auth_type, credentials_enc,
	priority, enabled, models_json, headers_json, token_expires_at,
	last_tested_at, last_test_status, last_test_error, created_at, updated_at`

func scanNode(sc rowScanner) (ProviderNode, error) {
	var (
		node        ProviderNode
		modelsJSON  string
		headersJSON string
	)
	if err := sc.Scan(
		&node.ID,
		&node.Name,
		&node.Protocol,
		&node.BaseURL,
		&modelsJSON,
		&headersJSON,
		&node.CreatedAt,
		&node.UpdatedAt,
	); err != nil {
		return ProviderNode{}, err
	}
	var err error
	if node.Models, err = decodeStringSlice(modelsJSON); err != nil {
		return ProviderNode{}, fmt.Errorf("providers: decode node models for %q: %w", node.ID, err)
	}
	if node.Headers, err = decodeStringMap(headersJSON); err != nil {
		return ProviderNode{}, fmt.Errorf("providers: decode node headers for %q: %w", node.ID, err)
	}
	return node, nil
}

func scanConnectionRecord(sc rowScanner) (connectionRecord, error) {
	var (
		rec            connectionRecord
		enabledInt     int
		modelsJSON     string
		headersJSON    string
		tokenExpiresAt sql.NullString
		lastTestedAt   sql.NullString
	)
	if err := sc.Scan(
		&rec.Summary.ID,
		&rec.Summary.ProviderID,
		&rec.Summary.Name,
		&rec.Summary.Protocol,
		&rec.Summary.BaseURL,
		&rec.Summary.AuthType,
		&rec.CredentialsEnc,
		&rec.Summary.Priority,
		&enabledInt,
		&modelsJSON,
		&headersJSON,
		&tokenExpiresAt,
		&lastTestedAt,
		&rec.Summary.LastTestStatus,
		&rec.Summary.LastTestError,
		&rec.Summary.CreatedAt,
		&rec.Summary.UpdatedAt,
	); err != nil {
		return connectionRecord{}, err
	}

	var err error
	if rec.Summary.Models, err = decodeStringSlice(modelsJSON); err != nil {
		return connectionRecord{}, fmt.Errorf("providers: decode connection models for %q: %w", rec.Summary.ID, err)
	}
	if rec.Summary.Headers, err = decodeStringMap(headersJSON); err != nil {
		return connectionRecord{}, fmt.Errorf("providers: decode connection headers for %q: %w", rec.Summary.ID, err)
	}
	rec.Summary.Enabled = enabledInt != 0
	rec.Summary.HasCredentials = rec.CredentialsEnc != ""
	if tokenExpiresAt.Valid && tokenExpiresAt.String != "" {
		s := tokenExpiresAt.String
		rec.Summary.TokenExpiresAt = &s
	}
	if lastTestedAt.Valid && lastTestedAt.String != "" {
		s := lastTestedAt.String
		rec.Summary.LastTestedAt = &s
	}
	return rec, nil
}

func (r *repository) ListNodes(ctx context.Context) ([]ProviderNode, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("providers: database is not initialized")
	}

	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT `+selectNodeColumns+` FROM provider_nodes ORDER BY created_at ASC, id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("providers: list provider_nodes: %w", err)
	}
	defer rows.Close()

	nodes := make([]ProviderNode, 0)
	for rows.Next() {
		node, err := scanNode(rows)
		if err != nil {
			return nil, fmt.Errorf("providers: scan provider_node: %w", err)
		}
		nodes = append(nodes, node)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("providers: iterate provider_nodes: %w", err)
	}
	return nodes, nil
}

func (r *repository) GetNode(ctx context.Context, id string) (ProviderNode, error) {
	if r == nil || r.db == nil {
		return ProviderNode{}, errors.New("providers: database is not initialized")
	}

	row := r.db.Reader().QueryRowContext(
		ctx,
		`SELECT `+selectNodeColumns+` FROM provider_nodes WHERE id = ?`,
		id,
	)
	node, err := scanNode(row)
	if errors.Is(err, sql.ErrNoRows) {
		return ProviderNode{}, ErrNodeNotFound
	}
	if err != nil {
		return ProviderNode{}, fmt.Errorf("providers: get provider_node %q: %w", id, err)
	}
	return node, nil
}

func (r *repository) CreateNode(ctx context.Context, node ProviderNode) (ProviderNode, error) {
	if r == nil || r.db == nil {
		return ProviderNode{}, errors.New("providers: database is not initialized")
	}

	modelsJSON, err := encodeStringSlice(node.Models)
	if err != nil {
		return ProviderNode{}, fmt.Errorf("providers: encode node models: %w", err)
	}
	headersJSON, err := encodeStringMap(node.Headers)
	if err != nil {
		return ProviderNode{}, fmt.Errorf("providers: encode node headers: %w", err)
	}

	now := time.Now().UTC().Format(rfc3339NanoFixed)
	_, err = r.db.Writer().ExecContext(
		ctx,
		`INSERT INTO provider_nodes (id, name, protocol, base_url, models_json, headers_json, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		node.ID,
		node.Name,
		node.Protocol,
		node.BaseURL,
		modelsJSON,
		headersJSON,
		now,
		now,
	)
	if err != nil {
		if isUniqueConstraintError(err) {
			return ProviderNode{}, ErrNodeAlreadyExists
		}
		return ProviderNode{}, fmt.Errorf("providers: insert provider_node: %w", err)
	}

	node.Models = cloneStringSlice(node.Models)
	node.Headers = cloneStringMap(node.Headers)
	node.CreatedAt = now
	node.UpdatedAt = now
	r.notifyChanged()
	return node, nil
}

func (r *repository) DeleteNode(ctx context.Context, id string) (string, error) {
	if r == nil || r.db == nil {
		return "", errors.New("providers: database is not initialized")
	}

	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM provider_nodes WHERE id = ?`, id)
	if err != nil {
		return "", fmt.Errorf("providers: delete provider_node %q: %w", id, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("providers: delete provider_node rows affected: %w", err)
	}
	if affected == 0 {
		return "", ErrNodeNotFound
	}
	r.notifyChanged()
	return time.Now().UTC().Format(rfc3339NanoFixed), nil
}

func (r *repository) ListConnections(ctx context.Context) ([]ConnectionSummary, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("providers: database is not initialized")
	}

	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT `+selectConnectionColumns+` FROM provider_connections ORDER BY priority ASC, created_at ASC, id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("providers: list provider_connections: %w", err)
	}
	defer rows.Close()

	items := make([]ConnectionSummary, 0)
	for rows.Next() {
		rec, err := scanConnectionRecord(rows)
		if err != nil {
			return nil, fmt.Errorf("providers: scan provider_connection: %w", err)
		}
		items = append(items, rec.Summary)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("providers: iterate provider_connections: %w", err)
	}
	return items, nil
}

func (r *repository) ListEnabledConnections(ctx context.Context) ([]contracts.ProviderConnection, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("providers: database is not initialized")
	}

	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT id, provider_id, name, protocol, base_url, auth_type, priority, enabled, models_json, headers_json
		 FROM provider_connections
		 WHERE enabled = 1
		 ORDER BY priority ASC, created_at ASC, id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("providers: list enabled connections: %w", err)
	}
	defer rows.Close()

	out := make([]contracts.ProviderConnection, 0)
	for rows.Next() {
		var (
			conn        contracts.ProviderConnection
			enabledInt  int
			modelsJSON  string
			headersJSON string
		)
		if err := rows.Scan(
			&conn.ID,
			&conn.ProviderID,
			&conn.Name,
			&conn.Protocol,
			&conn.BaseURL,
			&conn.AuthType,
			&conn.Priority,
			&enabledInt,
			&modelsJSON,
			&headersJSON,
		); err != nil {
			return nil, fmt.Errorf("providers: scan enabled connection: %w", err)
		}
		conn.Enabled = enabledInt != 0
		if conn.Models, err = decodeStringSlice(modelsJSON); err != nil {
			return nil, fmt.Errorf("providers: decode enabled connection models: %w", err)
		}
		if conn.ExtraHeaders, err = decodeStringMap(headersJSON); err != nil {
			return nil, fmt.Errorf("providers: decode enabled connection headers: %w", err)
		}
		out = append(out, conn)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("providers: iterate enabled connections: %w", err)
	}
	return out, nil
}

func (r *repository) GetConnectionSummary(ctx context.Context, id string) (ConnectionSummary, error) {
	rec, err := r.GetConnectionRecord(ctx, id)
	if err != nil {
		return ConnectionSummary{}, err
	}
	return rec.Summary, nil
}

func (r *repository) GetConnectionRecord(ctx context.Context, id string) (connectionRecord, error) {
	if r == nil || r.db == nil {
		return connectionRecord{}, errors.New("providers: database is not initialized")
	}

	row := r.db.Reader().QueryRowContext(
		ctx,
		`SELECT `+selectConnectionColumns+` FROM provider_connections WHERE id = ?`,
		id,
	)
	rec, err := scanConnectionRecord(row)
	if errors.Is(err, sql.ErrNoRows) {
		return connectionRecord{}, contracts.ErrConnectionNotFound
	}
	if err != nil {
		return connectionRecord{}, fmt.Errorf("providers: get provider_connection %q: %w", id, err)
	}
	return rec, nil
}

func (r *repository) CreateConnection(ctx context.Context, p createConnectionParams) (ConnectionSummary, error) {
	if r == nil || r.db == nil {
		return ConnectionSummary{}, errors.New("providers: database is not initialized")
	}

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return ConnectionSummary{}, fmt.Errorf("providers: begin create connection tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	summary, err := insertConnectionTx(ctx, tx, p)
	if err != nil {
		return ConnectionSummary{}, err
	}
	if err := tx.Commit(); err != nil {
		return ConnectionSummary{}, fmt.Errorf("providers: commit create connection tx: %w", err)
	}
	r.notifyChanged()
	return summary, nil
}

func insertConnectionTx(ctx context.Context, tx *sql.Tx, p createConnectionParams) (ConnectionSummary, error) {
	modelsJSON, err := encodeStringSlice(p.Models)
	if err != nil {
		return ConnectionSummary{}, fmt.Errorf("providers: encode connection models: %w", err)
	}
	headersJSON, err := encodeStringMap(p.Headers)
	if err != nil {
		return ConnectionSummary{}, fmt.Errorf("providers: encode connection headers: %w", err)
	}

	enabledInt := 0
	if p.Enabled {
		enabledInt = 1
	}
	var tokenExpVal any
	var tokenExpCopy *string
	if p.TokenExpiresAt != nil && *p.TokenExpiresAt != "" {
		s := *p.TokenExpiresAt
		tokenExpVal = s
		tokenExpCopy = &s
	}

	now := time.Now().UTC().Format(rfc3339NanoFixed)
	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO provider_connections (
			id, provider_id, name, protocol, base_url, auth_type, credentials_enc,
			priority, enabled, models_json, headers_json, token_expires_at,
			last_tested_at, last_test_status, last_test_error, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL, '', '', ?, ?)`,
		p.ID,
		p.ProviderID,
		p.Name,
		p.Protocol,
		p.BaseURL,
		p.AuthType,
		p.CredentialsEnc,
		p.Priority,
		enabledInt,
		modelsJSON,
		headersJSON,
		tokenExpVal,
		now,
		now,
	); err != nil {
		return ConnectionSummary{}, fmt.Errorf("providers: insert provider_connection: %w", err)
	}

	created := ConnectionSummary{
		ID:             p.ID,
		ProviderID:     p.ProviderID,
		Name:           p.Name,
		Protocol:       p.Protocol,
		BaseURL:        p.BaseURL,
		AuthType:       p.AuthType,
		Priority:       p.Priority,
		Enabled:        p.Enabled,
		Models:         cloneStringSlice(p.Models),
		Headers:        cloneStringMap(p.Headers),
		HasCredentials: p.CredentialsEnc != "",
		TokenExpiresAt: tokenExpCopy,
		LastTestedAt:   nil,
		LastTestStatus: "",
		LastTestError:  "",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	return created, nil
}

func (r *repository) UpdateConnection(ctx context.Context, id string, fn func(cur connectionRecord) (updateConnectionValues, error)) (ConnectionSummary, error) {
	if r == nil || r.db == nil {
		return ConnectionSummary{}, errors.New("providers: database is not initialized")
	}

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return ConnectionSummary{}, fmt.Errorf("providers: begin update connection tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	summary, err := updateConnectionTx(ctx, tx, id, fn)
	if err != nil {
		return ConnectionSummary{}, err
	}
	if err := tx.Commit(); err != nil {
		return ConnectionSummary{}, fmt.Errorf("providers: commit update connection tx: %w", err)
	}
	r.notifyChanged()
	return summary, nil
}

func updateConnectionTx(ctx context.Context, tx *sql.Tx, id string, fn func(cur connectionRecord) (updateConnectionValues, error)) (ConnectionSummary, error) {
	cur, err := scanConnectionRecord(tx.QueryRowContext(
		ctx,
		`SELECT `+selectConnectionColumns+` FROM provider_connections WHERE id = ?`,
		id,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return ConnectionSummary{}, contracts.ErrConnectionNotFound
	}
	if err != nil {
		return ConnectionSummary{}, fmt.Errorf("providers: read provider_connection %q for update: %w", id, err)
	}

	next, err := fn(cur)
	if err != nil {
		return ConnectionSummary{}, err
	}

	modelsJSON, err := encodeStringSlice(next.Models)
	if err != nil {
		return ConnectionSummary{}, fmt.Errorf("providers: encode updated connection models: %w", err)
	}
	headersJSON, err := encodeStringMap(next.Headers)
	if err != nil {
		return ConnectionSummary{}, fmt.Errorf("providers: encode updated connection headers: %w", err)
	}

	enabledInt := 0
	if next.Enabled {
		enabledInt = 1
	}
	var tokenExpVal any
	var tokenExpCopy *string
	if next.TokenExpiresAt != nil && *next.TokenExpiresAt != "" {
		s := *next.TokenExpiresAt
		tokenExpVal = s
		tokenExpCopy = &s
	}

	now := time.Now().UTC().Format(rfc3339NanoFixed)
	if _, err := tx.ExecContext(
		ctx,
		`UPDATE provider_connections
		 SET name = ?,
		     base_url = ?,
		     credentials_enc = ?,
		     priority = ?,
		     enabled = ?,
		     models_json = ?,
		     headers_json = ?,
		     token_expires_at = ?,
		     updated_at = ?
		 WHERE id = ?`,
		next.Name,
		next.BaseURL,
		next.CredentialsEnc,
		next.Priority,
		enabledInt,
		modelsJSON,
		headersJSON,
		tokenExpVal,
		now,
		id,
	); err != nil {
		return ConnectionSummary{}, fmt.Errorf("providers: update provider_connection %q: %w", id, err)
	}

	cur.Summary.Name = next.Name
	cur.Summary.BaseURL = next.BaseURL
	cur.Summary.Priority = next.Priority
	cur.Summary.Enabled = next.Enabled
	cur.Summary.Models = cloneStringSlice(next.Models)
	cur.Summary.Headers = cloneStringMap(next.Headers)
	cur.Summary.HasCredentials = next.CredentialsEnc != ""
	cur.Summary.TokenExpiresAt = tokenExpCopy
	cur.Summary.UpdatedAt = now
	return cur.Summary, nil
}

func (r *repository) DeleteConnection(ctx context.Context, id string) (providerID, deletedAt string, err error) {
	if r == nil || r.db == nil {
		return "", "", errors.New("providers: database is not initialized")
	}

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return "", "", fmt.Errorf("providers: begin delete connection tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	err = tx.QueryRowContext(ctx, `SELECT provider_id FROM provider_connections WHERE id = ?`, id).Scan(&providerID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", contracts.ErrConnectionNotFound
	}
	if err != nil {
		return "", "", fmt.Errorf("providers: read provider_connection %q for delete: %w", id, err)
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM provider_connections WHERE id = ?`, id); err != nil {
		return "", "", fmt.Errorf("providers: delete provider_connection %q: %w", id, err)
	}
	if err := tx.Commit(); err != nil {
		return "", "", fmt.Errorf("providers: commit delete connection tx: %w", err)
	}
	r.notifyChanged()
	return providerID, time.Now().UTC().Format(rfc3339NanoFixed), nil
}

func (r *repository) UpdateTestResult(ctx context.Context, id, testedAt, status, testErr string) error {
	if r == nil || r.db == nil {
		return errors.New("providers: database is not initialized")
	}

	res, err := r.db.Writer().ExecContext(
		ctx,
		`UPDATE provider_connections
		 SET last_tested_at = ?, last_test_status = ?, last_test_error = ?, updated_at = ?
		 WHERE id = ?`,
		testedAt,
		status,
		testErr,
		testedAt,
		id,
	)
	if err != nil {
		return fmt.Errorf("providers: update test result for %q: %w", id, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("providers: update test result rows affected: %w", err)
	}
	if affected == 0 {
		return contracts.ErrConnectionNotFound
	}
	return nil
}

func (r *repository) UpdateOAuthCredentials(ctx context.Context, id, credsEnc string, tokenExpiresAt *string) (string, error) {
	if r == nil || r.db == nil {
		return "", errors.New("providers: database is not initialized")
	}

	var tokenExpVal any
	if tokenExpiresAt != nil && *tokenExpiresAt != "" {
		tokenExpVal = *tokenExpiresAt
	}
	now := time.Now().UTC().Format(rfc3339NanoFixed)
	res, err := r.db.Writer().ExecContext(
		ctx,
		`UPDATE provider_connections
		 SET credentials_enc = ?, token_expires_at = ?, updated_at = ?
		 WHERE id = ?`,
		credsEnc,
		tokenExpVal,
		now,
		id,
	)
	if err != nil {
		return "", fmt.Errorf("providers: update oauth credentials for %q: %w", id, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("providers: update oauth credentials rows affected: %w", err)
	}
	if affected == 0 {
		return "", contracts.ErrConnectionNotFound
	}
	r.notifyChanged()
	return now, nil
}

func (r *repository) CreateOAuthSession(ctx context.Context, sess oauthSessionRecord) error {
	if r == nil || r.db == nil {
		return errors.New("providers: database is not initialized")
	}

	modelsJSON, err := encodeStringSlice(sess.Models)
	if err != nil {
		return fmt.Errorf("providers: encode oauth session models: %w", err)
	}

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("providers: begin create oauth session tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	now := time.Now().UTC().Format(rfc3339NanoFixed)
	if _, err := tx.ExecContext(ctx, `DELETE FROM oauth_sessions WHERE expires_at <= ?`, now); err != nil {
		return fmt.Errorf("providers: purge expired oauth sessions: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO oauth_sessions (
			state, provider_id, connection_name, code_verifier_enc, redirect_uri,
			client_id, client_secret_enc, base_url, models_json, expires_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		sess.State,
		sess.ProviderID,
		sess.ConnectionName,
		sess.CodeVerifierEnc,
		sess.RedirectURI,
		sess.ClientID,
		sess.ClientSecretEnc,
		sess.BaseURL,
		modelsJSON,
		sess.ExpiresAt,
		sess.CreatedAt,
	); err != nil {
		return fmt.Errorf("providers: insert oauth session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("providers: commit create oauth session tx: %w", err)
	}
	return nil
}

func (r *repository) GetOAuthSession(ctx context.Context, state string) (oauthSessionRecord, error) {
	if r == nil || r.db == nil {
		return oauthSessionRecord{}, errors.New("providers: database is not initialized")
	}

	var (
		sess       oauthSessionRecord
		modelsJSON string
	)
	err := r.db.Reader().QueryRowContext(
		ctx,
		`SELECT state, provider_id, connection_name, code_verifier_enc, redirect_uri,
		        client_id, client_secret_enc, base_url, models_json, expires_at, created_at
		 FROM oauth_sessions
		 WHERE state = ?`,
		state,
	).Scan(
		&sess.State,
		&sess.ProviderID,
		&sess.ConnectionName,
		&sess.CodeVerifierEnc,
		&sess.RedirectURI,
		&sess.ClientID,
		&sess.ClientSecretEnc,
		&sess.BaseURL,
		&modelsJSON,
		&sess.ExpiresAt,
		&sess.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return oauthSessionRecord{}, ErrOAuthSessionNotFound
	}
	if err != nil {
		return oauthSessionRecord{}, fmt.Errorf("providers: get oauth session: %w", err)
	}
	if sess.Models, err = decodeStringSlice(modelsJSON); err != nil {
		return oauthSessionRecord{}, fmt.Errorf("providers: decode oauth session models: %w", err)
	}
	return sess, nil
}

func (r *repository) DeleteOAuthSession(ctx context.Context, state string) error {
	if r == nil || r.db == nil {
		return errors.New("providers: database is not initialized")
	}
	_, err := r.db.Writer().ExecContext(ctx, `DELETE FROM oauth_sessions WHERE state = ?`, state)
	if err != nil {
		return fmt.Errorf("providers: delete oauth session: %w", err)
	}
	return nil
}

func (r *repository) CompleteOAuthExchange(
	ctx context.Context,
	state string,
	updateConnID string,
	createParams createConnectionParams,
	updateFn func(cur connectionRecord) (updateConnectionValues, error),
) (summary ConnectionSummary, action string, err error) {
	if r == nil || r.db == nil {
		return ConnectionSummary{}, "", errors.New("providers: database is not initialized")
	}

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return ConnectionSummary{}, "", fmt.Errorf("providers: begin oauth exchange tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	delRes, err := tx.ExecContext(ctx, `DELETE FROM oauth_sessions WHERE state = ?`, state)
	if err != nil {
		return ConnectionSummary{}, "", fmt.Errorf("providers: consume oauth session: %w", err)
	}
	affected, err := delRes.RowsAffected()
	if err != nil {
		return ConnectionSummary{}, "", fmt.Errorf("providers: oauth session rows affected: %w", err)
	}
	if affected == 0 {
		return ConnectionSummary{}, "", ErrOAuthSessionNotFound
	}

	if updateConnID != "" {
		summary, err = updateConnectionTx(ctx, tx, updateConnID, updateFn)
		if err != nil {
			return ConnectionSummary{}, "", err
		}
		action = "updated"
	} else {
		summary, err = insertConnectionTx(ctx, tx, createParams)
		if err != nil {
			return ConnectionSummary{}, "", err
		}
		action = "created"
	}

	if err := tx.Commit(); err != nil {
		return ConnectionSummary{}, "", fmt.Errorf("providers: commit oauth exchange tx: %w", err)
	}
	r.notifyChanged()
	return summary, action, nil
}

func encodeStringSlice(items []string) (string, error) {
	data, err := json.Marshal(cloneStringSlice(items))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func decodeStringSlice(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return []string{}, nil
	}
	var items []string
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil, err
	}
	if items == nil {
		return []string{}, nil
	}
	return items, nil
}

func encodeStringMap(m map[string]string) (string, error) {
	data, err := json.Marshal(cloneStringMap(m))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func decodeStringMap(raw string) (map[string]string, error) {
	if strings.TrimSpace(raw) == "" {
		return map[string]string{}, nil
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, err
	}
	if m == nil {
		return map[string]string{}, nil
	}
	return m, nil
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint failed") || strings.Contains(msg, "constraint_primarykey")
}
