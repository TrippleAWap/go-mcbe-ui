// Package pack generates a complete MCBE resource pack from screen definitions.
// It produces _ui_defs.json, pack.mcmeta, and all screen JSON files.
package pack

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
)

// Pack is a complete MCBE JSON UI resource pack.
type Pack struct {
	Namespace      string
	PackName       string
	Version        [3]int // major, minor, revision
	Screens        []*ScreenDef
	Definitions    map[string]string // screenID -> filename
	Textures       []string          // texture import paths for pack.mcmeta
	AnimationDefs  []string          // animation definition filenames for _ui_defs.json
	Transitions    []ScreenTransition
}

// ScreenDef represents a single screen within a pack.
type ScreenDef struct {
	ID       string // unique ID within the pack (becomes the filename without .json)
	Namespace string // JSON UI namespace; defaults to Pack.Namespace
	Control  *schema.Control
}

// ScreenTransition represents a transition between two screens.
type ScreenTransition struct {
	From        string
	To          string
	AnimationID string
	Duration    float64 // optional; 0 means use default
}

// AddTexture registers a texture import path for pack.mcmeta.
func (p *Pack) AddTexture(path string) {
	p.Textures = append(p.Textures, path)
}

// AddAnimationDef registers an animation definition file for _ui_defs.json.
func (p *Pack) AddAnimationDef(filename string) {
	p.AnimationDefs = append(p.AnimationDefs, filename)
}

// AddTransition registers a screen transition.
func (p *Pack) AddTransition(t ScreenTransition) {
	p.Transitions = append(p.Transitions, t)
}

// Validate checks the pack for common issues: missing screen references in
// transitions, and transition animation IDs that don't match any registered
// animation definition. Returns human-readable error strings.
func (p *Pack) Validate() []string {
	errs := make([]string, 0)

	// Build set of registered screen IDs.
	screenIDs := make(map[string]bool, len(p.Screens))
	for _, def := range p.Screens {
		screenIDs[def.ID] = true
	}

	// Build set of registered animation def filenames.
	animDefs := make(map[string]bool, len(p.AnimationDefs))
	for _, def := range p.AnimationDefs {
		animDefs[def] = true
	}

	for _, t := range p.Transitions {
		if !screenIDs[t.To] {
			errs = append(errs, fmt.Sprintf("transition %q -> %q: target screen %q not found", t.From, t.To, t.To))
		}
		if t.AnimationID != "" && !animDefs[t.AnimationID] {
			errs = append(errs, fmt.Sprintf("transition %q -> %q: animation def %q not registered", t.From, t.To, t.AnimationID))
		}
	}
	return errs
}

// New creates an empty Pack with sensible defaults.
func New(namespace, packName string) *Pack {
	return &Pack{
		Namespace:     namespace,
		PackName:      packName,
		Version:       [3]int{1, 0, 0},
		Screens:       make([]*ScreenDef, 0),
		Definitions:   make(map[string]string),
		Textures:      make([]string, 0),
		AnimationDefs: make([]string, 0),
		Transitions:   make([]ScreenTransition, 0),
	}
}

// AddScreen registers a screen. The control's Type should be "screen" or "panel".
// Returns the generated filename (e.g., "my_screen.json").
func (p *Pack) AddScreen(screen *schema.Screen) string {
	root := &screen.Root
	id := root.ID
	if id == "" {
		id = "screen_" + fmt.Sprintf("%d", len(p.Screens)+1)
		root.ID = id
	}

	filename := id + ".json"
	p.Screens = append(p.Screens, &ScreenDef{
		ID:        id,
		Namespace: screen.Namespace,
		Control:   root,
	})
	p.Definitions[id] = filename
	return filename
}

// AddScreenDef registers a screen with explicit filename.
func (p *Pack) AddScreenDef(def *ScreenDef) string {
	if def.Namespace == "" {
		def.Namespace = p.Namespace
	}
	p.Screens = append(p.Screens, def)
	p.Definitions[def.ID] = def.ID + ".json"
	return def.ID + ".json"
}

// Build generates all pack files into the target directory.
func (p *Pack) Build(targetDir string) error {
	// Create directory structure
	uiDir := filepath.Join(targetDir, "ui")
	if err := os.MkdirAll(uiDir, 0755); err != nil {
		return fmt.Errorf("create ui dir: %w", err)
	}

	// Generate _ui_defs.json
	if err := p.writeUIDefs(uiDir); err != nil {
		return fmt.Errorf("write _ui_defs.json: %w", err)
	}

	// Generate each screen file
	for _, def := range p.Screens {
		filename := p.Definitions[def.ID]
		data, err := json.MarshalIndent(def.Control, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal screen %s: %w", def.ID, err)
		}
		path := filepath.Join(uiDir, filename)
		if err := os.WriteFile(path, data, 0644); err != nil {
			return fmt.Errorf("write %s: %w", filename, err)
		}
	}

	// Generate pack.mcmeta
	if err := p.writePackMeta(targetDir); err != nil {
		return fmt.Errorf("write pack.mcmeta: %w", err)
	}

	return nil
}

func (p *Pack) writeUIDefs(dir string) error {
	type transitionDef struct {
		From        string  `json:"from"`
		To          string  `json:"to"`
		AnimationID string  `json:"animation_id"`
		Duration    float64 `json:"duration,omitempty"`
	}
	defs := struct {
		Format_version      int            `json:"format_version"`
		UI                  []string       `json:"ui"`
		Animations          []string       `json:"animations,omitempty"`
		Transitions         []transitionDef `json:"transitions,omitempty"`
		Script_api_version  int            `json:"script_api_version,omitempty"`
	}{
		Format_version:      1,
		Script_api_version:  1,
	}
	// Sort for deterministic output
	for _, def := range p.Screens {
		defs.UI = append(defs.UI, p.Definitions[def.ID])
	}
	sort.Strings(defs.UI)
	defs.Animations = append(defs.Animations, p.AnimationDefs...)
	sort.Strings(defs.Animations)
	defs.Transitions = append(defs.Transitions, func() []transitionDef {
		out := make([]transitionDef, len(p.Transitions))
		for i, t := range p.Transitions {
			out[i] = transitionDef{
				From:        t.From,
				To:          t.To,
				AnimationID: t.AnimationID,
				Duration:    t.Duration,
			}
		}
		return out
	}()...)

	data, err := json.MarshalIndent(defs, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, "_ui_defs.json"), data, 0644)
}

func (p *Pack) writePackMeta(dir string) error {
	meta := struct {
		Format_Version [2]int `json:"format_version"`
		Header         struct {
			Description string `json:"description"`
			Format    [2]int `json:"format"`
			Name      string `json:"name"`
			UUID      string `json:"uuid"`
			Version   [3]int `json:"version"`
		} `json:"header"`
		Imports []string `json:"imports,omitempty"`
		Modules []struct {
			Description string `json:"description"`
			Type        string `json:"type"`
			UUID        string `json:"uuid"`
			Version   [3]int `json:"version"`
		} `json:"modules"`
	}{
		Format_Version: [2]int{1, 0},
	}
	meta.Header.Description = p.PackName
	meta.Header.Format = [2]int{1, 0}
	meta.Header.Name = p.PackName
	meta.Header.Version = p.Version
	meta.Imports = append(meta.Imports, p.Textures...)
	meta.Modules = []struct {
		Description string `json:"description"`
		Type        string `json:"type"`
		UUID        string `json:"uuid"`
		Version   [3]int `json:"version"`
	}{{
		Description: "Resource pack for " + p.PackName,
		Type:        "resources",
		UUID:        generateUUID(p.PackName),
		Version:     p.Version,
	}}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, "pack.mcmeta"), data, 0644)
}

func generateUUID(name string) string {
	// Simple deterministic UUID from pack name for reproducibility.
	// In production you'd use a proper UUID generator.
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		hash(name), hash(name+":"), hash(name+"::"),
		hash(name+":::"), hash(name+"::::"))
}

func hash(s string) uint32 {
	var h uint32 = 2166136261
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}
