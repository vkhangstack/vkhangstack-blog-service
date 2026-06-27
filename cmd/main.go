package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/cache"
	storage "github.com/vkhangstack/hexagonal-architecture/internal/adapters/objectStorage"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/repository"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/snowflake"
	"github.com/vkhangstack/hexagonal-architecture/internal/config"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
	"github.com/vkhangstack/hexagonal-architecture/internal/migrations"
)

var (
	msgService          *services.MessengerService
	customerService     *services.CustomerService
	accountService      *services.AccountService
	firebaseService     *services.FirebaseService
	blogCategoryService *services.BlogCategoryService
	blogPostService     *services.BlogPostService
	tagService          *services.TagService
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	cfg := config.LoadConfig()
	if cfg.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	logger.SetupLogger()
	logger.Log.Info("Application starting")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(25)
	sqldb.SetConnMaxLifetime(5 * time.Minute)
	sqldb.SetConnMaxIdleTime(2 * time.Minute)
	db := bun.NewDB(sqldb, pgdialect.New())
	db.RegisterModel((*domain.BlogPostTag)(nil))

	redisCache, err := cache.NewRedisCache(cfg.Cache.Host, cfg.Cache.Password)
	if err != nil {
		logger.Log.WithError(err).Fatal("failed to connect to Redis")
	}

	ctx := context.Background()
	migrator := migrate.NewMigrator(db, migrations.Migrations)
	if err := migrator.Init(ctx); err != nil {
		panic(fmt.Sprintf("migration init: %v", err))
	}
	if _, err := migrator.Migrate(ctx); err != nil {
		panic(fmt.Sprintf("migration run: %v", err))
	}
	node, err := strconv.ParseInt(cfg.App.Node, 10, 64)
	if err != nil {
		logger.Log.WithError(err).Fatal("failed to parse node ID")
	}
	snowflakeNode := snowflake.NewNode(node)

	store := repository.NewDB(db, redisCache, nil, snowflakeNode)

	storageAdapter, err := storage.NewS3Adapter(ctx, storage.S3Config{
		AccessKeyID:     cfg.S3.AccessKey,
		SecretAccessKey: cfg.S3.SecretKey,
		Endpoint:        cfg.S3.Endpoint,
		PublicURL:       cfg.S3.PublicURL,
		Bucket:          cfg.S3.Bucket,
		UsePathStyle:    true,
	})
	if err != nil {
		logger.Log.WithError(err).Fatal("failed to initialize S3 adapter")
	}

	msgService = services.NewMessengerService(store)
	customerService = services.NewCustomerService(store)
	firebaseService = services.NewFirebaseService(store)
	accountService = services.NewAccountService(store)
	blogCategoryService = services.NewBlogCategoryService(store)
	blogPostService = services.NewBlogPostService(store)
	tagService = services.NewTagService(store)
	uploadService := services.NewUploadService(storageAdapter)
	rateLimiter := services.NewRateLimiter(10, 5) // Burst 10, refill 5 req/s

	accountService.CreateAccountRoot()

	InitRoutes(msgService, customerService, accountService, firebaseService, blogCategoryService, blogPostService, tagService, uploadService, rateLimiter)
}
