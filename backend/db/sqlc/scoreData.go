package db

import (
	"bitmoi/backend/utilities"
	"database/sql"
	"math"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Score struct {
	User       string  `json:"user"`
	Scoreid    string  `json:"scoreid"`
	Stage      int     `json:"stage"`
	Pairname   string  `json:"pairname"`
	Entrytime  string  `json:"entrytime"`
	Position   string  `json:"position"`
	Leverage   int     `json:"leverage"`
	Outtime    int     `json:"outtime"`
	Entryprice float64 `json:"entryprice"`
	Endprice   float64 `json:"endprice"`
	Pnl        float64 `json:"pnl"`
	Roe        float64 `json:"roe"`
	Pkey       string  `json:"-"`
}

type ScoreList struct {
	ScoreList []Score `json:"scorelist"`
}

type TotalData struct {
	User        string  `json:"user"`
	Displayname string  `json:"displayname"`
	PhotoUrl    string  `json:"photourl"`
	Scoreid     string  `json:"scoreid"`
	Balance     float64 `json:"balance"`
}

type TotalList struct {
	TotalList []TotalData `json:"totallist"`
}

type MoreInfoData struct {
	Comment    string        `json:"comment"`
	Scoreid    string        `json:"scoreid"`
	AvgLev     float64       `json:"avglev"`
	AvgPnl     float64       `json:"avgpnl"`
	AvgRoe     float64       `json:"avgroe"`
	StageArray []StageSimple `json:"stagearray"`
}

type StageSimple struct {
	Name string  `json:"name"`
	Date string  `json:"date"`
	Roe  float64 `json:"roe"`
}

var (
	userScoreDB *sql.DB
)

const (
	scoreDbPath    string = "./backend/db/scoreData/Score.db"
	stageTableName string = "userscore"
	totalTableName string = "totalscore"

	createStageTableQuery = `create table IF NOT EXISTS ` + stageTableName +
		` (
		user text NOT NULL,
		scoreid text,
		stage int,
		pairname text,
		entrytime text,
		position text,
		leverage integer,
		outtime integer,
		entryprice real,
		endprice real,
		pnl real,
		roe real,
		pkey string PRIMARY KEY
		)`

	createTotalTableQuery = `create table IF NOT EXISTS ` + totalTableName +
		` (
			user text NOT NULL,
			displayname text,
			photourl text,
			scoreid text,
			balance real,
			comment text
			)`

	insertStageQuery = `insert into ` + stageTableName +
		` (
		user,
		scoreid,
		stage,
		pairname,
		entrytime,
		position,
		leverage,
		outtime,
		entryprice,
		endprice,
		pnl,
		roe,
		pkey
	) values(?,?,?,?,?,?,?,?,?,?,?,?,?)`

	insertTotalQuery = `insert into ` + totalTableName +
		` (
		user,
		displayname,
		photourl,
		scoreid,
		balance,
		comment
	) values(?,?,?,?,?,?)`
)

func initStageScoreDB() *sql.DB {
	db, err := sql.Open("sqlite3", scoreDbPath)
	if userScoreDB == nil {
		utilities.Errchk(err)
		db.SetMaxOpenConns(1)
		files, err := filepath.Glob(scoreDbPath)
		utilities.Errchk(err)
		if len(files) == 0 {
			f, err := os.Create(scoreDbPath)
			utilities.Errchk(err)
			defer f.Close()
		}
		utilities.Errchk(err)
		_, err = db.Exec(createStageTableQuery)
		utilities.Errchk(err)
		_, err = db.Exec(createTotalTableQuery)
		utilities.Errchk(err)
	}
	userScoreDB = db
	return userScoreDB
}

func InsertStageScore(order OrderStruct, result ResultScore) {
	db := initStageScoreDB()
	defer db.Close()
	tx, err := db.Begin()
	utilities.Errchk(err)
	stmt, err := tx.Prepare(insertStageQuery)
	utilities.Errchk(err)
	var position string
	if order.IsLong {
		position = "LONG"
	} else {
		position = "SHORT"
	}
	var pkey string = order.Uid + order.ScoreId + result.Name
	if result.Isliquidated {
		result.Name += " (Liq.)"
	}
	_, err = stmt.Exec(order.Uid, order.ScoreId, order.Stage, result.Name, result.Entrytime, position, order.Leverage, result.OutHour, result.EntryPrice, result.EndPrice, result.Pnl, result.Roe, pkey)
	utilities.Errchk(err)
	tx.Commit()
	stmt.Close()
}

func SelectStageScoreDB(user string, limit int) ScoreList {
	db := initStageScoreDB()
	defer db.Close()
	stmt, err := db.Prepare(`select * from ` + stageTableName + " where user=? order by scoreid desc limit ?,?")
	utilities.Errchk(err)
	defer stmt.Close()
	rows, err := stmt.Query(user, (limit-1)*15, 15)
	utilities.Errchk(err)
	defer rows.Close()
	var sList ScoreList
	var s Score
	for rows.Next() {
		rows.Scan(&s.User, &s.Scoreid, &s.Stage, &s.Pairname, &s.Entrytime, &s.Position, &s.Leverage, &s.Outtime, &s.Entryprice, &s.Endprice, &s.Pnl, &s.Roe, &s.Pkey)
		sList.ScoreList = append(sList.ScoreList, s)
	}
	return sList
}

func InsertTotalScore(t TotalData) {
	db := initStageScoreDB()
	defer db.Close()
	stmt, err := db.Prepare(insertTotalQuery)
	utilities.Errchk(err)
	_, err = stmt.Exec(t.User, t.Displayname, t.PhotoUrl, t.Scoreid, math.Floor(t.Balance*100)/100, "")
	stmt.Close()
	utilities.Errchk(err)

	stmt3, err := db.Prepare(`select balance,scoreid from ` + totalTableName + ` where user=? order by balance desc limit 1`)
	utilities.Errchk(err)
	row := stmt3.QueryRow(t.User)
	utilities.Errchk(err)
	var b float64
	var s string
	row.Scan(&b, &s)

	stmt2, err := db.Prepare(`delete from ` + totalTableName + ` where user=? and balance <?`)
	utilities.Errchk(err)
	_, err = stmt2.Exec(t.User, b)
	stmt2.Close()
	utilities.Errchk(err)

	deleteMyscoreOver200(db, t.User, s)
}

func SelectTotalScoreDB() TotalList {
	db := initStageScoreDB()
	defer db.Close()
	rows, err := db.Query(`select user,displayname,photourl,scoreid,balance from ` + totalTableName + ` order by balance desc limit 30`)
	utilities.Errchk(err)
	defer rows.Close()
	var tList TotalList
	var b TotalData
	for rows.Next() {
		rows.Scan(&b.User, &b.Displayname, &b.PhotoUrl, &b.Scoreid, &b.Balance)
		tList.TotalList = append(tList.TotalList, b)
	}
	return tList
}

func deleteMyscoreOver200(db *sql.DB, user, scoreid string) {

	var idx int
	stmt, err := db.Prepare(`select count(scoreid) from userscore where user=?`)
	utilities.Errchk(err)
	row := stmt.QueryRow(user)
	utilities.Errchk(err)
	stmt.Close()
	row.Scan(&idx)

	if idx > 200 {
		stmt2, err := db.Prepare(`delete from userscore where scoreid in (select scoreid from userscore where scoreid != ? order by scoreid asc)`)
		utilities.Errchk(err)
		_, err = stmt2.Exec(scoreid)
		utilities.Errchk(err)
	}
}

func SendMoreInfo(user, scoreid string) MoreInfoData {
	db := initStageScoreDB()
	defer db.Close()
	stmt, err := db.Prepare(`select * from userscore where user=? and scoreid=?`)
	utilities.Errchk(err)
	rows, err := stmt.Query(user, scoreid)
	utilities.Errchk(err)
	stmt.Close()
	var m MoreInfoData
	var s Score
	for rows.Next() {
		rows.Scan(&s.User, &s.Scoreid, &s.Stage, &s.Pairname, &s.Entrytime, &s.Position, &s.Leverage, &s.Outtime, &s.Entryprice, &s.Endprice, &s.Pnl, &s.Roe, &s.Pkey)
		m.AvgPnl += s.Pnl
		m.AvgRoe += s.Roe
		m.AvgLev += float64(s.Leverage)
		m.StageArray = append(m.StageArray, StageSimple{s.Pairname, s.Entrytime, s.Roe})
	}
	rows.Close()
	m.Scoreid = s.Scoreid
	m.AvgPnl = math.Floor((m.AvgPnl/float64(len(m.StageArray)))*100) / 100
	m.AvgRoe = math.Floor((m.AvgRoe/float64(len(m.StageArray)))*100) / 100
	m.AvgLev = math.Floor((m.AvgLev/float64(len(m.StageArray)))*10) / 10

	stmt2, err := db.Prepare(`select comment from totalscore where user=?`)
	utilities.Errchk(err)
	row := stmt2.QueryRow(user)
	row.Scan(&m.Comment)
	stmt2.Close()

	return m
}

func UpdateComment(comment, user string) error {
	db := initStageScoreDB()
	defer db.Close()
	stmt, err := db.Prepare("update totalscore set comment = ? where user = ? ")
	utilities.Errchk(err)
	_, err = stmt.Exec(comment, user)
	stmt.Close()

	return err
}
