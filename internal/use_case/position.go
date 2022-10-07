package use_case

import (
	"context"
	"ykjam/new_tx/internal/entity"
	"ykjam/new_tx/internal/repo"
	"ykjam/new_tx/pkg/logger"
)

type positionRepo interface {
	PositionAdd(name string, priority int) (*entity.Position, repo.Operation)
	PositionUpdate(item *entity.Position, name string, priority int) repo.Operation
	PositionChangeState(item *entity.Position, state entity.State) repo.Operation
	PositionById(ctx context.Context, id int) (*entity.Position, error)
	PositionList(ctx context.Context) ([]*entity.Position, error)

	RunInTx(ctx context.Context, ops ...repo.Operation) error
	UserLogAdd(ai *entity.ActionInfo, action entity.UserAction) repo.Operation
}

type Position struct {
	r positionRepo
}

func NewPosition(r positionRepo) *Position {
	return &Position{
		r: r,
	}
}

func (u *Position) Add(ctx context.Context, ai *entity.ActionInfo, name string, priority int) (*entity.Position, error) {

	log := logger.Get()
	var err error
	var item *entity.Position
	var depOp, userLogOp repo.Operation

	item, depOp = u.r.PositionAdd(name, priority)
	userLogOp = u.r.UserLogAdd(ai, entity.UserActionPositionAdd)

	err = u.r.RunInTx(ctx, depOp, userLogOp)
	if err != nil {
		log.Error("in positionRepo.RunInTx: %s", err)
		return nil, ErrInternalServerError
	}
	return item, nil

}

func (u *Position) Update(ctx context.Context, ai *entity.ActionInfo, id int, name string, priority int) error {

	log := logger.Get()
	var err error
	var item *entity.Position
	var depOp, userLogOp repo.Operation

	item, err = u.r.PositionById(ctx, id)
	if err != nil {
		log.Error("in positionRepo.PositionById: %s", err)
		return ErrInternalServerError
	}
	if item == nil {
		log.Error("position not found")
		return ErrNotFound
	}

	depOp = u.r.PositionUpdate(item, name, priority)
	userLogOp = u.r.UserLogAdd(ai, entity.UserActionPositionUpdate)

	err = u.r.RunInTx(ctx, depOp, userLogOp)
	if err != nil {
		log.Error("error in positionRepo.RunInTx: %s", err)
		return ErrInternalServerError
	}
	return nil

}

func (u *Position) ChangeState(ctx context.Context, ai *entity.ActionInfo, id int, state entity.State) error {

	log := logger.Get()
	var err error
	var item *entity.Position
	var depOp, userLogOp repo.Operation

	item, err = u.r.PositionById(ctx, id)
	if err != nil {
		log.Error("in positionRepo.PositionById: %s", err)
		return ErrInternalServerError
	}
	if item == nil {
		log.Error("position not found")
		return ErrNotFound
	}

	depOp = u.r.PositionChangeState(item, state)
	userLogOp = u.r.UserLogAdd(ai, entity.UserActionPositionChangeState)

	err = u.r.RunInTx(ctx, depOp, userLogOp)
	if err != nil {
		log.Error("error in repo.RunInTx: %s", err)
		return ErrInternalServerError
	}
	return nil

}

func (u *Position) GetId(ctx context.Context, id int) (*entity.Position, error) {

	log := logger.Get()
	var err error
	var item *entity.Position

	item, err = u.r.PositionById(ctx, id)
	if err != nil {
		log.Error("in positionRepo.PositionById: %s", err)
		return nil, ErrInternalServerError
	}
	if item == nil {
		log.Error("position not found")
		return nil, ErrNotFound
	}

	return item, nil

}

func (u *Position) List(ctx context.Context) ([]*entity.Position, error) {

	log := logger.Get()
	var err error
	var items []*entity.Position

	items, err = u.r.PositionList(ctx)
	if err != nil {
		log.Error("in positionRepo.UserList: %s", err)
		return nil, ErrInternalServerError
	}

	return items, nil

}
