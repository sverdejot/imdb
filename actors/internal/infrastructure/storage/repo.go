package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/sverdejot/imdb-micro/actors/internal/domain"
)

const BATCH_SIZE = 5_000

type MysqlActorRepository struct {
	db *sql.DB
}

func NewMysqlActorRepository(db *sql.DB) *MysqlActorRepository {
	return &MysqlActorRepository{db}
}

func (r *MysqlActorRepository) BulkInsert(actors []domain.Actor) (int64, error) {
	if (len(actors) == 0) {
		return 0, nil
	}
	args := make([]string, 0, len(actors))
	vals := make([]any, 0, len(actors) * 4)

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

	for idx := range BATCH_SIZE {
		args = append(args, fmt.Sprintf("($%d, $%d, $%d, $%d)", (idx*4)+1, (idx*4)+2, (idx*4+3), (idx*4)+4))
	}

	bulkStmtStr := "INSERT INTO actors (id, name, birth_year, death_year) VALUES"
	no_batches := len(actors) / BATCH_SIZE
	residual_batch := len(actors) % BATCH_SIZE

	bulkStmt, _ := r.db.Prepare(fmt.Sprintf("%s %s", bulkStmtStr, strings.Join(args, ",")))
	residualStmt, _ := r.db.Prepare(fmt.Sprintf("%s %s", bulkStmtStr, strings.Join(args[:residual_batch], ",")))

	var total_imported int64
	for i := range no_batches {
		lower_bound, upper_bound := i*BATCH_SIZE*4, (i+1)*BATCH_SIZE*4
		res, err := bulkStmt.Exec(vals[lower_bound:upper_bound]...)
		if err != nil {
			return total_imported, err
		}
		rows_imported_in_batch, _ := res.RowsAffected()
		total_imported += rows_imported_in_batch
	}

	res, err := residualStmt.Exec(vals[(no_batches*BATCH_SIZE)*4:]...)
	if err != nil {
		return total_imported, err
	}
	rows_imported_in_residual_batch, _ := res.RowsAffected()

	return total_imported + rows_imported_in_residual_batch, nil
}
