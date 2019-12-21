package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/initiate"
)

/*
Middleware declares middlware as alias of HandlerFunc
*/
type Middleware func(http.HandlerFunc) http.HandlerFunc

/*
WriteAccessLog creates the access
log string and writes it to file.
*/
func WriteAccessLog(r *http.Request) {
	start := time.Now()
	logString := time.Now().Format(time.RFC1123) + "*" + r.Method + "*" + r.URL.Path + "*" + r.RemoteAddr + "*" + r.UserAgent() + "*" + strconv.Itoa(int(time.Since(start).Nanoseconds())) + "\n"
	file, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		helpers.HandleError(err)
	}
	_, err = file.WriteString(logString)
	if err != nil {
		helpers.HandleError(err)
	}
	file.Close()
}

/*
Logging middleware logs the request.
*/
func Logging() Middleware {

	return func(f http.HandlerFunc) http.HandlerFunc {

		return func(w http.ResponseWriter, r *http.Request) {
			go WriteAccessLog(r)

			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

/*
NeuterAndLog is used for static files and
returns a not found if the request ends in a '/'
(meaning they want to look at the directory and not view a file).
Also logs the request.
*/
func NeuterAndLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		go WriteAccessLog(r)

		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

/*
AdminPage checks if the user is an admin.
If the user is, it grants access otherwise
returns a not found.
*/
func AdminPage() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			if !helpers.UserIsAdmin(r) {
				http.NotFound(w, r)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

/*
NewConfigPage returns the new config html page
*/
func NewConfigPage(w http.ResponseWriter, r *http.Request) {
	filename := "./templates/initiate/new_config.html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		helpers.HandleError(err)
	}
	fmt.Fprint(w, string(body))
}

/*
CreateConfigFile creates the actual
config file based on the input from the user.
*/
func CreateConfigFile(w http.ResponseWriter, r *http.Request) {
	var config helpers.Configuration
	config.DatabaseUser = r.FormValue("database_user")
	config.DatabasePassword = r.FormValue("database_password")
	config.DatabaseHost = r.FormValue("database_host")
	config.DatabasePort = r.FormValue("database_port")
	config.DatabaseName = r.FormValue("database_name")

	// User doesn't have to specify database name.
	// If the user doesn't, we make it opencourseplatform.
	if len(config.DatabaseName) == 0 {
		jsonData, err := json.MarshalIndent(config, "", "\t")
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}
		err = ioutil.WriteFile(helpers.SettingsFileName, jsonData, 0644)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}
		db, err := helpers.CreateDBHandler()
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}
		defer db.Close()
		config.DatabaseName = "opencourseplatform"
		err = initiate.CreateDatabase(db, config.DatabaseName)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}
	}

	jsonData, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	err = ioutil.WriteFile(helpers.SettingsFileName, jsonData, 0644)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	err = initiate.CreateTables(db)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*
ConfigFileExists is a middleware
that checks if the config file exists.
If it doesn't, we display the create
config file page.
*/
func ConfigFileExists() Middleware {

	return func(f http.HandlerFunc) http.HandlerFunc {

		return func(w http.ResponseWriter, r *http.Request) {

			_, err := ioutil.ReadFile(helpers.SettingsFileName)
			if err != nil {
				log.Println(r.Method)
				if r.Method == "GET" {
					NewConfigPage(w, r)
				} else if r.Method == "POST" {
					CreateConfigFile(w, r)
				}
			} else {
				f(w, r)
			}
		}
	}
}
