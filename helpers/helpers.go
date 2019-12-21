package helpers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

/*
Path is used in toolbar and a list of path
should be the way to get to the current URL.
*/
type Path struct {
	Name string
	Link string
}

/*
SettingsFileName is the name of the settings file.
Could be customized to anything you want
but _should_ end in .json
*/
const SettingsFileName = "settings.json"

var (
	templates = template.Must(ParseTemplates(), nil)
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key = []byte("super-secret-key")
	// Store ...
	Store = sessions.NewCookieStore(key)
)

/*
Configuration is the data stored in the
settings.json file in a go struct.
*/
type Configuration struct {
	DatabaseUser     string
	DatabasePassword string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
}

// HandleError logs all errors and writes them to file.
func HandleError(err error) {
	go func() {
		start := time.Now()
		logString := time.Now().Format(time.RFC1123) + "*" + err.Error() + "*" + strconv.Itoa(int(time.Since(start).Nanoseconds())) + "\n"
		file, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		_, err = file.WriteString(logString)
		if err != nil {
			log.Println(err)
		}
		file.Close()
		log.Println(logString)
	}()
}

// ToolbarPage contains the data used in the toolbar items
type ToolbarPage struct {
	Title string
	Slug  string
}

// FooterPage contains the data used in the footer items
type FooterPage struct {
	Title    string
	Slug     string
	Category string
}

/*
Markdown converts the markdown to HTML.
To be used inside templates.
*/
func Markdown(arg string) template.HTML {
	unsafe := blackfriday.MarkdownCommon([]byte(arg))
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	output := p.SanitizeBytes(unsafe)

	return template.HTML(string(output))
}

/*
Function map of the functions to be used within
HTML templates.
*/
var funcMap = template.FuncMap{
	"Markdown": Markdown,
}

/*
CheckIfStringContainsHTML checks
if the string contains the .html suffix.
*/
func CheckIfStringContainsHTML(path string, templ *template.Template) error {
	var err error
	if strings.Contains(path, ".html") {
		_, err = templ.ParseFiles(path)
		if err != nil {
			HandleError(err)
		}
	}

	return err
}

/*
ParseTemplates parses the HTML templates in the
templates folder. Including any HTML templates
in subdirectories.
*/
func ParseTemplates() *template.Template {
	templ := template.New("").Funcs(funcMap)
	basePath := os.Getenv("BASEPATH")
	err := filepath.Walk(basePath+"templates", func(path string, info os.FileInfo, err error) error {
		err = CheckIfStringContainsHTML(path, templ)
		return err
	})

	if err != nil {
		HandleError(err)
	}

	return templ
}

/*
GetUsernameFromRequest gets the username from
the storage and returns it or an error.
*/
func GetUsernameFromRequest(r *http.Request) (string, error) {
	sessions, err := Store.Get(r, "session")
	if err != nil {
		return "", err
	}
	if sessions.Values["username"] == nil {
		return "", errors.New("Username not present in session")
	}
	return sessions.Values["username"].(string), nil
}

/*
CreateDBHandler returns a database handler
to the database exactly as defined in
the settings.json file.
*/
func CreateDBHandler() (*sql.DB, error) {
	var config Configuration
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(path + "/" + SettingsFileName)
	if err != nil {
		// For testing because it's going to look for the settings file
		// for example in the "helpers" directory when running helpers_test etc.
		data, err = ioutil.ReadFile(filepath.Dir(path) + "/" + SettingsFileName)
		if err != nil {
			return nil, err
		}
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	db, err := CreateDBHandlerWithDB(config.DatabaseName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

/*
CreateDBHandlerWithDB returns a database handler
to the database as defined in the settings.json
file without regards to the database_name config
option. Is very useful for unit tests.
*/
func CreateDBHandlerWithDB(database string) (*sql.DB, error) {
	var config Configuration
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(path + "/" + SettingsFileName)
	if err != nil {
		// For testing because it's going to look for the settings file
		// for example in the "helpers" directory when running helpers_test etc.
		data, err = ioutil.ReadFile(filepath.Dir(path) + "/" + SettingsFileName)
		if err != nil {
			return nil, err
		}
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	databaseURL :=
		config.DatabaseUser + ":" +
			config.DatabasePassword +
			"@(" + config.DatabaseHost +
			":" + config.DatabasePort +
			")/" + database +
			"?parseTime=true"

	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

/*
CheckIfUserIsPartOfGroup checks if the user sent
in is a part of the group sent in.
*/
func CheckIfUserIsPartOfGroup(db *sql.DB, username string, group string) (bool, error) {
	// String can't be empty since that is what the database
	// would return if no rows would be found.
	usernameCopy := "_____"

	query := `
	SELECT users.username
	FROM users
	INNER JOIN group_members
		ON group_members.user_id = users.id
	INNER JOIN user_groups
		ON user_groups.id = group_members.group_id
	WHERE users.username = ? AND user_groups.name = ?`
	err := db.QueryRow(query, username, group).Scan(&usernameCopy)
	if err != nil {
		return false, err
	}
	return username == usernameCopy, nil
}

/*
UserIsAdmin function checks if the user
is an administrator and returns
the result as a boolean.
*/
func UserIsAdmin(r *http.Request) bool {
	// If there is an error that is due to no username being present
	// so we can ignore the error.
	username, err := GetUsernameFromRequest(r)
	if err != nil {
		return false
	}

	db, err := CreateDBHandler()
	if err != nil {
		HandleError(err)
	}

	isAdmin, err := CheckIfUserIsPartOfGroup(db, username, "Administrator")
	if err != nil {
		HandleError(err)
	}

	return isAdmin
}

/*
ConvertInterfaceToHashmap converts an interface to a hashmap.
*/
func ConvertInterfaceToHashmap(inInterface interface{}) (map[string]interface{}, error) {
	hashmap := make(map[string]interface{})
	if inInterface != nil {
		inrec, err := json.Marshal(inInterface)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(inrec, &hashmap)
	}
	return hashmap, nil
}

/*
GetSettings returns all of the settings
from the database as a hashmap.
*/
func GetSettings(db *sql.DB) (map[string]string, error) {
	rows, err := db.Query(`
	SELECT option_name, option_value
	FROM settings
	`) // check err
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings map[string]string
	settings = make(map[string]string)
	for rows.Next() {
		var key string
		var value string
		err = rows.Scan(&key, &value) // check err
		if err != nil {
			return settings, err
		}
		settings[key] = value
	}
	err = rows.Err() // check err
	if err != nil {
		return settings, err
	}
	return settings, nil
}

/*
GetToolbarPages returns all of the toolbar pages
from the database.
*/
func GetToolbarPages(db *sql.DB) ([]ToolbarPage, error) {
	rows, err := db.Query(`
	SELECT pages.title, pages.slug
	FROM pages
	INNER JOIN toolbar
		ON toolbar.page_id = pages.id
	ORDER BY toolbar.id ASC
	`) // check err
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var toolbarPages []ToolbarPage
	for rows.Next() {
		var page ToolbarPage
		err = rows.Scan(&page.Title, &page.Slug) // check err
		if err != nil {
			return toolbarPages, err
		}
		toolbarPages = append(toolbarPages, page)
	}
	err = rows.Err() // check err
	if err != nil {
		return toolbarPages, err
	}
	return toolbarPages, nil
}

/*
GetFooterPages returns all of
the footer pages from the database
*/
func GetFooterPages(db *sql.DB) ([]FooterPage, error) {
	rows, err := db.Query(`
	SELECT pages.title, pages.slug, footer_categories.name
	FROM pages
	INNER JOIN footer
		ON footer.page_id = pages.id
	INNER JOIN footer_categories
		ON footer.footer_category_id = footer_categories.id
	`) // check err
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var footerPages []FooterPage
	for rows.Next() {
		var page FooterPage
		err = rows.Scan(&page.Title, &page.Slug, &page.Category) // check err
		if err != nil {
			return footerPages, err
		}
		footerPages = append(footerPages, page)
	}
	err = rows.Err() // check err
	if err != nil {
		return footerPages, err
	}
	return footerPages, nil
}

/*
GetSettingFromName returns the setting
value from the database by the setting name.
*/
func GetSettingFromName(db *sql.DB, name string) (string, error) {
	var value string

	query := `SELECT option_value FROM settings WHERE option_name = ?`
	err := db.QueryRow(query, name).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

/*
RenderTemplate renders the HTML file as sent
in by the 'tmpl' variable and uses the 'data' interface
as well as adds more data necessary for running the application.
*/
func RenderTemplate(r *http.Request, w http.ResponseWriter, tmpl string, data interface{}) {
	hashmap, err := ConvertInterfaceToHashmap(data)
	if err != nil {
		HandleError(err)
		InternalServerError(w)
		return
	}

	db, err := CreateDBHandler()
	if err != nil {
		HandleError(err)
		InternalServerError(w)
		return
	}

	hashmap["Settings"], err = GetSettings(db)
	if err != nil {
		HandleError(err)
		InternalServerError(w)
		return
	}

	hashmap["Toolbar"], err = GetToolbarPages(db)
	if err != nil {
		HandleError(err)
		InternalServerError(w)
		return
	}

	hashmap["Footer"], err = GetFooterPages(db)
	if err != nil {
		HandleError(err)
		InternalServerError(w)
		return
	}

	// If there is an error that is due to no username being present
	// so we can ignore the error.
	hashmap["Email"], _ = GetUsernameFromRequest(r)
	hashmap["UserIsAdmin"] = UserIsAdmin(r)

	defer db.Close()

	err = templates.ExecuteTemplate(w, tmpl+".html", hashmap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		HandleError(err)
		InternalServerError(w)
		return
	}
}

/*
GenerateRandomStringURLSafe returns a URL-safe, base64 encoded
securely generated random string.
It will return an error if the system's secure random
number generator fails to function correctly, in which
case the caller should not continue.
*/
func GenerateRandomStringURLSafe(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

/*
GenerateSlug generates a slug from an input name.
*/
func GenerateSlug(input string) string {
	re := regexp.MustCompile("[^a-z0-9]+")
	return strings.Trim(re.ReplaceAllString(strings.ToLower(input), "-"), "-")
}

// InternalServerError returns an internal server error.
func InternalServerError(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", 500)
}
