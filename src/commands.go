package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	name := r.URL.Query().Get("name")

	var cmd Command
	if len(name) != 0 {
		err := db.QueryRow("SELECT * FROM commands WHERE Name=?", name).Scan(&cmd.Name, &cmd.Aliases, &cmd.Permissions, &cmd.Description, &cmd.DynamicDescription, &cmd.GlobalCooldown, &cmd.Cooldown, &cmd.Testing, &cmd.OfflineOnly, &cmd.OnlineOnly, &cmd.Count)

		switch {
		case err == sql.ErrNoRows:
			returnEndpointError(w, "Couldn't find the name: "+name, 404)

		case err != nil:
			log.Fatal(err)

		default:
			json.NewEncoder(w).Encode(cmd)
		}
		return
	}

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

/* OTF Commands */
func getAllOTFCommands(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	nameQ := r.URL.Query().Get("name")

	var otf OTFCommand
	if len(nameQ) != 0 {
		err := db.QueryRow("SELECT * FROM otf WHERE Name=?", nameQ).Scan(&otf.Name, &otf.Response, &otf.Creator, &otf.Count)

		switch {
		case err == sql.ErrNoRows:
			returnEndpointError(w, "Couldn't find the name: "+nameQ, 404)

		case err != nil:
			log.Fatal(err)

		default:
			json.NewEncoder(w).Encode(otf)
		}
		return
	}

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
