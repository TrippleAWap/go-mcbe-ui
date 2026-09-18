package registry_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/controls"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/controls/components"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/pack"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/registry"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/transitions"
)

// TestIntegrationRoundTrip verifies that a control built with the library can be
// marshaled to JSON, loaded back via FromJSON, and produces an equivalent structure.
func TestIntegrationRoundTrip(t *testing.T) {
	original := controls.NewPanel().
		ID("root").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100)).
		ClipsChildren()

	child := controls.NewLabel().
		ID("child").
		Anchor(schema.TopLeft).
		Size(schema.Raw(50, 30)).
		Text("Hello").
		FontSize(schema.FontSizeNormal).
		Localize()

	original.AddControl(child.Control().ID)

	// Marshal to JSON.
	data, err := json.MarshalIndent(original.Control(), "", "  ")
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	// Unmarshal back.
	restored, err := registry.FromJSON(data)
	if err != nil {
		t.Fatalf("fromJSON failed: %v", err)
	}

	// Verify key fields match.
	if restored.ID != "root" {
		t.Errorf("restored ID = %q, want root", restored.ID)
	}
	if restored.Type != "panel" {
		t.Errorf("restored Type = %q, want panel", restored.Type)
	}
	if restored.ClipsChildren != original.Control().ClipsChildren {
		t.Error("ClipsChildren mismatch after round-trip")
	}
	if len(restored.Controls) != 1 || restored.Controls[0] != "child" {
		t.Errorf("Controls = %v, want [child]", restored.Controls)
	}
}

// TestIntegrationPackOutput verifies the full pipeline: build screens, register,
// validate, inject animations, generate pack, and verify output files.
func TestIntegrationPackOutput(t *testing.T) {
	menu := controls.NewScreen().
		ID("menu").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100)).
		AbsorbsInput()

	game := controls.NewScreen().
		ID("game").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100)).
		AbsorbsInput()

	playBtn := controls.NewButton().
		ID("play_btn").
		Anchor(schema.Center).
		Size(schema.Raw(100, 50)).
		DefaultControl("btn_play")

	rgstry := registry.New()
	rgstry.AddScreen("menu", menu.Control()).
		AddScreen("game", game.Control()).
		AddScreen("play_btn", playBtn.Control()).
		AddAnimationDef("fade").
		AddLink("play_btn", "game", string(transitions.Fade), 0.5)

	out, err := rgstry.ToPack("test_addon")
	if err != nil {
		t.Fatalf("ToPack failed: %v", err)
	}
	if len(out.Screens) != 3 {
		t.Errorf("expected 3 screens, got %d", len(out.Screens))
	}
	if len(out.Transitions) != 1 {
		t.Errorf("expected 1 transition, got %d", len(out.Transitions))
	}

	// Build the pack to disk and verify files.
	dir := t.TempDir()
	p := pack.New("test_addon", "Test Pack")
	p.AddAnimationDef("animations/fade.anim.json")
	for _, def := range out.Screens {
		p.AddScreenDef(&pack.ScreenDef{
			ID:        def.ID,
			Namespace: def.Namespace,
			Control:   def.Control,
		})
	}
	for _, tr := range out.Transitions {
		p.AddTransition(pack.ScreenTransition{
			From:        tr.From,
			To:          tr.To,
			AnimationID: tr.AnimationID,
			Duration:    tr.Duration,
		})
	}

	if err := p.Build(dir); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Verify _ui_defs.json structure.
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
	if !ok || len(uiList) != 3 {
		t.Fatalf("expected 3 UI entries, got %v", defs["ui"])
	}

	transList, ok := defs["transitions"].([]interface{})
	if !ok || len(transList) != 1 {
		t.Fatalf("expected 1 transition, got %v", defs["transitions"])
	}
	t0 := transList[0].(map[string]interface{})
	if t0["from"] != "play_btn" || t0["to"] != "game" || t0["animation_id"] != string(transitions.Fade) {
		t.Errorf("transition mismatch: %v", t0)
	}

	// Verify play_btn screen has animation injected.
	btnPath := filepath.Join(dir, "ui", "play_btn.json")
	btnData, err := os.ReadFile(btnPath)
	if err != nil {
		t.Fatalf("read play_btn.json: %v", err)
	}
	var btnObj map[string]interface{}
	if err := json.Unmarshal(btnData, &btnObj); err != nil {
		t.Fatalf("parse play_btn.json: %v", err)
	}
	anims, ok := btnObj["anims"].([]interface{})
	if !ok || len(anims) != 1 {
		t.Fatalf("expected 1 anim on play_btn, got %v", btnObj["anims"])
	}
	anim0 := anims[0].(map[string]interface{})
	if anim0["anim_type"] != "alpha" {
		t.Errorf("anim type = %v, want alpha", anim0["anim_type"])
	}
	if anim0["duration"] != float64(0.5) {
		t.Errorf("anim duration = %v, want 0.5", anim0["duration"])
	}
}

// TestIntegrationNestedResolution verifies that ResolveIDs correctly reports
// unresolved references and that RegisterNested resolves them.
func TestIntegrationNestedResolution(t *testing.T) {
	parent := controls.NewPanel().
		ID("parent").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100)).
		AddControl("inner_btn")

	innerBtn := controls.NewButton().
		ID("inner_btn").
		Anchor(schema.Center).
		Size(schema.Raw(50, 30))

	rgstry := registry.New()
	rgstry.AddScreen("parent", parent.Control())

	// Before registering the nested control, it should be unresolved.
	unresolved := rgstry.ResolveIDs()
	if len(unresolved) != 1 || unresolved[0] != "inner_btn" {
		t.Errorf("expected [inner_btn], got %v", unresolved)
	}

	errs := rgstry.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for unresolved control")
	}

	// After registering as nested, it should resolve.
	rgstry.RegisterNested("inner_btn", innerBtn.Control())
	unresolved = rgstry.ResolveIDs()
	if len(unresolved) != 0 {
		t.Errorf("expected 0 unresolved after RegisterNested, got %v", unresolved)
	}

	errs = rgstry.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no validation errors, got: %v", errs)
	}
}

// TestIntegrationComponentsPackOutput verifies that component screens produce
// valid pack output when built through the registry-to-pack pipeline.
func TestIntegrationComponentsPackOutput(t *testing.T) {
	shopScreen := components.ShopScreen()
	confirmScreen := components.ConfirmDialog("Exit?", "Are you sure?")
	inventoryScreen := components.InventoryScreenWithData([]components.InventorySlot{
		{ItemID: "minecraft:diamond", Count: 64, SlotID: "diamond_slot"},
		{ItemID: "minecraft:iron_ingot", Count: 32},
	})

	menuScreen := controls.NewScreen().
		ID("menu").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100)).
		AbsorbsInput()

	rgstry := registry.New()
	rgstry.AddScreen("menu", menuScreen.Control()).
		AddScreen("shop", shopScreen.Control()).
		AddScreen("confirm", confirmScreen.Control()).
		AddScreen("inventory", inventoryScreen.Control())

	dir := t.TempDir()
	p := pack.New("test_addon", "Component Pack")
	for id, ctrl := range rgstry.Screens() {
		p.AddScreenDef(&pack.ScreenDef{
			ID:        id,
			Namespace: "test_addon",
			Control:   ctrl,
		})
	}

	if err := p.Build(dir); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Verify all screen files were written and are valid JSON.
	for id, ctrl := range rgstry.Screens() {
		path := filepath.Join(dir, "ui", id+".json")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("missing screen file %s: %v", id, err)
		}
		var obj map[string]interface{}
		if err := json.Unmarshal(data, &obj); err != nil {
			t.Fatalf("invalid JSON in %s: %v", id, err)
		}
		typ, _ := obj["type"].(string)
		if typ != "screen" && typ != "panel" {
			t.Errorf("%s: type = %q, want screen or panel", id, typ)
		}
		// Verify the control has a non-empty ID.
		rootID, _ := obj["id"].(string)
		if rootID == "" {
			t.Errorf("%s: root control has no ID", id)
		}
		_ = ctrl
	}

	// Verify _ui_defs.json lists all screens.
	defsData, err := os.ReadFile(filepath.Join(dir, "ui", "_ui_defs.json"))
	if err != nil {
		t.Fatalf("read _ui_defs.json: %v", err)
	}
	var defs map[string]interface{}
	if err := json.Unmarshal(defsData, &defs); err != nil {
		t.Fatalf("parse _ui_defs.json: %v", err)
	}
	uiList, ok := defs["ui"].([]interface{})
	if !ok || len(uiList) != 4 {
		t.Fatalf("expected 4 UI entries, got %v", defs["ui"])
	}
}

// TestIntegrationComponentRoundTrip verifies that component output round-trips
// through marshal/unmarshal without losing key structure.
func TestIntegrationComponentRoundTrip(t *testing.T) {
	shop := components.ShopScreen()
	data, err := json.Marshal(shop.Control())
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	restored, err := registry.FromJSON(data)
	if err != nil {
		t.Fatalf("fromJSON failed: %v", err)
	}
	if restored.ID != "shop_screen" {
		t.Errorf("restored ID = %q, want shop_screen", restored.ID)
	}
	if restored.Type != "screen" {
		t.Errorf("restored Type = %q, want screen", restored.Type)
	}
}

// TestIntegrationTransitionValidation verifies that ToPack rejects transitions
// with unregistered animation definitions.
func TestIntegrationTransitionValidation(t *testing.T) {
	menu := controls.NewScreen().ID("menu").Control()
	game := controls.NewScreen().ID("game").Control()
	btn := controls.NewScreen().ID("btn").Control()

	rgstry := registry.New()
	rgstry.AddScreen("menu", menu).
		AddScreen("game", game).
		AddScreen("btn", btn).
		AddLink("btn", "game", "unknown_anim", 0.3)

	_, err := rgstry.ToPack("test_addon")
	if err == nil {
		t.Fatal("expected error for unregistered animation def")
	}
	if !strings.Contains(err.Error(), "unknown_anim") {
		t.Errorf("expected error mentioning 'unknown_anim', got: %v", err)
	}
}
