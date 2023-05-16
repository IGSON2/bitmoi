package db

type PriceData struct {
	Name  string  `json:"-"`
	Open  float64 `json:"open"`
	Close float64 `json:"close"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Time  int64   `json:"time"`
}

type VolumeData struct {
	Value float64 `json:"value"`
	Time  int64   `json:"time"`
	Color string  `json:"color"`
}

// func putAllChart(intN int, intU string, allChart map[string]CandleData) {
// 	dbname := dbNameByCase(intN, intU)
// 	newdbpointer := openDB(dbname)
// 	defer newdbpointer.Close()
// 	err := newdbpointer.Update(func(t *bolt.Tx) error {
// 		allChartBucket := t.Bucket([]byte(allChartBucket))
// 		err := allChartBucket.Put([]byte(allChartKey), encode(allChart))
// 		return err
// 	})
// 	utilities.Errchk(err)
// }

// func encode(anystruct interface{}) []byte {
// 	var buffer bytes.Buffer
// 	encoder := gob.NewEncoder(&buffer)
// 	err := encoder.Encode(anystruct)
// 	utilities.Errchk(err)
// 	return buffer.Bytes()
// }

// func dbNameByCase(intN int, intU string) string {
// 	switch intU {
// 	case "m":
// 		switch intN {
// 		case 5:
// 			return DbName_5M
// 		case 15:
// 			return DbName_15M
// 		}
// 	case "h":
// 		switch intN {
// 		case 1:
// 			return DbName_1H
// 		case 4:
// 			return DbName_4H
// 		}
// 	case "d":
// 		return DbName_1D

// 	}
// 	return ""
// }

// func getAllchart(start, end int64, intN int, intU string) map[string]CandleData {
// 	var allChart = make(map[string]CandleData)

// 	cycleNum := cycleByCase(start, end, intN, intU)
// 	for i := 1; i <= cycleNum; i++ {
// 		if i%3 == 0 {
// 			time.Sleep(1 * time.Minute)
// 		}
// 		Xcandles = i * 1000
// 		getOneIntvChart(intN, intU, allPairs)
// 		for j := 0; j < len(allPairs); j++ {
// 			var oneCoinOfOneIntv CandleData
// 			infoByCase(intN, intU, &oneCoinOfOneIntv)
// 			if oneCoinOfOneIntv.PData == nil || oneCoinOfOneIntv.VData == nil {
// 				continue
// 			}
// 			name := oneCoinOfOneIntv.PData[0].Name
// 			var tempPdatas []PriceData = allChart[name].PData
// 			var tempVdatas []VolumeData = allChart[name].VData
// 			tempPdatas = append(oneCoinOfOneIntv.PData, tempPdatas...)
// 			tempVdatas = append(oneCoinOfOneIntv.VData, tempVdatas...)
// 			allChart[name] = CandleData{tempPdatas, tempVdatas}
// 			fmt.Println(name, " Done.")
// 		}
// 		fmt.Printf("%s - %d / %d Done.\n", utilities.MakeInterval(intN, intU), i, cycleNum)
// 	}
// 	return allChart
// }

// func cycleByCase(start, end int64, intN int, intU string) int {
// 	var howHours int
// 	var cyclenum int
// 	termsByHours := (end - start) / (1000 * 60 * 60)
// 	howHours = int(termsByHours/250) + 1

// 	fmt.Println("TermHours : ", termsByHours)
// 	switch intU {
// 	case "m":
// 		switch intN {
// 		case 5:
// 			cyclenum = 3 * howHours
// 		case 15:
// 			cyclenum = howHours
// 		}
// 	case "h":
// 		switch intN {
// 		case 1:
// 			cyclenum = int(howHours/4) + 1
// 		case 4:
// 			cyclenum = int(howHours/16) + 1
// 		}
// 	case "d":
// 		cyclenum = int(howHours/96) + 1
// 	}
// 	fmt.Println("CycleNum : ", cyclenum)
// 	return cyclenum
// }

// func infoByCase(intN int, intU string, emptyInfo *CandleData) {

// 	switch intU {
// 	case "m":
// 		switch intN {
// 		case 5:
// 			*emptyInfo = <-FiveMinuteCh
// 		case 15:
// 			*emptyInfo = <-FifMinuteCh
// 		}
// 	case "h":
// 		switch intN {
// 		case 1:
// 			*emptyInfo = <-OneHourCh
// 		case 4:
// 			*emptyInfo = <-FourHourCh
// 		}

// 	case "d":
// 		*emptyInfo = <-OneDayCh
// 	}
// }

// func UpdateDB() {
// 	DB_5M := AC.InitAllchart(FiveM)
// 	DB_15M := AC.InitAllchart(FifM)
// 	DB_1H := AC.InitAllchart(OneH)
// 	DB_4H := AC.InitAllchart(FourH)
// 	DB_1D := AC.InitAllchart(OneD)

// 	// fmt.Printf(
// 	// 	"BEFORE_1H : %s ~ %s\tLen : %d\n",
// 	// 	utilities.TransMilli((*AC.oneD)["BTCUSDT"].PData[0].Time),
// 	// 	utilities.TransMilli((*AC.oneD)["BTCUSDT"].PData[len((*AC.oneD)["BTCUSDT"].PData)-1].Time),
// 	// 	len((*AC.oneD)["BTCUSDT"].PData),
// 	// )

// 	start := int64(1568073600000)
// 	// start := (*DB_1D)["BTCUSDT"].PData[len((*DB_1D)["BTCUSDT"].PData)-1].Time
// 	end := yesterday9AM()
// 	sliced5M := sliceChart(start, end, 5, "m")
// 	sliced15M := sliceChart(start, end, 15, "m")
// 	sliced1H := sliceChart(start, end, 1, "h")
// 	sliced4H := sliceChart(start, end, 4, "h")
// 	sliced1D := sliceChart(start, end, 1, "d")

// 	for _, name := range allPairs {
// 		temp5MPdata := append((*DB_5M)[name].PData, sliced5M[name].PData...)
// 		temp5MVdata := append((*DB_5M)[name].VData, sliced5M[name].VData...)
// 		(*DB_5M)[name] = CandleData{temp5MPdata, temp5MVdata}

// 		temp15MPdata := append((*DB_15M)[name].PData, sliced15M[name].PData...)
// 		temp15MVdata := append((*DB_15M)[name].VData, sliced15M[name].VData...)
// 		(*DB_15M)[name] = CandleData{temp15MPdata, temp15MVdata}

// 		temp1HPdata := append((*DB_1H)[name].PData, sliced1H[name].PData...)
// 		temp1HVdata := append((*DB_1H)[name].VData, sliced1H[name].VData...)
// 		(*DB_1H)[name] = CandleData{temp1HPdata, temp1HVdata}

// 		temp4HPdata := append((*DB_4H)[name].PData, sliced4H[name].PData...)
// 		temp4HVdata := append((*DB_4H)[name].VData, sliced4H[name].VData...)
// 		(*DB_4H)[name] = CandleData{temp4HPdata, temp4HVdata}

// 		temp1DPdata := append((*DB_1D)[name].PData, sliced1D[name].PData...)
// 		temp1DVdata := append((*DB_1D)[name].VData, sliced1D[name].VData...)
// 		(*DB_1D)[name] = CandleData{temp1DPdata, temp1DVdata}
// 	}

// 	putAllChart(5, "m", *DB_5M)
// 	putAllChart(15, "m", *DB_15M)
// 	putAllChart(1, "h", *DB_1H)
// 	putAllChart(4, "h", *DB_4H)
// 	putAllChart(1, "d", *DB_1D)
// 	fmt.Printf("%s ~ %s \n", utilities.TransMilli((*DB_5M)["DYDXUSDT"].PData[0].Time), utilities.TransMilli((*DB_5M)["DYDXUSDT"].PData[len((*DB_5M)["DYDXUSDT"].PData)-1].Time))
// 	fmt.Printf("%s ~ %s \n", utilities.TransMilli((*DB_15M)["DYDXUSDT"].PData[0].Time), utilities.TransMilli((*DB_15M)["DYDXUSDT"].PData[len((*DB_15M)["DYDXUSDT"].PData)-1].Time))
// 	fmt.Printf("%s ~ %s \n", utilities.TransMilli((*DB_1H)["DYDXUSDT"].PData[0].Time), utilities.TransMilli((*DB_1H)["DYDXUSDT"].PData[len((*DB_1H)["DYDXUSDT"].PData)-1].Time))
// 	fmt.Printf("%s ~ %s \n", utilities.TransMilli((*DB_4H)["DYDXUSDT"].PData[0].Time), utilities.TransMilli((*DB_4H)["DYDXUSDT"].PData[len((*DB_4H)["DYDXUSDT"].PData)-1].Time))
// 	fmt.Printf("%s ~ %s \n", utilities.TransMilli((*DB_1D)["DYDXUSDT"].PData[0].Time), utilities.TransMilli((*DB_1D)["DYDXUSDT"].PData[len((*DB_1D)["DYDXUSDT"].PData)-1].Time))
// }

// func sliceChart(start, end int64, intN int, intU string) map[string]CandleData {
// 	tilTodayChart := getAllchart(start, end, intN, intU)
// 	time.Sleep(1 * time.Minute)
// 	return tilTodayChart
// 	// 	var slicedChart map[string]CandleData = map[string]CandleData{}
// 	// 	var sliceStart int
// 	// 	var sliceEnd int

// 	// Outer:
// 	// 	for _, name := range allPairs {
// 	// 		var endbool bool = false
// 	// 		var startbool bool = false

// 	// 		for i := len(tilTodayChart[name].PData) - 1; i > 0; i-- {
// 	// 			if tilTodayChart[name].PData == nil {
// 	// 				fmt.Printf("TillTodayChart[%s] is Nil.\n", name)
// 	// 				break
// 	// 			}
// 	// 			if tilTodayChart[name].PData[i].Time == end {
// 	// 				sliceEnd = i
// 	// 				endbool = true
// 	// 			}
// 	// 			j := len(tilTodayChart[name].PData) - (i + 1)
// 	// 			if tilTodayChart[name].PData[j].Time == start {
// 	// 				sliceStart = j
// 	// 				startbool = true
// 	// 			}
// 	// 			if endbool && startbool {
// 	// 				break
// 	// 			}
// 	// 			if i == 1 && (!endbool || !startbool) {
// 	// 				fmt.Printf("TillTodayChart[%s]'s lenth is not enough.\nEndbool : %v / Startbool : %v\nStart : %d / End : %d\n%d ~ %d\n",
// 	// 					name, endbool, startbool, start, end, tilTodayChart[name].PData[0].Time, tilTodayChart[name].PData[len(tilTodayChart[name].PData)-1].Time)
// 	// 				continue Outer
// 	// 			}
// 	// 		}
// 	// 		slicedChart[name] = CandleData{
// 	// 			PData: tilTodayChart[name].PData[sliceStart+1 : sliceEnd+1],
// 	// 			VData: tilTodayChart[name].VData[sliceStart+1 : sliceEnd+1],
// 	// 		}
// 	// 	}
// 	// 	time.Sleep(1 * time.Minute)
// 	// 	return slicedChart
// }
