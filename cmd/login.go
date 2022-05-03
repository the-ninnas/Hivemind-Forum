package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
)

//Checks if inserted username exists in DB and if passwords match.
//Creates login session cookie
func login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("static/*.html"))

	var pagedata PageData
	var err error
	pagedata.Categories, err = getCategories()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "500.html", pagedata)
		return
	}

	//If you try to access /login through URL, it will redirect to "/" (main page)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//Parse inserted data
	error := r.ParseForm()
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		tmpl.ExecuteTemplate(w, "400.html", pagedata)
		return
	}

	//Define inserted username and password
	usernameLogin := r.FormValue("username")
	passwordLogin := r.FormValue("password")

	//Find username in account DB and check if passwords match
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "500.html", pagedata)
		return
	}

	//Get all info from database associated to the inserted username
	rows, err := db.Query("SELECT * FROM accounts WHERE username = ?", usernameLogin)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "500.html", pagedata)
		return
	}
	defer rows.Close()

	var id int
	var username string
	var email string
	var password string

	//Add everything to struct
	for rows.Next() {
		rows.Scan(&id, &username, &email, &password)
		user = User{
			ID:       id,
			Username: username,
			Email:    email,
			Password: password,
		}
	}
	//If username doesn't match, user is not created, so len(user) == 0
	if user.Username != usernameLogin {

		pagedata.Message = "This username is not registered"
		fmt.Println("This username is not registered")
		tmpl.ExecuteTemplate(w, "index.html", pagedata)
		return
	}

	//If username is in database, match passwords
	if passwordLogin == user.Password {
		fmt.Println("Login successful!")

	} else {
		pagedata.Message = "You entered a wrong password!"
		fmt.Println("You entered a wrong password!")
		tmpl.ExecuteTemplate(w, "index.html", pagedata)
		return
	}

	pagedata.Username = user.Username
	//If login is successful, create cookie
	createSessionForUser(w, r, user.ID)
	tmpl.ExecuteTemplate(w, "index.html", pagedata)
}
