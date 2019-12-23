package authentication

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/initiate"
	_ "github.com/go-sql-driver/mysql"
)

func TestEmailIsValid(t *testing.T) {
	if EmailIsValid("hejhej") {
		t.Errorf("EmailIsValid failed. 'hejhej' is not a valid email")
	}

	if !EmailIsValid("hejhej@hej.se") {
		t.Errorf("EmailIsValid failed. 'hejhej@hej.se' is a valid email")
	}
}

func TestGetUserIDFromUsername(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO users (id, username, email, password, verified, created_at)
	VALUES
		(1, 'johndoe', 'johndoe@gmail.com', '$2a$10$rNQZA0Ibm7RyXUS1aFCPCe5SGS8L/1aTev2Ej6g5MFLYQPXLYJ1T6', 1, '2019-11-11 10:11:18');		
	`)
	if err != nil {
		t.Errorf("Insertion of session in database failed. Error message: %s", err.Error())
		return
	}

	id, err := GetUserIDFromUsername(db, "johndoe")
	if err != nil {
		t.Errorf("GetSessions failed. Error message: %s", err.Error())
		return
	}

	if id != 1 {
		t.Errorf("GetUserIDFromUsername failed. Pages are not equal")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestInsertActivationToken(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = InsertActivationToken(db, 1, "hejhej")
	if err != nil {
		t.Errorf("InsertActivationToken failed. Error message %s", err.Error())
		return
	}

	var token string
	query := `SELECT token FROM email_verification WHERE user_id = 1`
	err = db.QueryRow(query).Scan(&token)
	if err != nil {
		t.Errorf("InsertActivationToken failed. Error message: %s", err.Error())
		return
	}

	if token != "hejhej" {
		t.Errorf("InsertActivationToken failed. Pages are not equal")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestInsertNewUser(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	id, err := InsertNewUser(db, "email", "username", "password")
	if err != nil {
		t.Errorf("InsertNewUser failed. Error message %s", err.Error())
		return
	}

	if id != 1 {
		t.Errorf("InsertNewUser failed. ID is not correct.")
		return
	}

	var email string
	var username string
	var password string
	query := `SELECT email, username, password FROM users WHERE id = ?`
	err = db.QueryRow(query, id).Scan(&email, &username, &password)
	if err != nil {
		t.Errorf("InsertNewUser failed. Error message: %s", err.Error())
		return
	}

	if email != "email" {
		t.Errorf("InsertNewUser failed. Pages are not equal")
		return
	}

	if username != "username" {
		t.Errorf("InsertNewUser failed. Pages are not equal")
		return
	}

	if password != "password" {
		t.Errorf("InsertNewUser failed. Pages are not equal")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestGetHost(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	host, err := GetHost(db)
	if err != nil {
		t.Errorf("GetHost failed. Error message %s", err.Error())
		return
	}

	if host != "" {
		t.Errorf("GetHost failed.")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestGetRegistrationMailSubject(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	subject, err := GetRegistrationMailSubject(db)
	if err != nil {
		t.Errorf("GetRegistrationMailSubject failed. Error message %s", err.Error())
		return
	}

	if subject != "" {
		t.Errorf("GetRegistrationMailSubject failed.")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestGetPostmarkClient(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	client, err := GetPostmarkClient(db)
	if err != nil {
		t.Errorf("GetPostmarkClient failed. Error message %s", err.Error())
		return
	}

	if client.Token != "" {
		t.Errorf("GetPostmarkClient failed.")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestGetMailSender(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	sender, err := GetMailSender(db)
	if err != nil {
		t.Errorf("GetMailSender failed. Error message %s", err.Error())
		return
	}

	if sender != "" {
		t.Errorf("GetMailSender failed.")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestGetPasswordHashFromUsername(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO users (id, username, email, password, verified, created_at)
	VALUES
		(1, 'johndoe', 'johndoe@gmail.com', '1aTev2Ej6g5MFLYQPXLYJ1T6', 1, '2019-11-11 10:11:18');		
	`)
	if err != nil {
		t.Errorf("Insertion of session in database failed. Error message: %s", err.Error())
		return
	}

	hash, err := GetPasswordHashFromUsername(db, "johndoe")
	if err != nil {
		t.Errorf("GetPasswordHashFromUsername failed. Error message: %s", err.Error())
		return
	}

	if hash != "1aTev2Ej6g5MFLYQPXLYJ1T6" {
		t.Errorf("GetPasswordHashFromUsername failed. Pages are not equal")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestSignIn(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SignIn)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestAuthenticating(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Authenticating)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestSignUp(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SignUp)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestDBForAuthenticating(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Authenticating)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := ""
	if rr.Body.String() == expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	err = initiate.WriteSettingsFile(config)
	if err != nil {
		t.Fatal(err)
	}

}

func TestDBForSignUpNewUser(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SignUpNewUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := ""
	if rr.Body.String() == expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	err = initiate.WriteSettingsFile(config)
	if err != nil {
		t.Fatal(err)
	}

}
