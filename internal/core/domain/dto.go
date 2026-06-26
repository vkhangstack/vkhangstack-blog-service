package domain

import "time"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Profile struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}

type LoginResponse struct {
	ID           uint64   `json:"_"`
	Email        string   `json:""`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         *Profile `json:"user"`
}

type CreateBlogCategoryRequest struct {
	Name        string  `json:"name"        binding:"required"`
	Slug        string  `json:"slug"        binding:"required"`
	Description *string `json:"description"`
	ParentID    *uint64 `json:"parent_id"`
}

type UpdateBlogCategoryRequest struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
	ParentID    *uint64 `json:"parent_id"`
	IsActive    *bool   `json:"is_active"`
}

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type CreateBlogPostRequest struct {
	Title         string     `json:"title"           binding:"required"`
	Slug          string     `json:"slug"            binding:"required"`
	Excerpt       *string    `json:"excerpt"`
	Content       string     `json:"content"         binding:"required"`
	CoverImageURL *string    `json:"cover_image_url"`
	CategoryID    *uint64    `json:"category_id,omitempty"`
	TagIDs        []uint64   `json:"tag_ids"`
	Status        PostStatus `json:"status"`
	ScheduledAt   *time.Time `json:"scheduled_at"`
}

type UpdateBlogPostRequest struct {
	Title         *string     `json:"title"`
	Slug          *string     `json:"slug"`
	Excerpt       *string     `json:"excerpt"`
	Content       *string     `json:"content"`
	CoverImageURL *string     `json:"cover_image_url"`
	CategoryID    *uint64     `json:"category_id,omitempty"`
	TagIDs        []uint64    `json:"tag_ids"`
	Status        *PostStatus `json:"status"`
	ScheduledAt   *time.Time  `json:"scheduled_at"`
}

type BlogPostFilter struct {
	Status     string  `form:"status"`
	CategoryID *uint64 `form:"category_id"`
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
