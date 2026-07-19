package domain

import (
	"time"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Profile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	ID           string   `json:"-"`
	Email        string   `json:"-"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         *Profile `json:"user"`
}

// CursorPaginationResponse represents a paginated response with cursor
type CursorPaginationResponse struct {
	Items      interface{} `json:"items"`
	NextCursor *string     `json:"next_cursor,omitempty"`
	HasMore    bool        `json:"has_more"`
}

// Pagination holds page/limit query parameters and computes the SQL offset.
type Pagination struct {
	Page  int `form:"page"  json:"page"`
	Limit int `form:"limit" json:"limit"`
}

// CursorPagination represents cursor-based pagination parameters
type CursorPagination struct {
	Cursor string `form:"cursor" json:"cursor"`
	Limit  int    `form:"limit"  json:"limit"`
}

type CreateBlogCategoryRequest struct {
	Name        string  `json:"name"        binding:"required"`
	Slug        string  `json:"slug"        binding:"required"`
	Description *string `json:"description"`
	ParentID    *string `json:"parent_id"`
}

type UpdateBlogCategoryRequest struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
	ParentID    *string `json:"parent_id"`
	IsActive    *bool   `json:"is_active"`
}

type CreateTagRequest struct {
	Name string   `json:"name" binding:"required"`
	Slug string   `json:"slug" binding:"required"`
	Type *TagType `json:"type" binding:"required"`
}

type CreateBlogPostRequest struct {
	Title         string          `json:"title"           binding:"required"`
	Slug          string          `json:"slug"            binding:"required"`
	Excerpt       *string         `json:"excerpt"`
	Content       string          `json:"content"         binding:"required"`
	CoverImageURL *string         `json:"cover_image_url"`
	CategoryID    *string         `json:"category_id,omitempty"`
	TagIDs        []string        `json:"tag_ids"`
	Status        PostStatus      `json:"status"`
	ScheduledAt   *time.Time      `json:"scheduled_at"`
	LexicalState  *string         `json:"lexical_state"`
	Type          *PostType       `json:"type,omitempty"`
	Locale        *string         `json:"locale,omitempty"`
	Visibility    *PostVisibility `json:"visibility,omitempty"`
}

type UpdateBlogPostRequest struct {
	Title         *string         `json:"title"`
	Slug          *string         `json:"slug"`
	Excerpt       *string         `json:"excerpt"`
	Content       *string         `json:"content"`
	CoverImageURL *string         `json:"cover_image_url"`
	CategoryID    *string         `json:"category_id,omitempty"`
	TagIDs        []string        `json:"tag_ids"`
	Status        *PostStatus     `json:"status"`
	ScheduledAt   *time.Time      `json:"scheduled_at"`
	LexicalState  *string         `json:"lexical_state"`
	Locale        *string         `json:"locale,omitempty"`
	Visibility    *PostVisibility `json:"visibility,omitempty"`
}

type BlogPostFilter struct {
	Status     string  `form:"status"`
	CategoryID *string `form:"category_id"`
	Tag        string  `form:"tag"`
	Page       int     `form:"page"`
	Limit      int     `form:"limit"`
}

type BlogPostListResponse struct {
	Total int         `json:"total"`
	Posts []*BlogPost `json:"posts"`
}
type BlogUserSearchResult struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
}
type BlogPostBySlugResponse struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	Slug          string         `json:"slug"`
	Excerpt       *string        `json:"excerpt"`
	Content       string         `json:"content"`
	CoverImageURL *string        `json:"cover_image_url"`
	CategoryID    *string        `json:"category_id,omitempty"`
	Status        PostStatus     `json:"status"`
	ViewCount     uint64         `json:"view_count"`
	AuthorID      string         `json:"author_id"`
	Type          PostType       `json:"type"`
	Visibility    PostVisibility `json:"visibility"`
	Locale        *string        `json:"locale,omitempty"`
}

type UploadFileResponse struct {
	FileKey string `json:"file_key"`
	URL     string `json:"url"`
}

type UploadFileResponseTinyEditor struct {
	FileKey string `json:"file_key"`
	URL     string `json:"url"`
}

type CreateTaskRequest struct {
	Title        string       `json:"title"        binding:"required"`
	Status       TaskStatus   `json:"status"`
	Label        TaskLabel    `json:"label"        binding:"required"`
	Priority     TaskPriority `json:"priority"`
	HTML         *string      `json:"html"`
	Lexical      *string      `json:"lexical"`
	Description  *string      `json:"description"`
	DueAt        *time.Time   `json:"due_at"`
	AssigneeID   *string      `json:"assignee_id"`
	EnableNotice *bool        `json:"enable_notice"`
	ReminderAt   *time.Time   `json:"reminder_at"`
	TypeReminder *int         `json:"type_reminder"`
}

type UpdateTaskRequest struct {
	Title        *string       `json:"title"`
	Status       *TaskStatus   `json:"status"`
	Label        *TaskLabel    `json:"label"`
	Priority     *TaskPriority `json:"priority"`
	HTML         *string       `json:"html"`
	Lexical      *string       `json:"lexical"`
	Description  *string       `json:"description"`
	DueAt        *time.Time    `json:"due_at"`
	AssigneeID   *string       `json:"assignee_id"`
	EnableNotice *bool         `json:"enable_notice"`
	ReminderAt   *time.Time    `json:"reminder_at"`
	TypeReminder *int          `json:"type_reminder"`
}

type TaskListResponse struct {
	Total int     `json:"total"`
	Tasks []*Task `json:"tasks"`
}

type TaskFilter struct {
	Status    string `form:"status"`
	Label     string `form:"label"`
	Priority  string `form:"priority"`
	CreatedBy string `form:"created_by"`
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
}

type CreateNoteRequest struct {
	Title       string   `json:"title" binding:"required"`
	SourceURL   []string `json:"source_url"`
	Status      string   `json:"status"`
	HTML        *string  `json:"html"`
	Lexical     *string  `json:"lexical"`
	Description *string  `json:"description"`
	TagIDs      []string `json:"tag_ids"`
}

type UpdateNoteRequest struct {
	Title       *string  `json:"title,omitempty"`
	SourceURL   []string `json:"source_url"`
	Status      *string  `json:"status"`
	HTML        *string  `json:"html"`
	Lexical     *string  `json:"lexical"`
	Description *string  `json:"description"`
	TagIDs      []string `json:"tag_ids"`
}

type NoteFilter struct {
	Status    *string `form:"status"`
	Title     *string `form:"title"`
	Page      int     `form:"page"`
	Limit     int     `form:"limit"`
	CreatedBy *string `form:"created_by"`
}

type CreateAccountRequest struct {
	FirstName   string        `json:"first_name" binding:"required"`
	LastName    string        `json:"last_name"  binding:"required"`
	Username    string        `json:"username"   binding:"required"`
	Email       *string       `json:"email"      binding:"omitempty,email"`
	PhoneNumber *string       `json:"phone_number"`
	Password    string        `json:"password"   binding:"required"`
	Role        string        `json:"role"`
	Status      AccountStatus `json:"status"`
}

type UpdateAccountRequest struct {
	FirstName   *string        `json:"first_name"`
	LastName    *string        `json:"last_name"`
	Username    *string        `json:"username"`
	Email       *string        `json:"email"`
	PhoneNumber *string        `json:"phone_number"`
	Password    *string        `json:"password"`
	Role        *string        `json:"role"`
	Status      *AccountStatus `json:"status"`
}

type InviteAccountRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role"`
	Desc  string `json:"desc"`
}

// AccountResponse is the account shape returned by the /v1/account endpoints — password omitted.
type AccountResponse struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Status      string    `json:"status"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RolePermissions is returned by GET /v1/account/roles/:role/permissions.
type RolePermissions struct {
	Role        string               `json:"role"`
	Permissions []ResourcePermission `json:"permissions"`
}

// RoleInfo is a single entry in GET /v1/account/roles.
type RoleInfo struct {
	Name string `json:"name"`
}

// UpdateRolePermissionsRequest is the body of PUT /v1/account/roles/:role/permissions.
type UpdateRolePermissionsRequest struct {
	Permissions []ResourcePermission `json:"permissions" binding:"required"`
}

type TagView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type NoteHasTag struct {
	*Note
	Tags []*TagView `json:"tags"`
}

type NoteListResponse struct {
	Total int           `json:"total"`
	Notes []*NoteHasTag `json:"notes"`
}

type CreateDrawingRequest struct {
	Title string `json:"title" binding:"required"`
}

type DrawingResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SceneElement struct {
	ID              string     `json:"id"`
	Type            string     `json:"type"`
	X               float64    `json:"x"`
	Y               float64    `json:"y"`
	Width           float64    `json:"width"`
	Height          float64    `json:"height"`
	Angle           float64    `json:"angle"`
	StrokeColor     string     `json:"stroke_color"`
	BackgroundColor string     `json:"background_color"`
	FillStyle       string     `json:"fill_style"`
	StrokeWidth     int        `json:"stroke_width"`
	StrokeStyle     string     `json:"stroke_style"`
	Roughness       int        `json:"roughness"`
	Opacity         int        `json:"opacity"`
	GroupIDs        []string   `json:"group_ids"`
	FrameID         *string    `json:"frame_id"`
	Index           string     `json:"index"`
	Roundness       *Roundness `json:"roundness"`
	Seed            int64      `json:"seed"`
	Version         int        `json:"version"`
	VersionNonce    int64      `json:"version_nonce"`
	IsDeleted       bool       `json:"is_deleted"`
	BoundElements   []any      `json:"bound_elements"`
	Updated         int64      `json:"updated"`
	Link            *string    `json:"link"`
	Locked          bool       `json:"locked"`
}

type Roundness struct {
	Type int `json:"type"`
}

type SceneAppState struct {
	ViewBackgroundColor        string    `json:"view_background_color"`
	CurrentItemStrokeColor     string    `json:"current_item_stroke_color"`
	CurrentItemBackgroundColor string    `json:"current_item_background_color"`
	CurrentItemFillStyle       string    `json:"current_item_fill_style"`
	CurrentItemStrokeWidth     int       `json:"current_item_stroke_width"`
	CurrentItemStrokeStyle     string    `json:"current_item_stroke_style"`
	CurrentItemRoughness       int       `json:"current_item_roughness"`
	CurrentItemOpacity         int       `json:"current_item_opacity"`
	CurrentItemFontFamily      int       `json:"current_item_font_family"`
	CurrentItemFontSize        int       `json:"current_item_font_size"`
	CurrentItemTextAlign       string    `json:"current_item_text_align"`
	ScrollX                    float64   `json:"scroll_x"`
	ScrollY                    float64   `json:"scroll_y"`
	Zoom                       ZoomState `json:"zoom"`
}

type ZoomState struct {
	Value float64 `json:"value"`
}

type UpdateDrawingRequest struct {
	Title    *string        `json:"title"`
	Elements []JSON         `json:"elements"`
	AppState JSON           `json:"app_state"`
	Files    map[string]any `json:"files"`
}

type DrawingFilter struct {
	Title *string `form:"title"`
	Page  int     `form:"page"`
	Limit int     `form:"limit"`
}

type ListDrawingResponse struct {
	Drawings []*DrawingResponse `json:"drawings"`
	Total    int                `json:"total"`
}

type DrawingDetailsResponse struct {
	ID        string         `json:"id"`
	Title     string         `json:"title"`
	Elements  []JSON         `json:"elements"`
	AppState  JSON           `json:"app_state"`
	Files     map[string]any `json:"files"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// Timetable Entries
type CreateTimetableEntryRequest struct {
	Subject   string  `json:"subject" binding:"required"`
	DayOfWeek int     `json:"day_of_week" binding:"required"`
	StartTime string  `json:"start_time" binding:"required"`
	EndTime   string  `json:"end_time" binding:"required"`
	Color     *string `json:"color"`
	Note      *string `json:"note"`
}

type UpdateTimetableEntryRequest struct {
	Subject   *string `json:"subject"`
	DayOfWeek *int    `json:"day_of_week"`
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
	Color     *string `json:"color"`
	Note      *string `json:"note"`
}

type TimetableEntryFilter struct {
	DayOfWeek *int `form:"day_of_week"`
	Page      int  `form:"page"`
	Limit     int  `form:"limit"`
}

type TimetableEntryResponse struct {
	ID        string     `json:"id"`
	Subject   string     `json:"subject"`
	DayOfWeek int        `json:"day_of_week"`
	StartTime string     `json:"start_time"`
	EndTime   string     `json:"end_time"`
	Color     string     `json:"color"`
	Note      *string    `json:"note"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// NotificationChannelInput is a single channel entry accepted from the client when
// updating notification settings — which channel, whether it's enabled, and the
// token/address to deliver to (push token, bot chat id, webhook URL, phone, email, ...).
type NotificationChannelInput struct {
	Channel NotificationChannelType `json:"channel" binding:"required,oneof=email sms slack_bot zalo_bot push webhook"`
	Enabled bool                    `json:"enabled"`
	Token   string                  `json:"token"`
}

// UpdateNotificationSettingRequest replaces the full set of channels for a user.
type UpdateNotificationSettingRequest struct {
	Channels []NotificationChannelInput `json:"channels" binding:"required,dive"`
}

type NotificationSettingResponse struct {
	ID        string                `json:"id"`
	UserID    string                `json:"user_id"`
	Channels  []NotificationChannel `json:"channels"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// CreateChannelVerificationRequest asks the server to generate a linking code for a
// channel that can't accept a client-supplied token directly (e.g. zalo_bot).
type CreateChannelVerificationRequest struct {
	Channel NotificationChannelType `json:"channel" binding:"required,oneof=zalo_bot"`
}

// ChannelVerificationResponse is the code/instructions the client shows the user to copy
// and send from within the target channel's own app (e.g. as a Zalo message to our bot).
type ChannelVerificationResponse struct {
	Channel   NotificationChannelType `json:"channel"`
	Code      string                  `json:"code"`
	Message   string                  `json:"message"`
	ExpiresAt time.Time               `json:"expires_at"`
}

// ChannelDeepLinkResponse is the static deep link that opens a channel's own app straight
// to a chat with our bot (e.g. the Zalo bot deep link). It isn't persisted anywhere — it's
// read directly from config — the client renders it as a QR code so the user can jump to
// the chat by scanning instead of searching for the bot manually.
type ChannelDeepLinkResponse struct {
	Channel  NotificationChannelType `json:"channel"`
	DeepLink string                  `json:"deep_link"`
}

type NotificationResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationListResponse struct {
	Total         int                     `json:"total"`
	Notifications []*NotificationResponse `json:"notifications"`
}

type NotificationFilter struct {
	Type  string `form:"type"`
	Read  *bool  `form:"read"`
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
}
type CreateNotificationRequest struct {
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

type UpdateNotificationRequest struct {
	Read *bool `json:"read"`
}
