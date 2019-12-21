package admin

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

// FilesPageData ...
type FilesPageData struct {
	Title         string
	Paths         []helpers.Path
	Files         []string
	NewItemButton bool
	NewItemText   string
	NewItemLink   string
}

// Files ...
func Files(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./static")
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	var filesList []string
	for _, f := range files {
		if f.Mode().IsDir() == false {
			filesList = append(filesList, f.Name())
		}
	}
	p := &FilesPageData{
		Title:         "Files",
		Paths:         []helpers.Path{{Name: "Admin", Link: "/admin"}},
		Files:         filesList,
		NewItemButton: true,
		NewItemText:   "Add new File",
		NewItemLink:   "/files/new",
	}
	helpers.RenderTemplate(r, w, "admin_files", p)
}

// NewFile ...
func NewFile(w http.ResponseWriter, r *http.Request) {
	p := &Page{
		Title: "New file",
		Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}, {Name: "Files", Link: "/admin/files"}},
	}
	helpers.RenderTemplate(r, w, "admin_new_file", p)
}

// UploadNewFile ...
func UploadNewFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		helpers.HandleError(err)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		helpers.HandleError(err)
		return
	}

	fileNameParts := strings.Split(handler.Filename, ".")
	fileName := ""
	for index := 0; index < len(fileNameParts)-1; index++ {
		fileName = fileName + fileNameParts[index]
	}

	fileName = helpers.GenerateSlug(fileName) + "." + fileNameParts[len(fileNameParts)-1]

	defer file.Close()
	f, err := os.OpenFile("./static/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		helpers.HandleError(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	http.Redirect(w, r, "/admin/files/"+fileName, http.StatusSeeOther)
}

// SingleFile ...
func SingleFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["file"]
	p := &Page{
		Title: slug,
		Paths: []helpers.Path{{Name: "Admin", Link: "/admin"}, {Name: "Files", Link: "/admin/files"}},
	}
	helpers.RenderTemplate(r, w, "admin_file", p)
}

// RemoveFile ...
func RemoveFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["file"]
	var err = os.Remove("./static/" + slug)
	if err != nil {
		helpers.HandleError(err)
		return
	}
	http.Redirect(w, r, "/admin/files", http.StatusSeeOther)
}
