package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/impact-eintr/WebKits/erbac"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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

func main() {
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
	rbac := erbac.New()
	permissions := make(erbac.Permissions)

	// Build roles and add them to goRBAC instance
	for rid, pids := range jsonRoles {
		role := erbac.NewStdRole(rid)
		for _, pid := range pids {
			_, ok := permissions[pid]
			if !ok {
				permissions[pid] = erbac.NewStdPermission(pid)
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

	if rbac.IsGranted("root", permissions["add-table"], nil) {
		log.Println("root can add table")
	}

	if !rbac.IsGranted("manager", permissions["add-record"], nil) {
		log.Println("manager can not add all record")
	}

	if !rbac.IsGranted("user", permissions["add-all-record"], nil) {
		log.Println("user can not add all record")
	}

	// Check if `nobody` can add text
	// `nobody` is not exist in goRBAC at the moment
	//if !rbac.IsGranted("nobody", permissions["read-text"], nil) {
	//	log.Println("Nobody can't read text")
	//}
	//// Add `nobody` and assign `read-text` permission
	//nobody := erbac.NewStdRole("nobody")
	//permissions["read-text"] = erbac.NewStdPermission("read-text")
	//nobody.Assign(permissions["read-text"])
	//rbac.Add(nobody)
	//// Check if `nobody` can read text again
	//if rbac.IsGranted("nobody", permissions["read-text"], nil) {
	//	log.Println("Nobody can read text")
	//}

	// Persist the change
	// map[RoleId]PermissionIds
	//jsonOutputRoles := make(map[string][]string)
	//// map[RoleId]ParentIds
	//jsonOutputInher := make(map[string][]string)
	//SaveJsonHandler := func(r erbac.Role, parents []string) error {
	//	// WARNING: Don't use erbac.RBAC instance in the handler,
	//	// otherwise it causes deadlock.
	//	permissions := make([]string, 0)
	//	for _, p := range r.(*erbac.StdRole).Permissions() {
	//		permissions = append(permissions, p.ID())
	//	}
	//	jsonOutputRoles[r.ID()] = permissions
	//	jsonOutputInher[r.ID()] = parents
	//	return nil
	//}
	//if err := erbac.Walk(rbac, SaveJsonHandler); err != nil {
	//	log.Fatalln(err)
	//}

	//// Save roles information
	//if err := SaveJson("new-roles.json", &jsonOutputRoles); err != nil {
	//	log.Fatal(err)
	//}
	//// Save inheritance information
	//if err := SaveJson("new-inher.json", &jsonOutputInher); err != nil {
	//	log.Fatal(err)
	//}
}
