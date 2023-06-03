package db

import "strconv"

const (
	FiveM = "5m"
	FifM  = "15m"
	OneH  = "1h"
	FourH = "4h"
	OneD  = "1d"
)

// 시간단위와 기간단위를 입력해 string 타입으로 변환합니다. 주로 NewKlineService 함수의 인자로 사용합니다.
func MakeInterval(a int, str string) string {
	switch str {
	case "d", "D":
		return OneD
	case "h", "H":
		switch a {
		case 4:
			return FourH
		case 1:
			return OneH
		}
	case "m", "M":
		return FifM
	}
	return "invalid"
}

func ParseInterval(interval string) (int, string) {
	timeStr, unit := interval[:len(interval)-1], interval[len(interval)-1]
	time, err := strconv.Atoi(timeStr)
	if err != nil {
		return 4, "h"
	}
	return time, string(unit)
}

func CalculateTerm(interval string, waitingTerm int) int {
	switch interval {
	case OneD:
		return waitingTerm
	case FourH:
		return waitingTerm * 6
	case OneH:
		return waitingTerm * 24
	case FifM:
		return waitingTerm * 96
	case FiveM:
		return waitingTerm * 288
	}
	return 0
}
