package use_case

import (
	"context"
	"ykjam/new_tx/internal/entity"
	"ykjam/new_tx/internal/repo"
	"ykjam/new_tx/pkg/logger"
)

type departmentRepo interface {
	DepartmentAdd(name string, priority int) (*entity.Department, repo.Operation)
	DepartmentUpdate(item *entity.Department, name string, priority int) repo.Operation
	DepartmentChangeState(item *entity.Department, state entity.State) repo.Operation
	DepartmentById(ctx context.Context, id int) (*entity.Department, error)
	DepartmentList(ctx context.Context) ([]*entity.Department, error)

	RunInTx(ctx context.Context, ops ...repo.Operation) error
	UserLogAdd(ai *entity.ActionInfo, action entity.UserAction) repo.Operation
}

type Department struct {
	r departmentRepo
}

func NewDepartment(r departmentRepo) *Department {
	return &Department{
		r: r,
	}
}

func (u *Department) Add(ctx context.Context, ai *entity.ActionInfo, name string, priority int) (*entity.Department, error) {

	log := logger.Get()
	var err error
	var item *entity.Department
	var depOp, userLogOp repo.Operation

	item, depOp = u.r.DepartmentAdd(name, priority)
	userLogOp = u.r.UserLogAdd(ai, entity.UserActionDepartmentAdd)

	err = u.r.RunInTx(ctx, depOp, userLogOp)
	if err != nil {
		log.Error("in departmentRepo.RunInTx: %s", err)
		return nil, ErrInternalServerError
	}
	return item, nil

}

func (u *Department) Update(ctx context.Context, ai *entity.ActionInfo, id int, name string, priority int) error {

	log := logger.Get()
	var err error
	var item *entity.Department
	var depOp, userLogOp repo.Operation

	item, err = u.r.DepartmentById(ctx, id)
	if err != nil {
		log.Error("in departmentRepo.DepartmentById: %s", err)
		return ErrInternalServerError
	}
	if item == nil {
		log.Error("department not found")
		return ErrNotFound
	}

	depOp = u.r.DepartmentUpdate(item, name, priority)
	userLogOp = u.r.UserLogAdd(ai, entity.UserActionDepartmentUpdate)

	err = u.r.RunInTx(ctx, depOp, userLogOp)
	if err != nil {
		log.Error("error in departmentRepo.RunInTx: %s", err)
		return ErrInternalServerError
	}
	return nil

}

func (u *Department) ChangeState(ctx context.Context, ai *entity.ActionInfo, id int, state entity.State) error {

	log := logger.Get()
	var err error
	var item *entity.Department
	var depOp, userLogOp repo.Operation

	item, err = u.r.DepartmentById(ctx, id)
	if err != nil {
		log.Error("in departmentRepo.DepartmentById: %s", err)
		return ErrInternalServerError
	}
	if item == nil {
		log.Error("department not found")
		return ErrNotFound
	}

	depOp = u.r.DepartmentChangeState(item, state)
	userLogOp = u.r.UserLogAdd(ai, entity.UserActionDepartmentChangeState)

	err = u.r.RunInTx(ctx, depOp, userLogOp)
	if err != nil {
		log.Error("error in repo.RunInTx: %s", err)
		return ErrInternalServerError
	}
	return nil

}

func (u *Department) GetId(ctx context.Context, id int) (*entity.Department, error) {

	log := logger.Get()
	var err error
	var item *entity.Department

	item, err = u.r.DepartmentById(ctx, id)
	if err != nil {
		log.Error("in departmentRepo.DepartmentById: %s", err)
		return nil, ErrInternalServerError
	}
	if item == nil {
		log.Error("department not found")
		return nil, ErrNotFound
	}

	return item, nil

}

func (u *Department) List(ctx context.Context) ([]*entity.Department, error) {

	log := logger.Get()
	var err error
	var items []*entity.Department

	items, err = u.r.DepartmentList(ctx)
	if err != nil {
		log.Error("in departmentRepo.UserList: %s", err)
		return nil, ErrInternalServerError
	}

	return items, nil

}
