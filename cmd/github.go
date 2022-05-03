package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

var githubClientID string = ""
var githubClientSecret string = ""

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scopes"`
}

type githubData struct {
	Email      string
	Primary    bool
	Verified   bool
	Visibility string
}

// Log in with Github
func githubAuth(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("static/*.html"))

	var pagedata PageData
	code := getGithubCode(w, r)
	accessToken := getGithubAccessToken(code)
	githubEmail, githubUsername := getGithubData(accessToken)

	// Check accounts DB if email already exists. If yes, then use the registered username, ID, etc. If not, then create new row with provided username and email.
	err, userID, usernameFromDB := githubAccount(githubEmail, githubUsername)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "500.html", pagedata)
		return
	}

	if usernameFromDB == "" {
		user.Username = githubUsername
	} else {
		user.Username = usernameFromDB
	}

	pagedata.Username = user.Username
	pagedata.Categories, err = getCategories()
	//If login is successful, create cookie
	createSessionForUser(w, r, userID)
	tmpl.ExecuteTemplate(w, "index.html", pagedata)
	fmt.Println(pagedata.Username, "logged in with Github")
}

// User is redirected back from github to our page, with a temporary code included in the URL, which we extract in this func
func getGithubCode(w http.ResponseWriter, r *http.Request) string {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	return r.FormValue("code")
}

//We exchange the code to the acess token
func getGithubAccessToken(code string) string {
	var v OAuthAccessResponse

	//Create request URL.
	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", githubClientID, githubClientSecret, code)

	//Create request object
	req, errReq := http.NewRequest("POST", reqURL, nil)
	if errReq != nil {
		log.Panic("Request creation failed!")
	}

	// Add request header. So we recieve the response in json.
	req.Header.Set("Accept", "application/json")

	// Send HTTP request
	res, errResp := http.DefaultClient.Do(req)
	if errResp != nil {
		log.Panic("Request failed")
	}

	newDec := json.NewDecoder(res.Body)
	err := newDec.Decode(&v)
	if err != nil {
		log.Panic("github oauth error: ", err)
	}
	return v.AccessToken
}

// We use the acess token to fetch user data from github
func getGithubData(acessToken string) (string, string) {
	// Create GET request
	reqEmail, errReq := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if errReq != nil {
		log.Panic("GET REQUEST creation failed!")
	}
	reqUser, errReq := http.NewRequest("GET", "https://api.github.com/user", nil)
	if errReq != nil {
		log.Panic("GET REQUEST creation failed!")
	}

	// Set authorization header
	authHeader := fmt.Sprintf("token %s", acessToken)
	reqEmail.Header.Set("Authorization", authHeader)
	reqUser.Header.Set("Authorization", authHeader)

	// Make request
	resEmail, errResp := http.DefaultClient.Do(reqEmail)
	if errResp != nil {
		log.Panic("Request failed")
	}
	resUser, errResp := http.DefaultClient.Do(reqUser)
	if errResp != nil {
		log.Panic("Request failed")
	}

	// Read response as a byte slice
	resBodyEmail, _ := ioutil.ReadAll(resEmail.Body)
	resBodyUser, _ := ioutil.ReadAll(resUser.Body)

	var email []githubData
	var result map[string]interface{}
	var githubEmail string
	var githubUsername string

	json.Unmarshal([]byte(resBodyEmail), &email)
	json.Unmarshal([]byte(resBodyUser), &result)

	for _, v := range email {
		if v.Primary == true {
			githubEmail = v.Email
		}
	}

	for i, v := range result {
		if i == "login" {
			githubUsername = v.(string)
		}
	}
	return githubEmail, githubUsername
}

func githubAccount(githubEmail, githubUsername string) (error, int, string) {
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
	_, errReg := stmt.Exec(githubUsername, githubEmail)

	//If this email is already in use, then get the username and ID from the DB
	if errReg != nil {
		if errReg.Error() == "UNIQUE constraint failed: accounts.email" {
			rows, err := db.Query("SELECT * FROM accounts WHERE email = ?", githubEmail)
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
		rows, err := db.Query("SELECT * FROM accounts WHERE email = ?", githubEmail)
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
