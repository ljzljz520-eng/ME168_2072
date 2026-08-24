package analytics

import (
	"math"
	"sort"
	"time"
)

type Forecast struct {
	NextWeek   float64
	Lower      float64
	Upper      float64
	Confidence string
}

func ForecastVisits(points []TrendPoint) Forecast {
	if len(points) == 0 {
		return Forecast{Confidence: "none"}
	}
	values := make([]float64, len(points))
	for i, p := range points {
		values[i] = float64(p.Visits)
	}
	avg := mean(values)
	slope := trend(values)
	next := math.Max(0, avg+slope*7)
	spread := deviation(values)
	confidence := "medium"
	if spread < 1 {
		confidence = "high"
	}
	if spread > avg && avg > 0 {
		confidence = "low"
	}
	return Forecast{NextWeek: next, Lower: math.Max(0, next-spread*2), Upper: next + spread*2, Confidence: confidence}
}
func mean(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}
func trend(xs []float64) float64 {
	if len(xs) < 2 {
		return 0
	}
	sum := 0.0
	for i := 1; i < len(xs); i++ {
		sum += xs[i] - xs[i-1]
	}
	return sum / float64(len(xs)-1)
}
func deviation(xs []float64) float64 {
	if len(xs) < 2 {
		return 0
	}
	m := mean(xs)
	sum := 0.0
	for _, x := range xs {
		d := x - m
		sum += d * d
	}
	return math.Sqrt(sum / float64(len(xs)-1))
}
func ForecastBySlot(points map[string][]TrendPoint) map[string]Forecast {
	out := map[string]Forecast{}
	for slot, rows := range points {
		out[slot] = ForecastVisits(rows)
	}
	return out
}
func FillMissingDays(points []TrendPoint, start, end time.Time) []TrendPoint {
	by := map[string]TrendPoint{}
	for _, p := range points {
		by[p.Day.Format("2006-01-02")] = p
	}
	out := []TrendPoint{}
	for day := start.Truncate(24 * time.Hour); !day.After(end); day = day.Add(24 * time.Hour) {
		if p, ok := by[day.Format("2006-01-02")]; ok {
			out = append(out, p)
		} else {
			out = append(out, TrendPoint{Day: day})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Day.Before(out[j].Day) })
	return out
}
