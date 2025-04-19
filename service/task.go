package service

import (
	"errors"

	"github.com/Raghulds/Go_REST_API_MUX/entity"
	"github.com/Raghulds/Go_REST_API_MUX/repository"
)

type TaskService interface {
	Validate(Task *entity.Task) error
	GetTasks() ([]*entity.Task, error)
	CreateTask(Task *entity.Task) (bool, error)
}

var taskRepository repository.TaskRepository

type TaskServiceImpl struct{}

func NewTaskService(repo repository.TaskRepository) TaskService {
	taskRepository = repo
	return &TaskServiceImpl{}
}

func (t *TaskServiceImpl) Validate(task *entity.Task) error {
	if task == nil {
		return errors.New("")
	}
	if task.Name == "" {
		return errors.New("title cannot be empty")
	}
	return nil
}

func (t *TaskServiceImpl) GetTasks() ([]*entity.Task, error) {
	return taskRepository.GetTasks()
}

func (t *TaskServiceImpl) CreateTask(task *entity.Task) (bool, error) {

	validationErr := t.Validate(task)
	if validationErr != nil {
		return false, validationErr
	}
	_, err := taskRepository.CreateTask(task)
	if err != nil {
		return false, err
	}
	return true, nil
}
