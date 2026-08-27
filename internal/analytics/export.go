package analytics

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func WriteDashboardJSON(w io.Writer, d Dashboard) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(d)
}
func WriteAttendanceCSV(w io.Writer, items []AttendanceMetric) error {
	cw := csv.NewWriter(w)
	if e := cw.Write([]string{"member_id", "visits", "rated", "average_rating", "last_visit"}); e != nil {
		return e
	}
	for _, m := range items {
		row := []string{m.MemberID, strconv.Itoa(m.Visits), strconv.Itoa(m.Rated), strconv.FormatFloat(m.AverageRating, 'f', 2, 64), m.LastVisit.Format("2006-01-02")}
		if e := cw.Write(row); e != nil {
			return e
		}
	}
	cw.Flush()
	return cw.Error()
}
func WriteClassCSV(w io.Writer, items []ClassMetric) error {
	cw := csv.NewWriter(w)
	if e := cw.Write([]string{"class_id", "attendance", "capacity", "utilization", "rating"}); e != nil {
		return e
	}
	for _, m := range items {
		if e := cw.Write([]string{m.ClassID, strconv.Itoa(m.Attendance), strconv.Itoa(m.Capacity), strconv.FormatFloat(m.Utilization, 'f', 2, 64), strconv.FormatFloat(m.Rating, 'f', 2, 64)}); e != nil {
			return e
		}
	}
	cw.Flush()
	return cw.Error()
}
func FormatDashboard(d Dashboard) string {
	return fmt.Sprintf("members=%d classes=%d recommendations=%d average=%.2f", d.Members, d.Classes, d.Recommendations, d.AverageScore)
}
func TopMembers(items []AttendanceMetric, n int) []AttendanceMetric {
	if n < 0 {
		n = 0
	}
	if n > len(items) {
		n = len(items)
	}
	out := append([]AttendanceMetric(nil), items...)
	return out[:n]
}
