package importcsv

import (
	"encoding/csv"
	"gymrecommend/internal/model"
	"io"
	"strconv"
	"strings"
	"time"
)

func ReadMembers(r io.Reader) ([]model.Member, error) {
	rows, e := csv.NewReader(r).ReadAll()
	if e != nil {
		return nil, e
	}
	out := []model.Member{}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 6 {
			continue
		}
		rating, _ := strconv.ParseFloat(row[4], 64)
		out = append(out, model.Member{ID: row[0], Name: row[1], Goals: strings.Split(row[2], ","), PreferredSlot: row[3], CoachRating: rating, Package: row[5], CreatedAt: time.Now()})
	}
	return out, nil
}
func ReadClasses(r io.Reader) ([]model.Class, error) {
	rows, e := csv.NewReader(r).ReadAll()
	if e != nil {
		return nil, e
	}
	out := []model.Class{}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 7 {
			continue
		}
		cap, _ := strconv.Atoi(row[6])
		out = append(out, model.Class{ID: row[0], Title: row[1], Kind: row[2], Slot: row[3], Coach: row[4], Level: row[5], Capacity: cap})
	}
	return out, nil
}
func WriteRecommendations(w io.Writer, items []model.Recommendation) error {
	cw := csv.NewWriter(w)
	if e := cw.Write([]string{"id", "member_id", "class_id", "kind", "score", "reason"}); e != nil {
		return e
	}
	for _, r := range items {
		if e := cw.Write([]string{r.ID, r.MemberID, r.ClassID, r.Kind, strconv.FormatFloat(r.Score, 'f', 3, 64), r.Reason}); e != nil {
			return e
		}
	}
	cw.Flush()
	return cw.Error()
}
