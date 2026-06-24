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
