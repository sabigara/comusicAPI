package main

import (
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sabigara/comusicAPI/http"
	"github.com/sabigara/comusicAPI/interactor"
	"github.com/sabigara/comusicAPI/mock"
	"github.com/sabigara/comusicAPI/mysql"
)

func openDB() *sqlx.DB {
	DSN, ok := os.LookupEnv("DSN")
	if !ok {
		panic("No DSN provided as environment variable.")
	}
	dsn := strings.Split(DSN, "://")
	if len(dsn) != 2 {
		panic("Malformed DSN.")
	}
	db := sqlx.MustConnect(dsn[0], dsn[1])

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)
	return db
}

func inject() {
	db := openDB()
	profileRepository := mysql.NewProfileRepository(db)
	profileUsecase := interactor.NewProfileUsecase(profileRepository)
	profileHandler := http.NewProfileHandler(profileUsecase)

	studioRepository := mysql.NewStudioRepository(db)
	studioUsecase := interactor.NewStudioUsecase(studioRepository)
	studioHandler := http.NewStudioHandler(studioUsecase)

	songRepository := mysql.NewSongRepository(db)
	songUsecase := interactor.NewSongUsecase(songRepository)
	songHandler := http.NewSongHandler(songUsecase)

	verRepository := mysql.NewVersionRepository(db)
	verUsecase := interactor.NewVersionUsecase(verRepository)
	verHandler := http.NewVersionHandler(verUsecase)

	http.SetHandlers(
		profileHandler,
		studioHandler,
		songHandler,
		verHandler,
	)
	http.SetAuthenticate(mock.Aunthenticate)
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
