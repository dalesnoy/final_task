package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const Dateformat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	date, err := time.Parse(Dateformat, dstart)
	if err != nil {
		return "", errors.New("Некорректная дата")
	}

	if repeat == "" {
		return "", errors.New("Некорректная дата")
	}
	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("Некорректный repeat")
		}

		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > 400 {
			return "", errors.New("Некорректное интервал дней")
		}

		for {
			date = date.AddDate(0, 0, n)
			if date.After(now) {
				break
			}
		}
		return date.Format(Dateformat), nil
	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}

		}
		return date.Format(Dateformat), nil
	default:
		return "", errors.New("Некорректный repeat")
	}
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		if now, err = time.Parse(Dateformat, nowStr); err != nil {
			http.Error(w, "Некорректный now", http.StatusBadRequest)
			return
		}
	}

	result, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(result))
}
