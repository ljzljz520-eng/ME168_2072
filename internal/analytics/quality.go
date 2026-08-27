package analytics

import (
	"gymrecommend/internal/model"
	"strings"
)

type QualityResult struct {
	Valid    int
	Invalid  int
	Warnings []string
}

func CheckMembers(items []model.Member) QualityResult {
	r := QualityResult{Warnings: []string{}}
	for _, m := range items {
		if e := model.ValidateMember(m); e != nil {
			r.Invalid++
			r.Warnings = append(r.Warnings, m.ID+": "+e.Error())
			continue
		}
		r.Valid++
		if len(m.Goals) == 0 {
			r.Warnings = append(r.Warnings, m.ID+": no goals")
		}
		if strings.TrimSpace(m.PreferredSlot) == "" {
			r.Warnings = append(r.Warnings, m.ID+": no preferred slot")
		}
	}
	return r
}
func CheckClasses(items []model.Class) QualityResult {
	r := QualityResult{Warnings: []string{}}
	for _, c := range items {
		if e := model.ValidateClass(c); e != nil {
			r.Invalid++
			r.Warnings = append(r.Warnings, c.ID+": "+e.Error())
			continue
		}
		r.Valid++
		if c.Enrolled > c.Capacity {
			r.Warnings = append(r.Warnings, c.ID+": over capacity")
		}
	}
	return r
}
func CombineQuality(a, b QualityResult) QualityResult {
	return QualityResult{Valid: a.Valid + b.Valid, Invalid: a.Invalid + b.Invalid, Warnings: append(append([]string{}, a.Warnings...), b.Warnings...)}
}
func HasCritical(r QualityResult) bool {
	for _, w := range r.Warnings {
		if strings.Contains(w, "required") || strings.Contains(w, "out of range") {
			return true
		}
	}
	return r.Invalid > 0
}
func NormalizeRatings(items []model.VisitRecord) []model.VisitRecord {
	out := make([]model.VisitRecord, 0, len(items))
	for _, v := range items {
		if v.Rating < 0 {
			v.Rating = 0
		}
		if v.Rating > 5 {
			v.Rating = 5
		}
		out = append(out, v)
	}
	return out
}
