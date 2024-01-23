package todo_controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nileshnk/golang-todo-app/controllers/db_controller"
	Types "github.com/nileshnk/golang-todo-app/types"
)

var tasks []Types.Task

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		http.Error(w, "UserId not found in context", http.StatusInternalServerError)
		return
	}
	fmt.Println(userId)

	allTasks, getErr := getAllTasks(userId)
	if getErr != nil {
		fmt.Println("Error getting all tasks")
		fmt.Println(getErr)
		w.Write([]byte("Error getting all tasks"))
		return
	}
	w.Header().Set("Content-Type", "application/json")

	jsonData, jsonDataError := json.Marshal(allTasks)
	if jsonDataError != nil {
		panic(jsonDataError)
	}
	w.Write(jsonData)
}

func AddTask(w http.ResponseWriter, r *http.Request) {

	var task Types.Task
	json.NewDecoder(r.Body).Decode(&task)
	task.UserID = r.Context().Value("userId").(string)
	addTask(&task)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func ChangeTaskStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	taskId, err := strconv.Atoi(id)

	if err != nil {
		fmt.Println("error parsing id from string to int")
		fmt.Println(err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	// requestBody := make(map[string]interface{})
	// decoder.Decode(&requestBody)
	// fmt.Println(requestBody["completed"])
	type completed struct {
		Completed bool `json:"completed"`
	}
	var taskCompleted completed

	decoder.Decode(&taskCompleted)
	updateTaskStatus(int64(taskId), taskCompleted.Completed)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
	fmt.Println(tasks)
}

func EditTaskText(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	taskId, err := strconv.Atoi(id)
	if err != nil {
		w.Write([]byte("not a valid id"))
		return
	}

	decoder := json.NewDecoder(r.Body)

	type text struct {
		Text string `json:"text"`
	}
	var taskText text

	decoder.Decode(&taskText)

	updateTaskText(int64(taskId), taskText.Text)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	taskId, err := strconv.Atoi(id)
	if err != nil {
		w.Write([]byte("not a valid id"))
	}

	deleteTask(int64(taskId))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
	fmt.Println(tasks)
}

func getAllTasks(userId string) ([]Types.Task, error) {
	fmt.Println("Getting all tasks")

	// Ensure there is a valid database connection
	if db_controller.DBInstance == nil {
		return nil, errors.New("database connection is nil")
	}

	// Execute the query
	rows, err := db_controller.DBInstance.Query("SELECT * FROM todo WHERE user_id=$1;", userId)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	// Create a slice to hold the retrieved tasks
	var tasks []Types.Task

	// Iterate through the result set and populate the tasks slice
	for rows.Next() {
		var task Types.Task
		err := rows.Scan(&task.Id, &task.UserID, &task.Text, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	fmt.Println(tasks)

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}

	return tasks, nil
}

func getTask(taskId int64) Types.Task {
	for i := 0; i < len(tasks); i++ {
		if tasks[i].Id == taskId {
			return tasks[i]
		}
	}
	return Types.Task{}
}

func addTask(task *Types.Task) (Types.AppResponse, error) {
	fmt.Println("Adding a task")

	// Ensure there is a valid database connection
	if db_controller.DBConnectError != nil {
		return Types.AppResponse{}, errors.New("database connection is nil")
	}

	// Execute the INSERT query
	err := db_controller.DBInstance.QueryRow("INSERT INTO todo (text, completed, user_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at", task.Text, task.Completed, task.UserID).Scan(&task.Id, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		log.Println("Error executing query:", err)
		return Types.AppResponse{}, err
	}

	return Types.AppResponse{Success: true, Message: "Task added successfully!", Data: task}, nil
}

func updateTaskText(taskId int64, newTaskText string) (*Types.Task, error) {

	// Ensure there is a valid database connection
	if db_controller.DBInstance == nil {
		return nil, errors.New("database connection is nil")
	}

	// Execute the query
	var id int64
	err := db_controller.DBInstance.QueryRow("UPDATE todo SET text=$1 WHERE id=$2 RETURNING id;", newTaskText, taskId).Scan(&id)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}

	return &Types.Task{Id: id, Text: newTaskText}, nil
}

func updateTaskStatus(taskId int64, status bool) (*Types.Task, error) {
	// Ensure there is a valid database connection
	if db_controller.DBInstance == nil {
		return nil, errors.New("database connection is nil")
	}

	// Execute the query
	var id int64
	err := db_controller.DBInstance.QueryRow("UPDATE todo SET completed=$1 WHERE id=$2 RETURNING id;", status, taskId).Scan(&id)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}

	return &Types.Task{Id: id, Completed: status}, nil
}

func deleteTask(taskId int64) (*Types.Task, error) {

	if db_controller.DBInstance == nil {
		return nil, errors.New("database connection is nil")
	}

	// Execute the query
	var id int64
	err := db_controller.DBInstance.QueryRow("DELETE FROM todo WHERE id=$1 RETURNING id;", taskId).Scan(&id)
	if err != nil {
		log.Println("Error executing query:", err)
		if err == sql.ErrNoRows {
			return nil, errors.New("No task with that id")
		}
		return nil, err
	}

	return &Types.Task{Id: id}, nil

}
