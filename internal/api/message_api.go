package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type display_state struct {
	mu          sync.RWMutex
	displayText string
}

var state display_state

func StartApiServer() func() {
	http.HandleFunc("/display", handleDisplay)

	return func() {
		fmt.Println("API server listening on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}
}

func GetDisplayState() []string {
	state.mu.RLock()
	defer state.mu.RUnlock()

	if state.displayText != "" {
		return []string{"", state.displayText}
	}

	return []string{}
}

func handleDisplay(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")

	state.mu.Lock()
	state.displayText = text
	state.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "display updated"})
}
