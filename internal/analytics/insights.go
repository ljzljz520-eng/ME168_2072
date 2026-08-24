package analytics

import (
	"fmt"
	"gymrecommend/internal/model"
	"sort"
	"strings"
)

type Insight struct {
	Kind   string
	Title  string
	Detail string
	Score  float64
}

func GenerateInsights(members []model.Member, classes []model.Class, recs []model.Recommendation, visits []model.VisitRecord) []Insight {
	out := []Insight{}
	d := BuildDashboard(members, classes, recs, visits)
	if d.AverageScore < .5 {
		out = append(out, Insight{"quality", "Recommendation quality needs attention", fmt.Sprintf("Average score is %.2f", d.AverageScore), 1 - d.AverageScore})
	} else {
		out = append(out, Insight{"quality", "Recommendations are healthy", fmt.Sprintf("Average score is %.2f", d.AverageScore), d.AverageScore})
	}
	segments := SegmentMembers(members, visits)
	for _, s := range segments {
		if s.Name == "new" {
			out = append(out, Insight{"engagement", "Welcome new members", fmt.Sprintf("%d members need onboarding", len(s.Members)), float64(len(s.Members))})
		}
		if s.Name == "champion" {
			out = append(out, Insight{"engagement", "Recognize loyal members", fmt.Sprintf("%d members are champions", len(s.Members)), float64(len(s.Members))})
		}
	}
	alerts := SortAlerts(append(DetectClassAlerts(classes, visits), DetectMemberAlerts(members, visits)...))
	for _, a := range alerts {
		if a.Severity == "warning" {
			out = append(out, Insight{"alert", a.Code, a.Message, 1})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}
func SummarizeInsights(items []Insight) string {
	if len(items) == 0 {
		return "No insights"
	}
	parts := make([]string, 0, len(items))
	for _, i := range items {
		parts = append(parts, i.Title)
	}
	return strings.Join(parts, "; ")
}
func FilterInsights(items []Insight, kind string) []Insight {
	out := []Insight{}
	for _, i := range items {
		if kind == "" || i.Kind == kind {
			out = append(out, i)
		}
	}
	return out
}
func HighestInsight(items []Insight) Insight {
	if len(items) == 0 {
		return Insight{}
	}
	best := items[0]
	for _, i := range items[1:] {
		if i.Score > best.Score {
			best = i
		}
	}
	return best
}
func InsightCount(items []Insight) map[string]int {
	out := map[string]int{}
	for _, i := range items {
		out[i.Kind]++
	}
	return out
}

func InsightsForMember(items []Insight, memberID string) []Insight {
	out := []Insight{}
	for _, item := range items {
		if memberID == "" || strings.Contains(item.Detail, memberID) {
			out = append(out, item)
		}
	}
	return out
}

func ScoreInsight(item Insight, weights map[string]float64) float64 {
	weight := weights[item.Kind]
	if weight == 0 {
		weight = 1
	}
	if item.Score < 0 {
		return 0
	}
	return item.Score * weight
}

func ReScore(items []Insight, weights map[string]float64) []Insight {
	out := append([]Insight(nil), items...)
	for i := range out {
		out[i].Score = ScoreInsight(out[i], weights)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}

func Actionable(items []Insight) []Insight {
	out := []Insight{}
	for _, item := range items {
		if item.Kind == "alert" || item.Kind == "quality" {
			out = append(out, item)
		}
	}
	return out
}

func DeduplicateInsights(items []Insight) []Insight {
	seen := map[string]bool{}
	out := []Insight{}
	for _, item := range items {
		key := item.Kind + ":" + item.Title
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, item)
	}
	return out
}

func InsightKinds(items []Insight) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, item := range items {
		if !seen[item.Kind] {
			seen[item.Kind] = true
			out = append(out, item.Kind)
		}
	}
	sort.Strings(out)
	return out
}

func InsightNeedsReview(item Insight) bool {
	if item.Kind == "alert" {
		return true
	}
	if item.Kind == "quality" && item.Score > .5 {
		return true
	}
	return false
}
