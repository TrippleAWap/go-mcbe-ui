# go-mcbe-ui

A Go library for building Minecraft Bedrock Edition (MCBE) JSON UI packs.

## Installation

```bash
go get github.com/trippleawap/go-mcbe-ui
```

## Architecture

```
schema/        Core types: Relative, Anchor, Control, Bind, Anim, Screen, layout helpers
controls/      Type-gated builders: Panel, Label, Image, Button, etc.
  components/  Reusable screens: ConfirmDialog, ShopScreen, ServerForm, HUDBar, Inventory, Crafting, Minimap
registry/      Screen registry: manages screens, validates control references, loads existing packs
pack/          Pack generator: _ui_defs.json, screen files, manifest.json
validate/      Schema-aware validation
```

## Quick Start

```go
import (
    "github.com/trippleawap/go-mcbe-ui/mcbeui/controls"
    "github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
)

root := controls.NewScreen().
    ID("my_screen").
    Anchor(schema.Center).
    Size(schema.Raw(100, 100)).
    AbsorbsInput()

dialog := controls.NewPanel().
    ID("dialog").
    Anchor(schema.Center).
    Pos(schema.Raw(0, 0)).
    Size(schema.Raw(600, 400)).
    ClipsChildren()

title := controls.NewLabel().
    ID("title").
    Anchor(schema.TopMiddle).
    Pos(schema.Raw(0, 20)).
    Size(schema.Raw(500, 60)).
    Text("Hello MCBE").
    FontSize(schema.FontSizeLarge).
    Align(schema.TextAlignCenter).
    Localize()

dialog.AddControl(title.Control().ID)
root.AddControl(dialog.Control().ID)
```

Each control type only exposes methods valid for that type. Calling `Text()` on a `Panel` won't compile.

## Relative Measurements

```go
schema.Raw(100, 200)           // [100, 200] — pixels
schema.Raw("fill", "parent")   // ["fill", "parent"]
schema.Raw(100, "100%")        // [100, "100%"] — mixed
schema.Raw("100%c", "100%cm")  // content-relative
schema.Raw("100%x", "100%y")   // cross-axis
schema.Raw("default")          // single-axis: ["default"]
schema.Raw(100, "100% - 4px")  // arithmetic expression
```

Single-axis helpers:

```go
schema.Pixel(100)              // [100]
schema.Percent(75)             // ["75%"]
schema.Named("fill")           // ["fill"]
schema.Pixel(100).SetY(schema.Named("fill")) // [100, "fill"]
```

## Layout Helpers

```go
schema.Pad(10, 5)                          // [10, 5] — inset offset for padding
schema.Margin(5, 5)                        // [5, 5] — outward offset for margins
schema.CenterIn(800, 600, 400, 300)        // [200, 150] — center a child in a parent
schema.HalfSize(100, 80)                   // [50, 40]
schema.Fill()                              // ["fill", "fill"]
schema.FillX()                             // ["fill", "parent"]
schema.FillY()                             // ["parent", "fill"]
```

## Build a Complete Pack with Registry

```go
import (
    "github.com/trippleawap/go-mcbe-ui/mcbeui/controls"
    "github.com/trippleawap/go-mcbe-ui/mcbeui/pack"
    "github.com/trippleawap/go-mcbe-ui/mcbeui/registry"
)

// Build screens
menu := controls.NewScreen().ID("menu").Anchor(schema.Center).Size(schema.Raw(100, 100))
game := controls.NewScreen().ID("game").Anchor(schema.Center).Size(schema.Raw(100, 100))
playBtn := controls.NewButton().ID("play_btn").Control()

// Create registry and add screens (fluent — errors collected, check via Errors())
rgstry := registry.New()
rgstry.AddScreen("menu", menu.Control()).
       AddScreen("game", game.Control()).
       AddScreen("play_btn", playBtn)

// ToPack validates the registry (screens and control references)
out, err := rgstry.ToPack("my_addon")
if err != nil {
    log.Fatalf("validation failed: %v", err)
}

// Generate the resource pack to disk
p := pack.New("my_addon", "My Resource Pack")
for _, def := range out.Screens {
    p.AddScreenDef(&pack.ScreenDef{
        ID:        def.ID,
        Namespace: def.Namespace,
        Control:   def.Control,
    })
}
if err := p.Build("./my_pack_output"); err != nil {
    log.Fatal(err)
}
```

## Loading Existing Packs

```go
// Load screen JSON files from a directory into a registry
rgstry := registry.New()
rgstry.LoadFromDir("./existing_pack/ui")

// Or load a single file
rgstry.LoadScreen("my_screen", "./existing_pack/ui/my_screen.json")

// Parse _ui_defs.json to get screen filenames
files, _ := registry.FromUIDefs(data)
```

## Nested Controls Limitation

MCBE JSON UI stores child references as string IDs (`Controls []string`). This means
controls nested inside other controls cannot be automatically resolved by ID unless
they are registered as separate screens or added via `RegisterNested()`:

```go
// If a button is nested inside a panel and you need to reference it:
rgstry.RegisterNested("nested_button", someButton.Control())
```

## Reusable Components

```go
import "github.com/trippleawap/go-mcbe-ui/mcbeui/controls/components"

dialog   := components.ConfirmDialog("Exit?", "Are you sure?")
shop     := components.ShopScreen()
form     := components.ServerForm("Rules", "Be nice.")
hud      := components.HUDBar()
inv      := components.InventoryScreenWithData([]components.InventorySlot{
    {ItemID: "minecraft:diamond", Count: 64},
})
craft    := components.CraftingScreen()
minimap  := components.MinimapScreen(64)
```

All components handle empty/zero inputs gracefully (defaults are used).

## Supported Controls

| Builder | Methods |
|---------|---------|
| `controls.NewPanel()` | `Pos`, `Size`, `ClipsChildren`, `Contained`, `AddControl` |
| `controls.NewLabel()` | `Text`, `FontSize`, `FontScale`, `Align`, `Localize`, `Shadow` |
| `controls.NewImage()` | `Texture`, `UV`, `NineSlice`, `Color`, `KeepRatio`, `Fill` |
| `controls.NewButton()` | `DefaultControl`, `HoverControl`, `PressedControl` |
| `controls.NewStackPanel()` | `Horizontal`, `Vertical` |
| `controls.NewGrid()` | `Dimensions`, `MaxItems`, `ItemTemplate`, `FillDirection` |
| `controls.NewToggle()` | `Name`, `RadioGroup`, `OnButton`, `OffButton`, `CheckedControl` |
| `controls.NewSlider()` | `TrackButton`, `SelectedButton`, `Steps` |
| `controls.NewEditBox()` | `BoxName`, `TextType`, `MaxLength`, `Multiline` |
| `controls.NewScrollView()` | `ViewPort`, `Content`, `ScrollbarBox`, `ScrollbarAlwaysVisible` |
| `controls.NewScreen()` | `Modal`, `AbsorbsInput`, `RenderGameBehind` |
| `controls.NewCustom()` | `Renderer` |
| `controls.NewCollectionPanel()` | `Collection` |
| `controls.NewDropdown()` | `Name` |
| `controls.NewSelectionWheel()` | `Slices` |

Every builder has `.Control()` to get the underlying `*schema.Control`.

## Data Bindings & Animations

```go
// Bindings
label.Bindings(
    schema.Bind{BindingType: schema.BindingTypeView, SourcePropertyName: "src", TargetPropertyName: "tgt"},
)

// Animations
img.Animations(
    schema.AnimAlpha(0.5, 0.0, 1.0),
    schema.AnimSize(1.0, schema.Vector2{X: 0, Y: 0}, schema.Vector2{X: 100, Y: 100}),
)
```

## Validation

```go
import "github.com/trippleawap/go-mcbe-ui/mcbeui/validate"

// Basic validation (anchors, sizes, bindings, animations)
res := validate.ValidateScreen(screen)
if res.IsValid() { ... }

// Deep validation with cross-control reference checking
allControls := map[string]*schema.Control{
    "btn": button.Control(),
    "inner": innerPanel.Control(),
}
deepRes := validate.ValidateScreenWithControls(screen, allControls)
```

## Registry Features

The `registry` package provides:
- **Screen management**: `AddScreen()` is fluent — errors are collected and checkable via `Errors()`
- **Nested controls**: `RegisterNested()` adds controls not registered as top-level screens
- **Loading**: `LoadScreen()` and `LoadFromDir()` read existing JSON UI files into the registry
- **Validation**: `Validate()` checks that all `Controls[]` references resolve to a registered screen or nested control
- **Pack output**: `ToPack(namespace)` validates everything and returns `*PackOutput` with screens
- **Flat control index**: `CollectAllControls()` resolves all control IDs across screens via iterative fixpoint
- **Unresolved ID detection**: `ResolveIDs()` reports Controls[] references that have no matching screen/nested control

## Project Structure

```
mcbeui/
  schema/        Anchor, Relative, Control, Bind, Anim, Screen, layout helpers
  controls/      Type-gated builders (Panel, Label, Image, Button, ...)
    components/  Reusable screens (ConfirmDialog, ShopScreen, ServerForm, HUDBar, Inventory, Crafting, Minimap)
  registry/      Screen registry with validation and pack loading
  pack/          Pack generator (_ui_defs.json, screen files, manifest.json)
  validate/      Schema-aware validation
examples/        pack_example — full workflow demo
```

## Limitations

- **`Controls []string`**: Child references are string IDs, not pointers. Nested controls
  can only be resolved if they are registered as screens or added via `RegisterNested()`.
  Use `ResolveIDs()` to discover unresolved references before calling `Validate()`.
