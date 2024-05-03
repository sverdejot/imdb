package domain

type Actor struct {
	Id int
	Name string
	BirthYear int
	DeathYear int
}

func NewActor(id int, name string, birth, death int) Actor {
	return Actor{
		Id: id,
		Name: name,
		BirthYear: birth,
		DeathYear: death,
	}
}
