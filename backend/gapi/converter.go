package gapi

import "bitmoi/backend/gapi/pb"

func convertGetCandlesRes(o *OnePairChart) *pb.CandlesResponse {
	return &pb.CandlesResponse{
		Name:       o.Name,
		OneChart:   o.OneChart,
		BtcRatio:   o.BtcRatio,
		EntryTime:  o.EntryTime,
		EntryPrice: o.EntryPrice,
		Identifier: o.Identifier,
	}
}
