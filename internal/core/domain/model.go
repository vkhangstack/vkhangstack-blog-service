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
	bun.BaseModel       `bun:"table:accounts"`
	ID                  string     `bun:"id,pk"`
	Username            string     `bun:"username,notnull"`
	Password            string     `bun:"password,notnull"`
	Email               *string    `bun:"email"`
	FullName            string     `bun:"full_name,notnull"`
	Role                string     `bun:"role,notnull"`
	IsActive            bool       `bun:"is_active,notnull,default:false"`
	LastLogin           *time.Time `bun:"last_login"`
	BlockedAt           *time.Time `bun:"blocked_at"`
	FailedLoginAttempts int        `bun:"failed_login_attempts,notnull,default:0"`
	CreatedAt           time.Time  `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt           time.Time  `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt           time.Time  `bun:",soft_delete,nullzero"`
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
	ID            uint64     `bun:"id,pk"`
	UserID        string     `bun:"user_id,notnull"`
	Jti           string     `bun:"jti,notnull"`
	IsRevoked     bool       `bun:"is_revoked,notnull,default:false"`
	CreatedAt     time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	ExpiresAt     *time.Time `bun:"expires_at"`
}

type FirebaseInfo struct {
	bun.BaseModel `bun:"table:firebase_infos"`
	ClientID      string `bun:"client_id"`
	PrivateID     string `bun:"private_id"`
}

type Tag struct {
	bun.BaseModel `bun:"table:blog_tags,alias:t"`
	ID            string    `bun:"id,pk,type:varchar(20)"              json:"id"`
	Name          string    `bun:"name,notnull,type:varchar(100)" json:"name"`
	Slug          string    `bun:"slug,notnull,type:varchar(100)" json:"slug"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at"`
}

type BlogPostTag struct {
	bun.BaseModel `bun:"table:blog_post_tags"`
	PostID        string `bun:"post_id,pk,type:varchar(20)"`
	TagID         string `bun:"tag_id,pk,type:varchar(20)"`
}

type BlogCategory struct {
	bun.BaseModel `bun:"table:blog_categories,alias:bc"`
	ID            string    `bun:"id,pk,type:varchar(20)"                              json:"id"`
	Name          string    `bun:"name,notnull,type:varchar(255)"                 json:"name"`
	Slug          string    `bun:"slug,notnull,type:varchar(255)"                 json:"slug"`
	Description   *string   `bun:"description,nullzero,type:text"                 json:"description"`
	ParentID      *string   `bun:"parent_id,nullzero,type:varchar(20)"                 json:"parent_id"`
	IsActive      bool      `bun:"is_active,notnull,default:true,type:boolean"    json:"is_active"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at"`
	DeletedAt     time.Time `bun:"deleted_at,soft_delete,nullzero,type:timestamptz"                       json:"-"`
}

type BlogPost struct {
	bun.BaseModel `bun:"table:blog_posts,alias:bp"`
	ID            string         `bun:"id,pk,type:varchar(20)"                               json:"id"`
	Title         string         `bun:"title,notnull,type:varchar(500)"                 json:"title"`
	Slug          string         `bun:"slug,notnull,type:varchar(500)"                  json:"slug"`
	Excerpt       *string        `bun:"excerpt,nullzero,type:text"                      json:"excerpt"`
	Content       string         `bun:"content,notnull,type:text"                       json:"content"`
	CoverImageURL *string        `bun:"cover_image_url,nullzero,type:varchar(1000)"     json:"cover_image_url"`
	CategoryID    *string        `bun:"category_id,nullzero,type:varchar(20)"                json:"category_id"`
	Status        PostStatus     `bun:"status,notnull,default:'draft',type:varchar(50)" json:"status"`
	PublishedAt   *time.Time     `bun:"published_at,nullzero,type:timestamptz"          json:"published_at"`
	ScheduledAt   *time.Time     `bun:"scheduled_at,nullzero,type:timestamptz"          json:"scheduled_at"`
	ViewCount     uint64         `bun:"view_count,notnull,default:0,type:bigint"        json:"view_count"`
	LexicalState  *string        `bun:"lexical_state,nullzero,type:text"                 json:"lexical_state"`
	AuthorID      string         `bun:"author_id,notnull,type:varchar(20)"                   json:"author_id"`
	Type          PostType       `bun:"type,notnull,default:'post',type:varchar(50)"    json:"type"`
	Visibility    PostVisibility `bun:"visibility,notnull,default:'public',type:varchar(50)" json:"visibility"`
	Locale        *string        `bun:"locale,nullzero,type:varchar(10)"                 json:"locale"`
	CreatedAt     time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt     time.Time      `bun:"updated_at,nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at"`
	DeletedAt     time.Time      `bun:"deleted_at,soft_delete,nullzero,type:timestamptz"                       json:"-"`
}

type BlogCategoryWithPostCount struct {
	*BlogCategory
	PostCount int `bun:"post_count" json:"post_count"`
}

type Task struct {
	bun.BaseModel `bun:"table:tasks,alias:t"`
	ID            string       `bun:"id,pk,type:varchar(20)"                              json:"id"`
	TaskID        string       `bun:"task_id,notnull,type:varchar(50),unique"             json:"task_id"`
	Title         string       `bun:"title,notnull,type:varchar(500)"                     json:"title"`
	Status        TaskStatus   `bun:"status,notnull,default:'todo',type:varchar(50)"      json:"status"`
	Label         TaskLabel    `bun:"label,notnull,type:varchar(50)"                      json:"label"`
	Priority      TaskPriority `bun:"priority,notnull,default:'medium',type:varchar(50)"  json:"priority"`
	HTML          *string      `bun:"html,nullzero,type:text"                             json:"html"`
	Lexical       *string      `bun:"lexical,nullzero,type:text"                          json:"lexical"`
	Description   *string      `bun:"description,nullzero,type:text"                      json:"description"`
	DueAt         *time.Time   `bun:"due_at,nullzero,type:timestamptz"                    json:"due_at"`
	AssigneeID    *string      `bun:"assignee_id,nullzero,type:varchar(36)"               json:"assignee_id"`
	EnableNotice  *bool        `bun:"enable_notice,nullzero,type:bool,default:false"      json:"enable_notice"`
	ReminderAt    *time.Time   `bun:"reminder_at,nullzero,type:timestamptz"               json:"reminder_at"`
	TypeReminder  *int         `bun:"type_reminder,nullzero,type:int,default:0"           json:"type_reminder"`
	CreatedAt     time.Time    `bun:"created_at,nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt     time.Time    `bun:"updated_at,nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at"`
	DeletedAt     time.Time    `bun:"deleted_at,soft_delete,nullzero,type:timestamptz"                       json:"-"`
}
