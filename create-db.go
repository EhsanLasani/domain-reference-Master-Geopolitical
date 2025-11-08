package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

func main() {
	// Connect to postgres database to create geopolitical database
	dsn := "host=localhost port=5432 user=postgres password=@Salman2021 dbname=postgres sslmode=disable"
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open connection:", err)
	}
	defer db.Close()
	
	// Create geopolitical database
	_, err = db.Exec("CREATE DATABASE geopolitical")
	if err != nil {
		fmt.Printf("Database might already exist: %v\n", err)
	} else {
		fmt.Println("✅ Created geopolitical database!")
	}
}