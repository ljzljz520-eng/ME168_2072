package model

import "time"

type Member struct {
	ID, Name      string
	Goals         []string
	PreferredSlot string
	CoachRating   float64
	Package       string
	CreatedAt     time.Time
}
type Class struct {
	ID, Title, Kind, Slot, Coach, Level string
	Capacity                            int
	Enrolled                            int
}
type Recommendation struct {
	ID, MemberID, ClassID, Kind, Reason string
	Score                               float64
	CreatedAt                           time.Time
}
type Purchase struct {
	ID, MemberID, Package, Status string
	Sessions                      int
	CreatedAt                     time.Time
}
type VisitRecord struct {
	ID, MemberID, ClassID string
	AttendedAt            time.Time
	Rating                int
	Notes                 string
}
type GoalProfile struct {
	MemberID     string
	Goals        []string
	Availability []string
	Intensity    string
	UpdatedAt    time.Time
}
type Audit struct {
	ID, Action, EntityID, Detail string
	At                           time.Time
}

func (m Member) HasGoal(goal string) bool {
	for _, g := range m.Goals {
		if g == goal {
			return true
		}
	}
	return false
}
func (c Class) IsAvailable() bool       { return c.Enrolled < c.Capacity }
func (r Recommendation) IsStrong() bool { return r.Score >= 0.75 }
func (p Purchase) Remaining() int {
	if p.Sessions < 0 {
		return 0
	}
	return p.Sessions
}
