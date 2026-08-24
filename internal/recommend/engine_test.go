package recommend

import (
	"gymrecommend/internal/model"
	"testing"
)

func TestRankPrefersGoal(t *testing.T) {
	e := NewEngine()
	m := model.Member{ID: "m", Goals: []string{"strength"}, PreferredSlot: "morning"}
	xs := e.Rank(m, []model.Class{{ID: "a", Title: "Strength Lab", Kind: "group", Slot: "morning", Capacity: 2}, {ID: "b", Title: "Calm", Kind: "group", Slot: "evening", Capacity: 2}}, nil)
	if xs[0].ClassID != "a" {
		t.Fatal(xs)
	}
}
