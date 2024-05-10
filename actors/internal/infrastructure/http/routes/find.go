package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sverdejot/imdb/actors/internal/application/usecases"
)

func NewFindActorHandler(uc *usecases.FindActor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("woops, bad request")
			return
		}

		actor, found := uc.Execute(ctx, id)
		if !found {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("not found")
		}

		var actorDto struct {
			Id    int    `json:"id"`
			Name  string `json:"name"`
			Dates struct {
				Birth int `json:"birth"`
				Death int `json:"death"`
			} `json:"dates"`
			Titles []int `json:"titles"`
		}

		actorDto.Id = actor.Id
		actorDto.Name = actor.Name
		actorDto.Dates.Birth = actor.BirthYear
		actorDto.Dates.Death = actor.DeathYear
		actorDto.Titles = actor.Titles

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(actorDto)
	}
}
