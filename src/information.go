package main

import (
	"encoding/json"
	"net/http"
)

type EndpointData struct {
	Endpoints []EndpointInformation `json:"data"`
}

type EndpointInformation struct {
	Name        string
	Scope       string
	Description string
	URL         string
}

func emoteEndpoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	endpoints := []EndpointInformation{
		{
			Name:        "bttv",
			Scope:       "global",
			Description: "Global emotes.",
			URL:         "https://api.betterttv.net/3/cached/emotes/global",
		},
		{
			Name:        "bttv",
			Scope:       "global",
			Description: "Badges for BTTV, mostly for BTTV dev/support/mod teams.",
			URL:         "https://api.betterttv.net/3/cached/badges",
		},
		{
			Name:        "bttv",
			Scope:       "channel",
			Description: "Fetch channel specific emotes by user ID.",
			URL:         "https://api.betterttv.net/3/cached/users/twitch/000000",
		},
		{
			Name:        "ffz",
			Scope:       "channel",
			Description: "Fetch channel specific emotes by user ID.",
			URL:         "https://api.betterttv.net/3/cached/frankerfacez/users/twitch/000000",
		},
	}

	json.NewEncoder(w).Encode(endpoints)
}
