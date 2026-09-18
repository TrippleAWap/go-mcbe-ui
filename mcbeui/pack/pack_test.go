package pack

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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
	uiList, ok := defs["ui"].([]interface{})
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

	metaPath := filepath.Join(dir, "pack.mcmeta")
	metaData, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("read pack.mcmeta: %v", err)
	}
	var meta map[string]interface{}
	if err := json.Unmarshal(metaData, &meta); err != nil {
		t.Fatalf("parse pack.mcmeta: %v", err)
	}
	if fv, ok := meta["format_version"].([]interface{}); !ok || len(fv) < 1 || int(fv[0].(float64)) != 1 {
		t.Errorf("expected format_version[0]=1, got %v", meta["format_version"])
	}
	header, ok := meta["header"].(map[string]interface{})
	if !ok {
		t.Fatal("expected header in pack.mcmeta")
	}
	if header["name"] != "My Test Pack" {
		t.Errorf("expected pack name 'My Test Pack', got %v", header["name"])
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
	uiList, _ := defs["ui"].([]interface{})
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
}

func TestAddTexture(t *testing.T) {
	p := New("texture_test", "Texture Pack")
	p.AddTexture("textures/items/test_item")
	p.AddTexture("textures/blocks/test_block")

	p.AddScreenDef(&ScreenDef{
		ID:        "screen1",
		Namespace: "texture_test",
		Control:   &schema.Control{ID: "screen1", Type: "screen"},
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
	if _, ok := defs["animations"]; ok {
		t.Error("expected no animations key in _ui_defs.json")
	}
	if _, ok := defs["transitions"]; ok {
		t.Error("expected no transitions key in _ui_defs.json")
	}

	metaPath := filepath.Join(dir, "pack.mcmeta")
	metaData, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("read pack.mcmeta: %v", err)
	}
	var meta map[string]interface{}
	if err := json.Unmarshal(metaData, &meta); err != nil {
		t.Fatalf("parse pack.mcmeta: %v", err)
	}
	imports, ok := meta["imports"].([]interface{})
	if !ok {
		t.Fatal("expected imports array in pack.mcmeta")
	}
	if len(imports) != 2 {
		t.Fatalf("expected 2 imports, got %d", len(imports))
	}
	expectedImports := []string{"textures/items/test_item", "textures/blocks/test_block"}
	for i, exp := range expectedImports {
		if imports[i] != exp {
			t.Errorf("expected import[%d]=%q, got %v", i, exp, imports[i])
		}
	}
}

func TestAddAnimationDef(t *testing.T) {
	p := New("anim_test", "Animation Pack")
	p.AddAnimationDef("animations/test.anim.json")
	p.AddAnimationDef("animations/menu_open.anim.json")

	p.AddScreenDef(&ScreenDef{
		ID:        "screen1",
		Namespace: "anim_test",
		Control:   &schema.Control{ID: "screen1", Type: "screen"},
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
	animList, ok := defs["animations"].([]interface{})
	if !ok {
		t.Fatal("expected animations array in _ui_defs.json")
	}
	if len(animList) != 2 {
		t.Fatalf("expected 2 animations, got %d", len(animList))
	}
	expectedAnims := []string{"animations/menu_open.anim.json", "animations/test.anim.json"}
	for i, exp := range expectedAnims {
		if animList[i] != exp {
			t.Errorf("expected animations[%d]=%q, got %v", i, exp, animList[i])
		}
	}
	if _, ok := defs["transitions"]; ok {
		t.Error("expected no transitions key in _ui_defs.json")
	}
}

func TestAddTransition(t *testing.T) {
	p := New("trans_test", "Transition Pack")
	p.AddTransition(ScreenTransition{From: "menu", To: "game", AnimationID: "transition_slide"})
	p.AddTransition(ScreenTransition{From: "game", To: "menu", AnimationID: "transition_slide_back"})

	p.AddScreenDef(&ScreenDef{
		ID:        "menu",
		Namespace: "trans_test",
		Control:   &schema.Control{ID: "menu", Type: "screen"},
	})
	p.AddScreenDef(&ScreenDef{
		ID:        "game",
		Namespace: "trans_test",
		Control:   &schema.Control{ID: "game", Type: "screen"},
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
	transList, ok := defs["transitions"].([]interface{})
	if !ok {
		t.Fatal("expected transitions array in _ui_defs.json")
	}
	if len(transList) != 2 {
		t.Fatalf("expected 2 transitions, got %d", len(transList))
	}
	m0 := transList[0].(map[string]interface{})
	m1 := transList[1].(map[string]interface{})
	// Transitions are in registration order
	var first, second map[string]interface{}
	if m0["from"] == "menu" {
		first, second = m0, m1
	} else {
		first, second = m1, m0
	}
	if first["from"] != "menu" || first["to"] != "game" || first["animation_id"] != "transition_slide" {
		t.Errorf("first transition mismatch: %v", first)
	}
	if second["from"] != "game" || second["to"] != "menu" || second["animation_id"] != "transition_slide_back" {
		t.Errorf("second transition mismatch: %v", second)
	}
}

func TestCombinedFeatures(t *testing.T) {
	p := New("combo_test", "Combined Pack")
	p.AddTexture("textures/ui/bg_texture")
	p.AddTexture("textures/ui/button_texture")
	p.AddAnimationDef("animations/button_press.anim.json")
	p.AddAnimationDef("animations/screen_fade.anim.json")
	p.AddTransition(ScreenTransition{From: "main_menu", To: "game", AnimationID: "animations/screen_fade.anim.json"})
	p.AddTransition(ScreenTransition{From: "game", To: "main_menu", AnimationID: "animations/screen_fade.anim.json"})

	p.AddScreenDef(&ScreenDef{
		ID:        "main_menu",
		Namespace: "combo_test",
		Control:   &schema.Control{ID: "main_menu", Type: "screen"},
	})
	p.AddScreenDef(&ScreenDef{
		ID:        "game",
		Namespace: "combo_test",
		Control:   &schema.Control{ID: "game", Type: "screen"},
	})

	dir := t.TempDir()
	if err := p.Build(dir); err != nil {
		t.Fatalf("Build: %v", err)
	}

	// Check _ui_defs.json
	defsData, err := os.ReadFile(filepath.Join(dir, "ui", "_ui_defs.json"))
	if err != nil {
		t.Fatalf("read _ui_defs.json: %v", err)
	}
	var defs map[string]interface{}
	if err := json.Unmarshal(defsData, &defs); err != nil {
		t.Fatalf("parse _ui_defs.json: %v", err)
	}

	// Should have ui, animations, transitions
	if ui, ok := defs["ui"].([]interface{}); !ok || len(ui) != 2 {
		t.Errorf("expected 2 UI entries, got %v", defs["ui"])
	}
	if anims, ok := defs["animations"].([]interface{}); !ok || len(anims) != 2 {
		t.Errorf("expected 2 animations, got %v", defs["animations"])
	}
	if trans, ok := defs["transitions"].([]interface{}); !ok || len(trans) != 2 {
		t.Errorf("expected 2 transitions, got %v", defs["transitions"])
	}

	// Check pack.mcmeta
	metaData, err := os.ReadFile(filepath.Join(dir, "pack.mcmeta"))
	if err != nil {
		t.Fatalf("read pack.mcmeta: %v", err)
	}
	var meta map[string]interface{}
	if err := json.Unmarshal(metaData, &meta); err != nil {
		t.Fatalf("parse pack.mcmeta: %v", err)
	}
	if imports, ok := meta["imports"].([]interface{}); !ok || len(imports) != 2 {
		t.Errorf("expected 2 imports in pack.mcmeta, got %v", meta["imports"])
	}
	if _, ok := meta["modules"]; !ok {
		t.Error("expected modules in pack.mcmeta")
	}
}

func TestValidateMissingTargetScreen(t *testing.T) {
	p := New("val_test", "Validation Pack")
	p.AddTransition(ScreenTransition{From: "menu", To: "missing_screen", AnimationID: "fade"})

	_ = p.AddScreenDef(&ScreenDef{
		ID:        "menu",
		Namespace: "val_test",
		Control:   &schema.Control{ID: "menu", Type: "screen"},
	})

	errs := p.Validate()
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0], "missing_screen") {
		t.Errorf("expected error mentioning 'missing_screen', got: %s", errs[0])
	}
}

func TestValidateMissingAnimDef(t *testing.T) {
	p := New("val_test", "Validation Pack")
	p.AddAnimationDef("animations/fade.anim.json")
	p.AddTransition(ScreenTransition{From: "menu", To: "game", AnimationID: "animations/missing.anim.json"})
	p.AddTransition(ScreenTransition{From: "game", To: "menu", AnimationID: "animations/fade.anim.json"})

	_ = p.AddScreenDef(&ScreenDef{ID: "menu", Namespace: "val_test", Control: &schema.Control{ID: "menu", Type: "screen"}})
	_ = p.AddScreenDef(&ScreenDef{ID: "game", Namespace: "val_test", Control: &schema.Control{ID: "game", Type: "screen"}})

	errs := p.Validate()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0], "missing.anim.json") {
		t.Errorf("expected error mentioning 'missing.anim.json', got: %s", errs[0])
	}
}

func TestValidateAllGood(t *testing.T) {
	p := New("val_test", "Validation Pack")
	p.AddAnimationDef("animations/fade.anim.json")
	p.AddTransition(ScreenTransition{From: "menu", To: "game", AnimationID: "animations/fade.anim.json"})

	_ = p.AddScreenDef(&ScreenDef{ID: "menu", Namespace: "val_test", Control: &schema.Control{ID: "menu", Type: "screen"}})
	_ = p.AddScreenDef(&ScreenDef{ID: "game", Namespace: "val_test", Control: &schema.Control{ID: "game", Type: "screen"}})

	errs := p.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}
