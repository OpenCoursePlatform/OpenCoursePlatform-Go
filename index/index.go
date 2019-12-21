package index

import (
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/course"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
)

/*
PageData contains the data for
rendering the index page.
*/
type PageData struct {
	Title   string
	Paths   []helpers.Path
	Courses []course.Course
}

/*
Page endpoint returns the index page.
*/
func Page(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	courses, err := course.GetCourses(db)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	coursePage := PageData{"Home", []helpers.Path{{Name: "Home", Link: "/"}}, courses}

	helpers.RenderTemplate(r, w, "index", coursePage)
}
