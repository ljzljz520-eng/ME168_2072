package storage

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"gymrecommend/internal/model"
	"path/filepath"
	"sync"
	"time"
)

var buckets = []string{"members", "classes", "recommendations", "purchases", "visits", "profiles", "audits"}

type Store struct {
	db *bbolt.DB
	mu sync.RWMutex
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(filepath.Clean(path), 0600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, b := range buckets {
			if _, e := tx.CreateBucketIfNotExists([]byte(b)); e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
func (s *Store) put(bucket, key string, v any) error {
	data, e := json.Marshal(v)
	if e != nil {
		return e
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(bucket)).Put([]byte(key), data) })
}
func (s *Store) get(bucket, key string, v any) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if b == nil {
			return fmt.Errorf("not found")
		}
		return json.Unmarshal(b, v)
	})
}
func (s *Store) list(bucket string, out func([]byte) error) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(bucket)).ForEach(func(_, v []byte) error {
			if v == nil {
				return nil
			}
			return out(v)
		})
	})
}
func (s *Store) SaveMember(m model.Member) error { return s.put("members", m.ID, m) }
func (s *Store) Member(id string) (model.Member, error) {
	var m model.Member
	e := s.get("members", id, &m)
	return m, e
}
func (s *Store) SaveClass(c model.Class) error { return s.put("classes", c.ID, c) }
func (s *Store) Class(id string) (model.Class, error) {
	var c model.Class
	e := s.get("classes", id, &c)
	return c, e
}
func (s *Store) SaveRecommendation(r model.Recommendation) error {
	return s.put("recommendations", r.ID, r)
}
func (s *Store) SavePurchase(p model.Purchase) error   { return s.put("purchases", p.ID, p) }
func (s *Store) SaveVisit(v model.VisitRecord) error   { return s.put("visits", v.ID, v) }
func (s *Store) SaveProfile(p model.GoalProfile) error { return s.put("profiles", p.MemberID, p) }
func (s *Store) SaveAudit(a model.Audit) error         { return s.put("audits", a.ID, a) }
func (s *Store) ListMembers() ([]model.Member, error) {
	var out []model.Member
	e := s.list("members", func(b []byte) error {
		var m model.Member
		if x := json.Unmarshal(b, &m); x != nil {
			return x
		}
		out = append(out, m)
		return nil
	})
	return out, e
}
func (s *Store) ListClasses() ([]model.Class, error) {
	var out []model.Class
	e := s.list("classes", func(b []byte) error {
		var c model.Class
		if x := json.Unmarshal(b, &c); x != nil {
			return x
		}
		out = append(out, c)
		return nil
	})
	return out, e
}
func (s *Store) ListRecommendations(member string) ([]model.Recommendation, error) {
	var out []model.Recommendation
	e := s.list("recommendations", func(b []byte) error {
		var r model.Recommendation
		if x := json.Unmarshal(b, &r); x != nil {
			return x
		}
		if member == "" || r.MemberID == member {
			out = append(out, r)
		}
		return nil
	})
	return out, e
}
