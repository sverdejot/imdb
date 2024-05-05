package internal

import (
	"net/http"

	"github.com/sverdejot/imdb-micro/actors/internal/application/usecases"
	"github.com/sverdejot/imdb-micro/actors/internal/infrastructure/http/routes"
)

type App struct {
	Router *http.ServeMux

	find *usecases.FindActor
}

func (a *App) AddRoutes() {
	a.Router.Handle("GET /actors/{id}", routes.NewFindActorHandler(a.find))
}

func NewApp(find *usecases.FindActor) *App {
	return &App{
		Router: http.NewServeMux(),

		find: find,
	}
}
