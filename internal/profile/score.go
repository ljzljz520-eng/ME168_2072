package profile

import (
	"gymrecommend/internal/model"
	"math"
)

func Consistency(h History) float64 {
	if len(h.Visits) < 2 {
		return float64(len(h.Visits))
	}
	days := map[string]bool{}
	for _, v := range h.Visits {
		days[v.AttendedAt.Format("2006-01-02")] = true
	}
	return math.Min(1, float64(len(days))/float64(len(h.Visits)))
}
func Satisfaction(h History) float64 {
	sum, n := 0, 0
	for _, v := range h.Visits {
		if v.Rating > 0 {
			sum += v.Rating
			n++
		}
	}
	if n == 0 {
		return 0
	}
	return float64(sum) / (float64(n) * 5)
}
func PackageFit(m model.Member, h History) string {
	if len(h.Visits) >= 8 {
		return "unlimited"
	}
	if m.CoachRating >= 4.5 {
		return "private-10"
	}
	return "group-10"
}
func RiskFlags(m model.Member, h History) []string {
	flags := []string{}
	if len(h.Visits) == 0 {
		flags = append(flags, "inactive")
	}
	if Satisfaction(h) < .5 && len(h.Visits) > 0 {
		flags = append(flags, "low-satisfaction")
	}
	if m.Package == "" {
		flags = append(flags, "no-package")
	}
	return flags
}
