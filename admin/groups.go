package admin

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

// GroupData ...
type GroupData struct {
	Name string
}

// GroupsData ...
type GroupsData struct {
	Title         string
	Paths         []helpers.Path
	Groups        []GroupData
	NewItemButton bool
	NewItemText   string
	NewItemLink   string
}

/*
GetGroups ...
*/
func GetGroups(db *sql.DB) ([]GroupData, error) {
	var groups []GroupData

	rows, err := db.Query(`SELECT name FROM user_groups`)
	if err != nil {
		return groups, err
	}
	defer rows.Close()

	for rows.Next() {
		var group GroupData
		err = rows.Scan(&group.Name)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	err = rows.Err()
	if err != nil {
		return groups, err
	}
	return groups, nil
}

// Groups ...
func Groups(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	groups, err := GetGroups(db)
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	p := &GroupsData{"Groups", []helpers.Path{{Name: "Admin", Link: "/admin"}}, groups, true, "New group", "/groups/new"}
	helpers.RenderTemplate(r, w, "admin_groups", p)
}

// GroupPageData ...
type GroupPageData struct {
	Title string
	Paths []helpers.Path
	Group GroupData
}

/*
GetGroupByName ...
*/
func GetGroupByName(db *sql.DB, name string) (GroupData, error) {
	var group GroupData

	query := `SELECT name FROM user_groups WHERE name = ?`
	err := db.QueryRow(query, name).Scan(&group.Name)
	if err != nil {
		return group, err
	}
	return group, nil
}

// Group ...
func Group(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["group"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	group, err := GetGroupByName(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &GroupPageData{group.Name, []helpers.Path{{Name: "Admin", Link: "/admin"}}, group}
	helpers.RenderTemplate(r, w, "admin_group", p)
}

/*
UpdateGroupInDB ...
*/
func UpdateGroupInDB(db *sql.DB, newName, oldName string) error {
	insForm, err := db.Prepare("UPDATE user_groups SET name=? WHERE name=?")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(newName, oldName)
	if err != nil {
		return err
	}
	return nil
}

// UpdateGroup ...
func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["group"]
	value := r.FormValue("name")

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	err = UpdateGroupInDB(db, value, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	http.Redirect(w, r, "/admin/groups/"+value, http.StatusSeeOther)
}

// NewGroup ...
func NewGroup(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "New group", Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}}}
	helpers.RenderTemplate(r, w, "admin_new_group", p)
}

/*
InsertNewGroupInDB ...
*/
func InsertNewGroupInDB(db *sql.DB, name string) error {
	insForm, err := db.Prepare("INSERT INTO user_groups (name) VALUES (?)")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(name)
	if err != nil {
		return err
	}

	return nil
}

// InsertNewGroup ...
func InsertNewGroup(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	name := r.FormValue("name")
	err = InsertNewGroupInDB(db, name)
	if err != nil {
		helpers.HandleError(err)
	}
	http.Redirect(w, r, "/admin/groups/"+name, http.StatusSeeOther)
}
