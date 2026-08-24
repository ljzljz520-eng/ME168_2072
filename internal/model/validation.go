package model

import "fmt"

func ValidateMember(m Member) error {
	if m.ID == "" {
		return fmt.Errorf("member id required")
	}
	if m.Name == "" {
		return fmt.Errorf("member name required")
	}
	if m.CoachRating < 0 || m.CoachRating > 5 {
		return fmt.Errorf("rating out of range")
	}
	return nil
}
func ValidateClass(c Class) error {
	if c.ID == "" || c.Title == "" {
		return fmt.Errorf("class identity required")
	}
	if c.Capacity <= 0 {
		return fmt.Errorf("capacity must be positive")
	}
	if c.Kind != "group" && c.Kind != "private" {
		return fmt.Errorf("unknown class kind")
	}
	return nil
}
func ValidateRecommendation(r Recommendation) error {
	if r.MemberID == "" || r.ClassID == "" {
		return fmt.Errorf("recommendation links required")
	}
	if r.Score < 0 || r.Score > 1 {
		return fmt.Errorf("score out of range")
	}
	return nil
}
func NormalizeGoals(goals []string) []string {
	out := make([]string, 0, len(goals))
	seen := map[string]bool{}
	for _, g := range goals {
		if g != "" && !seen[g] {
			out = append(out, g)
			seen[g] = true
		}
	}
	return out
}
func SlotCompatible(want, have string) bool {
	if want == "" || have == "" {
		return false
	}
	return want == have
}
