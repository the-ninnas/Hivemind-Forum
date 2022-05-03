package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
)

var googleClientID string = ""
var googleClientSecret string = ""

type googleData struct {
	Id             string
	Email          string
	Verified_email bool
	Picture        string
}

// Log in with google
func googleAuth(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("static/*.html"))

	var pagedata PageData
	code := getGoogleCode(w, r)
	acessToken := getGoogleAccessToken(code)
	googleEmail := getGoogleData(acessToken)
	googleUsername := strings.Split(googleEmail, "@")[0]

	// Check accounts DB if email already exists. If yes, then use the registered username, ID, etc. If not, then create new row with provided username and email.
	err, userID, usernameFromDB := googleAccount(googleEmail, googleUsername)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "500.html", pagedata)
		return
	}

	if usernameFromDB == "" {
		user.Username = googleUsername
	} else {
		user.Username = usernameFromDB
	}

	user.Username = usernameFromDB
	pagedata.Username = user.Username
	pagedata.Categories, err = getCategories()

	//If login is successful, create cookie
	createSessionForUser(w, r, userID)
	tmpl.ExecuteTemplate(w, "index.html", pagedata)
	fmt.Println(pagedata.Username, "logged in with Google")
}

func getGoogleData(accessToken string) string {
	var googleData googleData

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		log.Panic("Request failed: ", err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(content), &googleData)

	return googleData.Email
}

func getGoogleAccessToken(code string) string {
	var v OAuthAccessResponse

	//Create request URL.
	reqURL := fmt.Sprintf("https://oauth2.googleapis.com/token")
	data := url.Values{}

	data.Set("code", code)
	data.Set("client_id", googleClientID)
	data.Set("client_secret", googleClientSecret)
	data.Set("redirect_uri", "http://localhost:8080/googlecallback")
	data.Set("grant_type", "authorization_code")

	res, err := http.PostForm(reqURL, data)
	if err != nil {
		log.Panic("Request failed ", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	bodyString := string(body)

	json.Unmarshal([]byte(bodyString), &v)
	return v.AccessToken
}

func getGoogleCode(w http.ResponseWriter, r *http.Request) string {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	return r.FormValue("code")
}

func googleAccount(googleEmail, googleUsername string) (error, int, string) {
	var id int
	var username string
	var email string
	var password string

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err, 0, ""
	}

	stmt, err := db.Prepare(`
    INSERT INTO accounts (username, email) values (?, ?);
    `)
	if err != nil {
		return err, 0, ""
	}

	//Exec adds data to database
	_, errReg := stmt.Exec(googleUsername, googleEmail)

	//If this email is already in use, then get the username and ID from the DB
	if errReg != nil {
		if errReg.Error() == "UNIQUE constraint failed: accounts.email" {
			rows, err := db.Query("SELECT * FROM accounts WHERE email = ?", googleEmail)
			if err != nil {
				return err, 0, ""
			}

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
		}
		// If the email was not in use already, then get the newly added username and ID from the DB.
	} else {
		rows, err := db.Query("SELECT * FROM accounts WHERE email = ?", googleEmail)
		if err != nil {
			return err, 0, ""
		}

		for rows.Next() {
			rows.Scan(&id, &username, &email, &password)
			user = User{
				ID:       id,
				Username: username,
				Email:    email,
				Password: password,
			}
		}
	}
	db.Close()

	return err, user.ID, user.Username
}
