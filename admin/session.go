package admin

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/course"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

/*
GetSessionBySlugs ...
*/
func GetSessionBySlugs(db *sql.DB, course, module, session string) (string, int, int, error) {
	var name string
	var sessionType int
	var sessionID int
	query :=
		`
	SELECT session.name, session.session_type, session.id
	FROM session
	INNER JOIN module
	ON module.id = session.module_id
	INNER JOIN courses
	ON courses.id = module.course_id
	WHERE courses.slug = ?
	AND module.slug = ?
	AND session.slug = ?
	`
	err := db.QueryRow(query, course, module, session).Scan(&name, &sessionType, &sessionID)
	if err != nil {
		return "", 0, 0, err
	}
	return name, sessionType, sessionID, nil
}

/*
SessionText ...
*/
type SessionText struct {
	Title string
	Paths []helpers.Path
	Name  string
	Text  string
}

/*
GetSessionText ...
*/
func GetSessionText(db *sql.DB, sessionID int) (string, error) {
	var text string
	query :=
		`
	SELECT session_text.text
	FROM session_text
	WHERE session_text.session_id = ?
	`
	err := db.QueryRow(query, sessionID).Scan(&text)
	if err != nil {
		return "", err
	}
	return text, nil
}

/*
SessionYoutube ...
*/
type SessionYoutube struct {
	Title   string
	Paths   []helpers.Path
	Name    string
	Text    string
	Youtube string
}

/*
SessionMultipleChoice ...
*/
type SessionMultipleChoice struct {
	Title     string
	Paths     []helpers.Path
	Name      string
	Questions []course.Question
}

// GetSessionsYoutube ...
func GetSessionsYoutube(db *sql.DB, sessionID int) (string, string, error) {
	var pageText string
	var youtube string
	// Query the database and scan the values into out variables. Don't forget to check for errors.
	query := `
	SELECT text, youtube_id
	FROM session_youtube
	WHERE session_id = ?
	`
	err := db.QueryRow(query, sessionID).Scan(&pageText, &youtube)
	if err != nil {
		return "", "", err
	}
	return pageText, youtube, nil
}

/*
HandleSessionText ...
*/
func HandleSessionText(w http.ResponseWriter, r *http.Request, db *sql.DB, name string, sessionID int) {
	sessionText, err := GetSessionText(db, sessionID)
	if err != nil {
		helpers.HandleError(err)
	}
	p := &SessionText{Title: "Text", Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}}, Name: name, Text: sessionText}
	helpers.RenderTemplate(r, w, "admin_session_text", p)
}

/*
HandleSessionYoutube ...
*/
func HandleSessionYoutube(w http.ResponseWriter, r *http.Request, db *sql.DB, name string, sessionID int) {
	text, youtube, err := GetSessionsYoutube(db, sessionID)
	if err != nil {
		helpers.HandleError(err)
	}
	p := &SessionYoutube{Title: "Text", Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}}, Name: name, Text: text, Youtube: youtube}
	helpers.RenderTemplate(r, w, "admin_session_youtube", p)
}

/*
HandleSessionMultipleChoice ...
*/
func HandleSessionMultipleChoice(w http.ResponseWriter, r *http.Request, db *sql.DB, name string, sessionID int) {
	questions, err := course.GetSessionsMultipleChoice(db, sessionID)
	if err != nil {
		helpers.HandleError(err)
	}
	p := &SessionMultipleChoice{Title: "Multiple Choice", Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}}, Name: name, Questions: questions}
	helpers.RenderTemplate(r, w, "admin_session_multiple_choice", p)
}

// Session ...
func Session(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	course := vars["course"]
	module := vars["module"]
	session := vars["session"]
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	name, sessionType, sessionID, err := GetSessionBySlugs(db, course, module, session)
	if err != nil {
		helpers.HandleError(err)
	}

	if sessionType == 1 {
		HandleSessionText(w, r, db, name, sessionID)
	} else if sessionType == 2 {
		HandleSessionYoutube(w, r, db, name, sessionID)
	} else if sessionType == 3 {
		HandleSessionMultipleChoice(w, r, db, name, sessionID)
	}
}

/*
UpdateSessionNameInDB ...
*/
func UpdateSessionNameInDB(db *sql.DB, name string, id int) error {
	insForm, err := db.Prepare("UPDATE session SET name=? WHERE id=?")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(name, id)
	if err != nil {
		return err
	}
	return nil
}

/*
UpdateSessionTextInDB ...
*/
func UpdateSessionTextInDB(db *sql.DB, text string, id int) error {
	insForm, err := db.Prepare("UPDATE session_text SET text=? WHERE session_id=?")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(text, id)
	if err != nil {
		return err
	}
	return nil
}

/*
UpdateSessionYoutubeInDB ...
*/
func UpdateSessionYoutubeInDB(db *sql.DB, text, youtube string, id int) error {
	insForm, err := db.Prepare("UPDATE session_youtube SET text=?, youtube_id=? WHERE session_id=?")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(text, youtube, id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateSession ...
func UpdateSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	course := vars["course"]
	module := vars["module"]
	session := vars["session"]
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	_, sessionType, sessionID, err := GetSessionBySlugs(db, course, module, session)
	if err != nil {
		helpers.HandleError(err)
	}

	newName := r.FormValue("name")

	err = UpdateSessionNameInDB(db, newName, sessionID)
	if err != nil {
		helpers.HandleError(err)
	}

	if sessionType == 1 {
		newText := r.FormValue("text")
		err = UpdateSessionTextInDB(db, newText, sessionID)
		if err != nil {
			helpers.HandleError(err)
		}
	} else if sessionType == 2 {
		newText := r.FormValue("text")
		newYoutube := r.FormValue("youtube")
		err = UpdateSessionYoutubeInDB(db, newText, newYoutube, sessionID)
		if err != nil {
			helpers.HandleError(err)
		}
	}
	http.Redirect(w, r, "/admin/courses/"+course+"/"+module+"/"+session, http.StatusSeeOther)
}

// NewSession ...
func NewSession(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "New session", Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}}}
	helpers.RenderTemplate(r, w, "admin_new_session", p)
}

// InsertNewSessionInDB ...
func InsertNewSessionInDB(db *sql.DB, name string, moduleSlug string, sessionType int) (string, int64, error) {
	slug := helpers.GenerateSlug(name)
	var moduleID int
	query :=
		`
		SELECT module.id
		FROM module
		WHERE module.slug = ?
		`
	err := db.QueryRow(query, moduleSlug).Scan(&moduleID)
	if err != nil {
		return "", 0, err
	}
	insForm, err := db.Prepare("INSERT INTO session (name, slug, module_id, session_type) VALUES (?, ?, ?, ?)")
	if err != nil {
		return "", 0, err
	}
	res, err := insForm.Exec(name, slug, moduleID, sessionType)
	if err != nil {
		return "", 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", 0, err
	}
	return slug, id, nil
}

// InsertNewSessionTextInDB ...
func InsertNewSessionTextInDB(db *sql.DB, sessionID int64, text string) error {
	insForm, err := db.Prepare("INSERT INTO session_text (session_id, text) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(sessionID, text)
	if err != nil {
		return err
	}
	return nil
}

// InsertNewSessionYoutubeInDB ...
func InsertNewSessionYoutubeInDB(db *sql.DB, sessionID int64, text, youtube string) error {
	insForm, err := db.Prepare("INSERT INTO session_youtube (session_id, text, youtube_id) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(sessionID, text, youtube)
	if err != nil {
		return err
	}
	return nil
}

// InsertNewSession ...
func InsertNewSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	course := vars["course"]
	module := vars["module"]
	name := r.FormValue("name")
	sessionTypeString := r.FormValue("session_type")
	sessionType, err := strconv.Atoi(sessionTypeString)
	if err != nil {
		helpers.HandleError(err)
	}
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	sessionSlug, sessionID, err := InsertNewSessionInDB(db, name, module, sessionType)
	if err != nil {
		helpers.HandleError(err)
	}
	if sessionType == 1 {
		text := r.FormValue("text")
		err = InsertNewSessionTextInDB(db, sessionID, text)
		if err != nil {
			helpers.HandleError(err)
		}
	} else if sessionType == 2 {
		text := r.FormValue("text")
		youtube := r.FormValue("youtube")
		err = InsertNewSessionYoutubeInDB(db, sessionID, text, youtube)
		if err != nil {
			helpers.HandleError(err)
		}
	}
	http.Redirect(w, r, "/admin/courses/"+course+"/"+module+"/"+sessionSlug, http.StatusSeeOther)
}
