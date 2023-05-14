package utilities

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
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

func Tojson(data interface{}) []byte {
	jsonBytes, err := json.Marshal(data)
	Errchk(err)
	return jsonBytes
}

func ReadText(filename string) []string {
	var dates []string
	files, err := filepath.Glob(filename)
	Errchk(err)
	f, err := os.Open(files[0])
	Errchk(err)
	fmt.Printf("Read Text : %s...\n", files[0])
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		dates = append(dates, scanner.Text())
	}
	return dates
}

func ToByte(data interface{}) []byte {
	var buffer bytes.Buffer
	gob.NewEncoder(&buffer).Encode(data)
	return buffer.Bytes()
}

func MakeRanNum(seedNum, minimum int) int {
	ranSeed := big.NewInt(int64(seedNum - minimum))
	ranBigNum, err := rand.Int(rand.Reader, ranSeed)
	Errchk(err)
	return int(ranBigNum.Int64()) + minimum
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

// 시간단위와 기간단위를 입력해 string 타입으로 변환합니다. 주로 NewKlineService 함수의 인자로 사용합니다.
func MakeInterval(a int, str string) string {
	return strconv.Itoa(a) + strings.ToLower(str)
}

func EntryTimeFormatter(entryTime int64) string {
	timeString := fmt.Sprintf("%d.%02d.%02d %02d:%02d",
		time.UnixMilli(entryTime).Year(),
		time.UnixMilli(entryTime).Month(),
		time.UnixMilli(entryTime).Day(),
		time.UnixMilli(entryTime).Hour(),
		time.UnixMilli(entryTime).Minute(),
	)
	return timeString[2:]
}
