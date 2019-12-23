package course

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

// Course contains the data for a single course
type Course struct {
	Name        string
	Description string
	Slug        string
	Category    string
}

// PageData contains the data to render the HTML course page.
type PageData struct {
	Title      string
	CourseSlug string
	Courses    []Course
	Course     Course
	Modules    []Module
	Paths      []helpers.Path
}

/*
GetCourses gets all of the courses from the database.
*/
func GetCourses(db *sql.DB) ([]Course, error) {
	rows, err := db.Query(`
		SELECT courses.name, courses.description, courses.slug, course_categories.name
		FROM courses
		JOIN course_categories
			ON courses.category_id = course_categories.id
	`) // check err
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		err = rows.Scan(&course.Name, &course.Description, &course.Slug, &course.Category) // check err
		if err != nil {
			return courses, err
		}
		courses = append(courses, course)
	}
	err = rows.Err() // check err
	if err != nil {
		return courses, err
	}
	return courses, nil
}

/*
GetCourse gets the current course from the database.
*/
func GetCourse(db *sql.DB, slug string) (Course, error) {
	var course Course
	err := db.QueryRow(`
		SELECT name, description, slug, name
		FROM courses
		WHERE slug = ?
	`, slug).Scan(&course.Name, &course.Description, &course.Slug, &course.Name) // check err
	if err != nil {
		return course, err
	}

	return course, nil
}

/*
GetModulesByCourseSlug gets all of the modules
associated with a course by the slug of the course.
*/
func GetModulesByCourseSlug(db *sql.DB, slug string) ([]Module, error) {
	rows, err := db.Query(`
		SELECT module.name, module.description, module.image_link, module.slug
		FROM module
		JOIN courses
		ON module.course_id = courses.id
		WHERE courses.slug = ?
	`, slug) // check err
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []Module
	for rows.Next() {
		var module Module
		err = rows.Scan(&module.Name, &module.Description, &module.ImageLink, &module.Slug) // check err
		if err != nil {
			return modules, err
		}
		modules = append(modules, module)
	}
	err = rows.Err() // check err
	if err != nil {
		return modules, err
	}
	return modules, nil
}

/*
Page returns the HTML page with
all of the modules for a course.
*/
func Page(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["course"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	course, err := GetCourse(db, slug)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	courses, err := GetCourses(db)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	modules, err := GetModulesByCourseSlug(db, slug)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	coursePage := PageData{
		Title:      course.Name,
		CourseSlug: slug,
		Course:     course,
		Courses:    courses,
		Modules:    modules,
		Paths:      []helpers.Path{{Name: "Courses", Link: "/"}},
	}

	helpers.RenderTemplate(r, w, "course", coursePage)
}
