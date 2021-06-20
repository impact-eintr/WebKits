package main

import (
	"fmt"
	"log"

	"github.com/impact-eintr/WebKits/erbac"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	rbac, permissions := erbac.BuildRBAC("./role.json", "./inher.json")

	if rbac.IsGranted("root", permissions["add-table"], nil) {
		log.Println("root can add table")
	}

	if rbac.IsGranted("root", permissions["add-all-record"], nil) {
		log.Println("root can add all record")
	}

	if rbac.IsGranted("root", permissions["add-record"], nil) {
		log.Println("root can add record")
	}

	if rbac.IsGranted("manager", permissions["add-record"], nil) {
		log.Println("manager can add record")
	}

	if !rbac.IsGranted("user", permissions["add-all-record"], nil) {
		log.Println("user can not add all record")
	}

	// Check if `nobody` can add text
	// `nobody` is not exist in goRBAC at the moment
	//if !rbac.IsGranted("nobody", permissions["read-text"], nil) {
	//	log.Println("Nobody can't read text")
	//}
	// Add `nobody` and assign `read-text` permission
	nobody := erbac.NewStdRole("nobody")
	permissions["read-record"] = erbac.NewStdPermission("read-record")

	fmt.Println(permissions)

	nobody.Assign(permissions["read-record"])
	rbac.Add(nobody)
	// Check if `nobody` can read text again
	if rbac.IsGranted("nobody", permissions["read-record"], nil) {
		log.Println("Nobody can read record")
	}

	// Persist the change
	// map[RoleId]PermissionIds
	jsonOutputRoles := make(map[string][]string)
	// map[RoleId]ParentIds
	jsonOutputInher := make(map[string][]string)
	SaveJsonHandler := func(r erbac.Role, parents []string) error {
		// WARNING: Don't use erbac.RBAC instance in the handler,
		// otherwise it causes deadlock.
		permissions := make([]string, 0)
		for _, p := range r.(*erbac.StdRole).Permissions() {
			permissions = append(permissions, p.ID())
		}
		jsonOutputRoles[r.ID()] = permissions
		jsonOutputInher[r.ID()] = parents
		return nil
	}
	if err := erbac.Walk(rbac, SaveJsonHandler); err != nil {
		log.Fatalln(err)
	}

	// Save roles information
	if err := erbac.SaveJson("new-roles.json", &jsonOutputRoles); err != nil {
		log.Fatal(err)
	}
	// Save inheritance information
	if err := erbac.SaveJson("new-inher.json", &jsonOutputInher); err != nil {
		log.Fatal(err)
	}
}
