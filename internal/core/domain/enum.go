package domain

const (
	JwtIssuerAccess  = "golang-hexagonal-access"
	JwtIssuerRefresh = "golang-hexagonal-refresh"
)

type ErrorCode int

const (
	ErrorCodeEmailExists            = -498
	ErrorCodeEmailNotExists         = -499
	ErrorCodeFullName               = -493
	ErrorCodeRole                   = -492
	ErrorCodeTenantID               = -491
	ErrorCodePassword               = -497
	ErrorCodeTokenNotFound          = -496
	ErrorCodeUserInactive           = -495
	ErrorCodeInsufficientPermission = -494

	ErrorCodeUserStatusNotFound = -482
	ErrorCodeUserNotFound       = -481
	ErrorCodeProductNotFound    = -480
	ErrorCodeInventoryNotFound  = -479
	ErrorCodeOrderNotFound      = -478

	ErrorCodePayloadBadRequest   = -400
	ErrorCodeUnAuthorization     = -401
	ErrorCodeForbidden           = -403
	ErrorCodeInternalServerError = -500
	ErrorCodeTooManyRequests     = -429
	ErrorCodeInvalidCredentials  = -413
)

const (
	RoleRoot  = "ROOT"
	RoleAdmin = "ADMIN"
	RoleStaff = "STAFF"
	RoleUser  = "USER"
)

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusScheduled PostStatus = "scheduled"
	PostStatusArchived  PostStatus = "archived"
)

const (
	ErrorCodeBlogCategoryNotFound = -470
	ErrorCodeBlogPostNotFound     = -469
	ErrorCodeBlogSlugExists       = -468
)

type PostType string

const (
	PostTypePost     PostType = "post"
	PostTypePage     PostType = "page"
	PostTypeTemplate PostType = "template"
)

type PostVisibility string

const (
	PostVisibilityPublic  PostVisibility = "public"
	PostVisibilityPrivate PostVisibility = "private"
	PostVisibilityMembers PostVisibility = "members"
)

type FailedLoginAttemptsNumber int

const (
	FailedLoginAttemptsNumberMax          FailedLoginAttemptsNumber = 5
	FailedLoginAttemptsNumberMin          FailedLoginAttemptsNumber = 0
	FailedLoginAttemptsNumberBlockMinutes FailedLoginAttemptsNumber = 15
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

type TaskPriority string

const (
	TaskPriorityLow      TaskPriority = "low"
	TaskPriorityMedium   TaskPriority = "medium"
	TaskPriorityHigh     TaskPriority = "high"
	TaskPriorityCritical TaskPriority = "critical"
)

type TaskLabel string

const (
	TaskLabelDocumentation TaskLabel = "documentation"
	TaskLabelFeature       TaskLabel = "feature"
	TaskLabelBugFix        TaskLabel = "bug_fix"
	TaskLabelRefactor      TaskLabel = "refactor"
	TaskLabelTesting       TaskLabel = "testing"
)

type SearchEngineIndexName string

const (
	SearchEngineIndexNamePosts SearchEngineIndexName = "posts"
	SearchEngineIndexNameUsers SearchEngineIndexName = "users"
)
