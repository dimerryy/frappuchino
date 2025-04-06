package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"hot-coffee/internal/server"
	"hot-coffee/internal/utils"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func main() {
	db, err := CheckDb()
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	defer db.Close()

	utils.DB = db
	rows, err := db.Query("SELECT order_id, customer_name, total_amount FROM orders")
	if err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var orderID int
		var customerName string
		var totalAmount float64

		if err := rows.Scan(&orderID, &customerName, &totalAmount); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		fmt.Printf("Order ID: %d, Customer: %s, Total Amount: %.2f\n", orderID, customerName, totalAmount)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Row iteration error: %v", err)
	}
	DB = db

	server.StartTheCafe()
}

func CheckDb() (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
