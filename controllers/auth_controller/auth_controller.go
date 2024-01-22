package auth_controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	Types "github.com/nileshnk/golang-todo-app/types"
)

// create signup and login functions

// create signup
var validate *validator.Validate = validator.New()

func SignUp(w http.ResponseWriter, r *http.Request) {
	// get user data from request body
	// var validateTrans = map[string]string{
	// 	"required": "The {0} field is required.",
	// 	"email":    "The {0} field must be a valid email address.",
	// 	"min":      "The {0} field must be at least {1} characters long.",
	// 	"max":      "The {0} field cannot exceed {1} characters in length.",
	// 	"eqfield":  "The {0} field must match the {1} field.",
	// 	"unique":   "The {0} field must be unique.",
	// 	// Add more translations as needed
	// }

	var userInfo Types.UserRegistration;
	json.NewDecoder(r.Body).Decode(&userInfo)

	// validate user data
	if err := validate.Struct(userInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Types.AppResponse{
			Success: false,
			Message: "Invalid user data",
			Data: "",//err.(validator.ValidationErrors).Translate(validateTrans),
		})
		// fmt.Println(err.(validator.ValidationErrors).Translate(validateTrans))
		fmt.Println(err)
		return
	}

	// name email password password_confirmation
	// validate user data
	// create user
	// return user data
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	// get user data from request body
	// validate user data
	// return user data
}