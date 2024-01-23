package auth_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nileshnk/golang-todo-app/controllers/db_controller"
	Types "github.com/nileshnk/golang-todo-app/types"
)

// create signup and login functions

// create signup
var validate *validator.Validate = validator.New()

func SignUp(w http.ResponseWriter, r *http.Request) {

	var userInfo Types.UserRegistration
	json.NewDecoder(r.Body).Decode(&userInfo)
	fmt.Println(userInfo)
	err := validate.Struct(userInfo)
	errs := translateError(err)
	if err != nil {
		// fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Types.AppResponse{
			Success: false,
			Message: "Invalid user data",
			Data:    errs, //err.(validator.ValidationErrors).Translate(validateTrans),
		})
		return

	}

	// fmt.Println(errs)
	fmt.Println("Signup called")
	if db_controller.DBInstance == nil {

		json.NewEncoder(w).Encode(Types.AppResponse{
			Success: false,
			Message: "Error connecting to database",
			Data:    nil,
		})
		return
	}

	// Execute the query
	rows, err := db_controller.DBInstance.Query("SELECT * FROM users WHERE email=$1;", userInfo.Email)
	if err != nil {
		log.Println("Error executing query:", err)
		json.NewEncoder(w).Encode(Types.AppResponse{
			Success: false,
			Message: "Internal Server Error 1",
			Data:    nil,
		})
		return
	}
	defer rows.Close()

	var users []Types.User

	for rows.Next() {
		var user Types.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			json.NewEncoder(w).Encode(Types.AppResponse{
				Success: false,
				Message: "Internal Server Error 2",
				Data:    nil,
			})
			return
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		json.NewEncoder(w).Encode(Types.AppResponse{
			Success: false,
			Message: "Internal Server Error 3",
			Data:    nil,
		})
		return
	}

	if len(users) > 0 {
		json.NewEncoder(w).Encode(Types.AppResponse{
			Success: false,
			Message: "User with same email already exists",
			Data:    nil,
		})
		return
	}

	// create user
	createUserId, createUserIdErr := uuid.NewRandom()
	if createUserIdErr != nil {
		log.Println("Error generating user id:", createUserIdErr)
		json.NewEncoder(w).Encode(Types.AppResponse{
			Success: false,
			Message: "Internal Server Error",
			Data:    nil,
		})
		return
	}

	_, createUserErr := db_controller.DBInstance.Exec("INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4);", createUserId, userInfo.Name, userInfo.Email, userInfo.Password)
	if createUserErr != nil {
		log.Println("Error executing query:", createUserErr)
		json.NewEncoder(w).Encode(Types.AppResponse{
			Success: false,
			Message: "Internal Server Error",
			Data:    nil,
		})
		return
	}

	fmt.Println("User created successfully")
	accessToken, tokenErr := CreateAccessToken(Types.TokenPayload{
		UserId: createUserId,
	})
	if tokenErr != nil {
		log.Println("Error creating access token:", tokenErr)
	}

	json.NewEncoder(w).Encode(Types.AppResponse{
		Success: true,
		Message: "User created successfully",
		Data:    map[string]string{"access_token": accessToken},
	})
	return
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	// get user data from request body
	var userInfo Types.SignInPayload
	json.NewDecoder(r.Body).Decode(&userInfo)

	if db_controller.DBInstance == nil {
		json.NewEncoder(w).Encode(Types.AppResponse{
			Success: false,
			Message: "Error connecting to database",
			Data:    nil,
		})
	}

	// Execute the query
	rows, err := db_controller.DBInstance.Query("SELECT * FROM users WHERE email=$1;", userInfo.Email)
	if err != nil {
		log.Println("Error executing query:", err)
		json.NewEncoder(w).Encode(Types.AppResponse{
			Success: false,
			Message: "Internal Server Error 1",
			Data:    nil,
		})
	}
	defer rows.Close()

	var users []Types.User

	for rows.Next() {
		var user Types.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			json.NewEncoder(w).Encode(Types.AppResponse{
				Success: false,
				Message: "Internal Server Error 2",
				Data:    nil,
			})
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		json.NewEncoder(w).Encode(Types.AppResponse{
			Success: false,
			Message: "Internal Server Error 3",
			Data:    nil,
		})
		return
	}

	if len(users) == 0 {
		json.NewEncoder(w).Encode(Types.AppResponse{
			Success: false,
			Message: "User doesn't exists",
			Data:    nil,
		})
		return
	}

	if len(users) > 1 {
		json.NewEncoder(w).Encode(Types.AppResponse{
			Success: false,
			Message: "Multiple users with same email exists",
			Data:    nil,
		})
		return
	}

	if users[0].Password != userInfo.Password {
		json.NewEncoder(w).Encode(Types.AppResponse{
			Success: false,
			Message: "Invalid password",
			Data:    nil,
		})
		return
	}

	accessToken, tokenErr := CreateAccessToken(Types.TokenPayload{
		UserId: users[0].Id,
	})
	if tokenErr != nil {
		log.Println("Error creating access token:", tokenErr)
	}

	ok, claims, validateErr := ValidateAccessToken(accessToken)
	if validateErr != nil {
		log.Println("Error validating access token:", validateErr)
	}
	fmt.Println(ok, claims)

	w.Header().Set("Authorization", "Bearer "+accessToken)
	w.Header().Add("Set-Cookie", "access_token="+accessToken+"; Path=/; HttpOnly; SameSite=Strict")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cookies", "access_token="+accessToken+"; Path=/; HttpOnly; SameSite=Strict")

	json.NewEncoder(w).Encode(Types.AppResponse{
		Success: true,
		Message: "User logged in successfully",
		Data:    map[string]string{"access_token": accessToken},
	})

	return

	// validate user data
	// return user data
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	// delete access token from cookies
	// return success message

	w.Header().Del("Authorization")
	w.Header().Del("Set-Cookie")
	w.Header().Del("Cookies")
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	json.NewEncoder(w).Encode(Types.AppResponse{
		Success: true,
		Message: "User logged out successfully",
		Data:    nil,
	})
}
