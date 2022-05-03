package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

//When loggin in, creates session
func createSessionForUser(w http.ResponseWriter, r *http.Request, userid int) {
	//Create UUID - universal unique identifier
	sessionToken := uuid.NewV4()
	time := time.Now().Add(time.Second * 120)
	//Sets cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  sessionToken.String(),
		Path:   "/",
		MaxAge: 600,
	})
	//Assign session in DB
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Delete old session if exists. IF using two browsers, will log out from first
	db.Exec(`DELETE FROM session WHERE user_id = ?`, userid)

	stmt, err := db.Prepare(`INSERT INTO session (time, user_id, uuid) values (?, ?, ?);`)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Exec adds data to database
	_, errReg := stmt.Exec(time, userid, sessionToken.String())
	db.Close()

	if errReg != nil {
		log.Panic(errReg)
	}
}

//Check if cookies exist, so we know if user is logged in or not.
func checkSession(w http.ResponseWriter, r *http.Request) bool {
	//Obtain cookie
	userid := user.ID

	c, err := r.Cookie("session_token")
	if err != nil {
		return false
	}
	//Compare sessionToken to uuid in session DB
	sessionToken := c.Value

	// Get session from our database
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	rows, err := db.Query("SELECT uuid FROM session WHERE user_id = ?", userid)
	if err != nil {
		log.Fatal(err)
		// return false
	}
	defer rows.Close()

	var uuid string
	var sessionUuid string

	for rows.Next() {
		rows.Scan(&uuid)
		sessionUuid = uuid
	}
	return sessionUuid == sessionToken
}
