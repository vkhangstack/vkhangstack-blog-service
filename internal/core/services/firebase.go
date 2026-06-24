package services

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"google.golang.org/api/option"
)

type FirebaseService struct {
	repo ports.FirebaseRepository
}

func NewFirebaseService(repo ports.FirebaseRepository) *FirebaseService {
	return &FirebaseService{
		repo: repo,
	}
}

func InitializeAppWithServiceAccount() *firebase.App {
	// [START initialize_app_service_account_golang]
	path := os.Getenv("FIREBASE_PATH")

	if path == "" {
		panic("firebase admin file not found!")
	}

	opt := option.WithCredentialsFile(path)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	// [END initialize_app_service_account_golang]

	return app
}

func initializeAppWithRefreshToken() *firebase.App {
	// [START initialize_app_refresh_token_golang]
	opt := option.WithCredentialsFile("path/to/refreshToken.json")
	config := &firebase.Config{ProjectID: "my-project-id"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	// [END initialize_app_refresh_token_golang]

	return app
}

func initializeAppDefault() *firebase.App {
	// [START initialize_app_default_golang]
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	// [END initialize_app_default_golang]

	return app
}

func initializeServiceAccountID() *firebase.App {
	// [START initialize_sdk_with_service_account_id]
	conf := &firebase.Config{
		ServiceAccountID: "my-client-id@my-project-id.iam.gserviceaccount.com",
	}
	app, err := firebase.NewApp(context.Background(), conf)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	// [END initialize_sdk_with_service_account_id]
	return app
}

func accessServicesSingleApp() (*auth.Client, error) {
	// [START access_services_single_app_golang]
	// Initialize default app
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Access auth service from the default app
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	// [END access_services_single_app_golang]

	return client, err
}

func accessServicesMultipleApp() (*auth.Client, error) {
	// [START access_services_multiple_app_golang]
	// Initialize the default app
	defaultApp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Initialize another app with a different config
	opt := option.WithCredentialsFile("service-account-other.json")
	otherApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Access Auth service from default app
	defaultClient, err := defaultApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	// Access auth service from other app
	otherClient, err := otherApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	// [END access_services_multiple_app_golang]
	// Avoid unused
	_ = defaultClient
	return otherClient, nil
}

func (f *FirebaseService) GetUser(ctx context.Context, uid string) *auth.UserRecord {
	return f.repo.GetUser(ctx, uid)
}

func (f *FirebaseService) ListUsers(ctx context.Context) {
	f.repo.ListUsers(ctx)
}
