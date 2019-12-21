package support

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/authentication"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/initiate"

	_ "github.com/go-sql-driver/mysql"
)

func TestDBForSupport(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Support)

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

func TestDBForInsertNewTicket(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(InsertNewTicket)

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

func TestDBForTicket(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Ticket)

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

func TestDBForInsertNewTicketResponse(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(InsertNewTicketResponse)

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

func TestGetTicketsByUsername(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	DROP TABLE users
	`)
	if err != nil {
		t.Errorf("Deletion of users table in database failed. Error message: %s", err.Error())
		return
	}

	tickets, err := GetTicketsByUsername(db, "johndoe")
	if err == nil {
		t.Errorf("GetUserSettings should have failed.")
		return
	}

	err = initiate.CreateUsers(db)
	if err != nil {
		t.Errorf("Creation of users table in database failed. Error message: %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO users (id, username, email, password, verified, created_at)
	VALUES
		(1, 'johndoe', 'johndoe@gmail.com', '$2a$10$rNQZA0Ibm7RyXUS1aFCPCe5SGS8L/1aTev2Ej6g5MFLYQPXLYJ1T6', 0, '2019-11-11 10:11:18');
	`)
	if err != nil {
		t.Errorf("Insertion of user in database failed. Error message: %s", err.Error())
		return
	}

	_, err = db.Exec(`
	DROP TABLE tickets
	`)
	if err != nil {
		t.Errorf("Deletion of tickets table in database failed. Error message: %s", err.Error())
		return
	}

	_, err = db.Exec(`
	CREATE TABLE tickets (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		user_id int(11) unsigned NOT NULL,
		topic varchar(255) NOT NULL,
		slug varchar(255) NOT NULL,
		solved varchar(255) NOT NULL,
		created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY slug (slug),
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		t.Errorf("Creation of tickets table in database failed. Error message: %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO tickets (id, user_id, topic, slug, solved, created)
	VALUES
		(1, 1, 'Can\'t log in', 'cant-log-in', 'testing', '2019-11-20 16:05:03');
	
	`)
	if err != nil {
		t.Errorf("Insertion of ticket in database failed. Error message: %s", err.Error())
		return
	}

	tickets, err = GetTicketsByUsername(db, "johndoe")
	if err == nil {
		t.Errorf("GetTicketsByUsername should have failed.")
		return
	}

	_, err = db.Exec(`
	DROP TABLE tickets
	`)
	if err != nil {
		t.Errorf("Deletion of tickets table in database failed. Error message: %s", err.Error())
		return
	}

	err = initiate.CreateTickets(db)
	if err != nil {
		t.Errorf("Creation of tickets table in database failed. Error message: %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO tickets (id, user_id, topic, slug, solved, created)
	VALUES
		(1, 1, 'Can\'t log in', 'cant-log-in', 1, '2019-11-20 16:05:03');
	
	`)
	if err != nil {
		t.Errorf("Insertion of ticket in database failed. Error message: %s", err.Error())
		return
	}
	_, err = db.Exec(`
	INSERT INTO ticket_responses (id, ticket_id, user_id, text, created)
	VALUES
	(1, 1, 1, 'Hey man I can\'t log in', '2019-11-20 16:05:49');
	`)
	if err != nil {
		t.Errorf("Insertion of ticket in database failed. Error message: %s", err.Error())
		return
	}
	tickets, err = GetTicketsByUsername(db, "johndoe")
	if err != nil {
		t.Errorf("GetUserSettings failed. Error message: %s", err.Error())
		return
	}

	ticketsCopy := []TicketData{{Title: "Can't log in", Slug: "cant-log-in", Solved: true}}

	if tickets[0] != ticketsCopy[0] {
		t.Errorf("Data is not correct. Does not match test data.")
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestInsertNewTicketInDB(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO users (id, username, email, password, verified, created_at)
	VALUES
		(1, 'johndoe', 'johndoe@gmail.com', '$2a$10$rNQZA0Ibm7RyXUS1aFCPCe5SGS8L/1aTev2Ej6g5MFLYQPXLYJ1T6', 0, '2019-11-11 10:11:18');
	`)
	if err != nil {
		t.Errorf("Insertion of user in database failed. Error message: %s", err.Error())
		return
	}

	_, err = InsertNewTicketInDB(db, "johndoe", "Can't log in", "Hey man I can't log in")
	if err != nil {
		t.Errorf("InsertNewTicketInDB failed. Error message: %s", err.Error())
		return
	}
	tickets, err := GetTicketsByUsername(db, "johndoe")
	if err != nil {
		t.Errorf("GetUserSettings failed. Error message: %s", err.Error())
		return
	}

	ticketsCopy := []TicketData{{Title: "Can't log in", Slug: "can-t-log-in", Solved: false}}

	if tickets[0] != ticketsCopy[0] {
		t.Errorf("Data is not correct. Does not match test data.")
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestGetTicketBySlug(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO users (id, username, email, password, verified, created_at)
	VALUES
		(1, 'johndoe', 'johndoe@gmail.com', '$2a$10$rNQZA0Ibm7RyXUS1aFCPCe5SGS8L/1aTev2Ej6g5MFLYQPXLYJ1T6', 0, '2019-11-11 10:11:18');
	`)
	if err != nil {
		t.Errorf("Insertion of user in database failed. Error message: %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO tickets (id, user_id, topic, slug, solved, created)
	VALUES
		(1, 1, 'Can\'t log in', 'cant-log-in', 1, '2019-11-20 16:05:03');
	
	`)
	if err != nil {
		t.Errorf("Insertion of ticket in database failed. Error message: %s", err.Error())
		return
	}
	_, err = db.Exec(`
	INSERT INTO ticket_responses (id, ticket_id, user_id, text, created)
	VALUES
	(1, 1, 1, 'Hey man I can\'t log in', '2019-11-20 16:05:49');
	`)
	if err != nil {
		t.Errorf("Insertion of ticket in database failed. Error message: %s", err.Error())
		return
	}
	id, tickets, err := GetTicketBySlug(db, "cant-log-in", "johndoe")
	if err != nil {
		t.Errorf("GetUserSettings failed. Error message: %s", err.Error())
		return
	}

	ticketsCopy := TicketData{Title: "Can't log in", Slug: "cant-log-in", Solved: true}

	if tickets != ticketsCopy {
		t.Errorf("Data is not correct. Does not match test data.")
	}

	if id != 1 {
		t.Errorf("Data is not correct. Does not match test data.")
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestGetTicketResponsesByTicketID(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO users (id, username, email, password, verified, created_at)
	VALUES
		(1, 'johndoe', 'johndoe@gmail.com', '$2a$10$rNQZA0Ibm7RyXUS1aFCPCe5SGS8L/1aTev2Ej6g5MFLYQPXLYJ1T6', 0, '2019-11-11 10:11:18');
	`)
	if err != nil {
		t.Errorf("Insertion of user in database failed. Error message: %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO tickets (id, user_id, topic, slug, solved, created)
	VALUES
		(1, 1, 'Can\'t log in', 'cant-log-in', 1, '2019-11-20 16:05:03');
	
	`)
	if err != nil {
		t.Errorf("Insertion of ticket in database failed. Error message: %s", err.Error())
		return
	}
	_, err = db.Exec(`
	INSERT INTO ticket_responses (id, ticket_id, user_id, text, created)
	VALUES
	(1, 1, 1, 'Hey man I can\'t log in', '2019-11-20 16:05:49');
	`)
	if err != nil {
		t.Errorf("Insertion of ticket in database failed. Error message: %s", err.Error())
		return
	}
	responses, err := GetTicketResponsesByTicketID(db, 1)
	if err != nil {
		t.Errorf("GetUserSettings failed. Error message: %s", err.Error())
		return
	}

	responsesCopy := []string{"Hey man I can't log in"}

	if responses[0] != responsesCopy[0] {
		t.Errorf("Data is not correct. Does not match test data.")
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestInsertNewTicketResponseInDB(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}
	err = InsertNewTicketResponseInDB(db, 1, 1, "Testing 1")
	if err != nil {
		t.Errorf("InsertNewTicketResponseInDB failed. Error message: %s", err.Error())
		return
	}
	err = InsertNewTicketResponseInDB(db, 1, 1, "Testing 2")
	if err != nil {
		t.Errorf("InsertNewTicketResponseInDB failed. Error message: %s", err.Error())
		return
	}
	err = InsertNewTicketResponseInDB(db, 1, 1, "Testing 3")
	if err != nil {
		t.Errorf("InsertNewTicketResponseInDB failed. Error message: %s", err.Error())
		return
	}
	err = InsertNewTicketResponseInDB(db, 1, 1, "Testing 4")
	if err != nil {
		t.Errorf("InsertNewTicketResponseInDB failed. Error message: %s", err.Error())
		return
	}
	responses, err := GetTicketResponsesByTicketID(db, 1)
	if err != nil {
		t.Errorf("GetTicketResponsesByTicketID failed. Error message: %s", err.Error())
		return
	}

	responsesCopy := []string{"Testing 1", "Testing 2", "Testing 3", "Testing 4"}

	for index := range responses {
		if responses[index] != responsesCopy[index] {
			t.Errorf("Data is not correct. Does not match test data.")
		}
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestSupportFailure(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Support)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestSupport(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	name, err := initiate.UpdateSettingToTestDB()
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Support)

	err = authentication.SetSessionUsername(rr, req, "axel")
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}

	err = initiate.UpdateDatabaseName(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewTicket(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NewTicket)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestInsertNewTicket(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(InsertNewTicket)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestTicket(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Ticket)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestInsertNewTicketResponse(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(InsertNewTicketResponse)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}
