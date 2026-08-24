package planning

import (
	"fmt"
	"strings"
	"time"
)

func FormatPlan(plan WeekPlan) string {
	lines := []string{fmt.Sprintf("week of %s", plan.Week.Format("2006-01-02"))}
	for _, s := range plan.Sessions {
		lines = append(lines, fmt.Sprintf("%s %s %s %s", s.Day.Format("Mon 2006-01-02"), s.ClassID, s.Goal, s.Status))
	}
	lines = append(lines, plan.Notes...)
	return strings.Join(lines, "\n")
}
func FormatDay(day CalendarDay) string {
	return fmt.Sprintf("%s planned=%d completed=%d cancelled=%d", day.Day.Format("2006-01-02"), day.Planned, day.Completed, day.Cancelled)
}
func ParseDay(value string) (time.Time, error) { return time.Parse("2006-01-02", value) }
func DateRange(start, end time.Time) []time.Time {
	out := []time.Time{}
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		out = append(out, d)
	}
	return out
}
func CountStatus(plan WeekPlan, status string) int {
	count := 0
	for _, s := range plan.Sessions {
		if s.Status == status {
			count++
		}
	}
	return count
}
func Reschedule(plan *WeekPlan, classID string, day time.Time) bool {
	for i := range plan.Sessions {
		if plan.Sessions[i].ClassID == classID {
			plan.Sessions[i].Day = day
			plan.Sessions[i].Status = "planned"
			return true
		}
	}
	return false
}
