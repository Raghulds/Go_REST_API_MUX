package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Raghulds/Go_REST_API_MUX/entity"
	"github.com/Raghulds/Go_REST_API_MUX/service"
)

var subtaskSvc service.SubTaskService

type SubTaskController interface {
	GetSubTasks(w http.ResponseWriter, r *http.Request)
	CreateSubTask(w http.ResponseWriter, r *http.Request)
}

type SubTaskControllerImpl struct{}

func NewSubTaskController(service service.SubTaskService) SubTaskController {
	subtaskSvc = service
	return &SubTaskControllerImpl{}
}

func (t *SubTaskControllerImpl) GetSubTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	subtasks, err := subtaskSvc.GetSubTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subtasks)
}

func (t *SubTaskControllerImpl) CreateSubTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newTask entity.SubTask
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	created, err := subtaskSvc.CreateSubTask(&newTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}
