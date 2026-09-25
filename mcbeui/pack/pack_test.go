package pack

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
)

func TestBasicPack(t *testing.T) {
	p := New("my_addon", "My Test Pack")

	root := &schema.Control{
		ID:   "root",
		Type: "screen",
	}
	p.AddScreenDef(&ScreenDef{
		ID:        "test_screen",
		Namespace: "my_addon",
		Control:   root,
	})

	dir := t.TempDir()
	if err := p.Build(dir); err != nil {
		t.Fatalf("Build: %v", err)
	}

	defsPath := filepath.Join(dir, "ui", "_ui_defs.json")
	defsData, err := os.ReadFile(defsPath)
	if err != nil {
		t.Fatalf("read _ui_defs.json: %v", err)
	}
	var defs map[string]interface{}
	if err := json.Unmarshal(defsData, &defs); err != nil {
		t.Fatalf("parse _ui_defs.json: %v", err)
	}
	uiList, ok := defs["ui_defs"].([]interface{})
	if !ok || len(uiList) != 1 {
		t.Fatalf("expected 1 UI entry, got %v", uiList)
	}
	if uiList[0] != "test_screen.json" {
		t.Errorf("expected test_screen.json, got %v", uiList[0])
	}

	screenPath := filepath.Join(dir, "ui", "test_screen.json")
	screenData, err := os.ReadFile(screenPath)
	if err != nil {
		t.Fatalf("read screen file: %v", err)
	}
	var rootObj map[string]interface{}
	if err := json.Unmarshal(screenData, &rootObj); err != nil {
		t.Fatalf("parse screen JSON: %v", err)
	}
	if rootObj["type"] != "screen" {
		t.Errorf("expected type screen, got %v", rootObj["type"])
	}

	metaPath := filepath.Join(dir, "manifest.json")
	metaData, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("read manifest.json: %v", err)
	}
	var meta map[string]interface{}
	if err := json.Unmarshal(metaData, &meta); err != nil {
		t.Fatalf("parse manifest.json: %v", err)
	}
	if fv, ok := meta["format_version"].(float64); !ok || int(fv) != 2 {
		t.Errorf("expected format_version=2, got %v", meta["format_version"])
	}
	header, ok := meta["header"].(map[string]interface{})
	if !ok {
		t.Fatal("expected header in manifest.json")
	}
	if header["name"] != "My Test Pack" {
		t.Errorf("expected pack name 'My Test Pack', got %v", header["name"])
	}
	if _, ok := header["min_engine_version"]; !ok {
		t.Error("expected min_engine_version in header")
	}
}

func TestUIDefsHasOnlyStandardKeys(t *testing.T) {
	p := New("std_test", "Standard Pack")
	p.AddScreenDef(&ScreenDef{
		ID:        "screen1",
		Namespace: "std_test",
		Control:   &schema.Control{ID: "screen1", Type: "screen"},
	})

	dir := t.TempDir()
	if err := p.Build(dir); err != nil {
		t.Fatalf("Build: %v", err)
	}

	defsData, err := os.ReadFile(filepath.Join(dir, "ui", "_ui_defs.json"))
	if err != nil {
		t.Fatalf("read _ui_defs.json: %v", err)
	}
	var defs map[string]interface{}
	if err := json.Unmarshal(defsData, &defs); err != nil {
		t.Fatalf("parse _ui_defs.json: %v", err)
	}

	for key := range defs {
		switch key {
		case "format_version", "ui_defs":
		default:
			t.Errorf("unexpected key %q in _ui_defs.json", key)
		}
	}
	if _, ok := defs["ui_defs"]; !ok {
		t.Error("expected ui_defs key in _ui_defs.json")
	}
	for _, extra := range []string{"animations", "transitions", "script_api_version"} {
		if _, ok := defs[extra]; ok {
			t.Errorf("expected no non-standard %q key in _ui_defs.json", extra)
		}
	}
}

func TestMultipleScreens(t *testing.T) {
	p := New("shop_addon", "Shop Pack")

	for _, id := range []string{"shop_main", "shop_confirm", "shop_success"} {
		p.AddScreenDef(&ScreenDef{
			ID:        id,
			Namespace: "shop_addon",
			Control:   &schema.Control{ID: id, Type: "screen"},
		})
	}

	dir := t.TempDir()
	if err := p.Build(dir); err != nil {
		t.Fatalf("Build: %v", err)
	}

	defsPath := filepath.Join(dir, "ui", "_ui_defs.json")
	defsData, err := os.ReadFile(defsPath)
	if err != nil {
		t.Fatalf("read _ui_defs.json: %v", err)
	}
	var defs map[string]interface{}
	if err := json.Unmarshal(defsData, &defs); err != nil {
		t.Fatalf("parse _ui_defs.json: %v", err)
	}
	uiList, _ := defs["ui_defs"].([]interface{})
	if len(uiList) != 3 {
		t.Fatalf("expected 3 UI entries, got %d", len(uiList))
	}
	for i := 1; i < len(uiList); i++ {
		a, b := uiList[i-1].(string), uiList[i].(string)
		if a > b {
			t.Errorf("UI list not sorted: %s > %s", a, b)
		}
	}
}

func TestDeterministicOutput(t *testing.T) {
	p1 := New("det_test", "Deterministic Pack")
	p1.AddScreenDef(&ScreenDef{ID: "alpha", Namespace: "det_test", Control: &schema.Control{ID: "alpha", Type: "screen"}})
	p1.AddScreenDef(&ScreenDef{ID: "beta", Namespace: "det_test", Control: &schema.Control{ID: "beta", Type: "screen"}})

	p2 := New("det_test", "Deterministic Pack")
	p2.AddScreenDef(&ScreenDef{ID: "beta", Namespace: "det_test", Control: &schema.Control{ID: "beta", Type: "screen"}})
	p2.AddScreenDef(&ScreenDef{ID: "alpha", Namespace: "det_test", Control: &schema.Control{ID: "alpha", Type: "screen"}})

	d1, _ := os.MkdirTemp("", "pack1")
	d2, _ := os.MkdirTemp("", "pack2")
	p1.Build(d1)
	p2.Build(d2)

	data1, _ := os.ReadFile(filepath.Join(d1, "ui", "_ui_defs.json"))
	data2, _ := os.ReadFile(filepath.Join(d2, "ui", "_ui_defs.json"))
	if string(data1) != string(data2) {
		t.Error("same screens should produce identical _ui_defs.json regardless of registration order")
	}

	m1, _ := os.ReadFile(filepath.Join(d1, "manifest.json"))
	m2, _ := os.ReadFile(filepath.Join(d2, "manifest.json"))
	if string(m1) != string(m2) {
		t.Error("same pack should produce identical manifest.json regardless of registration order")
	}
}

func TestManifestDistinctUUIDs(t *testing.T) {
	p := New("uuid_test", "UUID Pack")
	p.AddScreenDef(&ScreenDef{
		ID:        "screen1",
		Namespace: "uuid_test",
		Control:   &schema.Control{ID: "screen1", Type: "screen"},
	})

	dir := t.TempDir()
	if err := p.Build(dir); err != nil {
		t.Fatalf("Build: %v", err)
	}

	metaData, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		t.Fatalf("read manifest.json: %v", err)
	}
	var meta map[string]interface{}
	if err := json.Unmarshal(metaData, &meta); err != nil {
		t.Fatalf("parse manifest.json: %v", err)
	}

	header := meta["header"].(map[string]interface{})
	modules := meta["modules"].([]interface{})
	if len(modules) != 1 {
		t.Fatalf("expected 1 module, got %d", len(modules))
	}
	mod := modules[0].(map[string]interface{})

	if mod["type"] != "resources" {
		t.Errorf("expected module type resources, got %v", mod["type"])
	}
	headerUUID, _ := header["uuid"].(string)
	moduleUUID, _ := mod["uuid"].(string)
	if headerUUID == "" || moduleUUID == "" {
		t.Fatal("expected UUIDs in header and module")
	}
	if headerUUID == moduleUUID {
		t.Error("header and module UUIDs must be distinct")
	}
}
