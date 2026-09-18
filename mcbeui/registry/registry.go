// Package registry provides a screen registry that manages screens, validates
// navigation links, and wires transition animations into controls.
package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/pack"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
)

// Link describes a navigation transition from one control to another screen.
type Link struct {
	FromID   string
	ToScreen string
	AnimID   string
	Duration float64
	Easing   schema.Easing
}

// Registry manages a collection of screens and their navigation links.
type Registry struct {
	screens        map[string]*schema.Control
	nested         map[string]*schema.Control // manually-registered nested controls
	links          []Link
	animationDefs  []string
	errs           []string
}

// New creates an empty registry.
func New() *Registry {
	return &Registry{
		screens:       make(map[string]*schema.Control),
		nested:        make(map[string]*schema.Control),
		links:         make([]Link, 0),
		animationDefs: make([]string, 0),
	}
}

// AddScreen registers a screen with the given ID. Errors are collected
// internally and retrievable via Errors(). Returns r for chaining.
func (r *Registry) AddScreen(id string, ctrl *schema.Control) *Registry {
	if ctrl == nil {
		r.errs = append(r.errs, fmt.Sprintf("screen %q: nil control", id))
		return r
	}
	if _, exists := r.screens[id]; exists {
		r.errs = append(r.errs, fmt.Sprintf("screen %q: duplicate ID", id))
		return r
	}
	r.screens[id] = ctrl
	return r
}

// AddAnimationDef registers an animation definition file for transition validation.
// This is checked by ToPack() to ensure every transition references a known anim def.
func (r *Registry) AddAnimationDef(filename string) *Registry {
	r.animationDefs = append(r.animationDefs, filename)
	return r
}

// AnimationDefs returns all registered animation definition filenames.
func (r *Registry) AnimationDefs() []string { return r.animationDefs }
func (r *Registry) AddLink(fromID, toScreen, animID string, duration float64) *Registry {
	r.links = append(r.links, Link{
		FromID:   fromID,
		ToScreen: toScreen,
		AnimID:   animID,
		Duration: duration,
	})
	return r
}

// Links returns all registered navigation links.
func (r *Registry) Links() []Link { return r.links }

// Screens returns a copy of the registered screens.
func (r *Registry) Screens() map[string]*schema.Control {
	out := make(map[string]*schema.Control, len(r.screens))
	for k, v := range r.screens {
		out[k] = v
	}
	return out
}

// Nested returns a copy of the manually-registered nested controls.
func (r *Registry) Nested() map[string]*schema.Control {
	out := make(map[string]*schema.Control, len(r.nested))
	for k, v := range r.nested {
		out[k] = v
	}
	return out
}

// Validate checks that all links reference valid screens and controls that
// exist within those screens. Returns human-readable error strings.
func (r *Registry) Validate() []string {
	r.errs = nil

	// Build full ID->pointer index across all screens and nested controls.
	idx := CollectAllControls(r.screens, r.nested)

	// Check for unresolved Controls[] references.
	for _, ctrl := range r.screens {
		for _, childID := range ctrl.Controls {
			if childID == "" || idx[childID] != nil {
				continue
			}
			r.errs = append(r.errs, fmt.Sprintf("screen %q: unresolved control %q in Controls[]", ctrl.ID, childID))
		}
	}

	for _, link := range r.links {
		// Check target screen exists.
		if _, ok := r.screens[link.ToScreen]; !ok {
			r.errs = append(r.errs, fmt.Sprintf("link %q -> %q: target screen not found", link.FromID, link.ToScreen))
			continue
		}
		// Check source control exists.
		if link.FromID == "" {
			r.errs = append(r.errs, fmt.Sprintf("link to %q: empty source control ID", link.ToScreen))
			continue
		}
		if _, ok := idx[link.FromID]; !ok {
			r.errs = append(r.errs, fmt.Sprintf("link %q -> %q: source control %q not found", link.FromID, link.ToScreen, link.FromID))
		}
	}
	return r.errs
}

// InjectAnimations attaches transition animations to the controls that trigger them.
// This modifies the control pointers in-place. Call Validate() first.
func (r *Registry) InjectAnimations() {
	index := CollectAllControls(r.screens, r.nested)
	for _, link := range r.links {
		if ctrl, ok := index[link.FromID]; ok {
			ctrl.Anims = append(ctrl.Anims, buildTransitionAnim(link))
		}
	}
}

// Errors returns any validation errors from the last Validate() call.
func (r *Registry) Errors() []string { return r.errs }

// RegisterNested manually adds a control to the registry so it can be resolved
// by ID even if it is not registered as a top-level screen. This is useful for
// deeply nested controls that cannot be reached through Controls[] strings alone.
func (r *Registry) RegisterNested(id string, ctrl *schema.Control) *Registry {
	if ctrl != nil && id != "" {
		r.nested[id] = ctrl
	}
	return r
}

// ResolveIDs scans all registered screens' Controls[] arrays and returns the
// IDs that are referenced but cannot be resolved to a known control pointer.
// These are controls that exist only as string references without a matching
// screen or nested registration. Call this after adding all screens to discover
// what needs RegisterNested() or AddScreen().
func (r *Registry) ResolveIDs() []string {
	idx := CollectAllControls(r.screens, r.nested)
	resolved := make(map[string]bool, len(idx))
	for id := range idx {
		resolved[id] = true
	}
	var unresolved []string
	seen := make(map[string]bool)
	for _, ctrl := range r.screens {
		for _, childID := range ctrl.Controls {
			if childID == "" || resolved[childID] || seen[childID] {
				continue
			}
			seen[childID] = true
			unresolved = append(unresolved, childID)
		}
	}
	return unresolved
}

// LoadScreen reads a single screen JSON file and registers it under the given ID.
// The file should contain a single control object (the screen root).
func (r *Registry) LoadScreen(id, path string) *Registry {
	data, err := os.ReadFile(path)
	if err != nil {
		r.errs = append(r.errs, fmt.Sprintf("load %q: %v", id, err))
		return r
	}
	ctrl, err := FromJSON(data)
	if err != nil {
		r.errs = append(r.errs, fmt.Sprintf("load %q: %v", id, err))
		return r
	}
	r.AddScreen(id, ctrl)
	return r
}

// LoadFromDir reads all .json files from the given directory, loads each as a
// screen named after the filename (without extension), and registers them.
// _ui_defs.json is skipped. Returns the registry for chaining.
func (r *Registry) LoadFromDir(dir string) *Registry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		r.errs = append(r.errs, fmt.Sprintf("read dir %q: %v", dir, err))
		return r
	}
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "_ui_defs.json" {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		id := entry.Name()[:len(entry.Name())-5] // strip .json
		r.LoadScreen(id, filepath.Join(dir, entry.Name()))
	}
	return r
}

func (r *Registry) lastErr() error {
	if len(r.errs) == 0 {
		return nil
	}
	return fmt.Errorf("%s", r.errs[0])
}

// CollectAllControls returns a flat map of every control ID to its pointer
// across all screens and any manually-registered nested controls.
//
// Because MCBE JSON UI stores child references as string IDs (Controls []string),
// only controls registered as top-level screens or added via RegisterNested
// can be resolved. Controls nested inside Controls[] without their own screen
// registration cannot be found. Use ResolveIDs() to discover unresolved references.
func CollectAllControls(screens, nested map[string]*schema.Control) map[string]*schema.Control {
	idx := make(map[string]*schema.Control)
	// Seed with all screen roots.
	for id, sc := range screens {
		if sc != nil && id != "" {
			idx[id] = sc
		}
	}
	// Seed with manually-registered nested controls.
	for id, nc := range nested {
		if nc != nil && id != "" {
			idx[id] = nc
		}
	}
	// Iteratively resolve Controls[] references until stable.
	for {
		added := false
		for _, ctrl := range idx {
			for _, childID := range ctrl.Controls {
				if childID == "" || idx[childID] != nil {
					continue
				}
				// Only resolve if the child is registered as a screen or nested control.
				// Controls[] strings alone cannot be followed because there is no
				// pointer chain back to the child definition.
				if _, ok := screens[childID]; ok && idx[childID] == nil {
					idx[childID] = screens[childID]
					added = true
				}
				if _, ok := nested[childID]; ok && idx[childID] == nil {
					idx[childID] = nested[childID]
					added = true
				}
			}
		}
		if !added {
			break
		}
	}
	return idx
}

// buildTransitionAnim creates a schema.Anim from a Link's transition params.
func buildTransitionAnim(l Link) schema.Anim {
	if l.Duration <= 0 {
		l.Duration = 0.3
	}
	switch l.AnimID {
	case "fade", "":
		return schema.AnimAlpha(l.Duration, 0, 1)
	case "slide-left":
		return schema.AnimOffset(l.Duration,
			schema.Vector2{X: -200, Y: 0},
			schema.Vector2{X: 0, Y: 0})
	case "slide-right":
		return schema.AnimOffset(l.Duration,
			schema.Vector2{X: 200, Y: 0},
			schema.Vector2{X: 0, Y: 0})
	case "scale":
		return schema.AnimSize(l.Duration,
			schema.Vector2{X: 0, Y: 0},
			schema.Vector2{X: 1, Y: 1})
	default:
		return schema.AnimAlpha(l.Duration, 0, 1)
	}
}

// ToPack converts the registry into a pack output structure ready for
// serialization. It validates the registry (screens, controls, links), then
// validates transitions against registered animation definitions, and injects
// animations into triggering controls. Returns an error if any validation fails.
func (r *Registry) ToPack(namespace string) (*PackOutput, error) {
	errs := r.Validate()
	if len(errs) > 0 {
		return nil, fmt.Errorf("registry validation failed: %v", errs)
	}
	r.InjectAnimations()

	out := &PackOutput{
		Namespace:   namespace,
		Screens:     make([]ScreenDef, 0, len(r.screens)),
		Transitions: make([]TransDef, 0, len(r.links)),
	}

	for id, ctrl := range r.screens {
		out.Screens = append(out.Screens, ScreenDef{
			ID:        id,
			Namespace: namespace,
			Control:   ctrl,
		})
	}

	for _, link := range r.links {
		out.Transitions = append(out.Transitions, TransDef{
			From:        link.FromID,
			To:          link.ToScreen,
			AnimationID: link.AnimID,
			Duration:    link.Duration,
		})
	}

	// Validate transitions against registered animation defs.
	p := pack.New(namespace, namespace)
	for _, def := range out.Screens {
		p.AddScreenDef(&pack.ScreenDef{ID: def.ID, Namespace: def.Namespace, Control: def.Control})
	}
	for _, tr := range out.Transitions {
		p.AddTransition(pack.ScreenTransition{From: tr.From, To: tr.To, AnimationID: tr.AnimationID, Duration: tr.Duration})
	}
	for _, def := range r.animationDefs {
		p.AddAnimationDef(def)
	}
	if transErrs := p.Validate(); len(transErrs) > 0 {
		return nil, fmt.Errorf("transition validation failed: %v", transErrs)
	}

	return out, nil
}

// ScreenDef is a screen ready for pack generation.
type ScreenDef struct {
	ID        string
	Namespace string
	Control   *schema.Control
}

// TransDef is a screen-to-screen transition ready for _ui_defs.json.
type TransDef struct {
	From        string
	To          string
	AnimationID string
	Duration    float64
}

// PackOutput is the complete output of a registry build.
type PackOutput struct {
	Namespace   string
	Screens     []ScreenDef
	Transitions []TransDef
}

// FromJSON loads a control from raw JSON bytes. The JSON should contain a single
// control object (typically the root of a screen). Returns the control or an error.
func FromJSON(data []byte) (*schema.Control, error) {
	var c schema.Control
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("unmarshal control: %w", err)
	}
	return &c, nil
}

// FromUIDefs parses a _ui_defs.json file and returns the list of UI file
// references. Use LoadScreen to load each individual screen.
func FromUIDefs(data []byte) ([]string, error) {
	var defs struct {
		UI []string `json:"ui"`
	}
	if err := json.Unmarshal(data, &defs); err != nil {
		return nil, fmt.Errorf("unmarshal _ui_defs.json: %w", err)
	}
	return defs.UI, nil
}
