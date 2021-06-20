package erbac

type Permission interface {
	ID() string
	Match(Permission) bool
}

type Permissions map[string]Permission

type StdPermission struct {
	IDStr string
}

func NewStdPermission(id string) Permission {
	return &StdPermission{id}
}

func (sp *StdPermission) ID() string {
	return sp.IDStr
}

func (sp *StdPermission) Match(p Permission) bool {
	return sp.IDStr == p.ID()
}
