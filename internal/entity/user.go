package entity

import (
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           int
	Username     string
	Password     string
	Firstname    string
	Lastname     string
	DepartmentId int
	Department   string
	PositionId   int
	Position     string
	State        UserState
	CreateTs     time.Time
	UpdateTs     time.Time
	Version      int
}

type ActionInfo struct {
	// Id of user performing action
	UserId int
	// Username of user performing action
	Username string
	// RemoteAddr of user performing action
	Ip net.IP
	// Role of user performing action
	UserRole UserRole
	// Additional info parameters of types int, string, float32, float64, int64, bool
	params map[string]interface{}
}

func (a *ActionInfo) ParamAdd(k string, v interface{}) (err error) {
	if a.params == nil {
		a.params = make(map[string]interface{})
	}
	switch v.(type) {
	case int, string, float32, float64, int64, bool, uuid.UUID:
		a.params[k] = v
		return
	default:
		err = errors.New("unsupported value type")
		return
	}
}

func (a *ActionInfo) GetSupInfo() (message json.RawMessage) {
	message, _ = json.Marshal(a.params)
	return
}

type UserRole string

const (
	UserRoleSuperadmin UserRole = "SUPERADMIN"
	UserRoleAdmin      UserRole = "ADMIN"
	UserRoleRegistry   UserRole = "REGISTRY"
	UserRoleHR         UserRole = "HR"
	UserRoleUser       UserRole = "USER"
	UserRoleSystem     UserRole = "SYSTEM"
)

type UserState string

const (
	UserStateActive  UserState = "ACTIVE"
	UserStateBlocked UserState = "BLOCKED"
	UserStateDeleted UserState = "DELETED"
)

type TfaType string

const (
	TfaTypeInactive TfaType = "INACTIVE"
	TfaTypeSms      TfaType = "SMS"
	TfaTypeFido     TfaType = "FIDO"
	TfaTypeOtp      TfaType = "OTP"
	TfaTypeEmail    TfaType = "EMAIL"
)
