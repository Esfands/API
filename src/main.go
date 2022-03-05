package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Esfands/API/src/database"

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

	router := mux.NewRouter().StrictSlash(false)

	// Routes
	router.HandleFunc("/", homePoint).Methods("GET")

	// RPB
	router.HandleFunc("/rpb/commands", getAllCommands).Methods("GET")
	router.HandleFunc("/rpb/otf", getAllOTFCommands).Methods("GET")
	router.HandleFunc("/rpb/suggestions", getFeedback).Methods("GET")

	// Twitch
	router.HandleFunc("/twitch/resolve/{user}", getTwitchId).Methods("GET")
	router.HandleFunc("/twitch/emotes/{user}", getTwitchEmotes).Methods("GET")
	router.HandleFunc("/eventsub/", eventsubRecievedNotification).Methods("POST")

	// Misc
	router.HandleFunc("/emotes/endpoints", emoteEndpoints)

	// Middleware
	router.Use(logRequests)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println("API listening on " + port)
	log.Fatal(
		http.ListenAndServe(":"+port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(router)))
}

/* Home */
func homePoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode("WideHardo")
}

func StrToInt(str string) (int, error) {
	nonFractionalPart := strings.Split(str, ".")
	return strconv.Atoi(nonFractionalPart[0])
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

type ErrorMessage struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

func returnEndpointError(w http.ResponseWriter, reason string, code int) {
	data := ErrorMessage{
		Code:   code,
		Reason: reason,
	}

	json.NewEncoder(w).Encode(data)
}
