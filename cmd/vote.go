package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

func setCommentVote(commentId int) {

	var likecount int
	var dislikecount int

	db, err := sql.Open("sqlite3", "./database.db")

	if err != nil {
		fmt.Println("getCommentVote sql open error")
		log.Panic(err)
	}
	//Count likes
	rows, err := db.Query(`SELECT COUNT (*) FROM commentvotes WHERE type = "like" AND post_id = ?`, commentId)

	if err != nil {
		fmt.Println("getCommentVote dbQuery error")
		log.Panic(err)
	}

	for rows.Next() {
		rows.Scan(&likecount)
	}

	// Count dislikes
	rows, err = db.Query(`SELECT COUNT (*) FROM commentvotes WHERE type = "dislike" AND post_id = ?`, commentId)

	if err != nil {
		fmt.Println("getCommentVote dbQuery dislike error")
		log.Panic(err)
	}

	for rows.Next() {
		rows.Scan(&dislikecount)
	}
	total := likecount - dislikecount

	//Update likes int in the comment db
	stmt, err := db.Prepare(`UPDATE comments SET likes = ? WHERE ID = ?`)
	if err != nil {
		log.Fatal(err)
	}

	//Exec adds data to database
	_, errReg := stmt.Exec(total, commentId)
	if errReg != nil {
		fmt.Println(errReg)
	}
	db.Close()
}

func commentVote(like, dislike, commentIdString string) {

	userid := user.Username
	//Creates unique string, so user cannot like/dislike a post more than once
	uniqId := userid + commentIdString
	time := time.Now().Format("2006-01-02 15:04")

	commentid, err := strconv.Atoi(commentIdString)
	if err != nil {
		log.Panic(err)
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println("commentVote db open error")
		log.Panic(err)
	}
	//If user presses like button
	if like != "" {

		//Checks if user has already disliked the post, if so, deletes the dislike
		stmt, err := db.Prepare(`DELETE FROM commentvotes WHERE (type, uniq_id) = (?, ?)`)

		if err != nil {
			fmt.Println("DELETE dislike error")
			log.Panic(err)
		}

		stmt.Exec("dislike", uniqId)

		//Inserts like data to the DB
		stmt, err = db.Prepare(`INSERT INTO commentvotes (type, post_id, username, time, uniq_id) values (?, ?, ?, ?, ?)`)

		if err != nil {
			fmt.Println("commentVote dbPrepare error")
			log.Panic(err)
		}

		_, errVote := stmt.Exec(like, commentid, userid, time, uniqId)

		//If user has already liked a comment, it will not add new like to the DB, instead will give this error and delete the like.
		if errVote != nil {

			//If already liked, delete like.
			if errVote.Error() == "UNIQUE constraint failed: commentvotes.uniq_id" {
				stmt, err := db.Prepare(`DELETE FROM commentvotes WHERE uniq_id = ?`)
				if err != nil {
					fmt.Println("commentVote DELETE like error")
					log.Panic(err)
				}
				stmt.Exec(uniqId)
				setCommentVote(commentid)
				fmt.Println("Comment like deleted")
			}

		} else {
			setCommentVote(commentid)
			fmt.Println("Comment liked")
		}
	}

	// IF user presses dislike button
	if dislike != "" {

		//Checks if user has already disliked the post, if so, deletes the dislike
		stmt, err := db.Prepare(`DELETE FROM commentvotes WHERE (type, uniq_id) = (?, ?)`)
		if err != nil {
			fmt.Println("DELETE dislike error")
			log.Panic(err)
		}

		stmt.Exec("like", uniqId)

		stmt, err = db.Prepare(`INSERT INTO commentvotes (type, post_id, username, time, uniq_id) values (?, ?, ?, ?, ?)`)

		if err != nil {
			fmt.Println("postVote dbPrepare error")
			log.Panic(err)
		}
		_, errVote := stmt.Exec(dislike, commentid, userid, time, uniqId)

		//If user has already disliked a post, it will delete it.
		if errVote != nil {

			if errVote.Error() == "UNIQUE constraint failed: commentvotes.uniq_id" {
				stmt, err := db.Prepare(`DELETE FROM commentvotes WHERE uniq_id = ?`)

				if err != nil {
					fmt.Println("postVote DELETE dislike")
					log.Panic(err)
				}
				stmt.Exec(uniqId)
				setCommentVote(commentid)
				fmt.Println("Comment dislike deleted")
			}

		} else {
			setCommentVote(commentid)
			fmt.Println("Comment disliked")
		}
	}
}

func postVote(like, dislike, postIdString string) {

	username := user.Username
	//Creates unique string, so user cannot like/dislike a post more than once
	uniqId := username + postIdString
	time := time.Now().Format("2006-01-02 15:04")
	postid, err := strconv.Atoi(postIdString)
	if err != nil {
		fmt.Println("postvote atoi")
		log.Panic(err)
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println("postVote db open error")
		log.Panic(err)
	}

	//If user presses like button
	if like != "" {

		//Checks if user has already disliked the post, if so, deletes the dislike
		stmt, err := db.Prepare(`DELETE FROM postvotes WHERE (type, uniq_id) = (?, ?)`)

		if err != nil {
			fmt.Println("DELETE dislike error")
			log.Panic(err)
		}
		stmt.Exec("dislike", uniqId)

		//Inserts like data to the DB
		stmt, err = db.Prepare(`INSERT INTO postvotes (type, post_id, username, time, uniq_id) values (?, ?, ?, ?, ?)`)

		if err != nil {
			fmt.Println("postVote dbPrepare error")
			log.Panic(err)
		}
		_, errVote := stmt.Exec(like, postid, username, time, uniqId)

		//If user has already liked a post, it will not add new like to the DB, instead will give this error and delete the like.
		if errVote != nil {

			//If already liked, delete like.
			if errVote.Error() == "UNIQUE constraint failed: postvotes.uniq_id" {
				stmt, err := db.Prepare(`DELETE FROM postvotes WHERE uniq_id = ?`)
				if err != nil {
					fmt.Println("postVote DELETE like")
					log.Panic(err)
				}
				stmt.Exec(uniqId)
				setPostVote(postid)
				fmt.Println("Post like deleted")
			}

		} else {
			setPostVote(postid)
			fmt.Println("Post liked")
		}
	}

	// IF user presses dislike button
	if dislike != "" {
		//Checks if user has already disliked the post, if so, deletes the dislike
		stmt, err := db.Prepare(`DELETE FROM postvotes WHERE (type, uniq_id) = (?, ?)`)
		if err != nil {
			fmt.Println("DELETE dislike error")
			log.Panic(err)
		}
		stmt.Exec("like", uniqId)

		stmt, err = db.Prepare(`INSERT INTO postvotes (type, post_id, username, time, uniq_id) values (?, ?, ?, ?, ?)`)

		if err != nil {
			fmt.Println("postVote dbPrepare error")
			log.Panic(err)
		}
		_, errVote := stmt.Exec(dislike, postid, username, time, uniqId)

		//If user has already disliked a post, it will delete it.
		if errVote != nil {

			if errVote.Error() == "UNIQUE constraint failed: postvotes.uniq_id" {
				stmt, err := db.Prepare(`DELETE FROM postvotes WHERE uniq_id = ?`)

				if err != nil {
					fmt.Println("postVote DELETE dislike")
					log.Panic(err)
				}
				stmt.Exec(uniqId)
				//setPostVote (Likes - dislikes = Nr of votes)
				setPostVote(postid)
				fmt.Println("Post dislike deleted")
			}

		} else {
			//setPostVote (Likes - dislikes = Nr of votes)
			setPostVote(postid)
			fmt.Println("Post disliked")
		}
	}
}

//Get sum of likes and dislikes and update the likes of the current post
func setPostVote(postid int) {

	var likecount int
	var dislikecount int
	var total int
	db, err := sql.Open("sqlite3", "./database.db")

	if err != nil {
		fmt.Println("getPostVote sql open error")
		log.Panic(err)
	}
	//Count likes
	rows, err := db.Query(`SELECT COUNT (*) FROM postvotes WHERE type = "like" AND post_id = ?`, postid)

	if err != nil {
		fmt.Println("getPostVote dbQuery error")
		log.Panic(err)
	}

	for rows.Next() {
		rows.Scan(&likecount)
	}

	// Count dislikes
	rows, err = db.Query(`SELECT COUNT (*) FROM postvotes WHERE type = "dislike" AND post_id = ?`, postid)

	if err != nil {
		fmt.Println("getPostVote dbQuery dislike error")
		log.Panic(err)
	}

	for rows.Next() {
		rows.Scan(&dislikecount)
	}
	total = likecount - dislikecount

	stmt, err := db.Prepare(`UPDATE posts SET likes = ? WHERE ID = ?`)
	if err != nil {
		log.Fatal(err)
	}

	//Exec adds data to database
	_, errReg := stmt.Exec(total, postid)
	if errReg != nil {
		fmt.Println(errReg)
	}
	db.Close()
}
