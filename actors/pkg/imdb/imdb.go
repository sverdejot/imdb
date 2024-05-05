package imdb

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/sverdejot/imdb-micro/actors/internal/domain"
	"github.com/sverdejot/imdb-micro/actors/internal/infrastructure/storage"
)

func parseYear(s string) int {
	if s == "\\N" {
		return 0
	}
	y, _ := strconv.Atoi(s)
	return y
}

func ParseCsv(r io.Reader) (actors []domain.Actor, err error) {
	buffer := bufio.NewScanner(r)
	// discard header
	buffer.Scan()

	for buffer.Scan() {
		bLine := buffer.Text()
		line := strings.Split(string(bLine), "\t")

		id, err := strconv.Atoi(line[0][2:])
		if err != nil {
			return actors, err
		}

		actors = append(actors, domain.Actor{
			Id:        id,
			Name:      line[1],
			BirthYear: parseYear(line[2]),
			DeathYear: parseYear(line[3]),
		})
	}
	return actors, nil
}

func Import(path, connectionString string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	actors, err := ParseCsv(f)
	if err != nil {
		log.Println(fmt.Sprintf("error parsing actors (got %d): %v", len(actors), err))
	}

	repo := storage.NewMysqlActorRepository(db)

	rows, err := repo.BulkInsert(actors)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("total actors imported: ", rows)
}
