package db

import (
	"bitmoi/backend/utilities"
	"bytes"
	"encoding/gob"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

const (
	DB_Path        = "./backend/db/chartData/"
	DbName_5M      = "future_5M_Chart.db"
	DbName_15M     = "future_15M_Chart.db"
	DbName_1H      = "future_1H_Chart.db"
	DbName_4H      = "future_4H_Chart.db"
	DbName_1D      = "future_1D_Chart.db"
	FiveM          = "5m"
	FifM           = "15m"
	OneH           = "1h"
	FourH          = "4h"
	OneD           = "1d"
	allChartBucket = "allchartbucket"
	allChartKey    = "allChartKey"
)

type dbCharts map[string]CandleData

type AllChart struct {
	fivM  *dbCharts
	fifM  *dbCharts
	oneH  *dbCharts
	fourH *dbCharts
	oneD  *dbCharts
}

var AC = &AllChart{}

func openDB(dbName string) *bolt.DB {
	newDBPointer, err := bolt.Open(DB_Path+dbName, 0644, bolt.DefaultOptions)
	utilities.Errchk(err)
	newDBPointer.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(allChartBucket))
		return err
	})
	return newDBPointer
}

func (a *AllChart) InitAllchart(interval string) *dbCharts {
	switch interval {
	case "5m":
		a.fivM = a.fivM.initIntervalChart(DbName_5M)
		return a.fivM
	case "15m":
		a.fifM = a.fifM.initIntervalChart(DbName_15M)
		return a.fifM
	case "1h":
		a.oneH = a.oneH.initIntervalChart(DbName_1H)
		return a.oneH
	case "4h":
		a.fourH = a.fourH.initIntervalChart(DbName_4H)
		return a.fourH
	case "1d":
		a.oneD = a.oneD.initIntervalChart(DbName_1D)
		return a.oneD
	default:
		return nil
	}
}

func (d *dbCharts) initIntervalChart(dbName string) *dbCharts {
	if d == nil {
		fmt.Printf("Initiate %s...\n", dbName)
		newDBpointer := openDB(dbName)
		d = loadChart(newDBpointer)
		defer newDBpointer.Close()
	}
	return d
}

func loadChart(b *bolt.DB) *dbCharts {
	var allChartByte []byte
	err := b.View(func(t *bolt.Tx) error {
		allChartBucket := t.Bucket([]byte(allChartBucket))
		allChartByte = allChartBucket.Get([]byte(allChartKey))
		return nil
	})
	utilities.Errchk(err)
	if allChartByte == nil {
		return &dbCharts{}
	} else {
		return dbDataDecode(allChartByte)
	}
}

func dbDataDecode(dataBytes []byte) *dbCharts {
	var emptyAllChart dbCharts
	var buffer bytes.Buffer
	_, err := buffer.Write(dataBytes)
	utilities.Errchk(err)
	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&emptyAllChart)
	utilities.Errchk(err)
	return &emptyAllChart
}
