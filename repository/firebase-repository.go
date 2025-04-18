package repository

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/Raghulds/Go_REST_API_MUX/entity"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Service account credentials
var opt = option.WithCredentialsJSON([]byte(`{
    "type": "service_account",
    "project_id": "",
    "private_key_id": "",
    "private_key": "",
    "client_email": "",
    "client_id": "",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_x509_cert_url": "",
    "universe_domain": "googleapis.com"
  }`))

const (
	projectId             = "todo-cae0c"
	taskCollectionName    = "tasks"
	subTaskCollectionName = "subtasks"
)

type TaskRepositoryImpl struct{}
type SubTaskRepositoryImpl struct{}

func NewFireStoreRepository() (TaskRepository, SubTaskRepository) {
	return &TaskRepositoryImpl{}, &SubTaskRepositoryImpl{}
}

// TaskRepository
func (t *TaskRepositoryImpl) GetTasks() ([]*entity.Task, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId, opt)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil, err
	}
	defer client.Close()

	var tasks []*entity.Task
	iter := client.Collection(taskCollectionName).Documents(ctx)
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

		task := entity.Task{}
		task.Id = doc.Ref.ID

		if name, exists := doc.Data()["name"]; exists && name != nil {
			if name, ok := name.(string); ok {
				task.Name = name
			}
		}
		if completed, exists := doc.Data()["completed"]; exists && completed != nil {
			if completed, ok := completed.(bool); ok {
				task.Completed = completed
			}
		}

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
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (t *TaskRepositoryImpl) CreateTask(Task *entity.Task) (bool, error) {
	return false, nil
}

// SubTaskRepository
func (st *SubTaskRepositoryImpl) GetSubTasks() ([]*entity.SubTask, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId, opt)
	if err != nil {
		fmt.Println("Unable to get Firebase client connection")
		return nil, err
	}
	defer client.Close()

	var subtasks []*entity.SubTask
	iter := client.Collection(subTaskCollectionName).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("Error while iterating the subtasks")
			return subtasks, err
		}

		var task entity.SubTask
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

		subtasks = append(subtasks, &task)
	}
	return subtasks, nil
}

func (st *SubTaskRepositoryImpl) CreateSubTask(subTask *entity.SubTask) (bool, error) {
	return false, nil
}
