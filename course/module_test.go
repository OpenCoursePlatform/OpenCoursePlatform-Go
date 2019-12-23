package course

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/initiate"

	_ "github.com/go-sql-driver/mysql"
)

func TestGetModules(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO module (id, name, description, image_link, slug, course_id)
	VALUES
		(1, 'Getting started with Python', 'Getting started with Python', 'link', 'getting-started-with-python', 1),
		(2, 'Python Data Structures', 'Python Data Structures', 'link', 'python-data-structures', 1);
	`)
	if err != nil {
		t.Errorf("Insertion of module in database failed. Error message: %s", err.Error())
		return
	}

	modules, err := GetModules(db)
	if err != nil {
		t.Errorf("GetModules failed. Error message: %s", err.Error())
		return
	}
	sampleModules := []Module{
		{"Getting started with Python", "", "link", "getting-started-with-python"},
		{"Python Data Structures", "", "link", "python-data-structures"},
	}

	if sampleModules[0] != modules[0] {
		t.Errorf("GetModules failed. Structs are not equal")
		return
	}

	if sampleModules[1] != modules[1] {
		t.Errorf("GetModules failed. Structs are not equal")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestGetSessionsByModuleSlug(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO module (id, name, description, slug, course_id)
	VALUES
		(1, 'Getting started with Python', 'Getting started with Python', 'getting-started-with-python', 1);

		`)
	if err != nil {
		t.Errorf("Insertion of module in database failed. Error message: %s", err.Error())
		return
	}
	_, err = db.Exec(`
		INSERT INTO session (id, name, slug, module_id, session_type)
		VALUES
			(1, 'Introduction!', 'introduction', 1, 0);
	`)
	if err != nil {
		t.Errorf("Insertion of session in database failed. Error message: %s", err.Error())
		return
	}

	sessions, err := GetSessionsByModuleSlug(db, "getting-started-with-python")
	if err != nil {
		t.Errorf("GetSessionsByModuleSlug failed. Error message: %s", err.Error())
		return
	}

	sampleSessions := []Session{
		{Name: "Introduction!", Slug: "introduction"},
	}

	if sessions[0] != sampleSessions[0] {
		fmt.Println(sessions[0])
		t.Errorf("GetSessionsByModuleSlug failed. Structs are not equal")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestModulePage(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ModulePage)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestDBForModulePage(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ModulePage)

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
