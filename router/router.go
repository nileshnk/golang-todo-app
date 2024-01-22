package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	AuthController "github.com/nileshnk/golang-todo-app/controllers/auth_controller"
	TodoController "github.com/nileshnk/golang-todo-app/controllers/todo_controller"
)

func MainRouter(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir("public"))
		fs.ServeHTTP(w, r);
	})

	r.Route("/api", TodoApiRoutes);
	r.Route("/auth", AuthApiRoutes);
}

func TodoApiRoutes(r chi.Router) {

	r.Get("/tasks", TodoController.GetAllTasks)

	r.Post("/tasks", TodoController.AddTask)
	
	r.Patch("/tasks/{id}", TodoController.ChangeTaskStatus)

	r.Put("/tasks/{id}", TodoController.EditTaskText)

	r.Delete("/tasks/{id}", TodoController.DeleteTask)
}

func AuthApiRoutes(r chi.Router) {
	r.Post("/signup", AuthController.SignUp)
	r.Post("/signin", AuthController.SignIn)
}