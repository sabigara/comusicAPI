package main

import (
	"context"
	"os"
	"strings"
	"time"

	fb "firebase.google.com/go"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sabigara/comusicAPI/firebase"
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

func main() {
	addr := "0.0.0.0:1323"
	var debug bool
	if val := os.Getenv("DEBUG"); val == "true" {
		debug = true
	} else {
		debug = false
	}

	fbapp, err := fb.NewApp(context.Background(), nil)
	if err != nil {
		panic("Cannot initiate firebase App")
	}
	authUsecase := firebase.NewAuthUsecase(fbapp)

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
	fileRepository := mock.NewFileRepository()
	verUsecase := interactor.NewVersionUsecase(verRepository, fileRepository)
	verHandler := http.NewVersionHandler(verUsecase)

	trackRepository := mysql.NewTrackRepository(db)
	trackUsecase := &interactor.TrackUsecase{TrackRepository: trackRepository}
	// Skip creating TrackHandler, because TrackUsecase needs TakeUsecase
	// which also refers to TrackUsecase.

	takeRepository := mysql.NewTakeRepository(db)
	takeUsecase := interactor.NewTakeUsecase(
		trackUsecase,
		takeRepository,
		fileRepository,
	)
	takeHandler := http.NewTakeHandler(takeUsecase, fileRepository)

	// Inject TakeUsecase to TrackUsecase here to avoid circular reference.
	trackUsecase.TakeUsecase = takeUsecase
	trackHandler := http.NewTrackHandler(trackUsecase)

	inviteRepository := mysql.NewInvitationRepository(db)
	userRepository := firebase.NewUserRepository(fbapp)
	mailUsecase := mock.NewMailUsecase()
	inviteUsecase := interactor.NewInvitationUsecase(
		inviteRepository,
		userRepository,
		studioUsecase,
		songUsecase,
		mailUsecase,
	)
	inviteHandler := http.NewInvitationHandler(inviteUsecase)

	hooks := http.NewHooks(profileUsecase, studioUsecase)

	http.New(
		authUsecase,
		profileHandler,
		studioHandler,
		songHandler,
		verHandler,
		trackHandler,
		takeHandler,
		inviteHandler,
		hooks,
	).Start(addr, debug)
}
