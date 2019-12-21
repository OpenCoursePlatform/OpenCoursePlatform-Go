package admin

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
)

// Page ...
type Page struct {
	Title string
	Paths []helpers.Path
}

// PageStats ...
type PageStats struct {
	Path  string
	Views int
	Speed int
}

// IndexPageData ...
type IndexPageData struct {
	Title     string
	Paths     []helpers.Path
	PageViews int
	Errors    int
	Tickets   int
	Revenue   int
	Pages     []PageStats
}

// GetPageStats ...
func GetPageStats(db *sql.DB) ([]PageStats, error) {
	rows, err := db.Query(`SELECT path, COUNT(path), FLOOR(AVG(speed)) FROM access_logs GROUP BY path ORDER BY COUNT(path) DESC LIMIT 10`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pageStats []PageStats
	for rows.Next() {
		var pageStat PageStats
		err = rows.Scan(&pageStat.Path, &pageStat.Views, &pageStat.Speed)
		if err != nil {
			return pageStats, err
		}
		pageStats = append(pageStats, pageStat)
	}
	err = rows.Err()
	if err != nil {
		return pageStats, err
	}
	return pageStats, nil
}

// Index ...
func Index(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	pages, err := GetPageStats(db)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &IndexPageData{"Admin", []helpers.Path{{Name: "Admin", Link: "/admin"}}, 1024, 2, 0, 10, pages}
	defer parseLog()
	helpers.RenderTemplate(r, w, "admin_index", p)
}

/*
	router.HandleFunc("/admin/courses/{course}/{module}/{session}", middleware.Chain(admin.Index, middleware.Logging())).Methods("GET")
	router.HandleFunc("/admin/courses/{course}/{module}/{session}", middleware.Chain(admin.Index, middleware.Logging())).Methods("POST")
*/
