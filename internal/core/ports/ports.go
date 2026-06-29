package ports

import (
	"context"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

type AuthService interface {
	ValidateToken(authHeader string, jwtSecret string) (string, error)
	GenerateAuthTokens(userID string) (*domain.LoginResponse, error)
}

type MessengerService interface {
	CreateMessage(userID string, message domain.Message) error
	ReadMessage(id string) (*domain.Message, error)
	ReadMessages() ([]*domain.Message, error)
	UpdateMessage(id string, message domain.Message) error
	DeleteMessage(id string) error
}

type MessengerRepository interface {
	CreateMessage(userID string, message domain.Message) error
	ReadMessage(id string) (*domain.Message, error)
	ReadMessages() ([]*domain.Message, error)
	UpdateMessage(id string, message domain.Message) error
	DeleteMessage(id string) error
}

type CustomerService interface {
	CreateUser(email, password string) (*domain.Customer, error)
	ReadUser(id uint64) (*domain.Customer, error)
	ReadUsers() ([]*domain.Customer, error)
	UpdateUser(id, email, password string) error
	DeleteUser(id uint64) error
	LoginUser(email, password string) (*domain.LoginResponse, error)
	UpdateMembershipStatus(id uint64, status bool) error
}

type CustomerRepository interface {
	CreateUser(email, password string) (*domain.Customer, error)
	ReadUser(id uint64) (*domain.Customer, error)
	ReadUsers() ([]*domain.Customer, error)
	UpdateUser(id, email, password string) error
	DeleteUser(id uint64) error
	LoginUser(email, password string) (*domain.LoginResponse, error)
	UpdateMembershipStatus(id uint64, status bool) error
}

type AccountRepository interface {
	CreateAccount(account domain.Account) (*domain.Account, error)
	FindAccountByUsername(username string) (*domain.Account, error)
	LoginAccount(username, password string) (*string, error)
	ProfileAccount(userID string) (*domain.Account, error)
	CheckAccountExists(username string) (bool, error)
	CheckAccountIsBlocked(username string) (bool, error)
	CheckAccountTemporarilyBlocked(username string) (bool, error)
	SetAccountTemporarilyBlocked(username string, duration time.Duration) error
	SetAccountBlocked(username string, blocked bool) error
	IncrementFailedLoginAttempts(username string) error
	ResetFailedLoginAttempts(username string) error
}

type AccountService interface {
	CreateAccountRoot() error
	LoginAccount(username, password string) (*domain.LoginResponse, error)
	ProfileAccount(userID string) (*domain.Account, error)
}

type PaymentService interface {
	CreateCheckoutSession(userID string, payment domain.Payment) error
	// ProcessPaymentWithStripe(userID string, payment domain.Payment) error
}

type PaymentRepository interface {
	CreateCheckoutSession(userID string, payment domain.Payment) error
	// ProcessPaymentWithStripe(userID string, payment domain.Payment) error
}

type FirebaseRepository interface {
	// InitializeAppWithServiceAccount() *firebase.App
	GetUser(ctx context.Context, id string) *auth.UserRecord
	ListUsers(ctx context.Context)
}

type TagRepository interface {
	CreateTag(tag domain.Tag) (*domain.Tag, error)
	GetTagBySlug(slug string) (*domain.Tag, error)
	ListTags() ([]*domain.Tag, error)
	AttachTags(postID string, tagIDs []string) error
	DetachTags(postID string) error
}

type TagService interface {
	CreateTag(req domain.CreateTagRequest) (*domain.Tag, error)
	ListTags() ([]*domain.Tag, error)
}

type BlogCategoryRepository interface {
	CreateCategory(category domain.BlogCategory) (*domain.BlogCategory, error)
	GetCategory(id string) (*domain.BlogCategory, error)
	GetCategoryBySlug(slug string) (*domain.BlogCategory, error)
	ListCategories() ([]*domain.BlogCategoryWithPostCount, error)
	ListCategoriesCursor(cursor string, limit int) ([]*domain.BlogCategoryWithPostCount, *string, error)
	UpdateCategory(category domain.BlogCategory) (*domain.BlogCategory, error)
	DeleteCategory(id string) error
}

type BlogPostRepository interface {
	CreatePost(post domain.BlogPost, tagIDs []string) (*domain.BlogPost, error)
	GetPost(id string) (*domain.BlogPost, error)
	GetPostBySlug(slug string) (*domain.BlogPost, error)
	ListPosts(filter domain.BlogPostFilter) ([]*domain.BlogPost, int, error)
	ListPostsCursor(filter domain.BlogPostFilter, cursor string, limit int) ([]*domain.BlogPost, *string, int, error)
	UpdatePost(post domain.BlogPost, tagIDs []string) error
	DeletePost(id string) error
	IncrementViewCount(id string) error
	CountPostsByCategory(categoryID string) (int, error)
}

type BlogCategoryService interface {
	CreateCategory(req domain.CreateBlogCategoryRequest) (*domain.BlogCategory, error)
	GetCategory(id string) (*domain.BlogCategoryWithPostCount, error)
	ListCategories() ([]*domain.BlogCategory, error)
	UpdateCategory(id string, req domain.UpdateBlogCategoryRequest) (*domain.BlogCategory, error)
	DeleteCategory(id string) error
}

type BlogPostService interface {
	CreatePost(authorID string, req domain.CreateBlogPostRequest) (*domain.BlogPost, error)
	GetPost(id string) (*domain.BlogPost, error)
	GetPostBySlug(slug string) (*domain.BlogPost, error)
	ListPosts(filter domain.BlogPostFilter) ([]*domain.BlogPost, int, error)
	UpdatePost(id string, req domain.UpdateBlogPostRequest) error
	DeletePost(id string) error
	PublishPost(id string) error
}

type UploadService interface {
	UploadFile(ctx context.Context, fileName string, fileData []byte) (string, error)
	UploadFileWithBucket(ctx context.Context, bucketName string, fileName string, fileData []byte) (string, error)
	DeleteFile(ctx context.Context, fileKey string) error
	DeleteFileWithBucket(ctx context.Context, bucketName string, fileKey string) error
	PublicURL(key string, bucket string) string
}

type RateLimiter interface {
	Allow(ip string) bool
}

type TaskRepository interface {
	CreateTask(task domain.Task) (*domain.Task, error)
	GetTaskByID(id string) (*domain.Task, error)
	GetTaskByTaskID(taskID string) (*domain.Task, error)
	UpdateTask(id string, updates domain.Task) (*domain.Task, error)
	DeleteTask(id string) error
	ListTasks(filter domain.TaskFilter) ([]*domain.Task, int, error)
	ListTasksCursor(filter domain.TaskFilter, cursor string, limit int) ([]*domain.Task, *string, int, error)
	ListAllTasks() ([]*domain.Task, error)
	CountTasksByStatus(status domain.TaskStatus) (int, error)
	CountTasksByPriority(priority domain.TaskPriority) (int, error)
	GetCount() (int, error)
}

type TaskService interface {
	CreateTask(req domain.CreateTaskRequest) (*domain.Task, error)
	GetTask(id string) (*domain.Task, error)
	UpdateTask(id string, req domain.UpdateTaskRequest) (*domain.Task, error)
	DeleteTask(id string) error
	ListTasks(filter domain.TaskFilter) ([]*domain.Task, int, error)
	ListTasksCursor(filter domain.TaskFilter, cursor string, limit int) ([]*domain.Task, *string, int, error)
	ListAllTasks() ([]*domain.Task, error)
	GetTaskStatistics() (map[string]interface{}, error)
}
type SearchEngineRepository interface {
	IndexDocument(indexName string, document interface{}) error
	Search(indexName string, query string, limit int) ([]map[string]interface{}, error)
}

type SearchEngineService interface {
	IndexDocument(indexName string, document interface{}) error
	Search(indexName string, query string, limit int) ([]map[string]interface{}, error)
}
