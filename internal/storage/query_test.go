package storage

import (
	"gymrecommend/internal/model"
	"path/filepath"
	"testing"
)

func TestListClasses(t *testing.T) {
	s, e := Open(filepath.Join(t.TempDir(), "x.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	_ = s.SaveClass(model.Class{ID: "c", Title: "Flow", Kind: "group", Capacity: 2})
	xs, e := s.ListClasses()
	if e != nil || len(xs) != 1 {
		t.Fatal(xs, e)
	}
}
