package domain

type ActorRepository interface {
	BulkInsert([]Actor) (int, error)
}
