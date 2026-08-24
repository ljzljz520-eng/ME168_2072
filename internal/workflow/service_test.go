package workflow

import (
	"gymrecommend/internal/model"
	"gymrecommend/internal/storage"
	"path/filepath"
	"testing"
)

func setup(t *testing.T) *Service {
	s, e := storage.Open(filepath.Join(t.TempDir(), "x.db"))
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { s.Close() })
	svc := NewService(s)
	_ = svc.RegisterMember(model.Member{ID: "m", Name: "Ada", Goals: []string{"strength"}, PreferredSlot: "morning", CoachRating: 4.8})
	_ = svc.RegisterClass(model.Class{ID: "c", Title: "Strength Circuit", Kind: "group", Slot: "morning", Capacity: 10})
	return svc
}
func TestWorkflowAccept(t *testing.T) {
	s := setup(t)
	xs, e := s.Recommend("m")
	if e != nil || len(xs) != 1 {
		t.Fatal(xs, e)
	}
	if e = s.Confirm("m"); e != nil {
		t.Fatal(e)
	}
}

func TestWorkflowPublish(t *testing.T) {
	s := setup(t)
	items, e := s.Recommend("m")
	if e != nil || len(items) == 0 {
		t.Fatal(e)
	}
	q := NewQueue()
	if e = q.Enqueue(items[0]); e != nil {
		t.Fatal(e)
	}
	if e = q.Drain(func(model.Recommendation) error { return nil }); e != nil {
		t.Fatal(e)
	}
}

func TestWorkflowReopen(t *testing.T) {
	s := setup(t)
	if _, _, _, e := s.Snapshot("m"); e != nil {
		t.Fatal(e)
	}
}
