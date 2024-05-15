package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) sendAnotherInterval(a *AnotherIntervalRequest, c *fiber.Ctx) (*OnePairChart, error) {
	originInfo, err := utilities.DecodeIdentificationData(a.Identifier)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal chart identifier. err : %w", err)
	}

	cdd, err := s.selectStageChart(originInfo.Name, a.ReqInterval, originInfo.RefTimestamp+db.GetIntervalStep(db.OneH)-1, c.Context())
	if err != nil {
		return nil, fmt.Errorf("cannot make chart to reference timestamp. name : %s, interval : %s, err : %w", originInfo.Name, a.ReqInterval, err)
	}

	if a.ReqInterval == db.OneD || a.ReqInterval == db.FourH {
		s.cutImdExeedChart(originInfo.RefTimestamp, cdd)
	}

	var oc = &OnePairChart{
		Name:         originInfo.Name,
		OneChart:     cdd,
		EntryTime:    utilities.EntryTimeFormatter(cdd.PData[0].Time),
		BtcRatio:     0,
		refTimestamp: originInfo.RefTimestamp,
		interval:     a.ReqInterval,
	}

	if a.Mode == competition {
		oc.priceFactor = originInfo.PriceFactor
		oc.timeFactor = originInfo.TimeFactor
		oc.volumeFactor = originInfo.VolumeFactor
		oc.anonymization()
	} else {
		oc.addIdentifier()
		oc.EntryPrice = oc.OneChart.PData[0].Close
	}
	return oc, nil
}

// 작은단위 뷰 -> 큰단위 요청 => 큰단위 마지막 캔들 + 1Step < 이 작은단위 뷰 마지막 캔들의 timestamp
// 큰단위 뷰 -> 작은단위 뷰 => 큰단위 마지막 캔들의 timestamp를 - 작은단위 1Step 까지 캔들은 이미 반환되고 있음.
func (s *Server) cutImdExeedChart(refTimestamp int64, cdd *CandleData) {
	var cuttingCnt int
	pdatas := cdd.PData

	intv := db.GetIntervalFromRange(pdatas[1].Time, pdatas[0].Time+1)

	step := db.GetIntervalStep(intv)
	for i := 0; pdatas[i].Time+step > refTimestamp; i++ {
		cuttingCnt++
	}
	cdd.PData = cdd.PData[cuttingCnt:]
	cdd.VData = cdd.VData[cuttingCnt:]
}
