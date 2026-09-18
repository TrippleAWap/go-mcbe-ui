// Package transitions provides pre-built animation transitions for screen changes.
package transitions

import (
	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
)

// TransitionID is a named transition animation used in navigation links.
type TransitionID string

const (
	// Fade is a fade-in transition (alpha 0→1).
	Fade TransitionID = "fade"
	// SlideLeft is a slide-in from the left transition.
	SlideLeft TransitionID = "slide-left"
	// SlideRight is a slide-in from the right transition.
	SlideRight TransitionID = "slide-right"
	// Scale is a scale-up from center transition.
	Scale TransitionID = "scale"
)

// Transition describes a screen transition between two controls.
type Transition struct {
	From       string
	To         string
	Duration   float64
	Easing     schema.Easing
	FromOffset schema.Relative
	ToOffset   schema.Relative
}

// NewTransition creates a new Transition between two control names.
func NewTransition(from, to string) *Transition {
	return &Transition{
		From:     from,
		To:       to,
		Duration: 0.3,
		Easing:   schema.EasingLinear,
	}
}

// TransitionEase returns a transition that eases alpha from 0 to 1 with default offset.
func TransitionEase() *Transition {
	return &Transition{
		Duration: 0.3,
		Easing:   schema.EasingLinear,
	}
}

// TransitionSlideLeft returns a transition that slides in from the left.
func TransitionSlideLeft() *Transition {
	return &Transition{
		Duration:   0.3,
		Easing:     schema.EasingLinear,
		FromOffset: schema.Raw(-200, 0),
		ToOffset:   schema.Raw(0, 0),
	}
}

// TransitionSlideRight returns a transition that slides in from the right.
func TransitionSlideRight() *Transition {
	return &Transition{
		Duration:   0.3,
		Easing:     schema.EasingLinear,
		FromOffset: schema.Raw(200, 0),
		ToOffset:   schema.Raw(0, 0),
	}
}

// TransitionFade returns a transition that fades in.
func TransitionFade() *Transition {
	return &Transition{
		Duration: 0.3,
		Easing:   schema.EasingLinear,
	}
}

// TransitionScale returns a transition that scales up from center.
func TransitionScale() *Transition {
	return &Transition{
		Duration:   0.3,
		Easing:     schema.EasingLinear,
		FromOffset: schema.Raw(0, 0),
		ToOffset:   schema.Raw("fill", "fill"),
	}
}

// WithDuration sets the transition duration and returns the transition for chaining.
func (t *Transition) WithDuration(d float64) *Transition {
	t.Duration = d
	return t
}

// WithEasing sets the easing curve and returns the transition for chaining.
func (t *Transition) WithEasing(e schema.Easing) *Transition {
	t.Easing = e
	return t
}

// BuildTransitionAnim generates a schema.Anim from the Transition's parameters.
// BuildTransitionAnim generates a schema.Anim from the Transition's parameters.
// Duration and Easing are respected; defaults (0.3s, linear) are used if unset.
func BuildTransitionAnim(t *Transition) schema.Anim {
	duration := t.Duration
	if duration <= 0 {
		duration = 0.3
	}
	easing := t.Easing
	if easing == "" {
		easing = schema.EasingLinear
	}
	switch {
	case t.FromOffset.IsEmpty() && t.ToOffset.IsEmpty():
		// Alpha transition (ease / fade)
		return schema.AnimAlpha(duration, 0, 1)
	case t.FromOffset.PixelX() == 0 && t.FromOffset.PixelY() == 0 &&
		t.ToOffset.PixelX() == 0 && t.ToOffset.PixelY() == 0 &&
		!t.FromOffset.IsEmpty():
		// Scale transition — non-empty but zero-pixel offsets with named values
		anim := schema.AnimSize(duration, schema.Vector2{X: 0, Y: 0},
			schema.Vector2{X: 1, Y: 1})
		anim.Easing = easing
		return anim
	default:
		// Offset transition (slide left / right)
		anim := schema.AnimOffset(duration,
			schema.Vector2{X: float64(t.FromOffset.PixelX()), Y: float64(t.FromOffset.PixelY())},
			schema.Vector2{X: float64(t.ToOffset.PixelX()), Y: float64(t.ToOffset.PixelY())})
		anim.Easing = easing
		return anim
	}
}
