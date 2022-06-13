package apiWebsocket

import (
	db "bitmoi/backend/db/chartData"
)

const (
	competitionMode int = 1
	practiceMode    int = 2
)

type Message struct {
	Message string `json:"message"`
}

type SendJson struct {
	Charts []db.OnePairChart `json:"charts"`
	Temp   string            `json:"temp"`
}

// func HandleMessage(m Message, mCH *MessageCh) {
// 	mCH.m.Lock()
// 	defer mCH.m.Unlock()
// 	switch m.Message {
// 	case "competition":
// 		utilities.Errchk(mCH.conn.WriteJSON(db.SendCharts(competitionMode)))
// 	case "practice":
// 		utilities.Errchk(mCH.conn.WriteJSON(db.SendCharts(practiceMode)))

// 	}
// }
