// Package config handles the initialization and configuration of the database connection.
//
// It establishes a connection to a MySQL database using GORM, configures the connection pool,
// and performs automatic migrations for the Contact model.
//
// Author: Tri Wicaksono
// Website: https://triwicaksono.com
package config

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// DB is a global variable that holds the database connection instance.
// It is accessible throughout the application for executing database operations.
var DB *gorm.DB

// InitDB initializes the database connection using environment variables.
// It sets up the connection pool and performs automatic migrations for the Contact model.
//
// The function performs the following steps:
// 1. Retrieves database configuration from environment variables.
// 2. Constructs the Data Source Name (DSN) for MySQL connection.
// 3. Opens the database connection using GORM with a singular table naming strategy.
// 4. Configures the connection pool with specified limits.
// 5. Automatically migrates the Contact model to create or update the corresponding table.
//
// If any step fails, the function will panic with an appropriate error message.
func InitDB() {
	// Retrieve database configuration from environment variables with default values.
	dbUser := GetEnv("DB_USER", "root")
	dbPassword := GetEnv("DB_PASSWORD", "root")
	dbHost := GetEnv("DB_HOST", "127.0.0.1")
	dbPort := GetEnv("DB_PORT", "3306")
	dbName := GetEnv("DB_NAME", "aksara_batak_api")

	// Construct the Data Source Name (DSN) for MySQL connection.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error

	// Open the database connection using GORM with a singular table naming strategy.
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Configure the connection pool settings.
	sqlDB, err := DB.DB()
	if err != nil {
		panic("Failed to get database instance!")
	}

	sqlDB.SetMaxOpenConns(10)           // Maximum number of open connections to the database.
	sqlDB.SetMaxIdleConns(5)            // Maximum number of idle connections in the pool.
	sqlDB.SetConnMaxLifetime(time.Hour) // Maximum amount of time a connection may be reused.

	// Automatically migrate the Contact model to create or update the corresponding table.
	//if err := DB.AutoMigrate(
	//	&domain.User{},
	//	&domain.Product{},
	//	&domain.Category{},
	//	&domain.Cart{},
	//	&domain.CartItems{},
	//	&domain.Order{},
	//	&domain.OrderItem{},
	//	&domain.Payment{},
	//); err != nil {
	//	panic(fmt.Sprintf("Failed to migrate database: %v", err))
	//}
	//log.Println("Database migrated successfully!")
}
