package main

type PageData struct {
	Username        string
	Message         string
	Categories      []Category
	Category        string
	CategoryEncoded string
	Posts           []Post
	Post            Post
	Comments        []Comment
	Google          string
	Github          string
}

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

type Category struct {
	ID           int
	Title        string
	Description  string
	TitleEncoded string
}

type Post struct {
	ID       int
	Title    string
	Content  string
	Category string
	Likes    int
	Time     string
	Username string
	Picture  string
}
type Comment struct {
	ID       int
	Content  string
	Likes    int
	Time     string
	Postid   int
	Username string
}

var user User
