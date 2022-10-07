package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"time"
	"ykjam/new_tx/internal/entity"
	"ykjam/new_tx/pkg/logger"
)

const (
	sqlUserLogAdd = `INSERT INTO tbl_user_log(id, user_id, username, ip, action, action_ts, sup_info) 
					VALUES($1, $2, $3, $4, $5, $6, $7)`
)

func (r *Repo) UserLogAdd(ai *entity.ActionInfo, action entity.UserAction) Operation {
	var item *entity.UserLog

	opt := func(ctx context.Context, tx pgx.Tx) error {

		log := logger.Get()
		var err error
		now := time.Now().UTC().Round(time.Microsecond)

		item = &entity.UserLog{
			Id:       uuid.New(),
			UserId:   ai.UserId,
			Ip:       ai.Ip,
			Username: ai.Username,
			Action:   action,
			ActionTs: now,
			SupInfo:  ai.GetSupInfo(),
		}

		defer func() {
			if err != nil {
				item = nil
			}
		}()

		//sqlUserLogAdd = `INSERT INTO tbl_user_log(id, user_id, username, ip, action, action_ts, sup_info)
		//			VALUES($1, $2, $3, $4, $5, $6, $7)`
		var cmdTag pgconn.CommandTag
		cmdTag, err = tx.Exec(ctx, sqlUserLogAdd, item.Id, item.UserId, item.Username, item.Ip, item.Action,
			item.ActionTs, item.SupInfo)
		if err != nil {
			log.Error("error at sqlUserLogAdd: %s", err)
			return err
		}
		if cmdTag.RowsAffected() == 0 {
			err = ErrNoRowsAffected
			log.Error("no rows affected during insert: %s", err)
			return err
		}

		return nil
	}
	return opt
}
