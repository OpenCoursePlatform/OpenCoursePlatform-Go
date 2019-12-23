package server

import (
	"net/http"
	"time"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/admin"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/authentication"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/blog"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/course"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/forum"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/index"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/middleware"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/pages"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/podcast"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/settings"
	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/support"
	"github.com/gorilla/mux"
)

/*
Make creates the router and adds all of the paths.
*/
func Make() *http.Server {
	router := mux.NewRouter()

	router.HandleFunc("/", middleware.Chain(index.Page, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/sign/in", middleware.Chain(authentication.SignIn, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/sign/in", middleware.Chain(authentication.Authenticating, middleware.Logging(), middleware.ConfigFileExists())).Methods("POST")
	router.HandleFunc("/sign/out", middleware.Chain(authentication.SignOut, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/sign/up", middleware.Chain(authentication.SignUp, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/sign/up", middleware.Chain(authentication.SignUpNewUser, middleware.Logging(), middleware.ConfigFileExists())).Methods("POST")
	router.HandleFunc("/activate/{token}", middleware.Chain(authentication.ActivateAccount, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")

	router.HandleFunc("/settings", middleware.Chain(settings.UserPage, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/settings/email", middleware.Chain(settings.EmailPage, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/settings/email", middleware.Chain(settings.UpdateEmail, middleware.Logging(), middleware.ConfigFileExists())).Methods("POST")
	router.HandleFunc("/settings/password", middleware.Chain(settings.PasswordPage, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/settings/password", middleware.Chain(settings.UpdatePassword, middleware.Logging(), middleware.ConfigFileExists())).Methods("POST")

	router.HandleFunc("/forum", middleware.Chain(forum.Forum, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/forum/new", middleware.Chain(forum.NewForumTopic, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/forum/new", middleware.Chain(forum.InsertNewTopic, middleware.Logging(), middleware.ConfigFileExists())).Methods("POST")
	router.HandleFunc("/forum/{topic}", middleware.Chain(forum.PostPage, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/forum/{topic}", middleware.Chain(forum.AnswerPost, middleware.Logging(), middleware.ConfigFileExists())).Methods("POST")

	router.HandleFunc("/blog", middleware.Chain(blog.Blog, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/blog/{post}", middleware.Chain(blog.Post, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/courses/{course}", middleware.Chain(course.Page, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/courses/{course}/{module}", middleware.Chain(course.ModulePage, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/courses/{course}/{module}/{session}", middleware.Chain(course.SessionPage, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/courses/{course}/{module}/{session}", middleware.Chain(course.SessionAnswer, middleware.Logging(), middleware.ConfigFileExists())).Methods("POST")

	router.HandleFunc("/podcast", middleware.Chain(podcast.Podcast, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/podcast/{post}", middleware.Chain(podcast.Post, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")

	router.HandleFunc("/support", middleware.Chain(support.Support, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/support/new", middleware.Chain(support.NewTicket, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/support/new", middleware.Chain(support.InsertNewTicket, middleware.Logging(), middleware.ConfigFileExists())).Methods("POST")
	router.HandleFunc("/support/{ticket}", middleware.Chain(support.Ticket, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")
	router.HandleFunc("/support/{ticket}", middleware.Chain(support.InsertNewTicketResponse, middleware.Logging(), middleware.ConfigFileExists())).Methods("POST")

	router.HandleFunc("/admin", middleware.Chain(admin.Index, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/categories", middleware.Chain(admin.Categories, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/categories/new", middleware.Chain(admin.NewCategory, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/categories/new", middleware.Chain(admin.InsertNewCategory, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")
	router.HandleFunc("/admin/categories/{category}", middleware.Chain(admin.Category, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/categories/{category}", middleware.Chain(admin.UpdateCategory, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")

	router.HandleFunc("/admin/courses", middleware.Chain(admin.Courses, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/courses/new", middleware.Chain(admin.NewCourse, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/courses/new", middleware.Chain(admin.InsertNewCourse, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")
	router.HandleFunc("/admin/courses/{course}", middleware.Chain(admin.Course, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/courses/{course}", middleware.Chain(admin.UpdateCourse, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")
	router.HandleFunc("/admin/courses/{course}/delete", middleware.Chain(admin.DeleteCourse, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")

	router.HandleFunc("/admin/courses/{course}/new", middleware.Chain(admin.NewModule, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/courses/{course}/new", middleware.Chain(admin.InsertNewModule, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")
	router.HandleFunc("/admin/courses/{course}/{module}", middleware.Chain(admin.Module, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/courses/{course}/{module}", middleware.Chain(admin.UpdateModule, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")
	router.HandleFunc("/admin/courses/{course}/{module}/delete", middleware.Chain(admin.DeleteModule, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")

	router.HandleFunc("/admin/courses/{course}/{module}/new", middleware.Chain(admin.NewSession, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/courses/{course}/{module}/new", middleware.Chain(admin.InsertNewSession, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")
	router.HandleFunc("/admin/courses/{course}/{module}/{session}", middleware.Chain(admin.Session, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/courses/{course}/{module}/{session}", middleware.Chain(admin.UpdateSession, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")

	router.HandleFunc("/admin/blog", middleware.Chain(admin.Blog, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/blog/new", middleware.Chain(admin.NewBlogPost, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/blog/new", middleware.Chain(admin.InsertNewBlogPost, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")
	router.HandleFunc("/admin/blog/{post}", middleware.Chain(admin.BlogPost, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/blog/{post}", middleware.Chain(admin.UpdateBlogPost, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")

	router.HandleFunc("/admin/podcast", middleware.Chain(admin.Podcast, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/podcast/new", middleware.Chain(admin.NewPodcastPost, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/podcast/new", middleware.Chain(admin.InsertNewPodcastPost, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")
	router.HandleFunc("/admin/podcast/{post}", middleware.Chain(admin.PodcastPost, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/podcast/{post}", middleware.Chain(admin.UpdatePodcastPost, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")

	router.HandleFunc("/admin/settings", middleware.Chain(admin.Settings, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/settings/new", middleware.Chain(admin.NewSetting, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/settings/new", middleware.Chain(admin.InsertNewSetting, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")
	router.HandleFunc("/admin/settings/{setting}", middleware.Chain(admin.Setting, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/settings/{setting}", middleware.Chain(admin.UpdateSetting, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")

	router.HandleFunc("/admin/tickets", middleware.Chain(admin.Tickets, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/tickets/{ticket}", middleware.Chain(admin.Ticket, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/tickets/{ticket}", middleware.Chain(admin.UpdateTicket, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")

	router.HandleFunc("/admin/users", middleware.Chain(admin.Users, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/users/{user}", middleware.Chain(admin.User, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/users/{user}", middleware.Chain(admin.Index, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")

	router.HandleFunc("/admin/groups", middleware.Chain(admin.Groups, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/groups/new", middleware.Chain(admin.NewGroup, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/groups/new", middleware.Chain(admin.InsertNewGroup, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")
	router.HandleFunc("/admin/groups/{group}", middleware.Chain(admin.Group, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/groups/{group}", middleware.Chain(admin.UpdateGroup, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")

	router.HandleFunc("/admin/toolbar", middleware.Chain(admin.Toolbar, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/toolbar", middleware.Chain(admin.UpdateToolbar, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")

	router.HandleFunc("/admin/footer", middleware.Chain(admin.FooterCategories, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/footer/{group}", middleware.Chain(admin.FooterCategory, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/footer/{group}", middleware.Chain(admin.UpdateFooter, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")

	router.HandleFunc("/admin/pages", middleware.Chain(admin.Pages, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/pages/new", middleware.Chain(admin.NewPage, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/pages/new", middleware.Chain(admin.InsertNewPage, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")
	router.HandleFunc("/admin/pages/{page}", middleware.Chain(admin.SinglePage, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/pages/{page}", middleware.Chain(admin.UpdatePage, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")

	router.HandleFunc("/admin/files", middleware.Chain(admin.Files, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/files/new", middleware.Chain(admin.NewFile, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/files/new", middleware.Chain(admin.UploadNewFile, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")
	router.HandleFunc("/admin/files/{file}", middleware.Chain(admin.SingleFile, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("GET")
	router.HandleFunc("/admin/files/{file}", middleware.Chain(admin.RemoveFile, middleware.Logging(), middleware.ConfigFileExists(), middleware.AdminPage())).Methods("POST")

	fileServer := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", middleware.NeuterAndLog(fileServer)))
	router.HandleFunc("/{page}", middleware.Chain(pages.SinglePage, middleware.Logging(), middleware.ConfigFileExists())).Methods("GET")

	server := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        router,
	}
	return server
}
