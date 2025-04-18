package repository

import (
	"github.com/Raghulds/Go_REST_API_MUX/entity"
)

type TaskRepository interface {
	GetTasks() ([]*entity.Task, error)
	CreateTask(task *entity.Task) (bool, error)
}
