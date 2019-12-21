package course

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

/*
Session struct contains session data both
for rendering and for logic such as deciding
the session type.
*/
type Session struct {
	ID          int
	Name        string
	Slug        string
	SessionType int
}

/*
TextData is used for
rendering a text session.
Should not be used by other structs.
*/
type TextData struct {
	Title    string
	Paths    []helpers.Path
	Course   Course
	Module   Module
	Session  Session
	Sessions []Session
	Text     string
}

/*
YoutubeData is used for
rendering a YouTube session.
Should not be used by other structs.
*/
type YoutubeData struct {
	Title    string
	Paths    []helpers.Path
	Course   Course
	Module   Module
	Session  Session
	Sessions []Session
	Text     string
	Youtube  string
}

/*
Option struct contains question data
for rendering.
*/
type Option struct {
	Text      string
	IsCorrect bool
	Checked   bool
}

/*
Question struct contains question data
for rendering.
*/
type Question struct {
	ID      int
	Text    string
	Options []Option
}

/*
MultipleChoiceData is used for
rendering a Multiple-choice session.
Should not be used by other structs.
*/
type MultipleChoiceData struct {
	Title     string
	Paths     []helpers.Path
	Course    Course
	Module    Module
	Session   Session
	Sessions  []Session
	Questions []Question
}

/*
GetSessionsText gets the session text for
the session id. Should only ever be one
possible value unless there's a bug or
someone does a manual insert.
*/
func GetSessionsText(db *sql.DB, sessionID int) (string, error) {
	var pageText string
	query := `
	SELECT session_text.text
	FROM session_text
	WHERE session_text.session_id = ?
	`
	err := db.QueryRow(query, sessionID).Scan(&pageText)
	if err != nil {
		return "", err
	}
	return pageText, nil
}

/*
GetSessionsByModuleID gets all of the
sessions associated with one module.
*/
func GetSessionsByModuleID(db *sql.DB, moduleID int) ([]Session, error) {
	rows, err := db.Query(`
	SELECT name, slug
	FROM session
	WHERE module_id = ?
	`, moduleID) // check err
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var session Session
		err = rows.Scan(&session.Name, &session.Slug) // check err
		if err != nil {
			return sessions, err
		}
		sessions = append(sessions, session)
	}
	err = rows.Err() // check err
	if err != nil {
		return sessions, err
	}
	return sessions, nil
}

/*
GetSession gets the session data by the
session slug.
*/
func GetSession(db *sql.DB, sessionSlug string) (Session, int, error) {
	var session Session
	var moduleID int
	query := `
	SELECT id, name, session_type, module_id
	FROM session
	WHERE slug = ?
	`
	err := db.QueryRow(query, sessionSlug).Scan(&session.ID, &session.Name, &session.SessionType, &moduleID)
	if err != nil {
		return session, 0, err
	}
	return session, moduleID, nil
}

/*
GetSessionsYoutube gets the data for
a YouTube session by the session id.
*/
func GetSessionsYoutube(db *sql.DB, sessionID int) (string, string, error) {
	var pageText string
	var youtube string
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
GetSessionsMultipleChoice gets the data for
a Multiple-choice session by the session id.
*/
func GetSessionsMultipleChoice(db *sql.DB, sessionID int) ([]Question, error) {
	rows, err := db.Query(`
	SELECT id, question
	FROM session_multiple_choice
	WHERE session_id = ?
	`, sessionID) // check err
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var question Question
		err = rows.Scan(&question.ID, &question.Text) // check err
		if err != nil {
			return questions, err
		}
		questions = append(questions, question)
	}
	err = rows.Err() // check err
	if err != nil {
		return questions, err
	}

	for index := range questions {
		rows, err := db.Query(`
		SELECT answer, is_correct
		FROM session_multiple_choice_answers
		WHERE multiple_choice_id = ?
		`, questions[index].ID) // check err
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var options []Option
		for rows.Next() {
			var option Option
			err = rows.Scan(&option.Text, &option.IsCorrect) // check err
			if err != nil {
				return questions, err
			}
			options = append(options, option)
		}
		err = rows.Err() // check err
		if err != nil {
			return questions, err
		}
		questions[index].Options = options
	}
	return questions, nil
}

/*
SessionPage endpoint returns the actual
session data. Could be any type
of session.
*/
func SessionPage(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	course, err := GetCourse(db, mux.Vars(r)["course"])
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	module, err := GetModule(db, mux.Vars(r)["module"])
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}

	session, moduleID, err := GetSession(db, mux.Vars(r)["session"])
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	sessions, err := GetSessionsByModuleID(db, moduleID)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	if session.SessionType == 1 {
		coursePage := TextData{
			Title:    session.Name,
			Paths:    []helpers.Path{{Name: "Courses", Link: "/courses"}, {Name: course.Name, Link: "/courses/" + course.Slug}, {Name: module.Name, Link: "/courses/" + course.Slug + "/" + module.Slug}},
			Course:   course,
			Module:   module,
			Session:  session,
			Sessions: sessions,
		}
		coursePageText, err := GetSessionsText(db, session.ID)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}

		coursePage.Text = coursePageText
		helpers.RenderTemplate(r, w, "text", coursePage)
	} else if session.SessionType == 2 {
		coursePage := YoutubeData{
			Title:    session.Name,
			Paths:    []helpers.Path{{Name: "Courses", Link: "/courses"}, {Name: course.Name, Link: "/courses/" + course.Slug}, {Name: module.Name, Link: "/courses/" + course.Slug + "/" + module.Slug}},
			Course:   course,
			Module:   module,
			Session:  session,
			Sessions: sessions,
		}
		text, youtube, err := GetSessionsYoutube(db, session.ID)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}

		coursePage.Text = text
		coursePage.Youtube = youtube
		helpers.RenderTemplate(r, w, "youtube", coursePage)
	} else if session.SessionType == 3 {
		questions, err := GetSessionsMultipleChoice(db, session.ID)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}

		coursePage := MultipleChoiceData{
			Title:     session.Name,
			Paths:     []helpers.Path{{Name: "Courses", Link: "/courses"}, {Name: course.Name, Link: "/courses/" + course.Slug}, {Name: module.Name, Link: "/courses/" + course.Slug + "/" + module.Slug}},
			Course:    course,
			Module:    module,
			Session:   session,
			Sessions:  sessions,
			Questions: questions,
		}
		helpers.RenderTemplate(r, w, "multiple_choice", coursePage)
	}
}

// SessionAnswer ...
func SessionAnswer(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	session, moduleID, err := GetSession(db, mux.Vars(r)["session"])
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	sessions, err := GetSessionsByModuleID(db, moduleID)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	if session.SessionType == 3 {
		r.ParseForm()
		questions, err := GetSessionsMultipleChoice(db, session.ID)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}

		for i := range questions {
			for j := range questions[i].Options {
				answer := questions[i].Options[j].Text
				for k := range r.Form[questions[i].Text] {
					if answer == r.Form[questions[i].Text][k] {
						questions[i].Options[j].Checked = true
					}
				}
			}
		}

		coursePage := MultipleChoiceData{
			Title:     session.Name,
			Paths:     []helpers.Path{{Name: "Introduction to Python", Link: "/introduction-to-python"}},
			Sessions:  sessions,
			Questions: questions,
		}
		helpers.RenderTemplate(r, w, "multiple_choice_answered", coursePage)
	} else {
		http.NotFound(w, r)
		return
	}
}
