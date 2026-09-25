package controls

import (
	"encoding/json"
	"fmt"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
)

// ctrl is the shared base for all control types. Each method that modifies
// state returns *ctrl so type-gated wrappers can wrap the result.
type ctrl struct {
	c schema.Control
}

func (b *ctrl) Control() *schema.Control { return &b.c }
func (b *ctrl) MustJSON() []byte         { out, _ := json.MarshalIndent(b.c, "", "  "); return out }
func (b *ctrl) JSON() ([]byte, error)    { return json.MarshalIndent(b.c, "", "  ") }

// ---- Layout (shared) ----

func (b *ctrl) pos(v schema.Relative)     { b.c.Offset = schema.PtrRel(v) }
func (b *ctrl) size(v schema.Relative)    { b.c.Size = schema.PtrRel(v) }
func (b *ctrl) maxSize(v schema.Relative) { b.c.MaxSize = schema.PtrRel(v) }
func (b *ctrl) minSize(v schema.Relative) { b.c.MinSize = schema.PtrRel(v) }
func (b *ctrl) addControl(id string)      { b.c.Controls = append(b.c.Controls, id) }

// ---- Visibility / layer (shared) ----

func (b *ctrl) visible(v bool)  { b.c.Visible = v }
func (b *ctrl) enabled(v bool)  { b.c.Enabled = v }
func (b *ctrl) layer(n int)     { b.c.Layer = n }
func (b *ctrl) alpha(a float64) { b.c.Alpha = a }
func (b *ctrl) clipsChildren()  { b.c.ClipsChildren = true }
func (b *ctrl) contained()      { b.c.Contained = true }

// ---- Bindings & animations (shared) ----

func (b *ctrl) bindings(bindings ...schema.Bind) { b.c.Bindings = append(b.c.Bindings, bindings...) }
func (b *ctrl) animations(anims ...schema.Anim)  { b.c.Anims = append(b.c.Anims, anims...) }

// ---- Raw field setter ----

func (b *ctrl) raw(key string, val interface{}) {
	if b.c.Raw == nil {
		b.c.Raw = make(map[string]json.RawMessage)
	}
	raw, err := json.Marshal(val)
	if err != nil {
		panic(fmt.Sprintf("raw field %q: %v", key, err))
	}
	b.c.Raw[key] = raw
}

// =====================================================================
// Panel
// =====================================================================

// Panel builds a panel control.
type Panel struct{ *ctrl }

// NewPanel creates a new Panel builder.
func NewPanel() *Panel {
	return &Panel{&ctrl{c: schema.Control{Type: "panel"}}}
}

func (p *Panel) ID(id string) *Panel               { p.c.ID = id; return p }
func (p *Panel) Anchor(a schema.Anchor) *Panel     { p.c.AnchorFrom = a; p.c.AnchorTo = a; return p }
func (p *Panel) AnchorFrom(a schema.Anchor) *Panel { p.c.AnchorFrom = a; return p }
func (p *Panel) AnchorTo(a schema.Anchor) *Panel   { p.c.AnchorTo = a; return p }
func (p *Panel) Pos(v schema.Relative) *Panel      { p.c.Offset = schema.PtrRel(v); return p }
func (p *Panel) Size(v schema.Relative) *Panel     { p.c.Size = schema.PtrRel(v); return p }
func (p *Panel) MaxSize(v schema.Relative) *Panel  { p.c.MaxSize = schema.PtrRel(v); return p }
func (p *Panel) MinSize(v schema.Relative) *Panel  { p.c.MinSize = schema.PtrRel(v); return p }
func (p *Panel) AddControl(id string) *Panel       { p.c.Controls = append(p.c.Controls, id); return p }
func (p *Panel) Visible(v bool) *Panel             { p.c.Visible = v; return p }
func (p *Panel) Enabled(v bool) *Panel             { p.c.Enabled = v; return p }
func (p *Panel) Layer(n int) *Panel                { p.c.Layer = n; return p }
func (p *Panel) Alpha(a float64) *Panel            { p.c.Alpha = a; return p }
func (p *Panel) ClipsChildren() *Panel             { p.c.ClipsChildren = true; return p }
func (p *Panel) Contained() *Panel                 { p.c.Contained = true; return p }
func (p *Panel) Bindings(bindings ...schema.Bind) *Panel {
	p.c.Bindings = append(p.c.Bindings, bindings...)
	return p
}
func (p *Panel) Animations(anims ...schema.Anim) *Panel {
	p.c.Anims = append(p.c.Anims, anims...)
	return p
}
func (p *Panel) Raw(key string, val interface{}) *Panel {
	if p.c.Raw == nil {
		p.c.Raw = make(map[string]json.RawMessage)
	}
	raw, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	p.c.Raw[key] = raw
	return p
}

// =====================================================================
// Label
// =====================================================================

// Label builds a label control.
type Label struct{ *ctrl }

// NewLabel creates a new Label builder.
func NewLabel() *Label {
	return &Label{&ctrl{c: schema.Control{Type: "label"}}}
}

func (l *Label) ID(id string) *Label                 { l.c.ID = id; return l }
func (l *Label) Anchor(a schema.Anchor) *Label       { l.c.AnchorFrom = a; l.c.AnchorTo = a; return l }
func (l *Label) AnchorFrom(a schema.Anchor) *Label   { l.c.AnchorFrom = a; return l }
func (l *Label) AnchorTo(a schema.Anchor) *Label     { l.c.AnchorTo = a; return l }
func (l *Label) Pos(v schema.Relative) *Label        { l.c.Offset = schema.PtrRel(v); return l }
func (l *Label) Size(v schema.Relative) *Label       { l.c.Size = schema.PtrRel(v); return l }
func (l *Label) MaxSize(v schema.Relative) *Label    { l.c.MaxSize = schema.PtrRel(v); return l }
func (l *Label) MinSize(v schema.Relative) *Label    { l.c.MinSize = schema.PtrRel(v); return l }
func (l *Label) AddControl(id string) *Label         { l.c.Controls = append(l.c.Controls, id); return l }
func (l *Label) Text(t string) *Label                { l.c.Text = t; return l }
func (l *Label) FontSize(fs schema.FontSize) *Label  { l.c.FontSize = fs; return l }
func (l *Label) FontScale(s float64) *Label          { l.c.FontScaleFactor = s; return l }
func (l *Label) FontType(ft schema.FontType) *Label  { l.c.FontType = ft; return l }
func (l *Label) Align(a schema.TextAlignment) *Label { l.c.TextAlignment = a; return l }
func (l *Label) Localize() *Label                    { l.c.Localize = true; return l }
func (l *Label) Shadow() *Label                      { l.c.Shadow = true; return l }
func (l *Label) LinePadding(p float64) *Label        { l.c.LinePadding = p; return l }
func (l *Label) Visible(v bool) *Label               { l.c.Visible = v; return l }
func (l *Label) Layer(n int) *Label                  { l.c.Layer = n; return l }
func (l *Label) Alpha(a float64) *Label              { l.c.Alpha = a; return l }
func (l *Label) Bindings(bindings ...schema.Bind) *Label {
	l.c.Bindings = append(l.c.Bindings, bindings...)
	return l
}
func (l *Label) Animations(anims ...schema.Anim) *Label {
	l.c.Anims = append(l.c.Anims, anims...)
	return l
}
func (l *Label) Raw(key string, val interface{}) *Label {
	if l.c.Raw == nil {
		l.c.Raw = make(map[string]json.RawMessage)
	}
	raw, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	l.c.Raw[key] = raw
	return l
}

// =====================================================================
// Image
// =====================================================================

// Image builds an image control.
type Image struct{ *ctrl }

// NewImage creates a new Image builder.
func NewImage() *Image {
	return &Image{&ctrl{c: schema.Control{Type: "image"}}}
}

func (i *Image) ID(id string) *Image               { i.c.ID = id; return i }
func (i *Image) Anchor(a schema.Anchor) *Image     { i.c.AnchorFrom = a; i.c.AnchorTo = a; return i }
func (i *Image) AnchorFrom(a schema.Anchor) *Image { i.c.AnchorFrom = a; return i }
func (i *Image) AnchorTo(a schema.Anchor) *Image   { i.c.AnchorTo = a; return i }
func (i *Image) Pos(v schema.Relative) *Image      { i.c.Offset = schema.PtrRel(v); return i }
func (i *Image) Size(v schema.Relative) *Image     { i.c.Size = schema.PtrRel(v); return i }
func (i *Image) AddControl(id string) *Image       { i.c.Controls = append(i.c.Controls, id); return i }
func (i *Image) Texture(path string) *Image        { i.c.Texture = path; return i }
func (i *Image) UV(x, y float64) *Image            { v := schema.Vector2{X: x, Y: y}; i.c.UV = &v; return i }
func (i *Image) UVSize(w, h float64) *Image {
	v := schema.Vector2{X: w, Y: h}
	i.c.UVSize = &v
	return i
}
func (i *Image) NineSlice(top, right, bottom, left int) *Image {
	if i.c.Raw == nil {
		i.c.Raw = make(map[string]json.RawMessage)
	}
	i.c.Raw["nineslice_size"] = mustMarshal([]int{top, right, bottom, left})
	return i
}
func (i *Image) Color(r, g, b float64) *Image {
	if i.c.Raw == nil {
		i.c.Raw = make(map[string]json.RawMessage)
	}
	i.c.Raw["color"] = mustMarshal([]float64{r, g, b})
	return i
}
func (i *Image) KeepRatio() *Image { i.c.KeepRatio = true; return i }
func (i *Image) Bilinear() *Image  { i.c.Bilinear = true; return i }
func (i *Image) Fill() *Image      { i.c.Fill = true; return i }
func (i *Image) Grayscale() *Image { i.c.Grayscale = true; return i }
func (i *Image) Tiled() *Image     { i.c.Tiled = true; return i }
func (i *Image) TiledScale(sx, sy float64) *Image {
	v := schema.Vector2{X: sx, Y: sy}
	i.c.TiledScale = &v
	return i
}
func (i *Image) ClipDirection(d schema.ClipDirection) *Image { i.c.ClipDirection = d; return i }
func (i *Image) ClipRatio(r float64) *Image                  { i.c.ClipRatio = r; return i }
func (i *Image) Visible(v bool) *Image                       { i.c.Visible = v; return i }
func (i *Image) Layer(n int) *Image                          { i.c.Layer = n; return i }
func (i *Image) Alpha(a float64) *Image                      { i.c.Alpha = a; return i }
func (i *Image) Bindings(bindings ...schema.Bind) *Image {
	i.c.Bindings = append(i.c.Bindings, bindings...)
	return i
}
func (i *Image) Animations(anims ...schema.Anim) *Image {
	i.c.Anims = append(i.c.Anims, anims...)
	return i
}
func (i *Image) Raw(key string, val interface{}) *Image {
	if i.c.Raw == nil {
		i.c.Raw = make(map[string]json.RawMessage)
	}
	raw, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	i.c.Raw[key] = raw
	return i
}

// =====================================================================
// Button
// =====================================================================

// Button builds a button control.
type Button struct{ *ctrl }

// NewButton creates a new Button builder.
func NewButton() *Button {
	return &Button{&ctrl{c: schema.Control{Type: "button"}}}
}

func (b *Button) ID(id string) *Button               { b.c.ID = id; return b }
func (b *Button) Anchor(a schema.Anchor) *Button     { b.c.AnchorFrom = a; b.c.AnchorTo = a; return b }
func (b *Button) AnchorFrom(a schema.Anchor) *Button { b.c.AnchorFrom = a; return b }
func (b *Button) AnchorTo(a schema.Anchor) *Button   { b.c.AnchorTo = a; return b }
func (b *Button) Pos(v schema.Relative) *Button      { b.c.Offset = schema.PtrRel(v); return b }
func (b *Button) Size(v schema.Relative) *Button     { b.c.Size = schema.PtrRel(v); return b }
func (b *Button) AddControl(id string) *Button       { b.c.Controls = append(b.c.Controls, id); return b }
func (b *Button) DefaultControl(id string) *Button   { b.c.DefaultControl = id; return b }
func (b *Button) HoverControl(id string) *Button     { b.c.HoverControl = id; return b }
func (b *Button) PressedControl(id string) *Button   { b.c.PressedControl = id; return b }
func (b *Button) LockedControl(id string) *Button    { b.c.LockedControl = id; return b }
func (b *Button) Visible(v bool) *Button             { b.c.Visible = v; return b }
func (b *Button) Layer(n int) *Button                { b.c.Layer = n; return b }
func (b *Button) Bindings(bindings ...schema.Bind) *Button {
	b.c.Bindings = append(b.c.Bindings, bindings...)
	return b
}
func (b *Button) Animations(anims ...schema.Anim) *Button {
	b.c.Anims = append(b.c.Anims, anims...)
	return b
}
func (b *Button) Raw(key string, val interface{}) *Button {
	if b.c.Raw == nil {
		b.c.Raw = make(map[string]json.RawMessage)
	}
	raw, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	b.c.Raw[key] = raw
	return b
}

// =====================================================================
// StackPanel
// =====================================================================

// StackPanel builds a stack_panel control.
type StackPanel struct{ *ctrl }

// NewStackPanel creates a new StackPanel builder.
func NewStackPanel() *StackPanel {
	return &StackPanel{&ctrl{c: schema.Control{Type: "stack_panel"}}}
}

func (s *StackPanel) ID(id string) *StackPanel { s.c.ID = id; return s }
func (s *StackPanel) Anchor(a schema.Anchor) *StackPanel {
	s.c.AnchorFrom = a
	s.c.AnchorTo = a
	return s
}
func (s *StackPanel) AnchorFrom(a schema.Anchor) *StackPanel { s.c.AnchorFrom = a; return s }
func (s *StackPanel) AnchorTo(a schema.Anchor) *StackPanel   { s.c.AnchorTo = a; return s }
func (s *StackPanel) Pos(v schema.Relative) *StackPanel      { s.c.Offset = schema.PtrRel(v); return s }
func (s *StackPanel) Size(v schema.Relative) *StackPanel     { s.c.Size = schema.PtrRel(v); return s }
func (s *StackPanel) AddControl(id string) *StackPanel {
	s.c.Controls = append(s.c.Controls, id)
	return s
}
func (s *StackPanel) Horizontal() *StackPanel {
	s.c.Orientation = schema.OrientationHorizontal
	return s
}
func (s *StackPanel) Vertical() *StackPanel      { s.c.Orientation = schema.OrientationVertical; return s }
func (s *StackPanel) Visible(v bool) *StackPanel { s.c.Visible = v; return s }
func (s *StackPanel) Layer(n int) *StackPanel    { s.c.Layer = n; return s }
func (s *StackPanel) Bindings(bindings ...schema.Bind) *StackPanel {
	s.c.Bindings = append(s.c.Bindings, bindings...)
	return s
}
func (s *StackPanel) Animations(anims ...schema.Anim) *StackPanel {
	s.c.Anims = append(s.c.Anims, anims...)
	return s
}
func (s *StackPanel) Raw(key string, val interface{}) *StackPanel {
	if s.c.Raw == nil {
		s.c.Raw = make(map[string]json.RawMessage)
	}
	raw, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	s.c.Raw[key] = raw
	return s
}

// =====================================================================
// Grid
// =====================================================================

// Grid builds a grid control.
type Grid struct{ *ctrl }

// NewGrid creates a new Grid builder.
func NewGrid() *Grid {
	return &Grid{&ctrl{c: schema.Control{Type: "grid"}}}
}

func (g *Grid) ID(id string) *Grid           { g.c.ID = id; return g }
func (g *Grid) Anchor(a schema.Anchor) *Grid { g.c.AnchorFrom = a; g.c.AnchorTo = a; return g }
func (g *Grid) Pos(v schema.Relative) *Grid  { g.c.Offset = schema.PtrRel(v); return g }
func (g *Grid) Size(v schema.Relative) *Grid { g.c.Size = schema.PtrRel(v); return g }
func (g *Grid) AddControl(id string) *Grid   { g.c.Controls = append(g.c.Controls, id); return g }
func (g *Grid) Dimensions(cols, rows int) *Grid {
	v := schema.Vector2{X: float64(cols), Y: float64(rows)}
	g.c.GridDimensions = &v
	return g
}
func (g *Grid) MaxItems(n int) *Grid                           { g.c.MaximumGridItems = n; return g }
func (g *Grid) ItemTemplate(id string) *Grid                   { g.c.GridItemTemplate = id; return g }
func (g *Grid) FillDirection(d schema.GridFillDirection) *Grid { g.c.GridFillDirection = d; return g }
func (g *Grid) Visible(v bool) *Grid                           { g.c.Visible = v; return g }
func (g *Grid) Layer(n int) *Grid                              { g.c.Layer = n; return g }
func (g *Grid) Bindings(bindings ...schema.Bind) *Grid {
	g.c.Bindings = append(g.c.Bindings, bindings...)
	return g
}
func (g *Grid) Animations(anims ...schema.Anim) *Grid {
	g.c.Anims = append(g.c.Anims, anims...)
	return g
}
func (g *Grid) Raw(key string, val interface{}) *Grid {
	if g.c.Raw == nil {
		g.c.Raw = make(map[string]json.RawMessage)
	}
	raw, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	g.c.Raw[key] = raw
	return g
}

// =====================================================================
// Screen
// =====================================================================

// Screen builds a screen control.
type Screen struct{ *ctrl }

// NewScreen creates a new Screen builder.
func NewScreen() *Screen {
	return &Screen{&ctrl{c: schema.Control{Type: "screen"}}}
}

func (s *Screen) ID(id string) *Screen           { s.c.ID = id; return s }
func (s *Screen) Anchor(a schema.Anchor) *Screen { s.c.AnchorFrom = a; s.c.AnchorTo = a; return s }
func (s *Screen) Pos(v schema.Relative) *Screen  { s.c.Offset = schema.PtrRel(v); return s }
func (s *Screen) Size(v schema.Relative) *Screen { s.c.Size = schema.PtrRel(v); return s }
func (s *Screen) AddControl(id string) *Screen   { s.c.Controls = append(s.c.Controls, id); return s }
func (s *Screen) Modal() *Screen                 { s.c.IsModal = true; return s }
func (s *Screen) AbsorbsInput() *Screen          { s.c.AbsorbsInput = true; return s }
func (s *Screen) RenderGameBehind() *Screen      { s.c.RenderGameBehind = true; return s }
func (s *Screen) CloseOnPlayerHurt() *Screen     { s.c.CloseOnPlayerHurt = true; return s }
func (s *Screen) Bindings(bindings ...schema.Bind) *Screen {
	s.c.Bindings = append(s.c.Bindings, bindings...)
	return s
}
func (s *Screen) Animations(anims ...schema.Anim) *Screen {
	s.c.Anims = append(s.c.Anims, anims...)
	return s
}
func (s *Screen) Raw(key string, val interface{}) *Screen {
	if s.c.Raw == nil {
		s.c.Raw = make(map[string]json.RawMessage)
	}
	raw, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	s.c.Raw[key] = raw
	return s
}

// =====================================================================
// Toggle
// =====================================================================

// Toggle builds a toggle control.
type Toggle struct{ *ctrl }

func NewToggle() *Toggle {
	return &Toggle{&ctrl{c: schema.Control{Type: "toggle"}}}
}

func (t *Toggle) ID(id string) *Toggle                    { t.c.ID = id; return t }
func (t *Toggle) Anchor(a schema.Anchor) *Toggle          { t.c.AnchorFrom = a; t.c.AnchorTo = a; return t }
func (t *Toggle) Pos(v schema.Relative) *Toggle           { t.c.Offset = schema.PtrRel(v); return t }
func (t *Toggle) Size(v schema.Relative) *Toggle          { t.c.Size = schema.PtrRel(v); return t }
func (t *Toggle) Name(name string) *Toggle                { t.c.ToggleName = name; return t }
func (t *Toggle) RadioGroup(group string) *Toggle         { t.c.RadioToggleGroup = group; return t }
func (t *Toggle) OnButton(id string) *Toggle              { t.c.ToggleOnButton = id; return t }
func (t *Toggle) OffButton(id string) *Toggle             { t.c.ToggleOffButton = id; return t }
func (t *Toggle) CheckedControl(id string) *Toggle        { t.c.CheckedControl = id; return t }
func (t *Toggle) UncheckedControl(id string) *Toggle      { t.c.UncheckedControl = id; return t }
func (t *Toggle) CheckedHoverControl(id string) *Toggle   { t.c.CheckedHoverControl = id; return t }
func (t *Toggle) UncheckedHoverControl(id string) *Toggle { t.c.UncheckedHoverControl = id; return t }

// =====================================================================
// Slider
// =====================================================================

// Slider builds a slider control.
type Slider struct{ *ctrl }

func NewSlider() *Slider {
	return &Slider{&ctrl{c: schema.Control{Type: "slider"}}}
}

func (s *Slider) ID(id string) *Slider               { s.c.ID = id; return s }
func (s *Slider) Anchor(a schema.Anchor) *Slider     { s.c.AnchorFrom = a; s.c.AnchorTo = a; return s }
func (s *Slider) Pos(v schema.Relative) *Slider      { s.c.Offset = schema.PtrRel(v); return s }
func (s *Slider) Size(v schema.Relative) *Slider     { s.c.Size = schema.PtrRel(v); return s }
func (s *Slider) TrackButton(id string) *Slider      { s.c.SliderTrackButton = id; return s }
func (s *Slider) SelectedButton(id string) *Slider   { s.c.SliderSelectedButton = id; return s }
func (s *Slider) DeselectedButton(id string) *Slider { s.c.SliderDeselectedButton = id; return s }
func (s *Slider) Steps(n int) *Slider                { s.c.SliderSteps = n; return s }

// =====================================================================
// EditBox
// =====================================================================

// EditBox builds an edit_box control.
type EditBox struct{ *ctrl }

func NewEditBox() *EditBox {
	return &EditBox{&ctrl{c: schema.Control{Type: "edit_box"}}}
}

func (e *EditBox) ID(id string) *EditBox           { e.c.ID = id; return e }
func (e *EditBox) Anchor(a schema.Anchor) *EditBox { e.c.AnchorFrom = a; e.c.AnchorTo = a; return e }
func (e *EditBox) Pos(v schema.Relative) *EditBox  { e.c.Offset = schema.PtrRel(v); return e }
func (e *EditBox) Size(v schema.Relative) *EditBox { e.c.Size = schema.PtrRel(v); return e }
func (e *EditBox) BoxName(name string) *EditBox    { e.c.TextBoxName = name; return e }
func (e *EditBox) TextType(t string) *EditBox      { e.c.TextType = t; return e }
func (e *EditBox) MaxLength(n int) *EditBox        { e.c.MaxLength = n; return e }
func (e *EditBox) Multiline() *EditBox             { e.c.EnabledNewline = true; return e }

// =====================================================================
// ScrollView
// =====================================================================

// ScrollView builds a scroll_view control.
type ScrollView struct{ *ctrl }

func NewScrollView() *ScrollView {
	return &ScrollView{&ctrl{c: schema.Control{Type: "scroll_view"}}}
}

func (s *ScrollView) ID(id string) *ScrollView { s.c.ID = id; return s }
func (s *ScrollView) Anchor(a schema.Anchor) *ScrollView {
	s.c.AnchorFrom = a
	s.c.AnchorTo = a
	return s
}
func (s *ScrollView) Pos(v schema.Relative) *ScrollView    { s.c.Offset = schema.PtrRel(v); return s }
func (s *ScrollView) Size(v schema.Relative) *ScrollView   { s.c.Size = schema.PtrRel(v); return s }
func (s *ScrollView) ViewPort(id string) *ScrollView       { s.c.ScrollViewPort = id; return s }
func (s *ScrollView) Content(id string) *ScrollView        { s.c.ScrollContent = id; return s }
func (s *ScrollView) ScrollbarBox(id string) *ScrollView   { s.c.ScrollbarBox = id; return s }
func (s *ScrollView) ScrollbarTrack(id string) *ScrollView { s.c.ScrollbarTrack = id; return s }
func (s *ScrollView) ScrollbarAlwaysVisible() *ScrollView {
	s.c.ScrollbarAlwaysVisible = true
	return s
}

// =====================================================================
// Custom
// =====================================================================

// Custom builds a custom control with a renderer.
type Custom struct{ *ctrl }

func NewCustom() *Custom {
	return &Custom{&ctrl{c: schema.Control{Type: "custom"}}}
}

func (c *Custom) ID(id string) *Custom           { c.c.ID = id; return c }
func (c *Custom) Anchor(a schema.Anchor) *Custom { c.c.AnchorFrom = a; c.c.AnchorTo = a; return c }
func (c *Custom) Pos(v schema.Relative) *Custom  { c.c.Offset = schema.PtrRel(v); return c }
func (c *Custom) Size(v schema.Relative) *Custom { c.c.Size = schema.PtrRel(v); return c }
func (c *Custom) Renderer(name string) *Custom   { c.c.Renderer = name; return c }

// =====================================================================
// CollectionPanel
// =====================================================================

// CollectionPanel builds a collection_panel control.
type CollectionPanel struct{ *ctrl }

func NewCollectionPanel() *CollectionPanel {
	return &CollectionPanel{&ctrl{c: schema.Control{Type: "collection_panel"}}}
}

func (cp *CollectionPanel) ID(id string) *CollectionPanel { cp.c.ID = id; return cp }
func (cp *CollectionPanel) Anchor(a schema.Anchor) *CollectionPanel {
	cp.c.AnchorFrom = a
	cp.c.AnchorTo = a
	return cp
}
func (cp *CollectionPanel) Pos(v schema.Relative) *CollectionPanel {
	cp.c.Offset = schema.PtrRel(v)
	return cp
}
func (cp *CollectionPanel) Size(v schema.Relative) *CollectionPanel {
	cp.c.Size = schema.PtrRel(v)
	return cp
}
func (cp *CollectionPanel) Collection(name string) *CollectionPanel {
	cp.c.CollectionName = name
	return cp
}

// =====================================================================
// Dropdown
// =====================================================================

// Dropdown builds a dropdown control.
type Dropdown struct{ *ctrl }

func NewDropdown() *Dropdown {
	return &Dropdown{&ctrl{c: schema.Control{Type: "dropdown"}}}
}

func (d *Dropdown) ID(id string) *Dropdown           { d.c.ID = id; return d }
func (d *Dropdown) Anchor(a schema.Anchor) *Dropdown { d.c.AnchorFrom = a; d.c.AnchorTo = a; return d }
func (d *Dropdown) Pos(v schema.Relative) *Dropdown  { d.c.Offset = schema.PtrRel(v); return d }
func (d *Dropdown) Size(v schema.Relative) *Dropdown { d.c.Size = schema.PtrRel(v); return d }
func (d *Dropdown) Name(name string) *Dropdown {
	if d.c.Raw == nil {
		d.c.Raw = make(map[string]json.RawMessage)
	}
	d.c.Raw["dropdown_name"] = mustMarshal(name)
	return d
}

// =====================================================================
// SelectionWheel
// =====================================================================

// SelectionWheel builds a selection_wheel control.
type SelectionWheel struct{ *ctrl }

func NewSelectionWheel() *SelectionWheel {
	return &SelectionWheel{&ctrl{c: schema.Control{Type: "selection_wheel"}}}
}

func (sw *SelectionWheel) ID(id string) *SelectionWheel { sw.c.ID = id; return sw }
func (sw *SelectionWheel) Anchor(a schema.Anchor) *SelectionWheel {
	sw.c.AnchorFrom = a
	sw.c.AnchorTo = a
	return sw
}
func (sw *SelectionWheel) Pos(v schema.Relative) *SelectionWheel {
	sw.c.Offset = schema.PtrRel(v)
	return sw
}
func (sw *SelectionWheel) Size(v schema.Relative) *SelectionWheel {
	sw.c.Size = schema.PtrRel(v)
	return sw
}
func (sw *SelectionWheel) Slices(n int) *SelectionWheel {
	if sw.c.Raw == nil {
		sw.c.Raw = make(map[string]json.RawMessage)
	}
	sw.c.Raw["slice_count"] = mustMarshal(n)
	return sw
}

func mustMarshal(v interface{}) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("marshal failed: %v", err))
	}
	return b
}
