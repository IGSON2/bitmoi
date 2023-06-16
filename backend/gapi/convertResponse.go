package gapi

import "bitmoi/backend/gapi/pb"

func convertGetCandlesRes(o *OnePairChart) *pb.GetCandlesResponse {
	return &pb.GetCandlesResponse{
		Name:       o.Name,
		Onechart:   o.OneChart,
		Btcratio:   o.BtcRatio,
		Entrytime:  o.EntryTime,
		EntryPrice: o.EntryPrice,
		Identifier: o.Identifier,
	}
}
