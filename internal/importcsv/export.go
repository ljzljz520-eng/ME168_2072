package importcsv

import (
	"encoding/json"
	"gymrecommend/internal/model"
	"io"
)

func WriteJSON(w io.Writer, value any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(value)
}
func ReadProfile(data []byte) (model.GoalProfile, error) {
	var p model.GoalProfile
	e := json.Unmarshal(data, &p)
	return p, e
}
