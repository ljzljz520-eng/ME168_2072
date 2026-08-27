package httpapi

import (
	"gymrecommend/internal/storage"
	"gymrecommend/internal/workflow"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestHealth(t *testing.T) {
	s, e := storage.Open(filepath.Join(t.TempDir(), "x.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	rr := httptest.NewRecorder()
	New(workflow.NewService(s)).Handler().ServeHTTP(rr, httptest.NewRequest("GET", "/health", nil))
	if rr.Code != 200 {
		t.Fatal(rr.Code)
	}
}
