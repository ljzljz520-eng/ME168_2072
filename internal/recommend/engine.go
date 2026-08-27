package recommend

import (
	"fmt"
	"gymrecommend/internal/model"
	"sort"
	"strings"
	"time"
)

type Engine struct{ weights map[string]float64 }

func NewEngine() *Engine {
	return &Engine{weights: map[string]float64{"goal": 0.45, "slot": 0.25, "rating": 0.2, "kind": 0.1}}
}
func (e *Engine) Rank(m model.Member, classes []model.Class, visits []model.VisitRecord) []model.Recommendation {
	out := make([]model.Recommendation, 0, len(classes))
	for _, c := range classes {
		score, reason := e.score(m, c, visits)
		out = append(out, model.Recommendation{ID: fmt.Sprintf("%s-%s", m.ID, c.ID), MemberID: m.ID, ClassID: c.ID, Kind: c.Kind, Reason: reason, Score: score, CreatedAt: time.Now()})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}
func (e *Engine) score(m model.Member, c model.Class, visits []model.VisitRecord) (float64, string) {
	score := 0.0
	reasons := []string{}
	if model.SlotCompatible(m.PreferredSlot, c.Slot) {
		score += e.weights["slot"]
		reasons = append(reasons, "matches preferred time")
	}
	if c.IsAvailable() {
		score += 0.1
	} else {
		return score, "class is full"
	}
	if m.CoachRating > 4 && c.Kind == "private" {
		score += e.weights["rating"]
		reasons = append(reasons, "high coach quality preference")
	}
	if m.HasGoal("strength") && strings.Contains(strings.ToLower(c.Title), "strength") {
		score += e.weights["goal"]
		reasons = append(reasons, "supports strength goal")
	}
	if m.HasGoal("mobility") && strings.Contains(strings.ToLower(c.Title), "mobility") {
		score += e.weights["goal"]
		reasons = append(reasons, "supports mobility goal")
	}
	if len(visits) > 0 {
		score += 0.05
		reasons = append(reasons, "builds on visit history")
	}
	if score > 1 {
		score = 1
	}
	return score, strings.Join(reasons, "; ")
}
func (e *Engine) Evaluate(recs []model.Recommendation) map[string]float64 {
	result := map[string]float64{"count": float64(len(recs))}
	if len(recs) == 0 {
		return result
	}
	total := 0.0
	strong := 0
	for _, r := range recs {
		total += r.Score
		if r.IsStrong() {
			strong++
		}
	}
	result["average"] = total / float64(len(recs))
	result["strong"] = float64(strong)
	return result
}
func Explain(r model.Recommendation) string {
	if r.Reason == "" {
		return "general availability match"
	}
	return fmt.Sprintf("%s (%.0f%% fit)", r.Reason, r.Score*100)
}
