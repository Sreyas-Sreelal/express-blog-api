package main

import (
	"database/sql"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var database *sql.DB

func main() {

	END_POINT := API_HOST + ":" + API_PORT

	Cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	app := iris.New()
	app.Use(Cors)
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
	app.Post("/insert", InsertConentHandler)
	app.Get("/getposts", GetPostsHandler)
	app.Post("/login", UserLogin)
	app.Get("/logout", UserLogout)
	app.Run(iris.Addr(END_POINT))
}
