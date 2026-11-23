package vouchers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const timestampLayout = "2006-01-02T15:04:05.000000000Z"

var (
	errVoucherNotFound   = errors.New("vouchers: voucher not found")
	errVoucherCodeExists = errors.New("vouchers: voucher code already exists")
	errVoucherDisabled   = errors.New("vouchers: voucher is disabled")
	errVoucherExpired    = errors.New("vouchers: voucher has expired")
	errVoucherExhausted  = errors.New("vouchers: voucher redemption limit reached")
	errAlreadyClaimed    = errors.New("vouchers: claimer has already redeemed this voucher")
)

// Voucher represents a redeemable API-key issuance code.
type Voucher struct {
	ID             string   `json:"id"`
	Code           string   `json:"code"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	QuotaUSD       float64  `json:"quota_usd"`
	RateLimitRPM   int      `json:"rate_limit_rpm"`
	AllowedModels  []string `json:"allowed_models"`
	TeamID         string   `json:"team_id"`
	MaxRedemptions int      `json:"max_redemptions"`
	RedeemedCount  int      `json:"redeemed_count"`
	Enabled        bool     `json:"enabled"`
	ExpiresAt      string   `json:"expires_at"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

// Redemption records a successful voucher claim.
type Redemption struct {
	ID         string `json:"id"`
	VoucherID  string `json:"voucher_id"`
	Code       string `json:"code"`
	ClaimerID  string `json:"claimer_id"`
	KeyID      string `json:"key_id"`
	IPAddress  string `json:"ip_address"`
	RedeemedAt string `json:"redeemed_at"`
}

type createVoucherInput struct {
	Code           string
	Name           string
	Description    string
	QuotaUSD       float64
	RateLimitRPM   int
	AllowedModels  []string
	TeamID         string
	MaxRedemptions int
	Enabled        bool
	ExpiresAt      string
}

type patchVoucherInput struct {
	Name           *string
	Description    *string
	MaxRedemptions *int
	Enabled        *bool
	ExpiresAt      *string
}

type repo struct {
	db *db.DB
}

func newRepo(database *db.DB) *repo {
	return &repo{db: database}
}

func formatUTC(t time.Time) string {
	return t.UTC().Format(timestampLayout)
}

func newID(prefix string) string {
	var buf [10]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return fmt.Sprintf("%s_%d", prefix, time.Now().UTC().UnixNano())
	}
	return prefix + "_" + hex.EncodeToString(buf[:])
}

func generateCode() string {
	var buf [6]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return fmt.Sprintf("OMG-%d", time.Now().UTC().UnixNano()%1_000_000)
	}
	raw := strings.ToUpper(hex.EncodeToString(buf[:]))
	return "OMG-" + raw[:4] + "-" + raw[4:8] + "-" + raw[8:12]
}

func normalizeCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func (r *repo) list(ctx context.Context) ([]Voucher, error) {
	rows, err := r.db.Reader().QueryContext(ctx,
		`SELECT id, code, name, description, quota_usd, rate_limit_rpm, allowed_models_json,
		        team_id, max_redemptions, redeemed_count, enabled, expires_at, created_at, updated_at
		 FROM vouchers ORDER BY created_at DESC, id DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list vouchers: %w", err)
	}
	defer func() { _ = rows.Close() }()

	items := make([]Voucher, 0)
	for rows.Next() {
		v, err := scanVoucher(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	return items, rows.Err()
}

func (r *repo) getByID(ctx context.Context, id string) (Voucher, error) {
	row := r.db.Reader().QueryRowContext(ctx,
		`SELECT id, code, name, description, quota_usd, rate_limit_rpm, allowed_models_json,
		        team_id, max_redemptions, redeemed_count, enabled, expires_at, created_at, updated_at
		 FROM vouchers WHERE id = ?`, id,
	)
	return scanVoucherRow(row)
}

func (r *repo) create(ctx context.Context, in createVoucherInput) (Voucher, error) {
	now := formatUTC(time.Now())
	code := normalizeCode(in.Code)
	if code == "" {
		code = generateCode()
	}
	allowed := in.AllowedModels
	if allowed == nil {
		allowed = []string{}
	}
	allowedJSON, _ := json.Marshal(allowed)

	v := Voucher{
		ID:             newID("vch"),
		Code:           code,
		Name:           strings.TrimSpace(in.Name),
		Description:    strings.TrimSpace(in.Description),
		QuotaUSD:       in.QuotaUSD,
		RateLimitRPM:   in.RateLimitRPM,
		AllowedModels:  allowed,
		TeamID:         strings.TrimSpace(in.TeamID),
		MaxRedemptions: in.MaxRedemptions,
		RedeemedCount:  0,
		Enabled:        in.Enabled,
		ExpiresAt:      strings.TrimSpace(in.ExpiresAt),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	enabledInt := 0
	if v.Enabled {
		enabledInt = 1
	}

	_, err := r.db.Writer().ExecContext(ctx,
		`INSERT INTO vouchers
		 (id, code, name, description, quota_usd, rate_limit_rpm, allowed_models_json,
		  team_id, max_redemptions, redeemed_count, enabled, expires_at, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		v.ID, v.Code, v.Name, v.Description, v.QuotaUSD, v.RateLimitRPM, string(allowedJSON),
		v.TeamID, v.MaxRedemptions, v.RedeemedCount, enabledInt, v.ExpiresAt, v.CreatedAt, v.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return Voucher{}, errVoucherCodeExists
		}
		return Voucher{}, fmt.Errorf("insert voucher: %w", err)
	}
	return v, nil
}

func (r *repo) patch(ctx context.Context, id string, in patchVoucherInput) (Voucher, error) {
	cur, err := r.getByID(ctx, id)
	if err != nil {
		return Voucher{}, err
	}
	if in.Name != nil {
		cur.Name = strings.TrimSpace(*in.Name)
	}
	if in.Description != nil {
		cur.Description = strings.TrimSpace(*in.Description)
	}
	if in.MaxRedemptions != nil {
		cur.MaxRedemptions = *in.MaxRedemptions
	}
	if in.Enabled != nil {
		cur.Enabled = *in.Enabled
	}
	if in.ExpiresAt != nil {
		cur.ExpiresAt = strings.TrimSpace(*in.ExpiresAt)
	}
	cur.UpdatedAt = formatUTC(time.Now())

	enabledInt := 0
	if cur.Enabled {
		enabledInt = 1
	}
	_, err = r.db.Writer().ExecContext(ctx,
		`UPDATE vouchers SET name = ?, description = ?, max_redemptions = ?, enabled = ?, expires_at = ?, updated_at = ?
		 WHERE id = ?`,
		cur.Name, cur.Description, cur.MaxRedemptions, enabledInt, cur.ExpiresAt, cur.UpdatedAt, id,
	)
	if err != nil {
		return Voucher{}, fmt.Errorf("update voucher: %w", err)
	}
	return r.getByID(ctx, id)
}

func (r *repo) deleteByID(ctx context.Context, id string) error {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM vouchers WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete voucher: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errVoucherNotFound
	}
	return nil
}

func (r *repo) listRedemptions(ctx context.Context, voucherID string) ([]Redemption, error) {
	if _, err := r.getByID(ctx, voucherID); err != nil {
		return nil, err
	}
	rows, err := r.db.Reader().QueryContext(ctx,
		`SELECT id, voucher_id, code, claimer_id, key_id, ip_address, redeemed_at
		 FROM voucher_redemptions WHERE voucher_id = ? ORDER BY redeemed_at DESC, id DESC`,
		voucherID,
	)
	if err != nil {
		return nil, fmt.Errorf("list voucher redemptions: %w", err)
	}
	defer func() { _ = rows.Close() }()

	items := make([]Redemption, 0)
	for rows.Next() {
		var item Redemption
		if err := rows.Scan(
			&item.ID, &item.VoucherID, &item.Code, &item.ClaimerID,
			&item.KeyID, &item.IPAddress, &item.RedeemedAt,
		); err != nil {
			return nil, fmt.Errorf("scan redemption: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// reserveClaim atomically checks voucher validity, enforces per-actor uniqueness,
// and increments redeemed_count inside a single SQLite write transaction.
func (r *repo) reserveClaim(ctx context.Context, code, claimerID, ipAddress string) (Voucher, Redemption, error) {
	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return Voucher{}, Redemption{}, fmt.Errorf("begin claim tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	row := tx.QueryRowContext(ctx,
		`SELECT id, code, name, description, quota_usd, rate_limit_rpm, allowed_models_json,
		        team_id, max_redemptions, redeemed_count, enabled, expires_at, created_at, updated_at
		 FROM vouchers WHERE code = ?`, normalizeCode(code),
	)
	v, err := scanVoucherRow(row)
	if err != nil {
		return Voucher{}, Redemption{}, err
	}
	if !v.Enabled {
		return Voucher{}, Redemption{}, errVoucherDisabled
	}
	now := time.Now().UTC()
	if v.ExpiresAt != "" {
		exp, parseErr := time.Parse(time.RFC3339, v.ExpiresAt)
		if parseErr == nil && now.After(exp) {
			return Voucher{}, Redemption{}, errVoucherExpired
		}
	}
	if v.RedeemedCount >= v.MaxRedemptions {
		return Voucher{}, Redemption{}, errVoucherExhausted
	}

	nowStr := formatUTC(now)
	red := Redemption{
		ID:         newID("red"),
		VoucherID:  v.ID,
		Code:       v.Code,
		ClaimerID:  strings.TrimSpace(claimerID),
		IPAddress:  strings.TrimSpace(ipAddress),
		RedeemedAt: nowStr,
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO voucher_redemptions (id, voucher_id, code, claimer_id, key_id, ip_address, redeemed_at)
		 VALUES (?, ?, ?, ?, '', ?, ?)`,
		red.ID, red.VoucherID, red.Code, red.ClaimerID, red.IPAddress, red.RedeemedAt,
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return Voucher{}, Redemption{}, errAlreadyClaimed
		}
		return Voucher{}, Redemption{}, fmt.Errorf("insert redemption: %w", err)
	}

	res, err := tx.ExecContext(ctx,
		`UPDATE vouchers SET redeemed_count = redeemed_count + 1, updated_at = ?
		 WHERE id = ? AND redeemed_count < max_redemptions`,
		nowStr, v.ID,
	)
	if err != nil {
		return Voucher{}, Redemption{}, fmt.Errorf("increment voucher redeemed_count: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return Voucher{}, Redemption{}, errVoucherExhausted
	}

	if err := tx.Commit(); err != nil {
		return Voucher{}, Redemption{}, fmt.Errorf("commit claim tx: %w", err)
	}
	v.RedeemedCount++
	v.UpdatedAt = nowStr
	return v, red, nil
}

func (r *repo) attachRedemptionKey(ctx context.Context, redemptionID, keyID string) error {
	_, err := r.db.Writer().ExecContext(ctx,
		`UPDATE voucher_redemptions SET key_id = ? WHERE id = ?`,
		keyID, redemptionID,
	)
	return err
}

func (r *repo) rollbackClaim(ctx context.Context, voucherID, redemptionID string) {
	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer func() { _ = tx.Rollback() }()
	_, _ = tx.ExecContext(ctx, `DELETE FROM voucher_redemptions WHERE id = ?`, redemptionID)
	_, _ = tx.ExecContext(ctx, `UPDATE vouchers SET redeemed_count = MAX(0, redeemed_count - 1) WHERE id = ?`, voucherID)
	_ = tx.Commit()
}

func scanVoucher(rows *sql.Rows) (Voucher, error) {
	var v Voucher
	var allowedJSON string
	var enabled int
	if err := rows.Scan(
		&v.ID, &v.Code, &v.Name, &v.Description, &v.QuotaUSD, &v.RateLimitRPM, &allowedJSON,
		&v.TeamID, &v.MaxRedemptions, &v.RedeemedCount, &enabled, &v.ExpiresAt, &v.CreatedAt, &v.UpdatedAt,
	); err != nil {
		return Voucher{}, fmt.Errorf("scan voucher: %w", err)
	}
	v.Enabled = enabled != 0
	v.AllowedModels = []string{}
	if strings.TrimSpace(allowedJSON) != "" {
		_ = json.Unmarshal([]byte(allowedJSON), &v.AllowedModels)
	}
	return v, nil
}

func scanVoucherRow(row *sql.Row) (Voucher, error) {
	var v Voucher
	var allowedJSON string
	var enabled int
	if err := row.Scan(
		&v.ID, &v.Code, &v.Name, &v.Description, &v.QuotaUSD, &v.RateLimitRPM, &allowedJSON,
		&v.TeamID, &v.MaxRedemptions, &v.RedeemedCount, &enabled, &v.ExpiresAt, &v.CreatedAt, &v.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Voucher{}, errVoucherNotFound
		}
		return Voucher{}, fmt.Errorf("scan voucher row: %w", err)
	}
	v.Enabled = enabled != 0
	v.AllowedModels = []string{}
	if strings.TrimSpace(allowedJSON) != "" {
		_ = json.Unmarshal([]byte(allowedJSON), &v.AllowedModels)
	}
	return v, nil
}
