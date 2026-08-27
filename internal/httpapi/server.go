package httpapi

import (
	"encoding/json"
	"gymrecommend/internal/workflow"
	"net/http"
	"strings"
)

type Server struct{ svc *workflow.Service }

func New(s *workflow.Service) *Server { return &Server{svc: s} }
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/recommend", s.recommend)
	mux.HandleFunc("/summary", s.summary)
	return mux
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
func (s *Server) recommend(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("member"))
	items, e := s.svc.Recommend(id)
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(items)
}
func (s *Server) summary(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("member"))
	out, e := s.svc.Summary(id)
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(out)
}
