package analytics

import (
	"gymrecommend/internal/model"
	"sort"
)

type Benchmark struct {
	Name   string
	Value  float64
	Target float64
	Status string
}

func BuildBenchmarks(d Dashboard) []Benchmark {
	items := []Benchmark{{"member coverage", float64(d.Members), 100, ""}, {"recommendation quality", d.AverageScore, .75, ""}}
	for i := range items {
		items[i].Status = BenchmarkStatus(items[i].Value, items[i].Target)
	}
	return items
}
func BenchmarkStatus(value, target float64) string {
	if value >= target {
		return "on-target"
	}
	if value >= target*.8 {
		return "watch"
	}
	return "below-target"
}
func SortBenchmarks(items []Benchmark) []Benchmark {
	out := append([]Benchmark(nil), items...)
	sort.Slice(out, func(i, j int) bool { return out[i].Value/out[i].Target < out[j].Value/out[j].Target })
	return out
}
func Gap(item Benchmark) float64 {
	if item.Value >= item.Target {
		return 0
	}
	return item.Target - item.Value
}
func OverallStatus(items []Benchmark) string {
	for _, i := range items {
		if i.Status == "below-target" {
			return "action-needed"
		}
	}
	for _, i := range items {
		if i.Status == "watch" {
			return "watch"
		}
	}
	return "healthy"
}
func ClassBenchmark(classes []model.Class, visits []model.VisitRecord) []Benchmark {
	metrics := ClassMetrics(classes, visits)
	out := []Benchmark{}
	for _, m := range metrics {
		out = append(out, Benchmark{m.ClassID, m.Utilization, .7, BenchmarkStatus(m.Utilization, .7)})
	}
	return out
}
