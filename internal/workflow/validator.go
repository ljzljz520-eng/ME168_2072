package workflow

import (
	"fmt"
	"gymrecommend/internal/model"
)

func ValidateSequence(steps []string) error {
	if len(steps) < 4 {
		return fmt.Errorf("workflow requires four steps")
	}
	for _, s := range steps {
		if s == "" {
			return fmt.Errorf("empty step")
		}
	}
	return nil
}
func ValidateRecommendationSet(items []model.Recommendation) error {
	seen := map[string]bool{}
	for _, r := range items {
		if e := model.ValidateRecommendation(r); e != nil {
			return e
		}
		if seen[r.ID] {
			return fmt.Errorf("duplicate recommendation")
		}
		seen[r.ID] = true
	}
	return nil
}
