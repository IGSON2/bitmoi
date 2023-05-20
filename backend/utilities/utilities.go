package utilities

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
func StrToFloat(volstr string) float64 {
	volumefloat, err := strconv.ParseFloat(volstr, 64)
	Errchk(err)
	volfloatstr := fmt.Sprintf("%.5f", volumefloat)
	volfloatstrdot, err := strconv.ParseFloat(volfloatstr, 64)
	Errchk(err)
	return volfloatstrdot
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
