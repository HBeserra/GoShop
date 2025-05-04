package domain_test

import (
	"github.com/HBeserra/GoShop/domain"
	"testing"
)

func TestByteSize_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.ByteSize
		expected string
	}{
		{"ZeroBytes", domain.ByteSize(0), `"0B"`},
		{"OneByte", domain.ByteSize(1), `"1B"`},
		{"One kibibyte", domain.ByteSize(1024), `"1KiB"`},
		{"OneMebibyte", domain.ByteSize(1024 * 1024), `"1MiB"`},
		{"LargeValue", domain.ByteSize(1024 * 1024 * 1024 * 10), `"10GiB"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON failed: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func TestByteSize_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    domain.ByteSize
		expectError bool
	}{
		{"ValidBytes", `"0 B"`, domain.ByteSize(0), false},
		{"ValidKilobytes", `"1kB"`, domain.ByteSize(1000), false},
		{"ValidMegabytes", `"1MB"`, domain.ByteSize(1000 * 1000), false},
		{"ValidGigabytes", `"10GB"`, domain.ByteSize(1000 * 1000 * 1000 * 10), false},
		{"InvalidString", `"invalid"`, 0, true},
		{"InvalidNumber", `1234`, 0, true},
		{"EmptyString", `""`, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result domain.ByteSize
			err := result.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
			if !tt.expectError && result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestByteSize_String(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.ByteSize
		expected string
	}{
		{"ZeroBytes", domain.ByteSize(0), "0B"},
		{"OneKilobyte", domain.ByteSize(1024), "1KiB"},
		{"OneMegabyte", domain.ByteSize(1024 * 1024), "1MiB"},
		{"LargeValue", domain.ByteSize(1024 * 1024 * 1024 * 10), "10GiB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
