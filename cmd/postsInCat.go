package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
)

func templatePostsInCategory(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("static/*.html"))

	var pagedata PageData
	pagedata.Google = googleClientID
	pagedata.Github = githubClientID
	var err error

	// Get decoded category name from URL, user readable.
	category, _ := url.QueryUnescape(strings.TrimPrefix(r.URL.Path, "/c/"))

	if CategoryDoesNotExist(category) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// Single category variable to write category name on the webpage
	pagedata.Category = category

	url := r.URL.RequestURI()
	//Encoded category name for different buttons with links to this category
	pagedata.CategoryEncoded = strings.TrimPrefix(url, "/c/")

	//All categories for {{range .Categories}}
	pagedata.Categories, err = getCategories()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "500.html", pagedata)
		return
	}

	if r.URL.Path == "/c/" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	switch r.Method {

	case "GET":

		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.ExecuteTemplate(w, "400.html", pagedata)
			return
		}

		//Check if logged in
		if checkSession(w, r) {
			pagedata.Username = user.Username

			if strings.HasSuffix(url, "posted=true") {
				pagedata.Posts, err = getPostByUsernameAndCat(user.Username, category)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					tmpl.ExecuteTemplate(w, "500.html", pagedata)
					return
				}

				tmpl.ExecuteTemplate(w, "postsInCat.html", pagedata)

			} else if strings.HasSuffix(url, "liked=true") {
				pagedata.Posts, err = getPostByUserLikesAndCat(user.Username, category)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					tmpl.ExecuteTemplate(w, "500.html", pagedata)
					return
				}

				tmpl.ExecuteTemplate(w, "postsInCat.html", pagedata)

			} else {
				pagedata.Posts, err = getPostByCategoryName(category)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					tmpl.ExecuteTemplate(w, "500.html", pagedata)
					return
				}

				tmpl.ExecuteTemplate(w, "postsInCat.html", pagedata)
			}
		} else {
			pagedata.Posts, err = getPostByCategoryName(category)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				tmpl.ExecuteTemplate(w, "500.html", pagedata)
				return
			}

			tmpl.ExecuteTemplate(w, "postsInCat.html", pagedata)
		}

	case "POST":

		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.ExecuteTemplate(w, "400.html", pagedata)
			return
		}

		//Form values for like/dislike, post title/content
		title := r.FormValue("title")
		button := r.FormValue("submit")
		content := strings.Replace(r.FormValue("content"), "\n", "<br>", -1)
		postid := r.FormValue("postid")
		postlike := r.FormValue("postlike")
		postdislike := r.FormValue("postdislike")

		//If logged in
		if checkSession(w, r) {

			err := r.ParseForm()
			if err != nil {
				log.Panic(err)
			}
			pagedata.Username = user.Username

			//If like or dislike is pressed, will call postvote func
			if postlike != "" || postdislike != "" {
				postVote(postlike, postdislike, postid)
			}

			//If user wants to post
			if button == "submit-post" {

				if title == "" {
					pagedata.Message = "Title missing!"
					fmt.Println("Title missing!")

				} else {
					addPost(title, content, category)
				}
			}

			if strings.HasSuffix(url, "posted=true") {
				pagedata.Posts, err = getPostByUsernameAndCat(user.Username, category)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					tmpl.ExecuteTemplate(w, "500.html", pagedata)
					return
				}

				tmpl.ExecuteTemplate(w, "postsInCat.html", pagedata)
				return
			}

			if strings.HasSuffix(url, "liked=true") {
				pagedata.Posts, err = getPostByUserLikesAndCat(user.Username, category)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					tmpl.ExecuteTemplate(w, "500.html", pagedata)
					return
				}

				tmpl.ExecuteTemplate(w, "postsInCat.html", pagedata)
				return
			}

			//Get all the posts in that category, including the one just posted and execute template to render the page
			pagedata.Posts, err = getPostByCategoryName(category)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				tmpl.ExecuteTemplate(w, "500.html", pagedata)
				return
			}

			tmpl.ExecuteTemplate(w, "postsInCat.html", pagedata)

		} else {

			// If not logged in, cannot like/dislike
			if postlike != "" || postdislike != "" {

				pagedata.Message = "Pleale log in to like/dislike"
				fmt.Println("Please log in to like/dislike")
			}

			if button == "submit-post" {
				pagedata.Message = "Please log in to post"
				fmt.Println("Please log in to post")
			}

			pagedata.Posts, err = getPostByCategoryName(category)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				tmpl.ExecuteTemplate(w, "500.html", pagedata)
				return
			}

			tmpl.ExecuteTemplate(w, "postsInCat.html", pagedata)
		}
	}
}
