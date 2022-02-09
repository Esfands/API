package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Suggestions
type Suggestion struct {
	ID          int    `json:"ID"`
	Username    string `json:"Username"`
	DisplayName string `json:"DisplayName"`
	Message     string `json:"Message"`
	Type        string `json:"Type"`
	Status      string `json:"Status"`
}

/* Feedback */
func getFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var suggestions []Suggestion

	rows, err := db.Query("SELECT * FROM suggestions")
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var suggestion Suggestion
		err := rows.Scan(&suggestion.ID, &suggestion.Username, &suggestion.DisplayName, &suggestion.Message, &suggestion.Type, &suggestion.Status)
		if err != nil {
			fmt.Println("Error scanning: " + err.Error())
		}
		suggestions = append(suggestions, suggestion)
	}

	json.NewEncoder(w).Encode(suggestions)
}

func getFeedbackById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	key := vars["id"]

	var queryString = "SELECT * FROM suggestions WHERE ID='" + key + "'"
	row, err := db.Query(queryString)
	if err != nil {
		panic(err.Error())
	}
	defer row.Close()

	var suggestion Suggestion
	for row.Next() {
		err := row.Scan(&suggestion.ID, &suggestion.Username, &suggestion.DisplayName, &suggestion.Message, &suggestion.Type, &suggestion.Status)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(suggestion)
}
