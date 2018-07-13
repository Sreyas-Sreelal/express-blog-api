package main

import (
	"fmt"
	"github.com/kataras/iris"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

//PostData struct Information about new blog post
type PostData struct {
	Name    string `json:"Name"`
	Title   string `db:"Title" json:"Title"`
	Content string `db:"Content" json:"Content"`
	Date    string `db:"Date" json:"Date"`
}

//InsertConentHandler handler to inserts new blog content to database
func InsertConentHandler(ctx iris.Context) {
	Content := ctx.PostValue("Content")
	Title := ctx.PostValue("Title")
	Date := ctx.PostValue("Date")
	log.Printf("Title : %s content : %s date :%s", Title, Content, Date)

	Query := fmt.Sprintf(`
		INSERT INTO
		BLOG_POSTS(Title,Content,Date)
		VALUES('%s','%s','%s')
		`,
		Title,
		Content,
		Date)

	_, err := database.Exec(Query)

	if err != nil {
		log.Fatalf("Couldn't insert\n Error:%q", err)
		return
	} else {
		ctx.StatusCode(iris.StatusOK)
	}
}

//GetPostsHandler handler for fetching blog posts from database
func GetPostsHandler(ctx iris.Context) {
	Posts := []PostData{}
	Query := "SELECT Title,Content,Date FROM BLOG_POSTS ORDER BY PostId DESC"
	rows, err := database.Query(Query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("Failed fetching blog posts from database \nError: %q", err)
		return
	}
	var count = 0

	for rows.Next() {
		content := PostData{}
		err = rows.Scan(&content.Title, &content.Content, &content.Date)
		if err != nil {
			log.Fatal(err)
		}
		content.Name = USER_NAME
		Posts = append(Posts, content)
		count++
	}

	if count == 0 {
		ctx.StatusCode(iris.StatusNotFound)
	} else {
		ctx.JSON(&Posts)
	}
}

//UserLogin handles user' login
func UserLogin(ctx iris.Context) {
	PostName := ctx.PostValue("Name")
	PostPass := ctx.PostValue("Password")
	log.Printf("Username :%s password :%s", PostName, PostPass)
	if PostName != USER_NAME || PostPass != PASSWORD {
		ctx.StatusCode(iris.StatusForbidden)
	} else {
		ctx.StatusCode(iris.StatusOK)
	}
}

//UserLogout handles user' logout
func UserLogout(ctx iris.Context) {
	log.Printf("loggedout")
	//todo
}
