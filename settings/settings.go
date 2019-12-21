package settings

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/authentication"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"golang.org/x/crypto/bcrypt"
)

/*
UserPageData stored all of the
data connected to the user for the
profile and setting pages.
*/
type UserPageData struct {
	Title      string
	Username   string
	UserEmail  string
	Registered string
}

/*
PageData struct is used for the
data being rendered and should
not be used in other structs.
*/
type PageData struct {
	User  UserPageData
	Paths []helpers.Path
}

/*
GetUserSettings gets all of the user settings from
the database based on the username.
*/
func GetUserSettings(db *sql.DB, username string, title string) (UserPageData, error) {
	user := UserPageData{Title: title}
	var createdAt time.Time

	query := `SELECT username, email, created_at FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&user.Username, &user.UserEmail, &createdAt)
	if err != nil {
		// No user
		return user, err
	}
	user.Registered = createdAt.Format("January 02, 2006")

	return user, nil
}

/*
UserPage endpoint returns the private profile
page or general settings page for the user.
*/
func UserPage(w http.ResponseWriter, r *http.Request) {
	username, err := helpers.GetUsernameFromRequest(r)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	user, err := GetUserSettings(db, username, "Settings")
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}

	helpers.RenderTemplate(r, w, "settings", user)
}

/*
EmailPage endpoint displays the users email
and adds a possibility to update the email.
*/
func EmailPage(w http.ResponseWriter, r *http.Request) {
	username, err := helpers.GetUsernameFromRequest(r)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	user, err := GetUserSettings(db, username, "Email")
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}

	helpers.RenderTemplate(r, w, "settings_email", user)
}

/*
UpdateUserEmail updates the users email
in the database based on the username.
*/
func UpdateUserEmail(db *sql.DB, newEmail string, username string) error {
	insForm, err := db.Prepare("UPDATE users SET email=? WHERE username=?")
	if err != nil {
		return err
	}
	res, err := insForm.Exec(newEmail, username)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	// Will execute successfully even if you enter a username that doesn't exist.
	// So we need to check the affected rows and makes sure it's not 0
	if rowsAffected == 0 {
		return errors.New("No rows affected. One row should be affected")
	}
	return nil
}

/*
UpdateEmail endpoint parses the data from the form page
and updates the email in the database.
*/
func UpdateEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if !authentication.EmailIsValid(email) {
		helpers.HandleError(errors.New("Email is not valid"))
		// Should redirect back to the page with a flash message
		helpers.InternalServerError(w)
		return
	}

	username, err := helpers.GetUsernameFromRequest(r)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	defer db.Close()

	err = UpdateUserEmail(db, email, username)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	http.Redirect(w, r, "/settings/email", http.StatusSeeOther)
}

/*
PasswordPage endpoint is a page with
an empty form for updating the username.
Note: Cannot be served as a static HTML page
since we still need some data like username.
*/
func PasswordPage(w http.ResponseWriter, r *http.Request) {
	helpers.RenderTemplate(r, w, "settings_password", nil)
}

/*
UpdateUserPassword takes the raw password string
and hashes it with BCrypt and inserts the new password
into the database.
*/
func UpdateUserPassword(db *sql.DB, newPassword string, username string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	insForm, err := db.Prepare("UPDATE users SET password=? WHERE username=?")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(passwordHash, username)
	if err != nil {
		return err
	}
	return nil
}

/*
UpdatePassword parses the password form and
insert the new password in the database, hashed.
*/
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	username, err := helpers.GetUsernameFromRequest(r)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	err = UpdateUserPassword(db, password, username)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	http.Redirect(w, r, "/settings/password", http.StatusSeeOther)
}
