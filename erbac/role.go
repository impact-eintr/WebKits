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

// Permit returns true if the role has specific permission.
func (role *StdRole) Permit(p Permission) (rslt bool) {
}
