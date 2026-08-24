package analytics

import (
	"fmt"
	"gymrecommend/internal/model"
)

type Alert struct {
	Code     string
	Severity string
	Message  string
	MemberID string
	ClassID  string
}

func DetectClassAlerts(classes []model.Class, visits []model.VisitRecord) []Alert {
	counts := map[string]int{}
	for _, v := range visits {
		counts[v.ClassID]++
	}
	out := []Alert{}
	for _, c := range classes {
		used := counts[c.ID]
		if c.Capacity > 0 && used >= c.Capacity {
			out = append(out, Alert{"FULL", "warning", fmt.Sprintf("%s is at capacity", c.Title), "", c.ID})
		}
		if used == 0 {
			out = append(out, Alert{"EMPTY", "info", fmt.Sprintf("%s has no visits", c.Title), "", c.ID})
		}
	}
	return out
}
func DetectMemberAlerts(members []model.Member, visits []model.VisitRecord) []Alert {
	counts := map[string]int{}
	for _, v := range visits {
		counts[v.MemberID]++
	}
	out := []Alert{}
	for _, m := range members {
		if counts[m.ID] == 0 {
			out = append(out, Alert{"NO_VISIT", "warning", m.Name + " has no recorded visits", m.ID, ""})
			continue
		}
		if counts[m.ID] >= 10 {
			out = append(out, Alert{"HIGH_ENGAGEMENT", "info", m.Name + " is highly engaged", m.ID, ""})
		}
	}
	return out
}
func SortAlerts(alerts []Alert) []Alert {
	priority := map[string]int{"warning": 0, "info": 1}
	out := append([]Alert(nil), alerts...)
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if priority[out[j].Severity] < priority[out[i].Severity] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
func AlertSummary(alerts []Alert) map[string]int {
	out := map[string]int{}
	for _, a := range alerts {
		out[a.Severity]++
	}
	return out
}
