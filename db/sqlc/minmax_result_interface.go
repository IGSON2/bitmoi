package db

import (
	"context"
	"fmt"
)

type MinMaxInterface interface {
	Convert() (int64, int64)
}

func (g *Get1dMinMaxTimeRow) Convert() (min, max int64) {
	if g.Max == nil || g.Min == nil {
		return 0, 0
	}
	return g.Min.(int64), g.Max.(int64)
}

func (g *Get4hMinMaxTimeRow) Convert() (min, max int64) {
	if g.Max == nil || g.Min == nil {
		return 0, 0
	}
	return g.Min.(int64), g.Max.(int64)
}

func (g *Get1hMinMaxTimeRow) Convert() (min, max int64) {
	if g.Max == nil || g.Min == nil {
		return 0, 0
	}
	return g.Min.(int64), g.Max.(int64)
}

func (g *Get15mMinMaxTimeRow) Convert() (min, max int64) {
	if g.Max == nil || g.Min == nil {
		return 0, 0
	}
	return g.Min.(int64), g.Max.(int64)
}

func (g *Get5mMinMaxTimeRow) Convert() (min, max int64) {
	if g.Max == nil || g.Min == nil {
		return 0, 0
	}
	return g.Min.(int64), g.Max.(int64)
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
