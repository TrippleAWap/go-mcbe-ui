package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
)

func TestRaw(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want string
	}{
		{"pixels", []interface{}{100, 200}, "[100,200]"},
		{"strings", []interface{}{"fill", "parent"}, `["fill","parent"]`},
		{"mixed", []interface{}{100, "100%"}, "[100,\"100%\"]"},
		{"mixed reversed", []interface{}{"100%", 50}, `["100%",50]`},
		{"content relative", []interface{}{"100%c", "100%cm"}, `["100%c","100%cm"]`},
		{"cross axis", []interface{}{"100%x", "100%y"}, `["100%x","100%y"]`},
		{"expression", []interface{}{100, "100% - 4px"}, "[100,\"100% - 4px\"]"},
		{"default size", []interface{}{"default", "default"}, `["default","default"]`},
		{"int pixel", []interface{}{1, 2}, "[1,2]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(schema.MustJSON(schema.Raw(tt.args...)))
			if got != tt.want {
				t.Errorf("Raw(%v) = %s, want %s", tt.args, got, tt.want)
			}
		})
	}
}

func TestSingleAxis(t *testing.T) {
	r := schema.Pixel(200)
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(b) != "[200]" {
		t.Errorf("Pixel(200) = %s, want [200]", string(b))
	}

	r2 := schema.Percent(75)
	b2, err := json.Marshal(r2)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(b2) != `["75%"]` {
		t.Errorf("Percent(75) = %s, want [\"75%%\"]", string(b2))
	}

	r3 := schema.Named("fill")
	b3, err := json.Marshal(r3)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(b3) != `["fill"]` {
		t.Errorf("Named(\"fill\") = %s, want [\"fill\"]", string(b3))
	}
}

func TestChaining(t *testing.T) {
	r := schema.Pixel(100).SetY(schema.Named("fill"))
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(b) != `[100,"fill"]` {
		t.Errorf("chained = %s, want [100,\"fill\"]", string(b))
	}
}

func TestEmpty(t *testing.T) {
	var r schema.Relative
	if !r.IsEmpty() {
		t.Error("expected empty Relative")
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(b) != "null" {
		t.Errorf("empty Relative should marshal as null, got %s", string(b))
	}
}

func TestPixelAccessors(t *testing.T) {
	r := schema.Raw(42, 99)
	if r.PixelX() != 42 {
		t.Errorf("PixelX() = %v, want 42", r.PixelX())
	}
	if r.PixelY() != 99 {
		t.Errorf("PixelY() = %v, want 99", r.PixelY())
	}

	r2 := schema.Raw("100%", "50%")
	if r2.PixelX() != 0 || r2.PixelY() != 0 {
		t.Error("non-pixel axes should return 0")
	}

	r3 := schema.Raw(100, "100% - 4px")
	if r3.PixelX() != 100 || r3.PixelY() != 0 {
		t.Errorf("expected [100, 0], got [%v, %v]", r3.PixelX(), r3.PixelY())
	}
}

func TestOmitEmpty(t *testing.T) {
	c := schema.Control{ID: "empty", Type: "panel"}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if _, ok := obj["size"]; ok {
		t.Error("size should be omitted when empty")
	}
	if _, ok := obj["offset"]; ok {
		t.Error("offset should be omitted when empty")
	}
}

func TestControlMarshal(t *testing.T) {
	c := schema.Control{
		ID:       "ctrl",
		Type:     "screen",
		IsModal:  true,
		Layer:    5,
		ZOrder:   10,
		Renderer: "live_player_renderer",
		Size:     schema.PtrRel(schema.Raw("100%", "100%")),
		Offset:   schema.PtrRel(schema.Raw(0, 0)),
	}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if obj["is_modal"] != true {
		t.Error("expected is_modal true")
	}
	if obj["renderer"] != "live_player_renderer" {
		t.Errorf("expected renderer, got %v", obj["renderer"])
	}
	if obj["layer"] != float64(5) {
		t.Errorf("expected layer 5, got %v", obj["layer"])
	}
	if obj["z_order"] != float64(10) {
		t.Errorf("expected z_order 10, got %v", obj["z_order"])
	}
	size, ok := obj["size"].([]interface{})
	if !ok || len(size) != 2 || size[0] != "100%" || size[1] != "100%" {
		t.Errorf("expected size [\"100%%\",\"100%%\"], got %v", size)
	}
	// Zero-value Vector2 fields should be omitted
	if _, ok := obj["uv"]; ok {
		t.Error("uv should be omitted when zero")
	}
}

func TestControlMarshalWithRawFields(t *testing.T) {
	c := &schema.Control{
		ID:   "raw_test",
		Type: "panel",
	}
	c.Raw = map[string]json.RawMessage{
		"custom_field": json.RawMessage(`"custom_value"`),
		"nested":       json.RawMessage(`{"key":"val"}`),
	}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if obj["custom_field"] != "custom_value" {
		t.Errorf("custom_field = %v, want custom_value", obj["custom_field"])
	}
	nested, ok := obj["nested"].(map[string]interface{})
	if !ok {
		t.Fatalf("nested is not an object, got %T", obj["nested"])
	}
	if nested["key"] != "val" {
		t.Errorf("nested.key = %v, want val", nested["key"])
	}
	// Ensure existing fields are still present
	if obj["id"] != "raw_test" {
		t.Errorf("id = %v, want raw_test", obj["id"])
	}
}

func TestControlMarshalRawOverwritesExisting(t *testing.T) {
	// Raw fields should override marshaled standard fields.
	c := &schema.Control{
		ID:   "override",
		Type: "panel",
	}
	c.Raw = map[string]json.RawMessage{
		"type": json.RawMessage(`"custom_type"`),
	}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if obj["type"] != "custom_type" {
		t.Errorf("type = %v, want custom_type", obj["type"])
	}
}

func TestNewScreenDefaults(t *testing.T) {
	s := schema.NewScreen("test_ns")
	if s.Namespace != "test_ns" {
		t.Errorf("Namespace = %s, want test_ns", s.Namespace)
	}
	if s.RootPanel != "root" {
		t.Errorf("RootPanel = %s, want root", s.RootPanel)
	}
	if s.Root.ID != "root" {
		t.Errorf("Root.ID = %s, want root", s.Root.ID)
	}
	if s.Root.Type != "panel" {
		t.Errorf("Root.Type = %s, want panel", s.Root.Type)
	}
	if s.Root.AnchorFrom != schema.Center {
		t.Errorf("AnchorFrom = %s, want center", s.Root.AnchorFrom)
	}
	if s.Root.AnchorTo != schema.Center {
		t.Errorf("AnchorTo = %s, want center", s.Root.AnchorTo)
	}
	if s.Root.Size == nil {
		t.Fatal("Root.Size should not be nil")
	}
	if s.Root.Size.PixelX() != 1920 || s.Root.Size.PixelY() != 1080 {
		t.Errorf("Root.Size = [%d, %d], want [1920, 1080]",
			s.Root.Size.PixelX(), s.Root.Size.PixelY())
	}
}

func TestRelativeMixedPixelStringAxes(t *testing.T) {
	r := schema.Raw(100, "50%")
	if r.PixelX() != 100 {
		t.Errorf("PixelX() = %d, want 100", r.PixelX())
	}
	if r.PixelY() != 0 {
		t.Errorf("PixelY() = %d, want 0", r.PixelY())
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(b) != `[100,"50%"]` {
		t.Errorf("marshaled = %s, want [100,\"50%%\"]", string(b))
	}
}

func TestRelativeSingleAxisYOnly(t *testing.T) {
	// Single-axis with only Y set should marshal as [y].
	r := schema.Raw("fill")
	r = r.SetY(schema.Pixel(200))
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(b) != `["fill",200]` {
		t.Errorf("mixed single-axis = %s, want [\"fill\",200]", string(b))
	}
}

func TestAnimCopyIndependence(t *testing.T) {
	a1 := schema.AnimAlpha(1.0, 0.0, 1.0)
	a2 := a1.Copy()
	if a1.AnimType != a2.AnimType {
		t.Errorf("copy AnimType = %s, want %s", a2.AnimType, a1.AnimType)
	}
	if a1.Duration != a2.Duration {
		t.Errorf("copy Duration = %f, want %f", a2.Duration, a1.Duration)
	}
	// Mutating copy should not affect original.
	a2.Duration = 5.0
	if a1.Duration == 5.0 {
		t.Error("Anim.Copy() should return an independent copy; mutating copy affected original")
	}
}

func TestVector2MarshalUnmarshal(t *testing.T) {
	v := schema.Vector2{X: 3.5, Y: 7.2}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var v2 schema.Vector2
	if err := json.Unmarshal(b, &v2); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if v2.X != 3.5 || v2.Y != 7.2 {
		t.Errorf("round-trip = {%v, %v}, want {3.5, 7.2}", v2.X, v2.Y)
	}
}

func TestVector2ZeroMarshalOmitsFields(t *testing.T) {
	// A nil *Vector2 should omit both x and y in a Control.
	c := schema.Control{ID: "v", Type: "image"}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if _, ok := obj["x"]; ok {
		t.Error("zero Vector2 X should be omitted")
	}
	if _, ok := obj["y"]; ok {
		t.Error("zero Vector2 Y should be omitted")
	}
}

func TestAllAnchorsIsValid(t *testing.T) {
	anchors := []schema.Anchor{
		schema.TopLeft, schema.TopCenter, schema.TopRight, schema.TopMiddle,
		schema.CenterLeft, schema.Center, schema.CenterRight,
		schema.BottomLeft, schema.BottomCenter, schema.BottomRight, schema.BottomMiddle,
		schema.MiddleLeft, schema.MiddleCenter, schema.MiddleRight,
	}
	for _, a := range anchors {
		if !a.IsValid() {
			t.Errorf("anchor %q should be valid", a)
		}
	}
	if schema.Anchor("unknown").IsValid() {
		t.Error("unknown anchor should be invalid")
	}
}

func TestRelativeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		json string
		want string // expected marshaled output
	}{
		{`[100, 200]`, "[100,200]"},
		{`["fill", "parent"]`, `["fill","parent"]`},
		{`[100, "75%"]`, "[100,\"75%\"]"},
		{`["default"]`, `["default"]`},
		{`null`, "null"},
	}
	for _, tt := range tests {
		t.Run(tt.json, func(t *testing.T) {
			var r schema.Relative
			if err := json.Unmarshal([]byte(tt.json), &r); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			// Round-trip: marshal should produce the same JSON.
			b, err := json.Marshal(r)
			if err != nil {
				t.Fatalf("marshal failed: %v", err)
			}
			if string(b) != tt.want {
				t.Errorf("marshal = %s, want %s", string(b), tt.want)
			}
		})
	}
}

func TestRawWithFloat64(t *testing.T) {
	r := schema.Raw(1.5, 2.5)
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(b) != "[1.5,2.5]" {
		t.Errorf("got %s, want [1.5,2.5]", string(b))
	}
}

func TestRawPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected Raw(0 args) to panic")
		}
	}()
	schema.Raw()
}

func TestRawPanicThreeArgs(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected Raw(1,2,3) to panic")
		}
	}()
	schema.Raw(1, 2, 3)
}

func TestAllEasingsAreValid(t *testing.T) {
	easings := []schema.Easing{
		schema.EasingLinear, schema.EasingInQuad, schema.EasingOutQuad, schema.EasingInOutQuad,
		schema.EasingInCubic, schema.EasingOutCubic, schema.EasingInOutCubic,
		schema.EasingInQuart, schema.EasingOutQuart, schema.EasingInOutQuart,
		schema.EasingInQuint, schema.EasingOutQuint, schema.EasingInOutQuint,
		schema.EasingInSine, schema.EasingOutSine, schema.EasingInOutSine,
		schema.EasingInExpo, schema.EasingOutExpo, schema.EasingInOutExpo,
		schema.EasingInCirc, schema.EasingOutCirc, schema.EasingInOutCirc,
		schema.EasingInBack, schema.EasingOutBack, schema.EasingInOutBack,
		schema.EasingInElastic, schema.EasingOutElastic, schema.EasingInOutElastic,
		schema.EasingInBounce, schema.EasingOutBounce, schema.EasingInOutBounce,
	}
	for _, e := range easings {
		if string(e) == "" {
			t.Errorf("expected non-empty easing name")
		}
	}
}
