package erbac

import (
	"encoding/json"
	"errors"
	"os"
)

type WalkHandler func(Role, []string) error

func Walk(rbac *RBAC, h WalkHandler) (err error) {
	if h == nil {
		return errors.New("WalkHandler is nil")
	}

	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()

	for id := range rbac.roles {
		var parents []string
		r := rbac.roles[id]

		for parent := range rbac.parents[id] {
			// fmt.Println("id: ", id, "parent: ", parent)
			parents = append(parents, parent)
		}
		if err := h(r, parents); err != nil {
			return err
		}
	}

	return
}

func LoadJson(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return j son.NewDecoder(f).Decode(v)
}

func SaveJson(filename string, v interface{}) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(v)
}

func InherCircle(rbac *RBAC) (err error) {
	rbac.mutex.Lock()

	skipped := make(map[string]struct{}, len(rbac.roles))
	var stack []string

	for id := range rbac.roles {
		if err = dfs(rbac, id, skipped, stack); err != nil {
			break
		}
	}

	rbac.mutex.Unlock()
	return err

}

func dfs(rbac *RBAC, id string, skipped map[string]struct{}, stack []string) error {
	if _, ok := skipped[id]; ok {
		return nil
	}

	for _, item := range stack {
		if item == id {
			return ErrFounfdCircle
		}
	}

	parents := rbac.parents[id]
	if len(parents) == 0 {
		stack = nil
		skipped[id] = Empty
		return nil
	}

	stack = append(stack, id)
	for pid := range parents {
		if err := dfs(rbac, pid, skipped, stack); err != nil {
			return err
		}
	}
	return nil

}

func AnyGranted(rbac *RBAC, roles []string, permission Permission,
	assert AssertionFunc) (r bool) {
	rbac.mutex.Lock()
	for _, role := range roles {
		if rbac.isGranted(role, permission, assert) {
			r = true
			break
		}
	}
	rbac.mutex.Unlock()
	return r

}

func AllGranted(rbac *RBAC, roles []string, permission Permission,
	assert AssertionFunc) (r bool) {
	rbac.mutex.Lock()

	for _, role := range roles {
		if !rbac.isGranted(role, permission, assert) {
			r = true
			break
		}
	}
	rbac.mutex.Unlock()
	return !r

}
