package services

import (
	"context"
	"encoding/json"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var FirebaseAuth *auth.Client

func InitFirebase() {
	creds := map[string]string{
		"type":         "service_account",
		"project_id":   os.Getenv("FIREBASE_PROJECT_ID"),
		"private_key":  os.Getenv("FIREBASE_PRIVATE_KEY"),
		"client_email": os.Getenv("FIREBASE_CLIENT_EMAIL"),
		"token_uri":    "https://oauth2.googleapis.com/token",
	}
	credsJSON, err := json.Marshal(creds)
	if err != nil {
		log.Fatalf("Error al preparar las credenciales: %v", err)
	}

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(credsJSON))
	if err != nil {
		log.Fatalf("Error al iniciar Firebase: %v", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error al obtener Firebase Auth: %v", err)
	}

	FirebaseAuth = authClient
	log.Println("✅ Firebase inicializado correctamente.")
}
