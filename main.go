package main

import (
	"database/sql"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/basicauth"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

var database *sql.DB

func main() {

	END_POINT := API_HOST + ":" + API_PORT

	Cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
	})

	authConfig := basicauth.Config{
		Users:   map[string]string{USER_NAME: PASSWORD},
		Realm:   "Authorization Required",
		Expires: time.Duration(30) * time.Minute,
	}

	authentication := basicauth.New(authConfig)
	app := iris.New()
	app.Use(Cors)

	/*
		|For production use|

		app.RegisterView(iris.HTML("./public", ".html"))
		app.StaticWeb("/static/", "./static/")
		app.StaticWeb("/", ".")
		app.Get("/", func(ctx iris.Context) {
			ctx.View("index.html")
		})
	*/
	var err error
	database, err = sql.Open("sqlite3", "./express-blog.db")
	if err != nil {
		log.Fatal("Couldn't create express-blog.db")
	}

	defer database.Close()

	query := `CREATE TABLE IF NOT EXISTS blog_posts (
				PostId INTEGER PRIMARY KEY AUTOINCREMENT,
				Title VARCHAR(45),
				Content VARCHAR(45),
				Date TEXT
			)`

	_, err = database.Exec(query)
	if err != nil {
		log.Fatalf("Failed to execute create table query %s \nError:%q", query, err)
	}
	app.Post("/v1/insert", authentication, InsertConentHandler)
	app.Get("/v1/getposts", GetPostsHandler)
	app.Post("/v1/login", UserLogin)
	app.Run(iris.Addr(END_POINT))
}
