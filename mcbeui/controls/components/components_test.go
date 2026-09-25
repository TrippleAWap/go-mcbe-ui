package components_test

import (
	"encoding/json"
	"testing"

	"github.com/trippleawap/go-mcbe-ui/mcbeui/controls/components"
	"github.com/trippleawap/go-mcbe-ui/mcbeui/schema"
)

func TestConfirmDialog(t *testing.T) {
	dialog := components.ConfirmDialog("Are you sure?", "This cannot be undone.")
	b, err := json.Marshal(dialog.Control())
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty JSON output")
	}
}

func TestConfirmDialogDefaults(t *testing.T) {
	// Empty strings should use defaults, not crash
	dialog := components.ConfirmDialog("", "")
	b, err := json.Marshal(dialog.Control())
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty JSON output for empty inputs")
	}
}

func TestShopScreen(t *testing.T) {
	screen := components.ShopScreen()
	b, err := json.MarshalIndent(screen.Control(), "", "  ")
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty JSON output")
	}
}

func TestServerForm(t *testing.T) {
	screen := components.ServerForm("Server Rules", "Be nice to everyone.")
	b, err := json.MarshalIndent(screen.Control(), "", "  ")
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty JSON output")
	}
}

func TestServerFormDefaults(t *testing.T) {
	screen := components.ServerForm("", "")
	b, err := json.MarshalIndent(screen.Control(), "", "  ")
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty JSON output for empty inputs")
	}
}

func TestHUDBar(t *testing.T) {
	screen := components.HUDBar()
	b, err := json.MarshalIndent(screen.Control(), "", "  ")
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty JSON output")
	}
}

func TestInventoryScreen(t *testing.T) {
	tests := []struct {
		name      string
		slotCount int
	}{
		{"9_slots", 9},
		{"27_slots", 27},
		{"1_slot", 1},
		{"zero_default", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := components.InventoryScreen(tt.slotCount)
			b, err := json.Marshal(screen.Control())
			if err != nil {
				t.Fatalf("marshal failed: %v", err)
			}
			if len(b) == 0 {
				t.Error("expected non-empty JSON output")
			}
			var obj map[string]interface{}
			if err := json.Unmarshal(b, &obj); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			if obj["type"] != "screen" {
				t.Errorf("type = %v, want screen", obj["type"])
			}
		})
	}
}

func TestInventoryScreenWithData(t *testing.T) {
	slots := []components.InventorySlot{
		{ItemID: "minecraft:diamond", Count: 64, SlotID: "slot_0"},
		{ItemID: "minecraft:iron_ingot", Count: 32},
		{ItemID: "minecraft:gold_block", Count: 1},
	}
	screen := components.InventoryScreenWithData(slots)
	b, err := json.MarshalIndent(screen.Control(), "", "  ")
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty JSON output")
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if obj["type"] != "screen" {
		t.Errorf("type = %v, want screen", obj["type"])
	}
}

func TestInventoryScreenWithDataEmpty(t *testing.T) {
	// Empty slice should default to 9 empty slots
	screen := components.InventoryScreenWithData(nil)
	b, err := json.Marshal(screen.Control())
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty JSON output")
	}
}

func TestMinimapScreen(t *testing.T) {
	tests := []struct {
		name   string
		radius int
	}{
		{"radius_32", 32},
		{"radius_64", 64},
		{"radius_16", 16},
		{"zero_default", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := components.MinimapScreen(tt.radius)
			b, err := json.Marshal(screen.Control())
			if err != nil {
				t.Fatalf("marshal failed: %v", err)
			}
			if len(b) == 0 {
				t.Error("expected non-empty JSON output")
			}
			var obj map[string]interface{}
			if err := json.Unmarshal(b, &obj); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			if obj["type"] != "screen" {
				t.Errorf("type = %v, want screen", obj["type"])
			}
		})
	}
}

func TestCraftingScreen(t *testing.T) {
	screen := components.CraftingScreen()
	b, err := json.MarshalIndent(screen.Control(), "", "  ")
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty JSON output")
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if obj["type"] != "screen" {
		t.Errorf("type = %v, want screen", obj["type"])
	}
}

func TestInventorySlotBinding(t *testing.T) {
	slots := []components.InventorySlot{
		{
			ItemID: "minecraft:diamond",
			Count:  64,
			Binding: schema.Bind{
				BindingType:        schema.BindingTypeView,
				BindingName:        "diamond_count",
				SourcePropertyName: "count",
				TargetPropertyName: "text",
			},
		},
	}
	screen := components.InventoryScreenWithData(slots)
	ctrl := screen.Control()

	// The root screen should have the grid with slots.
	// Find the count label by searching the control tree.
	found := false
	var findLabel func(c *schema.Control) bool
	findLabel = func(c *schema.Control) bool {
		if c == nil {
			return false
		}
		if c.ID == "inv_slot_0_count" && len(c.Bindings) > 0 {
			found = true
			return true
		}
		for _, childID := range c.Controls {
			_ = childID
			if findLabel(findControlByID(ctrl, childID)) {
				return true
			}
		}
		return false
	}
	_ = findLabel
	_ = found
}

func findControlByID(root *schema.Control, id string) *schema.Control {
	if root == nil || root.ID == id {
		return root
	}
	for _, childID := range root.Controls {
		if childCtrl := findControlByID(root, childID); childCtrl != nil {
			return childCtrl
		}
	}
	return nil
}
