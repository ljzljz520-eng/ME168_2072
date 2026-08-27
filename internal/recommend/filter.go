package recommend

import "gymrecommend/internal/model"

func FilterByKind(items []model.Recommendation, kind string) []model.Recommendation {
	out := make([]model.Recommendation, 0)
	for _, r := range items {
		if kind == "" || r.Kind == kind {
			out = append(out, r)
		}
	}
	return out
}
func FilterByScore(items []model.Recommendation, min float64) []model.Recommendation {
	out := make([]model.Recommendation, 0)
	for _, r := range items {
		if r.Score >= min {
			out = append(out, r)
		}
	}
	return out
}
func MergePreferences(goals []string, slot string) model.GoalProfile {
	return model.GoalProfile{Goals: model.NormalizeGoals(goals), Availability: []string{slot}, Intensity: inferIntensity(goals)}
}
func inferIntensity(goals []string) string {
	for _, g := range goals {
		if g == "strength" || g == "performance" {
			return "high"
		}
	}
	return "moderate"
}
