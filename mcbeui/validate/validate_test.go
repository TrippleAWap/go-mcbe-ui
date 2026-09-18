package validate_test

import (
	"testing"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/controls"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/validate"
)

func TestValidateNilScreen(t *testing.T) {
	r := validate.ValidateScreen(nil)
	if r.IsValid() {
		t.Error("expected validation to fail for nil screen")
	}
}

func TestValidateValidScreen(t *testing.T) {
	root := controls.NewScreen().ID("test").Anchor(schema.Center).Size(schema.Raw(1920, 1080))
	screen := &schema.Screen{
		Namespace: "test",
		RootPanel: "root",
		Root:      *root.Control(),
	}
	r := validate.ValidateScreen(screen)
	if !r.IsValid() {
		t.Errorf("expected valid screen, got: %s", r.String())
	}
}

func TestValidateNegativeSize(t *testing.T) {
	c := controls.NewPanel().
		ID("bad").
		Size(schema.Raw(-10, -20)).
		Control()
	r := validate.ValidateControl(c)
	if len(r.Warnings) == 0 {
		t.Error("expected warning for negative size")
	}
}

func TestValidateInvalidAnchor(t *testing.T) {
	c := controls.NewPanel().
		ID("bad_anchor").
		AnchorFrom("not_an_anchor").
		Control()
	r := validate.ValidateControl(c)
	if len(r.Errors) == 0 {
		t.Error("expected error for invalid anchor")
	}
}

func TestValidateMissingBindingName(t *testing.T) {
	c := controls.NewPanel().
		ID("bad_bind").
		Bindings(schema.Bind{
			BindingType: schema.BindingTypeView,
		}).
		Control()
	r := validate.ValidateControl(c)
	if len(r.Warnings) == 0 {
		t.Error("expected warning for missing binding name")
	}
}

func TestValidateViewBindingMissingProps(t *testing.T) {
	c := controls.NewPanel().
		ID("bad_view_bind").
		Bindings(schema.Bind{
			BindingType: schema.BindingTypeView,
			BindingName: "test",
		}).
		Control()
	r := validate.ValidateControl(c)
	if len(r.Errors) == 0 {
		t.Error("expected error for incomplete view binding")
	}
}

func TestValidateCollectionBindingMissingCollection(t *testing.T) {
	c := controls.NewPanel().
		ID("bad_col_bind").
		Bindings(schema.Bind{
			BindingType: schema.BindingTypeCollection,
			BindingName: "test",
		}).
		Control()
	r := validate.ValidateControl(c)
	if len(r.Errors) == 0 {
		t.Error("expected error for collection binding without collection name")
	}
}

func TestValidateInvalidAnimType(t *testing.T) {
	c := controls.NewPanel().
		ID("bad_anim").
		Animations(schema.Anim{
			Duration: 1.0,
		}).
		Control()
	r := validate.ValidateControl(c)
	if len(r.Errors) == 0 {
		t.Error("expected error for animation without anim_type")
	}
}

func TestValidateHighLayer(t *testing.T) {
	c := controls.NewPanel().
		ID("high_layer").
		Layer(150).
		Control()
	r := validate.ValidateControl(c)
	if len(r.Warnings) == 0 {
		t.Error("expected warning for layer > 100")
	}
}

func TestValidateMissingReference(t *testing.T) {
	root := controls.NewPanel().
		ID("root").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100))

	inner := controls.NewPanel().ID("inner")
	btn := controls.NewButton().
		ID("btn").
		Anchor(schema.TopLeft).
		Size(schema.Raw(50, 50)).
		DefaultControl("missing_btn")

	root.AddControl(inner.Control().ID)
	root.AddControl(btn.Control().ID)

	allControls := map[string]*schema.Control{
		"inner": inner.Control(),
		"btn":   btn.Control(),
	}
	screen := &schema.Screen{
		Namespace: "test",
		RootPanel: "root",
		Root:      *root.Control(),
	}
	r := validate.ValidateScreenWithControls(screen, allControls)
	if len(r.Errors) == 0 {
		t.Error("expected error for missing referenced control")
	}
}

func TestValidateExistingReference(t *testing.T) {
	root := controls.NewPanel().
		ID("root").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100))

	inner := controls.NewPanel().ID("inner")
	btn := controls.NewButton().
		ID("btn").
		Anchor(schema.TopLeft).
		Size(schema.Raw(50, 50)).
		DefaultControl("inner")

	root.AddControl(inner.Control().ID)
	root.AddControl(btn.Control().ID)

	allControls := map[string]*schema.Control{
		"inner": inner.Control(),
		"btn":   btn.Control(),
	}
	screen := &schema.Screen{
		Namespace: "test",
		RootPanel: "root",
		Root:      *root.Control(),
	}
	r := validate.ValidateScreenWithControls(screen, allControls)
	if !r.IsValid() {
		t.Errorf("expected no errors, got: %s", r.String())
	}
}

func TestResultString(t *testing.T) {
	r := &validate.Result{
		Errors:   []validate.Error{{ControlID: "x", Message: "bad"}},
		Warnings: []validate.Warning{{ControlID: "y", Message: "warn"}},
	}
	s := r.String()
	if s == "validation passed" {
		t.Error("expected non-trivial string output")
	}
}

func TestResultIsValid(t *testing.T) {
	r := &validate.Result{}
	if !r.IsValid() {
		t.Error("empty result should be valid")
	}

	r.Errors = append(r.Errors, validate.Error{Message: "fail"})
	if r.IsValid() {
		t.Error("result with errors should not be valid")
	}
}

func TestValidateNilControl(t *testing.T) {
	r := validate.ValidateControl(nil)
	if !r.IsValid() {
		t.Errorf("expected nil control to produce no errors, got: %s", r.String())
	}
}

func TestValidateEmptyID(t *testing.T) {
	c := controls.NewPanel().
		// No ID set
		Anchor(schema.Center).
		Size(schema.Raw(100, 100)).
		Control()
	r := validate.ValidateControl(c)
	if len(r.Errors) == 0 {
		t.Error("expected error for control with empty ID")
	}
}

func TestValidateMultipleInvalidAnchors(t *testing.T) {
	c := controls.NewPanel().
		ID("bad_anchors").
		AnchorFrom("not_valid").
		AnchorTo("also_not_valid").
		Control()
	r := validate.ValidateControl(c)
	if len(r.Errors) < 2 {
		t.Errorf("expected at least 2 errors for two invalid anchors, got %d: %s", len(r.Errors), r.String())
	}
}

func TestValidateNegativeLayerWarning(t *testing.T) {
	c := controls.NewPanel().
		ID("neg_layer").
		Layer(-5).
		Control()
	r := validate.ValidateControl(c)
	// Layer < 0 does not trigger a warning in current code (only > 100 does).
	// This test documents the behavior: negative layer should not error.
	if len(r.Errors) != 0 {
		t.Errorf("expected no errors for negative layer, got: %s", r.String())
	}
}

func TestValidateMissingBindingTargetPropertyName(t *testing.T) {
	c := controls.NewPanel().
		ID("bind_target").
		Bindings(schema.Bind{
			BindingType:      schema.BindingTypeView,
			BindingName:      "test_bind",
			SourcePropertyName: "source_prop",
			// TargetPropertyName missing
		}).
		Control()
	r := validate.ValidateControl(c)
	if len(r.Errors) == 0 {
		t.Error("expected error for view binding missing target_property_name")
	}
}

func TestValidateAnimWithoutDurationWarning(t *testing.T) {
	c := controls.NewPanel().
		ID("anim_dur").
		Animations(schema.Anim{
			AnimType: schema.AnimTypeAlpha,
			// Duration is 0 (default)
		}).
		Control()
	r := validate.ValidateControl(c)
	if len(r.Warnings) == 0 {
		t.Error("expected warning for animation with zero duration")
	}
}
