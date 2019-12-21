package index

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/initiate"
	_ "github.com/go-sql-driver/mysql"
)

func TestDB(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Page)

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

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Page)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	templates := template.Must(helpers.ParseTemplates(), nil)
	var buf bytes.Buffer

	err = templates.ExecuteTemplate(&buf, "index.html", nil)
	if err == nil {
		t.Error("ExecuteTemplate should have failed.")
	}
	expected := buf.String()
	if rr.Body.String() == expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestIndex(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}
	_, err = db.Exec(`
	INSERT INTO courses (id, name, description, slug, category_id)
	VALUES
		(2, 'Fundamentals of Computing', 'Fundamentals of Computing', 'fundamentals-of-computing', 2),
		(3, 'Simple Data', 'Simple Data', 'simple-data', 2);	
	`)
	if err != nil {
		t.Errorf("Insertion of courses in database failed. Error message: %s", err.Error())
		return
	}

	_, err = db.Exec(`
		INSERT INTO course_categories (id, name, slug)
		VALUES
		(2, 'Computer Science', 'computer-science');
	`)
	if err != nil {
		t.Errorf("Insertion of course categories in database failed. Error message: %s", err.Error())
		return
	}
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	name, err := initiate.UpdateSettingToTestDB()
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Page)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	templates := template.Must(helpers.ParseTemplates(), nil)
	var buf bytes.Buffer

	err = templates.ExecuteTemplate(&buf, "index.html", nil)
	if err == nil {
		t.Error("ExecuteTemplate should have failed.")
	}
	expected := buf.String()
	if rr.Body.String() == expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
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
