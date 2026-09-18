package transitions_test

import (
	"encoding/json"
	"testing"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/transitions"
)

func TestNewTransition(t *testing.T) {
	tr := transitions.NewTransition("from_ctrl", "to_ctrl")
	if tr.From != "from_ctrl" {
		t.Errorf("From = %q, want %q", tr.From, "from_ctrl")
	}
	if tr.To != "to_ctrl" {
		t.Errorf("To = %q, want %q", tr.To, "to_ctrl")
	}
	if tr.Duration != 0.3 {
		t.Errorf("Duration = %v, want 0.3", tr.Duration)
	}
	if tr.Easing != schema.EasingLinear {
		t.Errorf("Easing = %q, want %q", tr.Easing, schema.EasingLinear)
	}
}

func TestTransitionEase(t *testing.T) {
	tr := transitions.TransitionEase()
	if tr.Duration != 0.3 {
		t.Errorf("Duration = %v, want 0.3", tr.Duration)
	}
	if tr.Easing != schema.EasingLinear {
		t.Errorf("Easing = %q, want %q", tr.Easing, schema.EasingLinear)
	}
	if !tr.FromOffset.IsEmpty() || !tr.ToOffset.IsEmpty() {
		t.Error("Ease should have empty offsets")
	}
}

func TestTransitionSlideLeft(t *testing.T) {
	tr := transitions.TransitionSlideLeft()
	if tr.FromOffset.PixelX() != -200 {
		t.Errorf("FromOffset.X = %v, want -200", tr.FromOffset.PixelX())
	}
	if tr.ToOffset.PixelX() != 0 {
		t.Errorf("ToOffset.X = %v, want 0", tr.ToOffset.PixelX())
	}
}

func TestTransitionSlideRight(t *testing.T) {
	tr := transitions.TransitionSlideRight()
	if tr.FromOffset.PixelX() != 200 {
		t.Errorf("FromOffset.X = %v, want 200", tr.FromOffset.PixelX())
	}
	if tr.ToOffset.PixelX() != 0 {
		t.Errorf("ToOffset.X = %v, want 0", tr.ToOffset.PixelX())
	}
}

func TestTransitionFade(t *testing.T) {
	tr := transitions.TransitionFade()
	if tr.Duration != 0.3 {
		t.Errorf("Duration = %v, want 0.3", tr.Duration)
	}
	if !tr.FromOffset.IsEmpty() || !tr.ToOffset.IsEmpty() {
		t.Error("Fade should have empty offsets")
	}
}

func TestTransitionScale(t *testing.T) {
	tr := transitions.TransitionScale()
	if tr.FromOffset.PixelX() != 0 || tr.FromOffset.PixelY() != 0 {
		t.Errorf("FromOffset = %v, want [0, 0]", tr.FromOffset)
	}
	// Scale uses named values "fill", so PixelX/PixelY return 0
	b, _ := json.Marshal(tr.ToOffset)
	if string(b) != `["fill","fill"]` {
		t.Errorf("ToOffset = %s, want [\"fill\",\"fill\"]", string(b))
	}
}

func TestBuildTransitionAnimEase(t *testing.T) {
	tr := transitions.TransitionEase()
	anim := transitions.BuildTransitionAnim(tr)
	if anim.AnimType != schema.AnimTypeAlpha {
		t.Errorf("AnimType = %q, want %q", anim.AnimType, schema.AnimTypeAlpha)
	}
	if anim.Duration != 0.3 {
		t.Errorf("Duration = %v, want 0.3", anim.Duration)
	}
}

func TestBuildTransitionAnimFade(t *testing.T) {
	tr := transitions.TransitionFade()
	anim := transitions.BuildTransitionAnim(tr)
	if anim.AnimType != schema.AnimTypeAlpha {
		t.Errorf("AnimType = %q, want %q", anim.AnimType, schema.AnimTypeAlpha)
	}
}

func TestBuildTransitionAnimSlideLeft(t *testing.T) {
	tr := transitions.TransitionSlideLeft()
	anim := transitions.BuildTransitionAnim(tr)
	if anim.AnimType != schema.AnimTypeOffset {
		t.Errorf("AnimType = %q, want %q", anim.AnimType, schema.AnimTypeOffset)
	}
}

func TestBuildTransitionAnimSlideRight(t *testing.T) {
	tr := transitions.TransitionSlideRight()
	anim := transitions.BuildTransitionAnim(tr)
	if anim.AnimType != schema.AnimTypeOffset {
		t.Errorf("AnimType = %q, want %q", anim.AnimType, schema.AnimTypeOffset)
	}
}

func TestBuildTransitionAnimScale(t *testing.T) {
	tr := transitions.TransitionScale()
	anim := transitions.BuildTransitionAnim(tr)
	if anim.AnimType != schema.AnimTypeSize {
		t.Errorf("AnimType = %q, want %q", anim.AnimType, schema.AnimTypeSize)
	}
}

func TestBuildTransitionAnimJSON(t *testing.T) {
	tr := transitions.TransitionEase()
	anim := transitions.BuildTransitionAnim(tr)
	b, err := json.Marshal(anim)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if obj["anim_type"] != "alpha" {
		t.Errorf("anim_type = %v, want alpha", obj["anim_type"])
	}
}
