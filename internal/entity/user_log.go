package entity

import (
	"encoding/json"
	"github.com/google/uuid"
	"net"
	"time"
)

type UserLog struct {
	Id       uuid.UUID
	UserId   int
	Username string
	Ip       net.IP
	Action   UserAction
	ActionTs time.Time
	SupInfo  json.RawMessage
}

type UserAction string

const (
	UserActionLogIn  UserAction = "user-login"
	UserActionLogOut UserAction = "user-logout"

	UserActionUserAdd         UserAction = "user-add"
	UserActionUserUpdate      UserAction = "user-update"
	UserActionUserChangeState UserAction = "user-change-state"

	UserActionDepartmentAdd         UserAction = "dep-add"
	UserActionDepartmentUpdate      UserAction = "dep-update"
	UserActionDepartmentChangeState UserAction = "dep-change-state"

	UserActionPositionAdd         UserAction = "pos-add"
	UserActionPositionUpdate      UserAction = "pos-update"
	UserActionPositionChangeState UserAction = "pos-change-state"
)
