package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"bitmoi/backend/utilities/common"
)

type CandlesInterface interface {
	Interval() string
	Name() string
	EntryTime() string
	InitCandleData() *CandleData
	InitAdvancedCandleData() []PriceVolumeData
}

type Candles1dSlice []db.Candles1d

func (c *Candles1dSlice) EntryTime() string {
	slice := ([]db.Candles1d)(*c)
	return utilities.EntryTimeFormatter(slice[len(slice)-1].Time)
}

func (c *Candles1dSlice) Interval() string {
	return db.OneD
}

func (c *Candles1dSlice) Name() string {
	return ([]db.Candles1d)(*c)[0].Name
}

func (c *Candles1dSlice) InitCandleData() *CandleData {
	var pDataSlice []PriceData
	var vDataSlice []VolumeData
	for _, candle := range *c {
		pDataSlice = append(pDataSlice, PriceData{
			Open:  common.FloorDecimal(candle.Open),
			Close: common.FloorDecimal(candle.Close),
			High:  common.FloorDecimal(candle.High),
			Low:   common.FloorDecimal(candle.Low),
			Time:  candle.Time,
		})

		vDataSlice = append(vDataSlice, VolumeData{
			Value: common.FloorDecimal(candle.Volume),
			Time:  candle.Time,
			Color: candle.Color,
		})
	}
	return &CandleData{pDataSlice, vDataSlice}
}

func (c *Candles1dSlice) InitAdvancedCandleData() []PriceVolumeData {
	var pvDataSlice []PriceVolumeData
	for _, row := range *c {
		newPvData := PriceVolumeData{
			Open:   row.Open,
			Close:  row.Close,
			High:   row.High,
			Low:    row.Low,
			Volume: row.Volume,
			Time:   (row.Time + 7200) * 1000,
		}
		pvDataSlice = append([]PriceVolumeData{newPvData}, pvDataSlice...)
	}
	return pvDataSlice
}

type Candles4hSlice []db.Candles4h

func (c *Candles4hSlice) EntryTime() string {
	slice := ([]db.Candles4h)(*c)
	return utilities.EntryTimeFormatter(slice[len(slice)-1].Time)
}

func (c *Candles4hSlice) Interval() string {
	return db.FourH
}

func (c *Candles4hSlice) Name() string {
	return ([]db.Candles4h)(*c)[0].Name
}

func (c *Candles4hSlice) InitCandleData() *CandleData {
	var pDataSlice []PriceData
	var vDataSlice []VolumeData

	for _, candle := range *c {
		pDataSlice = append(pDataSlice, PriceData{
			Open:  common.FloorDecimal(candle.Open),
			Close: common.FloorDecimal(candle.Close),
			High:  common.FloorDecimal(candle.High),
			Low:   common.FloorDecimal(candle.Low),
			Time:  candle.Time,
		})

		vDataSlice = append(vDataSlice, VolumeData{
			Value: common.FloorDecimal(candle.Volume),
			Time:  candle.Time,
			Color: candle.Color,
		})
	}
	return &CandleData{pDataSlice, vDataSlice}
}

func (c *Candles4hSlice) InitAdvancedCandleData() []PriceVolumeData {
	var pvDataSlice []PriceVolumeData
	for _, row := range *c {
		newPvData := PriceVolumeData{
			Open:   row.Open,
			Close:  row.Close,
			High:   row.High,
			Low:    row.Low,
			Volume: row.Volume,
			Time:   row.Time * 1000,
		}
		pvDataSlice = append([]PriceVolumeData{newPvData}, pvDataSlice...)
	}
	return pvDataSlice
}

type Candles1hSlice []db.Candles1h

func (c *Candles1hSlice) EntryTime() string {
	slice := ([]db.Candles1h)(*c)
	return utilities.EntryTimeFormatter(slice[len(slice)-1].Time)
}

func (c *Candles1hSlice) Interval() string {
	return db.OneH
}

func (c *Candles1hSlice) Name() string {
	return ([]db.Candles1h)(*c)[0].Name
}

func (c *Candles1hSlice) InitCandleData() *CandleData {
	var pDataSlice []PriceData
	var vDataSlice []VolumeData

	for _, candle := range *c {
		pDataSlice = append(pDataSlice, PriceData{
			Open:  common.FloorDecimal(candle.Open),
			Close: common.FloorDecimal(candle.Close),
			High:  common.FloorDecimal(candle.High),
			Low:   common.FloorDecimal(candle.Low),
			Time:  candle.Time,
		})

		vDataSlice = append(vDataSlice, VolumeData{
			Value: common.FloorDecimal(candle.Volume),
			Time:  candle.Time,
			Color: candle.Color,
		})
	}
	return &CandleData{pDataSlice, vDataSlice}
}

func (c *Candles1hSlice) InitAdvancedCandleData() []PriceVolumeData {
	var pvDataSlice []PriceVolumeData
	for _, row := range *c {
		newPvData := PriceVolumeData{
			Open:   row.Open,
			Close:  row.Close,
			High:   row.High,
			Low:    row.Low,
			Volume: row.Volume,
			Time:   row.Time * 1000,
		}
		pvDataSlice = append([]PriceVolumeData{newPvData}, pvDataSlice...)
	}
	return pvDataSlice
}

type Candles15mSlice []db.Candles15m

func (c *Candles15mSlice) EntryTime() string {
	slice := ([]db.Candles15m)(*c)
	return utilities.EntryTimeFormatter(slice[len(slice)-1].Time)
}

func (c *Candles15mSlice) Interval() string {
	return db.FifM
}

func (c *Candles15mSlice) Name() string {
	return ([]db.Candles15m)(*c)[0].Name
}

func (c *Candles15mSlice) InitCandleData() *CandleData {
	var pDataSlice []PriceData
	var vDataSlice []VolumeData

	for _, candle := range *c {
		pDataSlice = append(pDataSlice, PriceData{
			Open:  common.FloorDecimal(candle.Open),
			Close: common.FloorDecimal(candle.Close),
			High:  common.FloorDecimal(candle.High),
			Low:   common.FloorDecimal(candle.Low),
			Time:  candle.Time,
		})

		vDataSlice = append(vDataSlice, VolumeData{
			Value: common.FloorDecimal(candle.Volume),
			Time:  candle.Time,
			Color: candle.Color,
		})
	}
	return &CandleData{pDataSlice, vDataSlice}
}

func (c *Candles15mSlice) InitAdvancedCandleData() []PriceVolumeData {
	var pvDataSlice []PriceVolumeData
	for _, row := range *c {
		newPvData := PriceVolumeData{
			Open:   row.Open,
			Close:  row.Close,
			High:   row.High,
			Low:    row.Low,
			Volume: row.Volume,
			Time:   row.Time * 1000,
		}
		pvDataSlice = append([]PriceVolumeData{newPvData}, pvDataSlice...)
	}
	return pvDataSlice
}

type Candles5mSlice []db.Candles5m

func (c *Candles5mSlice) EntryTime() string {
	slice := ([]db.Candles5m)(*c)
	return utilities.EntryTimeFormatter(slice[len(slice)-1].Time)
}

func (c *Candles5mSlice) Interval() string {
	return db.FifM
}

func (c *Candles5mSlice) Name() string {
	return ([]db.Candles5m)(*c)[0].Name
}

func (c *Candles5mSlice) InitCandleData() *CandleData {
	var pDataSlice []PriceData
	var vDataSlice []VolumeData

	for _, candle := range *c {
		pDataSlice = append(pDataSlice, PriceData{
			Open:  common.FloorDecimal(candle.Open),
			Close: common.FloorDecimal(candle.Close),
			High:  common.FloorDecimal(candle.High),
			Low:   common.FloorDecimal(candle.Low),
			Time:  candle.Time,
		})

		vDataSlice = append(vDataSlice, VolumeData{
			Value: common.FloorDecimal(candle.Volume),
			Time:  candle.Time,
			Color: candle.Color,
		})
	}
	return &CandleData{pDataSlice, vDataSlice}
}

func (c *Candles5mSlice) InitAdvancedCandleData() []PriceVolumeData {
	var pvDataSlice []PriceVolumeData
	for _, row := range *c {
		newPvData := PriceVolumeData{
			Open:   row.Open,
			Close:  row.Close,
			High:   row.High,
			Low:    row.Low,
			Volume: row.Volume,
			Time:   row.Time * 1000,
		}
		pvDataSlice = append([]PriceVolumeData{newPvData}, pvDataSlice...)
	}
	return pvDataSlice
}
