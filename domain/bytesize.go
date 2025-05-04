package domain

import (
	"encoding/json"
	"github.com/docker/go-units"
)

type ByteSize int64

func (b *ByteSize) MarshalJSON() ([]byte, error) {
	str := units.BytesSize(float64(*b))
	return json.Marshal(str)
}
func (b *ByteSize) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	val, err := units.FromHumanSize(str)
	if err != nil {
		return err
	}
	*b = ByteSize(val)
	return nil
}

func (b ByteSize) String() string {
	return units.BytesSize(float64(b))
}
