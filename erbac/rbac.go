package erbac

type RBAC struct {
	mutex sync.Mutex
	roles Roles
	parents map 
}
