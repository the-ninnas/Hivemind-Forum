package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"text/template"
)

func register(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("static/*.html"))

	var pagedata PageData

	switch r.Method {

	//If user goes to registration page
	case "GET":
		if checkSession(w, r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			tmpl.ExecuteTemplate(w, "reg.html", pagedata)
		}

	//If user presses register button
	case "POST":
		//Parses data inserted to the form
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.ExecuteTemplate(w, "400.html", pagedata)
			return
		}

		//Checks username length, email validity and password length
		errorMsg := validateUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"), r.FormValue("confirmPassword"))
		//If user validation fails, print errorMsg, stay on registration page
		if errorMsg != "" {
			pagedata.Message = errorMsg
			fmt.Println(errorMsg)
			tmpl.ExecuteTemplate(w, "reg.html", pagedata)
			return
		}

		//Add data to DB
		db, err := sql.Open("sqlite3", "./database.db")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			tmpl.ExecuteTemplate(w, "500.html", pagedata)
			return
		}
		stmt, err := db.Prepare(`
			INSERT INTO accounts (username, email, password) values (?, ?, ?);
		`)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			tmpl.ExecuteTemplate(w, "500.html", pagedata)
			return
		}

		//Exec adds data to database
		_, errReg := stmt.Exec(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))
		db.Close()

		//Checks if username or email is already in Database
		//Message should be displayed to user?
		if errReg != nil {
			if errReg.Error() == "UNIQUE constraint failed: accounts.username" {
				errorMsg = "Username already in use. Please choose another username."
			}
			if errReg.Error() == "UNIQUE constraint failed: accounts.email" {
				errorMsg = "Email already in use. Please enter a different email."
			}
			fmt.Println(errorMsg)
			pagedata.Message = errorMsg
			tmpl.ExecuteTemplate(w, "reg.html", pagedata)
			return
		}
		//If registration succeeds, go back to main page
		//Display message to user?
		fmt.Println("Registration successful!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

//Validate registration data
func validateUser(username, email, password, confirmPassword string) string {

	//Username with length of 5-12
	if len(username) > 15 {
		return "Username too long!"
	}
	if len(username) < 5 {
		return "Username too short!"
	}

	//Check if email is valid. (example@something == VALID)
	validEmail, _ := mail.ParseAddress(email)
	if validEmail == nil {
		return "Email address not valid!"
	}

	if len(password) == 0 {
		return "Password missing! This can be done in html with <\\input required>."
	}
	if len(password) < 8 {
		return "Password is too short! Enter at least 8 characters."
	}
	if password != confirmPassword {
		return "Passwords are not matching! Try again."
	}
	return ""
}
