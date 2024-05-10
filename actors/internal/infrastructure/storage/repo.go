package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/sverdejot/imdb/actors/internal/domain"
)

const BULK_BATCH_SIZE = 5_000

type MysqlActorRepository struct {
	db *sql.DB
}

func NewMysqlActorRepository(db *sql.DB) *MysqlActorRepository {
	return &MysqlActorRepository{db}
}

func (r *MysqlActorRepository) insertTitlesIds(a domain.Actor) error {
	if len(a.Titles) <= 0 {
		return nil
	}
	args := []string{}
	vals := []any{}

	for i, id := range a.Titles {
		args = append(args, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		vals = append(vals, a.Id, id)
	}


	var strs []string
	for _, id := range a.Titles {
		strs = append(strs, fmt.Sprintf("(%s)", strconv.Itoa(id)))
	}
	
	_, err :=r.db.Exec(fmt.Sprintf("INSERT INTO titles (id) VALUES %s ON CONFLICT DO NOTHING", strings.Join(strs, ",")))
	if err != nil {
		log.Fatal("while inserting title: ", err)
	}
	_, err = r.db.Exec(fmt.Sprintf("INSERT INTO actors_titles (actor_id, title_id) VALUES %s ON CONFLICT DO NOTHING", strings.Join(args, ",")), vals...)
	if err != nil {
		log.Fatal("while inserting m2m: ", err)
	}
	return nil
}

func (r *MysqlActorRepository) BulkInsert(_ context.Context, actors []domain.Actor) (int64, error) {
	if len(actors) == 0 {
		return 0, nil
	}
	args := make([]string, 0, len(actors))
	vals := make([]any, 0, len(actors)*4)

	for _, actor := range actors {
		vals = append(vals, actor.Id)
		vals = append(vals, actor.Name)
		vals = append(vals, actor.BirthYear)
		if actor.DeathYear == 0 {
			vals = append(vals, sql.NullInt32{})
		} else {
			vals = append(vals, actor.DeathYear)
		}
	}

	for idx := range BULK_BATCH_SIZE {
		args = append(args, fmt.Sprintf("($%d, $%d, $%d, $%d)", (idx*4)+1, (idx*4)+2, (idx*4+3), (idx*4)+4))
	}

	bulkStmtStr := "INSERT INTO actors (id, name, birth_year, death_year) VALUES"
	no_batches := len(actors) / BULK_BATCH_SIZE
	residual_batch := len(actors) % BULK_BATCH_SIZE

	bulkStmt, _ := r.db.Prepare(fmt.Sprintf("%s %s", bulkStmtStr, strings.Join(args, ",")))
	residualStmt, _ := r.db.Prepare(fmt.Sprintf("%s %s", bulkStmtStr, strings.Join(args[:residual_batch], ",")))

	var total_imported int64
	var lower_bound, upper_bound int
	for i := range no_batches {
		lower_bound, upper_bound = i*BULK_BATCH_SIZE*4, (i+1)*BULK_BATCH_SIZE*4
		res, err := bulkStmt.Exec(vals[lower_bound:upper_bound]...)
		if err != nil {
			return total_imported, err
		}
		rows_imported_in_batch, _ := res.RowsAffected()
		total_imported += rows_imported_in_batch
		log.Println("done inserting batch: ", i)
	}

	res, err := residualStmt.Exec(vals[(no_batches*BULK_BATCH_SIZE)*4:]...)
	if err != nil {
		return total_imported, err
	}
	rows_imported_in_residual_batch, _ := res.RowsAffected()

	log.Println("importing titles")
	var batches int
	for _, actor := range actors {
		r.insertTitlesIds(actor)
		batches++
		if batches == 50_000 {
			log.Println("done another 50.000")
			batches = 0
		}
	}

	return total_imported + rows_imported_in_residual_batch, nil
}

func (r *MysqlActorRepository) Find(_ context.Context, id int) (domain.Actor, bool) {
	actor := domain.Actor{}

	err := r.db.QueryRow(`SELECT id, name, birth_year, death_year FROM actors WHERE id=$1;`, id).Scan(
		&actor.Id,
		&actor.Name,
		&actor.BirthYear,
		&actor.DeathYear,
	)

	if err == sql.ErrNoRows {
		log.Println("no rows")
		return domain.Actor{}, false
	}


	rows, err :=r.db.Query(`SELECT title_id FROM actors_titles WHERE actor_id = $1`, actor.Id)
	if err != nil {
		log.Println(err)
		return actor, true
	}
	titles := []int{}
	for rows.Next() {
		var title int	
		rows.Scan(&title)
		titles = append(titles, title)
	}

	actor.Titles = titles

	return actor, true
}
