package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ResponseCommon struct {
	StatusCode   int
	Header       http.Header
	Error        string `json:"error"`
	ErrorStatus  int    `json:"status"`
	ErrorMessage string `json:"message"`
}

type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
	CreatedAt       Time   `json:"created_at"`
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

/* func makeTwitchRequest(w http.ResponseWriter, url string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer"+os.Getenv("YBD_TOKEN"))
	req.Header.Add("Client-Id", os.Getenv("YBD_ID"))

	response, err := client.Do(req)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	return body
} */

func getTwitchUserId(w http.ResponseWriter, endpoint string) (data ManyUsers) {
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/users?"+endpoint, nil)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("YBD_TOKEN"))
	req.Header.Add("Client-Id", os.Getenv("YBD_ID"))

	response, err := client.Do(req)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	userInfo := ManyUsers{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		log.Fatal(err)
	}

	data = userInfo
	return
}

/* Twitch */
func getTwitchId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	parts := strings.Split(r.URL.Path, "/")
	userId := parts[3]

	// Get the endpoint by checking if param is a number or string
	var endPoint string
	val, err := strconv.Atoi(userId)
	if err != nil {
		endPoint = "login=" + userId
	} else {
		endPoint = "id=" + strconv.Itoa(val)
	}

	userInfo := getTwitchUserId(w, endPoint)

	json.NewEncoder(w).Encode(userInfo)
}

func getTwitchEmotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	parts := strings.Split(r.URL.Path, "/")
	userParam := parts[3]

	var userId int
	valParam, err := strconv.Atoi(userParam)
	if err != nil {
		userInfo := getTwitchUserId(w, "login="+userParam)
		userBody := userInfo.Users[0]
		userId, err = strconv.Atoi(userBody.ID)
		if err != nil {
			panic(err.Error())
		}
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
