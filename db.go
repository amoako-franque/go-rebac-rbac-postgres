package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=password dbname=rebac_rbac_db port=5432 sslmode=disable"
	}

	// Try to connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// If database doesn't exist, try to create it
	if err != nil && strings.Contains(err.Error(), "does not exist") {
		log.Println("Database does not exist, attempting to create it...")

		// Extract database name from DSN
		dbName := extractDBName(dsn)
		if dbName != "" {
			// Connect to postgres database to create the target database
			adminDSN := strings.Replace(dsn, " dbname="+dbName, " dbname=postgres", 1)
			adminDB, adminErr := gorm.Open(postgres.Open(adminDSN), &gorm.Config{})
			if adminErr != nil {
				return fmt.Errorf("failed to connect to postgres database: %w. Please create the database manually: CREATE DATABASE %s;", adminErr, dbName)
			}

			// Create the database using raw SQL
			sqlDB, _ := adminDB.DB()
			createDBQuery := fmt.Sprintf(`CREATE DATABASE "%s";`, dbName)
			_, execErr := sqlDB.Exec(createDBQuery)
			if execErr != nil {
				// Database might already exist, that's okay
				if !strings.Contains(execErr.Error(), "already exists") {
					log.Printf("Warning: Could not create database: %v", execErr)
				}
			} else {
				log.Printf("Database '%s' created successfully", dbName)
			}
			sqlDB.Close()

			// Try connecting again
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
		}
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	// AutoMigrate all models
	if err := DB.AutoMigrate(
		&Permission{},
		&Role{},
		&RolePermission{},
		&User{},
		&UserRole{},
		&PatientRecord{},
		&Relationship{},
	); err != nil {
		return err
	}

	log.Println("Database connected and migrated successfully")
	return nil
}

// extractDBName extracts the database name from a PostgreSQL DSN string
func extractDBName(dsn string) string {
	// Look for "dbname=" in the DSN
	parts := strings.Split(dsn, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "dbname=") {
			return strings.TrimPrefix(part, "dbname=")
		}
	}
	return ""
}
