package admin

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

// UserData ...
type UserData struct {
	Username string
	Email    string
}

// UsersData ...
type UsersData struct {
	Title string
	Paths []helpers.Path
	Users []UserData
}

/*
GetUsers ...
*/
func GetUsers(db *sql.DB) ([]UserData, error) {
	rows, err := db.Query(`SELECT username, email FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserData
	for rows.Next() {
		var user UserData
		err = rows.Scan(&user.Username, &user.Email)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return users, err
	}
	return users, nil
}

// Users ...
func Users(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	users, err := GetUsers(db)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &UsersData{"Users", []helpers.Path{{Name: "Admin", Link: "/admin"}}, users}
	helpers.RenderTemplate(r, w, "admin_users", p)
}

// GroupUserData ...
type GroupUserData struct {
	Name         string
	UserIsMember bool
}

// UserPageData ...
type UserPageData struct {
	Title  string
	Paths  []helpers.Path
	User   UserData
	Groups []GroupUserData
}

/*
GetGroupsUserIsPartOf ...
*/
func GetGroupsUserIsPartOf(db *sql.DB, username string) ([]GroupUserData, error) {
	rows, err := db.Query(`
		SELECT name,
		(
			SELECT name
			FROM user_groups
			LEFT OUTER JOIN group_members
			ON user_groups.id = group_members.group_id
			LEFT OUTER JOIN users
			ON group_members.user_id = users.id
			WHERE users.username = ?
		) IS NOT NULL as userIsMember
		FROM user_groups;
	`, username)
	if err != nil {
		return nil, err
	}
	var groups []GroupUserData
	for rows.Next() {
		var group GroupUserData
		err = rows.Scan(&group.Name, &group.UserIsMember)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	err = rows.Err()
	if err != nil {
		return groups, err
	}
	defer rows.Close()
	return groups, nil
}

/*
GetUserByUsername ...
*/
func GetUserByUsername(db *sql.DB, username string) (UserData, error) {
	var user UserData

	query := `SELECT username, email FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&user.Username, &user.Email)
	if err != nil {
		// No user
		return user, err
	}
	return user, nil
}

// User ...
func User(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["user"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	groups, err := GetGroupsUserIsPartOf(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	user, err := GetUserByUsername(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	data := &UserPageData{user.Username, []helpers.Path{{Name: "Admin", Link: "/admin"}}, user, groups}
	helpers.RenderTemplate(r, w, "admin_user", data)
}
