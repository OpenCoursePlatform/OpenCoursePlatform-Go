package settings

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/authentication"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/initiate"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

func TestDBForUserPage(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UserPage)

	err = authentication.SetSessionUsername(rr, req, "axel")
	if err != nil {
		t.Fatal(err)
	}

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

func TestDBForEmailPage(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(EmailPage)

	err = authentication.SetSessionUsername(rr, req, "axel")
	if err != nil {
		t.Fatal(err)
	}

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

func TestDBForUpdateEmail(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateEmail)

	err = authentication.SetSessionUsername(rr, req, "axel")
	if err != nil {
		t.Fatal(err)
	}

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

func TestDBForUpdatePassword(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdatePassword)

	err = authentication.SetSessionUsername(rr, req, "axel")
	if err != nil {
		t.Fatal(err)
	}

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

func TestGetUserSettings(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO users (id, username, email, password, verified, created_at)
	VALUES (1, 'johndoe', 'johndoe@gmail.com', '$2a$10$rNQZA0Ibm7RyXUS1aFCPCe5SGS8L/1aTev2Ej6g5MFLYQPXLYJ1T6', 0, '2019-11-11 10:11:18');
	`)
	if err != nil {
		t.Errorf("Insertion of user in database failed. Error message: %s", err.Error())
		return
	}

	user, err := GetUserSettings(db, "johndoe", "testing")
	if err != nil {
		t.Errorf("GetUserSettings failed. Error message: %s", err.Error())
		return
	}

	createdAt, err := time.Parse("2006-01-02 15:04:05", "2019-11-11 10:11:18")
	if err != nil {
		t.Errorf("Parsing time failed. If it fails it's probably due to the time package. Error message: %s", err.Error())
		return
	}

	userCopy := UserPageData{Title: "testing", Username: "johndoe", UserEmail: "johndoe@gmail.com", Registered: createdAt.Format("January 02, 2006")}

	if user != userCopy {
		t.Errorf("Data is not correct. Error message: %s", err.Error())
	}

	_, err = GetUserSettings(db, "johndoeeee", "testing")
	if err == nil {
		t.Errorf("GetUserSettings failed. User should not have been found but was.")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestUpdateUserEmail(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO users (id, username, email, password, verified, created_at)
	VALUES (1, 'johndoe', 'johndoe@gmail.com', '$2a$10$rNQZA0Ibm7RyXUS1aFCPCe5SGS8L/1aTev2Ej6g5MFLYQPXLYJ1T6', 0, '2019-11-11 10:11:18');
	`)
	if err != nil {
		t.Errorf("Insertion of user in database failed. Error message: %s", err.Error())
		return
	}

	createdAt, err := time.Parse("2006-01-02 15:04:05", "2019-11-11 10:11:18")
	if err != nil {
		t.Errorf("Parsing time failed. If it fails it's probably due to the time package. Error message: %s", err.Error())
		return
	}

	err = UpdateUserEmail(db, "testing@test.com", "johndoe")
	if err != nil {
		t.Errorf("UpdateUserEmail failed. Error message: %s", err.Error())
		return
	}

	err = UpdateUserEmail(db, "testing@test.com", "johndoeeeee")
	if err == nil {
		t.Errorf("UpdateUserEmail failed. User should not have been found but was.")
		return
	}

	user, err := GetUserSettings(db, "johndoe", "testing")
	if err != nil {
		t.Errorf("GetUserSettings failed. Error message: %s", err.Error())
		return
	}

	userCopy := UserPageData{Title: "testing", Username: "johndoe", UserEmail: "testing@testing.com", Registered: createdAt.Format("January 02, 2006")}

	if user == userCopy {
		t.Errorf("Users do not match.")
	}

	userCopy = UserPageData{Title: "testing", Username: "johndoe", UserEmail: "testing@test.com", Registered: createdAt.Format("January 02, 2006")}
	if user != userCopy {
		t.Errorf("Users do not match.")
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestUpdateUserPassword(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO users (id, username, email, password, verified, created_at)
	VALUES (1, 'johndoe', 'johndoe@gmail.com', '$2a$10$rNQZA0Ibm7RyXUS1aFCPCe5SGS8L/1aTev2Ej6g5MFLYQPXLYJ1T6', 1, '2019-11-11 10:11:18');
	`)
	if err != nil {
		t.Errorf("Insertion of user in database failed. Error message: %s", err.Error())
		return
	}

	err = UpdateUserPassword(db, "testing", "johndoe")
	if err != nil {
		t.Errorf("UpdateUserPassword failed. Error message: %s", err.Error())
		return
	}

	passwordHash, err := authentication.GetPasswordHashFromUsername(db, "johndoe")
	if err != nil {
		t.Errorf("GetPasswordHashFromUsername failed. Error message: %s", err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte("testing"))
	if err != nil {
		t.Errorf("Passwords don't match. Test failed. Error message: %s", err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte("testinging"))
	if err == nil {
		t.Errorf("Passwords should not match but they do. Test failed.")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestUserPage(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UserPage)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestEmailPage(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(EmailPage)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestUpdateEmail(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateEmail)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestPasswordPage(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PasswordPage)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestUpdatePassword(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdatePassword)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}
