package erbac

import (
	"errors"
	"sync"
)

type RBAC struct {
	mutex   sync.RWMutex
	roles   Roles
	parents map[string]map[string]struct{}
}

var (
	ErrRoleExist    = errors.New("角色已经存在")
	ErrRoleNotExist = errors.New("角色不存在")
	ErrFounfdCircle = errors.New("发现环继承")
	Empty           = struct{}{}
)

type AssertionFunc func(*RBAC, string, Permission) bool

func NewRBAC() *RBAC {
	return &RBAC{
		roles:   make(Roles),
		parents: make(map[string]map[string]struct{}),
	}
}

// 给当前的rbac添加一个拥有权限的角色
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

// 移除一个角色
func (rbac *RBAC) Remove(id string) (err error) {
	rbac.mutex.Lock()

	if _, ok := rbac.roles[id]; ok {
		delete(rbac.roles, id)
		for rid, parents := range rbac.parents {
			if rid == id {
				delete(rbac.parents, rid)
				continue
			}
			for parent := range parents {
				if parent == id {
					delete(rbac.parents[rid], id)
					break
				}
			}
		}
	} else {
		err = ErrRoleExist
	}

	rbac.mutex.Unlock()
	return

}

func (rbac *RBAC) GetParents(id string) (parents []string, err error) {
	rbac.mutex.RLock()
	defer rbac.mutex.RUnlock()

	if _, ok := rbac.roles[id]; !ok {
		return nil, ErrRoleNotExist
	}

	if ids, ok := rbac.parents[id]; ok {
		for parent := range ids {
			parents = append(parents, parent)
		}
	}
	return

}

// 设置单个parent
func (rbac *RBAC) SetParent(id string, parent string) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()

	if _, ok := rbac.roles[id]; !ok {
		return ErrRoleNotExist
	}
	if _, ok := rbac.roles[parent]; !ok {
		return ErrRoleNotExist
	}
	if _, ok := rbac.parents[id]; !ok {
		rbac.parents[id] = make(map[string]struct{})
	}

	rbac.parents[id][parent] = Empty

	return InherCircle(rbac)

}

// 设置多个parent
func (rbac *RBAC) SetParents(id string, parents []string) error {
	rbac.mutex.Lock()
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
	rbac.mutex.Unlock()

	return InherCircle(rbac)

}

// 移除单个parent
func (rbac *RBAC) RemoveParent(id string, parent string) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()

	if _, ok := rbac.roles[id]; !ok {
		return ErrRoleNotExist
	}
	if _, ok := rbac.roles[parent]; !ok {
		return ErrRoleNotExist
	}
	delete(rbac.parents[id], parent)

	return nil

}

func (rbac *RBAC) Get(id string) (r Role, parents []string, err error) {
	rbac.mutex.RLock()
	rbac.mutex.RUnlock()
	return
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

func (rbac *RBAC) SaveUserRBAC(newRoleFile, newInherFile string) error {
	// Persist the change
	// map[RoleId]PermissionIds
	jsonOutputRoles := make(map[string][]string)
	// map[RoleId]ParentIds
	jsonOutputInher := make(map[string][]string)
	SaveJsonHandler := func(r Role, parents []string) error {
		// WARNING: Don't use erbac.RBAC instance in the handler,
		// otherwise it causes deadlock.
		permissions := make([]string, 0)
		for _, p := range r.(*StdRole).Permissions() {
			permissions = append(permissions, p.ID())
		}
		jsonOutputRoles[r.ID()] = permissions
		jsonOutputInher[r.ID()] = parents
		return nil
	}
	if err := Walk(rbac, SaveJsonHandler); err != nil {
		return err
	}

	// Save roles information
	if err := SaveJson(newRoleFile, &jsonOutputRoles); err != nil {
		return err
	}
	// Save inheritance information
	if err := SaveJson(newInherFile, &jsonOutputInher); err != nil {
		return err
	}
	return nil

}

// 从文件中构建erbac
func BuildRBAC(roleFile, inherFile string) (*RBAC, Permissions, error) {
	// map[RoleId]PermissionIds
	var jsonRoles map[string][]string

	// map[RoleId]ParentIds
	var jsonInher map[string][]string

	// Load roles information
	if err := LoadJson(roleFile, &jsonRoles); err != nil {
		return nil, nil, err
	}

	// Load inheritance information
	if err := LoadJson(inherFile, &jsonInher); err != nil {
		return nil, nil, err
	}

	rbac := NewRBAC()
	permissions := make(Permissions)

	// Build roles and add them to eRBAC instance
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
			return nil, nil, err
		}
	}
	return rbac, permissions, nil

}
