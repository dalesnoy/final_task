package main

import (
	"net/http"
	"os"
	"todo_task_final/pkg/db"
)

func main() {
	webDir := "./web"

	dbfile := os.Getenv("TODO_DBFILE")
	if dbfile == "" {
		dbfile = "scheduler.db"
	}

	if err := db.Init(dbfile); err != nil {
		panic(err)
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}

}
