package main

import (
	"log"
	"net/http"

	"forim/database"
	"forim/database/creatdatabase"
	hand "forim/handlers"
)

func main() {
	err := database.InitializeDB("./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()
	creatdatabase.Creatdb()

	srv := http.Server{
		Addr:    ":8080",
		Handler: routes(),
	}
	log.Println("Listening on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.HandleFunc("/", hand.Login)
	mux.HandleFunc("/register", hand.Register)
	mux.HandleFunc("/post", hand.GetHome)
	mux.HandleFunc("/comment", hand.GetComment)
	mux.HandleFunc("/post/create", hand.CreatePost)
	mux.HandleFunc("/like_post", hand.Like_post)

	mux.HandleFunc("/newcomment", hand.NewComment)
	mux.HandleFunc("/logout", database.Logout)

	return mux
}
