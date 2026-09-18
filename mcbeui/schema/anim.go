package schema

import (
	"encoding/json"
	"fmt"
)

// AnimBuilder builds a single animation entry with fluent chaining.
type AnimBuilder struct {
	a Anim
}

// NewAnim starts building an animation. Pass the animation type.
func NewAnim(typ AnimType) *AnimBuilder {
	return &AnimBuilder{a: Anim{AnimType: typ}}
}

// Anim returns the built animation.
func (ab *AnimBuilder) Anim() *Anim { return &ab.a }

// Duration sets the animation duration in seconds.
func (ab *AnimBuilder) Duration(d float64) *AnimBuilder { ab.a.Duration = d; return ab }

// From sets the starting value (number, vector array, etc.).
func (ab *AnimBuilder) From(v interface{}) *AnimBuilder {
	ab.a.From = mustMarshal(v)
	return ab
}

// To sets the ending value (number, vector array, etc.).
func (ab *AnimBuilder) To(v interface{}) *AnimBuilder {
	ab.a.To = mustMarshal(v)
	return ab
}

// Easing sets the interpolation curve.
func (ab *AnimBuilder) Easing(e Easing) *AnimBuilder { ab.a.Easing = e; return ab }

// Looping marks the animation as looping.
func (ab *AnimBuilder) Looping() *AnimBuilder { ab.a.Looping = true; return ab }

// Reversible marks the animation as reversible.
func (ab *AnimBuilder) Reversible() *AnimBuilder { ab.a.Reversible = true; return ab }

// Resettable marks the animation as resettable.
func (ab *AnimBuilder) Resettable() *AnimBuilder { ab.a.Resettable = true; return ab }

// PlayEvent sets the event played when the animation starts.
func (ab *AnimBuilder) PlayEvent(name string) *AnimBuilder { ab.a.PlayEvent = name; return ab }

// EndEvent sets the event played when the animation ends.
func (ab *AnimBuilder) EndEvent(name string) *AnimBuilder { ab.a.EndEvent = name; return ab }

// DestroyAtEnd destroys the control when the animation finishes.
func (ab *AnimBuilder) DestroyAtEnd() *AnimBuilder { ab.a.DestroyAtEnd = true; return ab }

// --- Convenience constructors ---

// AnimAlpha creates an alpha animation.
func AnimAlpha(duration, from, to float64) Anim {
	return NewAnim(AnimTypeAlpha).Duration(duration).From(from).To(to).Anim().Copy()
}

// AnimSize creates a size animation.
func AnimSize(duration float64, from, to Vector2) Anim {
	return NewAnim(AnimTypeSize).Duration(duration).
		From([]float64{from.X, from.Y}).To([]float64{to.X, to.Y}).Anim().Copy()
}

// AnimOffset creates an offset animation.
func AnimOffset(duration float64, from, to Vector2) Anim {
	return NewAnim(AnimTypeOffset).Duration(duration).
		From([]float64{from.X, from.Y}).To([]float64{to.X, to.Y}).Anim().Copy()
}

// AnimColor creates a color animation.
func AnimColor(duration float64, from, to []float64) Anim {
	return NewAnim(AnimTypeColor).Duration(duration).From(from).To(to).Anim().Copy()
}

// AnimUV creates a UV animation.
func AnimUV(duration float64, from, to Vector2) Anim {
	return NewAnim(AnimTypeUV).Duration(duration).
		From([]float64{from.X, from.Y}).To([]float64{to.X, to.Y}).Anim().Copy()
}

// --- Helper ---

// Copy returns a value copy of the animation.
func (a *Anim) Copy() Anim { return *a }

func mustMarshal(v interface{}) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal animation value: %v", err))
	}
	return b
}

// Anim represents a single animation entry.
type Anim struct {
	AnimType                 AnimType        `json:"anim_type"`
	Duration                 float64         `json:"duration,omitempty"`
	Next                     string          `json:"next,omitempty"`
	DestroyAtEnd             bool            `json:"destroy_at_end,omitempty"`
	PlayEvent                string          `json:"play_event,omitempty"`
	EndEvent                 string          `json:"end_event,omitempty"`
	StartEvent               string          `json:"start_event,omitempty"`
	ResetEvent               string          `json:"reset_event,omitempty"`
	Easing                   Easing          `json:"easing,omitempty"`
	From                     json.RawMessage `json:"from,omitempty"`
	To                       json.RawMessage `json:"to,omitempty"`
	InitialUV                json.RawMessage `json:"initial_uv,omitempty"`
	FPS                      float64         `json:"fps,omitempty"`
	FrameCount               int             `json:"frame_count,omitempty"`
	FrameStep                float64         `json:"frame_step,omitempty"`
	Reversible               bool            `json:"reversible,omitempty"`
	Resettable               bool            `json:"resettable,omitempty"`
	ScaleFromStartingAlpha   bool            `json:"scale_from_starting_alpha,omitempty"`
	Activated                bool            `json:"activated,omitempty"`
	Looping                  bool            `json:"looping,omitempty"`
}
