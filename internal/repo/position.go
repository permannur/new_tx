package repo

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"ykjam/new_tx/internal/entity"
	"ykjam/new_tx/pkg/logger"
)

const (
	sqlPositionAdd         = `INSERT INTO tbl_position(name, state, priority) VALUES($1, $2, $3) RETURNING id`
	sqlPositionUpdate      = `UPDATE tbl_position SET name=$3, priority=$4 WHERE id=$1 AND state!=$2`
	sqlPositionChangeState = `UPDATE tbl_position SET state=$3 WHERE id=$1 AND state!=$2`
	sqlPositionById        = `SELECT id, name, state, priority FROM tbl_position WHERE id=$1 and state!=$2`
	sqlPositionList        = `SELECT id, name, state, priority FROM tbl_position WHERE state!=$1 ORDER BY priority ASC`
)

func (r *Repo) PositionAdd(name string, priority int) (*entity.Position, Operation) {

	var item *entity.Position

	opt := func(ctx context.Context, tx pgx.Tx) error {

		log := logger.Get()
		var err error

		item = &entity.Position{
			Name:     name,
			State:    entity.StateEnabled,
			Priority: priority,
		}

		defer func() {
			if err != nil {
				item = nil
			}
		}()

		//sqlPositionAdd = `INSERT INTO tbl_position(name, state, priority) VALUES($1, $2, $3) RETURNING id`
		row := tx.QueryRow(ctx, sqlPositionAdd, item.Name, item.State, item.Priority)
		err = row.Scan(&item.Id)
		if err != nil {
			log.Error("Error at sqlPositionAdd: %s", err)
			return err
		}

		return nil
	}

	return item, opt
}

func (r *Repo) PositionUpdate(item *entity.Position, name string, priority int) Operation {

	return func(ctx context.Context, tx pgx.Tx) error {

		log := logger.Get()
		var err error

		//sqlPositionUpdate = `UPDATE tbl_position SET name=$3, priority=$4 WHERE id=$1 AND state!=$2`
		var cmdTag pgconn.CommandTag
		cmdTag, err = tx.Exec(ctx, sqlPositionUpdate, item.Id, entity.StateDeleted, name, priority)
		if err != nil {
			log.Error("error at sqlPositionUpdate: %s", err)
			return err
		}
		if cmdTag.RowsAffected() == 0 {
			err = ErrNoRowsAffected
			log.Error("no rows affected during update: %s", err)
			return err
		}

		item.Name = name
		item.Priority = priority

		return nil
	}
}

func (r *Repo) PositionChangeState(item *entity.Position, state entity.State) Operation {

	return func(ctx context.Context, tx pgx.Tx) error {

		log := logger.Get()
		var err error

		//sqlPositionChangeState = `UPDATE tbl_position SET state=$3 WHERE id=$1 AND state!=$2`
		var cmdTag pgconn.CommandTag
		cmdTag, err = tx.Exec(ctx, sqlPositionChangeState, item.Id, entity.StateDeleted, state)
		if err != nil {
			log.Error("error at sqlPositionChangeState: %s", err)
			return err
		}
		if cmdTag.RowsAffected() == 0 {
			err = ErrNoRowsAffected
			log.Error("no rows affected during change state: %s", err)
			return err
		}

		item.State = state

		return nil
	}
}

func (r *Repo) PositionById(ctx context.Context, id int) (*entity.Position, error) {

	log := logger.Get()
	var err error
	item := &entity.Position{}

	//sqlPositionById = `SELECT id, name, state, priority FROM tbl_position WHERE id=$1 and state!=$2`
	row := r.Pool.QueryRow(ctx, sqlPositionById, id, entity.StateDeleted)
	err = row.Scan(&item.Id, &item.Name, &item.State, &item.Priority)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("No Result")
			return nil, nil
		}
		log.Error("error in sqlPositionById: %s", err)
		return nil, err
	}
	return item, nil
}

func (r *Repo) PositionList(ctx context.Context) ([]*entity.Position, error) {

	log := logger.Get()
	var err error
	items := make([]*entity.Position, 0)

	//sqlPositionList = `SELECT id, name, state, priority FROM tbl_position WHERE state!=$1 ORDER BY priority ASC`
	var rows pgx.Rows
	rows, err = r.Pool.Query(ctx, sqlPositionList, entity.StateDeleted)
	if err != nil {
		log.Error("error in sqlPositionList: %s", err)
		return nil, err
	}
	for rows.Next() {
		item := &entity.Position{}
		err = rows.Scan(&item.Id, &item.Name, &item.State, &item.Priority)
		if err != nil {
			log.Error("Error at rows.Scan: %s", err)
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
