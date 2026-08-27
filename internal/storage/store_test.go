package storage

import (
	"gymrecommend/internal/model"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	p := filepath.Join(t.TempDir(), "gym.db")
	s, e := Open(p)
	if e != nil {
		t.Fatal(e)
	}
	m := model.Member{ID: "m1", Name: "Ada", CoachRating: 4}
	if e = s.SaveMember(m); e != nil {
		t.Fatal(e)
	}
	s.Close()
	s, e = Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	got, e := s.Member("m1")
	if e != nil || got.Name != "Ada" {
		t.Fatalf("reopen got %+v %v", got, e)
	}
}
