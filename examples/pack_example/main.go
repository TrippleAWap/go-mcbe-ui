// Package main provides the pack generation example.
//
//go:generate go run ../../generate/generate_screen.go --input ../../generate/sample_screen.json --output ../../generate/generated/screen.go --name SampleScreen
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/controls"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/controls/components"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/pack"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/registry"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/transitions"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/validate"
)

func main() {
	p := pack.New("my_addon", "My Resource Pack")

	// =====================================================
	// 1. Build screens from scratch
	// =====================================================
	mainMenu := controls.NewScreen().
		ID("main_menu").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100)).
		AbsorbsInput()

	playBtn := controls.NewButton().
		ID("play_btn").
		Anchor(schema.Center).
		Pos(schema.Raw(0, 20)).
		Size(schema.Raw(200, 50)).
		DefaultControl("btn_play").
		HoverControl("btn_play_hover")

	settingsBtn := controls.NewButton().
		ID("settings_btn").
		Anchor(schema.Center).
		Pos(schema.Raw(0, 80)).
		Size(schema.Raw(200, 50)).
		DefaultControl("btn_settings").
		HoverControl("btn_settings_hover")

	settingsScreen := controls.NewScreen().
		ID("settings_screen").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100)).
		AbsorbsInput()

	backBtn := controls.NewButton().
		ID("back_btn").
		Anchor(schema.BottomCenter).
		Pos(schema.Raw(0, -40)).
		Size(schema.Raw(200, 50)).
		DefaultControl("btn_back")

	// =====================================================
	// 2. Reusable components (added directly to pack — not in registry)
	//    Components have nested Controls[] that cannot be resolved
	//    without RegisterNested, so they bypass the registry.
	// =====================================================
	_ = components.ConfirmDialog("Exit?", "Really quit?") // used below
	_ = components.InventoryScreenWithData([]components.InventorySlot{
		{ItemID: "minecraft:diamond", Count: 64, SlotID: "diamond_slot"},
	})

	// =====================================================
	// 3. Screen navigation with registry
	//    Only screens with fully-resolvable Controls[] should go here.
	// =====================================================
	rgstry := registry.New()

	rgstry.AddScreen("main_menu", mainMenu.Control()).
		AddScreen("settings_screen", settingsScreen.Control()).
		AddScreen("play_btn", playBtn.Control()).
		AddScreen("settings_btn", settingsBtn.Control()).
		AddScreen("back_btn", backBtn.Control()).
		AddAnimationDef(string(transitions.Fade)).
		AddAnimationDef(string(transitions.SlideLeft))

	// Add navigation links using typed transition constants
	rgstry.AddLink("play_btn", "settings_screen", string(transitions.Fade), 0.5).
		AddLink("settings_btn", "settings_screen", string(transitions.Fade), 0.3).
		AddLink("back_btn", "main_menu", string(transitions.SlideLeft), 0.3)

	// Validate navigation links (also validates transitions against anim defs)
	if errs := rgstry.Validate(); len(errs) > 0 {
		log.Fatalf("Registry validation failed: %v", errs)
	}
	rgstry.InjectAnimations()

	// Register screens from registry output
	for id, ctrl := range rgstry.Screens() {
		p.AddScreenDef(&pack.ScreenDef{
			ID:        id,
			Namespace: "my_addon",
			Control:   ctrl,
		})
	}

	// Register transitions
	for _, tr := range rgstry.Links() {
		p.AddTransition(pack.ScreenTransition{
			From:        tr.FromID,
			To:          tr.ToScreen,
			AnimationID: tr.AnimID,
			Duration:    tr.Duration,
		})
	}

	// =====================================================
	// 4. Additional component screens (added directly to pack)
	// =====================================================
	shop := components.ShopScreen()
	p.AddScreenDef(&pack.ScreenDef{
		ID:        "shop",
		Namespace: "my_addon",
		Control:   shop.Control(),
	})

	serverForm := components.ServerForm("Server Rules", "Be respectful to all players.")
	p.AddScreenDef(&pack.ScreenDef{
		ID:        "server_rules",
		Namespace: "my_addon",
		Control:   serverForm.Control(),
	})

	craft := components.CraftingScreen()
	p.AddScreenDef(&pack.ScreenDef{
		ID:        "crafting",
		Namespace: "my_addon",
		Control:   craft.Control(),
	})

	minimap := components.MinimapScreen(64)
	p.AddScreenDef(&pack.ScreenDef{
		ID:        "minimap",
		Namespace: "my_addon",
		Control:   minimap.Control(),
	})

	hud := components.HUDBar()
	p.AddScreenDef(&pack.ScreenDef{
		ID:        "hud_bar",
		Namespace: "my_addon",
		Control:   hud.Control(),
	})

	// =====================================================
	// 5. Pack metadata
	// =====================================================
	p.AddTexture("textures/ui/inventory/slot")
	p.AddTexture("textures/ui/hud/health_bar")
	p.AddTexture("textures/ui/minimap/circle_bg")

	// =====================================================
	// 6. Validate all screens
	// =====================================================
	for _, def := range p.Screens {
		screen := &schema.Screen{Namespace: def.Namespace, RootPanel: def.ID, Root: *def.Control}
		if r := validate.ValidateScreen(screen); !r.IsValid() {
			log.Printf("Validation warnings for %s: %s", def.ID, r.String())
		}
	}

	// =====================================================
	// 7. Generate
	// =====================================================
	outDir := "./my_pack_output"
	os.RemoveAll(outDir)
	if err := p.Build(outDir); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Pack generated to %s/\n", outDir)
	fmt.Println("Screens:")
	for _, def := range p.Screens {
		fmt.Printf("  ui/%s.json\n", def.ID)
	}
	fmt.Println("  ui/_ui_defs.json")
	fmt.Println("  pack.mcmeta")
	fmt.Printf("\nNavigation links: %d\n", len(rgstry.Links()))
	fmt.Printf("Registered textures: %d\n", len(p.Textures))
	fmt.Printf("Animation defs: %d\n", len(p.AnimationDefs))
	fmt.Printf("Transitions: %d\n", len(p.Transitions))
}
