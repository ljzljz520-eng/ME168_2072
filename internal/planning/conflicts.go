package planning

import (
	"time"
)

func FindConflicts(items []SessionPlan) map[string][]string {
	out := map[string][]string{}
	for i, a := range items {
		for j := i + 1; j < len(items); j++ {
			b := items[j]
			if sameDay(a.Day, b.Day) {
				out[a.ClassID] = append(out[a.ClassID], b.ClassID)
				out[b.ClassID] = append(out[b.ClassID], a.ClassID)
			}
		}
	}
	return out
}
func sameDay(a, b time.Time) bool { return a.Format("2006-01-02") == b.Format("2006-01-02") }
func ShiftConflicts(plan *WeekPlan) int {
	conf := FindConflicts(plan.Sessions)
	moved := 0
	for i := range plan.Sessions {
		if len(conf[plan.Sessions[i].ClassID]) > 0 {
			plan.Sessions[i].Day = plan.Sessions[i].Day.AddDate(0, 0, 1)
			moved++
		}
	}
	return moved
}
func IsBalanced(plan WeekPlan) bool {
	if len(plan.Sessions) < 2 {
		return true
	}
	counts := map[string]int{}
	for _, s := range plan.Sessions {
		counts[s.Goal]++
	}
	return len(counts) > 1
}
func GoalCounts(plan WeekPlan) map[string]int {
	out := map[string]int{}
	for _, s := range plan.Sessions {
		out[s.Goal]++
	}
	return out
}
