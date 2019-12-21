package admin

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/blog"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

// BlogPostData ...
type BlogPostData struct {
	Title string
	Text  string
	Slug  string
}

// BlogPageData ...
type BlogPageData struct {
	Title         string
	Paths         []helpers.Path
	Posts         []blog.PostData
	NewItemButton bool
	NewItemText   string
	NewItemLink   string
}

// Blog ...
func Blog(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	posts, err := blog.GetBlogPosts(db)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &BlogPageData{"Blog", []helpers.Path{{Name: "Admin", Link: "/admin"}}, posts, true, "New Blog Post", "/blog/new"}
	helpers.RenderTemplate(r, w, "admin_blog", p)
}

// BlogPostPageData ...
type BlogPostPageData struct {
	Title string
	Paths []helpers.Path
	Post  blog.PostData
}

// BlogPost ...
func BlogPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["post"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	post, err := blog.GetBlogPostBySlug(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &BlogPostPageData{post.Title, []helpers.Path{{Name: "Admin", Link: "/admin"}}, post}
	helpers.RenderTemplate(r, w, "admin_blog_post", p)
}

/*
UpdateBlogPostInDB ...
*/
func UpdateBlogPostInDB(db *sql.DB, title, text, slug string) error {
	insForm, err := db.Prepare("UPDATE blog_posts SET title=?, text=? WHERE slug=?")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(title, text, slug)
	if err != nil {
		return err
	}
	return nil
}

// UpdateBlogPost ...
func UpdateBlogPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["post"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	err = UpdateBlogPostInDB(db, r.FormValue("title"), r.FormValue("text"), slug)
	if err != nil {
		helpers.HandleError(err)
	}

	http.Redirect(w, r, "/admin/blog/"+slug, http.StatusSeeOther)
}

// NewBlogPost ...
func NewBlogPost(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "New blog post", Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}}}
	helpers.RenderTemplate(r, w, "admin_new_blog_post", p)
}

/*
InsertNewBlogPostInDB ...
*/
func InsertNewBlogPostInDB(db *sql.DB, title, text string) (string, error) {
	slug := helpers.GenerateSlug(title)
	insForm, err := db.Prepare("INSERT INTO blog_posts (title, text, slug) VALUES (?, ?, ?)")
	if err != nil {
		return "", err
	}
	_, err = insForm.Exec(title, text, slug)
	if err != nil {
		return "", err
	}

	return slug, nil
}

// InsertNewBlogPost ...
func InsertNewBlogPost(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	title := r.FormValue("title")
	text := r.FormValue("text")
	slug, err := InsertNewBlogPostInDB(db, title, text)
	if err != nil {
		helpers.HandleError(err)
	}
	http.Redirect(w, r, "/admin/blog/"+slug, http.StatusSeeOther)
}
