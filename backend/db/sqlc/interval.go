package db

import (
	"strconv"
	"time"
)

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

func CalculateWaitingTerm(interval string, waitingTerm int) int {
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
	default:
		return 0
	}
}

func GetIntervalStep(interval string) int64 {
	switch interval {
	case OneD:
		return int64(24 * time.Hour.Seconds())
	case FourH:
		return int64(4 * time.Hour.Seconds())
	case OneH:
		return int64(1 * time.Hour.Seconds())
	case FifM:
		return int64(15 * time.Minute.Seconds())
	case FiveM:
		return int64(5 * time.Minute.Seconds())
	default:
		return 0
	}
}

// FiveM 임시 제외
func GetAnotherIntervals(interval string) []string {
	switch interval {
	case OneD:
		return []string{FourH, OneH, FifM}
	case FourH:
		return []string{OneD, OneH, FifM}
	case OneH:
		return []string{OneD, FourH, FifM}
	case FifM:
		return []string{OneD, FourH, OneH}
	default:
		return []string{}
	}
}

func GetIntervalFromRange(from, to int64) string {
	diff := to - from
	if diff > 86400 {
		return OneD
	} else if diff > 14400 {
		return FourH
	} else if diff > 3600 {
		return OneH
	} else if diff > 900 {
		return FifM
	}
	return ""
}

func ConvertResolution(resolution string) string {
	switch resolution {
	case "1D":
		return OneD
	case "240":
		return FourH
	case "60":
		return OneH
	case "15":
		return FifM
	default:
		return ""
	}
}
