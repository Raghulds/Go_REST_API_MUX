package main

import (
	"context"

	"github.com/Raghulds/Go_REST_API_MUX/controller"
	"github.com/Raghulds/Go_REST_API_MUX/helpers"
	"github.com/Raghulds/Go_REST_API_MUX/repository"
	"github.com/Raghulds/Go_REST_API_MUX/router"
	"github.com/Raghulds/Go_REST_API_MUX/service"
)

var (
	ctx            = context.Background()
	firebaseClient = helpers.ConnectToFirebaseAndGetClient(ctx)

	taskRepo          = repository.NewTaskRepository(firebaseClient)
	subTaskRepo       = repository.NewSubTaskRepository(firebaseClient)
	taskService       = service.NewTaskService(taskRepo)
	subTaskService    = service.NewSubTaskService(subTaskRepo)
	taskController    = controller.NewTaskController(taskService)
	subTaskController = controller.NewSubTaskController(subTaskService)
)

func main() {
	router := router.NewMuxRouter()

	// Ping route
	router.GET("/ping", taskController.Ping)

	// Task routes
	router.GET("/tasks", taskController.GetTasks)
	router.POST("/tasks", taskController.CreateTask)

	// Sub task routes
	router.GET("/subtasks", subTaskController.GetSubTasks)
	router.POST("/subtasks", subTaskController.CreateSubTask)

	router.SERVE(":8080")
}
