package course

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

/*
Module struct stores
module data for rendering.
*/
type Module struct {
	Name        string
	Description string
	ImageLink   string
	Slug        string
}

/*
ModulePageData stores all data for
rendering the module HTML page.
*/
type ModulePageData struct {
	Title      string
	ModuleSlug string
	Course     Course
	// Modules in sidebar
	Module   Module
	Modules  []Module
	Sessions []Session
	Paths    []helpers.Path
}

/*
GetModules gets all of the modules from the database.
*/
func GetModules(db *sql.DB, slug string) ([]Module, error) {
	rows, err := db.Query(`
	SELECT module.name, module.image_link, module.slug
	FROM module
	INNER JOIN courses
		ON module.course_id = courses.id
	WHERE courses.slug = ?`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []Module
	for rows.Next() {
		var module Module
		err = rows.Scan(&module.Name, &module.ImageLink, &module.Slug)
		if err != nil {
			return modules, err
		}
		modules = append(modules, module)
	}
	err = rows.Err()
	if err != nil {
		return modules, err
	}
	return modules, err
}

/*
GetModule gets the current module from the database.
*/
func GetModule(db *sql.DB, slug string) (Module, error) {
	var module Module
	err := db.QueryRow(`
		SELECT name, image_link, slug
		FROM module
		WHERE slug = ?
	`, slug).Scan(&module.Name, &module.ImageLink, &module.Slug) // check err
	if err != nil {
		return module, err
	}

	return module, nil
}

/*
GetSessionsByModuleSlug gets all of the sessions related
to a module by the module slug from the database.
*/
func GetSessionsByModuleSlug(db *sql.DB, moduleSlug, courseSlug string) ([]Session, error) {
	rows, err := db.Query(`
		SELECT session.name, session.slug
		FROM session
		JOIN module
			ON session.module_id = module.id
		JOIN courses
			ON module.course_id = courses.id
		WHERE module.slug = ?
		AND courses.slug = ?
	`, moduleSlug, courseSlug) // check err
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var session Session
		err = rows.Scan(&session.Name, &session.Slug) // check err
		if err != nil {
			return sessions, err
		}
		sessions = append(sessions, session)
	}
	err = rows.Err() // check err
	if err != nil {
		return sessions, err
	}
	return sessions, nil
}

/*
ModulePage endpoint returns the module page.
*/
func ModulePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["module"]
	courseSlug := vars["course"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	course, err := GetCourse(db, courseSlug)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	module, err := GetModule(db, slug)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	modules, err := GetModules(db, courseSlug)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	sessions, err := GetSessionsByModuleSlug(db, slug, courseSlug)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
	}

	coursePage := ModulePageData{
		Title:      module.Name,
		ModuleSlug: slug,
		Course:     course,
		Module:     module,
		Modules:    modules,
		Sessions:   sessions,
		Paths:      []helpers.Path{{Name: "Courses", Link: "/courses"}, {Name: course.Name, Link: "/courses/" + course.Slug}},
	}

	helpers.RenderTemplate(r, w, "module", coursePage)
}
