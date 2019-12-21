package forum

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/initiate"
	_ "github.com/go-sql-driver/mysql"
)

func TestGetForumTopicsAndPost(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
		INSERT INTO forum_posts (id, text, topic_id, author_id)
		VALUES
		(5, 'I am new and inexperienced in oro crm/platform and this is my first project I am working on.\n\nI have created a new bundle with menu but that menu item doesn’t showing up anymore. Can anyone help me with the exact issue please? I tried with different approaches but no luck yet.', 4, 1);
	`)
	if err != nil {
		t.Errorf("Insertion of forum_post in database failed. Error message: %s", err.Error())
		return
	}
	_, err = db.Exec(`
		INSERT INTO forum_posts (id, text, topic_id, author_id)
		VALUES
		(6, 'Please follow the documentation https://doc.oroinc.com/backend/navigation/\nAnd don’t forget to clear the cache after adding a new configuration file to make the application know about it.', 4, 1);
	`)
	if err != nil {
		t.Errorf("Insertion of forum_post in database failed. Error message: %s", err.Error())
		return
	}

	_, err = db.Exec(`
		INSERT INTO forum_topics (id, title, slug)
		VALUES
		(4, 'Custom Menu not showing up', 'custom-menu');
	`)
	if err != nil {
		t.Errorf("Insertion of forum_topics in database failed. Error message: %s", err.Error())
		return
	}

	forumPosts, err := GetForumTopicsAndPost(db)
	if err != nil {
		t.Errorf("GetForumTopicsAndPost failed. Error message %s", err.Error())
		return
	}

	sampleForumPosts := []FPost{
		{"Custom Menu not showing up", "I am new and inexperienced in oro crm/platform and this is my first project I am working on.\n\nI have created a new bundle with menu but that menu item doesn’t showing up anymore. Can anyone help me with the exact issue please? I tried with different approaches but no luck yet.", "custom-menu"},
		{"Custom Menu not showing up", "Please follow the documentation https://doc.oroinc.com/backend/navigation/\nAnd don’t forget to clear the cache after adding a new configuration file to make the application know about it.", "custom-menu"},
	}

	if forumPosts[0] != sampleForumPosts[0] {
		t.Errorf("Data is not equal to sample data.")
	}

	if forumPosts[0] == sampleForumPosts[1] {
		t.Errorf("Data should not be equal to sample data")
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestGetForumPosts(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
		INSERT INTO forum_posts (id, text, topic_id, author_id)
		VALUES
		(5, 'I am new and inexperienced in oro crm/platform and this is my first project I am working on.\n\nI have created a new bundle with menu but that menu item doesn’t showing up anymore. Can anyone help me with the exact issue please? I tried with different approaches but no luck yet.', 4, 1);
	`)
	if err != nil {
		t.Errorf("Insertion of forum_post in database failed. Error message: %s", err.Error())
		return
	}
	_, err = db.Exec(`
		INSERT INTO forum_posts (id, text, topic_id, author_id)
		VALUES
		(6, 'Please follow the documentation https://doc.oroinc.com/backend/navigation/\nAnd don’t forget to clear the cache after adding a new configuration file to make the application know about it.', 4, 1);
	`)
	if err != nil {
		t.Errorf("Insertion of forum_post in database failed. Error message: %s", err.Error())
		return
	}

	_, err = db.Exec(`
		INSERT INTO forum_topics (id, title, slug)
		VALUES
		(4, 'Custom Menu not showing up', 'custom-menu');
	`)
	if err != nil {
		t.Errorf("Insertion of forum_topics in database failed. Error message: %s", err.Error())
		return
	}

	forumPosts, err := GetForumPosts(db, "custom-menu")
	if err != nil {
		t.Errorf("GetForumTopicsAndPost failed. Error message %s", err.Error())
		return
	}

	sampleForumPosts := []FPost{
		{"Custom Menu not showing up", "I am new and inexperienced in oro crm/platform and this is my first project I am working on.\n\nI have created a new bundle with menu but that menu item doesn’t showing up anymore. Can anyone help me with the exact issue please? I tried with different approaches but no luck yet.", "custom-menu"},
		{"Custom Menu not showing up", "Please follow the documentation https://doc.oroinc.com/backend/navigation/\nAnd don’t forget to clear the cache after adding a new configuration file to make the application know about it.", "custom-menu"},
	}

	if forumPosts[0] != sampleForumPosts[0] {
		t.Errorf("Data is not equal to sample data.")
	}

	if forumPosts[1] != sampleForumPosts[1] {
		t.Errorf("Data should not be equal to sample data")
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestPostPage(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostPage)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestForum(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Forum)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestDBForPostPage(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostPage)

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

func TestDBForForum(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Forum)

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

func TestDBForNewForumTopic(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NewForumTopic)

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

func TestDBForInsertNewTopic(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(InsertNewTopic)

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

func TestDBForAnswerPost(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AnswerPost)

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
