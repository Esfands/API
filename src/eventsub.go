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
)

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
