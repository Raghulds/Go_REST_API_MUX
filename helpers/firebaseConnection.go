package helpers

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

const (
	projectId = "todo-cae0c"
)

func ConnectToFirebaseAndGetClient(ctx context.Context) *firestore.Client {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	credentialsPath := "firebase-credentials.json"
	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		log.Fatalf("Firebase credentials file not found - %s", credentialsPath)
	}

	opt := option.WithCredentialsFile(credentialsPath)
	client, err := firestore.NewClient(ctx, projectId, opt)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}
