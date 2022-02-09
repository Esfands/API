package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Commands
type Command struct {
	Name               string `json:"Name"`
	Aliases            string `json:"Aliases"`
	Permissions        string `json:"Permissions"`
	GlobalCooldown     int    `json:"GlobalCooldown"`
	Cooldown           int    `json:"Cooldown"`
	Description        string `json:"Description"`
	DynamicDescription string `json:"DynamicDescription"`
	Testing            bool   `json:"Testing"`
	OfflineOnly        bool   `json:"OfflineOnly"`
	OnlineOnly         bool   `json:"OnlineOnly"`
	Count              int    `json:"Count"`
}

// OTF Commands
type OTFCommand struct {
	Name     string `json:"Name"`
	Response string `json:"Response"`
	Creator  string `json:"Creator"`
	Count    int    `json:"Count"`
}

/* Commands */
func getAllCommands(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var commands []Command

	rows, err := db.Query("SELECT * FROM commands")
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {

		var cmd Command
		err := rows.Scan(&cmd.Name, &cmd.Aliases, &cmd.Permissions, &cmd.Description, &cmd.DynamicDescription, &cmd.GlobalCooldown, &cmd.Cooldown, &cmd.Testing, &cmd.OfflineOnly, &cmd.OnlineOnly, &cmd.Count)
		if err != nil {
			fmt.Println("error scanning: " + err.Error())
		}
		commands = append(commands, cmd)
	}
	json.NewEncoder(w).Encode(commands)
}

func getCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	key := vars["name"]

	var queryString = "SELECT * FROM commands WHERE Name='" + key + "'"
	fmt.Println(queryString)
	row, err := db.Query(queryString)
	if err != nil {
		panic(err.Error())
	}
	defer row.Close()

	var cmd Command
	for row.Next() {
		err := row.Scan(&cmd.Name, &cmd.Aliases, &cmd.Permissions, &cmd.Description, &cmd.DynamicDescription, &cmd.GlobalCooldown, &cmd.Cooldown, &cmd.Testing, &cmd.OfflineOnly, &cmd.OnlineOnly, &cmd.Count)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(cmd)
}

/* OTF Commands */
func getAllOTFCommands(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var otfCommands []OTFCommand

	rows, err := db.Query("SELECT * FROM otf")
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var otf OTFCommand
		err := rows.Scan(&otf.Name, &otf.Response, &otf.Creator, &otf.Count)
		if err != nil {
			fmt.Println("Error scanning: " + err.Error())
		}
		otfCommands = append(otfCommands, otf)
	}

	json.NewEncoder(w).Encode(otfCommands)
}

func getOTFCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	key := vars["name"]

	var queryString = "SELECT * FROM otf WHERE Name='" + key + "'"
	row, err := db.Query(queryString)
	if err != nil {
		panic(err.Error())
	}
	defer row.Close()

	var otf OTFCommand
	for row.Next() {
		err := row.Scan(&otf.Name, &otf.Response, &otf.Creator, &otf.Count)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(otf)
}
