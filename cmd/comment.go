package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// comments.html page
func commentTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("static/*.html"))

	var pagedata PageData
	pagedata.Google = googleClientID
	pagedata.Github = githubClientID
	var err error

	//get post id
	postid, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comments/postid="))
	if err != nil {
		log.Panic(err)
	}

	pagedata.Post, err = getPostByID(postid)
	pagedata.Comments, err = getCommentByPostId(postid)

	//All categories for {{range .Categories}}
	pagedata.Categories, err = getCategories()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.ExecuteTemplate(w, "500.html", pagedata)
		return
	}

	//Encoded category name for link back to that Category
	pagedata.CategoryEncoded = url.QueryEscape(pagedata.Post.Category)

	switch r.Method {
	//If user presses on a post link, will go to post with comments page
	//a href method has to be GET, or use form
	case "GET":

		if checkSession(w, r) {
			pagedata.Username = user.Username
			tmpl.ExecuteTemplate(w, "comments.html", pagedata)

		} else {
			tmpl.ExecuteTemplate(w, "comments.html", pagedata)
		}

		//Inserting comment on the comments.html form
	case "POST":

		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.ExecuteTemplate(w, "400.html", pagedata)
			return
		}
		//Form values for like/dislike
		content := strings.Replace(r.FormValue("comment"), "\n", "<br>", -1)
		commentid := r.FormValue("commentid")
		commentlike := r.FormValue("commentlike")
		commentdislike := r.FormValue("commentdislike")

		//If user is logged in, allow commenting, and voting
		if checkSession(w, r) {

			//If like or dislike is pressed, will call postvote func
			if commentlike != "" || commentdislike != "" {
				commentVote(commentlike, commentdislike, commentid)
			}

			//add comment
			if content != "" {
				addComment(content, postid)
			}
			pagedata.Username = user.Username
			pagedata.Comments, err = getCommentByPostId(postid)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				tmpl.ExecuteTemplate(w, "500.html", pagedata)
				return
			}

			tmpl.ExecuteTemplate(w, "comments.html", pagedata)

		} else {

			// If user wants to comment without logging in, send to main/login page
			if content != "" {
				pagedata.Message = "Please log in to comment"
				fmt.Println("Please log in to comment")
			}

			if commentlike != "" || commentdislike != "" {
				pagedata.Message = "Please log in to vote"
				fmt.Println("Please log in to vote")
			}

			tmpl.ExecuteTemplate(w, "comments.html", pagedata)
		}
	}
}

//getCommentsByPostId
func getCommentByPostId(ID int) ([]Comment, error) {

	var comments []Comment
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM comments WHERE post_id = ?", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comment Comment
	//Add everything to struct. Maybe dont need everything, will delete the rest later.
	for rows.Next() {

		rows.Scan(&comment.ID, &comment.Content, &comment.Likes, &comment.Time, &comment.Postid, &comment.Username)
		comments = append(comments, comment)
	}
	return comments, nil
}

func addComment(content string, postid int) error {
	//Open database
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`
		INSERT INTO comments (content, likes, time, post_id, username) values (?, ?, ?, ?, ?);
	`)
	if err != nil {
		return err
	}

	likes := 0 //inital likes is 0
	time := time.Now().Format("2006-01-02 15:04")

	//Exec adds data to database
	_, errPost := stmt.Exec(content, likes, time, postid, user.Username)
	if errPost != nil {
		return err
	}
	fmt.Println("Commented")

	return nil
}
