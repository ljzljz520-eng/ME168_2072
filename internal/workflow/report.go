package workflow

import (
	"fmt"
	"gymrecommend/internal/model"
	"sort"
	"strings"
)

type Report struct {
	MemberID  string
	Total     int
	Average   float64
	Top       model.Recommendation
	Narrative string
}

func BuildReport(member string, items []model.Recommendation) Report {
	r := Report{MemberID: member, Total: len(items)}
	if len(items) == 0 {
		r.Narrative = "No matching classes"
		return r
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Score > items[j].Score })
	r.Top = items[0]
	sum := 0.0
	for _, x := range items {
		sum += x.Score
	}
	r.Average = sum / float64(len(items))
	r.Narrative = fmt.Sprintf("Top choice %s: %s", r.Top.ClassID, strings.TrimSpace(r.Top.Reason))
	return r
}
func RenderReport(r Report) string {
	return fmt.Sprintf("member=%s total=%d average=%.2f top=%s narrative=%s", r.MemberID, r.Total, r.Average, r.Top.ClassID, r.Narrative)
}
