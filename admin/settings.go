package admin

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

// SettingData ...
type SettingData struct {
	Name  string
	Value string
}

// SettingsData ...
type SettingsData struct {
	Title         string
	Paths         []helpers.Path
	Options       []SettingData
	NewItemButton bool
	NewItemText   string
	NewItemLink   string
}

/*
GetSettings ...
*/
func GetSettings(db *sql.DB) ([]SettingData, error) {
	rows, err := db.Query(`SELECT option_name, option_value FROM settings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []SettingData
	for rows.Next() {
		var setting SettingData
		err = rows.Scan(&setting.Name, &setting.Value)
		if err != nil {
			return settings, err
		}
		settings = append(settings, setting)
	}
	err = rows.Err()
	if err != nil {
		return settings, err
	}

	return settings, nil
}

// Settings ...
func Settings(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	settings, err := GetSettings(db)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &SettingsData{"Settings", []helpers.Path{{Name: "Admin", Link: "/admin"}}, settings, true, "New Setting", "/settings/new"}
	helpers.RenderTemplate(r, w, "admin_settings", p)
}

// SettingPageData ...
type SettingPageData struct {
	Title  string
	Paths  []helpers.Path
	Option SettingData
}

/*
GetSettingFromName ...
*/
func GetSettingFromName(db *sql.DB, name string) (SettingData, error) {
	var setting SettingData

	query := `SELECT option_name, option_value FROM settings WHERE option_name = ?`
	err := db.QueryRow(query, name).Scan(&setting.Name, &setting.Value)
	if err != nil {
		return setting, err
	}
	return setting, nil
}

// Setting ...
func Setting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["setting"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	setting, err := GetSettingFromName(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &SettingPageData{setting.Name, []helpers.Path{{Name: "Admin", Link: "/admin"}}, setting}
	helpers.RenderTemplate(r, w, "admin_setting", p)
}

/*
UpdateSettingWithName ...
*/
func UpdateSettingWithName(db *sql.DB, name, value string) error {
	insForm, err := db.Prepare("UPDATE settings SET option_value=? WHERE option_name=?")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(value, name)
	if err != nil {
		return err
	}
	return nil
}

// UpdateSetting ...
func UpdateSetting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["setting"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	err = UpdateSettingWithName(db, slug, r.FormValue("value"))
	if err != nil {
		helpers.HandleError(err)
	}

	http.Redirect(w, r, "/admin/settings/"+slug, http.StatusSeeOther)
}

// NewSetting ...
func NewSetting(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "New setting", Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}}}
	helpers.RenderTemplate(r, w, "admin_new_setting", p)
}

/*
InsertNewSettingInDB ...
*/
func InsertNewSettingInDB(db *sql.DB, key, value string) error {
	insForm, err := db.Prepare("INSERT INTO settings (option_name, option_value) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(key, value)
	if err != nil {
		return err
	}

	return nil
}

// InsertNewSetting ...
func InsertNewSetting(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	key := r.FormValue("key")
	value := r.FormValue("value")
	err = InsertNewSettingInDB(db, key, value)
	if err != nil {
		helpers.HandleError(err)
	}
	http.Redirect(w, r, "/admin/settings/"+key, http.StatusSeeOther)
}
