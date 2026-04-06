package main

import (
	"log"
	"net/http"
	"os"
	"todo_task_final/pkg/api"
	"todo_task_final/pkg/db"
)

func main() {

	dbfile := os.Getenv("TODO_DBFILE")
	if dbfile == "" {
		dbfile = "scheduler.db"
	}

	if err := db.Init(dbfile); err != nil {
		panic(err)
	}
	defer db.DB.Close()

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	log.Printf("Сервер запущен на порту %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
