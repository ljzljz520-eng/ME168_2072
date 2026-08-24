package analytics

import (
	"gymrecommend/internal/model"
	"sort"
)

type Comparison struct {
	ClassID   string
	Current   float64
	Previous  float64
	Change    float64
	Direction string
}

func CompareUtilization(current, previous []ClassMetric) []Comparison {
	old := map[string]float64{}
	for _, p := range previous {
		old[p.ClassID] = p.Utilization
	}
	out := []Comparison{}
	for _, c := range current {
		p := old[c.ClassID]
		change := c.Utilization - p
		dir := "flat"
		if change > .01 {
			dir = "up"
		} else if change < -.01 {
			dir = "down"
		}
		out = append(out, Comparison{c.ClassID, c.Utilization, p, change, dir})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Change > out[j].Change })
	return out
}
func CompareScores(current, previous []model.Recommendation) float64 {
	if len(current) == 0 || len(previous) == 0 {
		return 0
	}
	avg := func(xs []model.Recommendation) float64 {
		sum := 0.0
		for _, x := range xs {
			sum += x.Score
		}
		return sum / float64(len(xs))
	}
	return avg(current) - avg(previous)
}
func ChangeLabel(change float64) string {
	if change > .1 {
		return "material improvement"
	}
	if change < -.1 {
		return "material decline"
	}
	return "stable"
}
func RankChanges(items []Comparison, limit int) []Comparison {
	if limit < 0 {
		limit = 0
	}
	if limit > len(items) {
		limit = len(items)
	}
	out := append([]Comparison(nil), items...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Change == out[j].Change {
			return out[i].ClassID < out[j].ClassID
		}
		return out[i].Change > out[j].Change
	})
	return out[:limit]
}
