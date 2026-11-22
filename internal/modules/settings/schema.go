package settings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	// SettingTypeEnum identifies a constrained string enumeration setting.
	SettingTypeEnum = "enum"
	// SettingTypeBoolean identifies a strict boolean setting.
	SettingTypeBoolean = "boolean"
	// SettingTypeInteger identifies a bounded integer setting.
	SettingTypeInteger = "integer"

	codeInvalidSettingKey   = "invalid_setting_key"
	codeInvalidSettingValue = "invalid_setting_value"
)

// ValidationError carries a machine-readable API error code and message for
// setting key or typed value validation failures.
type ValidationError struct {
	Code    string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// SettingItem represents a single typed setting merged from built-in schema
// metadata and any persisted database override.
type SettingItem struct {
	Key          string   `json:"key"`
	Type         string   `json:"type"`
	Value        any      `json:"value"`
	DefaultValue any      `json:"default_value"`
	Revision     int64    `json:"revision"`
	UpdatedAt    string   `json:"updated_at"`
	UpdatedBy    string   `json:"updated_by,omitempty"`
	Description  string   `json:"description"`
	Options      []string `json:"options,omitempty"`
	Min          *int64   `json:"min,omitempty"`
	Max          *int64   `json:"max,omitempty"`
}

// HistoryItem represents a single revision history record from settings_history.
type HistoryItem struct {
	ID         int64  `json:"id"`
	SettingKey string `json:"setting_key"`
	OldValue   any    `json:"old_value"`
	NewValue   any    `json:"new_value"`
	Revision   int64  `json:"revision"`
	ChangedBy  string `json:"changed_by"`
	Reason     string `json:"reason"`
	CreatedAt  string `json:"created_at"`
}

type settingDefinition struct {
	Key          string
	Type         string
	DefaultValue any
	Description  string
	Options      []string
	Min          *int64
	Max          *int64
}

func int64Ptr(v int64) *int64 {
	return &v
}

var builtinDefinitions = []settingDefinition{
	{
		Key:          "ui.theme",
		Type:         SettingTypeEnum,
		DefaultValue: "dark",
		Description:  "Dashboard color theme mode.",
		Options:      []string{"dark", "light", "system"},
	},
	{
		Key:          "ui.locale",
		Type:         SettingTypeEnum,
		DefaultValue: "en",
		Description:  "Dashboard display language locale.",
		Options:      []string{"en", "id"},
	},
	{
		Key:          "usage.payload_logging_enabled",
		Type:         SettingTypeBoolean,
		DefaultValue: false,
		Description:  "Enable capturing request and response payload bodies for usage inspection.",
	},
	{
		Key:          "usage.payload_retention_days",
		Type:         SettingTypeInteger,
		DefaultValue: int64(7),
		Description:  "Number of days to retain captured request and response payloads.",
		Min:          int64Ptr(1),
		Max:          int64Ptr(90),
	},
	{
		Key:          "usage.metadata_retention_days",
		Type:         SettingTypeInteger,
		DefaultValue: int64(90),
		Description:  "Number of days to retain request usage metadata and cost aggregates.",
		Min:          int64Ptr(1),
		Max:          int64Ptr(3650),
	},
	{
		Key:          "gateway.default_timeout_seconds",
		Type:         SettingTypeInteger,
		DefaultValue: int64(120),
		Description:  "Default upstream request timeout in seconds when no per-route override is set.",
		Min:          int64Ptr(5),
		Max:          int64Ptr(600),
	},
	{
		Key:          "gateway.allow_response_chain_fallback",
		Type:         SettingTypeBoolean,
		DefaultValue: true,
		Description:  "Allow falling back to alternative upstream targets when a response chain target fails.",
	},
}

var builtinDefinitionMap = func() map[string]settingDefinition {
	m := make(map[string]settingDefinition, len(builtinDefinitions))
	for _, def := range builtinDefinitions {
		m[def.Key] = def
	}
	return m
}()

func lookupDefinition(key string) (settingDefinition, bool) {
	def, ok := builtinDefinitionMap[key]
	return def, ok
}

func (d settingDefinition) defaultJSON() string {
	switch d.Type {
	case SettingTypeEnum:
		s, _ := d.DefaultValue.(string)
		b, err := json.Marshal(s)
		if err != nil {
			panic(fmt.Sprintf("settings: marshal default enum %q: %v", d.Key, err))
		}
		return string(b)
	case SettingTypeBoolean:
		b, _ := d.DefaultValue.(bool)
		if b {
			return "true"
		}
		return "false"
	case SettingTypeInteger:
		n, _ := d.DefaultValue.(int64)
		return strconv.FormatInt(n, 10)
	default:
		panic(fmt.Sprintf("settings: unsupported setting type %q for %q", d.Type, d.Key))
	}
}

func (d settingDefinition) buildItem(value any, revision int64, updatedAt, updatedBy string) SettingItem {
	item := SettingItem{
		Key:          d.Key,
		Type:         d.Type,
		Value:        value,
		DefaultValue: d.DefaultValue,
		Revision:     revision,
		UpdatedAt:    updatedAt,
		UpdatedBy:    updatedBy,
		Description:  d.Description,
	}
	if len(d.Options) > 0 {
		item.Options = make([]string, len(d.Options))
		copy(item.Options, d.Options)
	}
	if d.Min != nil {
		v := *d.Min
		item.Min = &v
	}
	if d.Max != nil {
		v := *d.Max
		item.Max = &v
	}
	return item
}

func (d settingDefinition) validateRaw(raw json.RawMessage) (val any, canonicalJSON string, err error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, "", &ValidationError{
			Code:    codeInvalidSettingValue,
			Message: fmt.Sprintf("setting %q requires a non-empty value", d.Key),
		}
	}

	switch d.Type {
	case SettingTypeEnum:
		if trimmed[0] != '"' {
			return nil, "", &ValidationError{
				Code:    codeInvalidSettingValue,
				Message: fmt.Sprintf("setting %q must be a string (%s)", d.Key, strings.Join(d.Options, ", ")),
			}
		}
		var s string
		if err := json.Unmarshal(trimmed, &s); err != nil {
			return nil, "", &ValidationError{
				Code:    codeInvalidSettingValue,
				Message: fmt.Sprintf("setting %q must be a valid JSON string", d.Key),
			}
		}
		if !slices.Contains(d.Options, s) {
			return nil, "", &ValidationError{
				Code:    codeInvalidSettingValue,
				Message: fmt.Sprintf("setting %q must be one of: %s", d.Key, strings.Join(d.Options, ", ")),
			}
		}
		canonical, err := json.Marshal(s)
		if err != nil {
			return nil, "", fmt.Errorf("settings: marshal enum %q: %w", d.Key, err)
		}
		return s, string(canonical), nil

	case SettingTypeBoolean:
		switch string(trimmed) {
		case "true":
			return true, "true", nil
		case "false":
			return false, "false", nil
		default:
			return nil, "", &ValidationError{
				Code:    codeInvalidSettingValue,
				Message: fmt.Sprintf("setting %q must be a boolean", d.Key),
			}
		}

	case SettingTypeInteger:
		if !isStrictJSONInteger(trimmed) {
			return nil, "", &ValidationError{
				Code:    codeInvalidSettingValue,
				Message: fmt.Sprintf("setting %q must be an integer", d.Key),
			}
		}
		n, err := strconv.ParseInt(string(trimmed), 10, 64)
		if err != nil {
			return nil, "", &ValidationError{
				Code:    codeInvalidSettingValue,
				Message: fmt.Sprintf("setting %q must be a valid 64-bit integer", d.Key),
			}
		}
		if d.Min != nil && n < *d.Min {
			return nil, "", &ValidationError{
				Code:    codeInvalidSettingValue,
				Message: fmt.Sprintf("setting %q must be between %d and %d", d.Key, *d.Min, *d.Max),
			}
		}
		if d.Max != nil && n > *d.Max {
			return nil, "", &ValidationError{
				Code:    codeInvalidSettingValue,
				Message: fmt.Sprintf("setting %q must be between %d and %d", d.Key, *d.Min, *d.Max),
			}
		}
		return n, strconv.FormatInt(n, 10), nil

	default:
		return nil, "", fmt.Errorf("settings: unsupported type %q for %q", d.Type, d.Key)
	}
}

func isStrictJSONInteger(raw []byte) bool {
	if len(raw) == 0 {
		return false
	}
	start := 0
	if raw[0] == '-' {
		if len(raw) == 1 {
			return false
		}
		start = 1
	}
	for i := start; i < len(raw); i++ {
		if raw[i] < '0' || raw[i] > '9' {
			return false
		}
	}
	return true
}

func decodeStoredValue(key, rawJSON string) (any, error) {
	if def, ok := lookupDefinition(key); ok {
		val, _, err := def.validateRaw(json.RawMessage(rawJSON))
		if err != nil {
			return nil, fmt.Errorf("settings: decode stored value for %q: %w", key, err)
		}
		return val, nil
	}
	var generic any
	if err := json.Unmarshal([]byte(rawJSON), &generic); err != nil {
		return nil, fmt.Errorf("settings: decode legacy value for %q: %w", key, err)
	}
	return generic, nil
}
