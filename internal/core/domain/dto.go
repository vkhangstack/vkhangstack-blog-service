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
	Status   string `form:"status"`
	Label    string `form:"label"`
	Priority string `form:"priority"`
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
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
	Elements  []JSON    `json:"elements"`
	AppState  JSON      `json:"app_state"`
	Files     JSON      `json:"files"`
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
