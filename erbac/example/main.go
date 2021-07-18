package main

import (
	"log"

	"github.com/impact-eintr/WebKits/erbac"
)

func main() {
	rbac, permissions, err := erbac.BuildRBAC("./roles.json", "./inher.json")
	if err != nil {
		log.Fatalln(err)
	}

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

	nobody.Assign(permissions["read-record"])
	rbac.Add(nobody)
	// Check if `nobody` can read text again
	if rbac.IsGranted("nobody", permissions["read-record"], nil) {
		log.Println("Nobody can read record")
	}

	err = rbac.SaveUserRBAC("newRoles.json", "newInher.json")
	if err != nil {
		log.Fatalln(err)
	}

}
