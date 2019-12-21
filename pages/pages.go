package pages

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

/*
PageData contains data for the single page.
*/
type PageData struct {
	Title string
	Text  string
	Slug  string
}

/*
SinglePageData contains data for a
single page and general data
*/
type SinglePageData struct {
	Title string
	Paths []helpers.Path
	Page  PageData
}

/*
GetPage gets a single page from the database by
its slug.
*/
func GetPage(db *sql.DB, slug string) (PageData, error) {
	var page PageData
	page.Slug = slug

	query := `SELECT title, content FROM pages WHERE slug = ?`
	err := db.QueryRow(query, "/"+slug).Scan(&page.Title, &page.Text)
	if err != nil {
		return page, err
	}
	return page, nil
}

/*
SinglePage returns a single page to the user.
The type of page is 'semi-static' with a possibility
to edit the data in the database but should be used for
pages such as about etc.
*/
func SinglePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["page"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	page, err := GetPage(db, slug)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}

	p := &SinglePageData{page.Title, []helpers.Path{{Name: "Home", Link: "/"}}, page}
	helpers.RenderTemplate(r, w, "page", p)
}
