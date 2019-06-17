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
	db.Measurement = "ticks_1076225_5m"
	f, _ := db.GetLowestLow(1)
	t.Error(f)

}

// func TestGetOpen(t *testing.T) {
// 	db := getDB()
// 	db.GetMarketOpenPrice("2019-06-07T09:11:00Z")

// }

func TestInsertTrade(t *testing.T) {
	db := getDB()
	db.Measurement = "trade"
	db.InsertTrade("12345", "BUY")
	//t.Error("hi")

}

func TestGetLastTrade(t *testing.T) {
	db := getDB()
	res, _ := db.GetLastTrade("12345")
	t.Error(res)

}

func TestGetPointsCount(t *testing.T) {
	db := getDB()
	db.Measurement = "ticks_1076225_5m"
	res, _ := db.GetPointsCount()
	t.Error(res)

}
