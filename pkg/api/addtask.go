package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"todo_task_final/pkg/db"
)

func writeJSON(w http.ResponseWriter, data any) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func writeErrorJSON(w http.ResponseWriter, msg string) {
	writeJSON(w, map[string]string{"error": msg})
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeErrorJSON(w, "Ошибка сериализации данных")
		return
	}

	if task.Title == "" {
		writeErrorJSON(w, "Заголовок задачи не указан")
		return
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(Dateformat)
	}

	t, err := time.Parse(Dateformat, task.Date)
	if err != nil {
		writeErrorJSON(w, "Некорректная дата")
		return
	}

	var next string

	if task.Repeat != "" {
		next, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeErrorJSON(w, err.Error())
			return
		}
	}

	nowDate, _ := time.Parse(Dateformat, now.Format(Dateformat))
	if t.Before(nowDate) {
		if task.Repeat == "" {
			task.Date = now.Format(Dateformat)
		} else {
			task.Date = next
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeErrorJSON(w, err.Error())
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}
