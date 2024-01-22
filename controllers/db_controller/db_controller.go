package db_controller

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	Types "github.com/nileshnk/golang-todo-app/types"
)

var (
	DBInstance     *sql.DB
	DBConnectError error
)

// ConnectToDatabase establishes a connection to the database.
func ConnectToDatabase() (*sql.DB, error) {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		return nil, errors.New("DATABASE_URL environment variable not set")
	}

	DBInstance, DBConnectError = sql.Open("postgres", connString)
	if DBConnectError != nil {
		log.Fatalf("Failed to connect to the database: %v", DBConnectError)
	}

	// Ping the database to ensure the connection is valid
	if err := DBInstance.Ping(); err != nil {
		DBInstance.Close()
		return nil, fmt.Errorf("Failed to ping the database: %v", err)
	}

	fmt.Println("Connected to the database")
	return DBInstance, nil
}

// GetDBInstance returns the existing database connection instance.
func GetDBInstance() *sql.DB {
	return DBInstance
}


func ApplyMigrations(db *sql.DB)  (Types.AppResponse, error){
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "todo"
	}
	if migrationsPath == "" {
		migrationsPath = "/database/migrations";
	}

	driver, driverErr := postgres.WithInstance(db, &postgres.Config{})
	if driverErr != nil {
		log.Fatal(driverErr)
		return Types.AppResponse{Success: false, Message: "Error applying migrations"}, driverErr
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:/%s", migrationsPath),
		dbName, driver)
	if err != nil {
		log.Fatal(err)
		return Types.AppResponse{Success: false, Message: "Error applying migrations"}, err
	}
	
	if err := m.Up(); err != nil {

		if err == migrate.ErrNoChange {
			return Types.AppResponse{Success: true, Message: "No migrations to apply"}, nil
		}

		return Types.AppResponse{Success: false, Message: "Error applying migrations"}, err
	}

	return Types.AppResponse{Success: true, Message: "Migrations applied successfully"}, nil
}

