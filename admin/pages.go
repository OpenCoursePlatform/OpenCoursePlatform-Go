package admin

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

// PageData ...
type PageData struct {
	Title string
	Text  string
	Slug  string
}

// PagesData ...
type PagesData struct {
	Title         string
	Paths         []helpers.Path
	Pages         []PageData
	NewItemButton bool
	NewItemText   string
	NewItemLink   string
}

/*
GetPages ...
*/
func GetPages(db *sql.DB) ([]PageData, error) {
	rows, err := db.Query(`SELECT title, content, slug FROM pages ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []PageData
	for rows.Next() {
		var page PageData
		err = rows.Scan(&page.Title, &page.Text, &page.Slug) // check err
		if err != nil {
			return pages, err
		}
		pages = append(pages, page)
	}
	err = rows.Err() // check err
	if err != nil {
		return pages, err
	}
	return pages, nil
}

// Pages ...
func Pages(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	pages, err := GetPages(db)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &PagesData{"Pages", []helpers.Path{{Name: "Admin", Link: "/admin"}}, pages, true, "New Page", "/pages/new"}
	helpers.RenderTemplate(r, w, "admin_pages", p)
}

// SinglePageData ...
type SinglePageData struct {
	Title string
	Paths []helpers.Path
	Page  PageData
}

/*
GetPage ...
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

// SinglePage ...
func SinglePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["page"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	page, err := GetPage(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &SinglePageData{page.Title, []helpers.Path{{Name: "Admin", Link: "/admin"}}, page}
	helpers.RenderTemplate(r, w, "admin_page", p)
}

/*
UpdatePageInDB ...
*/
func UpdatePageInDB(db *sql.DB, title, text, slug string) error {
	insForm, err := db.Prepare("UPDATE pages SET title=?, content=? WHERE slug=?")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(title, text, "/"+slug)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePage ...
func UpdatePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["page"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	err = UpdatePageInDB(db, r.FormValue("title"), r.FormValue("text"), slug)
	if err != nil {
		helpers.HandleError(err)
	}

	http.Redirect(w, r, "/admin/pages/"+slug, http.StatusSeeOther)
}

// NewPage ...
func NewPage(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "New Page", Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}}}
	helpers.RenderTemplate(r, w, "admin_new_page", p)
}

/*
InsertNewPageInDB ...
*/
func InsertNewPageInDB(db *sql.DB, title, text, slug string) error {
	insForm, err := db.Prepare("INSERT INTO pages (title, content, slug) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(title, text, "/"+slug)
	if err != nil {
		return err
	}

	return nil
}

// InsertNewPage ...
func InsertNewPage(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	title := r.FormValue("title")
	text := r.FormValue("text")
	slug := r.FormValue("slug")
	err = InsertNewPageInDB(db, title, text, slug)
	if err != nil {
		helpers.HandleError(err)
	}
	http.Redirect(w, r, "/admin/pages/"+slug, http.StatusSeeOther)
}
