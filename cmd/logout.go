package main

import (
	"database/sql"
	"log"
	"net/http"
)

func logout(w http.ResponseWriter, r *http.Request) {
	userid := user.ID

	//If logged in in two brosers, will not log out from second
	if checkSession(w, r) {
		//Delete cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  "",
			MaxAge: -1,
		})
		//Delete session from db
		db, err := sql.Open("sqlite3", "./database.db")
		if err != nil {
			log.Panic(err)
		}
		db.Exec(`DELETE FROM session WHERE user_id = ?`, userid)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
