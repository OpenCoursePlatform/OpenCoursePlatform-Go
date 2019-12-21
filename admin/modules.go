package admin

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/course"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

// SessionData ...
type SessionData struct {
	Name string
	Slug string
}

// ModulePageData ...
type ModulePageData struct {
	Title         string
	Paths         []helpers.Path
	Module        ModuleData
	CourseSlug    string
	Courses       []course.Course
	Sessions      []SessionData
	NewItemButton bool
	NewItemText   string
	NewItemLink   string
}

/*
GetModuleAndCourseBySlug ...
*/
func GetModuleAndCourseBySlug(db *sql.DB, slug string) (ModuleData, error) {
	var module ModuleData

	query := `
	SELECT module.name, module.description, module.slug, courses.name
	FROM module
	JOIN courses
		ON courses.id = module.course_id
	WHERE module.slug = ?`
	err := db.QueryRow(query, slug).Scan(&module.Name, &module.Description, &module.Slug, &module.Course)
	if err != nil {
		return module, err
	}
	return module, nil
}

/*
GetCourses ...
*/
func GetCourses(db *sql.DB) ([]course.Course, error) {
	rows, err := db.Query(`
	SELECT name, slug
	FROM courses`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []course.Course
	for rows.Next() {
		var course course.Course
		err = rows.Scan(&course.Name, &course.Slug)
		if err != nil {
			return courses, err
		}
		courses = append(courses, course)
	}
	err = rows.Err()
	if err != nil {
		return courses, err
	}
	return courses, nil
}

/*
GetSessionsByModuleSlug ...
*/
func GetSessionsByModuleSlug(db *sql.DB, slug string) ([]SessionData, error) {
	rows, err := db.Query(`
	SELECT session.name, session.slug
	FROM session
	INNER JOIN module
		ON module.id = session.module_id
	WHERE module.slug = ?`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []SessionData
	for rows.Next() {
		var session SessionData
		err = rows.Scan(&session.Name, &session.Slug)
		if err != nil {
			return sessions, err
		}
		sessions = append(sessions, session)
	}
	err = rows.Err()
	if err != nil {
		return sessions, err
	}
	return sessions, nil
}

// Module ...
func Module(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["module"]
	course := vars["course"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	module, err := GetModuleAndCourseBySlug(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	courses, err := GetCourses(db)
	if err != nil {
		helpers.HandleError(err)
	}

	sessions, err := GetSessionsByModuleSlug(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &ModulePageData{module.Name, []helpers.Path{{Name: "Admin", Link: "/admin"}}, module, course, courses, sessions, true, "New Session", "/courses/" + course + "/" + slug + "/new"}
	helpers.RenderTemplate(r, w, "admin_module", p)
}

/*
GetCourseFromName ...
*/
func GetCourseFromName(db *sql.DB, name string) (int, string, error) {
	var courseID int
	var courseSlug string

	query := `SELECT id, slug FROM courses WHERE name = ?`
	err := db.QueryRow(query, name).Scan(&courseID, &courseSlug)
	if err != nil {
		return courseID, courseSlug, err
	}
	return courseID, courseSlug, nil

}

/*
UpdateModuleInDB ...
*/
func UpdateModuleInDB(db *sql.DB, name, description, slug string, courseID int) error {
	insForm, err := db.Prepare("UPDATE module SET name=?, description=?, course_id=? WHERE slug=?")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(name, description, courseID, slug)
	if err != nil {
		return err
	}
	return nil
}

// UpdateModule ...
func UpdateModule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["module"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	course := r.FormValue("course")

	courseID, courseSlug, err := GetCourseFromName(db, course)
	if err != nil {
		helpers.HandleError(err)
	}

	err = UpdateModuleInDB(db, r.FormValue("name"), r.FormValue("description"), slug, courseID)
	if err != nil {
		helpers.HandleError(err)
	}

	http.Redirect(w, r, "/admin/courses/"+courseSlug+"/"+slug, http.StatusSeeOther)
}

// NewModule ...
func NewModule(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "New module", Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}}}
	helpers.RenderTemplate(r, w, "admin_new_module", p)
}

/*
InsertNewModuleInDB ...
*/
func InsertNewModuleInDB(db *sql.DB, name, description, course string) (string, error) {
	slug := helpers.GenerateSlug(name)
	var courseID int
	query := `
	SELECT courses.id
	FROM courses
	WHERE courses.slug = ?`
	err := db.QueryRow(query, course).Scan(&courseID)
	if err != nil {
		return "", err
	}
	insForm, err := db.Prepare("INSERT INTO module (name, description, slug, course_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	res, err := insForm.Exec(name, description, slug, courseID)
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	var module string
	query = `
	SELECT module.slug
	FROM module
	WHERE module.id = ?`
	err = db.QueryRow(query, id).Scan(&module)
	if err != nil {
		return "", err
	}
	return module, nil
}

// InsertNewModule ...
func InsertNewModule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	course := vars["course"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	name := r.FormValue("name")
	description := r.FormValue("description")
	module, err := InsertNewModuleInDB(db, name, description, course)
	if err != nil {
		helpers.HandleError(err)
	}
	http.Redirect(w, r, "/admin/courses/"+course+"/"+module, http.StatusSeeOther)
}
