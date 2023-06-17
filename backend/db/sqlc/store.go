package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	SelectMinMaxTime(interval, name string, c context.Context) (int64, int64, error)
}

type SqlStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SqlStore{
		Queries: New(db),
		db:      db,
	}
}

func (s *SqlStore) SelectMinMaxTime(interval, name string, c context.Context) (int64, int64, error) {
	switch interval {
	case OneD:
		r, err := s.Get1dMinMaxTime(c, name)
		min, max := convTimestamp(&r)
		return min, max, err
	case FourH:
		r, err := s.Get4hMinMaxTime(c, name)
		min, max := convTimestamp(&r)
		return min, max, err
	case OneH:
		r, err := s.Get1hMinMaxTime(c, name)
		min, max := convTimestamp(&r)
		return min, max, err
	case FifM:
		r, err := s.Get15mMinMaxTime(c, name)
		min, max := convTimestamp(&r)
		return min, max, err
	case FiveM:
		r, err := s.Get5mMinMaxTime(c, name)
		min, max := convTimestamp(&r)
		return min, max, err
	}
	return 0, 0, fmt.Errorf("invalid interval %s", interval)
}

func convTimestamp(m MinMaxInterface) (min, max int64) {
	return m.Convert()
}
