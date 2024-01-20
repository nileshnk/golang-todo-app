package todo_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)
type Task struct {
	Id int `json:"id"`
	Text string `json:"text"`
	Completed bool `json:"completed"`
}
var tasks []Task

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	allTasks := getAllTasks();
	
	w.Header().Set("Content-Type", "application/json")

	jsonData, jsonDataError := json.Marshal(allTasks);
	if jsonDataError != nil {
		panic(jsonDataError)
	}
	w.Write(jsonData)
}

// func AddTask(w http.ResponseWriter, r *http.Request) *http.HandlerFunc {

// }

// func ChangeTaskStatus(w http.ResponseWriter, r *http.Request) *http.HandlerFunc {

// }

// func EditTaskText(w http.ResponseWriter, r *http.Request) *http.HandlerFunc {

// }

// func DeleteTask(w http.ResponseWriter, r *http.Request) *http.HandlerFunc {

// }

func getAllTasks() []Task {
	return tasks
}

func getTask(taskId int) Task {
	for i := 0; i < len(tasks); i++ {
		if tasks[i].Id == taskId {
			return tasks[i]
		}
	}
	return Task{}
}

func addTask(task Task) {
	tasks = append(tasks, task)
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