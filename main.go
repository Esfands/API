package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"rpb-api3/database"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var db *sql.DB

var client *http.Client

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db = database.Connect()
	if db == nil {
		return
	}

	client = &http.Client{
		Timeout: time.Second * 10,
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePoint).Methods("GET")

	router.HandleFunc("/commands", getAllCommands).Methods("GET")
	router.HandleFunc("/commands/{name}", getCommand).Methods("GET")

	router.HandleFunc("/otf", getAllOTFCommands).Methods("GET")
	router.HandleFunc("/otf/{name}", getOTFCommand).Methods("GET")

	router.HandleFunc("/suggestions", getFeedback).Methods("GET")
	router.HandleFunc("/suggestions/{id}", getFeedbackById).Methods("GET")

	router.HandleFunc("/subathon/chatters", getSubathonActiveChatters).Methods("GET")
	router.HandleFunc("/subathon/giftedsubs", getSubathonGiftedSubs).Methods("GET")
	router.HandleFunc("/subathon/bitsdonated", getSubathonBitsDonated).Methods("GET")

	router.HandleFunc("/twitch/id", getTwitchId).Methods("GET")
	router.HandleFunc("/twitch/emotes", getTwitchEmotes).Methods("GET")

	router.HandleFunc("/eventsub/", eventsubRecievedNotification).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println("API listening on " + port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}

/* Home */
func homePoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "WideHardo")
}
