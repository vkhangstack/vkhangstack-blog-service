package domain

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Message struct {
	bun.BaseModel `bun:"table:messages"`
	ID            string `bun:"id,pk"          json:"id"`
	UserID        string `bun:"user_id,notnull" json:"user_id"`
	Body          string `bun:"body,notnull"    json:"body"`
}

type Base struct {
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

func (b *Base) BeforeInsert(_ context.Context, _ *bun.InsertQuery) error {
	now := time.Now()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	b.UpdatedAt = now
	return nil
}

func (b *Base) BeforeUpdate(_ context.Context, _ *bun.UpdateQuery) error {
	b.UpdatedAt = time.Now()
	return nil
}

type Account struct {
	bun.BaseModel `bun:"table:accounts"`
	ID            uint64    `bun:"id,pk"`
	Username      string    `bun:"username,notnull"`
	Password      string    `bun:"password,notnull"`
	Email         *string   `bun:"email"`
	FullName      string    `bun:"full_name,notnull"`
	Role          string    `bun:"role,notnull"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     time.Time `bun:",soft_delete,nullzero"`
}

type Store struct {
	bun.BaseModel `bun:"table:stores"`
	ID            uint64    `bun:"id,pk"`
	Name          string    `bun:"name,notnull"`
	Description   string    `bun:"description"`
	Address       string    `bun:"address,notnull"`
	ImageURL      string    `bun:"image_url"`
	Phone         string    `bun:"phone,notnull"`
	Email         string    `bun:"email"`
	IsActive      bool      `bun:"is_active,notnull"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     time.Time `bun:",soft_delete,nullzero"`
}

type Customer struct {
	bun.BaseModel `bun:"table:customers"`
	ID            uint64    `bun:"id,pk"`
	Email         string    `bun:"email,notnull"`
	ExternalID    string    `bun:"external_id"`
	Password      string    `bun:"password"`
	IsActive      bool      `bun:"is_active,notnull"`
	Balance       float64   `bun:"balance"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     time.Time `bun:",soft_delete,nullzero"`
}

type Category struct {
	bun.BaseModel `bun:"table:categories"`
	ID            uint64                 `bun:"id,pk"`
	Name          string                 `bun:"name,notnull"`
	Description   string                 `bun:"description,notnull"`
	ImageURL      string                 `bun:"image_url,notnull"`
	IsActive      bool                   `bun:"is_active,notnull"`
	Subcategories map[string]interface{} `bun:"subcategories,type:jsonb"` // PostgreSQL JSONB
	CreatedAt     time.Time              `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time              `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     time.Time              `bun:",soft_delete,nullzero"`
}

type Product struct {
	bun.BaseModel `bun:"table:products"`
	ID            uint64    `bun:"id,pk"`
	Name          string    `bun:"name,notnull"`
	Description   string    `bun:"description,notnull"`
	CostPrice     float64   `bun:"cost_price,notnull"`
	SalePrice     float64   `bun:"sale_price"`
	IsActive      bool      `bun:"is_active,notnull"`
	ImageURL      string    `bun:"image_url"`
	CategoryID    uint64    `bun:"category_id,notnull"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     time.Time `bun:",soft_delete,nullzero"`
}

type Staff struct {
	bun.BaseModel    `bun:"table:staff"`
	ID               uint64    `bun:"id,pk"`
	FullName         string    `bun:"full_name,notnull"`
	StoreID          uint64    `bun:"store_id,notnull"`
	Role             string    `bun:"role,notnull"`
	Email            string    `bun:"email"`
	AuthenticationID string    `bun:"authentication_id,notnull"`
	IsActive         bool      `bun:"is_active,notnull"`
	CreatedAt        time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt        time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt        time.Time `bun:",soft_delete,nullzero"`
}

type Order struct {
	bun.BaseModel `bun:"table:orders"`
	ID            uint64    `bun:"id,pk"`
	CustomerID    uint64    `bun:"customer_id,notnull"`
	StoreID       uint64    `bun:"store_id,notnull"`
	TotalPrice    float64   `bun:"total_price,notnull"`
	Status        string    `bun:"status,notnull"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     time.Time `bun:",soft_delete,nullzero"`
}

type OrderItem struct {
	bun.BaseModel `bun:"table:order_items"`
	ID            uint64    `bun:"id,pk"`
	OrderID       uint64    `bun:"order_id,notnull"`
	ProductID     uint64    `bun:"product_id,notnull"`
	Quantity      int       `bun:"quantity,notnull"`
	UnitPrice     float64   `bun:"unit_price,notnull"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     time.Time `bun:",soft_delete,nullzero"`
}

type Inventory struct {
	bun.BaseModel `bun:"table:inventories"`
	ID            uint64    `bun:"id,pk"`
	ProductID     uint64    `bun:"product_id,notnull"`
	StoreID       uint64    `bun:"store_id,notnull"`
	Quantity      int       `bun:"quantity,notnull"`
	ReorderLevel  int       `bun:"reorder_level,notnull"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     time.Time `bun:",soft_delete,nullzero"`
}

type Payment struct {
	bun.BaseModel `bun:"table:payments"`
	ID            uint64    `bun:"id,pk"`
	CustomerID    uint64    `bun:"customer_id,notnull"`
	BuyerInfo     *Customer `bun:"rel:belongs-to,join:customer_id=id"`
	CheckoutID    string    `bun:"checkout_id,notnull"`
	OrderID       uint64    `bun:"order_id,notnull"`
}

type OrderInfo struct {
	bun.BaseModel `bun:"table:order_infos"`
	ID            uint64 `bun:"id,pk"`
	OrderID       uint64 `bun:"order_id"`
	SellerAccount string `bun:"seller_account"`
	Amount        string `bun:"amount"`
	Currency      string `bun:"currency"`
	Status        string `bun:"status"`
}

type JtiSession struct {
	bun.BaseModel `bun:"table:jti_sessions"`
	ID            uint64 `bun:"id,pk"`
	UserID        string `bun:"user_id,notnull"`
	Jti           string `bun:"jti,notnull"`
}

type FirebaseInfo struct {
	bun.BaseModel `bun:"table:firebase_infos"`
	ClientID      string `bun:"client_id"`
	PrivateID     string `bun:"private_id"`
}

type Tag struct {
	bun.BaseModel `bun:"table:tags,alias:t"`
	ID            uint64    `bun:"id,pk,type:bigint"              json:"id"`
	Name          string    `bun:"name,notnull,type:varchar(100)" json:"name"`
	Slug          string    `bun:"slug,notnull,type:varchar(100)" json:"slug"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
}

type BlogPostTag struct {
	bun.BaseModel `bun:"table:blog_post_tags"`
	PostID        uint64    `bun:"post_id,pk,type:bigint"`
	Post          *BlogPost `bun:"rel:belongs-to,join:post_id=id"`
	TagID         uint64    `bun:"tag_id,pk,type:bigint"`
	Tag           *Tag      `bun:"rel:belongs-to,join:tag_id=id"`
}

type BlogCategory struct {
	bun.BaseModel `bun:"table:blog_categories,alias:bc"`
	ID            uint64         `bun:"id,pk,type:bigint"                              json:"id"`
	Name          string         `bun:"name,notnull,type:varchar(255)"                 json:"name"`
	Slug          string         `bun:"slug,notnull,type:varchar(255)"                 json:"slug"`
	Description   *string        `bun:"description,nullzero,type:text"                 json:"description,omitempty"`
	ParentID      *uint64        `bun:"parent_id,nullzero,type:bigint"                 json:"parent_id,omitempty"`
	Parent        *BlogCategory  `bun:"rel:belongs-to,join:parent_id=id"               json:"parent,omitempty"`
	Children      []BlogCategory `bun:"rel:has-many,join:id=parent_id"                 json:"children,omitempty"`
	IsActive      bool           `bun:"is_active,notnull,default:true,type:boolean"    json:"is_active"`
	CreatedAt     time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt     time.Time      `bun:"updated_at,nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at"`
	DeletedAt     time.Time      `bun:"deleted_at,soft_delete,nullzero,type:timestamptz"                       json:"-"`
}

type BlogPost struct {
	bun.BaseModel `bun:"table:blog_posts,alias:bp"`
	ID            uint64        `bun:"id,pk,type:bigint"                               json:"id"`
	Title         string        `bun:"title,notnull,type:varchar(500)"                 json:"title"`
	Slug          string        `bun:"slug,notnull,type:varchar(500)"                  json:"slug"`
	Excerpt       *string       `bun:"excerpt,nullzero,type:text"                      json:"excerpt,omitempty"`
	Content       string        `bun:"content,notnull,type:text"                       json:"content"`
	CoverImageURL *string       `bun:"cover_image_url,nullzero,type:varchar(1000)"     json:"cover_image_url,omitempty"`
	CategoryID    *uint64       `bun:"category_id,nullzero,type:bigint"                json:"category_id,omitempty"`
	Category      *BlogCategory `bun:"rel:belongs-to,join:category_id=id"              json:"category,omitempty"`
	Tags          []*Tag        `bun:"m2m:blog_post_tags,join:Post=Tag"                json:"tags,omitempty"`
	Status        PostStatus    `bun:"status,notnull,default:'draft',type:varchar(50)" json:"status"`
	PublishedAt   *time.Time    `bun:"published_at,nullzero,type:timestamptz"          json:"published_at,omitempty"`
	ScheduledAt   *time.Time    `bun:"scheduled_at,nullzero,type:timestamptz"          json:"scheduled_at,omitempty"`
	ViewCount     uint64        `bun:"view_count,notnull,default:0,type:bigint"        json:"view_count"`
	AuthorID      uint64        `bun:"author_id,notnull,type:bigint"                   json:"author_id"`
	CreatedAt     time.Time     `bun:"created_at,nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt     time.Time     `bun:"updated_at,nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at"`
	DeletedAt     time.Time     `bun:"deleted_at,soft_delete,nullzero,type:timestamptz"                       json:"-"`
}
