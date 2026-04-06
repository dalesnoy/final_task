package api

import (
	"encoding/json"
	"net/http"
	"time"
	"todo_task_final/pkg/db"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeErrorJSON(w, "ID задачи не указан", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeErrorJSON(w, "задача не найдена", http.StatusNotFound)
		return
	}
	writeJSON(w, task)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		writeErrorJSON(w, "Ошибка десериализации данных", http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeErrorJSON(w, "ID задачи не указан", http.StatusBadRequest)
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
		writeErrorJSON(w, "Ошибка парсинга даты", http.StatusBadRequest)
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
	err = db.UpdateTask(&task)
	if err != nil {
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{})
}
