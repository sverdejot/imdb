package bootstrap

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/sverdejot/imdb-micro/actors/internal"
	"github.com/sverdejot/imdb-micro/actors/internal/application/usecases"
	"github.com/sverdejot/imdb-micro/actors/internal/infrastructure/storage"
)

func Run() {
	db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/imdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	repo := storage.NewMysqlActorRepository(db)
	find := usecases.NewFindActorUseCase(repo)

	app := internal.NewApp(find)
	app.AddRoutes()

	srv := http.Server{
		Addr: ":8080",

		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,

		Handler: app.Router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
