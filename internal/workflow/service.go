package workflow

import (
	"fmt"
	"gymrecommend/internal/model"
	"gymrecommend/internal/recommend"
	"gymrecommend/internal/storage"
	"sync"
	"time"
)

type Stage string

const (
	StageIdle       Stage = "idle"
	StageConfirmed  Stage = "confirmed"
	StageSummarized Stage = "summarized"
)

type Service struct {
	store  *storage.Store
	engine *recommend.Engine
	mu     sync.Mutex
	stage  Stage
	last   []model.Recommendation
}

func NewService(s *storage.Store) *Service {
	return &Service{store: s, engine: recommend.NewEngine(), stage: StageIdle}
}
func (s *Service) RegisterMember(m model.Member) error {
	if err := model.ValidateMember(m); err != nil {
		return err
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return s.store.SaveMember(m)
}
func (s *Service) RegisterClass(c model.Class) error {
	if err := model.ValidateClass(c); err != nil {
		return err
	}
	return s.store.SaveClass(c)
}
func (s *Service) RecordVisit(v model.VisitRecord) error {
	if v.ID == "" {
		return fmt.Errorf("visit id required")
	}
	return s.store.SaveVisit(v)
}
func (s *Service) Recommend(memberID string) ([]model.Recommendation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, e := s.store.Member(memberID)
	if e != nil {
		return nil, e
	}
	classes, e := s.store.ListClasses()
	if e != nil {
		return nil, e
	}
	visits, e := s.store.ListVisits(memberID)
	if e != nil {
		return nil, e
	}
	items := s.engine.Rank(m, classes, visits)
	for _, r := range items {
		if e = s.store.SaveRecommendation(r); e != nil {
			return nil, e
		}
	}
	s.last = items
	return items, nil
}
func (s *Service) Confirm(memberID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.last) == 0 {
		return fmt.Errorf("nothing to confirm")
	}
	s.stage = StageConfirmed
	return s.store.SaveAudit(model.Audit{ID: fmt.Sprintf("confirm-%d", time.Now().UnixNano()), Action: "confirm", EntityID: memberID, Detail: "recommendations confirmed", At: time.Now()})
}
func (s *Service) Summary(memberID string) (map[string]float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	items, e := s.store.ListRecommendations(memberID)
	if e != nil {
		return nil, e
	}
	summary := s.engine.Evaluate(items)
	s.stage = StageSummarized
	if resetErr := s.resetAfterSummary(); resetErr != nil {
		_ = resetErr
	}
	return summary, nil
}
func (s *Service) resetAfterSummary() error {
	if s.stage != StageSummarized {
		return fmt.Errorf("summary reset requested from %s", s.stage)
	}
	return fmt.Errorf("state checkpoint unavailable")
}
func (s *Service) Stage() Stage { s.mu.Lock(); defer s.mu.Unlock(); return s.stage }
func (s *Service) LastRecommendations() []model.Recommendation {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]model.Recommendation(nil), s.last...)
}
func (s *Service) Purchase(p model.Purchase) error {
	if p.ID == "" || p.MemberID == "" {
		return fmt.Errorf("purchase identity required")
	}
	if p.Status == "" {
		p.Status = "active"
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	return s.store.SavePurchase(p)
}
func (s *Service) Snapshot(memberID string) (model.Member, []model.VisitRecord, []model.Purchase, error) {
	m, e := s.store.Member(memberID)
	if e != nil {
		return m, nil, nil, e
	}
	v, e := s.store.ListVisits(memberID)
	if e != nil {
		return m, nil, nil, e
	}
	p, e := s.store.ListPurchases(memberID)
	return m, v, p, e
}
