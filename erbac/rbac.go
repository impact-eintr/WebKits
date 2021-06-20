package erbac

import (
	"errors"
	"sync"
)

type RBAC struct {
	mutex   sync.Mutex
	roles   Roles
	parents map[string]map[string]struct{}
}

var (
	ErrRoleExist    = errors.New("角色已经存在")
	ErrRoleNotExist = errors.New("角色不存在")
	Empty           = struct{}{}
)

type AssertionFunc func(*RBAC, string, Permission) bool

func New() *RBAC {
	return &RBAC{
		roles:   make(Roles),
		parents: make(map[string]map[string]struct{}),
	}
}

func (rbac *RBAC) Add(r Role) (err error) {
	rbac.mutex.Lock()
	if _, ok := rbac.roles[r.ID()]; !ok {
		rbac.roles[r.ID()] = r
	} else {
		err = ErrRoleExist
	}
	rbac.mutex.Unlock()
	return

}

func (rbac *RBAC) SetParents(id string, parents []string) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()

	if _, ok := rbac.roles[id]; !ok {
		return ErrRoleNotExist
	}

	for _, parent := range parents {
		if _, ok := rbac.roles[parent]; !ok {
			return ErrRoleNotExist
		}
	}
	if _, ok := rbac.parents[id]; !ok {
		rbac.parents[id] = make(map[string]struct{})
	}

	for _, parent := range parents {
		rbac.parents[id][parent] = Empty
	}
	return nil

}

func (rbac *RBAC) IsGranted(id string, p Permission, assert AssertionFunc) (res bool) {
	rbac.mutex.Lock()
	res = rbac.isGranted(id, p, assert)
	rbac.mutex.Unlock()
	return
}

func (rbac *RBAC) isGranted(id string, p Permission, assert AssertionFunc) bool {
	if assert != nil && !assert(rbac, id, p) {
		return false
	}
	return rbac.recursionCheck(id, p)
}

// 循环检测
func (rbac *RBAC) recursionCheck(id string, p Permission) bool {
	if role, ok := rbac.roles[id]; ok {
		if role.Permit(p) {
			return true
		}
		if parents, ok := rbac.parents[id]; ok {
			for pID := range parents {
				if _, ok := rbac.roles[id]; ok {
					if rbac.recursionCheck(pID, p) {
						return true
					}
				}
			}
		}
	}
	return false
}
