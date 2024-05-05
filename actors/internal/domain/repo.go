package domain

type ActorRepository interface {
	BulkInsert(actors []Actor) (int64, error)
	Find(id int) (Actor, bool)
}
