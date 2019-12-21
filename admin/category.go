package admin

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

// CategoryData ...
type CategoryData struct {
	Name string
	Slug string
}

// CategoriesData ...
type CategoriesData struct {
	Title         string
	Paths         []helpers.Path
	Categories    []CategoryData
	NewItemButton bool
	NewItemText   string
	NewItemLink   string
}

/*
GetCategories ...
*/
func GetCategories(db *sql.DB) ([]CategoryData, error) {
	rows, err := db.Query(`SELECT name, slug FROM course_categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryData
	for rows.Next() {
		var category CategoryData
		err = rows.Scan(&category.Name, &category.Slug)
		if err != nil {
			return categories, err
		}
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		return categories, err
	}
	return categories, nil
}

// Categories ...
func Categories(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	categories, err := GetCategories(db)
	if err != nil {
		helpers.HandleError(err)
	}

	data := &CategoriesData{"Categories", []helpers.Path{{Name: "Admin", Link: "/admin"}}, categories, true, "New category", "/categories/new"}
	helpers.RenderTemplate(r, w, "admin_categories", data)
}

// CategoryPageData ...
type CategoryPageData struct {
	Title    string
	Paths    []helpers.Path
	Category CategoryData
}

/*
GetCategory ...
*/
func GetCategory(db *sql.DB, slug string) (CategoryData, error) {
	var category CategoryData

	query := `SELECT name, slug FROM course_categories WHERE slug = ?`
	err := db.QueryRow(query, slug).Scan(&category.Name, &category.Slug)
	if err != nil {
		// No user
		return category, err
	}
	return category, nil
}

// Category ...
func Category(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["category"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	category, err := GetCategory(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &CategoryPageData{category.Name, []helpers.Path{{Name: "Admin", Link: "/admin"}}, category}
	helpers.RenderTemplate(r, w, "admin_category", p)
}

/*
UpdateCategoryInDatabase ...
*/
func UpdateCategoryInDatabase(db *sql.DB, newName string, slug string) error {
	insForm, err := db.Prepare("UPDATE course_categories SET name=? WHERE slug=?")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(newName, slug)
	if err != nil {
		return err
	}

	return nil
}

// UpdateCategory ...
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["category"]
	name := r.FormValue("category_name")

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	err = UpdateCategoryInDatabase(db, name, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	http.Redirect(w, r, "/admin/categories/"+slug, http.StatusSeeOther)

}

// DeleteCategory ...
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "title", Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}}}
	helpers.RenderTemplate(r, w, "admin", p)
}

/*
NewCategory ...
*/
func NewCategory(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "New category", Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}}}
	helpers.RenderTemplate(r, w, "admin_new_category", p)
}

/*
InsertNewCategoryInDB ...
*/
func InsertNewCategoryInDB(db *sql.DB, name string) (string, error) {
	slug := helpers.GenerateSlug(name)
	insForm, err := db.Prepare("INSERT INTO course_categories (name, slug) VALUES (?, ?)")
	if err != nil {
		return "", err
	}
	_, err = insForm.Exec(name, slug)
	if err != nil {
		return "", err
	}

	return slug, nil
}

/*
InsertNewCategory ...
*/
func InsertNewCategory(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	name := r.FormValue("category_name")
	slug, err := InsertNewCategoryInDB(db, name)
	if err != nil {
		helpers.HandleError(err)
	}
	http.Redirect(w, r, "/admin/categories/"+slug, http.StatusSeeOther)
}
