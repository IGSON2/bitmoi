package api

import (
	"bitmoi/backend/utilities"
	"bitmoi/backend/utilities/common"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) sendAnotherInterval(a *AnotherIntervalRequest, c *fiber.Ctx) (*OnePairChart, error) {
	originInfo := new(utilities.IdentificationData)
	infoByte := utilities.DecryptByASE(a.Identifier)
	err := json.Unmarshal(infoByte, originInfo)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal chart identifier. err : %w", err)
	}
	cdd, err := s.selectStageChart(originInfo.Name, a.ReqInterval, originInfo.RefTimestamp, c)
	if err != nil {
		return nil, fmt.Errorf("cannot make chart to reference timestamp. name : %s, interval : %s, err : %w", originInfo.Name, a.ReqInterval, err)
	}

	ratio, err := s.calcBtcRatio(a.ReqInterval, originInfo.Name, originInfo.RefTimestamp, c)
	if err != nil {
		return nil, fmt.Errorf("cannot calculate btc ratio. name : %s, interval : %s, refTime : %d, err : %w",
			originInfo.Name, a.ReqInterval, originInfo.RefTimestamp, err)
	}

	var oc = &OnePairChart{
		Name:         originInfo.Name,
		OneChart:     cdd,
		EntryTime:    utilities.EntryTimeFormatter(cdd.PData[0].Time),
		BtcRatio:     common.CeilDecimal(ratio) * 100,
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
