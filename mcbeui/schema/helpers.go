// Package schema exports UI-building constructors and helpers.
package schema

import "encoding/json"

// NewVector2 creates a new Vector2.
func NewVector2(x, y float64) Vector2 {
	return Vector2{X: x, Y: y}
}

// NewVector3 creates a new Vector3 color.
func NewVector3(r, g, b float64) Vector3 {
	return Vector3{R: r, G: g, B: b}
}

// ColorRGB creates a raw JSON array from RGB floats (0-1).
func ColorRGB(r, g, b float64) json.RawMessage {
	bits, _ := json.Marshal([]float64{r, g, b})
	return bits
}

// ColorHex creates a raw JSON string from a hex color like "#FF0000".
func ColorHex(hex string) json.RawMessage {
	b, _ := json.Marshal(hex)
	return b
}

// FromJSON unmarshals raw JSON into a Control.
func FromJSON(data []byte) (*Control, error) {
	var c Control
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// MustJSON panics if json.Marshal fails.
func MustJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
