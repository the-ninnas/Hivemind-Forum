package main

import (
	"net/http"
	"text/template"
)

func templateIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("static/*.html"))

	var pagedata PageData
	pagedata.Google = googleClientID
	pagedata.Github = githubClientID
	var err error
	pagedata.Categories, err = getCategories()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "500.html", pagedata)
		return
	}

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		tmpl.ExecuteTemplate(w, "404.html", pagedata)
		return
	}

	//Check if cookies exist/user is logged in
	if checkSession(w, r) {

		pagedata.Username = user.Username
		tmpl.ExecuteTemplate(w, "index.html", pagedata)

	} else {
		// If not logged in
		tmpl.ExecuteTemplate(w, "index.html", pagedata)
	}
}
