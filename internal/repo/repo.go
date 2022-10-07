package repo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"math"
	"ykjam/new_tx/pkg/logger"
	"ykjam/new_tx/pkg/postgres"
)

var ErrNoRowsAffected = errors.New("no rows affected")

type Repo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Repo {
	return &Repo{pg}
}

func (r *Repo) RunInTx(ctx context.Context, ops ...Operation) error {

	log := logger.Get()

	var err error
	var tx pgx.Tx
	tx, err = r.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.Deferrable,
	})
	if err != nil {
		log.Error("%s", err)
		return err
	}

	defer func() {
		if err != nil {
			err1 := tx.Rollback(ctx)
			if err1 != nil {
				log.Error("error in rollback: %s", err1)
			}
		}
	}()

	for _, op := range ops {
		err = op(ctx, tx)
		if err != nil {
			log.Error("error in op: %s", err)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error("error in commit: %s", err)
	}

	return nil
}

func newVersion(currentVersion int) (newVersion int) {
	if currentVersion > math.MaxInt8 {
		newVersion = 0
	} else {
		newVersion = currentVersion + 1
	}
	return
}
