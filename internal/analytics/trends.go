package analytics

import (
	"gymrecommend/internal/model"
	"sort"
	"time"
)

type TrendPoint struct {
	Day           time.Time
	Visits        int
	AverageRating float64
}

func DailyTrend(visits []model.VisitRecord) []TrendPoint {
	days := map[string]*TrendPoint{}
	for _, v := range visits {
		day := v.AttendedAt.UTC().Truncate(24 * time.Hour)
		key := day.Format("2006-01-02")
		p := days[key]
		if p == nil {
			p = &TrendPoint{Day: day}
			days[key] = p
		}
		p.Visits++
		if v.Rating > 0 {
			p.AverageRating += float64(v.Rating)
		}
	}
	out := make([]TrendPoint, 0, len(days))
	for _, p := range days {
		if p.Visits > 0 {
			p.AverageRating /= float64(p.Visits)
		}
		out = append(out, *p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Day.Before(out[j].Day) })
	return out
}
func PeakDay(points []TrendPoint) TrendPoint {
	var best TrendPoint
	for _, p := range points {
		if p.Visits > best.Visits {
			best = p
		}
	}
	return best
}
func Retention(members []model.Member, visits []model.VisitRecord, window time.Duration) float64 {
	if len(members) == 0 {
		return 0
	}
	cut := time.Now().Add(-window)
	seen := map[string]bool{}
	for _, v := range visits {
		if v.AttendedAt.After(cut) {
			seen[v.MemberID] = true
		}
	}
	return float64(len(seen)) / float64(len(members))
}
func GoalDistribution(members []model.Member) map[string]int {
	out := map[string]int{}
	for _, m := range members {
		for _, g := range m.Goals {
			out[g]++
		}
	}
	return out
}
func SlotDistribution(members []model.Member) map[string]int {
	out := map[string]int{}
	for _, m := range members {
		if m.PreferredSlot != "" {
			out[m.PreferredSlot]++
		}
	}
	return out
}
func RatingHistogram(visits []model.VisitRecord) map[int]int {
	out := map[int]int{}
	for _, v := range visits {
		if v.Rating >= 1 && v.Rating <= 5 {
			out[v.Rating]++
		}
	}
	return out
}
