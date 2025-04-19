package repository

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/Raghulds/Go_REST_API_MUX/entity"
	"google.golang.org/api/iterator"
)

type SubTaskRepository interface {
	GetSubTasks() ([]*entity.SubTask, error)
	CreateSubTask(subTask *entity.SubTask) (*firestore.DocumentRef, error)
}

const (
	subTaskCollectionName = "subtasks"
)

type SubTaskRepositoryImpl struct {
	firestoreClientSubTask *firestore.Client
}

func NewSubTaskRepository(client *firestore.Client) SubTaskRepository {
	return &SubTaskRepositoryImpl{
		firestoreClientSubTask: client,
	}
}

func transformSubTask(doc *firestore.DocumentSnapshot) *entity.SubTask {
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

	return &task
}

// SubTaskRepository
func (st *SubTaskRepositoryImpl) GetSubTasks() ([]*entity.SubTask, error) {
	ctx := context.Background()

	var subtasks []*entity.SubTask
	iter := st.firestoreClientSubTask.Collection(subTaskCollectionName).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("Error while iterating the subtasks")
			return subtasks, err
		}

		var task = transformSubTask(doc)

		subtasks = append(subtasks, task)
	}
	return subtasks, nil
}

func (st *SubTaskRepositoryImpl) CreateSubTask(subTask *entity.SubTask) (*firestore.DocumentRef, error) {
	ctx := context.Background()
	var task map[string]interface{} = make(map[string]interface{})
	task["completed"] = subTask.Completed
	task["name"] = subTask.Name

	created, _, err := st.firestoreClientSubTask.Collection(subTaskCollectionName).Add(ctx, task)
	if subTask.ParentId != "" {
		parentRef := st.firestoreClientSubTask.Collection("tasks").Doc(subTask.ParentId)

		_, err = parentRef.Update(ctx, []firestore.Update{
			{
				Path:  "subtasks",
				Value: firestore.ArrayUnion(created),
			},
		})
		if err != nil {
			fmt.Printf("Warning: Failed to update parent task's subtasks array: %v", err)
		}
	}
	if err != nil {
		return nil, err
	}
	return created, nil
}
