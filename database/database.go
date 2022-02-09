package database

import (
	"database/sql"
	"fmt"
	"os"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("MARIA_USER")+":"+os.Getenv("MARIA_PASSWORD")+"@tcp("+os.Getenv("MARIA_HOST")+":"+os.Getenv("MARIA_PORT")+")/"+os.Getenv("MARIA_DATABASE")+"?parseTime=true")
	if err != nil {
		fmt.Println("[ERROR] Unable to connect to the database: ", err)
		return nil
	}

	fmt.Println("[INFO] Connected to database")
	return db
}
