package repository

import (
	"github.com/Raghulds/Go_REST_API_MUX/entity"
)

type SubTaskRepository interface {
	GetSubTasks() ([]*entity.SubTask, error)
	CreateSubTask(subTask *entity.SubTask) (bool, error)
}
