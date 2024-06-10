package db

type InsertQueryInterface interface {
	SetCandleInfo(open, close, high, low, volume float64, time int64, name, color string)
}

func (i *Insert1dCandlesParams) SetCandleInfo(open, close, high, low, volume float64, time int64, name, color string) {
	*i = Insert1dCandlesParams{
		Name:   name,
		Open:   open,
		Close:  close,
		High:   high,
		Low:    low,
		Volume: volume,
		Time:   time,
		Color:  color,
	}
}
func (i *Insert4hCandlesParams) SetCandleInfo(open, close, high, low, volume float64, time int64, name, color string) {
	*i = Insert4hCandlesParams{
		Name:   name,
		Open:   open,
		Close:  close,
		High:   high,
		Low:    low,
		Volume: volume,
		Time:   time,
		Color:  color,
	}
}
func (i *Insert1hCandlesParams) SetCandleInfo(open, close, high, low, volume float64, time int64, name, color string) {
	*i = Insert1hCandlesParams{
		Name:   name,
		Open:   open,
		Close:  close,
		High:   high,
		Low:    low,
		Volume: volume,
		Time:   time,
		Color:  color,
	}
}
func (i *Insert15mCandlesParams) SetCandleInfo(open, close, high, low, volume float64, time int64, name, color string) {
	*i = Insert15mCandlesParams{
		Name:   name,
		Open:   open,
		Close:  close,
		High:   high,
		Low:    low,
		Volume: volume,
		Time:   time,
		Color:  color,
	}
}
func (i *Insert5mCandlesParams) SetCandleInfo(open, close, high, low, volume float64, time int64, name, color string) {
	*i = Insert5mCandlesParams{
		Name:   name,
		Open:   open,
		Close:  close,
		High:   high,
		Low:    low,
		Volume: volume,
		Time:   time,
		Color:  color,
	}
}
