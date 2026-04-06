package api

import (
	"errors"
	"log"
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
			return "", errors.New("Некорректный интервал дней")
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

	case "w":
		if len(parts) != 2 {
			return "", errors.New("не указаны дни недели")
		}

		weekdays := make(map[time.Weekday]bool)
		for _, ds := range strings.Split(parts[1], ",") {
			d, err := strconv.Atoi(strings.TrimSpace(ds))
			if err != nil || d < 1 || d > 7 {
				return "", errors.New("некорректный день недели")
			}
			weekdays[time.Weekday(d%7)] = true
		}

		if date.Before(now) {
			date = now
		}

		for {
			date = date.AddDate(0, 0, 1)
			if weekdays[date.Weekday()] {
				break
			}
		}
		return date.Format(Dateformat), nil

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("некорректный формат m")
		}

		var days []int
		for _, ds := range strings.Split(parts[1], ",") {
			d, err := strconv.Atoi(strings.TrimSpace(ds))
			if err != nil || d == 0 || d < -2 || d > 31 {
				return "", errors.New("некорректный день месяца")
			}
			days = append(days, d)
		}

		months := make(map[int]bool)
		if len(parts) == 3 {
			for _, ms := range strings.Split(parts[2], ",") {
				m, err := strconv.Atoi(strings.TrimSpace(ms))
				if err != nil || m < 1 || m > 12 {
					return "", errors.New("некорректный месяц")
				}
				months[m] = true
			}
		}

		if date.Before(now) {
			date = now
		}

		for i := 0; i < 400*31; i++ {
			date = date.AddDate(0, 0, 1)

			if len(months) > 0 && !months[int(date.Month())] {
				continue
			}

			for _, d := range days {
				targetDay := d
				if d == -1 {
					targetDay = lastDay(date)
				} else if d == -2 {
					targetDay = lastDay(date) - 1
				}

				if date.Day() == targetDay {
					return date.Format(Dateformat), nil
				}
			}
		}
		return "", errors.New("не удалось найти подходящую дату")

	default:
		return "", errors.New("Некорректный repeat")
	}
}

func lastDay(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
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

	if _, err := w.Write([]byte(result)); err != nil {
		log.Printf("ERROR: ошибка записи ответа: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
