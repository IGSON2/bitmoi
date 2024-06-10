package db

type GetCandlesInterface interface {
	GetName() string
	GetTime() int64
	GetLimit() int32
	GetInterval() string
}

func (p *Get5mCandlesParams) GetName() string {
	return p.Name
}

func (p *Get5mCandlesParams) GetTime() int64 {
	return p.Time
}

func (p *Get5mCandlesParams) GetLimit() int32 {
	return p.Limit
}

func (p *Get5mCandlesParams) GetInterval() string {
	return FiveM
}

func (p *Get15mCandlesParams) GetName() string {
	return p.Name
}

func (p *Get15mCandlesParams) GetTime() int64 {
	return p.Time
}

func (p *Get15mCandlesParams) GetLimit() int32 {
	return p.Limit
}

func (p *Get15mCandlesParams) GetInterval() string {
	return FifM
}

func (p *Get1hCandlesParams) GetName() string {
	return p.Name
}

func (p *Get1hCandlesParams) GetTime() int64 {
	return p.Time
}

func (p *Get1hCandlesParams) GetLimit() int32 {
	return p.Limit
}

func (p *Get1hCandlesParams) GetInterval() string {
	return OneH
}

func (p *Get4hCandlesParams) GetName() string {
	return p.Name
}

func (p *Get4hCandlesParams) GetTime() int64 {
	return p.Time
}

func (p *Get4hCandlesParams) GetLimit() int32 {
	return p.Limit
}

func (p *Get4hCandlesParams) GetInterval() string {
	return FourH
}

func (p *Get1dCandlesParams) GetName() string {
	return p.Name
}

func (p *Get1dCandlesParams) GetTime() int64 {
	return p.Time
}

func (p *Get1dCandlesParams) GetLimit() int32 {
	return p.Limit
}

func (p *Get1dCandlesParams) GetInterval() string {
	return OneD
}
