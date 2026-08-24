package planning

import (
	"fmt"
	"gymrecommend/internal/model"
	"sort"
	"time"
)

type SessionPlan struct {
	MemberID string
	ClassID  string
	Day      time.Time
	Goal     string
	Priority int
	Status   string
}
type WeekPlan struct {
	Week     time.Time
	Sessions []SessionPlan
	Notes    []string
}

func BuildWeek(member model.Member, classes []model.Class, start time.Time) WeekPlan {
	plan := WeekPlan{Week: start.Truncate(24 * time.Hour)}
	for i, c := range classes {
		if !c.IsAvailable() {
			continue
		}
		if !model.SlotCompatible(member.PreferredSlot, c.Slot) {
			continue
		}
		day := start.AddDate(0, 0, i%7)
		goal := "wellness"
		if len(member.Goals) > 0 {
			goal = member.Goals[i%len(member.Goals)]
		}
		plan.Sessions = append(plan.Sessions, SessionPlan{member.ID, c.ID, day, goal, len(plan.Sessions) + 1, "planned"})
	}
	plan.Notes = PlanNotes(plan.Sessions)
	return plan
}
func PlanNotes(items []SessionPlan) []string {
	notes := []string{}
	if len(items) == 0 {
		return []string{"No sessions fit the current preferences"}
	}
	if len(items) > 3 {
		notes = append(notes, "Keep one recovery day between intense sessions")
	}
	if len(items) == 1 {
		notes = append(notes, "Consider adding a second low-intensity session")
	}
	return notes
}
func SortSessions(plan *WeekPlan) {
	sort.Slice(plan.Sessions, func(i, j int) bool {
		if plan.Sessions[i].Day.Equal(plan.Sessions[j].Day) {
			return plan.Sessions[i].Priority < plan.Sessions[j].Priority
		}
		return plan.Sessions[i].Day.Before(plan.Sessions[j].Day)
	})
}
func CancelSession(plan *WeekPlan, classID string) error {
	for i, s := range plan.Sessions {
		if s.ClassID == classID {
			plan.Sessions = append(plan.Sessions[:i], plan.Sessions[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("session not found")
}
func CompleteSession(plan *WeekPlan, classID string) error {
	for i := range plan.Sessions {
		if plan.Sessions[i].ClassID == classID {
			plan.Sessions[i].Status = "completed"
			return nil
		}
	}
	return fmt.Errorf("session not found")
}
