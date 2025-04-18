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

var opt = option.WithCredentialsJSON([]byte(`{
    "type": "service_account",
    "project_id": "todo-cae0c",
    "private_key_id": "8a154199f177cace25aa5a068ac04cc8c2fdcfe1",
    "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCxrsDPbYBamNIg\nAFN1/PTPNK/YIMNpbOzPlnzx0aoeuKWk71bGTqolrntIWHCJlBn9LjhkuqlIhCgo\nXAzegm1pH2BtB2ljFBgTM9WvEWwOlt4tIzEtSUG+ylNjUVAXA1DjIyZL5jcltDQP\nZTsQp5ccMXNyWqPHFlvt31TGYw1X99RGxGVur6Lwq3dozXjMdhBcEz2NQkW7Zj9T\n9UbzNwy27+umYhYJ/FTXHdkjnntDecdZZU8BkFpCzn97n9pGl4JdjBCcdG7aV3m9\nDUcGhLeV2I0OqcDfmS7/RfFTJw1AP06rHhMTA4fa2E5cWunQmIStQUzKnuCj9u9q\n6LUTRa3ZAgMBAAECggEABzC9Wn+Ry8jGqxccSYhtr0NwcMh9o/Ltn+Z1ArOzhNf/\n0w2A91g4XDtD/6zdjAG/bf6cd14S6u13Mw/rkOcaDqCfjjX4MrNd2IwW1AV6WBpM\nnirACxc9dv840ep42V+P2stBud1Ugrz4xN05v04q6DALnhdzLIO3/2ez5m3cfMBo\nS5zyU1fbUBfqtqCY5NFKjg3wmw4k6iCn4evza72WdeFF7HvSe4ZWmo7lv0jjPpLT\nfj5iBlDYqtzgIFHrKEX0X7THUNlhRtWaLWhWb9BeyqSDRpowp3sE3CWu27PQLpAL\nmQ6V31Jjt/g/ddG9QahAau1YLVWT1BlBKQnmBDSeEQKBgQDk2mtKZHPl4cKAySla\nBwKlz60Q9uW9jq7R60JtF/1RPO28jm8aluNAv2+TR82nhlgB+82U3GUADxtMaoF4\nnMqGRuix676AE0BJHTxE3UEjBB7PYtRq8nyYBSNsGSIUcRebyza9GB3dnSsttXNu\nJ6AOefS8GTZIHCVzIwyD72nGyQKBgQDGwnA0moK39Y/NwKKcJKTpDJAAvZ+h+zdA\nXzgMZ9bn91+tN3ckZDz2OkLesMnEC4EoHg3ZeRNFSKd15Huj5ebW7zRZsdy2TRfK\nxqSJcLlkhaFgOPt5mbI4lRd6zQAe3Vi+JXH9oZW/1AHcPpaWLbe4s62QbCBXmXyW\n9Dn2RehmkQKBgARjsBUgMhzhpKJluVZRthpKDm653ZQyLWY3VfHTuPca7RBlxvnC\nlR9DzLcNdINXD08SblIBnCpRH9vqWwteLoA+0e2/sMqyE9STK/nCKKKsTI77vUlD\n12HvD0vee4na1XIWhrk/wiri/dYFme4t8mL0sd39uc3ORGWd8XqCWAwJAoGBAJ8+\nCx3qKQ7v/B3x15ZYOZPKD9m4Exx5JEQ7xbESxPimlg42oQsUEE+KUCcQ5yZdvUYC\nBkCVo53f3uMygujGelL75SpzuQyJ6aT5z7uaB78E3U01ei1ruYFh2iT198HCv6Xg\nZFq7yjmdxzvJHWcHC+o8crOCHctxWoOq+oFYcyExAoGBAKvHI35x4zZRlPV82vSJ\n0NZUtpDVaQJM7gvmsZftbfHmPJ5SHr+vFnYFZz3YhYJpUYXUI3e/JqTSXlK8BlTq\njTSBUV2crCLuyaPqdAzIXFHL3KR1XnnqnFZH0YSRYujrrzBLv/ASZy1d15bkqvLt\nGixus7eLltlgfb13rhubc9sC\n-----END PRIVATE KEY-----\n",
    "client_email": "firebase-adminsdk-fbsvc@todo-cae0c.iam.gserviceaccount.com",
    "client_id": "109026341313450868245",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-fbsvc%40todo-cae0c.iam.gserviceaccount.com",
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
