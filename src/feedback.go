package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Suggestions
type Suggestion struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Message     string `json:"message"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
}

/* Feedback */
func getFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id := r.URL.Query().Get("id")

	if len(id) != 0 {
		fmt.Println("String is not empty")

		// If ID then fetch just that ID and return it
		var queryString = "SELECT * FROM suggestions WHERE ID='" + id + "'"
		row, err := db.Query(queryString)
		if err != nil {
			panic(err.Error())
		}
		defer row.Close()

		var suggestion Suggestion
		for row.Next() {
			err := row.Scan(&suggestion.ID, &suggestion.Username, &suggestion.DisplayName, &suggestion.Message, &suggestion.Type, &suggestion.Status, &suggestion.Priority)
			if err != nil {
				panic(err.Error())
			}
		}

		fmt.Println(suggestion)
		json.NewEncoder(w).Encode(suggestion)
		return
	}

	var suggestions []Suggestion

	rows, err := db.Query("SELECT * FROM suggestions")
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var suggestion Suggestion
		err := rows.Scan(&suggestion.ID, &suggestion.Username, &suggestion.DisplayName, &suggestion.Message, &suggestion.Type, &suggestion.Status, &suggestion.Priority)
		if err != nil {
			fmt.Println("Error scanning: " + err.Error())
		}
		suggestions = append(suggestions, suggestion)
	}

	json.NewEncoder(w).Encode(suggestions)
}
