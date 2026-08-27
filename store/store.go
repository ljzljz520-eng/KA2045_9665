package store

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"os"
	"sync"
	"workpay/domain"
)

var buckets = [][]byte{[]byte("records"), []byte("profiles"), []byte("events"), []byte("audits"), []byte("periods")}

type BoltStore struct {
	db *bolt.DB
	mu sync.RWMutex
}

func Open(path string) (*BoltStore, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	s := &BoltStore{db: db}
	err = db.Update(func(tx *bolt.Tx) error {
		for _, b := range buckets {
			if _, e := tx.CreateBucketIfNotExists(b); e != nil {
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
func (s *BoltStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
func (s *BoltStore) put(bucket, key string, v any) error {
	data, e := json.Marshal(v)
	if e != nil {
		return e
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("store closed")
	}
	return s.db.Update(func(tx *bolt.Tx) error { return tx.Bucket([]byte(bucket)).Put([]byte(key), data) })
}
func (s *BoltStore) get(bucket, key string, out any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("store closed")
	}
	return s.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if v == nil {
			return os.ErrNotExist
		}
		return json.Unmarshal(v, out)
	})
}
func (s *BoltStore) PutRecord(r domain.Record) error { return s.put("records", r.ID, r) }
func (s *BoltStore) GetRecord(id string) (domain.Record, error) {
	var r domain.Record
	e := s.get("records", id, &r)
	return r, e
}
func (s *BoltStore) PutProfile(p domain.Profile) error { return s.put("profiles", p.ID, p) }
func (s *BoltStore) GetProfile(id string) (domain.Profile, error) {
	var p domain.Profile
	e := s.get("profiles", id, &p)
	return p, e
}
func (s *BoltStore) PutEvent(e domain.Event) error   { return s.put("events", e.ID, e) }
func (s *BoltStore) PutAudit(a domain.Audit) error   { return s.put("audits", a.ID, a) }
func (s *BoltStore) PutPeriod(p domain.Period) error { return s.put("periods", p.ID, p) }
func (s *BoltStore) GetPeriod(id string) (domain.Period, error) {
	var p domain.Period
	e := s.get("periods", id, &p)
	return p, e
}
func (s *BoltStore) List(bucket string) ([][]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out [][]byte
	if s.db == nil {
		return nil, fmt.Errorf("store closed")
	}
	e := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucket)).ForEach(func(k, v []byte) error { out = append(out, append([]byte(nil), v...)); return nil })
	})
	return out, e
}
