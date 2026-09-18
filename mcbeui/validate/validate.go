// Package validate provides schema-aware validation for MCBE JSON UI controls and screens.
//
// It checks anchors, sizes, bindings, animations, and control references for
// common errors and warnings.
package validate

import (
	"fmt"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
)

var validTypes = map[string]bool{
	"panel": true, "stack_panel": true, "grid": true, "label": true, "image": true,
	"button": true, "toggle": true, "slider": true, "edit_box": true, "dropdown": true,
	"scroll_view": true, "scrollbar_box": true, "factory": true, "screen": true,
	"custom": true, "selection_wheel": true, "tab": true, "carousel_label": true,
	"grid_item": true, "input_panel": true,
}

var validModOps = map[string]bool{
	"insert_back": true, "insert_front": true, "insert_after": true, "insert_before": true,
	"move_back": true, "move_front": true, "move_after": true, "move_before": true,
	"swap": true, "remove": true, "replace": true,
}

func referenceFields(c *schema.Control) []string {
	return []string{
		c.DefaultControl, c.HoverControl, c.PressedControl, c.LockedControl,
		c.ToggleOnButton, c.ToggleOffButton,
		c.CheckedControl, c.UncheckedControl,
		c.CheckedHoverControl, c.UncheckedHoverControl,
		c.CheckedLockedControl, c.UncheckedLockedControl,
		c.SliderTrackButton, c.SliderSelectedButton, c.SliderDeselectedButton,
		c.SliderBoxControl, c.BackgroundControl, c.ProgressControl,
		c.ScrollViewPort, c.ScrollContent, c.ScrollbarBox, c.ScrollbarTrack,
		c.GridItemTemplate,
	}
}

// Warning represents a non-fatal layout concern.
type Warning struct {
	ControlID string
	Message   string
}

// Error represents a validation failure.
type Error struct {
	ControlID string
	Message   string
}

// Result holds validation output.
type Result struct {
	Errors   []Error
	Warnings []Warning
}

// IsValid returns true if there are no errors.
func (r *Result) IsValid() bool {
	return len(r.Errors) == 0
}

// ValidateScreen checks a screen for common issues including anchor validity,
// binding correctness, animation properties, and control type recognition.
func ValidateScreen(s *schema.Screen) *Result {
	r := &Result{}
	if s == nil {
		r.Errors = append(r.Errors, Error{Message: "screen is nil"})
		return r
	}
	r.checkControl(&s.Root, "")
	return r
}

// ValidateControl checks a single control and its properties.
func ValidateControl(c *schema.Control) *Result {
	r := &Result{}
	if c == nil {
		return r
	}
	r.checkControl(c, c.ID)
	return r
}

func (r *Result) checkControl(c *schema.Control, parentID string) {
	if c == nil {
		return
	}

	id := c.ID
	if id == "" {
		r.Errors = append(r.Errors, Error{ControlID: parentID, Message: "control has no ID"})
	}

	r.validateType(c)
	r.validateLayout(c)
	r.validateBindings(c)
	r.validateAnimations(c)
	r.validateReferences(c, parentID)
}

func (r *Result) validateType(c *schema.Control) {
	if c.Type != "" && !validTypes[c.Type] {
		r.Warnings = append(r.Warnings, Warning{
			ControlID: c.ID,
			Message:   fmt.Sprintf("unknown control type %q", c.Type),
		})
	}
}

func (r *Result) validateLayout(c *schema.Control) {
	if (c.Size != nil && c.Size.PixelX() < 0) || (c.Size != nil && c.Size.PixelY() < 0) {
		r.Warnings = append(r.Warnings, Warning{
			ControlID: c.ID,
			Message:   "negative size detected",
		})
	}

	if c.AnchorFrom != "" && !c.AnchorFrom.IsValid() {
		r.Errors = append(r.Errors, Error{
			ControlID: c.ID,
			Message:   fmt.Sprintf("invalid anchor_from %q", c.AnchorFrom),
		})
	}
	if c.AnchorTo != "" && !c.AnchorTo.IsValid() {
		r.Errors = append(r.Errors, Error{
			ControlID: c.ID,
			Message:   fmt.Sprintf("invalid anchor_to %q", c.AnchorTo),
		})
	}

	if c.Layer > 100 {
		r.Warnings = append(r.Warnings, Warning{
			ControlID: c.ID,
			Message:   fmt.Sprintf("layer %d exceeds recommended max of 100", c.Layer),
		})
	}
}

func (r *Result) validateBindings(c *schema.Control) {
	for i, b := range c.Bindings {
		if b.BindingName == "" {
			r.Warnings = append(r.Warnings, Warning{
				ControlID: c.ID,
				Message:   fmt.Sprintf("binding[%d] has no name", i),
			})
		}
		switch b.BindingType {
		case schema.BindingTypeView:
			if b.SourcePropertyName == "" || b.TargetPropertyName == "" {
				r.Errors = append(r.Errors, Error{
					ControlID: c.ID,
					Message:   "view binding requires source_property_name and target_property_name",
				})
			}
		case schema.BindingTypeCollection:
			if b.BindingCollectionName == "" {
				r.Errors = append(r.Errors, Error{
					ControlID: c.ID,
					Message:   "collection binding requires binding_collection_name",
				})
			}
		}
	}
}

func (r *Result) validateAnimations(c *schema.Control) {
	for i, a := range c.Anims {
		if a.AnimType == "" {
			r.Errors = append(r.Errors, Error{
				ControlID: c.ID,
				Message:   fmt.Sprintf("anim[%d] has no anim_type", i),
			})
		}
		if a.Duration <= 0 {
			r.Warnings = append(r.Warnings, Warning{
				ControlID: c.ID,
				Message:   fmt.Sprintf("anim[%d] has non-positive duration", i),
			})
		}
	}
}

func (r *Result) validateReferences(c *schema.Control, parentID string) {
	for i, m := range c.Modifications {
		if !validModOps[m.Operation] {
			r.Warnings = append(r.Warnings, Warning{
				ControlID: c.ID,
				Message:   fmt.Sprintf("modification[%d] has unknown operation %q", i, m.Operation),
			})
		}
	}
	for _, ref := range referenceFields(c) {
		if ref != "" {
			// Reference existence is checked at the screen level in ValidateReferences.
			_ = ref
		}
	}
}

// ValidateReferences checks that all referenced control IDs (e.g., DefaultControl,
// HoverControl, ToggleOnButton) exist within the screen. It validates the root
// control's references against known IDs from the root and its Controls array.
func ValidateReferences(s *schema.Screen) *Result {
	r := &Result{}
	ids := make(map[string]bool)
	collectIDsAll(&s.Root, ids)

	walkReferencesAll(&s.Root, func(ref string) {
		if !ids[ref] {
			r.Errors = append(r.Errors, Error{
				ControlID: "",
				Message:   fmt.Sprintf("referenced control %q not found in screen", ref),
			})
		}
	})
	return r
}

// ValidateScreenWithControls checks a screen using an explicit map of all controls.
// Pass controls built by the user so child control references can be resolved.
func ValidateScreenWithControls(s *schema.Screen, controls map[string]*schema.Control) *Result {
	r := &Result{}
	if s == nil {
		r.Errors = append(r.Errors, Error{Message: "screen is nil"})
		return r
	}
	// Collect all IDs from the explicit map.
	ids := make(map[string]bool)
	for id := range controls {
		r.checkControl(controls[id], "")
		ids[id] = true
	}
	// Also add IDs from the root's Controls array.
	for _, childID := range s.Root.Controls {
		ids[childID] = true
	}
	ids[s.Root.ID] = true

	// Check references on all provided controls.
	for id, c := range controls {
		walkReferencesAll(c, func(ref string) {
			if !ids[ref] {
				r.Errors = append(r.Errors, Error{
					ControlID: id,
					Message:   fmt.Sprintf("referenced control %q not found in screen", ref),
				})
			}
		})
	}

	// Also check the root.
	r.checkControl(&s.Root, "")
	walkReferencesAll(&s.Root, func(ref string) {
		if !ids[ref] {
			r.Errors = append(r.Errors, Error{
				ControlID: "",
				Message:   fmt.Sprintf("referenced control %q not found in screen", ref),
			})
		}
	})
	return r
}

func collectIDsAll(c *schema.Control, ids map[string]bool) {
	if c == nil {
		return
	}
	if c.ID != "" {
		ids[c.ID] = true
	}
	// Add all first-level child IDs from the Controls array.
	for _, childID := range c.Controls {
		ids[childID] = true
	}
	// Add IDs from Modifications targets.
	for _, mod := range c.Modifications {
		if mod.TargetControl != "" {
			ids[mod.TargetControl] = true
		}
		if mod.ControlName != "" {
			ids[mod.ControlName] = true
		}
	}
}

func walkReferencesAll(c *schema.Control, fn func(string)) {
	if c == nil {
		return
	}
	for _, ref := range referenceFields(c) {
		if ref != "" {
			fn(ref)
		}
	}
}

// String returns a human-readable summary of the result.
func (r *Result) String() string {
	if r.IsValid() && len(r.Warnings) == 0 {
		return "validation passed"
	}
	var out string
	for _, e := range r.Errors {
		out += fmt.Sprintf("ERROR [%s]: %s\n", e.ControlID, e.Message)
	}
	for _, w := range r.Warnings {
		out += fmt.Sprintf("WARN [%s]: %s\n", w.ControlID, w.Message)
	}
	return out
}
