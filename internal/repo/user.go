package repo

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"time"
	"ykjam/new_tx/internal/entity"
	"ykjam/new_tx/pkg/logger"
)

const (
	sqlUserAdd = `INSERT INTO tbl_user(username, password, firstname, lastname, department_id, position_id, state, 
                     create_ts, update_ts, version) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	sqlUserUpdate = `UPDATE tbl_user SET firstname=$3, lastname=$4, department_id=$5, position_id=$6, update_ts=$7, 
                    version=$8  WHERE id=$1 AND version=$2`

	sqlUserChangeState = `UPDATE tbl_user SET state=$3, update_ts=$4, version=$5 WHERE id=$1 AND version=$2`

	sqlUserById = `SELECT id, username, firstname, lastname, department_id, position_id, state, create_ts, update_ts, 
       version FROM tbl_user WHERE id=$1 AND state!=$2`

	sqlUserList = `SELECT a.id, username, firstname, lastname, department_id, b.name, position_id, c.name, 
       create_ts, update_ts, version FROM tbl_user a 
           INNER JOIN tbl_department b ON b.id=a.department_id
           INNER JOIN tbl_position c ON c.id=a.position_id  
               WHERE a.state!=$1 AND a.state!=$2 ORDER BY c.priority, a.username`
)

func (r *Repo) UserAdd(username, password, firstname, lastname string, departmentId, positionId int,
	state entity.UserState) (*entity.User, Operation) {

	var item *entity.User

	opt := func(ctx context.Context, tx pgx.Tx) error {

		log := logger.Get()
		var err error
		now := time.Now().UTC().Round(time.Microsecond)

		item = &entity.User{
			Username:     username,
			Password:     password,
			Firstname:    firstname,
			Lastname:     lastname,
			DepartmentId: departmentId,
			PositionId:   positionId,
			State:        state,
			CreateTs:     now,
			UpdateTs:     now,
			Version:      0,
		}

		defer func() {
			if err != nil {
				item = nil
			}
		}()

		//sqlUserAdd = `INSERT INTO tbl_user(username, password, firstname, lastname, department_id, position_id, state,
		//             create_ts, update_ts, version) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
		row := tx.QueryRow(ctx, sqlUserAdd, item.Username, item.Password, item.Firstname, item.Lastname,
			item.DepartmentId, item.PositionId, item.State, item.CreateTs, item.UpdateTs, item.Version)
		err = row.Scan(&item.Id)
		if err != nil {
			log.Error("Error at sqlUserAdd: %s", err)
			return err
		}

		return nil
	}

	return item, opt
}

func (r *Repo) UserUpdate(item *entity.User, firstname, lastname string, depId, posId int) Operation {

	return func(ctx context.Context, tx pgx.Tx) error {

		log := logger.Get()
		var err error
		now := time.Now().UTC().Round(time.Microsecond)
		nv := newVersion(item.Version)

		//sqlUserUpdate = `UPDATE tbl_user SET firstname=$3, lastname=$4, department_id=$5, position_id=$6, update_ts=$7,
		//            version=$8  WHERE id=$1 AND version=$2`
		var cmdTag pgconn.CommandTag
		cmdTag, err = tx.Exec(ctx, sqlUserUpdate, item.Id, item.Version, firstname, lastname, depId, posId, now, nv)
		if err != nil {
			log.Error("error at sqlUserUpdate: %s", err)
			return err
		}
		if cmdTag.RowsAffected() == 0 {
			err = ErrNoRowsAffected
			log.Error("no rows affected during update: %s", err)
			return err
		}

		item.Firstname = firstname
		item.Lastname = lastname
		item.DepartmentId = depId
		item.PositionId = posId
		item.UpdateTs = now
		item.Version = nv

		return nil
	}
}

func (r *Repo) UserChangeState(item *entity.User, state entity.UserState) Operation {

	return func(ctx context.Context, tx pgx.Tx) error {

		log := logger.Get()
		var err error
		now := time.Now().UTC().Round(time.Microsecond)
		nv := newVersion(item.Version)

		//sqlUserChangeState = `UPDATE tbl_user SET state=$3, update_ts=$4, version=$5 WHERE id=$1 AND version=$2`
		var cmdTag pgconn.CommandTag
		cmdTag, err = tx.Exec(ctx, sqlUserChangeState, item.Id, item.Version, state, now, nv)
		if err != nil {
			log.Error("error at sqlUserChangeState: %s", err)
			return err
		}
		if cmdTag.RowsAffected() == 0 {
			err = ErrNoRowsAffected
			log.Error("no rows affected during change state: %s", err)
			return err
		}

		item.State = state
		item.UpdateTs = now
		item.Version = nv

		return nil
	}
}

func (r *Repo) UserById(ctx context.Context, id int) (*entity.User, error) {

	log := logger.Get()
	var err error
	item := &entity.User{}

	//sqlUserById = `SELECT id, username, firstname, lastname, department_id, position_id, state, create_ts, update_ts,
	//   version FROM tbl_user WHERE id=$1 AND state!=$2`
	row := r.Pool.QueryRow(ctx, sqlUserById, id, entity.UserStateDeleted)
	err = row.Scan(&item.Id, &item.Username, &item.Firstname, &item.Lastname, &item.DepartmentId, &item.PositionId,
		&item.State, &item.CreateTs, &item.UpdateTs, &item.Version)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug("No Result")
			return nil, nil
		}
		log.Error("error in sqlUserById: %s", err)
		return nil, err
	}
	return item, nil
}

func (r *Repo) UserList(ctx context.Context) ([]*entity.User, error) {

	log := logger.Get()
	var err error
	items := make([]*entity.User, 0)

	//sqlUserList = `SELECT a.id, username, firstname, lastname, department_id, b.name, position_id, c.name,
	//   create_ts, update_ts, version FROM tbl_user a
	//       INNER JOIN tbl_department b ON b.id=a.department_id
	//       INNER JOIN tbl_position c ON c.id=a.position_id
	//           WHERE a.state!=$1 AND a.state!=$2 ORDER BY c.priority, a.username`
	var rows pgx.Rows
	rows, err = r.Pool.Query(ctx, sqlUserList, entity.UserStateDeleted, entity.UserStateBlocked)
	if err != nil {
		log.Error("error in sqlUserList: %s", err)
		return nil, err
	}
	for rows.Next() {
		item := &entity.User{}
		err = rows.Scan(&item.Id, &item.Username, &item.Firstname, &item.Lastname, &item.DepartmentId,
			&item.Department, &item.PositionId, &item.Position, &item.CreateTs, &item.UpdateTs, &item.Version)
		if err != nil {
			log.Error("Error at rows.Scan: %s", err)
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
