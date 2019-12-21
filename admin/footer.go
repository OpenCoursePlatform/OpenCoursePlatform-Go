package admin

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

// FooterCategoryData ...
type FooterCategoryData struct {
	Title string
	Slug  string
}

// FooterPage ...
type FooterPage struct {
	Title         string
	Paths         []helpers.Path
	Categories    []FooterCategoryData
	NewItemButton bool
	NewItemText   string
	NewItemLink   string
}

// GetFooterCategories ...
func GetFooterCategories(db *sql.DB) ([]FooterCategoryData, error) {
	rows, err := db.Query(`
	SELECT name, slug
	FROM footer_categories
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []FooterCategoryData
	for rows.Next() {
		var category FooterCategoryData
		err = rows.Scan(&category.Title, &category.Slug)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// FooterCategories ...
func FooterCategories(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	categories, err := GetFooterCategories(db)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &FooterPage{"Footer", []helpers.Path{{Name: "Admin", Link: "/admin"}}, categories, true, "New Footer Category", "/footer/new"}
	helpers.RenderTemplate(r, w, "admin_footer", p)
}

/*
GetFooterCategory ...
*/
func GetFooterCategory(db *sql.DB, slug string) (FooterCategoryData, int, error) {
	var category FooterCategoryData
	var id int
	category.Slug = slug
	query := `SELECT id, name FROM footer_categories WHERE slug = ?`
	err := db.QueryRow(query, slug).Scan(&id, &category.Title)
	if err != nil {
		return category, id, err
	}

	return category, id, nil
}

// GetFooterItems ...
func GetFooterItems(db *sql.DB, footerID int) ([]ToolbarData, error) {
	rows, err := db.Query(`
	SELECT pages.title, footer.footer_category_id IS NOT NULL
	FROM pages
	LEFT OUTER JOIN footer
	ON footer.page_id = pages.id
	AND footer.footer_category_id = ?
	`, footerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var toolbar []ToolbarData
	for rows.Next() {
		var item ToolbarData
		err = rows.Scan(&item.Title, &item.Active)
		if err != nil {
			return toolbar, err
		}
		toolbar = append(toolbar, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return toolbar, nil
}

// FooterCategoryPage ...
type FooterCategoryPage struct {
	Title         string
	Paths         []helpers.Path
	Category      FooterCategoryData
	ToolbarItems  []ToolbarData
	NewItemButton bool
	NewItemText   string
	NewItemLink   string
}

/*
FooterCategory ...
*/
func FooterCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["group"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	category, categoryID, err := GetFooterCategory(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	footerItems, err := GetFooterItems(db, categoryID)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &FooterCategoryPage{category.Title, []helpers.Path{{Name: "Admin", Link: "/admin"}, {Name: "Footer", Link: "/admin/footer"}}, category, footerItems, false, "", ""}
	helpers.RenderTemplate(r, w, "admin_footer_category", p)
}

// DropCategoryItems ...
func DropCategoryItems(db *sql.DB, categoryID int) error {
	insForm, err := db.Prepare("DELETE FROM footer WHERE footer_category_id = ?")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(categoryID)
	if err != nil {
		return err
	}
	return nil
}

/*
UpdateCategoryItems ...
*/
func UpdateCategoryItems(db *sql.DB, categoryID int, value string) error {
	var id int

	query := `SELECT id FROM pages WHERE title = ?`
	err := db.QueryRow(query, value).Scan(&id)
	if err != nil {
		return err
	}
	insForm, err := db.Prepare("INSERT INTO footer (footer_category_id, page_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE page_id=page_id")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(categoryID, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateFooter ...
func UpdateFooter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["group"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	_, categoryID, err := GetFooterCategory(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}
	err = r.ParseForm()
	if err != nil {
		helpers.HandleError(err)
	}
	err = DropCategoryItems(db, categoryID)
	if err != nil {
		helpers.HandleError(err)
	}
	for value := range r.Form {
		err = UpdateCategoryItems(db, categoryID, value)
		if err != nil {
			helpers.HandleError(err)
		}
	}
	http.Redirect(w, r, "/admin/footer/"+slug, http.StatusSeeOther)
}
