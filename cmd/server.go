package main

import (
	"database/sql"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sabigara/comusicAPI/http"
)

func openDB() *sql.DB {
	DSN, ok := os.LookupEnv("DSN")
	if !ok {
		panic("No DSN provided as environment variable.")
	}
	dsn := strings.Split(DSN, "://")
	if len(dsn) != 2 {
		panic("Malformed DSN.")
	}
	db, err := sql.Open(dsn[0], dsn[1])
	if err != nil {
		panic(err.Error)
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)
	return db
}

func inject() {

}

func main() {
	addr := "0.0.0.0:1323"
	var debug bool
	if val := os.Getenv("DEBUG"); val == "true" {
		debug = true
	} else {
		debug = false
	}
	inject()
	http.Start(addr, debug)
}
