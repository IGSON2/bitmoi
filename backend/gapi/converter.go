package gapi

import "bitmoi/backend/gapi/pb"

func convertGetCandlesRes(o *OnePairChart) *pb.GetCandlesResponse {
	return &pb.GetCandlesResponse{
		Name:       o.Name,
		Onechart:   o.OneChart,
		BtcRatio:   o.BtcRatio,
		EntryTime:  o.EntryTime,
		EntryPrice: o.EntryPrice,
		Identifier: o.Identifier,
	}
}
