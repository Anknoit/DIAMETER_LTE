package server

import (
	"epiphanyHSS/modules/common"
	"epiphanyHSS/modules/db/mysql"
	"log"
)

// InitializeServer initializes all necessary configurations and services
func InitializeServer() {

	// Initialize the Client configuration
	if err := common.LoadMMEConfig("conf/clients.yaml"); err != nil {
		log.Fatalf("Client configuration initialization error: %v", err)
	}
	log.Println("Client configuration loaded successfully.")

	// Initialize the RDS (database) configuration
	if err := common.LoadRDSConfig("conf/rds.yaml"); err != nil {
		log.Fatalf("RDS configuration initialization error: %v", err)
	}
	log.Println("RDS configuration loaded successfully.")

	// Initialize the MySQL database connection pool
	if err := mysql.InitializeDB(); err != nil {
		log.Fatalf("Database connection pool initialization error: %v", err)
	}
	log.Println("Database connection pool initialized successfully.")

}
