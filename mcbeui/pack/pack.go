// Package pack generates a complete MCBE resource pack from screen definitions.
// It produces _ui_defs.json, manifest.json, and all screen JSON files.
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
	Namespace   string
	PackName    string
	Version     [3]int // major, minor, revision
	Screens     []*ScreenDef
	Definitions map[string]string // screenID -> filename
}

// ScreenDef represents a single screen within a pack.
type ScreenDef struct {
	ID        string // unique ID within the pack (becomes the filename without .json)
	Namespace string // JSON UI namespace; defaults to Pack.Namespace
	Control   *schema.Control
}

// New creates an empty Pack with sensible defaults.
func New(namespace, packName string) *Pack {
	return &Pack{
		Namespace:   namespace,
		PackName:    packName,
		Version:     [3]int{1, 0, 0},
		Screens:     make([]*ScreenDef, 0),
		Definitions: make(map[string]string),
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

	// Generate manifest.json
	if err := p.writeManifest(targetDir); err != nil {
		return fmt.Errorf("write manifest.json: %w", err)
	}

	return nil
}

func (p *Pack) writeUIDefs(dir string) error {
	defs := struct {
		Format_version int      `json:"format_version"`
		UIDefs         []string `json:"ui_defs"`
	}{
		Format_version: 1,
	}
	// Sort for deterministic output
	for _, def := range p.Screens {
		defs.UIDefs = append(defs.UIDefs, p.Definitions[def.ID])
	}
	sort.Strings(defs.UIDefs)

	data, err := json.MarshalIndent(defs, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, "_ui_defs.json"), data, 0644)
}

func (p *Pack) writeManifest(dir string) error {
	meta := struct {
		Format_version int `json:"format_version"`
		Header         struct {
			Description        string `json:"description"`
			Name               string `json:"name"`
			UUID               string `json:"uuid"`
			Version            [3]int `json:"version"`
			Min_engine_version [3]int `json:"min_engine_version"`
		} `json:"header"`
		Modules []struct {
			Description string `json:"description"`
			Type        string `json:"type"`
			UUID        string `json:"uuid"`
			Version     [3]int `json:"version"`
		} `json:"modules"`
	}{
		Format_version: 2,
	}
	meta.Header.Description = p.PackName
	meta.Header.Name = p.PackName
	meta.Header.UUID = generateUUID(p.PackName)
	meta.Header.Version = p.Version
	meta.Header.Min_engine_version = [3]int{1, 20, 0}
	meta.Modules = []struct {
		Description string `json:"description"`
		Type        string `json:"type"`
		UUID        string `json:"uuid"`
		Version     [3]int `json:"version"`
	}{{
		Description: "Resource pack for " + p.PackName,
		Type:        "resources",
		UUID:        generateUUID("module:" + p.PackName),
		Version:     p.Version,
	}}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, "manifest.json"), data, 0644)
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
