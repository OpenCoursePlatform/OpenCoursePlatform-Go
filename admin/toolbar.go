package admin

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
)

// ToolbarData ...
type ToolbarData struct {
	Title  string
	Active bool
}

// ToolbarPage ...
type ToolbarPage struct {
	Title         string
	Paths         []helpers.Path
	ToolbarItems  []ToolbarData
	NewItemButton bool
	NewItemText   string
	NewItemLink   string
}

// GetToolbar ...
func GetToolbar(db *sql.DB) ([]ToolbarData, error) {
	rows, err := db.Query(`
	SELECT pages.title, toolbar.id IS NOT NULL
	FROM pages
	LEFT OUTER JOIN toolbar
	ON toolbar.page_id = pages.id
	`)
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

// Toolbar ...
func Toolbar(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	toolbar, err := GetToolbar(db)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &ToolbarPage{"Toolbar", []helpers.Path{{Name: "Admin", Link: "/admin"}}, toolbar, false, "", ""}
	helpers.RenderTemplate(r, w, "admin_toolbar", p)
}

/*
UpdateValue ...
*/
func UpdateValue(db *sql.DB, value string) error {
	var id int

	query := `SELECT id FROM pages WHERE title = ?`
	err := db.QueryRow(query, value).Scan(&id)
	if err != nil {
		return err
	}
	insForm, err := db.Prepare("INSERT INTO toolbar (page_id) VALUES (?) ON DUPLICATE KEY UPDATE page_id=page_id")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

// DropToolbar ...
func DropToolbar(db *sql.DB) error {
	_, err := db.Query("TRUNCATE TABLE toolbar")
	return err
}

// UpdateToolbar ...
func UpdateToolbar(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	err = r.ParseForm()
	if err != nil {
		helpers.HandleError(err)
	}
	err = DropToolbar(db)
	if err != nil {
		helpers.HandleError(err)
	}
	for value := range r.Form {
		err = UpdateValue(db, value)
		if err != nil {
			helpers.HandleError(err)
		}
	}
	http.Redirect(w, r, "/admin/toolbar", http.StatusSeeOther)
}
