package use_case

import (
	"context"
	"ykjam/new_tx/internal/entity"
	"ykjam/new_tx/internal/repo"
	"ykjam/new_tx/pkg/logger"
)

type userRepo interface {
	UserAdd(username, password, firstname, lastname string, departmentId, positionI int,
		state entity.UserState) (*entity.User, repo.Operation)
	UserUpdate(item *entity.User, firstname, lastname string, depId, posId int) repo.Operation
	UserChangeState(item *entity.User, state entity.UserState) repo.Operation
	UserById(ctx context.Context, id int) (*entity.User, error)
	UserList(ctx context.Context) ([]*entity.User, error)

	RunInTx(ctx context.Context, ops ...repo.Operation) error
	UserLogAdd(ai *entity.ActionInfo, action entity.UserAction) repo.Operation
}

type User struct {
	r userRepo
}

func NewUser(r userRepo) *User {
	return &User{
		r: r,
	}
}

func (u *User) Add(ctx context.Context, ai *entity.ActionInfo, username, password, firstname, lastname string,
	departmentId, positionId int, state entity.UserState) (*entity.User, error) {

	log := logger.Get()
	var err error
	var item *entity.User
	var userOp, userLogOp repo.Operation

	item, userOp = u.r.UserAdd(username, password, firstname, lastname, departmentId, positionId, state)
	userLogOp = u.r.UserLogAdd(ai, entity.UserActionUserAdd)

	err = u.r.RunInTx(ctx, userOp, userLogOp)
	if err != nil {
		log.Error("in userRepo.RunInTx: %s", err)
		return nil, ErrInternalServerError
	}
	return item, nil

}

func (u *User) Update(ctx context.Context, ai *entity.ActionInfo, id int, firstname, lastname string,
	departmentId, positionId int) error {

	log := logger.Get()
	var err error
	var item *entity.User
	var userOp, userLogOp repo.Operation

	item, err = u.r.UserById(ctx, id)
	if err != nil {
		log.Error("in userRepo.UserById: %s", err)
		return ErrInternalServerError
	}
	if item == nil {
		log.Error("user not found")
		return ErrNotFound
	}

	userOp = u.r.UserUpdate(item, firstname, lastname, departmentId, positionId)
	userLogOp = u.r.UserLogAdd(ai, entity.UserActionUserUpdate)

	err = u.r.RunInTx(ctx, userOp, userLogOp)
	if err != nil {
		log.Error("error in userRepo.RunInTx: %s", err)
		return ErrInternalServerError
	}
	return nil

}

func (u *User) ChangeState(ctx context.Context, ai *entity.ActionInfo, id int, state entity.UserState) error {

	log := logger.Get()
	var err error
	var item *entity.User
	var userOp, userLogOp repo.Operation

	item, err = u.r.UserById(ctx, id)
	if err != nil {
		log.Error("in userRepo.UserById: %s", err)
		return ErrInternalServerError
	}
	if item == nil {
		log.Error("user not found")
		return ErrNotFound
	}

	userOp = u.r.UserChangeState(item, state)
	userLogOp = u.r.UserLogAdd(ai, entity.UserActionUserChangeState)

	err = u.r.RunInTx(ctx, userOp, userLogOp)
	if err != nil {
		log.Error("error in userRepo.RunInTx: %s", err)
		return ErrInternalServerError
	}
	return nil

}

func (u *User) GetId(ctx context.Context, id int) (*entity.User, error) {

	log := logger.Get()
	var err error
	var item *entity.User

	item, err = u.r.UserById(ctx, id)
	if err != nil {
		log.Error("in userRepo.UserById: %s", err)
		return nil, ErrInternalServerError
	}
	if item == nil {
		log.Error("user not found")
		return nil, ErrNotFound
	}

	return item, nil

}

func (u *User) List(ctx context.Context) ([]*entity.User, error) {

	log := logger.Get()
	var err error
	var items []*entity.User

	items, err = u.r.UserList(ctx)
	if err != nil {
		log.Error("in userRepo.UserList: %s", err)
		return nil, ErrInternalServerError
	}

	return items, nil

}
