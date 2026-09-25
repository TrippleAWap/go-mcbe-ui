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
	// 2. Register screens in the registry (fluent — errors
	//    are collected and checked via Errors())
	// =====================================================
	rgstry := registry.New()

	rgstry.AddScreen("main_menu", mainMenu.Control()).
		AddScreen("settings_screen", settingsScreen.Control()).
		AddScreen("play_btn", playBtn.Control()).
		AddScreen("settings_btn", settingsBtn.Control()).
		AddScreen("back_btn", backBtn.Control())

	if errs := rgstry.Validate(); len(errs) > 0 {
		log.Fatalf("Registry validation failed: %v", errs)
	}

	// Add all registered screens to the pack.
	for id, ctrl := range rgstry.Screens() {
		p.AddScreenDef(&pack.ScreenDef{
			ID:        id,
			Namespace: "my_addon",
			Control:   ctrl,
		})
	}

	// =====================================================
	// 3. Component screens (added directly to the pack)
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
	// 4. Validate all screens
	// =====================================================
	for _, def := range p.Screens {
		screen := &schema.Screen{Namespace: def.Namespace, RootPanel: def.ID, Root: *def.Control}
		if r := validate.ValidateScreen(screen); !r.IsValid() {
			log.Printf("Validation warnings for %s: %s", def.ID, r.String())
		}
	}

	// =====================================================
	// 5. Generate the resource pack
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
	fmt.Println("  manifest.json")
}
