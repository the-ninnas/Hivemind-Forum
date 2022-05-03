package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"text/template"
)

//Get all categories
func getCategories() ([]Category, error) {
	var categories []Category
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`SELECT * FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var category Category
	//Add everything to struct. Maybe dont need everything, will delete the rest later.
	for rows.Next() {
		rows.Scan(&category.ID, &category.Title, &category.TitleEncoded, &category.Description)
		categories = append(categories, category)
	}
	return categories, nil
}

//Template for creating category page.
func templateCreateCategory(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("static/*.html"))

	var pagedata PageData
	var errorMsg string
	var err error

	//All categories for {{range .Categories}}
	pagedata.Categories, err = getCategories()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "500.html", pagedata)
		return
	}

	if checkSession(w, r) {
		pagedata.Username = user.Username

		switch r.Method {

		case "GET":

			tmpl.ExecuteTemplate(w, "createCat.html", pagedata)

		case "POST":

			err := r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				tmpl.ExecuteTemplate(w, "400.html", pagedata)
				return
			}

			title := r.FormValue("title")
			description := r.FormValue("description")
			// fmt.Println(title, description)

			//Add data to DB
			db, err := sql.Open("sqlite3", "./database.db")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				tmpl.ExecuteTemplate(w, "500.html", pagedata)
				return
			}
			stmt, err := db.Prepare(`
				INSERT INTO categories (title, title_encoded, description) values (?, ?, ?);
			`)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				tmpl.ExecuteTemplate(w, "500.html", pagedata)
				return
			}

			//Exec adds data to database
			_, errReg := stmt.Exec(title, url.QueryEscape(title), description)
			db.Close()

			if errReg != nil {

				if errReg.Error() == "UNIQUE constraint failed: categories.title" {
					errorMsg = "This category already exists"
				}

				fmt.Println(errorMsg)
				pagedata.Message = errorMsg
				tmpl.ExecuteTemplate(w, "category.html", pagedata)
				return
			}

			fmt.Println("New category created")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	} else {
		pagedata.Message = "Please log in to create a category!"
		pagedata.Categories, err = getCategories()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			tmpl.ExecuteTemplate(w, "500.html", pagedata)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func CategoryDoesNotExist(category string) bool {

	categories, err := getCategories()

	if err != nil {
		fmt.Println(err)
	}

	for _, v := range categories {
		if v.Title == category {
			return false
		}
	}
	return true
}
