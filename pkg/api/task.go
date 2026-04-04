package api

import (
	"net/http"
	"todo_task_final/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodGet:
		GetTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)

	default:
		writeErrorJSON(w, "Метод не поддерживается")
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeErrorJSON(w, "ID задачи не указан")
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeErrorJSON(w, err.Error())
		return
	}

	writeJSON(w, map[string]any{})

}
