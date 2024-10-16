package routes

import (
	"kenshilabs/controllers"
	"kenshilabs/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/signup", controllers.Signup)
	app.Post("/signin", controllers.Signin)
	app.Post("/signout", controllers.Signout)

	// Task management routes (require authentication)
	taskGroup := app.Group("/tasks", middlewares.Protected)
	taskGroup.Post("/", controllers.CreateTask)
	taskGroup.Get("/", controllers.GetTasks)
	taskGroup.Get("/:id", controllers.GetTaskByID)
	taskGroup.Put("/:id", controllers.UpdateTask)
	taskGroup.Delete("/:id", controllers.DeleteTask)
}
