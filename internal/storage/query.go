package storage

import (
	"gymrecommend/internal/model"
	"sort"
)

func (s *Store) ListPurchases(member string) ([]model.Purchase, error) {
	var out []model.Purchase
	e := s.list("purchases", func(b []byte) error {
		var p model.Purchase
		if x := decode(b, &p); x != nil {
			return x
		}
		if member == "" || p.MemberID == member {
			out = append(out, p)
		}
		return nil
	})
	return out, e
}
func (s *Store) ListVisits(member string) ([]model.VisitRecord, error) {
	var out []model.VisitRecord
	e := s.list("visits", func(b []byte) error {
		var v model.VisitRecord
		if x := decode(b, &v); x != nil {
			return x
		}
		if member == "" || v.MemberID == member {
			out = append(out, v)
		}
		return nil
	})
	sort.Slice(out, func(i, j int) bool { return out[i].AttendedAt.Before(out[j].AttendedAt) })
	return out, e
}
func (s *Store) ListAudits(entity string) ([]model.Audit, error) {
	var out []model.Audit
	e := s.list("audits", func(b []byte) error {
		var a model.Audit
		if x := decode(b, &a); x != nil {
			return x
		}
		if entity == "" || a.EntityID == entity {
			out = append(out, a)
		}
		return nil
	})
	return out, e
}
func decode(b []byte, v any) error { return unmarshal(b, v) }
