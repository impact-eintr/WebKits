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
	return json.NewDecoder(f).Decode(v)
}

func SaveJson(filename string, v interface{}) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(v)
}
