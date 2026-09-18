// Package components provides reusable MCBE UI screen patterns.
package components

import (
	"fmt"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/controls"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
)

// ConfirmDialog builds a centered confirmation dialog with a title, message, and Yes/No buttons.
// Returns the root panel containing the dialog. Add it to your screen with screen.AddControl(dialog.ID).
func ConfirmDialog(title, message string) *controls.Panel {
	if title == "" {
		title = "Confirm"
	}
	if message == "" {
		message = "Are you sure?"
	}
	bg := controls.NewPanel().
		ID("confirm_bg").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100)).
		Alpha(0.5)

	panel := controls.NewPanel().
		ID("confirm_panel").
		Anchor(schema.Center).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(400, 250))

	titleLbl := controls.NewLabel().
		ID("confirm_title").
		Anchor(schema.TopMiddle).
		Pos(schema.Raw(0, 20)).
		Size(schema.Raw(360, 50)).
		Text(title).
		FontSize(schema.FontSizeLarge).
		Align(schema.TextAlignCenter).
		Localize()

	msgLbl := controls.NewLabel().
		ID("confirm_msg").
		Anchor(schema.TopMiddle).
		Pos(schema.Raw(0, 80)).
		Size(schema.Raw(360, 100)).
		Text(message).
		FontSize(schema.FontSizeNormal).
		Align(schema.TextAlignCenter).
		Shadow()

	yesBtn := controls.NewButton().
		ID("confirm_yes").
		Anchor(schema.BottomMiddle).
		Pos(schema.Raw(-90, -40)).
		Size(schema.Raw(130, 50)).
		DefaultControl("btn_yes").
		HoverControl("btn_yes_hover")

	noBtn := controls.NewButton().
		ID("confirm_no").
		Anchor(schema.BottomMiddle).
		Pos(schema.Raw(90, -40)).
		Size(schema.Raw(130, 50)).
		DefaultControl("btn_no").
		HoverControl("btn_no_hover")

	panel.AddControl(titleLbl.Control().ID)
	panel.AddControl(msgLbl.Control().ID)
	panel.AddControl(yesBtn.Control().ID)
	panel.AddControl(noBtn.Control().ID)

	root := controls.NewPanel().
		ID("confirm_dialog_root").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100))
	root.AddControl(bg.Control().ID)
	root.AddControl(panel.Control().ID)

	return root
}

// ShopScreen builds a basic shop layout with a header, item list, and bottom bar.
func ShopScreen() *controls.Screen {
	screen := controls.NewScreen().
		ID("shop_screen").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100)).
		RenderGameBehind()

	header := controls.NewPanel().
		ID("shop_header").
		Anchor(schema.TopLeft).
		Size(schema.Raw(100, 60)).
		Contained()

	headerTitle := controls.NewLabel().
		ID("shop_header_title").
		Anchor(schema.TopMiddle).
		Pos(schema.Raw(0, 10)).
		Size(schema.Raw(300, 40)).
		Text("Shop").
		FontSize(schema.FontSizeLarge).
		Align(schema.TextAlignCenter).
		Localize()
	header.AddControl(headerTitle.Control().ID)

	itemList := controls.NewScrollView().
		ID("shop_items").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(20, 70)).
		Size(schema.Raw(360, 500)).
		ScrollbarAlwaysVisible()

	content := controls.NewStackPanel().
		ID("shop_item_list").
		Anchor(schema.TopLeft).
		Size(schema.Fill()).
		Vertical()

	itemTemplate := controls.NewStackPanel().
		ID("shop_item_template").
		Anchor(schema.TopLeft).
		Size(schema.Raw(100, 60)).
		Horizontal()

	nameLabel := controls.NewLabel().
		ID("shop_item_name").
		Anchor(schema.TopLeft).
		Pos(schema.Pad(10, 5)).
		Size(schema.Raw(200, 50)).
		FontSize(schema.FontSizeNormal)

	buyBtn := controls.NewButton().
		ID("shop_buy_btn").
		Anchor(schema.TopRight).
		Pos(schema.Margin(10, 5)).
		Size(schema.Raw(100, 50)).
		DefaultControl("btn_buy")

	itemTemplate.AddControl(nameLabel.Control().ID)
	itemTemplate.AddControl(buyBtn.Control().ID)
	content.AddControl(itemTemplate.Control().ID)
	itemList.Content(content.Control().ID)
	itemList.ViewPort("shop_view_port")

	screen.AddControl(header.Control().ID)
	screen.AddControl(itemList.Control().ID)

	return screen
}

// ServerForm builds a server message/disclaimer screen with a title, message, and acknowledge button.
func ServerForm(title, message string) *controls.Screen {
	if title == "" {
		title = "Server Message"
	}
	if message == "" {
		message = "Please read the rules."
	}
	screen := controls.NewScreen().
		ID("server_form").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100)).
		AbsorbsInput()

	container := controls.NewPanel().
		ID("server_form_container").
		Anchor(schema.Center).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(500, 350)).
		ClipsChildren()

	formTitle := controls.NewLabel().
		ID("server_form_title").
		Anchor(schema.TopMiddle).
		Pos(schema.Raw(0, 20)).
		Size(schema.Raw(460, 60)).
		Text(title).
		FontSize(schema.FontSizeExtraLarge).
		Align(schema.TextAlignCenter).
		Localize()

	formMessage := controls.NewScrollView().
		ID("server_form_message_scroll").
		Anchor(schema.TopMiddle).
		Pos(schema.Raw(0, 90)).
		Size(schema.Raw(460, 180)).
		ScrollbarAlwaysVisible()

	msgContent := controls.NewPanel().
		ID("server_form_message_content").
		Anchor(schema.TopLeft).
		Size(schema.Raw(460, 180))

	msgLabel := controls.NewLabel().
		ID("server_form_message_label").
		Anchor(schema.TopLeft).
		Pos(schema.Pad(15, 10)).
		Size(schema.Raw(430, 160)).
		Text(message).
		FontSize(schema.FontSizeSmall).
		Localize()

	msgContent.AddControl(msgLabel.Control().ID)
	formMessage.Content(msgContent.Control().ID)
	formMessage.ViewPort("server_form_view_port")

	ackBtn := controls.NewButton().
		ID("server_form_ack").
		Anchor(schema.BottomMiddle).
		Pos(schema.Raw(0, -40)).
		Size(schema.Raw(200, 50)).
		DefaultControl("btn_ack").
		HoverControl("btn_ack_hover")

	container.AddControl(formTitle.Control().ID)
	container.AddControl(formMessage.Control().ID)
	container.AddControl(ackBtn.Control().ID)
	screen.AddControl(container.Control().ID)

	return screen
}

// MinimapScreen builds a circular minimap overlay with a radial texture and directional arrow.
func MinimapScreen(radius int) *controls.Screen {
	if radius <= 0 {
		radius = 48
	}
	screen := controls.NewScreen().
		ID("minimap_screen").
		Anchor(schema.TopRight).
		Pos(schema.Raw(-radius-10, 10)).
		Size(schema.Raw(radius*2, radius*2)).
		AbsorbsInput()

	container := controls.NewPanel().
		ID("minimap_container").
		Anchor(schema.Center).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(radius*2, radius*2)).
		ClipsChildren()

	bg := controls.NewImage().
		ID("minimap_bg").
		Anchor(schema.Center).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(radius*2, radius*2)).
		Texture("textures/ui/minimap/circle_bg").
		Tiled()

	mapSurface := controls.NewPanel().
		ID("minimap_surface").
		Anchor(schema.Center).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(radius*2, radius*2))

	indicators := controls.NewStackPanel().
		ID("minimap_indicators").
		Anchor(schema.Center).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(radius*2, radius*2))

	arrow := controls.NewImage().
		ID("minimap_arrow").
		Anchor(schema.TopMiddle).
		Pos(schema.Raw(0, radius/2)).
		Size(schema.Raw(16, 16)).
		Texture("textures/ui/minimap/arrow")

	directionLabel := controls.NewLabel().
		ID("minimap_direction").
		Anchor(schema.BottomMiddle).
		Pos(schema.Raw(0, -radius/2-5)).
		Size(schema.Raw(60, 20)).
		Text("N").
		FontSize(schema.FontSizeTiny).
		Align(schema.TextAlignCenter)

	indicators.AddControl(arrow.Control().ID)
	indicators.AddControl(directionLabel.Control().ID)
	mapSurface.AddControl(indicators.Control().ID)
	container.AddControl(bg.Control().ID)
	container.AddControl(mapSurface.Control().ID)
	screen.AddControl(container.Control().ID)

	return screen
}

// CraftingScreen builds a crafting layout with a 2x2 grid on the left, an arrow button in the
// middle, and an output slot on the right. Includes a title label and a result label below the
// output slot.
func CraftingScreen() *controls.Screen {
	screen := controls.NewScreen().
		ID("crafting_screen").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100)).
		RenderGameBehind()

	title := controls.NewLabel().
		ID("crafting_title").
		Anchor(schema.TopMiddle).
		Pos(schema.Raw(0, 10)).
		Size(schema.Raw(300, 40)).
		Text("Crafting").
		FontSize(schema.FontSizeLarge).
		Align(schema.TextAlignCenter).
		Localize()
	screen.AddControl(title.Control().ID)

	// Left: 2x2 crafting grid container
	gridContainer := controls.NewPanel().
		ID("crafting_grid_container").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(20, 60)).
		Size(schema.Raw(140, 140))

	grid := controls.NewGrid().
		ID("crafting_grid").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(140, 140)).
		Dimensions(2, 2).
		FillDirection(schema.GridFillDown)

	slotSize := 60
	slotGap := 4
	positions := [4][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
	for _, pos := range positions {
		slot := controls.NewPanel().
			ID(fmt.Sprintf("craft_slot_%d_%d", pos[0], pos[1])).
			Anchor(schema.TopLeft).
			Pos(schema.Raw(pos[0]*(slotSize+slotGap), pos[1]*(slotSize+slotGap))).
			Size(schema.Raw(slotSize, slotSize))
		slotBg := controls.NewImage().
			ID(fmt.Sprintf("craft_slot_bg_%d_%d", pos[0], pos[1])).
			Anchor(schema.Center).
			Pos(schema.Raw(0, 0)).
			Size(schema.Raw(slotSize, slotSize)).
			Texture("textures/ui/inventory/slot")
		slot.AddControl(slotBg.Control().ID)
		grid.AddControl(slot.Control().ID)
	}
	gridContainer.AddControl(grid.Control().ID)
	screen.AddControl(gridContainer.Control().ID)

	// Middle: arrow button
	arrowBtn := controls.NewButton().
		ID("crafting_arrow").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(170, 95)).
		Size(schema.Raw(30, 50)).
		DefaultControl("btn_craft_arrow").
		HoverControl("btn_craft_arrow_hover")
	screen.AddControl(arrowBtn.Control().ID)

	// Right: output slot container
	outputContainer := controls.NewPanel().
		ID("crafting_output_container").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(210, 60)).
		Size(schema.Raw(80, 140))

	outputSlot := controls.NewPanel().
		ID("crafting_output_slot").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(80, 80))
	outputSlotBg := controls.NewImage().
		ID("crafting_output_bg").
		Anchor(schema.Center).
		Pos(schema.Raw(0, 0)).
		Size(schema.Raw(80, 80)).
		Texture("textures/ui/inventory/slot")
	outputSlot.AddControl(outputSlotBg.Control().ID)
	outputContainer.AddControl(outputSlot.Control().ID)
	screen.AddControl(outputContainer.Control().ID)

	resultLabel := controls.NewLabel().
		ID("crafting_result_label").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(210, 150)).
		Size(schema.Raw(80, 30)).
		Text("Result").
		FontSize(schema.FontSizeSmall).
		Align(schema.TextAlignCenter)
	screen.AddControl(resultLabel.Control().ID)

	return screen
}

// InventorySlot represents a single item in an inventory screen.
type InventorySlot struct {
	ItemID  string       // e.g. "minecraft:diamond"
	Count   int          // stack count
	SlotID  string       // custom control ID; auto-generated if empty
	Binding schema.Bind  // optional binding to keep count in sync with data
}

// InventoryScreenWithData builds an inventory grid populated with actual item data.
// Each slot shows its item ID as text and the stack count.
// slotCount must be positive; defaults to 27 if <= 0.
func InventoryScreenWithData(slots []InventorySlot) *controls.Screen {
	slotCount := len(slots)
	if slotCount <= 0 {
		slotCount = 9
	}
	cols := 9
	if slotCount <= 9 {
		cols = slotCount
	} else if slotCount <= 18 {
		cols = 9
	} else if slotCount <= 27 {
		cols = 9
	} else {
		cols = 9
	}
	rows := (slotCount + cols - 1) / cols

	screen := controls.NewScreen().
		ID("inventory_screen").
		Anchor(schema.Center).
		Size(schema.Raw(100, 100)).
		RenderGameBehind()

	title := controls.NewLabel().
		ID("inventory_title").
		Anchor(schema.TopMiddle).
		Pos(schema.Raw(0, 10)).
		Size(schema.Raw(400, 40)).
		Text("Inventory").
		FontSize(schema.FontSizeLarge).
		Align(schema.TextAlignCenter).
		Localize()
	screen.AddControl(title.Control().ID)

	slotSize := 40
	slotGap := 4
	totalW := cols*slotSize + (cols-1)*slotGap
	totalH := rows*slotSize + (rows-1)*slotGap
	startX := 10
	startY := 60

	grid := controls.NewGrid().
		ID("inventory_grid").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(startX, startY)).
		Size(schema.Raw(totalW, totalH)).
		Dimensions(cols, rows).
		FillDirection(schema.GridFillDown)

	for i, slot := range slots {
		id := slot.SlotID
		if id == "" {
			id = fmt.Sprintf("inv_slot_%d", i)
		}
		slotPanel := controls.NewPanel().
			ID(id).
			Anchor(schema.TopLeft).
			Pos(schema.Raw((i%cols)*(slotSize+slotGap), (i/cols)*(slotSize+slotGap))).
			Size(schema.Raw(slotSize, slotSize))

		slotBg := controls.NewImage().
			ID(id + "_bg").
			Anchor(schema.Center).
			Pos(schema.Raw(0, 0)).
			Size(schema.Raw(slotSize, slotSize)).
			Texture("textures/ui/inventory/slot")

		itemLabel := controls.NewLabel().
			ID(id + "_item").
			Anchor(schema.Center).
			Pos(schema.Raw(0, -5)).
			Size(schema.Raw(slotSize-4, 20)).
			Text(slot.ItemID).
			FontSize(schema.FontSizeTiny).
			Align(schema.TextAlignCenter)

		countLabel := controls.NewLabel().
			ID(id + "_count").
			Anchor(schema.BottomRight).
			Pos(schema.Raw(-2, -2)).
			Size(schema.Raw(20, 16)).
			FontSize(schema.FontSizeTiny).
			Align(schema.TextAlignRight)

		if slot.Binding.BindingType != "" {
			countLabel.Bindings(slot.Binding)
		} else {
			countLabel.Text(fmt.Sprintf("%d", slot.Count))
		}

		slotPanel.AddControl(slotBg.Control().ID)
		slotPanel.AddControl(itemLabel.Control().ID)
		slotPanel.AddControl(countLabel.Control().ID)
		grid.AddControl(slotPanel.Control().ID)
	}

	screen.AddControl(grid.Control().ID)
	return screen
}

// InventoryScreen is an alias for InventoryScreenWithData with no items.
// Use InventoryScreenWithData for a populated inventory.
func InventoryScreen(slotCount int) *controls.Screen {
	slots := make([]InventorySlot, 0, slotCount)
	for i := 0; i < slotCount; i++ {
		slots = append(slots, InventorySlot{ItemID: "empty", Count: 0})
	}
	return InventoryScreenWithData(slots)
}

// HUDBar builds a minimal HUD with health, hunger, and experience bars at the bottom center.
func HUDBar() *controls.Screen {
	screen := controls.NewScreen().
		ID("hud_bar").
		Anchor(schema.BottomLeft).
		Size(schema.Raw(100, 80)).
		AbsorbsInput()

	barContainer := controls.NewPanel().
		ID("hud_bar_container").
		Anchor(schema.BottomCenter).
		Pos(schema.Raw(0, 5)).
		Size(schema.Raw(360, 30))

	healthBg := controls.NewImage().
		ID("hud_health_bg").
		Anchor(schema.BottomLeft).
		Pos(schema.Margin(0, 0)).
		Size(schema.Raw(182, 20)).
		Texture("textures/ui/hud/health_bar")

	hungerBg := controls.NewImage().
		ID("hud_hunger_bg").
		Anchor(schema.BottomRight).
		Pos(schema.Margin(0, 0)).
		Size(schema.Raw(182, 20)).
		Texture("textures/ui/hud/hunger_bar")

	expBar := controls.NewPanel().
		ID("hud_exp_bar").
		Anchor(schema.TopLeft).
		Pos(schema.Raw(0, -30)).
		Size(schema.Raw(360, 5)).
		Contained()

	expFill := controls.NewPanel().
		ID("hud_exp_fill").
		Anchor(schema.TopLeft).
		Size(schema.Raw(100, 100))

	expBorder := controls.NewImage().
		ID("hud_exp_border").
		Anchor(schema.TopLeft).
		Size(schema.FillX()).
		Texture("textures/ui/hud/experience_bar_border").
		Tiled()

	expBar.AddControl(expFill.Control().ID)
	expBar.AddControl(expBorder.Control().ID)
	barContainer.AddControl(healthBg.Control().ID)
	barContainer.AddControl(hungerBg.Control().ID)
	barContainer.AddControl(expBar.Control().ID)
	screen.AddControl(barContainer.Control().ID)

	return screen
}
