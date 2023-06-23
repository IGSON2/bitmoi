package gapi

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/utilities"
	"context"
	"encoding/json"
	"fmt"
)

func (s *Server) sendAnotherInterval(a *pb.AnotherIntervalRequest, c context.Context) (*OnePairChart, error) {
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

	err = s.cutExceedChart(originInfo.Name, originInfo.RefTimestamp, c, cdd)
	if err != nil {
		return nil, fmt.Errorf("cannot cut exceed chart. name : %s, interval : %s, err : %w", originInfo.Name, a.ReqInterval, err)
	}

	var oc = &OnePairChart{
		Name:         originInfo.Name,
		OneChart:     cdd,
		EntryTime:    utilities.EntryTimeFormatter(cdd.PData[0].Time),
		BtcRatio:     ratio,
		refTimestamp: originInfo.RefTimestamp,
		interval:     a.ReqInterval,
	}

	if a.Mode == competition {
		oc.priceFactor = originInfo.PriceFactor
		oc.timeFactor = originInfo.TimeFactor
		oc.volumeFactor = originInfo.VolumeFactor
		oc.anonymization(int(a.Stage) - 1)
	} else {
		oc.addIdentifier()
		oc.EntryPrice = oc.OneChart.PData[0].Close
	}
	return oc, nil
}

// 최초 전송 차트의 interval이 1h이므로 이보다 작은 단위의 another interval 요청은
// reqTimestamp에 더 가까이 접근하여 entrytime보다 미래의 캔들을 불러온다. cutExceedChart은 이를 방지한다.
func (s *Server) cutExceedChart(name string, originTimestamp int64, c context.Context, cd *pb.CandleData) error {
	var cuttingCnt int
	pdatas := cd.PData
	entryTIme, err := s.store.Get1hEntryTimestamp(c, db.Get1hEntryTimestampParams{Name: name, Time: originTimestamp})
	if err != nil {
		return err
	}

	for i := 0; pdatas[i].Time > entryTIme; i++ {
		cuttingCnt++
	}
	cd.PData = cd.PData[cuttingCnt:]
	cd.VData = cd.VData[cuttingCnt:]
	return nil
}
