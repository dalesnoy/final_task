package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"todo_task_final/pkg/db"
)

func writeJSON(w http.ResponseWriter, data any) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func writeErrorJSON(w http.ResponseWriter, msg string, status int) {
	log.Printf("ERROR [%d]: %s", status, msg)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeErrorJSON(w, "Ошибка сериализации данных", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeErrorJSON(w, "Заголовок задачи не указан", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(Dateformat)
	}

	t, err := time.Parse(Dateformat, task.Date)
	if err != nil {
		writeErrorJSON(w, "Некорректная дата", http.StatusBadRequest)
		return
	}

	var next string

	if task.Repeat != "" {
		next, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeErrorJSON(w, err.Error(), http.StatusBadRequest)
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
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}
