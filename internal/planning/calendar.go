package planning

import (
	"sort"
	"time"
)

type CalendarDay struct {
	Day       time.Time
	Planned   int
	Completed int
	Cancelled int
}

func Calendar(plan WeekPlan) []CalendarDay {
	by := map[string]*CalendarDay{}
	for _, s := range plan.Sessions {
		key := s.Day.Format("2006-01-02")
		d := by[key]
		if d == nil {
			d = &CalendarDay{Day: s.Day}
			by[key] = d
		}
		switch s.Status {
		case "completed":
			d.Completed++
		case "cancelled":
			d.Cancelled++
		default:
			d.Planned++
		}
	}
	out := []CalendarDay{}
	for _, d := range by {
		out = append(out, *d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Day.Before(out[j].Day) })
	return out
}
func MarkCancelled(plan *WeekPlan, classID string) bool {
	for i := range plan.Sessions {
		if plan.Sessions[i].ClassID == classID {
			plan.Sessions[i].Status = "cancelled"
			return true
		}
	}
	return false
}
func Upcoming(plan WeekPlan, now time.Time) []SessionPlan {
	out := []SessionPlan{}
	for _, s := range plan.Sessions {
		if s.Day.After(now) && s.Status == "planned" {
			out = append(out, s)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Day.Before(out[j].Day) })
	return out
}
func WeekBounds(day time.Time) (time.Time, time.Time) {
	start := day.AddDate(0, 0, -int(day.Weekday()))
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	return start, start.AddDate(0, 0, 7)
}
