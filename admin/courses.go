package admin

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/course"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

// CoursesPageData ...
type CoursesPageData struct {
	Title         string
	Paths         []helpers.Path
	Courses       []course.Course
	NewItemButton bool
	NewItemText   string
	NewItemLink   string
}

// Courses ...
func Courses(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	courses, err := course.GetCourses(db)
	if err != nil {
		helpers.HandleError(err)
	}

	coursePage := CoursesPageData{"Courses", []helpers.Path{{Name: "Admin", Link: "/admin"}}, courses, true, "New Course", "/courses/new"}
	helpers.RenderTemplate(r, w, "admin_courses", coursePage)
}

// ModuleData ...
type ModuleData struct {
	Name        string
	Description string
	Slug        string
	Course      string
}

// CoursePageData ...
type CoursePageData struct {
	Title         string
	Paths         []helpers.Path
	Course        course.Course
	Categories    []CategoryData
	Modules       []ModuleData
	NewItemButton bool
	NewItemText   string
	NewItemLink   string
}

/*
GetCourse ...
*/
func GetCourse(db *sql.DB, slug string) (course.Course, error) {
	var course course.Course

	query := `SELECT courses.name, courses.description, courses.slug, course_categories.name FROM courses JOIN course_categories ON courses.category_id = course_categories.id WHERE courses.slug = ?`
	err := db.QueryRow(query, slug).Scan(&course.Name, &course.Description, &course.Slug, &course.Category)
	if err != nil {
		return course, err
	}
	return course, nil
}

/*
GetModulesByCourseSlug ...
*/
func GetModulesByCourseSlug(db *sql.DB, slug string) ([]ModuleData, error) {
	rows, err := db.Query(`SELECT module.name, module.description, module.slug FROM module INNER JOIN courses ON courses.id = module.course_id WHERE courses.slug = ?`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []ModuleData
	for rows.Next() {
		var module ModuleData
		err = rows.Scan(&module.Name, &module.Description, &module.Slug)
		if err != nil {
			return modules, err
		}
		modules = append(modules, module)
	}
	err = rows.Err()
	if err != nil {
		return modules, err
	}
	return modules, nil
}

// Course ...
func Course(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["course"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	categories, err := GetCategories(db)
	if err != nil {
		helpers.HandleError(err)
	}

	course, err := GetCourse(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	modules, err := GetModulesByCourseSlug(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &CoursePageData{course.Name, []helpers.Path{{Name: "Admin", Link: "/admin"}}, course, categories, modules, true, "New Module", "/courses/" + course.Slug + "/new"}
	helpers.RenderTemplate(r, w, "admin_course", p)
}

/*
GetCategoryIDFromName ...
*/
func GetCategoryIDFromName(db *sql.DB, name string) (int, error) {
	var categoryID int

	query := `SELECT id FROM course_categories WHERE name = ?`
	err := db.QueryRow(query, name).Scan(&categoryID)
	if err != nil {
		return categoryID, err
	}

	return categoryID, nil
}

/*
UpdateCourseInDB ...
*/
func UpdateCourseInDB(db *sql.DB, name, description, slug string, categoryID int) error {
	insForm, err := db.Prepare("UPDATE courses SET name=?, description=?, category_id=? WHERE slug=?")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(name, description, categoryID, slug)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCourse ...
func UpdateCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["course"]
	name := r.FormValue("name")
	category := r.FormValue("category")

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	categoryID, err := GetCategoryIDFromName(db, category)
	if err != nil {
		helpers.HandleError(err)
	}

	err = UpdateCourseInDB(db, name, r.FormValue("description"), slug, categoryID)
	if err != nil {
		helpers.HandleError(err)
	}

	http.Redirect(w, r, "/admin/courses/"+slug, http.StatusSeeOther)
}

/*
NewCoursePage ...
*/
type NewCoursePage struct {
	Title      string
	Paths      []helpers.Path
	Categories []CategoryData
}

// NewCourse ...
func NewCourse(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()
	categories, err := GetCategories(db)
	if err != nil {
		helpers.HandleError(err)
	}
	p := &NewCoursePage{Title: "New course", Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}}, Categories: categories}
	helpers.RenderTemplate(r, w, "admin_new_course", p)
}

/*
InsertCourseInDB ...
*/
func InsertCourseInDB(db *sql.DB, name, description string, categoryID int) (string, error) {
	slug := helpers.GenerateSlug(name)
	insForm, err := db.Prepare("INSERT INTO courses (name, description, category_id, slug) VALUES (?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	_, err = insForm.Exec(name, description, categoryID, slug)
	if err != nil {
		return "", err
	}
	return slug, nil
}

// InsertNewCourse ...
func InsertNewCourse(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	category := r.FormValue("category")

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	categoryID, err := GetCategoryIDFromName(db, category)
	if err != nil {
		helpers.HandleError(err)
	}

	slug, err := InsertCourseInDB(db, name, r.FormValue("description"), categoryID)
	if err != nil {
		helpers.HandleError(err)
	}

	http.Redirect(w, r, "/admin/courses/"+slug, http.StatusSeeOther)
}

// DeleteCourseInDB ...
func DeleteCourseInDB(db *sql.DB, slug string) error {
	var courseID int
	err := db.QueryRow(`SELECT id FROM courses WHERE slug = ?`, slug).Scan(&courseID)
	if err != nil {
		return err
	}
	rows, err := db.Query(`SELECT id FROM module WHERE course_id = ?`, courseID) // check err
	if err != nil {
		return err
	}
	defer rows.Close()
	var moduleIDS []int
	for rows.Next() {
		var moduleID int
		err = rows.Scan(&moduleID) // check err
		if err != nil {
			return err
		}
		moduleIDS = append(moduleIDS, moduleID)
	}
	for index := range moduleIDS {
		form, err := db.Prepare("DELETE FROM session WHERE module_id = ?")
		if err != nil {
			return err
		}
		_, err = form.Exec(moduleIDS[index])
		if err != nil {
			return err
		}
		form, err = db.Prepare("DELETE FROM module WHERE id = ?")
		if err != nil {
			return err
		}
		_, err = form.Exec(moduleIDS[index])
		if err != nil {
			return err
		}
	}
	form, err := db.Prepare("DELETE FROM courses WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = form.Exec(courseID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCourse ...
func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["course"]
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	err = DeleteCourseInDB(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	http.Redirect(w, r, "/admin/courses", http.StatusSeeOther)
}
