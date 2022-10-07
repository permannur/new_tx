package repo

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"ykjam/new_tx/internal/entity"
	"ykjam/new_tx/pkg/logger"
)

const (
	sqlDepartmentAdd         = `INSERT INTO tbl_department(name, state, priority) VALUES($1, $2, $3) RETURNING id`
	sqlDepartmentUpdate      = `UPDATE tbl_department SET name=$3, priority=$4 WHERE id=$1 AND state!=$2`
	sqlDepartmentChangeState = `UPDATE tbl_department SET state=$3 WHERE id=$1 AND state!=$2`
	sqlDepartmentById        = `SELECT id, name, state, priority FROM tbl_department WHERE id=$1 and state!=$2`
	sqlDepartmentList        = `SELECT id, name, state, priority FROM tbl_department WHERE state!=$1 ORDER BY priority ASC`
)

func (r *Repo) DepartmentAdd(name string, priority int) (*entity.Department, Operation) {

	var item *entity.Department

	opt := func(ctx context.Context, tx pgx.Tx) error {

		log := logger.Get()
		var err error

		item = &entity.Department{
			Name:     name,
			State:    entity.StateEnabled,
			Priority: priority,
		}

		defer func() {
			if err != nil {
				item = nil
			}
		}()

		//sqlDepartmentAdd = `INSERT INTO tbl_department(name, state, priority) VALUES($1, $2, $3) RETURNING id`
		row := tx.QueryRow(ctx, sqlDepartmentAdd, item.Name, item.State, item.Priority)
		err = row.Scan(&item.Id)
		if err != nil {
			log.Error("Error at sqlDepartmentAdd: %s", err)
			return err
		}

		return nil
	}

	return item, opt
}

func (r *Repo) DepartmentUpdate(item *entity.Department, name string, priority int) Operation {

	return func(ctx context.Context, tx pgx.Tx) error {

		log := logger.Get()
		var err error

		//sqlDepartmentUpdate = `UPDATE tbl_department SET name=$3, priority=$4 WHERE id=$1 AND state!=$2`
		var cmdTag pgconn.CommandTag
		cmdTag, err = tx.Exec(ctx, sqlDepartmentUpdate, item.Id, entity.StateDeleted, name, priority)
		if err != nil {
			log.Error("error at sqlDepartmentUpdate: %s", err)
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

func (r *Repo) DepartmentChangeState(item *entity.Department, state entity.State) Operation {

	return func(ctx context.Context, tx pgx.Tx) error {

		log := logger.Get()
		var err error

		//sqlDepartmentChangeState = `UPDATE tbl_department SET state=$3 WHERE id=$1 AND state!=$2`
		var cmdTag pgconn.CommandTag
		cmdTag, err = tx.Exec(ctx, sqlDepartmentChangeState, item.Id, entity.StateDeleted, state)
		if err != nil {
			log.Error("error at sqlDepartmentChangeState: %s", err)
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

func (r *Repo) DepartmentById(ctx context.Context, id int) (*entity.Department, error) {

	log := logger.Get()
	var err error
	item := &entity.Department{}

	//sqlDepartmentById = `SELECT id, name, state, priority FROM tbl_department WHERE id=$1 and state!=$2`
	row := r.Pool.QueryRow(ctx, sqlDepartmentById, id, entity.StateDeleted)
	err = row.Scan(&item.Id, &item.Name, &item.State, &item.Priority)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("No Result")
			return nil, nil
		}
		log.Error("error in sqlDepartmentById: %s", err)
		return nil, err
	}
	return item, nil
}

func (r *Repo) DepartmentList(ctx context.Context) ([]*entity.Department, error) {

	log := logger.Get()
	var err error
	items := make([]*entity.Department, 0)

	//sqlDepartmentList = `SELECT id, name, state, priority FROM tbl_department WHERE state!=$1 ORDER BY priority ASC`
	var rows pgx.Rows
	rows, err = r.Pool.Query(ctx, sqlDepartmentList, entity.StateDeleted)
	if err != nil {
		log.Error("error in sqlDepartmentList: %s", err)
		return nil, err
	}
	for rows.Next() {
		item := &entity.Department{}
		err = rows.Scan(&item.Id, &item.Name, &item.State, &item.Priority)
		if err != nil {
			log.Error("Error at rows.Scan: %s", err)
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
