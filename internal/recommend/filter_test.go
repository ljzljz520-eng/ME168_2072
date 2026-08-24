package recommend

import (
	"gymrecommend/internal/model"
	"testing"
)

func TestFilters(t *testing.T) {
	xs := []model.Recommendation{{ID: "a", Kind: "group", Score: .8}, {ID: "b", Kind: "private", Score: .4}}
	if len(FilterByKind(xs, "private")) != 1 || len(FilterByScore(xs, .7)) != 1 {
		t.Fatal("filter")
	}
}
