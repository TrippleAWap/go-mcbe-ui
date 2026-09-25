package schema

import (
	"encoding/json"
	"fmt"
)

// Relative is a measurement value that can be pixels, percentages, or named constants.
type Relative struct {
	x, y interface{}
}

// Raw creates a Relative from raw dimensions (pixel values, percentages, expressions).
func Raw(dims ...interface{}) Relative {
	switch len(dims) {
	case 1:
		return Relative{x: dims[0]}
	case 2:
		return Relative{x: dims[0], y: dims[1]}
	default:
		panic(fmt.Sprintf("Raw requires 1 or 2 arguments, got %d", len(dims)))
	}
}

// Pixel creates a pixel-relative measurement.
func Pixel(n int) Relative {
	return Relative{x: n}
}

// Percent creates a percentage-relative measurement.
func Percent(n int) Relative {
	return Relative{x: fmt.Sprintf("%d%%", n)}
}

// Named creates a named-relative measurement (e.g., "fill", "parent").
func Named(name string) Relative {
	return Relative{x: name}
}

// SetY sets the Y component on an existing Relative.
func (r Relative) SetY(y Relative) Relative {
	r.y = y.x
	return r
}

// SetX sets the X component on an existing Relative.
func (r Relative) SetX(x Relative) Relative {
	r.x = x.x
	return r
}

// PixelX returns the pixel value of the X axis, or 0 if not a pixel.
func (r Relative) PixelX() int {
	switch v := r.x.(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		return 0
	}
}

// PixelY returns the pixel value of the Y axis, or 0 if not a pixel.
func (r Relative) PixelY() int {
	if r.y == nil {
		return r.PixelX()
	}
	switch v := r.y.(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		return 0
	}
}

// IsEmpty returns true if both axes are unset.
func (r Relative) IsEmpty() bool {
	return r.x == nil && r.y == nil
}

// MarshalJSON serializes a Relative as a JSON array.
func (r Relative) MarshalJSON() ([]byte, error) {
	if r.IsEmpty() {
		return []byte("null"), nil
	}
	switch {
	case r.y == nil:
		return json.Marshal([]interface{}{r.x})
	case r.x == nil:
		return json.Marshal([]interface{}{r.y})
	default:
		return json.Marshal([]interface{}{r.x, r.y})
	}
}

// UnmarshalJSON deserializes a Relative from a JSON array.
func (r *Relative) UnmarshalJSON(data []byte) error {
	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch len(v) {
	case 0:
		r.x, r.y = nil, nil
	case 1:
		r.x = v[0]
		r.y = nil
	case 2:
		r.x = v[0]
		r.y = v[1]
	default:
		return fmt.Errorf("Relative: expected 1 or 2 values, got %d", len(v))
	}
	return nil
}

// PtrRel converts a Relative to a pointer for use in struct literals.
func PtrRel(r Relative) *Relative {
	return &r
}

// ---- Layout Helpers ----

// Pad returns an offset that pushes a control inward by the given amount on each side.
// Useful for creating padding inside a parent panel.
func Pad(h, v int) Relative {
	return Raw(h, v)
}

// Margin returns an offset that pushes a control away from its anchor point.
func Margin(h, v int) Relative {
	return Raw(h, v)
}

// CenterIn returns an offset that centers a child inside its parent.
// The parent must use anchor_from=parent and anchor_to=parent with size.
func CenterIn(parentW, parentH int, childW, childH int) Relative {
	return Raw((parentW-childW)/2, (parentH-childH)/2)
}

// HalfSize returns a size that is half of the given dimensions.
func HalfSize(w, h int) Relative {
	return Raw(w/2, h/2)
}

// Fill returns a relative that fills 100% of the parent on both axes.
func Fill() Relative {
	return Raw("fill", "fill")
}

// FillX returns a relative that fills 100% horizontally.
func FillX() Relative {
	return Raw("fill", "parent")
}

// FillY returns a relative that fills 100% vertically.
func FillY() Relative {
	return Raw("parent", "fill")
}

// BindingType represents the type of a data binding.
type BindingType string

const (
	BindingTypeView       BindingType = "view_binding"
	BindingTypeCollection BindingType = "collection_binding"
	BindingTypeVariable   BindingType = "variable_binding"
)

// BindingCondition represents the condition under which a binding is active.
type BindingCondition string

const (
	BindingConditionAlways BindingCondition = "always"
	BindingConditionWhen   BindingCondition = "when"
)

// Anchor represents a UI anchor point.
type Anchor string

const (
	TopLeft      Anchor = "top_left"
	TopCenter    Anchor = "top_center"
	TopRight     Anchor = "top_right"
	TopMiddle    Anchor = "top_middle"
	CenterLeft   Anchor = "center_left"
	Center       Anchor = "center"
	CenterRight  Anchor = "center_right"
	BottomLeft   Anchor = "bottom_left"
	BottomCenter Anchor = "bottom_center"
	BottomRight  Anchor = "bottom_right"
	BottomMiddle Anchor = "bottom_middle"
	MiddleLeft   Anchor = "middle_left"
	MiddleCenter Anchor = "middle_center"
	MiddleRight  Anchor = "middle_right"
)

// IsValid returns true if the anchor is a recognized value.
func (a Anchor) IsValid() bool {
	return a == TopLeft || a == TopCenter || a == TopRight ||
		a == TopMiddle || a == CenterLeft || a == Center || a == CenterRight ||
		a == BottomLeft || a == BottomCenter || a == BottomRight ||
		a == BottomMiddle || a == MiddleLeft || a == MiddleCenter || a == MiddleRight
}

// Vector2 represents a two-dimensional vector (e.g. size, offset, UV).
type Vector2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Vector3 represents an RGB color vector.
type Vector3 struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
}

// FontSize represents the font size for a label.
type FontSize string

const (
	FontSizeTiny       FontSize = "tiny"
	FontSizeSmall      FontSize = "small"
	FontSizeNormal     FontSize = "normal"
	FontSizeLarge      FontSize = "large"
	FontSizeExtraLarge FontSize = "extra_large"
)

// FontType represents the font family for a label.
type FontType string

const (
	FontTypeDefault FontType = "default"
	FontTypeUniform FontType = "uniform"
)

// TextAlignment represents text alignment within a control.
type TextAlignment string

const (
	TextAlignLeft   TextAlignment = "left"
	TextAlignCenter TextAlignment = "center"
	TextAlignRight  TextAlignment = "right"
)

// Orientation represents the layout direction of a stack panel.
type Orientation string

const (
	OrientationHorizontal Orientation = "horizontal"
	OrientationVertical   Orientation = "vertical"
)

// ClipDirection represents the direction a control clips its children.
type ClipDirection string

const (
	ClipDirectionNone       ClipDirection = "none"
	ClipDirectionHorizontal ClipDirection = "horizontal"
	ClipDirectionVertical   ClipDirection = "vertical"
	ClipDirectionBoth       ClipDirection = "both"
)

// GridFillDirection represents the fill direction for a grid.
type GridFillDirection string

const (
	GridFillDown  GridFillDirection = "down"
	GridFillRight GridFillDirection = "right"
)

// GridRescalingType represents how a grid rescales its items.
type GridRescalingType string

const (
	GridRescalingNone           GridRescalingType = "none"
	GridRescalingMaintainAspect GridRescalingType = "maintain_aspect"
)

// AnchorFrom represents the source anchor.
type AnchorFrom Anchor

// AnchorTo represents the target anchor.
type AnchorTo Anchor

// AnimType represents the type of animation.
type AnimType string

const (
	AnimTypeAlpha     AnimType = "alpha"
	AnimTypeSize      AnimType = "size"
	AnimTypeOffset    AnimType = "offset"
	AnimTypeColor     AnimType = "color"
	AnimTypeUV        AnimType = "uv"
	AnimTypeRotation  AnimType = "rotation"
	AnimTypeAnimation AnimType = "animation"
)

// Easing represents the interpolation easing curve.
type Easing string

const (
	EasingLinear       Easing = "linear"
	EasingInQuad       Easing = "in_quad"
	EasingOutQuad      Easing = "out_quad"
	EasingInOutQuad    Easing = "in_out_quad"
	EasingInCubic      Easing = "in_cubic"
	EasingOutCubic     Easing = "out_cubic"
	EasingInOutCubic   Easing = "in_out_cubic"
	EasingInQuart      Easing = "in_quart"
	EasingOutQuart     Easing = "out_quart"
	EasingInOutQuart   Easing = "in_out_quart"
	EasingInQuint      Easing = "in_quint"
	EasingOutQuint     Easing = "out_quint"
	EasingInOutQuint   Easing = "in_out_quint"
	EasingInSine       Easing = "in_sine"
	EasingOutSine      Easing = "out_sine"
	EasingInOutSine    Easing = "in_out_sine"
	EasingInExpo       Easing = "in_expo"
	EasingOutExpo      Easing = "out_expo"
	EasingInOutExpo    Easing = "in_out_expo"
	EasingInCirc       Easing = "in_circ"
	EasingOutCirc      Easing = "out_circ"
	EasingInOutCirc    Easing = "in_out_circ"
	EasingInBack       Easing = "in_back"
	EasingOutBack      Easing = "out_back"
	EasingInOutBack    Easing = "in_out_back"
	EasingInElastic    Easing = "in_elastic"
	EasingOutElastic   Easing = "out_elastic"
	EasingInOutElastic Easing = "in_out_elastic"
	EasingInBounce     Easing = "in_bounce"
	EasingOutBounce    Easing = "out_bounce"
	EasingInOutBounce  Easing = "in_out_bounce"
)

// Bind represents a single data binding entry.
type Bind struct {
	BindingName             string           `json:"binding_name,omitempty"`
	BindingNameOverride     string           `json:"binding_name_override,omitempty"`
	BindingType             BindingType      `json:"binding_type,omitempty"`
	BindingCollectionName   string           `json:"binding_collection_name,omitempty"`
	BindingCollectionPrefix string           `json:"binding_collection_prefix,omitempty"`
	BindingCondition        BindingCondition `json:"binding_condition,omitempty"`
	SourceControlName       string           `json:"source_control_name,omitempty"`
	SourcePropertyName      string           `json:"source_property_name,omitempty"`
	TargetPropertyName      string           `json:"target_property_name,omitempty"`
	ResolveSiblingScope     bool             `json:"resolve_sibling_scope,omitempty"`
}

// Modification represents a UI modification operation.
type Modification struct {
	ArrayName     string          `json:"array_name"`
	ControlName   string          `json:"control_name"`
	Operation     string          `json:"operation"`
	Value         json.RawMessage `json:"value,omitempty"`
	Where         string          `json:"where,omitempty"`
	Target        string          `json:"target,omitempty"`
	TargetControl string          `json:"target_control,omitempty"`
}

// Variable represents a $variable declaration.
type Variable struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// Control is the base struct for all Bedrock JSON UI controls.
type Control struct {
	ID string `json:"id,omitempty"`

	// Type is the Bedrock control type (panel, label, button, image, etc.)
	Type string `json:"type,omitempty"`

	// --- Core ---
	Controls      []string       `json:"controls,omitempty"`
	Variables     []Variable     `json:"variables,omitempty"`
	Modifications []Modification `json:"modifications,omitempty"`
	Ignored       bool           `json:"ignored,omitempty"`

	// --- Control ---
	Visible              bool            `json:"visible,omitempty"`
	Enabled              bool            `json:"enabled,omitempty"`
	Layer                int             `json:"layer,omitempty"`
	ZOrder               int             `json:"z_order,omitempty"`
	Alpha                float64         `json:"alpha,omitempty"`
	PropagateAlpha       bool            `json:"propagate_alpha,omitempty"`
	ClipsChildren        bool            `json:"clips_children,omitempty"`
	AllowClipping        bool            `json:"allow_clipping,omitempty"`
	ClipOffset           *Vector2        `json:"clip_offset,omitempty"`
	ClipStateChangeEvent string          `json:"clip_state_change_event,omitempty"`
	EnableScissorTest    bool            `json:"enable_scissor_test,omitempty"`
	PropertyBag          json.RawMessage `json:"property_bag,omitempty"`
	Selected             bool            `json:"selected,omitempty"`
	UseChildAnchors      bool            `json:"use_child_anchors,omitempty"`
	CollectionIndex      int             `json:"collection_index,omitempty"`
	Debug                string          `json:"debug,omitempty"`

	// --- Layout ---
	AnchorFrom              Anchor    `json:"anchor_from,omitempty"`
	AnchorTo                Anchor    `json:"anchor_to,omitempty"`
	Size                    *Relative `json:"size,omitempty"`
	MaxSize                 *Relative `json:"max_size,omitempty"`
	MinSize                 *Relative `json:"min_size,omitempty"`
	Offset                  *Relative `json:"offset,omitempty"`
	Contained               bool      `json:"contained,omitempty"`
	Draggable               bool      `json:"draggable,omitempty"`
	FollowsCursor           bool      `json:"follows_cursor,omitempty"`
	InheritMaxSiblingWidth  bool      `json:"inherit_max_sibling_width,omitempty"`
	InheritMaxSiblingHeight bool      `json:"inherit_max_sibling_height,omitempty"`
	UseAnchoredOffset       bool      `json:"use_anchored_offset,omitempty"`

	// --- Data Binding ---
	Bindings              []Bind           `json:"bindings,omitempty"`
	BindingType           BindingType      `json:"binding_type,omitempty"`
	BindingName           string           `json:"binding_name,omitempty"`
	BindingCollectionName string           `json:"binding_collection_name,omitempty"`
	BindingCondition      BindingCondition `json:"binding_condition,omitempty"`

	// --- Sprite / Image ---
	Texture                  string          `json:"texture,omitempty"`
	UV                       *Vector2        `json:"uv,omitempty"`
	UVSize                   *Vector2        `json:"uv_size,omitempty"`
	NineSliceSize            json.RawMessage `json:"nineslice_size,omitempty"`
	BaseSize                 *Vector2        `json:"base_size,omitempty"`
	Color                    json.RawMessage `json:"color,omitempty"`
	Tiled                    bool            `json:"tiled,omitempty"`
	TiledScale               *Vector2        `json:"tiled_scale,omitempty"`
	ClipDirection            ClipDirection   `json:"clip_direction,omitempty"`
	ClipRatio                float64         `json:"clip_ratio,omitempty"`
	ClipPixelPerfect         bool            `json:"clip_pixelperfect,omitempty"`
	KeepRatio                bool            `json:"keep_ratio,omitempty"`
	Fill                     bool            `json:"fill,omitempty"`
	Grayscale                bool            `json:"grayscale,omitempty"`
	Bilinear                 bool            `json:"bilinear,omitempty"`
	PixelPerfect             bool            `json:"pixel_perfect,omitempty"`
	ForceTextureReload       bool            `json:"force_texture_reload,omitempty"`
	ZipFolder                string          `json:"zip_folder,omitempty"`
	TextureFileSystem        string          `json:"texture_file_system,omitempty"`
	AllowDebugMissingTexture bool            `json:"allow_debug_missing_texture,omitempty"`
	FitToWidth               bool            `json:"$fit_to_width,omitempty"`

	// --- Text / Label ---
	Text                  string          `json:"text,omitempty"`
	LockedColor           json.RawMessage `json:"locked_color,omitempty"`
	FontSize              FontSize        `json:"font_size,omitempty"`
	FontScaleFactor       float64         `json:"font_scale_factor,omitempty"`
	FontType              FontType        `json:"font_type,omitempty"`
	BackupFontType        string          `json:"backup_font_type,omitempty"`
	TextAlignment         TextAlignment   `json:"text_alignment,omitempty"`
	Alignment             string          `json:"alignment,omitempty"`
	Localize              bool            `json:"localize,omitempty"`
	LinePadding           float64         `json:"line_padding,omitempty"`
	Shadow                bool            `json:"shadow,omitempty"`
	HideHyphen            bool            `json:"hide_hyphen,omitempty"`
	LockedAlpha           float64         `json:"locked_alpha,omitempty"`
	EnableProfanityFilter bool            `json:"enable_profanity_filter,omitempty"`
	NotifyOnEllipses      []string        `json:"notify_on_ellipses,omitempty"`
	NotifyEllipsesSibling string          `json:"notify_ellipses_sibling,omitempty"`
	UsePlaceHolder        bool            `json:"use_place_holder,omitempty"`
	PlaceHolderText       string          `json:"place_holder_text,omitempty"`
	PlaceHolderTextColor  json.RawMessage `json:"place_holder_text_color,omitempty"`

	// --- Stack Panel ---
	Orientation Orientation `json:"orientation,omitempty"`

	// --- Grid ---
	GridDimensions       *Vector2          `json:"grid_dimensions,omitempty"`
	MaximumGridItems     int               `json:"maximum_grid_items,omitempty"`
	GridDimensionBinding string            `json:"grid_dimension_binding,omitempty"`
	GridItemTemplate     string            `json:"grid_item_template,omitempty"`
	GridFillDirection    GridFillDirection `json:"grid_fill_direction,omitempty"`
	GridRescalingType    GridRescalingType `json:"grid_rescaling_type,omitempty"`
	PreCachedItemCount   int               `json:"precached_grid_item_count,omitempty"`
	GridPosition         *Vector2          `json:"grid_position,omitempty"`

	// --- Button ---
	DefaultControl string `json:"default_control,omitempty"`
	HoverControl   string `json:"hover_control,omitempty"`
	PressedControl string `json:"pressed_control,omitempty"`
	LockedControl  string `json:"locked_control,omitempty"`

	// --- Toggle ---
	ToggleName             string `json:"toggle_name,omitempty"`
	RadioToggleGroup       string `json:"radio_toggle_group,omitempty"`
	ToggleDefaultState     bool   `json:"toggle_default_state,omitempty"`
	ToggleOnButton         string `json:"toggle_on_button,omitempty"`
	ToggleOffButton        string `json:"toggle_off_button,omitempty"`
	CheckedControl         string `json:"checked_control,omitempty"`
	UncheckedControl       string `json:"unchecked_control,omitempty"`
	CheckedHoverControl    string `json:"checked_hover_control,omitempty"`
	UncheckedHoverControl  string `json:"unchecked_hover_control,omitempty"`
	CheckedLockedControl   string `json:"checked_locked_control,omitempty"`
	UncheckedLockedControl string `json:"unchecked_locked_control,omitempty"`
	ResetOnFocusLost       bool   `json:"reset_on_focus_lost,omitempty"`

	// --- Slider ---
	SliderTrackButton         string  `json:"slider_track_button,omitempty"`
	SliderSmallDecreaseButton string  `json:"slider_small_decrease_button,omitempty"`
	SliderSmallIncreaseButton string  `json:"slider_small_increase_button,omitempty"`
	SliderSteps               int     `json:"slider_steps,omitempty"`
	SliderDirection           string  `json:"slider_direction,omitempty"`
	SliderTimeout             float64 `json:"slider_timeout,omitempty"`
	SliderCollectionName      string  `json:"slider_collection_name,omitempty"`
	SliderName                string  `json:"slider_name,omitempty"`
	SliderSelectOnHover       bool    `json:"slider_select_on_hover,omitempty"`
	SliderSelectedButton      string  `json:"slider_selected_button,omitempty"`
	SliderDeselectedButton    string  `json:"slider_deselected_button,omitempty"`
	SliderBoxControl          string  `json:"slider_box_control,omitempty"`
	BackgroundControl         string  `json:"background_control,omitempty"`
	BackgroundHoverControl    string  `json:"background_hover_control,omitempty"`
	ProgressControl           string  `json:"progress_control,omitempty"`
	ProgressHoverControl      string  `json:"progress_hover_control,omitempty"`

	// --- Edit Box ---
	TextBoxName                   string `json:"text_box_name,omitempty"`
	TextEditBoxGridCollectionName string `json:"text_edit_box_grid_collection_name,omitempty"`
	TextType                      string `json:"text_type,omitempty"`
	MaxLength                     int    `json:"max_length,omitempty"`
	EnabledNewline                bool   `json:"enabled_newline,omitempty"`
	ConstrainToRect               bool   `json:"constrain_to_rect,omitempty"`
	TextControl                   string `json:"text_control,omitempty"`
	PlaceHolderControl            string `json:"place_holder_control,omitempty"`
	CanBeDeselected               bool   `json:"can_be_deselected,omitempty"`
	AlwaysListening               bool   `json:"always_listening,omitempty"`
	VirtualKeyboardBufferControl  string `json:"virtual_keyboard_buffer_control,omitempty"`

	// --- Scroll View ---
	ScrollViewPort                 string  `json:"scroll_view_port,omitempty"`
	ScrollContent                  string  `json:"scroll_content,omitempty"`
	ScrollbarBox                   string  `json:"scrollbar_box,omitempty"`
	ScrollbarTrack                 string  `json:"scrollbar_track,omitempty"`
	ScrollSpeed                    float64 `json:"scroll_speed,omitempty"`
	JumpToBottomOnUpdate           bool    `json:"jump_to_bottom_on_update,omitempty"`
	AllowScrollEvenWhenContentFits bool    `json:"allow_scroll_even_when_content_fits,omitempty"`
	ScrollbarTrackButton           string  `json:"scrollbar_track_button,omitempty"`
	ScrollbarTouchButton           string  `json:"scrollbar_touch_button,omitempty"`
	GestureControlEnabled          bool    `json:"gesture_control_enabled,omitempty"`
	AlwaysHandleScrolling          bool    `json:"always_handle_scrolling,omitempty"`
	TouchMode                      bool    `json:"touch_mode,omitempty"`
	ScrollBoxAndTrackPanel         string  `json:"scroll_box_and_track_panel,omitempty"`
	ScrollbarAlwaysVisible         bool    `json:"scrollbar_always_visible,omitempty"`

	// --- Custom Renderer ---
	Renderer           string  `json:"renderer,omitempty"`
	CameraTiltDegrees  float64 `json:"camera_tilt_degrees,omitempty"`
	StartingRotation   float64 `json:"starting_rotation,omitempty"`
	UseSelectedSkin    bool    `json:"use_selected_skin,omitempty"`
	UseUUID            bool    `json:"use_uuid,omitempty"`
	UseSkinGUIScale    bool    `json:"use_skin_gui_scale,omitempty"`
	UsePlayerPaperdoll bool    `json:"use_player_paperdoll,omitempty"`
	Rotation           float64 `json:"rotation,omitempty"`
	AnimationLooped    bool    `json:"animation_looped,omitempty"`

	// --- Screen ---
	IsModal                         bool   `json:"is_modal,omitempty"`
	AbsorbsInput                    bool   `json:"absorbs_input,omitempty"`
	RenderGameBehind                bool   `json:"render_game_behind,omitempty"`
	CloseOnPlayerHurt               bool   `json:"close_on_player_hurt,omitempty"`
	GamepadCursor                   bool   `json:"gamepad_cursor,omitempty"`
	RenderOnlyWhenTopmost           bool   `json:"render_only_when_topmost,omitempty"`
	ScreenNotFlushable              bool   `json:"screen_not_flushable,omitempty"`
	AlwaysAcceptsInput              bool   `json:"always_accepts_input,omitempty"`
	IsShowingMenu                   bool   `json:"is_showing_menu,omitempty"`
	ShouldStealMouse                bool   `json:"should_steal_mouse,omitempty"`
	LowFrequencyRendering           bool   `json:"low_frequency_rendering,omitempty"`
	ScreenDrawsLast                 bool   `json:"screen_draws_last,omitempty"`
	ForceRenderBelow                bool   `json:"force_render_below,omitempty"`
	SendTelemetry                   bool   `json:"send_telemetry,omitempty"`
	CacheScreen                     bool   `json:"cache_screen,omitempty"`
	LoadScreenImmediately           bool   `json:"load_screen_immediately,omitempty"`
	GamepadCursorDeflectionMode     string `json:"gamepad_cursor_deflection_mode,omitempty"`
	ShouldBeSkippedDuringAutomation bool   `json:"should_be_skipped_during_automation,omitempty"`
	VRMode                          bool   `json:"vr_mode,omitempty"`

	// --- Animations ---
	Anims                  []Anim `json:"anims,omitempty"`
	DisableAnimFastForward bool   `json:"disable_anim_fast_forward,omitempty"`
	AnimationResetName     string `json:"animation_reset_name,omitempty"`

	// --- Collection ---
	CollectionName string `json:"collection_name,omitempty"`

	// --- Input ---
	Modal                           bool   `json:"modal,omitempty"`
	InlineModal                     bool   `json:"inline_modal,omitempty"`
	AlwaysListenToInput             bool   `json:"always_listen_to_input,omitempty"`
	AlwaysHandlePointer             bool   `json:"always_handle_pointer,omitempty"`
	AlwaysHandleControllerDirection bool   `json:"always_handle_controller_direction,omitempty"`
	HoverEnabled                    bool   `json:"hover_enabled,omitempty"`
	PreventTouchInput               bool   `json:"prevent_touch_input,omitempty"`
	ConsumeEvent                    bool   `json:"consume_event,omitempty"`
	ConsumeHoverEvents              bool   `json:"consume_hover_events,omitempty"`
	GestureTrackingButton           string `json:"gesture_tracking_button,omitempty"`

	// --- Sound ---
	SoundName   string  `json:"sound_name,omitempty"`
	SoundVolume float64 `json:"sound_volume,omitempty"`
	SoundPitch  float64 `json:"sound_pitch,omitempty"`

	// --- Focus ---
	FocusEnabled       bool   `json:"focus_enabled,omitempty"`
	FocusWrapEnabled   bool   `json:"focus_wrap_enabled,omitempty"`
	FocusMagnetEnabled bool   `json:"focus_magnet_enabled,omitempty"`
	FocusIdentifier    string `json:"focus_identifier,omitempty"`
	FocusChangeDown    string `json:"focus_change_down,omitempty"`
	FocusChangeUp      string `json:"focus_change_up,omitempty"`
	FocusChangeLeft    string `json:"focus_change_left,omitempty"`
	FocusChangeRight   string `json:"focus_change_right,omitempty"`
	FocusContainer     string `json:"focus_container,omitempty"`
	UseLastFocus       bool   `json:"use_last_focus,omitempty"`

	// --- Gradient ---
	GradientDirection string          `json:"gradient_direction,omitempty"`
	GradientColor1    json.RawMessage `json:"color1,omitempty"`
	GradientColor2    json.RawMessage `json:"color2,omitempty"`

	// --- TTS ---
	TTSName          string `json:"tts_name,omitempty"`
	TTSControlHeader string `json:"tts_control_header,omitempty"`

	// Raw allows arbitrary additional fields for future-proofing.
	Raw map[string]json.RawMessage `json:"-"`
}

// MarshalJSON provides custom marshaling so that zero-value optional fields
// are omitted and raw fields are merged into the output.
func (c *Control) MarshalJSON() ([]byte, error) {
	type Alias Control
	aux := &struct {
		*Alias
	}{Alias: (*Alias)(c)}

	b, err := json.Marshal(aux)
	if err != nil {
		return nil, err
	}

	if len(c.Raw) == 0 {
		return b, nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(b, &obj); err != nil {
		return nil, err
	}
	for k, v := range c.Raw {
		obj[k] = v
	}
	return json.Marshal(obj)
}

// AddControl appends a child control ID.
func (c *Control) AddControl(id string) {
	c.Controls = append(c.Controls, id)
}

// Screen represents the top-level Bedrock JSON UI screen definition.
type Screen struct {
	Namespace string  `json:"namespace"`
	RootPanel string  `json:"root_panel"`
	Root      Control `json:"ui"`
}

// NewScreen creates a new Screen with the given namespace.
func NewScreen(namespace string) *Screen {
	return &Screen{
		Namespace: namespace,
		RootPanel: "root",
		Root: Control{
			ID:         "root",
			Type:       "panel",
			AnchorFrom: Center,
			AnchorTo:   Center,
			Size:       PtrRel(Raw(1920, 1080)),
		},
	}
}
