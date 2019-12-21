package blog

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"

	"github.com/gorilla/mux"
)

// PostData contains the data for a blog post
type PostData struct {
	Title string
	Text  string
	Slug  string
}

// Page contains all of the data to render the blog index
type Page struct {
	Title string
	Paths []helpers.Path
	Posts []PostData
}

/*
GetBlogPosts gets all of the
blog posts from the database.
*/
func GetBlogPosts(db *sql.DB) ([]PostData, error) {
	rows, err := db.Query(`SELECT title, text, slug FROM blog_posts ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostData
	for rows.Next() {
		var post PostData
		err = rows.Scan(&post.Title, &post.Text, &post.Slug) // check err
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	err = rows.Err() // check err
	if err != nil {
		return posts, err
	}
	return posts, nil
}

/*
Blog returns the blog index
page with all of the blog posts
*/
func Blog(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()
	setting, err := helpers.GetSettingFromName(db, "BlogActivated")
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	// We check if blog is enabled
	if setting == "1" {
		blogPosts, err := GetBlogPosts(db)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}

		blogPage := Page{"Blog", []helpers.Path{{Name: "Blog", Link: "/blog"}}, blogPosts}

		helpers.RenderTemplate(r, w, "blog", blogPage)
	} else {
		http.NotFound(w, r)
	}
}

// PostPage contains all of the data for a single blog post
type PostPage struct {
	Title string
	Paths []helpers.Path
	Post  PostData
	Posts []PostData
}

/*
GetBlogPostBySlug returns the blog post
associated with the slug variable sent in.
*/
func GetBlogPostBySlug(db *sql.DB, slug string) (PostData, error) {
	var post PostData
	post.Slug = slug

	query := `SELECT title, text FROM blog_posts WHERE slug = ?`
	err := db.QueryRow(query, slug).Scan(&post.Title, &post.Text)
	if err != nil {
		return post, err
	}
	return post, nil
}

/*
Post renders a HTML page
with a single blog post.
*/
func Post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["post"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()
	setting, err := helpers.GetSettingFromName(db, "BlogActivated")
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	// We check if blog is enabled
	if setting == "1" {
		blogPosts, err := GetBlogPosts(db)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}

		post, err := GetBlogPostBySlug(db, slug)
		if err != nil {
			helpers.HandleError(err)
			helpers.InternalServerError(w)
			return
		}

		blogPage := PostPage{
			Title: post.Title,
			Paths: []helpers.Path{{Name: "Blog", Link: "/blog"}},
			Posts: blogPosts,
			Post:  post,
		}

		helpers.RenderTemplate(r, w, "blog_post", blogPage)
	} else {
		http.NotFound(w, r)
	}
}
