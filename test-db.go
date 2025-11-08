package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

func main() {
	// Test connection with your exact credentials
	dsn := "host=localhost port=5432 user=postgres password=@Salman2021 dbname=referencemaster sslmode=disable"
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open connection:", err)
	}
	defer db.Close()
	
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	
	fmt.Println("✅ Successfully connected to referencemaster database!")
	
	// Test if geopolitical database exists
	dsn2 := "host=localhost port=5432 user=postgres password=@Salman2021 dbname=geopolitical sslmode=disable"
	db2, err := sql.Open("postgres", dsn2)
	if err == nil {
		if err := db2.Ping(); err == nil {
			fmt.Println("✅ geopolitical database also exists!")
		} else {
			fmt.Println("❌ geopolitical database does not exist")
		}
		db2.Close()
	}
}