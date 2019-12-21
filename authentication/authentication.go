package authentication

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"regexp"
	"text/template"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/mattevans/postmark-go"
	"golang.org/x/crypto/bcrypt"
)

/*
Page is a generic HTML page with
only the necessary data.
*/
type Page struct {
	Title string
	Paths []helpers.Path
}

/*
SignOut endpoint removes the
username from the session.
*/
func SignOut(w http.ResponseWriter, r *http.Request) {
	sessions, err := helpers.Store.Get(r, "session")
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	// Revoke users authentication
	sessions.Values["username"] = nil
	sessions.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// SignIn endpoint displays a HTML page.
func SignIn(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Sign in", Paths: []helpers.Path{{Name: "Sign in", Link: "/sign/in"}}}
	helpers.RenderTemplate(r, w, "signin", p)
}

/*
GetPasswordHashFromUsername returns the password hash
from the database associated with the username sent in.
*/
func GetPasswordHashFromUsername(db *sql.DB, username string) (string, error) {
	var passwordHash string
	query := `SELECT password FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&passwordHash)
	if err != nil {
		return passwordHash, err
	}
	return passwordHash, nil
}

/*
SetSessionUsername sets the session username for
the user to the username variable sent in.
*/
func SetSessionUsername(w http.ResponseWriter, r *http.Request, username string) error {
	sessions, err := helpers.Store.Get(r, "session")
	if err != nil {
		return err
	}

	sessions.Values["username"] = username
	err = sessions.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

/*
Authenticating endpoint parses the form data the user
sent in and authenticates the user if it is valid.
*/
func Authenticating(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	passwordHash, err := GetPasswordHashFromUsername(db, email)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		// Passwords don't match
		http.Redirect(w, r, "/sign/in", http.StatusSeeOther)
	} else {
		err = SetSessionUsername(w, r, email)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// SignUp endpoint returns the Sign up HTML page.
func SignUp(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Sign up", Paths: []helpers.Path{{Name: "Sign up", Link: "/sign/up"}}}
	helpers.RenderTemplate(r, w, "signup", p)
}

/*
VerificationMail contains the data to
render the verification mail HTML page.
*/
type VerificationMail struct {
	Host            string
	ActivationToken string
}

/*
GetMailSender gets the MailSender value in
the settings table in the database.
*/
func GetMailSender(db *sql.DB) (string, error) {
	var mailSender string
	query := `SELECT option_value FROM settings WHERE option_name = 'MailSender'`
	err := db.QueryRow(query).Scan(&mailSender)
	if err != nil {
		// No user
		return "", err
	}

	return mailSender, nil
}

/*
GetPostmarkClient gets the PostmarkToken value in
the settings table in the database.
*/
func GetPostmarkClient(db *sql.DB) (*postmark.Client, error) {
	var postmarkToken string
	query := `SELECT option_value FROM settings WHERE option_name = 'PostmarkToken'`
	err := db.QueryRow(query).Scan(&postmarkToken)
	if err != nil {
		return nil, err
	}

	auth := &http.Client{
		Transport: &postmark.AuthTransport{Token: postmarkToken},
	}
	client := postmark.NewClient(auth)

	return client, nil
}

/*
GetRegistrationMailSubject gets the RegistrationMailSubject value in
the settings table in the database.
*/
func GetRegistrationMailSubject(db *sql.DB) (string, error) {
	var registrationMailSubject string
	query := `SELECT option_value FROM settings WHERE option_name = 'RegistrationMailSubject'`
	err := db.QueryRow(query).Scan(&registrationMailSubject)
	if err != nil {
		return "", err
	}
	return registrationMailSubject, nil
}

/*
GetHost gets the SiteUrl value in
the settings table in the database.
*/
func GetHost(db *sql.DB) (string, error) {
	var host string
	query := `SELECT option_value FROM settings WHERE option_name = 'SiteUrl'`
	err := db.QueryRow(query).Scan(&host)
	if err != nil {
		return "", err
	}
	return host, nil
}

/*
EmailIsValid checks if the email is valid and returns it in boolean form.
*/
func EmailIsValid(email string) bool {
	const emailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

/*
InsertNewUser inserts the user into the database.
*/
func InsertNewUser(db *sql.DB, email string, username string, password string) (int64, error) {
	res, err := db.Exec(`INSERT INTO users (email, username, password) VALUES (?, ?, ?)`, email, username, password)
	if err != nil {
		return 0, err
	}
	userID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return userID, nil
}

/*
InsertActivationToken inserts the activation token into the database
*/
func InsertActivationToken(db *sql.DB, userID int64, activationToken string) error {
	_, err := db.Exec(`INSERT INTO email_verification (user_id, token) VALUES (?, ?)`, userID, activationToken)
	if err != nil {
		return err
	}
	return nil
}

/*
SignUpNewUser endpoint parses the form for signing up and inserts
the user into the database and sends a verification email.
*/
func SignUpNewUser(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()
	email := r.FormValue("email")
	if !EmailIsValid(email) {
		helpers.HandleError(errors.New("Email is not valid"))
	}

	password := r.FormValue("password")
	username := r.FormValue("username")
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	userID, err := InsertNewUser(db, email, username, string(passwordBytes))
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	mailSender, err := GetMailSender(db)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	client, err := GetPostmarkClient(db)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	registrationMailSubject, err := GetRegistrationMailSubject(db)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	host, err := GetHost(db)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	activationToken, err := helpers.GenerateRandomStringURLSafe(32)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	err = InsertActivationToken(db, userID, activationToken)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	verificationMail := VerificationMail{
		ActivationToken: activationToken,
		Host:            host,
	}

	tmpl := template.Must(template.ParseFiles("./templates/mail/registration_mail.html"))

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, verificationMail)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	result := tpl.String()

	emailReq := &postmark.Email{
		From:       mailSender,
		To:         email,
		Subject:    registrationMailSubject,
		HTMLBody:   result,
		TrackOpens: true,
	}

	_, _, err = client.Email.Send(emailReq)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*
GetUserIDFromUsername gets the userID from database
that matches the username variable sent in.
*/
func GetUserIDFromUsername(db *sql.DB, username string) (int, error) {
	var userID int
	query := `SELECT id FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
