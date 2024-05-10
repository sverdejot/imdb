package domain

import "context"

type ActorRepository interface {
	BulkInsert(ctx context.Context, actors []Actor) (int64, error)
	Find(ctx context.Context, id int) (Actor, bool)
}
