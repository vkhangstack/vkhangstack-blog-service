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
	ErrorCodeAccountNotFound    = -404

	ErrorCodePayloadBadRequest   = -400
	ErrorCodeUnAuthorization     = -401
	ErrorCodeForbidden           = -403
	ErrorCodeInternalServerError = -500
	ErrorCodeTooManyRequests     = -429
	ErrorCodeInvalidCredentials  = -413
)

const (
	RoleRoot    = "root"
	RoleAdmin   = "admin"
	RoleStaff   = "staff"
	RoleAnalyst = "anylytics"
	RoleGuest   = "guest"
	RoleUser    = "user"
)

// KnownRoles is the canonical list of RBAC roles the system recognizes (excludes RoleGuest,
// which is a request-context default rather than an assignable role).
var KnownRoles = []string{RoleRoot, RoleAdmin, RoleStaff, RoleAnalyst, RoleUser}

// IsKnownRole reports whether role is one of KnownRoles.
func IsKnownRole(role string) bool {
	for _, r := range KnownRoles {
		if r == role {
			return true
		}
	}
	return false
}

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

type AccountStatus string

const (
	AccountStatusActive    AccountStatus = "active"
	AccountStatusInactive  AccountStatus = "inactive"
	AccountStatusInvited   AccountStatus = "invited"
	AccountStatusSuspended AccountStatus = "suspended"
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
	SearchEngineIndexNameTasks SearchEngineIndexName = "tasks"
	SearchEngineIndexNamePages SearchEngineIndexName = "pages"
	SearchEngineIndexNameNotes SearchEngineIndexName = "notes"
)

type TagType string

const (
	TagTypePost TagType = "post"
	TagTypeNote TagType = "note"
	TagTypeTask TagType = "task"
	TagTypePage TagType = "page"
	TagTypeUser TagType = "user"
)
