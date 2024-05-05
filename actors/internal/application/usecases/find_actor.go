package usecases

import "github.com/sverdejot/imdb-micro/actors/internal/domain"

type FindActor struct {
	repo domain.ActorRepository
}

func NewFindActorUseCase(repo domain.ActorRepository) *FindActor {
	return &FindActor{repo}
}

func (uc *FindActor) Execute(id int) (domain.Actor, bool) {
	return uc.repo.Find(id)
}
