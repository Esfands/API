package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type User struct {
	ID              string    `json:"id"`
	Login           string    `json:"login"`
	DisplayName     string    `json:"display_name"`
	Type            string    `json:"type"`
	BroadcasterType string    `json:"broadcaster_type"`
	Description     string    `json:"description"`
	ProfileImageURL string    `json:"profile_image_url"`
	OfflineImageURL string    `json:"offline_image_url"`
	ViewCount       int       `json:"view_count"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
}

type ManyUsers struct {
	Users []User `json:"data"`
}

type UsersResponse struct {
	Data ManyUsers
}

type Emote struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Images     EmoteImage `json:"images"`
	Tier       string     `json:"tier"`
	EmoteType  string     `json:"emote_type"`
	EmoteSetId string     `json:"emote_set_id"`
}

type EmoteImage struct {
	Url1x string `json:"url_1x"`
	Url2x string `json:"url_2x"`
	Url4x string `json:"url_4x"`
}

type ManyEmotes struct {
	Emotes []Emote `json:"data"`
}

/* Twitch */
func getTwitchId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	keys, ok := r.URL.Query()["user"]

	if !ok || len(keys[0]) < 1 {
		json.NewEncoder(w).Encode("URL param 'user' is missing")
		return
	}

	userId := keys[0]

	// Get the endpoint by checking if param is a number or string
	var endPoint string
	val, err := strconv.Atoi(userId)
	if err != nil {
		endPoint = "login=" + userId
	} else {
		fmt.Println(val)
		endPoint = "id=" + string(val)
	}

	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/users?"+endPoint, nil)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("YBD_TOKEN"))
	req.Header.Add("Client-Id", os.Getenv("YBD_ID"))

	response, err := client.Do(req)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	var userInfo ManyUsers
	if err := json.Unmarshal(body, &userInfo); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(userInfo)
}

func getTwitchEmotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	keys, ok := r.URL.Query()["user"]

	if !ok || len(keys[0]) < 1 {
		json.NewEncoder(w).Encode("URL param 'user' is missing")
		return
	}

	userParam := keys[0]

	var userId int
	valParam, err := strconv.Atoi(userParam)
	if err != nil {
		// TODO: Make it so it fetches the ID since they put in a username, for now it returns an error
		json.NewEncoder(w).Encode(userParam + " is not a valid ID")
		return
	} else {
		// It's an ID so get the emotes
		userId = valParam
	}

	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/chat/emotes?broadcaster_id="+strconv.Itoa(userId), nil)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("YBD_TOKEN"))
	req.Header.Add("Client-Id", os.Getenv("YBD_ID"))

	response, err := client.Do(req)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	var emoteSets ManyEmotes
	if err := json.Unmarshal(body, &emoteSets); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(emoteSets)
}
