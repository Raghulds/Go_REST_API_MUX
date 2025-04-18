package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Raghulds/Go_REST_API_MUX/entity"
	"github.com/Raghulds/Go_REST_API_MUX/service"
)

var taskSvc service.TaskService

// Controller interface
type TaskController interface {
	Ping(w http.ResponseWriter, r *http.Request)
	GetTasks(w http.ResponseWriter, r *http.Request)
	CreateTask(w http.ResponseWriter, r *http.Request)
}

type TaskControllerImpl struct{}

func NewTaskController(taskService service.TaskService) TaskController {
	taskSvc = taskService
	return &TaskControllerImpl{}
}

func (t *TaskControllerImpl) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pong"))
}

func (t *TaskControllerImpl) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tasks, err := taskSvc.GetTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (t *TaskControllerImpl) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newTask entity.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	validationError := taskSvc.Validate(&newTask)
	if validationError != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validationError.Error()))
		return
	}

	created, err := taskSvc.CreateTask(&newTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}
