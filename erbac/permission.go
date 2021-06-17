package erbac

type Permission interface {
	ID() string
	Match(Permission) bool
}

type Permissions map[string]Permission
