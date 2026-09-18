package controls_test

import (
	"encoding/json"
	"testing"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/controls"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
)

func TestPanelBuilder(t *testing.T) {
	p := controls.NewPanel().
		ID("test_panel").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100)).
		Layer(5).
		Visible(false)

	c := p.Control()
	if c.ID != "test_panel" {
		t.Errorf("ID = %s, want test_panel", c.ID)
	}
	if c.Type != "panel" {
		t.Errorf("Type = %s, want panel", c.Type)
	}
	if c.Layer != 5 {
		t.Errorf("Layer = %d, want 5", c.Layer)
	}
	if c.Visible {
		t.Error("Visible should be false")
	}
}

func TestLabelBuilder(t *testing.T) {
	l := controls.NewLabel().
		ID("msg").
		Anchor(schema.TopLeft).
		Size(schema.Raw(200, 40)).
		Text("Hello World").
		FontSize(schema.FontSizeNormal).
		Localize()

	c := l.Control()
	if c.Text != "Hello World" {
		t.Errorf("Text = %s, want Hello World", c.Text)
	}
	if !c.Localize {
		t.Error("Localize should be true")
	}
	if c.FontSize != schema.FontSizeNormal {
		t.Errorf("FontSize = %s, want normal", c.FontSize)
	}
}

func TestImageBuilder(t *testing.T) {
	i := controls.NewImage().
		ID("icon").
		Anchor(schema.TopLeft).
		Size(schema.Raw(32, 32)).
		Texture("textures/ui/test_icon").
		KeepRatio()

	c := i.Control()
	if c.Texture != "textures/ui/test_icon" {
		t.Errorf("Texture = %s, want textures/ui/test_icon", c.Texture)
	}
	if !c.KeepRatio {
		t.Error("KeepRatio should be true")
	}
}

func TestButtonBuilder(t *testing.T) {
	b := controls.NewButton().
		ID("confirm").
		Anchor(schema.Center).
		Size(schema.Raw(150, 50)).
		DefaultControl("confirm_default").
		HoverControl("confirm_hover")

	c := b.Control()
	if c.DefaultControl != "confirm_default" {
		t.Errorf("DefaultControl = %s", c.DefaultControl)
	}
	if c.HoverControl != "confirm_hover" {
		t.Errorf("HoverControl = %s", c.HoverControl)
	}
}

func TestStackPanelBuilder(t *testing.T) {
	sp := controls.NewStackPanel().
		ID("container").
		Anchor(schema.TopLeft).
		Size(schema.Raw(200, 100)).
		Horizontal()

	c := sp.Control()
	if c.Orientation != schema.OrientationHorizontal {
		t.Errorf("Orientation = %s, want horizontal", c.Orientation)
	}
}

func TestScreenBuilder(t *testing.T) {
	s := controls.NewScreen().
		ID("my_screen").
		Anchor(schema.Center).
		Size(schema.Raw(1920, 1080)).
		AbsorbsInput()

	c := s.Control()
	if c.IsModal {
		t.Error("IsModal should be false by default")
	}
	if !c.AbsorbsInput {
		t.Error("AbsorbsInput should be true")
	}
}

func TestJSONOutput(t *testing.T) {
	p := controls.NewPanel().
		ID("root").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100))

	b, err := json.MarshalIndent(p.Control(), "", "  ")
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if obj["id"] != "root" {
		t.Errorf("id = %v, want root", obj["id"])
	}
	if obj["type"] != "panel" {
		t.Errorf("type = %v, want panel", obj["type"])
	}
}

func TestAddControl(t *testing.T) {
	root := controls.NewPanel().ID("root")
	child := controls.NewPanel().ID("child")

	root.AddControl(child.Control().ID)
	if len(root.Control().Controls) != 1 {
		t.Error("expected child control to be added")
	}
}

func TestRawField(t *testing.T) {
	p := controls.NewPanel().
		ID("raw_test").
		Raw("custom_prop", map[string]interface{}{"key": "value"})

	b, err := json.Marshal(p.Control())
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(b) == "" {
		t.Error("expected non-empty JSON output")
	}
}

func TestBuildAlias(t *testing.T) {
	p := controls.NewPanel().ID("alias_test")
	c1 := p.Control()
	c2 := p.Control()
	if c1.ID != c2.ID {
		t.Error("Control() and Build() should return same control")
	}
}

func TestLayoutHelpers(t *testing.T) {
	tests := []struct {
		name string
		got  schema.Relative
		want string
	}{
		{"pad", schema.Pad(10, 20), "[10,20]"},
		{"margin", schema.Margin(5, 5), "[5,5]"},
		{"center_in", schema.CenterIn(800, 600, 400, 300), "[200,150]"},
		{"half_size", schema.HalfSize(100, 80), "[50,40]"},
		{"fill", schema.Fill(), `["fill","fill"]`},
		{"fill_x", schema.FillX(), `["fill","parent"]`},
		{"fill_y", schema.FillY(), `["parent","fill"]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.got)
			if err != nil {
				t.Fatalf("marshal failed: %v", err)
			}
			if string(b) != tt.want {
				t.Errorf("got %s, want %s", string(b), tt.want)
			}
		})
	}
}

func TestScrollViewScrollbarAlwaysVisible(t *testing.T) {
	s := controls.NewScrollView().
		ID("scroll").
		Anchor(schema.TopLeft).
		Size(schema.Raw(200, 200)).
		ScrollbarAlwaysVisible()

	c := s.Control()
	if !c.ScrollbarAlwaysVisible {
		t.Error("ScrollbarAlwaysVisible should be true")
	}
}

func TestToggleBuilder(t *testing.T) {
	tg := controls.NewToggle().
		ID("toggle1").
		Name("my_toggle").
		RadioGroup("radio_grp").
		OnButton("on_btn").
		OffButton("off_btn").
		CheckedControl("checked_ctrl").
		UncheckedControl("unchecked_ctrl").
		CheckedHoverControl("checked_hover").
		UncheckedHoverControl("unchecked_hover")

	c := tg.Control()
	if c.ToggleName != "my_toggle" {
		t.Errorf("ToggleName = %s, want my_toggle", c.ToggleName)
	}
	if c.RadioToggleGroup != "radio_grp" {
		t.Errorf("RadioToggleGroup = %s, want radio_grp", c.RadioToggleGroup)
	}
	if c.ToggleOnButton != "on_btn" {
		t.Errorf("ToggleOnButton = %s, want on_btn", c.ToggleOnButton)
	}
	if c.ToggleOffButton != "off_btn" {
		t.Errorf("ToggleOffButton = %s, want off_btn", c.ToggleOffButton)
	}
	if c.CheckedControl != "checked_ctrl" {
		t.Errorf("CheckedControl = %s, want checked_ctrl", c.CheckedControl)
	}
	if c.UncheckedControl != "unchecked_ctrl" {
		t.Errorf("UncheckedControl = %s, want unchecked_ctrl", c.UncheckedControl)
	}
	if c.CheckedHoverControl != "checked_hover" {
		t.Errorf("CheckedHoverControl = %s, want checked_hover", c.CheckedHoverControl)
	}
	if c.UncheckedHoverControl != "unchecked_hover" {
		t.Errorf("UncheckedHoverControl = %s, want unchecked_hover", c.UncheckedHoverControl)
	}
}

func TestSliderBuilder(t *testing.T) {
	sl := controls.NewSlider().
		ID("slider1").
		TrackButton("track_btn").
		SelectedButton("selected_btn").
		DeselectedButton("deselected_btn").
		Steps(10)

	c := sl.Control()
	if c.SliderTrackButton != "track_btn" {
		t.Errorf("SliderTrackButton = %s, want track_btn", c.SliderTrackButton)
	}
	if c.SliderSelectedButton != "selected_btn" {
		t.Errorf("SliderSelectedButton = %s, want selected_btn", c.SliderSelectedButton)
	}
	if c.SliderDeselectedButton != "deselected_btn" {
		t.Errorf("SliderDeselectedButton = %s, want deselected_btn", c.SliderDeselectedButton)
	}
	if c.SliderSteps != 10 {
		t.Errorf("SliderSteps = %d, want 10", c.SliderSteps)
	}
}

func TestEditBoxBuilder(t *testing.T) {
	eb := controls.NewEditBox().
		ID("edit1").
		BoxName("my_box").
		TextType("numeric").
		MaxLength(42).
		Multiline()

	c := eb.Control()
	if c.TextBoxName != "my_box" {
		t.Errorf("TextBoxName = %s, want my_box", c.TextBoxName)
	}
	if c.TextType != "numeric" {
		t.Errorf("TextType = %s, want numeric", c.TextType)
	}
	if c.MaxLength != 42 {
		t.Errorf("MaxLength = %d, want 42", c.MaxLength)
	}
	if !c.EnabledNewline {
		t.Error("Multiline should set EnabledNewline to true")
	}
}

func TestGridBuilder(t *testing.T) {
	g := controls.NewGrid().
		ID("grid1").
		Dimensions(3, 4).
		MaxItems(24).
		ItemTemplate("item_template").
		FillDirection(schema.GridFillDown)

	c := g.Control()
	if c.GridDimensions == nil {
		t.Fatal("GridDimensions should not be nil")
	}
	if c.GridDimensions.X != 3 || c.GridDimensions.Y != 4 {
		t.Errorf("GridDimensions = {%v, %v}, want {3, 4}", c.GridDimensions.X, c.GridDimensions.Y)
	}
	if c.MaximumGridItems != 24 {
		t.Errorf("MaximumGridItems = %d, want 24", c.MaximumGridItems)
	}
	if c.GridItemTemplate != "item_template" {
		t.Errorf("GridItemTemplate = %s, want item_template", c.GridItemTemplate)
	}
	if c.GridFillDirection != schema.GridFillDown {
		t.Errorf("GridFillDirection = %s, want down", c.GridFillDirection)
	}
}

func TestCustomBuilder(t *testing.T) {
	cst := controls.NewCustom().
		ID("custom1").
		Renderer("live_player_renderer")

	ctrl := cst.Control()
	if ctrl.Renderer != "live_player_renderer" {
		t.Errorf("Renderer = %s, want live_player_renderer", ctrl.Renderer)
	}
	if ctrl.Type != "custom" {
		t.Errorf("Type = %s, want custom", ctrl.Type)
	}
}

func TestCollectionPanelBuilder(t *testing.T) {
	cp := controls.NewCollectionPanel().
		ID("coll1").
		Collection("my_collection")

	c := cp.Control()
	if c.CollectionName != "my_collection" {
		t.Errorf("CollectionName = %s, want my_collection", c.CollectionName)
	}
	if c.Type != "collection_panel" {
		t.Errorf("Type = %s, want collection_panel", c.Type)
	}
}

func TestDropdownBuilderNameRawField(t *testing.T) {
	dd := controls.NewDropdown().
		ID("dropdown1").
		Name("my_dropdown")

	c := dd.Control()
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if obj["dropdown_name"] != "my_dropdown" {
		t.Errorf("dropdown_name = %v, want my_dropdown", obj["dropdown_name"])
	}
}

func TestSelectionWheelBuilderSlicesRawField(t *testing.T) {
	sw := controls.NewSelectionWheel().
		ID("wheel1").
		Slices(5)

	c := sw.Control()
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if obj["slice_count"] != float64(5) {
		t.Errorf("slice_count = %v, want 5", obj["slice_count"])
	}
}

func TestBuildAndControlReturnIdenticalPointer(t *testing.T) {
	// Each control type stores its schema.Control by value inside the builder,
	// so Build() and Control() both return &b.c — the same pointer.
	pan := controls.NewPanel().ID("ptr_test")
	if pan.Control() != pan.Control() {
		t.Error("Panel: Build() and Control() should return the same pointer")
	}

	lab := controls.NewLabel().ID("ptr_test")
	if lab.Control() != lab.Control() {
		t.Error("Label: Build() and Control() should return the same pointer")
	}

	img := controls.NewImage().ID("ptr_test")
	if img.Control() != img.Control() {
		t.Error("Image: Build() and Control() should return the same pointer")
	}

	btn := controls.NewButton().ID("ptr_test")
	if btn.Control() != btn.Control() {
		t.Error("Button: Build() and Control() should return the same pointer")
	}

	sp := controls.NewStackPanel().ID("ptr_test")
	if sp.Control() != sp.Control() {
		t.Error("StackPanel: Build() and Control() should return the same pointer")
	}

	grd := controls.NewGrid().ID("ptr_test")
	if grd.Control() != grd.Control() {
		t.Error("Grid: Build() and Control() should return the same pointer")
	}

	scr := controls.NewScreen().ID("ptr_test")
	if scr.Control() != scr.Control() {
		t.Error("Screen: Build() and Control() should return the same pointer")
	}

	tgl := controls.NewToggle().ID("ptr_test")
	if tgl.Control() != tgl.Control() {
		t.Error("Toggle: Build() and Control() should return the same pointer")
	}

	sld := controls.NewSlider().ID("ptr_test")
	if sld.Control() != sld.Control() {
		t.Error("Slider: Build() and Control() should return the same pointer")
	}

	eb := controls.NewEditBox().ID("ptr_test")
	if eb.Control() != eb.Control() {
		t.Error("EditBox: Build() and Control() should return the same pointer")
	}

	sv := controls.NewScrollView().ID("ptr_test")
	if sv.Control() != sv.Control() {
		t.Error("ScrollView: Build() and Control() should return the same pointer")
	}

	cst := controls.NewCustom().ID("ptr_test")
	if cst.Control() != cst.Control() {
		t.Error("Custom: Build() and Control() should return the same pointer")
	}

	cp := controls.NewCollectionPanel().ID("ptr_test")
	if cp.Control() != cp.Control() {
		t.Error("CollectionPanel: Build() and Control() should return the same pointer")
	}

	dd := controls.NewDropdown().ID("ptr_test")
	if dd.Control() != dd.Control() {
		t.Error("Dropdown: Build() and Control() should return the same pointer")
	}

	sw := controls.NewSelectionWheel().ID("ptr_test")
	if sw.Control() != sw.Control() {
		t.Error("SelectionWheel: Build() and Control() should return the same pointer")
	}
}

func TestPanelAllMethods(t *testing.T) {
	p := controls.NewPanel().
		ID("all_panel").
		Anchor(schema.TopLeft).
		AnchorFrom(schema.TopCenter).
		AnchorTo(schema.BottomRight).
		Pos(schema.Raw(10, 20)).
		Size(schema.Raw(100, 200)).
		MaxSize(schema.Raw(500, 400)).
		MinSize(schema.Raw(50, 30)).
		Visible(false).
		Enabled(false).
		Layer(3).
		Alpha(0.75).
		ClipsChildren().
		Contained().
		AddControl("child1").
		Bindings(schema.Bind{BindingType: schema.BindingTypeView, SourcePropertyName: "src", TargetPropertyName: "tgt"}).
		Animations(schema.AnimAlpha(0.5, 0, 1)).
		Raw("custom_key", "custom_val")

	c := p.Control()
	if c.AnchorFrom != schema.TopCenter {
		t.Errorf("AnchorFrom = %q, want top_center", c.AnchorFrom)
	}
	if c.AnchorTo != schema.BottomRight {
		t.Errorf("AnchorTo = %q, want bottom_right", c.AnchorTo)
	}
	if c.Offset == nil || c.Offset.PixelX() != 10 || c.Offset.PixelY() != 20 {
		t.Errorf("Offset = %v, want [10, 20]", c.Offset)
	}
	if c.Size == nil || c.Size.PixelX() != 100 || c.Size.PixelY() != 200 {
		t.Errorf("Size = %v, want [100, 200]", c.Size)
	}
	if c.MaxSize == nil || c.MaxSize.PixelX() != 500 {
		t.Errorf("MaxSize = %v, want [500, 400]", c.MaxSize)
	}
	if c.MinSize == nil || c.MinSize.PixelY() != 30 {
		t.Errorf("MinSize = %v, want [50, 30]", c.MinSize)
	}
	if c.Visible || c.Enabled || c.Layer != 3 || c.Alpha != 0.75 {
		t.Errorf("Visibility/layer/alpha mismatch")
	}
	if !c.ClipsChildren || !c.Contained {
		t.Error("ClipsChildren/Contained should be true")
	}
	if len(c.Controls) != 1 || c.Controls[0] != "child1" {
		t.Errorf("Controls = %v, want [child1]", c.Controls)
	}
	if len(c.Bindings) != 1 {
		t.Errorf("Bindings length = %d, want 1", len(c.Bindings))
	}
	if len(c.Anims) != 1 {
		t.Errorf("Anims length = %d, want 1", len(c.Anims))
	}
	if c.Raw == nil || c.Raw["custom_key"] == nil {
		t.Error("Raw field missing custom_key")
	}
}

func TestLabelAllMethods(t *testing.T) {
	l := controls.NewLabel().
		ID("all_label").
		Anchor(schema.BottomLeft).
		Pos(schema.Raw(5, 10)).
		Size(schema.Raw(200, 40)).
		Text("Hello MCBE").
		FontSize(schema.FontSizeLarge).
		FontScale(1.5).
		FontType(schema.FontTypeUniform).
		Align(schema.TextAlignRight).
		Localize().
		Shadow().
		LinePadding(2.0).
		Visible(false).
		Layer(1).
		Alpha(0.8).
		Bindings(schema.Bind{BindingType: schema.BindingTypeView, SourcePropertyName: "a", TargetPropertyName: "b"})

	c := l.Control()
	if c.Text != "Hello MCBE" {
		t.Errorf("Text = %q, want Hello MCBE", c.Text)
	}
	if c.FontSize != schema.FontSizeLarge {
		t.Errorf("FontSize = %q, want large", c.FontSize)
	}
	if c.FontScaleFactor != 1.5 {
		t.Errorf("FontScaleFactor = %v, want 1.5", c.FontScaleFactor)
	}
	if c.FontType != schema.FontTypeUniform {
		t.Errorf("FontType = %q, want uniform", c.FontType)
	}
	if c.TextAlignment != schema.TextAlignRight {
		t.Errorf("TextAlignment = %q, want right", c.TextAlignment)
	}
	if !c.Localize || !c.Shadow || c.LinePadding != 2.0 {
		t.Error("Localize/Shadow/LinePadding mismatch")
	}
}

func TestImageAllMethods(t *testing.T) {
	i := controls.NewImage().
		ID("all_image").
		Anchor(schema.Center).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(64, 64)).
		Texture("textures/ui/test_icon").
		UV(0.0, 0.0).
		UVSize(32, 32).
		NineSlice(4, 4, 4, 4).
		Color(1.0, 0.5, 0.25).
		KeepRatio().
		Bilinear().
		Fill().
		Grayscale().
		Tiled().
		TiledScale(1.5, 2.0).
		ClipDirection(schema.ClipDirectionBoth).
		ClipRatio(0.5).
		Visible(false).
		Layer(2).
		Alpha(0.9).
		Raw("extra", "value")

	c := i.Control()
	if c.Texture != "textures/ui/test_icon" {
		t.Errorf("Texture = %q, want textures/ui/test_icon", c.Texture)
	}
	if c.UV == nil || c.UV.X != 0 || c.UV.Y != 0 {
		t.Errorf("UV = %v, want [0, 0]", c.UV)
	}
	if c.UVSize == nil || c.UVSize.X != 32 {
		t.Errorf("UVSize = %v, want [32, 32]", c.UVSize)
	}
	if !c.KeepRatio || !c.Bilinear || !c.Fill || !c.Grayscale || !c.Tiled {
		t.Error("image flags mismatch")
	}
	if c.TiledScale == nil || c.TiledScale.X != 1.5 {
		t.Errorf("TiledScale = %v, want [1.5, 2.0]", c.TiledScale)
	}
	if c.ClipDirection != schema.ClipDirectionBoth {
		t.Errorf("ClipDirection = %q, want both", c.ClipDirection)
	}
	if c.ClipRatio != 0.5 {
		t.Errorf("ClipRatio = %v, want 0.5", c.ClipRatio)
	}
}

func TestButtonAllMethods(t *testing.T) {
	b := controls.NewButton().
		ID("all_btn").
		Anchor(schema.TopMiddle).
		Pos(schema.Raw(0, 10)).
		Size(schema.Raw(150, 50)).
		DefaultControl("btn_default").
		HoverControl("btn_hover").
		PressedControl("btn_pressed").
		LockedControl("btn_locked")

	c := b.Control()
	if c.DefaultControl != "btn_default" {
		t.Errorf("DefaultControl = %q", c.DefaultControl)
	}
	if c.HoverControl != "btn_hover" {
		t.Errorf("HoverControl = %q", c.HoverControl)
	}
	if c.PressedControl != "btn_pressed" {
		t.Errorf("PressedControl = %q", c.PressedControl)
	}
	if c.LockedControl != "btn_locked" {
		t.Errorf("LockedControl = %q", c.LockedControl)
	}
}

func TestStackPanelAllMethods(t *testing.T) {
	s := controls.NewStackPanel().
		ID("all_stack").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(200, 400)).
		Horizontal().
		Visible(false).
		Layer(1)

	c := s.Control()
	if c.Orientation != schema.OrientationHorizontal {
		t.Errorf("Orientation = %q, want horizontal", c.Orientation)
	}
}

func TestGridAllMethods(t *testing.T) {
	g := controls.NewGrid().
		ID("all_grid").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(300, 300)).
		Dimensions(3, 4).
		MaxItems(24).
		ItemTemplate("grid_item").
		FillDirection(schema.GridFillDown).
		Visible(false).
		Layer(1)

	c := g.Control()
	if c.GridDimensions == nil || c.GridDimensions.X != 3 || c.GridDimensions.Y != 4 {
		t.Errorf("GridDimensions = %v, want {3, 4}", c.GridDimensions)
	}
	if c.MaximumGridItems != 24 {
		t.Errorf("MaximumGridItems = %d, want 24", c.MaximumGridItems)
	}
	if c.GridItemTemplate != "grid_item" {
		t.Errorf("GridItemTemplate = %q, want grid_item", c.GridItemTemplate)
	}
	if c.GridFillDirection != schema.GridFillDown {
		t.Errorf("GridFillDirection = %q, want down", c.GridFillDirection)
	}
}

func TestScreenAllMethods(t *testing.T) {
	s := controls.NewScreen().
		ID("all_screen").
		Anchor(schema.BottomCenter).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(1920, 1080)).
		Modal().
		AbsorbsInput().
		RenderGameBehind().
		CloseOnPlayerHurt()

	c := s.Control()
	if !c.IsModal || !c.AbsorbsInput || !c.RenderGameBehind || !c.CloseOnPlayerHurt {
		t.Error("screen flags mismatch")
	}
}

func TestToggleAllMethods(t *testing.T) {
	tg := controls.NewToggle().
		ID("all_toggle").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(50, 50)).
		Name("my_toggle").
		RadioGroup("radio_grp").
		OnButton("on_btn").
		OffButton("off_btn").
		CheckedControl("checked").
		UncheckedControl("unchecked")

	c := tg.Control()
	if c.ToggleName != "my_toggle" {
		t.Errorf("ToggleName = %q", c.ToggleName)
	}
	if c.RadioToggleGroup != "radio_grp" {
		t.Errorf("RadioToggleGroup = %q", c.RadioToggleGroup)
	}
	if c.ToggleOnButton != "on_btn" || c.ToggleOffButton != "off_btn" {
		t.Errorf("Toggle buttons mismatch")
	}
	if c.CheckedControl != "checked" || c.UncheckedControl != "unchecked" {
		t.Errorf("Checked/Unchecked controls mismatch")
	}
}

func TestSliderAllMethods(t *testing.T) {
	sl := controls.NewSlider().
		ID("all_slider").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(200, 30)).
		TrackButton("track").
		SelectedButton("selected").
		DeselectedButton("deselected").
		Steps(10)

	c := sl.Control()
	if c.SliderTrackButton != "track" {
		t.Errorf("SliderTrackButton = %q", c.SliderTrackButton)
	}
	if c.SliderSelectedButton != "selected" {
		t.Errorf("SliderSelectedButton = %q", c.SliderSelectedButton)
	}
	if c.SliderDeselectedButton != "deselected" {
		t.Errorf("SliderDeselectedButton = %q", c.SliderDeselectedButton)
	}
	if c.SliderSteps != 10 {
		t.Errorf("SliderSteps = %d, want 10", c.SliderSteps)
	}
}

func TestEditBoxAllMethods(t *testing.T) {
	eb := controls.NewEditBox().
		ID("all_edit").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(200, 30)).
		BoxName("my_box").
		TextType("numeric").
		MaxLength(42).
		Multiline()

	c := eb.Control()
	if c.TextBoxName != "my_box" {
		t.Errorf("TextBoxName = %q", c.TextBoxName)
	}
	if c.TextType != "numeric" {
		t.Errorf("TextType = %q", c.TextType)
	}
	if c.MaxLength != 42 {
		t.Errorf("MaxLength = %d, want 42", c.MaxLength)
	}
	if !c.EnabledNewline {
		t.Error("Multiline should set EnabledNewline")
	}
}

func TestScrollViewAllMethods(t *testing.T) {
	sv := controls.NewScrollView().
		ID("all_scroll").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(300, 400)).
		ViewPort("vp").
		Content("content").
		ScrollbarBox("sb").
		ScrollbarTrack("st").
		ScrollbarAlwaysVisible()

	c := sv.Control()
	if c.ScrollViewPort != "vp" {
		t.Errorf("ScrollViewPort = %q", c.ScrollViewPort)
	}
	if c.ScrollContent != "content" {
		t.Errorf("ScrollContent = %q", c.ScrollContent)
	}
	if c.ScrollbarBox != "sb" {
		t.Errorf("ScrollbarBox = %q", c.ScrollbarBox)
	}
	if c.ScrollbarTrack != "st" {
		t.Errorf("ScrollbarTrack = %q", c.ScrollbarTrack)
	}
	if !c.ScrollbarAlwaysVisible {
		t.Error("ScrollbarAlwaysVisible should be true")
	}
}

func TestCustomAllMethods(t *testing.T) {
	cst := controls.NewCustom().
		ID("all_custom").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(100, 100)).
		Renderer("live_player_renderer")

	c := cst.Control()
	if c.Renderer != "live_player_renderer" {
		t.Errorf("Renderer = %q", c.Renderer)
	}
	if c.Type != "custom" {
		t.Errorf("Type = %q, want custom", c.Type)
	}
}

func TestCollectionPanelAllMethods(t *testing.T) {
	cp := controls.NewCollectionPanel().
		ID("all_collection").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(200, 200)).
		Collection("my_collection")

	c := cp.Control()
	if c.CollectionName != "my_collection" {
		t.Errorf("CollectionName = %q", c.CollectionName)
	}
	if c.Type != "collection_panel" {
		t.Errorf("Type = %q, want collection_panel", c.Type)
	}
}

func TestDropdownBuilder(t *testing.T) {
	dd := controls.NewDropdown().
		ID("all_dropdown").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(150, 30)).
		Name("my_dropdown")

	c := dd.Control()
	b, err := c.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if obj["dropdown_name"] != "my_dropdown" {
		t.Errorf("dropdown_name = %v, want my_dropdown", obj["dropdown_name"])
	}
}

func TestSelectionWheelBuilder(t *testing.T) {
	sw := controls.NewSelectionWheel().
		ID("all_wheel").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(200, 50)).
		Slices(5)

	c := sw.Control()
	b, err := c.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if obj["slice_count"] != float64(5) {
		t.Errorf("slice_count = %v, want 5", obj["slice_count"])
	}
}
