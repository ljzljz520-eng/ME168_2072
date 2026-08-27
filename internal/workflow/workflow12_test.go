package workflow

import "testing"

func TestWorkflow12(t *testing.T) {
	s := setup(t)
	if _, e := s.Recommend("m"); e != nil {
		t.Fatal(e)
	}
	if e := s.Confirm("m"); e != nil {
		t.Fatal(e)
	}
	got, e := s.Summary("m")
	if e != nil {
		t.Fatal(e)
	}
	if got["count"] != 1 {
		t.Fatalf("summary count %.0f", got["count"])
	}
	if s.Stage() != StageIdle {
		t.Fatalf("summary left workflow in %s", s.Stage())
	}
	if _, e = s.Recommend("m"); e != nil {
		t.Fatal(e)
	}
	next, e := s.Summary("m")
	if e != nil {
		t.Fatal(e)
	}
	if next["count"] != 1 {
		t.Fatalf("expected independent summary, got %.0f", next["count"])
	}
}
