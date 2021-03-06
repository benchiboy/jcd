package dbcomm

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db   *sql.DB
	ccdb *sql.DB
)

func InitDB(dbUrl string, ccdbUrl string, idleConns int, openConns int) {
	var err error
	db, err = sql.Open("mysql", dbUrl)
	if err != nil {
		log.Println("Open database error:", err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Println("Ping database error:", err)
		return
	}
	db.SetMaxIdleConns(idleConns)
	db.SetMaxOpenConns(openConns)
	log.Println("Database Connected successful!")

	ccdb, err = sql.Open("mysql", ccdbUrl)
	if err != nil {
		log.Println("Open database error:", err)
		return
	}
	if err = ccdb.Ping(); err != nil {
		log.Println("Ping database error:", err)
		return
	}
	ccdb.SetMaxIdleConns(idleConns)
	ccdb.SetMaxOpenConns(openConns)
	log.Println("CC Database Connected successful!")

}

func GetDB() *sql.DB {
	return db
}

func GetCCDB() *sql.DB {
	return ccdb
}
