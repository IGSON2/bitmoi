package utilities

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Errchk(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func ReadText(filename string) []string {
	var datas []string
	files, err := filepath.Glob(filename)
	Errchk(err)
	f, err := os.Open(files[0])
	Errchk(err)
	defer f.Close()
	fmt.Printf("Read Text : %s...\n", files[0])
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		datas = append(datas, scanner.Text())
	}
	return datas
}

func ToByte(data interface{}) []byte {
	var buffer bytes.Buffer
	gob.NewEncoder(&buffer).Encode(data)
	return buffer.Bytes()
}

// string 타입을 float64 타입으로 변환합니다.
func StrToFloat(strData string) float64 {
	floatData, err := strconv.ParseFloat(strData, 64)
	Errchk(err)
	return math.Floor(floatData*10000) / 10000
}

// float64 타입을 string 타입으로 변환합니다.
func FloatToStr(floatnum float64) string {
	floatStr := strconv.FormatFloat(floatnum, 'f', 5, 64)
	return floatStr
}

// m 까지의 Unixmilli시간을 알아보기 쉽게 재 포멧 합니다.
func TransMilli(m int64) string {
	date := fmt.Sprintf("%02d년%02d월%02d일%02d시%02d분", time.UnixMilli(m).Year(), time.UnixMilli(m).Month(), time.UnixMilli(m).Day(), time.UnixMilli(m).Hour(), time.UnixMilli(m).Minute())
	return date
}

func EntryTimeFormatter(entryTime int64) string {
	timeString := fmt.Sprintf("%d.%02d.%02d %02d:%02d",
		time.Unix(entryTime, 0).Year(),
		time.Unix(entryTime, 0).Month(),
		time.Unix(entryTime, 0).Day(),
		time.Unix(entryTime, 0).Hour(),
		time.Unix(entryTime, 0).Minute(),
	)
	return timeString[2:]
}

// 어제 오전 9시의 Unix Millsecond timestamp를 반환합니다.
func Yesterday9AM() int64 {
	nineAMmilli := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 9, 1, 0, 0, time.Local).UnixMilli()
	return nineAMmilli
}

func FindDiffPair(allPair, history []string) string {
	var ranName string
Outer:
	for {
		ranName = allPair[MakeRanInt(0, len(allPair))]
		var sameHere bool = false
		for _, name := range history {
			if ranName == name {
				sameHere = true
			}
		}
		if !sameHere {
			break Outer
		}
	}
	return ranName
}

func Splitnames(names string) []string {
	var splited []string
	if names != "" {
		splitted := strings.Split(names, ",")
		for _, str := range splitted {
			if str != "" && !strings.Contains(str, "STAGE") {
				splited = append(splited, str)
			}
		}
	}
	return splited
}

func cycleByCase(start, end int64, intN int, intU string) int {
	var howHours int
	var cyclenum int
	termsByHours := (end - start) / (1000 * 60 * 60)
	howHours = int(termsByHours/250) + 1

	fmt.Println("TermHours : ", termsByHours)
	switch intU {
	case "m":
		switch intN {
		case 5:
			cyclenum = 3 * howHours
		case 15:
			cyclenum = howHours
		}
	case "h":
		switch intN {
		case 1:
			cyclenum = int(howHours/4) + 1
		case 4:
			cyclenum = int(howHours/16) + 1
		}
	case "d":
		cyclenum = int(howHours/96) + 1
	}
	fmt.Println("CycleNum : ", cyclenum)
	return cyclenum
}
