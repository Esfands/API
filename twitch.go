package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

// Represents a subscription
type EventSubSubscription struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Status    string            `json:"status"`
	Condition EventSubCondition `json:"condition"`
	Transport EventSubTransport `json:"transport"`
	CreatedAt Time              `json:"created_at"`
	Cost      int               `json:"cost"`
}

type EventSubCondition struct {
	BroadcasterUserID     string `json:"broadcaster_user_id"`
	FromBroadcasterUserID string `json:"from_broadcaster_user_id"`
	ToBroadcasterUserID   string `json:"to_broadcaster_user_id"`
	RewardID              string `json:"reward_id"`
	ClientID              string `json:"client_id"`
	ExtensionClientID     string `json:"extension_client_id"`
	UserID                string `json:"user_id"`
}

// Transport for the subscription, currently the only supported Method is "webhook". Secret must be between 10 and 100 characters
type EventSubTransport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

type Pagination struct {
	Cursor string `json:"cursor"`
}

// Twitch Response for getting all current subscriptions
type ManyEventSubSubscriptions struct {
	TotalCost             int                    `json:"total_cost"`
	MaxTotalCost          int                    `json:"max_total_cost"`
	EventSubSubscriptions []EventSubSubscription `json:"data"`
	Pagination            Pagination             `json:"pagination"`
}

type EventSubChannelPredictionBeginEvent struct {
	ID                   string            `json:"id"`
	BroadcasterUserID    string            `json:"broadcaster_user_id"`
	BroadcasterUserLogin string            `json:"broadcaster_user_login"`
	BroadcasterUserName  string            `json:"broadcaster_user_name"`
	Title                string            `json:"title"`
	Outcomes             []EventSubOutcome `json:"outcomes"`
	StartedAt            Time              `json:"started_at"`
	LocksAt              Time              `json:"locks_at"`
}

// Data for a channel prediction progress event
type EventSubChannelPredictionProgressEvent = EventSubChannelPredictionBeginEvent

// Data for a channel prediction lock event
type EventSubChannelPredictionLockEvent struct {
	ID                   string            `json:"id"`
	BroadcasterUserID    string            `json:"broadcaster_user_id"`
	BroadcasterUserLogin string            `json:"broadcaster_user_login"`
	BroadcasterUserName  string            `json:"broadcaster_user_name"`
	Title                string            `json:"title"`
	WinningOutcomeID     string            `json:"winning_outcome_id"`
	Outcomes             []EventSubOutcome `json:"outcomes"`
	Status               string            `json:"status"`
	StartedAt            Time              `json:"started_at"`
	LockedAt             Time              `json:"locked_at"`
}

// Data for a channel prediction end event
type EventSubChannelPredictionEndEvent struct {
	ID                   string            `json:"id"`
	BroadcasterUserID    string            `json:"broadcaster_user_id"`
	BroadcasterUserLogin string            `json:"broadcaster_user_login"`
	BroadcasterUserName  string            `json:"broadcaster_user_name"`
	Title                string            `json:"title"`
	WinningOutcomeID     string            `json:"winning_outcome_id"`
	Outcomes             []EventSubOutcome `json:"outcomes"`
	Status               string            `json:"status"`
	StartedAt            Time              `json:"started_at"`
	EndedAt              Time              `json:"eneded_at"`
}

type EventSubTopPredictor struct {
	UserID            string `json:"user_id"`
	UserLogin         string `json:"user_login"`
	UserName          string `json:"user_name"`
	ChannelPointWon   int    `json:"channel_points_won"`
	ChannelPointsUsed int    `json:"channel_points_used"`
}

type EventSubOutcome struct {
	ID            string                 `json:"id"`
	Title         string                 `json:"title"`
	Color         string                 `json:"color"`
	Users         int                    `json:"users"`
	ChannelPoints int                    `json:"channel_points"`
	TopPredictors []EventSubTopPredictor `json:"top_predictors"`
}

// Data for a stream online notification
type EventSubStreamOnlineEvent struct {
	ID                   string `json:"id"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Type                 string `json:"type"`
	StartedAt            Time   `json:"started_at"`
}

// Data for a stream offline notification
type EventSubStreamOfflineEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
}

// Data for a channel update notification
type EventSubChannelUpdateEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Title                string `json:"title"`
	Language             string `json:"language"`
	CategoryID           string `json:"category_id"`
	CategoryName         string `json:"category_name"`
	IsMature             bool   `json:"is_mature"`
}

// Data for a channel poll begin event
type EventSubChannelPollBeginEvent struct {
	ID                   string                      `json:"id"`
	BroadcasterUserID    string                      `json:"broadcaster_user_id"`
	BroadcasterUserLogin string                      `json:"broadcaster_user_login"`
	BroadcasterUserName  string                      `json:"broadcaster_user_name"`
	Title                string                      `json:"title"`
	Choices              []PollChoice                `json:"choices"`
	BitsVoting           EventSubBitVoting           `json:"bits_voting"`
	ChannelPointsVoting  EventSubChannelPointsVoting `json:"channel_points_voting"`
	StartedAt            Time                        `json:"started_at"`
	EndsAt               Time                        `json:"ends_at"`
}

type PollChoice struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	BitsVotes          int    `json:"bits_votes"`
	ChannelPointsVotes int    `json:"channel_points_votes"`
	Votes              int    `json:"votes"`
}

type EventSubBitVoting struct {
	IsEnabled     bool `json:"is_enabled"`
	AmountPerVote int  `json:"amount_per_vote"`
}

type EventSubChannelPollEndEvent struct {
	ID                   string                      `json:"id"`
	BroadcasterUserID    string                      `json:"broadcaster_user_id"`
	BroadcasterUserLogin string                      `json:"broadcaster_user_login"`
	BroadcasterUserName  string                      `json:"broadcaster_user_name"`
	Title                string                      `json:"title"`
	Choices              []PollChoice                `json:"choices"`
	BitsVoting           EventSubBitVoting           `json:"bits_voting"`
	ChannelPointsVoting  EventSubChannelPointsVoting `json:"channel_points_voting"`
	Status               string                      `json:"status"`
	StartedAt            Time                        `json:"started_at"`
	EndedAt              Time                        `json:"ended_at"`
}

type EventSubChannelPointsVoting = EventSubBitVoting

// Data for a channel poll progress event, it's the same as the channel poll begin event
type EventSubChannelPollProgressEvent = EventSubChannelPollBeginEvent

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

// Verify message from EventSub
func VerifyEventSubNotification(secret string, header http.Header, message string) bool {
	hmacMessage := []byte(fmt.Sprintf("%s%s%s", header.Get("Twitch-Eventsub-Message-Id"), header.Get("Twitch-Eventsub-Message-Timestamp"), message))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(hmacMessage)
	hmacsha256 := fmt.Sprintf("sha256=%s", hex.EncodeToString(mac.Sum(nil)))
	return hmacsha256 == header.Get("Twitch-Eventsub-Message-Signature")
}

type eventsubNotification struct {
	Subscription EventSubSubscription `json:"subscription"`
	Challenge    string               `json:"challenge"`
	Event        json.RawMessage      `json:"event"`
}

// Handles all eventsub events.
func eventsubRecievedNotification(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	defer r.Body.Close()
	// Verify that Twitch sent the message
	if !VerifyEventSubNotification(os.Getenv("EVENTSUB_SECRET"), r.Header, string(body)) {
		log.Println("No valid signature on subscription")
		return
	} else {
		log.Println("Verified signature on subscription")
	}
	var vals eventsubNotification
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&vals)
	if err != nil {
		log.Println(err)
		return
	}
	// if there's a challenge in the request, respond with only the challenge to verify your eventsub.
	if vals.Challenge != "" {
		w.Write([]byte(vals.Challenge))
		return
	}

	eventType := bytes.NewBuffer([]byte(vals.Subscription.Type)).String()

	if eventType == "stream.online" {
		var streamOnline EventSubStreamOnlineEvent
		err := json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamOnline)
		if err != nil {
			panic(err.Error())
		}

		log.Printf("%s just went online.", streamOnline.BroadcasterUserLogin)
		w.WriteHeader(200)
		w.Write([]byte("ok"))
		return

	} else if eventType == "stream.offline" {
		var streamOffline EventSubStreamOnlineEvent
		err := json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamOffline)
		if err != nil {
			panic(err.Error())
		}

		log.Printf("%s just went offline.", streamOffline.BroadcasterUserLogin)
		w.WriteHeader(200)
		w.Write([]byte("ok"))
		return

	} else if eventType == "stream.update" {
		fmt.Println("stream.update")
		var streamUpdate EventSubChannelUpdateEvent
		err := json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamUpdate)
		if err != nil {
			panic(err.Error())
		}

		log.Printf("%s just went offline.", streamUpdate.BroadcasterUserLogin)
		w.WriteHeader(200)
		w.Write([]byte("ok"))
		return

	} else if eventType == "channel.prediction.begin" {
		fmt.Println("channel.prediction.begin")
		var streamPredictionBegin EventSubChannelPredictionBeginEvent
		err := json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamPredictionBegin)
		if err != nil {
			panic(err.Error())
		}

		log.Printf("%s just started a prediction: %s", streamPredictionBegin.BroadcasterUserLogin, streamPredictionBegin.Title)
		w.WriteHeader(200)
		w.Write([]byte("ok"))
		return

	} else if eventType == "channel.prediction.progress" {
		fmt.Println("channel.prediction.progress")
		var streamPredictionProgress EventSubChannelPredictionProgressEvent
		err := json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamPredictionProgress)
		if err != nil {
			panic(err.Error())
		}

		log.Printf("%s just started a prediction: %s", streamPredictionProgress.BroadcasterUserLogin, streamPredictionProgress.Title)
		w.WriteHeader(200)
		w.Write([]byte("ok"))
		return

	} else if eventType == "channel.prediction.lock" {
		fmt.Println("channel.prediction.lock")
		var streamPredictionLock EventSubChannelPredictionLockEvent
		err := json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamPredictionLock)
		if err != nil {
			panic(err.Error())
		}

		log.Printf("%s just locked a prediction: %s", streamPredictionLock.BroadcasterUserLogin, streamPredictionLock.Title)
		w.WriteHeader(200)
		w.Write([]byte("ok"))
		return

	} else if eventType == "channel.prediction.end" {
		fmt.Println("channel.prediction.end")
		var streamPredictionEnd EventSubChannelPredictionEndEvent
		err := json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamPredictionEnd)
		if err != nil {
			panic(err.Error())
		}

		log.Printf("%s just ended a prediction: %s", streamPredictionEnd.BroadcasterUserLogin, streamPredictionEnd.Title)
		w.WriteHeader(200)
		w.Write([]byte("ok"))
		return

	} else if eventType == "channel.poll.begin" {
		fmt.Println("channel.poll.begin")
		var streamPollBegan EventSubChannelPollBeginEvent
		err := json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamPollBegan)
		if err != nil {
			panic(err.Error())
		}

		log.Printf("%s just started a poll: %s", streamPollBegan.BroadcasterUserLogin, streamPollBegan.Title)
		w.WriteHeader(200)
		w.Write([]byte("ok"))
		return

	} else if eventType == "channel.poll.progress" {
		fmt.Println("channel.poll.progress")
		var streamPollProgress EventSubChannelPollProgressEvent
		err := json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamPollProgress)
		if err != nil {
			panic(err.Error())
		}

		log.Printf("%s progess for poll: %s", streamPollProgress.BroadcasterUserLogin, streamPollProgress.Title)
		w.WriteHeader(200)
		w.Write([]byte("ok"))
		return

	} else if eventType == "channel.poll.end" {
		fmt.Println("channel.poll.end")
		var streamPollEnd EventSubChannelPollEndEvent
		err := json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamPollEnd)
		if err != nil {
			panic(err.Error())
		}

		log.Printf("%s ended poll: %s", streamPollEnd.BroadcasterUserLogin, streamPollEnd.Title)
		w.WriteHeader(200)
		w.Write([]byte("ok"))
		return

	}
}
