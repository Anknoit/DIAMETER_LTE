package mysql

import (
	"database/sql"
	"epiphanyHSS/modules/common"
	"fmt"
	"log"
	"time"

	_ "epiphanyHSS/modules/go-mysql" // MySQL driver for Go
)

// Global variable to hold the database connection pool
var DB *sql.DB

// InitializeDB initializes the database connection pool using RDS settings from the common package
func InitializeDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		common.RDS.Database.User,
		common.RDS.Database.Password,
		common.RDS.Database.Host,
		common.RDS.Database.Port,
		common.RDS.Database.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		// log.Printf("Database connection error: %v", err)
		return err
	}

	// Configure the connection pool settings
	db.SetMaxOpenConns(common.RDS.Pool.MaxOpenConns)
	db.SetMaxIdleConns(common.RDS.Pool.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(common.RDS.Pool.ConnMaxLifetime) * time.Second)

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		// log.Printf("Database ping error: %v", err)
		db.Close()
		return err
	}

	DB = db
	return nil
}

// CheckDBStats logs the current statistics of the database connection pool
func CheckDBStats() {
	if DB == nil {
		// log.Println("Database is not initialized.")
		return
	}

	stats := DB.Stats()
	log.Printf("Open Connections: %d, Idle Connections: %d, In Use Connections: %d, Wait Count: %d",
		stats.OpenConnections, stats.Idle, stats.InUse, stats.WaitCount)

	// log.Println("Database stats logged successfully.")
}

// CloseDB closes the database connection pool
func CloseDB() {
	if DB != nil {
		DB.Close()
		// log.Println("Database connection pool closed.")
	}
}
