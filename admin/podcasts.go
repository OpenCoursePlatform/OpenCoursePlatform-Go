package admin

import (
	"database/sql"
	"io"
	"net/http"
	"os"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/podcast"
	"github.com/gorilla/mux"
)

// PodcastPostData ...
type PodcastPostData struct {
	Title string
	Text  string
	Slug  string
}

// PodcastPageData ...
type PodcastPageData struct {
	Title         string
	Paths         []helpers.Path
	Posts         []podcast.PostData
	NewItemButton bool
	NewItemText   string
	NewItemLink   string
}

// Podcast ...
func Podcast(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	posts, err := podcast.GetPodcastPosts(db)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &PodcastPageData{"Podcast", []helpers.Path{{Name: "Admin", Link: "/admin"}}, posts, true, "New Podcast Post", "/podcast/new"}
	helpers.RenderTemplate(r, w, "admin_podcast", p)
}

// PodcastPostPageData ...
type PodcastPostPageData struct {
	Title string
	Paths []helpers.Path
	Post  podcast.PostData
}

// PodcastPost ...
func PodcastPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["post"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	post, err := podcast.GetPodcastPost(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &PodcastPostPageData{"Post", []helpers.Path{{Name: "Admin", Link: "/admin"}}, post}
	helpers.RenderTemplate(r, w, "admin_podcast_post", p)
}

/*
UpdatePodcastPostInDB ...
*/
func UpdatePodcastPostInDB(db *sql.DB, title, text, file, slug string) error {
	insForm, err := db.Prepare("UPDATE podcast_posts SET title=?, text=?, file=? WHERE slug=?")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(title, text, file, slug)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePodcastPost ...
func UpdatePodcastPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["post"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		helpers.HandleError(err)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		helpers.HandleError(err)
		return
	}
	defer file.Close()
	f, err := os.OpenFile("./static/podcasts/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		helpers.HandleError(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	err = UpdatePodcastPostInDB(db, r.FormValue("title"), r.FormValue("text"), handler.Filename, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	http.Redirect(w, r, "/admin/podcast/"+slug, http.StatusSeeOther)
}

// NewPodcastPost ...
func NewPodcastPost(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "New podcast post", Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}}}
	helpers.RenderTemplate(r, w, "admin_new_podcast_post", p)
}

/*
InsertNewPodcastPostInDB ...
*/
func InsertNewPodcastPostInDB(db *sql.DB, title, text, file string) (string, error) {
	slug := helpers.GenerateSlug(title)
	insForm, err := db.Prepare("INSERT INTO podcast_posts (title, text, slug, file) VALUES (?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	_, err = insForm.Exec(title, text, slug, file)
	if err != nil {
		return "", err
	}

	return slug, nil
}

// InsertNewPodcastPost ...
func InsertNewPodcastPost(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	title := r.FormValue("title")
	text := r.FormValue("text")
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		helpers.HandleError(err)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		helpers.HandleError(err)
		return
	}
	defer file.Close()
	f, err := os.OpenFile("./static/podcasts/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		helpers.HandleError(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	slug, err := InsertNewPodcastPostInDB(db, title, text, handler.Filename)
	if err != nil {
		helpers.HandleError(err)
	}
	http.Redirect(w, r, "/admin/podcast/"+slug, http.StatusSeeOther)
}
