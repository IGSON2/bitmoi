package db

import (
	"bitmoi/backend/utilities"
	"database/sql"
	"log"
	"os"
	"testing"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	c := utilities.GetConfig("../../../.")
	var err error
	testDB, err = sql.Open(c.DBDriver, c.DBSource)
	if err != nil {
		log.Fatal("Can't open th db : ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
