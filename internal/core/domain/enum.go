package domain

import "errors"

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

	ErrorCodeNotificationNotFound        = -467
	ErrorCodeNotificationSettingNotFound = -466
	ErrorCodeChannelVerificationInvalid  = -465
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

type QuizExamType string

const (
	QuizExamTypeIELTS   QuizExamType = "ielts"
	QuizExamTypeTOEIC   QuizExamType = "toeic"
	QuizExamTypeGeneral QuizExamType = "general"
)

type QuizSkill string

const (
	QuizSkillListening  QuizSkill = "listening"
	QuizSkillReading    QuizSkill = "reading"
	QuizSkillGrammar    QuizSkill = "grammar"
	QuizSkillVocabulary QuizSkill = "vocabulary"
	QuizSkillMixed      QuizSkill = "mixed"
)

type QuizStatus string

const (
	QuizStatusDraft     QuizStatus = "draft"
	QuizStatusPublished QuizStatus = "published"
	QuizStatusArchived  QuizStatus = "archived"
)

type FlashcardDeckStatus string

const (
	FlashcardDeckStatusDraft     FlashcardDeckStatus = "draft"
	FlashcardDeckStatusPublished FlashcardDeckStatus = "published"
	FlashcardDeckStatusArchived  FlashcardDeckStatus = "archived"
)

// FlashcardRating is the Anki-style grade a learner gives a card during review,
// driving the SM-2 scheduling calculation.
type FlashcardRating string

const (
	FlashcardRatingAgain FlashcardRating = "again"
	FlashcardRatingHard  FlashcardRating = "hard"
	FlashcardRatingGood  FlashcardRating = "good"
	FlashcardRatingEasy  FlashcardRating = "easy"
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

// NotificationChannelType is a delivery channel a user can enable in their notification settings.
type NotificationChannelType string

const (
	NotificationChannelEmail    NotificationChannelType = "email"
	NotificationChannelSMS      NotificationChannelType = "sms"
	NotificationChannelSlackBot NotificationChannelType = "slack_bot"
	NotificationChannelZaloBot  NotificationChannelType = "zalo_bot"
	NotificationChannelPush     NotificationChannelType = "push"
	NotificationChannelWebhook  NotificationChannelType = "webhook"
)

// KnownNotificationChannels is the canonical list of channels notification settings accept.
var KnownNotificationChannels = []NotificationChannelType{
	NotificationChannelEmail,
	NotificationChannelSMS,
	NotificationChannelSlackBot,
	NotificationChannelZaloBot,
	NotificationChannelPush,
	NotificationChannelWebhook,
}

// IsKnownNotificationChannel reports whether channel is one of KnownNotificationChannels.
func IsKnownNotificationChannel(channel NotificationChannelType) bool {
	for _, c := range KnownNotificationChannels {
		if c == channel {
			return true
		}
	}
	return false
}

// VerifiableNotificationChannels are channels whose token can't be supplied directly by the
// client — instead the user copies a generated code into the channel's own app (e.g. sends it
// as a message to our Zalo bot), and the channel's webhook reports back the code plus the
// sender's platform-specific ID, which becomes the channel's token once verified.
var VerifiableNotificationChannels = []NotificationChannelType{
	NotificationChannelZaloBot,
}

// IsVerifiableNotificationChannel reports whether channel is one of VerifiableNotificationChannels.
func IsVerifiableNotificationChannel(channel NotificationChannelType) bool {
	for _, c := range VerifiableNotificationChannels {
		if c == channel {
			return true
		}
	}
	return false
}

// ChannelVerificationStatus tracks the lifecycle of a NotificationChannelVerification code.
type ChannelVerificationStatus string

const (
	ChannelVerificationPending  ChannelVerificationStatus = "pending"
	ChannelVerificationVerified ChannelVerificationStatus = "verified"
	ChannelVerificationExpired  ChannelVerificationStatus = "expired"
)

// ErrZaloBotMessage indicates a webhook update was sent by the bot itself (e.g. its own
// reply echoed back), not by an end user — callers should skip it, not treat it as an error.
var ErrZaloBotMessage = errors.New("zalo webhook: message is from the bot itself")
