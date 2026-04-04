package api

import (
	"net/http"
	"time"
	"todo_task_final/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeErrorJSON(w, "ID задачи не указан")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeErrorJSON(w, err.Error())
		return
	}

	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			writeErrorJSON(w, err.Error())
			return
		}
		writeJSON(w, map[string]any{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeErrorJSON(w, err.Error())
		return
	}

	err = db.UpdateDate(id, next)
	if err != nil {
		writeErrorJSON(w, err.Error())
		return
	}

	writeJSON(w, map[string]any{})

}
