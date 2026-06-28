package domain

import "time"

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
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
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

type UploadFileResponse struct {
	FileKey string `json:"file_key"`
	URL     string `json:"url"`
}

type UploadFileResponseTinyEditor struct {
	FileKey string `json:"file_key"`
	URL     string `json:"url"`
}

type CreateTaskRequest struct {
	TaskID      string       `json:"task_id"      binding:"required"`
	Title       string       `json:"title"        binding:"required"`
	Status      TaskStatus   `json:"status"`
	Label       TaskLabel    `json:"label"        binding:"required"`
	Priority    TaskPriority `json:"priority"`
	HTML        *string      `json:"html"`
	Lexical     *string      `json:"lexical"`
	Description *string      `json:"description"`
}

type UpdateTaskRequest struct {
	Title       *string       `json:"title"`
	Status      *TaskStatus   `json:"status"`
	Label       *TaskLabel    `json:"label"`
	Priority    *TaskPriority `json:"priority"`
	HTML        *string       `json:"html"`
	Lexical     *string       `json:"lexical"`
	Description *string       `json:"description"`
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
