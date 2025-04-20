package currency

import (
	"encoding/json"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    BRL
		expected string
	}{
		{"integer value", BRL(123456), `"1234.56"`},
		{"zero value", BRL(0), `"0.00"`},
		{"negative value", BRL(-12345), `"-123.45"`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.input.MarshalJSON()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result) != test.expected {
				t.Errorf("expected: %s, got: %s", test.expected, result)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected BRL
		hasError bool
	}{
		{"valid value", `"1234.56"`, BRL(123456), false},
		{"zero value", `"0.00"`, BRL(0), false},
		{"negative value", `"-123.45"`, BRL(-12345), false},
		{"invalid format", `"1234,56"`, 0, true},
		{"not a number", `"abc"`, 0, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result BRL
			err := json.Unmarshal([]byte(test.input), &result)
			if test.hasError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result != test.expected {
					t.Errorf("expected: %v, got: %v", test.expected, result)
				}
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		input    BRL
		expected string
	}{
		{"integer value", BRL(123456), "1234.56"},
		{"zero value", BRL(0), "0.00"},
		{"negative value", BRL(-12345), "-123.45"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := test.input.String(); result != test.expected {
				t.Errorf("expected: %s, got: %s", test.expected, result)
			}
		})
	}
}

func TestCents(t *testing.T) {
	tests := []struct {
		name     string
		input    BRL
		expected int64
	}{
		{"positive value", BRL(123456), 123456},
		{"zero value", BRL(0), 0},
		{"negative value", BRL(-12345), -12345},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := test.input.Cents(); result != test.expected {
				t.Errorf("expected: %d, got: %d", test.expected, result)
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		name     string
		input    BRL
		expected bool
	}{
		{"zero value", BRL(0), true},
		{"positive value", BRL(12345), false},
		{"negative value", BRL(-12345), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := test.input.IsZero(); result != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     BRL
		expected BRL
	}{
		{"positive values", BRL(12345), BRL(67890), BRL(80235)},
		{"zero value", BRL(12345), BRL(0), BRL(12345)},
		{"negative value", BRL(12345), BRL(-4567), BRL(7788)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := test.a.Add(test.b); result != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		name     string
		a, b     BRL
		expected BRL
	}{
		{"positive values", BRL(12345), BRL(6789), BRL(5556)},
		{"zero value", BRL(12345), BRL(0), BRL(12345)},
		{"negative value", BRL(12345), BRL(-4567), BRL(16912)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := test.a.Sub(test.b); result != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}

func TestMul(t *testing.T) {
	tests := []struct {
		name     string
		a, b     BRL
		expected BRL
	}{
		{"positive values", BRL(123), BRL(456), BRL(56088)},
		{"zero value", BRL(12345), BRL(0), BRL(0)},
		{"negative value", BRL(123), BRL(-2), BRL(-246)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := test.a.Mul(test.b); result != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}

func TestDiv(t *testing.T) {
	tests := []struct {
		name      string
		a, b      BRL
		expected  BRL
		expectErr bool
	}{
		{"positive values", BRL(1234), BRL(2), BRL(617), false},
		{"division by zero", BRL(1234), BRL(0), 0, true},
		{"negative value", BRL(1234), BRL(-2), BRL(-617), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.a.Div(test.b)
			if test.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result != test.expected {
					t.Errorf("expected: %v, got: %v", test.expected, result)
				}
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name     string
		input    BRL
		n        int
		expected []BRL
	}{
		{"even split", BRL(1000), 4, []BRL{250, 250, 250, 250}},
		{"uneven split", BRL(1001), 4, []BRL{251, 250, 250, 250}},
		{"zero value", BRL(0), 4, []BRL{0, 0, 0, 0}},
		{"zero parts", BRL(1000), 0, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.input.Split(test.n)
			if len(result) != len(test.expected) {
				t.Fatalf("expected length: %d, got: %d", len(test.expected), len(result))
			}
			for i, v := range result {
				if v != test.expected[i] {
					t.Errorf("at index %d: expected: %v, got: %v", i, test.expected[i], v)
				}
			}
		})
	}
}
