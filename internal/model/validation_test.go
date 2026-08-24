package model

import "testing"

func TestNormalizeGoals(t *testing.T) {
	g := NormalizeGoals([]string{"strength", "", "strength", "mobility"})
	if len(g) != 2 {
		t.Fatalf("got %d", len(g))
	}
	if !SlotCompatible("morning", "morning") {
		t.Fatal("slot")
	}
}
