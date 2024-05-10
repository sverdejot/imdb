package usecases

import (
	"context"

	"github.com/sverdejot/imdb/actors/internal/domain"
)

type FindActor struct {
	repo domain.ActorRepository
}

func NewFindActorUseCase(repo domain.ActorRepository) *FindActor {
	return &FindActor{repo}
}

func (uc *FindActor) Execute(ctx context.Context, id int) (domain.Actor, bool) {
	return uc.repo.Find(ctx, id)
}
