package podcast

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"

	"github.com/gorilla/mux"
)

/*
PostData contains data for
an individual podcast post.
To be compared to a blog post,
but with an additional field for
the file.
*/
type PostData struct {
	Title string
	Text  string
	File  string
	Slug  string
}

/*
Page struct is used for the
podcast page data being
rendered and should not be
used in other structs.
*/
type Page struct {
	Title string
	Paths []helpers.Path
	Posts []PostData
}

/*
GetPodcastPosts returns all of the
podcasts post from the database.
*/
func GetPodcastPosts(db *sql.DB) ([]PostData, error) {
	rows, err := db.Query(`SELECT title, text, file, slug FROM podcast_posts ORDER BY id DESC`) // check err
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostData
	for rows.Next() {
		var post PostData
		err = rows.Scan(&post.Title, &post.Text, &post.File, &post.Slug) // check err
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	err = rows.Err() // check err
	if err != nil {
		return []PostData{}, err
	}
	return posts, nil

}

/*
Podcast endpoints displays all podcast posts.
*/
func Podcast(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	setting, err := helpers.GetSettingFromName(db, "PodcastActivated")
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	// Checks if podcasts are enabled
	if setting == "1" {
		podcastPosts, err := GetPodcastPosts(db)
		if err != nil {
			helpers.HandleError(err)
			http.NotFound(w, r)
			return
		}

		helpers.RenderTemplate(r, w, "podcast", Page{"Podcast", []helpers.Path{{Name: "Podcast", Link: "/podcast"}}, podcastPosts})
	} else {
		http.NotFound(w, r)
		return
	}
}

/*
PostPage contains the data to render
a single podcast post
*/
type PostPage struct {
	Title string
	Paths []helpers.Path
	Post  PostData
	// Is used for sidebar
	Posts []PostData
}

/*
GetPodcastPost retrieves a single podcast
post from the database by its slug and.
*/
func GetPodcastPost(db *sql.DB, slug string) (PostData, error) {
	var post PostData
	post.Slug = slug
	query := `SELECT title, text, file FROM podcast_posts WHERE slug = ?`
	err := db.QueryRow(query, slug).Scan(&post.Title, &post.Text, &post.File)
	if err != nil {
		return post, err
	}
	return post, nil
}

/*
Post endpoints displays a single podcast post.
*/
func Post(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()
	setting, err := helpers.GetSettingFromName(db, "PodcastActivated")
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	// Checks if podcasts are enabled
	if setting == "1" {
		podcastPosts, err := GetPodcastPosts(db)
		if err != nil {
			helpers.HandleError(err)
			http.NotFound(w, r)
			return
		}

		slug := mux.Vars(r)["post"]

		post, err := GetPodcastPost(db, slug)
		if err != nil {
			helpers.HandleError(err)
			http.NotFound(w, r)
			return
		}

		podcastPage := PostPage{post.Title, []helpers.Path{{Name: "Podcast", Link: "/podcast"}, {Name: "Post", Link: "/podcast"}}, post, podcastPosts}

		helpers.RenderTemplate(r, w, "podcast_post", podcastPage)
	} else {
		http.NotFound(w, r)
		return
	}
}
