package forum

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/authentication"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

/*
GetForumTopicsAndPost returns forums
topics with their initial forum post.
*/
func GetForumTopicsAndPost(db *sql.DB) ([]FPost, error) {
	rows, err := db.Query(`
	SELECT forum_topics.title, forum_topics.slug, forum_posts.text
	FROM forum_topics 
	INNER JOIN forum_posts ON forum_posts.id = (
		SELECT id
		FROM forum_posts as fp
		WHERE fp.topic_id = forum_topics.id
		LIMIT 1
	)
	LIMIT 10
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forumTopics []FPost
	for rows.Next() {
		var topic FPost
		err = rows.Scan(&topic.Title, &topic.Slug, &topic.Text)
		if err != nil {
			return forumTopics, err
		}
		forumTopics = append(forumTopics, topic)
	}
	err = rows.Err()
	if err != nil {
		return forumTopics, err
	}

	return forumTopics, nil
}

/*
GetForumPosts returns forum topics with
all forum posts by the forum topic slug
from the database.
*/
func GetForumPosts(db *sql.DB, slug string) ([]FPost, error) {
	rows, err := db.Query(`
	SELECT forum_topics.title, forum_topics.slug, forum_posts.text
	FROM forum_topics 
	INNER JOIN forum_posts ON forum_topics.id = forum_posts.topic_id
	WHERE forum_topics.slug = ?
	`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forumPosts []FPost
	for rows.Next() {
		var post FPost
		err = rows.Scan(&post.Title, &post.Slug, &post.Text)
		if err != nil {
			return forumPosts, err
		}
		forumPosts = append(forumPosts, post)
	}
	err = rows.Err()
	if err != nil {
		return forumPosts, err
	}
	return forumPosts, nil
}

/*
PostPage endpoint returns a forum
topic with all of its posts.
*/
func PostPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["topic"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	setting, err := helpers.GetSettingFromName(db, "ForumActivated")
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	// We check if the forum is activated
	if setting == "1" {
		forumTopics, err := GetForumTopicsAndPost(db)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}

		forumPosts, err := GetForumPosts(db, slug)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}

		forumPage := Post{"Forum post", []helpers.Path{{Name: "Forum", Link: "/forum"}, {Name: "Post", Link: "/forum"}}, forumTopics, forumPosts}

		helpers.RenderTemplate(r, w, "post", forumPage)
	} else {
		http.NotFound(w, r)
	}
}

/*
Post contains the data to render
the HTML page for a single post.
*/
type Post struct {
	Title  string
	Paths  []helpers.Path
	Topics []FPost
	Posts  []FPost
}

// FPost contains the data for a single Forum post.
type FPost struct {
	Title string
	Text  string
	Slug  string
}

/*
PostPageData contains the data to render
the HTML page for the forum index.
*/
type PostPageData struct {
	Title string
	Paths []helpers.Path
	Posts []FPost
}

// Forum endpoint returns the forum index.
func Forum(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()
	setting, err := helpers.GetSettingFromName(db, "ForumActivated")
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	// We check if the forum is activated
	if setting == "1" {

		forumPosts, err := GetForumTopicsAndPost(db)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}
		forumPage := PostPageData{"Forum", []helpers.Path{{Name: "Forum", Link: "/forum"}}, forumPosts}

		helpers.RenderTemplate(r, w, "forum_index", forumPage)
	} else {
		http.NotFound(w, r)
	}
}

/*
Page is a generic struct for a general HTML page.
Only contains the necessary data to render the page.
*/
type Page struct {
	Title string
	Paths []helpers.Path
}

// NewForumTopic returns the HTML page we use to create a forum topic.
func NewForumTopic(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()
	setting, err := helpers.GetSettingFromName(db, "ForumActivated")
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	// We check if the forum is activated
	if setting == "1" {
		supportPage := Page{"New Forum Topic", []helpers.Path{{Name: "Forum", Link: "/forum"}, {Name: "New", Link: "/forum/new"}}}

		helpers.RenderTemplate(r, w, "forum_new_topic", supportPage)
	} else {
		http.NotFound(w, r)
	}
}

/*
InsertNewTopicInDB inserts a new forum topic
and the associated forum post into the database.
*/
func InsertNewTopicInDB(db *sql.DB, username, title, text string) (string, error) {
	slug := helpers.GenerateSlug(title)
	userID, err := authentication.GetUserIDFromUsername(db, username)
	if err != nil {
		return "", err
	}
	insForm, err := db.Prepare("INSERT INTO forum_topics (title, slug) VALUES (?, ?)")
	if err != nil {
		return "", err
	}
	res, err := insForm.Exec(title, slug)
	if err != nil {
		return "", err
	}
	insForm, err = db.Prepare("INSERT INTO forum_posts (topic_id, author_id, text) VALUES (?, ?, ?)")
	if err != nil {
		return "", err
	}
	lastInsertedID, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	_, err = insForm.Exec(lastInsertedID, userID, text)
	if err != nil {
		return "", err
	}
	return slug, nil
}

/*
InsertNewTopic endpoint parses the form data and
inserts a new forum topic with an associated post.
*/
func InsertNewTopic(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()
	setting, err := helpers.GetSettingFromName(db, "ForumActivated")
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	// We check if the forum is activated
	if setting == "1" {
		title := r.FormValue("title")
		text := r.FormValue("text")
		username, err := helpers.GetUsernameFromRequest(r)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}
		slug, err := InsertNewTopicInDB(db, username, title, text)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}
		http.Redirect(w, r, "/forum/"+slug, http.StatusSeeOther)
	} else {
		http.NotFound(w, r)
	}
}

/*
InsertNewAnswer inserts a new forum post as an answer to the forum
topic associated with the topic by the topic slug into the database.
*/
func InsertNewAnswer(db *sql.DB, username, text, topic string) error {
	var topicID int
	query :=
		`
		SELECT id
		FROM forum_topics
		WHERE slug = ?
		`
	err := db.QueryRow(query, topic).Scan(&topicID)
	if err != nil {
		return err
	}
	userID, err := authentication.GetUserIDFromUsername(db, username)
	if err != nil {
		return err
	}
	insForm, err := db.Prepare("INSERT INTO forum_posts (topic_id, author_id, text) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(topicID, userID, text)
	if err != nil {
		return err
	}
	return nil
}

// AnswerPost parses and inserts the data for a new forum topic answer.
func AnswerPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topic := vars["topic"]
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()
	setting, err := helpers.GetSettingFromName(db, "ForumActivated")
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	// We check if the forum is activated
	if setting == "1" {
		text := r.FormValue("text")
		username, err := helpers.GetUsernameFromRequest(r)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}
		err = InsertNewAnswer(db, username, text, topic)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}
		http.Redirect(w, r, "/forum/"+topic, http.StatusSeeOther)
	} else {
		http.NotFound(w, r)
	}
}
