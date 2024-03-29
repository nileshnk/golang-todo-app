package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/nileshnk/golang-todo-app/controllers/auth_controller"
	"github.com/nileshnk/golang-todo-app/controllers/db_controller"
	Router "github.com/nileshnk/golang-todo-app/router"
)

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		panic(envErr)
	}

	DBInstance, dbConnectErr := db_controller.ConnectToDatabase()
	if dbConnectErr != nil {
		panic(dbConnectErr)
	}
	log.Println("Connected to database!")

	migrationsStatus, migrationsErr := db_controller.ApplyMigrations(DBInstance)
	if migrationsErr != nil {
		fmt.Println(migrationsErr)
		// panic(migrationsErr)
	}
	log.Println(migrationsStatus.Message)

	RouteHandler := chi.NewRouter()

	RouteHandler.Use(auth_controller.AuthMiddleware)

	RouteHandler.Route("/", Router.MainRouter)

	PORT := os.Getenv("APP_PORT")
	if PORT == "" {
		PORT = "4000"
	}
	ServerAddress := fmt.Sprintf("127.0.0.1:%s", PORT)
	http.ListenAndServe(ServerAddress, RouteHandler)

}
