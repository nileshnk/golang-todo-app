package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)
var tasks []Task


func main() {
	Router := chi.NewRouter()
	Router.Route("/", mainRouter)

	// createServer := func() *http.Server {
	// 	return &http.Server{
	// 		Addr: ":8080",
	// 	}	
	// }

	// createServer().ListenAndServe()
	ServerAddress := "127.0.0.1:4000"
	http.ListenAndServe(ServerAddress, Router)

}

func mainRouter(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir("public"))
		fs.ServeHTTP(w, r);
	})

	r.Route("/api", apiRouter);
	tasks = append(tasks, Task{Id: 1, Text: "Task 1", Completed: false})
	tasks = append(tasks, Task{Id: 2, Text: "Task 2", Completed: false})

}

type Task struct {
	Id int `json:"id"`
	Text string `json:"text"`
	Completed bool `json:"completed"`
}

func apiRouter(r chi.Router) {
	r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		allTasks := getAllTasks()

		jsonData, jsonDataError := json.Marshal(allTasks);
		if jsonDataError != nil {
			panic(jsonDataError)
		}
		w.Write(jsonData)
	})

	r.Post("/tasks", func(w http.ResponseWriter, r *http.Request) {
		var task Task
		json.NewDecoder(r.Body).Decode(&task)

		addTask(task)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
		
		// fmt.Println(tasks)
	})

	r.Patch("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		fmt.Println("status changed 1")
		fmt.Println(id)
		taskId, err := strconv.Atoi(id)

		if err != nil {
			fmt.Println("error parsing id from string to int")
			fmt.Println(err);
			w.Write([]byte("not a valid id"));
			return 
		}
		fmt.Println("status changed 2")
		decoder := json.NewDecoder(r.Body)
		// requestBody := make(map[string]interface{})
		// decoder.Decode(&requestBody)
		// fmt.Println(requestBody["completed"])
		type completed struct {
			Completed bool `json:"completed"`
		}
		var taskCompleted completed

		fmt.Println("status changed")
		decoder.Decode(&taskCompleted)
		updateTaskStatus(taskId, taskCompleted.Completed)
		fmt.Println(tasks)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
		fmt.Println(tasks)
	});

	r.Put("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
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
		// fmt.Println(tasks)
	});

	r.Delete("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		taskId, err := strconv.Atoi(id)
		if err != nil {
			w.Write([]byte("not a valid id"));
		}

		deleteTask(taskId)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
		fmt.Println(tasks)
	}	);
}

// update the functions after implementing the database

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