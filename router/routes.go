package routes

import (
	"github.com/go-chi/chi/v5"
	TodoController "github.com/nileshnk/golang-todo-app/controller/todo-controller"
)

func ApiRouter(r chi.Router) {

	r.Get("/tasks", TodoController.GetAllTasks)

	r.Post("/tasks", TodoController.AddTask)
	
	r.Patch("/tasks/{id}", TodoController.ChangeTaskStatus)

	r.Put("/tasks/{id}", TodoController.EditTaskText)

	r.Delete("/tasks/{id}", TodoController.DeleteTask)
	// r.Post("/tasks", func(w http.ResponseWriter, r *http.Request) {
	// 	var task Task
	// 	json.NewDecoder(r.Body).Decode(&task)

	// 	addTask(task)

	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(task)
		
	// 	// fmt.Println(tasks)
	// })

	// r.Patch("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	id := chi.URLParam(r, "id")
	// 	fmt.Println("status changed 1")
	// 	fmt.Println(id)
	// 	taskId, err := strconv.Atoi(id)

	// 	if err != nil {
	// 		fmt.Println("error parsing id from string to int")
	// 		fmt.Println(err);
	// 		w.Write([]byte("not a valid id"));
	// 		return 
	// 	}
	// 	fmt.Println("status changed 2")
	// 	decoder := json.NewDecoder(r.Body)
	// 	// requestBody := make(map[string]interface{})
	// 	// decoder.Decode(&requestBody)
	// 	// fmt.Println(requestBody["completed"])
	// 	type completed struct {
	// 		Completed bool `json:"completed"`
	// 	}
	// 	var taskCompleted completed

	// 	fmt.Println("status changed")
	// 	decoder.Decode(&taskCompleted)
	// 	updateTaskStatus(taskId, taskCompleted.Completed)
	// 	fmt.Println(tasks)

	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(tasks)
	// 	fmt.Println(tasks)
	// });

	// r.Put("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	id := chi.URLParam(r, "id")
	// 	taskId, err := strconv.Atoi(id)
	// 	if err != nil {
	// 		w.Write([]byte("not a valid id"));
	// 		return
	// 	}

	// 	decoder := json.NewDecoder(r.Body)
		
	// 	type text struct {
	// 		Text string `json:"text"`
	// 	}
	// 	var taskText text
		
	// 	decoder.Decode(&taskText)

	// 	updateTaskText(taskId, taskText.Text)
		
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(tasks)
	// 	// fmt.Println(tasks)
	// });

	// r.Delete("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	id := chi.URLParam(r, "id")
	// 	taskId, err := strconv.Atoi(id)
	// 	if err != nil {
	// 		w.Write([]byte("not a valid id"));
	// 	}

	// 	deleteTask(taskId)

	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(tasks)
	// 	fmt.Println(tasks)
	// }	);
}
