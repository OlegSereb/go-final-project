// pkg/api/nextdate.go
package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is empty: task will be deleted")
	}

	startDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid dstart format: %w", err)
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("invalid repeat format")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid 'd' format: expected 'd <number>'")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("invalid number in 'd' rule")
		}
		if days <= 0 || days > 400 {
			return "", errors.New("number of days must be between 1 and 400")
		}

		date := startDate
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid 'y' format: no extra arguments allowed")
		}

		date := startDate
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}

	default:
		return "", errors.New("unsupported repeat rule: only 'd N' and 'y' are supported in basic mode")
	}
}

// NextDateHandler обрабатывает GET /api/nextdate
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	if dateStr == "" {
		http.Error(w, "missing 'date' parameter", http.StatusBadRequest)
		return
	}
	if repeatStr == "" {
		http.Error(w, "missing 'repeat' parameter", http.StatusBadRequest)
		return
	}

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "invalid 'now' format, expected YYYYMMDD", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(next))
}
