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
	return math.Ceil(floatData*1000000) / 1000000
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

func SplitPairnames(names string) []string {
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

// SplitAndTrim splits input separated by a comma
// and trims excessive white space from the substrings.
func SplitAndTrim(input string) (ret []string) {
	l := strings.Split(input, ",")
	for _, r := range l {
		if r = strings.TrimSpace(r); r != "" {
			ret = append(ret, r)
		}
	}
	return ret
}

func GenerateEmailMessage(user, url string) string {
	return fmt.Sprintf(`<h3>%s</h3>님 안녕하세요! </br>
	바로 시작하는 모의투자 비트모이에 가입하신 것을 환영합니다!<br/>
	아래의 인증 링크를 클릭하여 인증을 완료해주세요.<br/>
	<h3><a href="%s" style="color:black">인증하기</a></h3>
	`, user, url)

}

func ToDBTimestamp(binance int64) int64 { return (binance / 1000) + 32400 }
func ToBinanceTimestamp(db int64) int64 { return (db - 32400) * 1000 }
