package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/satori/go.uuid"
)

func main() {

	createDB()

	port := ":8080"
	fmt.Println("Click here: http://localhost" + port)

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	http.HandleFunc("/", templateIndex)
	http.HandleFunc("/post", templateCreatePost)
	http.HandleFunc("/c/", templatePostsInCategory)
	http.HandleFunc("/category", templateCreateCategory)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/comments/", commentTemplate)
	http.HandleFunc("/googlecallback", googleAuth)
	http.HandleFunc("/githubcallback", githubAuth)

	//For later use
	// http.HandleFunc("/user/", templateUser)
	// http.HandleFunc("/posted/", templatePosted)
	// http.HandleFunc("/liked/", templateLiked)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
