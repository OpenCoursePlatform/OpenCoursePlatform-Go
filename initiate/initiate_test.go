package initiate

import (
	"testing"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	_ "github.com/go-sql-driver/mysql"
)

func TestTests(t *testing.T) {
	_, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = Tests()
	if err == nil {
		t.Errorf("Initiation of tests should have failed. Did not.")
		return
	}

	db, err := helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestFinishTests(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	err = FinishTests(db)
	if err == nil {
		t.Errorf("FinishTests should have failed. Did not.")
		return
	}
}

func TestCreateDatabase(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateDatabase(db, "newTestingDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "newTestingDB")
	if err == nil {
		t.Errorf("CreateDatabase should have failed. Did not.")
		return
	}
	_, err = db.Exec(`
		DROP DATABASE newTestingDB
	`)
	if err != nil {
		t.Errorf("Deletion of database failed. Error message: %s", err.Error())
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateAccessLogs(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateAccessLogs(db)
	if err == nil {
		t.Errorf("CreateAccessLogs should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateAccessLogs(db)
	if err != nil {
		t.Errorf("CreateAccessLogs failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateBlogPosts(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateBlogPosts(db)
	if err == nil {
		t.Errorf("CreateBlogPosts should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateBlogPosts(db)
	if err != nil {
		t.Errorf("CreateBlogPosts failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateCourseCategories(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateCourseCategories(db)
	if err == nil {
		t.Errorf("CreateCourseCategories should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateCourseCategories(db)
	if err != nil {
		t.Errorf("CreateCourseCategories failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateCourses(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateCourses(db)
	if err == nil {
		t.Errorf("CreateCourses should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateCourses(db)
	if err != nil {
		t.Errorf("CreateCourses failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateEmailVerification(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateEmailVerification(db)
	if err == nil {
		t.Errorf("CreateEmailVerification should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateEmailVerification(db)
	if err != nil {
		t.Errorf("CreateEmailVerification failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateFooter(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateFooter(db)
	if err == nil {
		t.Errorf("CreateFooter should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateFooter(db)
	if err != nil {
		t.Errorf("CreateFooter failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateFooterCategories(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateFooterCategories(db)
	if err == nil {
		t.Errorf("CreateFooterCategories should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateFooterCategories(db)
	if err != nil {
		t.Errorf("CreateFooterCategories failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateForumPosts(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateForumPosts(db)
	if err == nil {
		t.Errorf("CreateForumPosts should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateForumPosts(db)
	if err != nil {
		t.Errorf("CreateForumPosts failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateForumTopics(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateForumTopics(db)
	if err == nil {
		t.Errorf("CreateForumTopics should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateForumTopics(db)
	if err != nil {
		t.Errorf("CreateForumTopics failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateGroupMembers(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateGroupMembers(db)
	if err == nil {
		t.Errorf("CreateGroupMembers should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateGroupMembers(db)
	if err != nil {
		t.Errorf("CreateGroupMembers failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateGroupPermissions(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateGroupPermissions(db)
	if err == nil {
		t.Errorf("CreateGroupPermissions should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateGroupPermissions(db)
	if err != nil {
		t.Errorf("CreateGroupPermissions failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateModule(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateModule(db)
	if err == nil {
		t.Errorf("CreateModule should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateModule(db)
	if err != nil {
		t.Errorf("CreateModule failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreatePages(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreatePages(db)
	if err == nil {
		t.Errorf("CreatePages should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreatePages(db)
	if err != nil {
		t.Errorf("CreatePages failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreatePermissions(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreatePermissions(db)
	if err == nil {
		t.Errorf("CreatePermissions should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreatePermissions(db)
	if err != nil {
		t.Errorf("CreatePermissions failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}
func TestCreatePodcastPost(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreatePodcastPost(db)
	if err == nil {
		t.Errorf("CreatePodcastPost should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreatePodcastPost(db)
	if err != nil {
		t.Errorf("CreatePodcastPost failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateSession(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateSession(db)
	if err == nil {
		t.Errorf("CreateSession should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateSession(db)
	if err != nil {
		t.Errorf("CreateSession failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateSessionText(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateSessionText(db)
	if err == nil {
		t.Errorf("CreateSessionText should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateSessionText(db)
	if err != nil {
		t.Errorf("CreateSessionText failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateSessionYoutube(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateSessionYoutube(db)
	if err == nil {
		t.Errorf("CreateSessionYoutube should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateSessionYoutube(db)
	if err != nil {
		t.Errorf("CreateSessionYoutube failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateSettings(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateSettings(db)
	if err == nil {
		t.Errorf("CreateSettings should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateSettings(db)
	if err != nil {
		t.Errorf("CreateSettings failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateTicketResponses(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateTicketResponses(db)
	if err == nil {
		t.Errorf("CreateTicketResponses should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateTicketResponses(db)
	if err != nil {
		t.Errorf("CreateTicketResponses failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateTickets(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateTickets(db)
	if err == nil {
		t.Errorf("CreateTickets should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateTickets(db)
	if err != nil {
		t.Errorf("CreateTickets failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateToolbar(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateToolbar(db)
	if err == nil {
		t.Errorf("CreateToolbar should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateToolbar(db)
	if err != nil {
		t.Errorf("CreateToolbar failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateUserGroups(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateUserGroups(db)
	if err == nil {
		t.Errorf("CreateUserGroups should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateUserGroups(db)
	if err != nil {
		t.Errorf("CreateUserGroups failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateUserPermissions(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateUserPermissions(db)
	if err == nil {
		t.Errorf("CreateUserPermissions should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateUserPermissions(db)
	if err != nil {
		t.Errorf("CreateUserPermissions failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateUsers(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateUsers(db)
	if err == nil {
		t.Errorf("CreateUsers should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateDatabase(db, "testDB")
	if err != nil {
		t.Errorf("CreateDatabase failed. Error message %s", err.Error())
		return
	}
	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		t.Errorf("CreateDBHandlerWithDB failed. Error message %s", err.Error())
		return
	}
	err = CreateUsers(db)
	if err != nil {
		t.Errorf("CreateUsers failed. Error message %s", err.Error())
		return
	}
	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}

func TestCreateTables(t *testing.T) {
	db, err := Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	err = CreateTables(db)
	if err == nil {
		t.Errorf("CreateTables should have failed. Did not.")
		return
	}
	err = FinishTests(db)
	if err != nil {
		t.Errorf("FinishTests failed. Error message %s", err.Error())
		return
	}
}
