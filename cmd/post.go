package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

func getPostByUsernameAndCat(username, category string) ([]Post, error) {
	var posts []Post

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
	SELECT * 
	FROM posts 
	WHERE username = ? 
	AND category = ? 
	ORDER BY time DESC
	`, username, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var post Post

	for rows.Next() {
		rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Likes, &post.Time, &post.Username, &post.Picture)
		posts = append(posts, post)
	}
	return posts, nil
}

//For user posted filter
func getPostByUsername(username string) ([]Post, error) {
	var posts []Post

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
	SELECT * 
	FROM posts 
	WHERE username = ?
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var post Post

	for rows.Next() {
		rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Likes, &post.Time, &post.Username, &post.Picture)
		posts = append(posts, post)
	}
	return posts, nil
}

//For user liked filter
func getPostByUserLikesAndCat(username, category string) ([]Post, error) {
	var posts []Post

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
		SELECT
			p.ID, title, content, category, likes, p.time, p.username, p.picture
		FROM posts p, postvotes pv
		WHERE p.ID = pv.post_id
		AND pv.username = ? AND pv.type = "like" AND p.category = ?
		ORDER BY pv.time DESC
		`, username, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var post Post

	for rows.Next() {
		rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Likes, &post.Time, &post.Username, &post.Picture)
		posts = append(posts, post)
	}
	return posts, nil
}

//For user liked filter
func getPostByUserLikes(username string) ([]Post, error) {
	var posts []Post

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
		SELECT
			p.ID, title, content, category, likes, p.time, p.username, p.picture
		FROM posts p, postvotes pv
		WHERE p.ID = pv.post_id
		AND pv.username = ? AND pv.type = "like" 
		ORDER BY pv.time DESC
		`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var post Post

	for rows.Next() {
		rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Likes, &post.Time, &post.Username, &post.Picture)
		posts = append(posts, post)
	}
	return posts, nil
}

func getPostByID(ID int) (Post, error) {
	var post Post

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return post, err
	}

	rows, err := db.Query(`SELECT * FROM posts WHERE ID = ?`, ID)
	if err != nil {
		return post, err
	}
	defer rows.Close()

	//Add everything to struct. Maybe dont need everything, will delete the rest later.
	for rows.Next() {
		rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Likes, &post.Time, &post.Username, &post.Picture)

	}
	return post, nil
}

//If you press on a category on the main page, it will look for all the posts under that category
func getPostByCategoryName(category string) ([]Post, error) {
	var posts []Post
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
	SELECT * 
	FROM posts 
	WHERE category = ?
	ORDER BY time DESC
	`, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var post Post

	for rows.Next() {
		rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Likes, &post.Time, &post.Username, &post.Picture)
		posts = append(posts, post)
	}
	return posts, nil
}

func templateCreatePost(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("static/*.html"))

	var pagedata PageData
	var err error
	pagedata.Categories, err = getCategories()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "500.html", pagedata)
		return
	}

	switch r.Method {

	//If user wants to access post page trough URL
	case "GET":

		//Checks if user is logged in
		if checkSession(w, r) {

			pagedata.Username = user.Username
			tmpl.ExecuteTemplate(w, "postInManyCat.html", pagedata)

		} else {
			// If user wants to access post page without logging in, send to main/login page
			pagedata.Message = "Please log in"
			http.Redirect(w, r, "/", http.StatusSeeOther)
			fmt.Println("Please log in")
		}

	case "POST":

		if checkSession(w, r) {

			pagedata.Username = user.Username
			title := r.FormValue("title")
			content := r.FormValue("content")
			category := r.Form["category"]

			if title == "" {
				pagedata.Message = "Missing title!"
				tmpl.ExecuteTemplate(w, "postInManyCat.html", pagedata)
				return
			}
			if len(category) == 0 {
				pagedata.Message = "Please choose or create a category!"
				tmpl.ExecuteTemplate(w, "postInManyCat.html", pagedata)
				return
			}

			//Upload image handler
			picture, uploadErr := uploadImageHandler(w, r)
			if uploadErr != "" {
				pagedata.Message = uploadErr
				tmpl.ExecuteTemplate(w, "postInManyCat.html", pagedata)
				return
			}

			addPostManyCat(title, content, picture, category)
			http.Redirect(w, r, "/", http.StatusSeeOther)

		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func addPostManyCat(title, content, picture string, category []string) {

	likes := 0
	time := time.Now().Format("2006-01-02 15:04")

	//Open database
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := db.Prepare(`
		INSERT INTO posts (title, content, category, likes, time, username, picture) values (?, ?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		log.Fatal(err)
	}

	//Exec adds data to database
	for _, cat := range category {
		_, errPost := stmt.Exec(title, content, cat, likes, time, user.Username, picture)
		if errPost != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Posting succeeded")
}

func addPost(title, content, category string) {
	likes := 0
	time := time.Now().Format("2006-01-02 15:04")

	//Open database
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := db.Prepare(`
		INSERT INTO posts (title, content, category, likes, time, username) values (?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		log.Fatal(err)
	}

	//Exec adds data to database
	_, errPost := stmt.Exec(title, content, category, likes, time, user.Username)
	if errPost != nil {
		log.Fatal(err)
	}
	fmt.Println("Posted")
}
