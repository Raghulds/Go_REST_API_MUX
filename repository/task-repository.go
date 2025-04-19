package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/Raghulds/Go_REST_API_MUX/entity"
	"google.golang.org/api/iterator"
)

type TaskRepository interface {
	GetTasks() ([]*entity.Task, error)
	CreateTask(task *entity.Task) (*firestore.DocumentRef, error)
	GetTaskById(id string) (*entity.Task, error)
}

const (
	taskCollectionName = "tasks"
)

type TaskRepositoryImpl struct {
	firestoreClientTask *firestore.Client
}

func NewTaskRepository(client *firestore.Client) TaskRepository {
	return &TaskRepositoryImpl{
		firestoreClientTask: client,
	}
}

func transformTask(doc *firestore.DocumentSnapshot) *entity.Task {
	var task entity.Task
	task.Id = doc.Ref.ID

	if name, ok := doc.Data()["name"]; ok && name != nil {
		if n, ok := name.(string); ok {
			task.Name = n
		}
	}

	if completed, ok := doc.Data()["completed"]; ok && completed != nil {
		if c, ok := completed.(bool); ok {
			task.Completed = c
		}
	}

	return &task
}

// TaskRepository
func (t *TaskRepositoryImpl) GetTasks() ([]*entity.Task, error) {
	ctx := context.Background()

	var tasks []*entity.Task
	iter := t.firestoreClientTask.Collection(taskCollectionName).Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			break
		}

		task := transformTask(doc)

		if subtasksRefs, exists := doc.Data()["subtasks"]; exists && subtasksRefs != nil {
			if subTasks, ok := subtasksRefs.([]interface{}); ok {
				stSlice := make([]*entity.SubTask, 0, len(subTasks))
				for _, st := range subTasks {
					if docRef, ok := st.(*firestore.DocumentRef); ok {
						subtask := &entity.SubTask{}
						stMap, err := docRef.Get(ctx)
						if err != nil {
							log.Fatalf("Failed to get document: %v", err)
						}

						stMapData := stMap.Data()
						subtask.Id = docRef.ID
						if name, ok := stMapData["name"].(string); ok {
							subtask.Name = name
						}
						if completed, ok := stMapData["completed"].(bool); ok {
							subtask.Completed = completed
						}
						stSlice = append(stSlice, subtask)
					}
				}
				task.Subtasks = stSlice
			}
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t *TaskRepositoryImpl) CreateTask(Task *entity.Task) (*firestore.DocumentRef, error) {
	ctx := context.Background()
	var task map[string]interface{} = make(map[string]interface{})
	task["completed"] = Task.Completed
	task["name"] = Task.Name

	fmt.Printf("Create - %+v", task)
	created, _, err := t.firestoreClientTask.Collection(taskCollectionName).Add(ctx, task)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (st *TaskRepositoryImpl) GetTaskById(id string) (*entity.Task, error) {
	ctx := context.Background()
	doc, err := st.firestoreClientTask.Collection(taskCollectionName).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, errors.New("subtask not found")
	}

	task := transformTask(doc)

	return task, nil
}
