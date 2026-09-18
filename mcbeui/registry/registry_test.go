package registry_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/controls"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/registry"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
)

func TestNewRegistry(t *testing.T) {
	r := registry.New()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestAddScreen(t *testing.T) {
	r := registry.New()
	screen := controls.NewScreen().ID("menu").Control()
	r.AddScreen("menu", screen)
	if len(r.Screens()) != 1 {
		t.Errorf("expected 1 screen, got %d", len(r.Screens()))
	}
}

func TestAddScreenFluent(t *testing.T) {
	r := registry.New()
	out := r.AddScreen("a", controls.NewScreen().ID("a").Control()).
		AddScreen("b", controls.NewScreen().ID("b").Control())
	if out != r {
		t.Error("AddScreen should return the registry for chaining")
	}
	if len(r.Screens()) != 2 {
		t.Errorf("expected 2 screens, got %d", len(r.Screens()))
	}
}

func TestAddDuplicateScreen(t *testing.T) {
	r := registry.New()
	screen := controls.NewScreen().ID("menu").Control()
	r.AddScreen("menu", screen).AddScreen("menu", screen)
	errs := r.Errors()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestAddNilScreen(t *testing.T) {
	r := registry.New()
	r.AddScreen("menu", nil)
	errs := r.Errors()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestAddLink(t *testing.T) {
	r := registry.New()
	r.AddScreen("menu", controls.NewScreen().ID("menu").Control()).
		AddScreen("game", controls.NewScreen().ID("game").Control()).
		AddLink("play_btn", "game", "fade", 0.5)

	links := r.Links()
	if len(links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(links))
	}
	if links[0].FromID != "play_btn" {
		t.Errorf("FromID = %q, want play_btn", links[0].FromID)
	}
	if links[0].ToScreen != "game" {
		t.Errorf("ToScreen = %q, want game", links[0].ToScreen)
	}
}

func TestValidateMissingTarget(t *testing.T) {
	r := registry.New()
	r.AddScreen("menu", controls.NewScreen().ID("menu").Control())
	r.AddLink("btn", "missing_screen", "", 0)

	errs := r.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for missing target screen")
	}
}

func TestValidateEmptyFromID(t *testing.T) {
	r := registry.New()
	r.AddScreen("menu", controls.NewScreen().ID("menu").Control())
	r.AddScreen("game", controls.NewScreen().ID("game").Control())
	r.AddLink("", "game", "", 0)

	errs := r.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for empty FromID")
	}
}

func TestValidateMissingSourceControl(t *testing.T) {
	r := registry.New()
	r.AddScreen("menu", controls.NewScreen().ID("menu").Control())
	r.AddScreen("game", controls.NewScreen().ID("game").Control())
	r.AddLink("nonexistent_btn", "game", "fade", 0)

	errs := r.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for missing source control")
	}
}

func TestValidateAllGood(t *testing.T) {
	r := registry.New()
	r.AddScreen("menu", controls.NewScreen().ID("menu").Control())
	r.AddScreen("game", controls.NewScreen().ID("game").Control())
	r.AddScreen("play_btn", controls.NewScreen().ID("play_btn").Control())
	r.AddLink("play_btn", "game", "fade", 0.5)

	errs := r.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}

func TestInjectAnimations(t *testing.T) {
	playScreen := controls.NewScreen().ID("play_btn").Control()
	menuRoot := controls.NewScreen().ID("menu").Control()

	r := registry.New()
	r.AddScreen("menu", menuRoot).AddScreen("play_btn", playScreen)
	r.AddLink("play_btn", "game", "fade", 0.5)

	// Validate first (game screen doesn't exist, so Validate fails)
	// We test injection directly.
	index := registry.CollectAllControls(r.Screens(), r.Nested())
	btn, ok := index["play_btn"]
	if !ok {
		t.Fatal("play_btn not found in index")
	}
	if len(btn.Anims) != 0 {
		t.Fatalf("expected 0 anims before injection, got %d", len(btn.Anims))
	}

	r.InjectAnimations()

	btn2, ok := index["play_btn"]
	if !ok {
		t.Fatal("play_btn not found after injection")
	}
	if len(btn2.Anims) != 1 {
		t.Fatalf("expected 1 anim after injection, got %d", len(btn2.Anims))
	}
	if btn2.Anims[0].AnimType != schema.AnimTypeAlpha {
		t.Errorf("anim type = %q, want alpha", btn2.Anims[0].AnimType)
	}
}

func TestCollectAllControls(t *testing.T) {
	menuRoot := controls.NewScreen().ID("menu").Control()
	btnCtrl := controls.NewScreen().ID("btn").Control()

	r := registry.New()
	r.AddScreen("menu", menuRoot).AddScreen("btn", btnCtrl)

	index := registry.CollectAllControls(r.Screens(), r.Nested())
	if len(index) != 2 {
		t.Errorf("expected 2 controls, got %d", len(index))
	}
	if _, ok := index["menu"]; !ok {
		t.Error("expected 'menu' in index")
	}
	if _, ok := index["btn"]; !ok {
		t.Error("expected 'btn' in index")
	}
}

func TestRegisterNested(t *testing.T) {
	nested := controls.NewButton().ID("nested_btn").Control()
	r := registry.New()
	r.AddScreen("parent", controls.NewScreen().ID("parent").Control())
	r.RegisterNested("nested_btn", nested)

	index := registry.CollectAllControls(r.Screens(), r.Nested())
	if _, ok := index["nested_btn"]; !ok {
		t.Error("expected 'nested_btn' in index after RegisterNested")
	}
}

func TestToPack(t *testing.T) {
	r := registry.New()
	r.AddScreen("menu", controls.NewScreen().ID("menu").Control()).
		AddScreen("game", controls.NewScreen().ID("game").Control())
	r.AddLink("play_btn", "game", "fade", 0.5)

	_, err := r.ToPack("my_addon")
	if err == nil {
		t.Error("expected error: source control 'play_btn' not in any screen")
	}

	btnCtrl := controls.NewScreen().ID("play_btn").Control()
	r.AddScreen("play_btn", btnCtrl)
	r.AddAnimationDef("fade")

	out, err := r.ToPack("my_addon")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Screens) != 3 {
		t.Errorf("expected 3 screens, got %d", len(out.Screens))
	}
	if len(out.Transitions) != 1 {
		t.Errorf("expected 1 transition, got %d", len(out.Transitions))
	}
	if out.Transitions[0].From != "play_btn" {
		t.Errorf("transition from = %q, want play_btn", out.Transitions[0].From)
	}
}

func TestFromJSON(t *testing.T) {
	data := []byte(`{"id":"test","type":"panel","visible":true}`)
	ctrl, err := registry.FromJSON(data)
	if err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}
	if ctrl.ID != "test" {
		t.Errorf("ID = %q, want test", ctrl.ID)
	}
	if !ctrl.Visible {
		t.Error("Visible should be true")
	}
}

func TestFromUIDefs(t *testing.T) {
	data := []byte(`{"format_version":1,"ui":["screen_a.json","screen_b.json"],"script_api_version":1}`)
	files, err := registry.FromUIDefs(data)
	if err != nil {
		t.Fatalf("FromUIDefs failed: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	if files[0] != "screen_a.json" {
		t.Errorf("files[0] = %q, want screen_a.json", files[0])
	}
}

func TestLoadScreen(t *testing.T) {
	tmpDir := t.TempDir()
	screenJSON := []byte(`{"id":"main","type":"screen","anchor_from":"center","anchor_to":"center","size":[100,100]}`)
	path := filepath.Join(tmpDir, "main.json")
	if err := os.WriteFile(path, screenJSON, 0644); err != nil {
		t.Fatal(err)
	}

	r := registry.New()
	r.LoadScreen("main", path)
	if len(r.Errors()) > 0 {
		t.Fatalf("load errors: %v", r.Errors())
	}
	screens := r.Screens()
	if _, ok := screens["main"]; !ok {
		t.Error("expected 'main' screen to be loaded")
	}
	if screens["main"].ID != "main" {
		t.Errorf("screen ID = %q, want main", screens["main"].ID)
	}
}

func TestLoadFromDir(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "screen_a.json"), []byte(`{"id":"a","type":"screen"}`), 0644)
	os.WriteFile(filepath.Join(tmpDir, "screen_b.json"), []byte(`{"id":"b","type":"screen"}`), 0644)
	os.WriteFile(filepath.Join(tmpDir, "_ui_defs.json"), []byte(`{"ui":["screen_a.json"]}`), 0644)
	os.WriteFile(filepath.Join(tmpDir, "skip.txt"), []byte("not json"), 0644)

	r := registry.New()
	r.LoadFromDir(tmpDir)
	if len(r.Errors()) > 0 {
		t.Fatalf("load errors: %v", r.Errors())
	}
	screens := r.Screens()
	if _, ok := screens["screen_a"]; !ok {
		t.Error("expected 'screen_a' to be loaded")
	}
	if _, ok := screens["screen_b"]; !ok {
		t.Error("expected 'screen_b' to be loaded")
	}
	if _, ok := screens["_ui_defs"]; ok {
		t.Error("should skip _ui_defs.json")
	}
	if _, ok := screens["skip"]; ok {
		t.Error("should skip non-.json files")
	}
}

func TestResolveIDs(t *testing.T) {
	// parent references "inner" which is not registered anywhere.
	parent := controls.NewPanel().ID("parent").AddControl("inner")
	child := controls.NewPanel().ID("inner")

	r := registry.New()
	r.AddScreen("parent", parent.Control())
	r.AddScreen("inner", child.Control())

	unresolved := r.ResolveIDs()
	if len(unresolved) != 0 {
		t.Errorf("expected 0 unresolved, got %d: %v", len(unresolved), unresolved)
	}
}

func TestResolveIDsUnresolved(t *testing.T) {
	// parent references "btn" which is neither a screen nor nested.
	parent := controls.NewPanel().ID("parent").AddControl("btn")

	r := registry.New()
	r.AddScreen("parent", parent.Control())

	unresolved := r.ResolveIDs()
	if len(unresolved) != 1 || unresolved[0] != "btn" {
		t.Errorf("expected [btn], got %v", unresolved)
	}
}

func TestResolveIDsMixed(t *testing.T) {
	// parent references "known" (registered) and "unknown" (not registered).
	parent := controls.NewPanel().
		ID("parent").
		AddControl("known").
		AddControl("unknown")

	known := controls.NewPanel().ID("known")

	r := registry.New()
	r.AddScreen("parent", parent.Control())
	r.AddScreen("known", known.Control())

	unresolved := r.ResolveIDs()
	if len(unresolved) != 1 || unresolved[0] != "unknown" {
		t.Errorf("expected [unknown], got %v", unresolved)
	}
}

func TestValidateUnresolvedControls(t *testing.T) {
	// parent references "missing" which is not registered.
	parent := controls.NewPanel().ID("parent").AddControl("missing")

	r := registry.New()
	r.AddScreen("parent", parent.Control())

	errs := r.Validate()
	if len(errs) == 0 {
		t.Error("expected validation error for unresolved control")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e, "missing") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error mentioning 'missing', got: %v", errs)
	}
}

func TestValidateWithResolvedNested(t *testing.T) {
	// parent references "btn" which is registered as nested.
	parent := controls.NewPanel().ID("parent").AddControl("btn")
	btn := controls.NewButton().ID("btn")

	r := registry.New()
	r.AddScreen("parent", parent.Control()).
		RegisterNested("btn", btn.Control())

	errs := r.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}

func TestCollectAllControlsResolvesChained(t *testing.T) {
	// a -> b -> c where each is a registered screen.
	a := controls.NewPanel().ID("a").AddControl("b")
	b := controls.NewPanel().ID("b").AddControl("c")
	c := controls.NewPanel().ID("c")

	idx := registry.CollectAllControls(
		map[string]*schema.Control{"a": a.Control(), "b": b.Control(), "c": c.Control()},
		map[string]*schema.Control{},
	)
	if len(idx) != 3 {
		t.Errorf("expected 3 controls, got %d", len(idx))
	}
	for _, id := range []string{"a", "b", "c"} {
		if _, ok := idx[id]; !ok {
			t.Errorf("expected %q in index", id)
		}
	}
}
