package analytics

import (
	"gymrecommend/internal/model"
	"strings"
)

type Segment struct {
	Name      string
	Members   []string
	Rationale string
}

func SegmentMembers(members []model.Member, visits []model.VisitRecord) []Segment {
	metrics := MemberMetrics(members, visits)
	groups := map[string][]string{}
	for _, m := range metrics {
		groups[EngagementBand(m)] = append(groups[EngagementBand(m)], m.MemberID)
	}
	names := []string{"new", "growing", "regular", "champion"}
	out := []Segment{}
	for _, name := range names {
		if ids := groups[name]; len(ids) > 0 {
			out = append(out, Segment{name, ids, segmentReason(name)})
		}
	}
	return out
}
func segmentReason(name string) string {
	switch name {
	case "new":
		return "welcome and orientation"
	case "growing":
		return "encourage a second weekly visit"
	case "regular":
		return "offer progression"
	default:
		return "recognize loyalty"
	}
}
func MatchSegment(m model.Member, s Segment) bool {
	for _, id := range s.Members {
		if id == m.ID {
			return true
		}
	}
	return false
}
func SearchSegments(segments []Segment, query string) []Segment {
	query = strings.ToLower(query)
	out := []Segment{}
	for _, s := range segments {
		if strings.Contains(strings.ToLower(s.Name), query) || strings.Contains(strings.ToLower(s.Rationale), query) {
			out = append(out, s)
		}
	}
	return out
}
func SegmentCounts(segments []Segment) map[string]int {
	out := map[string]int{}
	for _, s := range segments {
		out[s.Name] = len(s.Members)
	}
	return out
}
