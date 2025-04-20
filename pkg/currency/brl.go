package currency

import (
	"encoding/json"
	"fmt"
)

// BRL represents a monetary value in Brazilian Reais stored as an integer in centavos.
type BRL int64

// MarshalJSON customizes the JSON encoding for the BRL type, representing it as a string in the format "major.minor".
func (b BRL) MarshalJSON() ([]byte, error) {

	value := b / 100
	cents := b % 100
	if cents < 0 {
		cents = -cents
	}

	return []byte(fmt.Sprintf(`"%d.%02d"`, value, cents)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for the BRL type, parsing a stringified monetary value in cents.
func (b *BRL) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// Parse string in the format "%d.%02d" back to BRL (integer representation in cents).
	var major, minor int
	if _, err := fmt.Sscanf(str, "%d.%02d", &major, &minor); err != nil {
		return err
	}
	*b = BRL(major*100 + minor)
	return nil
}

// String formats the BRL value into a human-readable string using the format "major.minor" (e.g., "123.45").
func (b BRL) String() string {
	return fmt.Sprintf("%d.%02d", b/100, b%100)
}

// Cents returns the integer representation of the BRL amount in cents.
func (b BRL) Cents() int64 {
	return int64(b)
}

// IsZero checks if the BRL value is equal to zero.
func (b BRL) IsZero() bool {
	return b == 0
}

// Add returns the sum of the current BRL value and another BRL value.
func (b BRL) Add(other BRL) BRL {
	return b + other
}

// Sub subtracts the value of the given BRL from the receiver and returns the resulting BRL.
func (b BRL) Sub(other BRL) BRL {
	return b - other
}

// Mul multiplies the current BRL value by the given BRL value and returns the resulting BRL value.
func (b BRL) Mul(other BRL) BRL {
	return b * other
}

// Div divides the receiver BRL value by the provided BRL value and returns the result or an error if division by zero occurs.
func (b BRL) Div(other BRL) (BRL, error) {
	if other == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return b / other, nil
}

// Split splits the BRL amount into `n` parts, distributing the remainder evenly across the parts where applicable.
func (b BRL) Split(n int) []BRL {
	if n == 0 {
		return nil
	}

	parts := make([]BRL, n)

	total := b.Cents()

	// Calculate the base amount to distribute
	baseAmount := total / int64(n)

	// Calculate the remainder to distribute
	remainder := total % int64(n)

	// Distribute the amounts
	for i := 0; i < n; i++ {
		parts[i] = BRL(baseAmount)
		if int64(i) < remainder { // Distribute the remainder in a round-robin manner
			parts[i] += 1
		}
	}

	return parts
}
