package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// PromoCode struct for JSON serialization
type PromoCode struct {
	ID            int    `json:"id"`
	Code          string `json:"code"`
	IsBroadcasted bool   `json:"is_broadcasted"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

var db *sql.DB

// Setup the database connection
func setupDatabase() error {
	// Replace with your own DSN (Data Source Name)
	dsn := "root:@tcp(localhost:3306)/bandit" // No password, just use `root:`
	var err error
	db, err = sql.Open("mysql", dsn)
	return err // Don't forget to return the error if any
}

// Fetch latest promo codes from the database
func fetchLatestPromoCodes() ([]PromoCode, int, error) {
	// First, get the promo codes
	rows, err := db.Query("SELECT id, promo_code, is_broadcasted, created_at, updated_at FROM latest_promo_codes ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var promoCodes []PromoCode
	for rows.Next() {
		var promoCode PromoCode
		if err := rows.Scan(&promoCode.ID, &promoCode.Code, &promoCode.IsBroadcasted, &promoCode.CreatedAt, &promoCode.UpdatedAt); err != nil {
			return nil, 0, err
		}
		promoCodes = append(promoCodes, promoCode)
	}

	// Now, get the count of records where is_broadcasted == 0
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM latest_promo_codes WHERE is_broadcasted = 0").Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return promoCodes, count, nil
}

// InsertPromoCode inserts a new promo code into the latest_promo_codes table
func InsertPromoCode(promo string) error {
	// Check if the promo code already exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM latest_promo_codes WHERE promo_code = ?)", promo).Scan(&exists)
	if err != nil {
		return err // Return error if there's an issue with the query
	}

	if exists {
		return fmt.Errorf("promo code already exists") // Return an error if the promo code exists
	}

	// Insert the new promo code into the database
	_, err = db.Exec("INSERT INTO latest_promo_codes (promo_code, is_broadcasted, created_at, updated_at) VALUES (?, 0, NOW(), NOW())", promo)
	return err
}

// Close the database connection
func closeDatabase() {
	if db != nil {
		db.Close()
	}
}
