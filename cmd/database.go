package main

import (
	"database/sql"
	"fmt"
	"log"
)

//Create database
func createDB() {
	createAccountDB()
	createCommentDB()
	createPostDB()
	createCategoryDB()
	createSesssionDB()
	createPostVoteDB()
	createCommentVoteDB()
	createNotificationDB()
}

func createNotificationDB() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Panic(err)
	}
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS notifications (
		"notification_id"	INTEGER NOT NULL UNIQUE,
		"type" TEXT NOT NULL,
		"time" TEXT NOT NULL,
		"username_a" TEXT NOT NULL,
		"username_b" TEXT NOT NULL,
		"post_id" INTEGER NOT NULL,
		PRIMARY KEY("notification_id" AUTOINCREMENT)
		FOREIGN KEY (username_a) REFERENCES accounts (username)
		FOREIGN KEY (username_b) REFERENCES accounts (username)
		FOREIGN KEY (post_id) REFERENCES posts (ID)
	);
	`)

	if err != nil {
		fmt.Println(err.Error())
		log.Panic(err)
	}
	stmt.Exec()
}

func createAccountDB() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Panic(err)
	}
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS accounts (
		"ID"	INTEGER NOT NULL UNIQUE,
		"username"	TEXT NOT NULL UNIQUE,
		"email"	TEXT NOT NULL UNIQUE,
		"password"	TEXT,
		PRIMARY KEY("ID" AUTOINCREMENT)
	);
	`)

	if err != nil {
		log.Panic(err)
	}
	stmt.Exec()
}

func createPostDB() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Panic(err)
	}
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS posts (
		"ID"	INTEGER NOT NULL UNIQUE,
		"title"	TEXT NOT NULL,
		"content"	TEXT,
		"category"	TEXT NOT NULL,
		"likes" INTEGER,
		"time"	TEXT DEFAULT CURRENT_TIMESTAMP,
		"username" TEXT NOT NULL,
		"picture" TEXT,
		PRIMARY KEY("ID" AUTOINCREMENT),
		FOREIGN KEY (username) REFERENCES accounts (username),
		FOREIGN KEY (category) REFERENCES categories (title)
	);
	`)

	if err != nil {
		log.Panic(err)
	}
	stmt.Exec()
}

func createCategoryDB() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Panic(err)
	}
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS categories (
		"ID"	INTEGER NOT NULL UNIQUE,
		"title"	TEXT UNIQUE,
		"title_encoded" TEXT UNIQUE,
		"description" TEXT,
		PRIMARY KEY("ID" AUTOINCREMENT)
	);
	`)

	if err != nil {
		log.Panic(err)
	}
	stmt.Exec()
}

func createCommentDB() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Panic(err)
	}
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS comments (
		"ID"	INTEGER NOT NULL UNIQUE,
		"content"	TEXT,
		"likes" INTEGER NOT NULL,
		"time"	TEXT DEFAULT CURRENT_TIMESTAMP,
		"post_id" INTEGER NOT NULL,
		"username" TEXT NOT NULL,
		PRIMARY KEY("ID" AUTOINCREMENT),
		FOREIGN KEY (username) REFERENCES accounts (username),
		FOREIGN KEY (post_id) REFERENCES posts (ID)
	);
	`)

	if err != nil {
		log.Panic(err)
	}
	stmt.Exec()
}

func createSesssionDB() {

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Panic(err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS session (
		"ID"	INTEGER NOT NULL UNIQUE,
		"time"	DATETIME NOT NULL,
		"user_id"	INTEGER NOT NULL,
		"uuid"	TEXT NOT NULL,
		PRIMARY KEY ("ID" AUTOINCREMENT),
		FOREIGN KEY (user_id) REFERENCES accounts (ID)
	);`)

	if err != nil {
		log.Panic(err)
	}
	stmt.Exec()
}

func createPostVoteDB() {

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Panic(err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS postvotes (
		"postvote_id"	INTEGER NOT NULL UNIQUE,
		"type" STRING NOT NULL,
		"post_id" INTEGER NOT NULL,
		"username"	INTEGER NOT NULL,
		"time"	DATETIME NOT NULL,
		"uniq_id" STRING NOT NULL UNIQUE,
		PRIMARY KEY("postvote_id" AUTOINCREMENT),
		FOREIGN KEY(username) REFERENCES accounts (username)
		FOREIGN KEY(post_id) REFERENCES posts (ID)
	);`)

	if err != nil {
		log.Panic(err)
	}
	stmt.Exec()
}

func createCommentVoteDB() {

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Panic(err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS commentvotes (
		"commentvote_id"	INTEGER NOT NULL UNIQUE,
		"type" STRING NOT NULL,
		"post_id" INTEGER NOT NULL,
		"username"	INTEGER NOT NULL,
		"time"	DATETIME NOT NULL,
		"uniq_id" STRING NOT NULL UNIQUE,
		PRIMARY KEY("commentvote_id" AUTOINCREMENT),
		FOREIGN KEY (username) REFERENCES accounts (username)
		FOREIGN KEY (post_id) REFERENCES posts (ID)
	);`)

	if err != nil {
		log.Panic(err)
	}
	stmt.Exec()
}
