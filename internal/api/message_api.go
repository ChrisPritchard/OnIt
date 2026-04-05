package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"
)

type display_state struct {
	mu          sync.RWMutex
	alarmTime   *time.Time
	alarmTitle  string
	alarmDesc   string
	alarmAlert  bool
	displayText string
}

var state display_state

func StartApiServer() func() {
	http.HandleFunc("/alarm", handleAlarm)
	http.HandleFunc("/display", handleDisplay)
	http.HandleFunc("/clear", handleClear)

	return func() {
		fmt.Println("API server listening on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}
}

func GetDisplayState(currentTime time.Time, tick bool) []string {
	state.mu.RLock()
	defer state.mu.RUnlock()

	message := []string{}
	if state.displayText != "" {
		message = []string{"", state.displayText}
	}

	if state.alarmTime != nil {
		currentUnix := currentTime.Unix()
		alarmUnix := state.alarmTime.Unix()
		if currentUnix < alarmUnix+60 {
			timeSet := fmt.Sprintf("alarm set at %s:", state.alarmTime.Local().Format("03:04 PM"))
			text := []string{timeSet, " " + state.alarmTitle, " " + state.alarmDesc}
			if currentUnix >= alarmUnix && tick {
				text[1] = "*\x1b[7m" + text[1]
				text[2] = "*" + text[2] + "\x1b[0m"
			}
			message = slices.Concat(message, []string{""}, text)
		}
	}

	return message
}

func handleAlarm(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	timeStr := query.Get("time")
	title := query.Get("title")
	description := query.Get("description")
	alert := query.Get("alert") == "true"

	if timeStr == "" {
		http.Error(w, "missing time parameter", http.StatusBadRequest)
		return
	}

	timestampSec, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid time format", http.StatusBadRequest)
		return
	}
	timestamp := time.Unix(timestampSec, 0)

	state.mu.Lock()
	state.alarmTime = &timestamp
	state.alarmTitle = title
	state.alarmDesc = description
	state.alarmAlert = alert
	state.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "alarm set"})
}

func handleDisplay(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")

	state.mu.Lock()
	state.displayText = text
	state.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "display updated"})
}

func handleClear(w http.ResponseWriter, r *http.Request) {
	state.mu.Lock()
	state.alarmTime = nil
	state.alarmTitle = ""
	state.alarmDesc = ""
	state.alarmAlert = false
	state.displayText = ""
	state.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "cleared"})
}
