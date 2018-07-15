package main

import (
	"github.com/kataras/iris"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

//PostData struct Information about new blog post
type PostData struct {
	PostID  string `json:"PostID"`
	Name    string `json:"Name"`
	Title   string `db:"Title" json:"Title"`
	Content string `db:"Content" json:"Content"`
	Date    string `db:"Date" json:"Date"`
}

//InsertConentHandler handler to inserts new blog content to database
func InsertConentHandler(ctx iris.Context) {
	username, password, success := ctx.Request().BasicAuth()
	log.Printf("Username : %s Password : %s success : %t", username, password, success)
	Content := ctx.PostValue("Content")
	Title := ctx.PostValue("Title")
	Date := ctx.PostValue("Date")
	log.Printf("Title : %s content : %s date :%s", Title, Content, Date)

	statement, _ := database.Prepare("INSERT INTO	BLOG_POSTS(Title,Content,Date) VALUES(?,?,?)")
	_, err := statement.Exec(Title, Content, Date)

	if err != nil {
		log.Fatalf("Couldn't insert\n Error:%q", err)
		return
	}
	ctx.StatusCode(iris.StatusOK)
}

//GetPostsHandler handler for fetching blog posts from database
func GetPostsHandler(ctx iris.Context) {
	Posts := []PostData{}
	Query := "SELECT PostID,Title,Content,Date FROM BLOG_POSTS ORDER BY PostID DESC"
	rows, err := database.Query(Query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("Failed fetching blog posts from database \nError: %q", err)
		return
	}
	var count = 0

	for rows.Next() {
		content := PostData{}
		err = rows.Scan(&content.PostID, &content.Title, &content.Content, &content.Date)
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

//UserLoginHandler handles user' login
func UserLoginHandler(ctx iris.Context) {
	PostName := ctx.PostValue("Name")
	PostPass := ctx.PostValue("Password")
	log.Printf("Username :%s password :%s", PostName, PostPass)
	if PostName != USER_NAME || PostPass != PASSWORD {
		ctx.StatusCode(iris.StatusForbidden)
	} else {
		ctx.StatusCode(iris.StatusOK)
	}
}

//DeleteContentHandler handles post delete
func DeleteContentHandler(ctx iris.Context) {
	username, password, success := ctx.Request().BasicAuth()
	log.Printf("Username : %s Password : %s success : %t", username, password, success)
	PostID := ctx.PostValue("PostID")

	log.Printf("Id : %s", PostID)

	statement, _ := database.Prepare("DELETE FROM BLOG_POSTS WHERE PostID=?")
	_, err := statement.Exec(PostID)

	if err != nil {
		log.Fatalf("Couldn't insert\n Error:%q", err)
		return
	}
	ctx.StatusCode(iris.StatusOK)
}

//GetPostsByIDHandler handles getrequest for post with specified id
func GetPostsByIDHandler(ctx iris.Context) {
	PostID := ctx.Params().Get("postid")
	log.Printf("Id is %s ", PostID)
	statement, err := database.Prepare("SELECT Content,Title FROM BLOG_POSTS WHERE PostID=?")

	defer statement.Close()
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	content := PostData{}

	err = statement.QueryRow(PostID).Scan(&content.Content, &content.Title)
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(&content)

}
