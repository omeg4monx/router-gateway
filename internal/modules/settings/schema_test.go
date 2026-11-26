package settings

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestSettingDefinitionValidateRawIntegerBounds(t *testing.T) {
	t.Run("min only", func(t *testing.T) {
		def := settingDefinition{
			Key:  "test.min_only",
			Type: SettingTypeInteger,
			Min:  int64Ptr(10),
		}

		if val, canonical, err := def.validateRaw(json.RawMessage("10")); err != nil || val != int64(10) || canonical != "10" {
			t.Fatalf("validateRaw(10) = (%v, %q, %v), want (10, \"10\", nil)", val, canonical, err)
		}

		_, _, err := def.validateRaw(json.RawMessage("9"))
		var valErr *ValidationError
		if !errors.As(err, &valErr) {
			t.Fatalf("validateRaw(9) err = %v, want *ValidationError", err)
		}
		if !strings.Contains(valErr.Message, "10") || !strings.Contains(valErr.Message, "at least") {
			t.Fatalf("unexpected error message for min-only: %q", valErr.Message)
		}
	})

	t.Run("max only", func(t *testing.T) {
		def := settingDefinition{
			Key:  "test.max_only",
			Type: SettingTypeInteger,
			Max:  int64Ptr(100),
		}

		if val, canonical, err := def.validateRaw(json.RawMessage("100")); err != nil || val != int64(100) || canonical != "100" {
			t.Fatalf("validateRaw(100) = (%v, %q, %v), want (100, \"100\", nil)", val, canonical, err)
		}

		_, _, err := def.validateRaw(json.RawMessage("101"))
		var valErr *ValidationError
		if !errors.As(err, &valErr) {
			t.Fatalf("validateRaw(101) err = %v, want *ValidationError", err)
		}
		if !strings.Contains(valErr.Message, "100") || !strings.Contains(valErr.Message, "at most") {
			t.Fatalf("unexpected error message for max-only: %q", valErr.Message)
		}
	})

	t.Run("both min and max", func(t *testing.T) {
		def := settingDefinition{
			Key:  "test.bounded",
			Type: SettingTypeInteger,
			Min:  int64Ptr(1),
			Max:  int64Ptr(90),
		}

		for _, raw := range []string{"0", "91"} {
			_, _, err := def.validateRaw(json.RawMessage(raw))
			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Fatalf("validateRaw(%s) err = %v, want *ValidationError", raw, err)
			}
			want := `setting "test.bounded" must be between 1 and 90`
			if valErr.Message != want {
				t.Fatalf("validateRaw(%s) message = %q, want %q", raw, valErr.Message, want)
			}
		}
	})
}
