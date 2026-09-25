// Package registry provides a screen registry that manages screens, validates
// control references, and loads existing JSON UI files.
package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
)

// Registry manages a collection of screens and their nested controls.
type Registry struct {
	screens map[string]*schema.Control
	nested  map[string]*schema.Control // manually-registered nested controls
	errs    []string
}

// New creates an empty registry.
func New() *Registry {
	return &Registry{
		screens: make(map[string]*schema.Control),
		nested:  make(map[string]*schema.Control),
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

// Validate checks that all Controls[] references resolve to a registered
// screen or nested control. Returns human-readable error strings.
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

	return r.errs
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

// ToPack converts the registry into a pack output structure ready for pack
// generation. It validates the registry (screens and control references) and
// returns an error if any validation fails.
func (r *Registry) ToPack(namespace string) (*PackOutput, error) {
	errs := r.Validate()
	if len(errs) > 0 {
		return nil, fmt.Errorf("registry validation failed: %v", errs)
	}

	out := &PackOutput{
		Namespace: namespace,
		Screens:   make([]ScreenDef, 0, len(r.screens)),
	}

	for id, ctrl := range r.screens {
		out.Screens = append(out.Screens, ScreenDef{
			ID:        id,
			Namespace: namespace,
			Control:   ctrl,
		})
	}

	return out, nil
}

// ScreenDef is a screen ready for pack generation.
type ScreenDef struct {
	ID        string
	Namespace string
	Control   *schema.Control
}

// PackOutput is the complete output of a registry build.
type PackOutput struct {
	Namespace string
	Screens   []ScreenDef
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
		UIDefs []string `json:"ui_defs"`
	}
	if err := json.Unmarshal(data, &defs); err != nil {
		return nil, fmt.Errorf("unmarshal _ui_defs.json: %w", err)
	}
	return defs.UIDefs, nil
}
