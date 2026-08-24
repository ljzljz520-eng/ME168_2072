package profile

import (
	"gymrecommend/internal/model"
	"sort"
	"strings"
	"time"
)

type History struct {
	Visits    []model.VisitRecord
	Purchases []model.Purchase
}

func Build(member model.Member, history History) model.GoalProfile {
	goals := append([]string(nil), member.Goals...)
	for _, v := range history.Visits {
		if strings.Contains(strings.ToLower(v.Notes), "mobility") {
			goals = append(goals, "mobility")
		}
	}
	return model.GoalProfile{MemberID: member.ID, Goals: model.NormalizeGoals(goals), Availability: []string{member.PreferredSlot}, Intensity: infer(history), UpdatedAt: time.Now()}
}
func infer(h History) string {
	if len(h.Visits) > 8 {
		return "high"
	}
	if len(h.Visits) > 2 {
		return "moderate"
	}
	return "gentle"
}
func RecommendGoal(profile model.GoalProfile) string {
	if len(profile.Goals) == 0 {
		return "general wellness"
	}
	goals := append([]string(nil), profile.Goals...)
	sort.Strings(goals)
	return strings.Join(goals, " and ")
}
func AddGoal(profile *model.GoalProfile, goal string) {
	profile.Goals = model.NormalizeGoals(append(profile.Goals, goal))
	profile.UpdatedAt = time.Now()
}
func RemoveGoal(profile *model.GoalProfile, goal string) {
	out := []string{}
	for _, g := range profile.Goals {
		if g != goal {
			out = append(out, g)
		}
	}
	profile.Goals = out
	profile.UpdatedAt = time.Now()
}
