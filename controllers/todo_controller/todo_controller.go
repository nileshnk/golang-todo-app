package todo_controller

import (
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
	allTasks, getErr := getAllTasks();
	if getErr != nil {
		fmt.Println("Error getting all tasks")
		fmt.Println(getErr)
		w.Write([]byte("Error getting all tasks"))
		return
	}
	w.Header().Set("Content-Type", "application/json")

	jsonData, jsonDataError := json.Marshal(allTasks);
	if jsonDataError != nil {
		panic(jsonDataError)
	}
	w.Write(jsonData)
}

func AddTask(w http.ResponseWriter, r *http.Request) {

	var task Types.Task
		json.NewDecoder(r.Body).Decode(&task)

		addTask(task)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
}

func ChangeTaskStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	taskId, err := strconv.Atoi(id)

	if err != nil {
		fmt.Println("error parsing id from string to int")
		fmt.Println(err);
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
	updateTaskStatus(taskId, taskCompleted.Completed)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
	fmt.Println(tasks)
}

func EditTaskText(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	taskId, err := strconv.Atoi(id)
	if err != nil {
		w.Write([]byte("not a valid id"));
		return
	}

	decoder := json.NewDecoder(r.Body)
	
	type text struct {
		Text string `json:"text"`
	}
	var taskText text
	
	decoder.Decode(&taskText)

	updateTaskText(taskId, taskText.Text)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
		taskId, err := strconv.Atoi(id)
		if err != nil {
			w.Write([]byte("not a valid id"));
		}

		deleteTask(taskId)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
		fmt.Println(tasks)
}

func getAllTasks() ([]Types.Task, error) {
	fmt.Println("Getting all tasks")

	// Ensure there is a valid database connection
	if db_controller.DBInstance == nil {
		return nil, errors.New("database connection is nil")
	}

	// Execute the query
	rows, err := db_controller.DBInstance.Query("SELECT * FROM todo;")
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	// Create a slice to hold the retrieved tasks
	var tasks []Types.Task
	// var UserID sql.NullString
	// Iterate through the result set and populate the tasks slice
	for rows.Next() {
		var task Types.Task
		err := rows.Scan(&task.Id, &task.Text, &task.Completed, &task.UserID, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		// if UserID.Valid {
		// 	task.UserID = UserID.String
		// }
		tasks = append(tasks, task)
	}

	

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}

	return tasks, nil
}

func getTask(taskId int) Types.Task {
	for i := 0; i < len(tasks); i++ {
		if tasks[i].Id == taskId {
			return tasks[i]
		}
	}
	return Types.Task{}
}

func addTask(task Types.Task) (Types.AppResponse, error) {
	fmt.Println("Adding a task")

	// Ensure there is a valid database connection
	if db_controller.DBConnectError != nil {
		return Types.AppResponse{}, errors.New("database connection is nil")
	}
	fmt.Println(task);
	// Execute the INSERT query
	result, err := db_controller.DBInstance.Exec("INSERT INTO todo (text, completed) VALUES ($1, $2)", task.Text, task.Completed)
	if err != nil {
		log.Println("Error executing query:", err)
		return Types.AppResponse{}, err
	}

	// Check the number of rows affected to ensure the task was added
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return Types.AppResponse{}, err
	}

	if rowsAffected == 0 {
		return Types.AppResponse{Success: false, Message: "Failed to add task"}, nil
	}

	return Types.AppResponse{Success: true, Message: "Task added successfully!"}, nil
}

func updateTaskText(taskId int, newTaskTest string) {
	for i := 0; i < len(tasks); i++ {
		if tasks[i].Id == taskId {
			tasks[i].Text = newTaskTest
		}
	}
	fmt.Println(tasks)
}

func updateTaskStatus(taskId int, status bool) {
	for i := 0; i < len(tasks); i++ {
		if tasks[i].Id == taskId {
			tasks[i].Completed = status
		}
	}
	fmt.Println(tasks)
}

func deleteTask(taskId int) {
	for i := 0; i < len(tasks); i++ {
		if tasks[i].Id == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...)
		}
	}
}