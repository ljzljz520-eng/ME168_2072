package planning

import (
	"fmt"
	"time"
)

type Recurrence struct {
	ClassID  string
	MemberID string
	Weekdays []time.Weekday
	Weeks    int
	Start    time.Time
}

func ExpandRecurrence(r Recurrence) []SessionPlan {
	out := []SessionPlan{}
	if r.Weeks < 1 {
		return out
	}
	if len(r.Weekdays) == 0 {
		return out
	}
	for week := 0; week < r.Weeks; week++ {
		for _, weekday := range r.Weekdays {
			delta := (int(weekday) - int(r.Start.Weekday()) + 7) % 7
			day := r.Start.AddDate(0, 0, week*7+delta)
			out = append(out, SessionPlan{MemberID: r.MemberID, ClassID: r.ClassID, Day: day, Goal: "recurring", Priority: week + 1, Status: "planned"})
		}
	}
	return out
}
func ValidateRecurrence(r Recurrence) error {
	if r.ClassID == "" || r.MemberID == "" {
		return fmt.Errorf("recurrence identity required")
	}
	if r.Weeks < 1 || r.Weeks > 52 {
		return fmt.Errorf("weeks out of range")
	}
	seen := map[time.Weekday]bool{}
	for _, d := range r.Weekdays {
		if seen[d] {
			return fmt.Errorf("duplicate weekday")
		}
		seen[d] = true
	}
	return nil
}
func MergePlans(a, b WeekPlan) WeekPlan {
	out := a
	out.Sessions = append(append([]SessionPlan{}, a.Sessions...), b.Sessions...)
	out.Notes = append(append([]string{}, a.Notes...), b.Notes...)
	SortSessions(&out)
	return out
}
func Deduplicate(plan *WeekPlan) int {
	seen := map[string]bool{}
	out := []SessionPlan{}
	removed := 0
	for _, s := range plan.Sessions {
		key := s.ClassID + "|" + s.MemberID + "|" + s.Day.Format("2006-01-02")
		if seen[key] {
			removed++
			continue
		}
		seen[key] = true
		out = append(out, s)
	}
	plan.Sessions = out
	return removed
}
func NextSession(plan WeekPlan, now time.Time) (SessionPlan, bool) {
	var best SessionPlan
	found := false
	for _, s := range plan.Sessions {
		if s.Status == "planned" && s.Day.After(now) && (!found || s.Day.Before(best.Day)) {
			best = s
			found = true
		}
	}
	return best, found
}
func SessionSpan(plan WeekPlan) time.Duration {
	if len(plan.Sessions) < 2 {
		return 0
	}
	first, last := plan.Sessions[0].Day, plan.Sessions[0].Day
	for _, s := range plan.Sessions {
		if s.Day.Before(first) {
			first = s.Day
		}
		if s.Day.After(last) {
			last = s.Day
		}
	}
	return last.Sub(first)
}

func SessionsForGoal(plan WeekPlan, goal string) []SessionPlan {
	out := []SessionPlan{}
	for _, session := range plan.Sessions {
		if goal == "" || session.Goal == goal {
			out = append(out, session)
		}
	}
	return out
}

func MoveGoal(plan *WeekPlan, from, to string) int {
	changed := 0
	for i := range plan.Sessions {
		if plan.Sessions[i].Goal == from {
			plan.Sessions[i].Goal = to
			changed++
		}
	}
	return changed
}

func PriorityScore(session SessionPlan) int {
	score := session.Priority
	if session.Status == "planned" {
		score += 2
	}
	if session.Goal == "strength" {
		score += 3
	}
	if session.Goal == "recovery" {
		score -= 1
	}
	return score
}

func HighestPriority(plan WeekPlan) (SessionPlan, bool) {
	if len(plan.Sessions) == 0 {
		return SessionPlan{}, false
	}
	best := plan.Sessions[0]
	for _, session := range plan.Sessions[1:] {
		if PriorityScore(session) > PriorityScore(best) {
			best = session
		}
	}
	return best, true
}
