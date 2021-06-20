package erbac

import (
	"sync"
)

// Role is an interface.
type Role interface {
	ID() string
	Permit(Permission) bool
}

// Roles is a map
type Roles map[string]Role

type StdRole struct {
	sync.RWMutex
	IDStr       string `json:"id"`
	permissions Permissions
}

func NewStdRole(id string) *StdRole {
	role := &StdRole{
		IDStr:       id,
		permissions: make(Permissions),
	}
	return role
}

// ID returns the role's identity name.
func (role *StdRole) ID() string {
	return role.IDStr
}

// Permit 进行权限判断
func (role *StdRole) Permit(p Permission) (rslt bool) {
	role.RLock()
	for _, rp := range role.permissions {
		if rp.Match(p) {
			rslt = true
			break
		}
	}
	role.RUnlock()
	return
}

// Assign 分配权限
func (role *StdRole) Assign(p Permission) error {
	role.Lock()
	role.permissions[p.ID()] = p
	role.Unlock()
	return nil

}

// 展开stdRole的Permissions
func (role *StdRole) Permissions() []Permission {
	role.RLock()
	res := make([]Permission, 0, len(role.permissions))
	for _, p := range role.permissions {
		res = append(res, p)
	}
	role.RUnlock()
	return res

}
