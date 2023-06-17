package db

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
