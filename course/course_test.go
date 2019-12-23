package course

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/initiate"

	_ "github.com/go-sql-driver/mysql"
)

func TestGetCourses(t *testing.T) {
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

	courses, err := GetCourses(db)
	if err != nil {
		t.Errorf("GetCourses failed. Error message: %s", err.Error())
		return
	}
	sampleCourses := []Course{
		{"Fundamentals of Computing", "Fundamentals of Computing", "fundamentals-of-computing", "Computer Science"},
		{"Simple Data", "Simple Data", "simple-data", "Computer Science"},
	}

	if sampleCourses[0] != courses[0] {
		t.Errorf("GetCourses failed. Structs are not equal")
		return
	}

	if sampleCourses[0] != courses[0] {
		t.Errorf("GetCourses failed. Structs are not equal")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestGetModulesByCourseSlug(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO courses (id, name, description, slug, category_id)
	VALUES
		(1, 'Introduction to Programming', 'Introduction to Programming', 'introduction-to-programming', 1);
	`)
	if err != nil {
		t.Errorf("Insertion of courses in database failed. Error message: %s", err.Error())
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

	modules, err := GetModulesByCourseSlug(db, "introduction-to-programming")
	if err != nil {
		t.Errorf("GetModulesByCourseSlug failed. Error message: %s", err.Error())
		return
	}
	sampleModules := []Module{
		{Name: "Getting started with Python", Description: "Getting started with Python", ImageLink: "link", Slug: "getting-started-with-python"},
		{Name: "Python Data Structures", Description: "Python Data Structures", ImageLink: "link", Slug: "python-data-structures"},
	}

	if sampleModules[0] != modules[0] {
		t.Errorf("GetModulesByCourseSlug failed. Structs are not equal")
		return
	}

	if sampleModules[1] != modules[1] {
		t.Errorf("GetModulesByCourseSlug failed. Structs are not equal")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestPage(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
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
}

func TestDBForPage(t *testing.T) {
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
