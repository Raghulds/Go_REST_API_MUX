package service

import (
	"github.com/Raghulds/Go_REST_API_MUX/entity"
	"github.com/Raghulds/Go_REST_API_MUX/repository"
)

type SubTaskService interface {
	GetSubTasks() ([]*entity.SubTask, error)
	CreateSubTask(subTask *entity.SubTask) (bool, error)
}

var subtaskRepository repository.SubTaskRepository

type SubTaskServiceImpl struct{}

func NewSubTaskService(repo repository.SubTaskRepository) *SubTaskServiceImpl {
	subtaskRepository = repo
	return &SubTaskServiceImpl{}
}

func (t *SubTaskServiceImpl) GetSubTasks() ([]*entity.SubTask, error) {
	return subtaskRepository.GetSubTasks()
}

func (t *SubTaskServiceImpl) CreateSubTask(task *entity.SubTask) (bool, error) {
	return subtaskRepository.CreateSubTask(task)
}
