package service

import (
	"errors"

	"github.com/Raghulds/Go_REST_API_MUX/entity"
	"github.com/Raghulds/Go_REST_API_MUX/repository"
)

type SubTaskService interface {
	Validate(subTask *entity.SubTask) error
	GetSubTasks() ([]*entity.SubTask, error)
	CreateSubTask(subTask *entity.SubTask) (bool, error)
}

var subtaskRepository repository.SubTaskRepository

type SubTaskServiceImpl struct{}

func NewSubTaskService(repo repository.SubTaskRepository) *SubTaskServiceImpl {
	subtaskRepository = repo
	return &SubTaskServiceImpl{}
}

func (t *SubTaskServiceImpl) Validate(task *entity.SubTask) error {
	if task == nil {
		return errors.New("")
	}
	if task.Name == "" {
		return errors.New("title cannot be empty")
	}
	if task.ParentId != "" {
		_, err := taskRepository.GetTaskById(task.ParentId)
		if err != nil {
			return errors.New("parent task not found")
		}
	}

	return nil
}

func (t *SubTaskServiceImpl) GetSubTasks() ([]*entity.SubTask, error) {
	return subtaskRepository.GetSubTasks()
}

func (t *SubTaskServiceImpl) CreateSubTask(task *entity.SubTask) (bool, error) {

	validationErr := t.Validate(task)
	if validationErr != nil {
		return false, validationErr
	}

	_, err := subtaskRepository.CreateSubTask(task)
	if err != nil {
		return false, err
	}
	return true, nil
}
