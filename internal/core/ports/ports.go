package ports

import (
	"context"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

type AuthService interface {
	ValidateToken(authHeader string, jwtSecret string) (string, string, error)
	GenerateAuthTokens(userID string) (*domain.LoginResponse, error)
}

// AuthorizationService evaluates RBAC policy decisions and manages per-user permissions.
type AuthorizationService interface {
	IsAllowed(ctx context.Context, input domain.AuthzInput) (domain.AuthzResult, error)
	// SyncUserRole ensures g(userID, role) exists so the user inherits role policies.
	SyncUserRole(userID, role string) error
	// GrantPermission adds a direct user-level policy: p(userID, resource, action).
	GrantPermission(userID, resource, action string) error
	// RevokePermission removes a direct user-level policy.
	RevokePermission(userID, resource, action string) error
	// GetUserPermissions returns all direct policies for a user (not inherited from role).
	GetUserPermissions(userID string) [][]string
}

// MenuService builds the navigation menu filtered by role permissions.
type MenuService interface {
	GetMenu(ctx context.Context, role string) (*domain.MenuResponse, error)
}

type MenuRepository interface {
	CreateMenu(ctx context.Context, entry domain.MenuEntry) (*domain.MenuEntry, error)
	GetMenuByID(ctx context.Context, id string) (*domain.MenuEntry, error)
	UpdateMenu(ctx context.Context, id string, updates domain.MenuEntry) error
	DeleteMenu(ctx context.Context, id string) error
	ListMenus(ctx context.Context) ([]*domain.MenuEntry, error)
}

type MenuAdminService interface {
	CreateMenu(ctx context.Context, req domain.CreateMenuRequest) (*domain.MenuEntry, error)
	GetMenu(ctx context.Context, id string) (*domain.MenuEntry, error)
	UpdateMenu(ctx context.Context, id string, req domain.UpdateMenuRequest) error
	DeleteMenu(ctx context.Context, id string) error
	ListMenus(ctx context.Context) ([]*domain.MenuEntry, error)
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
	GetRoleByUserID(userID string) (string, error)
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
	CreateTag(ctx context.Context, tag domain.Tag) (*domain.Tag, error)
	GetTagBySlug(ctx context.Context, slug string) (*domain.Tag, error)
	ListTags(ctx context.Context, authorID string, tagType string) ([]*domain.Tag, error)
	AttachTags(ctx context.Context, postID string, tagIDs []string) error
	DetachTags(ctx context.Context, postID string) error
	GetTagByID(ctx context.Context, id string) (*domain.Tag, error)
	DeleteTag(ctx context.Context, id string) error
}

type TagService interface {
	CreateTag(ctx context.Context, authorID string, req domain.CreateTagRequest) (*domain.Tag, error)
	ListTags(ctx context.Context, authorID string, tagType string) ([]*domain.Tag, error)
	GetTagByID(ctx context.Context, id string) (*domain.Tag, error)
	DeleteTag(ctx context.Context, id string) error
}

type BlogCategoryRepository interface {
	CreateCategory(ctx context.Context, category domain.BlogCategory) (*domain.BlogCategory, error)
	GetCategory(ctx context.Context, id string) (*domain.BlogCategory, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*domain.BlogCategory, error)
	ListCategories(ctx context.Context) ([]*domain.BlogCategoryWithPostCount, error)
	ListCategoriesCursor(ctx context.Context, cursor string, limit int) ([]*domain.BlogCategoryWithPostCount, *string, error)
	UpdateCategory(ctx context.Context, category domain.BlogCategory) (*domain.BlogCategory, error)
	DeleteCategory(ctx context.Context, id string) error
}

type BlogPostRepository interface {
	CreatePost(ctx context.Context, post domain.BlogPost, tagIDs []string) (*domain.BlogPost, error)
	GetPost(ctx context.Context, id string) (*domain.BlogPost, error)
	GetPostBySlug(ctx context.Context, slug string) (*domain.BlogPost, error)
	ListPosts(ctx context.Context, filter domain.BlogPostFilter) ([]*domain.BlogPost, int, error)
	ListPostsCursor(ctx context.Context, filter domain.BlogPostFilter, cursor string, limit int) ([]*domain.BlogPost, *string, int, error)
	UpdatePost(ctx context.Context, post domain.BlogPost, tagIDs []string) error
	DeletePost(ctx context.Context, id string) error
	IncrementViewCount(id string) error
	CountPostsByCategory(ctx context.Context, categoryID string) (int, error)
}

type BlogCategoryService interface {
	CreateCategory(req domain.CreateBlogCategoryRequest) (*domain.BlogCategory, error)
	GetCategory(ctx context.Context, id string) (*domain.BlogCategoryWithPostCount, error)
	ListCategories(ctx context.Context) ([]*domain.BlogCategory, error)
	UpdateCategory(ctx context.Context, id string, req domain.UpdateBlogCategoryRequest) (*domain.BlogCategory, error)
	DeleteCategory(ctx context.Context, id string) error
}

type BlogPostService interface {
	CreatePost(authorID string, req domain.CreateBlogPostRequest) (*domain.BlogPost, error)
	GetPost(ctx context.Context, id string) (*domain.BlogPost, error)
	GetPostBySlug(ctx context.Context, slug string) (*domain.BlogPostBySlugResponse, error)
	ListPosts(ctx context.Context, filter domain.BlogPostFilter) ([]*domain.BlogPost, int, error)
	UpdatePost(ctx context.Context, id string, req domain.UpdateBlogPostRequest) error
	DeletePost(ctx context.Context, id string) error
	PublishPost(ctx context.Context, id string) error
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
	CreateTask(ctx context.Context, task domain.Task) (*domain.Task, error)
	GetTaskByID(ctx context.Context, id string) (*domain.Task, error)
	GetTaskByTaskID(ctx context.Context, taskID string) (*domain.Task, error)
	UpdateTask(ctx context.Context, id string, updates domain.Task) (*domain.Task, error)
	DeleteTask(ctx context.Context, id string) error
	ListTasks(ctx context.Context, filter domain.TaskFilter) ([]*domain.Task, int, error)
	ListTasksCursor(ctx context.Context, filter domain.TaskFilter, cursor string, limit int) ([]*domain.Task, *string, int, error)
	ListAllTasks(ctx context.Context) ([]*domain.Task, error)
	CountTasksByStatus(ctx context.Context, status domain.TaskStatus) (int, error)
	CountTasksByPriority(ctx context.Context, priority domain.TaskPriority) (int, error)
	GetCount(ctx context.Context) (int, error)
}

type TaskService interface {
	CreateTask(ctx context.Context, req domain.CreateTaskRequest) (*domain.Task, error)
	GetTask(ctx context.Context, id string) (*domain.Task, error)
	UpdateTask(ctx context.Context, id string, req domain.UpdateTaskRequest) (*domain.Task, error)
	DeleteTask(ctx context.Context, id string) error
	ListTasks(ctx context.Context, filter domain.TaskFilter) ([]*domain.Task, int, error)
	ListTasksCursor(ctx context.Context, filter domain.TaskFilter, cursor string, limit int) ([]*domain.Task, *string, int, error)
	ListAllTasks(ctx context.Context) ([]*domain.Task, error)
	GetTaskStatistics(ctx context.Context) (map[string]interface{}, error)
}
type SearchEngineRepository interface {
	IndexDocument(indexName string, document interface{}) error
	Search(indexName string, query string, limit int) ([]map[string]interface{}, error)
}

type SearchEngineService interface {
	IndexDocument(indexName string, document interface{}) error
	Search(indexName string, query string, limit int) ([]map[string]interface{}, error)
}

type NoteRepository interface {
	CreateNote(ctx context.Context, note domain.Note) (*domain.Note, error)
	GetNoteByID(ctx context.Context, id string) (*domain.Note, error)
	UpdateNote(ctx context.Context, id string, updates domain.Note) error
	DeleteNote(ctx context.Context, id string) error
	ListNotes(ctx context.Context, filter domain.NoteFilter) ([]*domain.NoteHasTag, int, error)
	ListNotesCursor(ctx context.Context, filter domain.NoteFilter, cursor string, limit int) ([]*domain.NoteHasTag, *string, int, error)
	AttachNoteTags(ctx context.Context, noteID string, tagIDs []string) error
	DetachNoteTags(ctx context.Context, noteID string) error
}

type NoteService interface {
	CreateNote(ctx context.Context, authorID string, req domain.CreateNoteRequest) (*domain.Note, error)
	GetNote(ctx context.Context, id string) (*domain.Note, error)
	ListNotes(ctx context.Context, filter domain.NoteFilter) ([]*domain.NoteHasTag, int, error)
	ListNotesCursor(ctx context.Context, filter domain.NoteFilter, cursor string, limit int) ([]*domain.NoteHasTag, *string, int, error)
	UpdateNote(ctx context.Context, id string, req domain.UpdateNoteRequest) error
	DeleteNote(ctx context.Context, id string) error
}

type DrawingRepository interface {
	CreateDrawing(ctx context.Context, drawing domain.Drawing) (*domain.Drawing, error)
	GetDrawingByID(ctx context.Context, id string) (*domain.Drawing, error)
	UpdateDrawing(ctx context.Context, id string, updates domain.Drawing) error
	DeleteDrawing(ctx context.Context, id string) error
	ListDrawings(ctx context.Context, filter domain.DrawingFilter) ([]*domain.Drawing, int, error)
	ListDrawingsCursor(ctx context.Context, filter domain.DrawingFilter, cursor string, limit int) ([]*domain.Drawing, *string, int, error)
}

type DrawingService interface {
	CreateDrawing(ctx context.Context, authorID string, req domain.CreateDrawingRequest) (*domain.Drawing, error)
	GetDrawing(ctx context.Context, id string) (*domain.DrawingResponse, error)
	ListDrawings(ctx context.Context, filter domain.DrawingFilter) ([]*domain.DrawingResponse, int, error)
	ListDrawingsCursor(ctx context.Context, filter domain.DrawingFilter, cursor string, limit int) ([]*domain.DrawingResponse, *string, int, error)
	UpdateDrawing(ctx context.Context, id string, req domain.UpdateDrawingRequest) error
	DeleteDrawing(ctx context.Context, id string) error
}

type TimetableRepository interface {
	CreateTimetableEntry(ctx context.Context, entry domain.TimetableEntry) (*domain.TimetableEntry, error)
	GetTimetableEntryByID(ctx context.Context, authorID, id string) (*domain.TimetableEntry, error)
	UpdateTimetableEntry(ctx context.Context, id string, updates domain.TimetableEntry) error
	DeleteTimetableEntry(ctx context.Context, id string) error
	ListTimetableEntries(ctx context.Context, authorID string, filter domain.TimetableEntryFilter) ([]*domain.TimetableEntry, int, error)
	ListTimetableEntriesCursor(ctx context.Context, authorID string, filter domain.TimetableEntryFilter, cursor string, limit int) ([]*domain.TimetableEntry, *string, int, error)
}

type TimetableService interface {
	CreateTimetableEntry(ctx context.Context, authorID string, req domain.CreateTimetableEntryRequest) (*domain.TimetableEntry, error)
	GetTimetableEntry(ctx context.Context, authorID, id string) (*domain.TimetableEntryResponse, error)
	ListTimetableEntries(ctx context.Context, authorID string, filter domain.TimetableEntryFilter) ([]*domain.TimetableEntryResponse, int, error)
	ListTimetableEntriesCursor(ctx context.Context, authorID string, filter domain.TimetableEntryFilter, cursor string, limit int) ([]*domain.TimetableEntryResponse, *string, int, error)
	UpdateTimetableEntry(ctx context.Context, id string, req domain.UpdateTimetableEntryRequest) error
	DeleteTimetableEntry(ctx context.Context, id string) error
}
