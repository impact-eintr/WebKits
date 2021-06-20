package erbac

import (
	"errors"
	"log"
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

func NewRBAC() *RBAC {
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

// 从文件中构建erbac
func BuildRBAC(roleFile, inherFile string) (*RBAC, Permissions) {
	// map[RoleId]PermissionIds
	var jsonRoles map[string][]string
	// map[RoleId]ParentIds
	var jsonInher map[string][]string
	// Load roles information
	if err := LoadJson("roles.json", &jsonRoles); err != nil {
		log.Fatal(err)
	}
	// Load inheritance information
	if err := LoadJson("inher.json", &jsonInher); err != nil {
		log.Fatal(err)
	}

	rbac := NewRBAC()
	permissions := make(Permissions)

	// Build roles and add them to goRBAC instance
	for rid, pids := range jsonRoles {
		role := NewStdRole(rid)
		for _, pid := range pids {
			_, ok := permissions[pid]
			if !ok {
				permissions[pid] = NewStdPermission(pid)
			}
			role.Assign(permissions[pid])
		}
		rbac.Add(role)
	}
	// Assign the inheritance relationship
	for rid, parents := range jsonInher {
		if err := rbac.SetParents(rid, parents); err != nil {
			log.Fatal(err)
		}
	}

	return rbac, permissions
}
