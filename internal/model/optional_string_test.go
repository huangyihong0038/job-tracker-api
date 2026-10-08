package model

import (
	"encoding/json"
	"testing"
)

func TestOptionalStringUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantPresent bool
		wantNil     bool
		wantValue   string
		wantError   bool
	}{
		{
			name:        "field is absent",
			body:        `{}`,
			wantPresent: false,
			wantNil:     true,
		},
		{
			name:        "field is null",
			body:        `{"applied_at": null}`,
			wantPresent: true,
			wantNil:     true,
		},
		{
			name:        "field contains a string",
			body:        `{"applied_at": "2026-10-07"}`,
			wantPresent: true,
			wantValue:   "2026-10-07",
		},
		{
			name:      "field has invalid type",
			body:      `{"applied_at": 123}`,
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var request UpdateApplicationRequest

			err := json.Unmarshal(
				[]byte(test.body),
				&request,
			)

			if test.wantError {
				if err == nil {
					t.Fatal(
						"expected JSON unmarshal error",
					)
				}
				return
			}

			if err != nil {
				t.Fatalf(
					"unexpected JSON unmarshal error: %v",
					err,
				)
			}

			field := request.AppliedAt

			if field.Present != test.wantPresent {
				t.Fatalf(
					"expected Present=%t, got %t",
					test.wantPresent,
					field.Present,
				)
			}

			if test.wantNil {
				if field.Value != nil {
					t.Fatalf(
						"expected nil Value, got %q",
						*field.Value,
					)
				}
				return
			}

			if field.Value == nil {
				t.Fatal(
					"expected non-nil Value",
				)
			}

			if *field.Value != test.wantValue {
				t.Fatalf(
					"expected Value %q, got %q",
					test.wantValue,
					*field.Value,
				)
			}
		})
	}
}
