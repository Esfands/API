package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Subathon
type SubathonStat struct {
	ID           int    `json:"ID"`
	Username     string `json:"Username"`
	MessageCount int    `json:"MessageCount"`
	GiftedSubs   int    `json:"GiftedSubs"`
	BitsDonated  int    `json:"BitsDonated"`
}

/* Subathon Stats */
func getSubathonActiveChatters(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	keys, ok := r.URL.Query()["offset"]

	if !ok || len(keys[0]) < 1 {
		// TODO: Add default route in case there isn't an offset
		return
	}

	offset, err := strconv.Atoi(keys[0])
	if err != nil {
		panic(err.Error())
	}
	var activeChatters []SubathonStat

	rows, err := db.Query("SELECT * FROM subathonstats ORDER BY MessageCount DESC LIMIT 25 OFFSET " + strconv.Itoa(offset-1))
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var messages SubathonStat
		err := rows.Scan(&messages.ID, &messages.Username, &messages.MessageCount, &messages.GiftedSubs, &messages.BitsDonated)
		if err != nil {
			fmt.Println("Error scanning: " + err.Error())
		}
		activeChatters = append(activeChatters, messages)
	}
	json.NewEncoder(w).Encode(activeChatters)
}

func getSubathonGiftedSubs(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	keys, ok := r.URL.Query()["offset"]

	if !ok || len(keys[0]) < 1 {
		// TODO: Add default route in case there isn't an offset
		return
	}

	offset, err := strconv.Atoi(keys[0])
	if err != nil {
		panic(err.Error())
	}
	var giftedSubs []SubathonStat

	rows, err := db.Query("SELECT * FROM subathonstats ORDER BY GiftedSubs DESC LIMIT 25 OFFSET " + strconv.Itoa(offset-1))
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var subs SubathonStat
		err := rows.Scan(&subs.ID, &subs.Username, &subs.MessageCount, &subs.GiftedSubs, &subs.BitsDonated)
		if err != nil {
			fmt.Println("Error scanning: " + err.Error())
		}
		giftedSubs = append(giftedSubs, subs)
	}
	json.NewEncoder(w).Encode(giftedSubs)
}

func getSubathonBitsDonated(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	keys, ok := r.URL.Query()["offset"]

	if !ok || len(keys[0]) < 1 {
		// TODO: Add default route in case there isn't an offset
		return
	}

	offset, err := strconv.Atoi(keys[0])
	if err != nil {
		panic(err.Error())
	}

	var bitsDonated []SubathonStat

	rows, err := db.Query("SELECT * FROM subathonstats ORDER BY BitsDonated DESC LIMIT 25 OFFSET " + strconv.Itoa(offset-1))
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var bits SubathonStat
		err := rows.Scan(&bits.ID, &bits.Username, &bits.MessageCount, &bits.GiftedSubs, &bits.BitsDonated)
		if err != nil {
			fmt.Println("Error scanning: " + err.Error())
		}
		bitsDonated = append(bitsDonated, bits)
	}
	json.NewEncoder(w).Encode(bitsDonated)
}
