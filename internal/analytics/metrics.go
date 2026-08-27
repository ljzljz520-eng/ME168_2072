package analytics

import (
	"gymrecommend/internal/model"
	"math"
	"sort"
	"time"
)

type AttendanceMetric struct {
	MemberID      string
	Visits        int
	Rated         int
	AverageRating float64
	LastVisit     time.Time
}
type ClassMetric struct {
	ClassID     string
	Attendance  int
	Capacity    int
	Utilization float64
	Rating      float64
}
type Dashboard struct {
	Members         int
	Classes         int
	Recommendations int
	AverageScore    float64
	Attendance      []AttendanceMetric
	ClassesByUse    []ClassMetric
}

func MemberMetrics(members []model.Member, visits []model.VisitRecord) []AttendanceMetric {
	result := make([]AttendanceMetric, 0, len(members))
	by := map[string][]model.VisitRecord{}
	for _, v := range visits {
		by[v.MemberID] = append(by[v.MemberID], v)
	}
	for _, m := range members {
		metric := AttendanceMetric{MemberID: m.ID}
		rows := by[m.ID]
		for _, v := range rows {
			metric.Visits++
			if v.Rating > 0 {
				metric.Rated++
				metric.AverageRating += float64(v.Rating)
			}
			if metric.LastVisit.IsZero() || v.AttendedAt.After(metric.LastVisit) {
				metric.LastVisit = v.AttendedAt
			}
		}
		if metric.Rated > 0 {
			metric.AverageRating /= float64(metric.Rated)
		}
		result = append(result, metric)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Visits > result[j].Visits })
	return result
}
func ClassMetrics(classes []model.Class, visits []model.VisitRecord) []ClassMetric {
	result := make([]ClassMetric, 0, len(classes))
	by := map[string][]model.VisitRecord{}
	for _, v := range visits {
		by[v.ClassID] = append(by[v.ClassID], v)
	}
	for _, c := range classes {
		metric := ClassMetric{ClassID: c.ID, Capacity: c.Capacity}
		rows := by[c.ID]
		for _, v := range rows {
			metric.Attendance++
			if v.Rating > 0 {
				metric.Rating += float64(v.Rating)
			}
		}
		if metric.Rating > 0 {
			metric.Rating /= float64(metric.Attendance)
		}
		if c.Capacity > 0 {
			metric.Utilization = math.Min(1, float64(metric.Attendance)/float64(c.Capacity))
		}
		result = append(result, metric)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Utilization > result[j].Utilization })
	return result
}
func BuildDashboard(members []model.Member, classes []model.Class, recs []model.Recommendation, visits []model.VisitRecord) Dashboard {
	d := Dashboard{Members: len(members), Classes: len(classes), Recommendations: len(recs)}
	for _, r := range recs {
		d.AverageScore += r.Score
	}
	if len(recs) > 0 {
		d.AverageScore /= float64(len(recs))
	}
	d.Attendance = MemberMetrics(members, visits)
	d.ClassesByUse = ClassMetrics(classes, visits)
	return d
}
func ActiveMemberCount(metrics []AttendanceMetric, since time.Time) int {
	count := 0
	for _, m := range metrics {
		if !m.LastVisit.IsZero() && m.LastVisit.After(since) {
			count++
		}
	}
	return count
}
func EngagementBand(m AttendanceMetric) string {
	if m.Visits == 0 {
		return "new"
	}
	if m.Visits < 3 {
		return "growing"
	}
	if m.Visits < 8 {
		return "regular"
	}
	return "champion"
}
