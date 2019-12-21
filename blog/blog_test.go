package blog

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/initiate"

	_ "github.com/go-sql-driver/mysql"
)

func TestGetBlogPosts(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO blog_posts (id, title, text, slug, published, created)
	VALUES
		(1, 'Hej igen!', '# Hej igen!!\r\n\r\nOch här är en ny mening', 'hej-igen', '2019-10-23 13:22:54', '2019-10-23 13:22:54');	
	
	`)
	if err != nil {
		t.Errorf("Insertion of session in database failed. Error message: %s", err.Error())
		return
	}

	posts, err := GetBlogPosts(db)
	if err != nil {
		t.Errorf("GetBlogPosts failed. Error message: %s", err.Error())
		return
	}

	samplePosts := []PostData{{"Hej igen!", "# Hej igen!!\r\n\r\nOch här är en ny mening", "hej-igen"}}

	if posts[0] != samplePosts[0] {
		t.Errorf("GetBlogPosts failed. Structs are not equal")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestGetBlogPost(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO blog_posts (id, title, text, slug, published, created)
	VALUES
		(1, 'Hej igen!', '# Hej igen!!\r\n\r\nOch här är en ny mening', 'hej-igen', '2019-10-23 13:22:54', '2019-10-23 13:22:54');	
	
	`)
	if err != nil {
		t.Errorf("Insertion of session in database failed. Error message: %s", err.Error())
		return
	}

	post, err := GetBlogPostBySlug(db, "hej-igen")
	if err != nil {
		t.Errorf("GetBlogPostBySlug failed. Error message: %s", err.Error())
		return
	}

	samplePost := PostData{"Hej igen!", "# Hej igen!!\r\n\r\nOch här är en ny mening", "hej-igen"}

	if post != samplePost {
		t.Errorf("GetBlogPostBySlug failed. Structs are not equal")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestBlog(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Blog)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestPost(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Post)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestDBForBlog(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Blog)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := ""
	if rr.Body.String() == expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	err = initiate.WriteSettingsFile(config)
	if err != nil {
		t.Fatal(err)
	}

}

func TestDBForPost(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Post)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := ""
	if rr.Body.String() == expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	err = initiate.WriteSettingsFile(config)
	if err != nil {
		t.Fatal(err)
	}

}
