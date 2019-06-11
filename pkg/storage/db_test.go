package storage

import "testing"

var (
	//DBUrl is the url to connect to the influx DB
	DBUrl = "http://localhost:8086"
	//StockDB is the main database to hold ticks information
	StockDB = "stockist"
)

func getDB() *DB {
	db := NewDB(DBUrl, StockDB, "")
	return db
}
func TestGetMaxHigh(t *testing.T) {
	db := getDB()
	db.GetMaxHigh()
	t.Error("hi")

}

func TestGetLowest(t *testing.T) {
	db := getDB()
	db.GetLowest()
	t.Error("hi")

}

func TestGetOpen(t *testing.T) {
	db := getDB()
	db.GetMarketOpenPrice("2019-06-07T09:11:00Z")

}
