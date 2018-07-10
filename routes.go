package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var (
	cookieUser = "express-user-cookie"
	sess       = sessions.New(sessions.Config{Cookie: cookieUser})
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
	if auth, _ := sess.Start(ctx).GetBoolean("loggined"); !auth {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}
	Query := fmt.Sprintf(`
		INSERT INTO
		BLOG_POSTS(Title,Content,Date)
		VALUES('%s','%s','%s')
		`,
		ctx.PostValue("Title"),
		ctx.PostValue("Content"),
		ctx.PostValue("Date"))

	_, err := database.Exec(Query)

	if err != nil {
		log.Fatalf("Couldn't insert\n Error:%q", err)
		return
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
	session := sess.Start(ctx)
	PostName := ctx.PostValue("Name")
	PostPass := ctx.PostValue("Password")
	log.Printf("Username :%s password :%s", PostName, PostPass)
	if PostName != USER_NAME || PostPass != PASSWORD {
		ctx.StatusCode(iris.StatusForbidden)
	} else {
		ctx.StatusCode(iris.StatusOK)
		session.Set("loggined", true)
	}

}

//UserLogout handles user' logout
func UserLogout(ctx iris.Context) {
	session := sess.Start(ctx)
	session.Set("loggined", false)
}
