package httpapi

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func methodAllowed(w http.ResponseWriter, allowed string) {
	w.Header().Set("Allow", allowed)
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
