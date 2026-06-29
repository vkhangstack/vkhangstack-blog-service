package main

import (
	"log"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/handler"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/http"
	"github.com/vkhangstack/hexagonal-architecture/internal/config"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
)

// InitRoutes initializes all application routes
func InitRoutes(
	msgService *services.MessengerService,
	customerService *services.CustomerService,
	accountService *services.AccountService,
	firebaseService *services.FirebaseService,
	blogCategoryService *services.BlogCategoryService,
	blogPostService *services.BlogPostService,
	tagService *services.TagService,
	taskService *services.TaskService,
	uploadService *services.UploadService,
	rateLimiter *services.RateLimiter,
	searchEngineService *services.SearchEngineService,
) {
	// Create routers
	router := gin.Default()
	// router2 := gin.Default()

	// Register profiling
	pprof.Register(router)
	// pprof.Register(router2)

	// Initialize handlers
	messageHandler := handler.NewMessageHandler(*msgService)
	customerHandler := handler.NewUserHandler(*customerService, *firebaseService)
	loginHandler := handler.NewLoginHandler(*accountService)
	blogHandler := handler.NewBlogHandler(blogCategoryService, blogPostService, searchEngineService)
	tagHandler := handler.NewTagHandler(tagService)
	taskHandler := handler.NewTaskHandler(taskService)
	uploadHandler := handler.NewUploadHandler(uploadService)

	// Setup route groups
	setupV1Routes(router, messageHandler, customerHandler, loginHandler, blogHandler, tagHandler, taskHandler, uploadHandler, rateLimiter)
	// setupV2Routes(router2, customerHandler)

	// Start servers
	startServers(router, nil)
}

// setupV1Routes configures v1 API routes
func setupV1Routes(
	router *gin.Engine,
	messageHandler *handler.MessageHandler,
	customerHandler *handler.UserHandler,
	loginHandler *handler.LoginHandler,
	blogHandler *handler.BlogHandler,
	tagHandler *handler.TagHandler,
	taskHandler *handler.TaskHandler,
	uploadHandler *handler.UploadHandler,
	rateLimiter *services.RateLimiter,
) {

	// Health check route
	router.GET("/health", http.TracingMiddleware(), handler.NewHealthHandler().HealthCheck)

	v1 := router.Group("/v1")
	{
		// Message routes
		messages := v1.Group("/messages")
		messages.Use(http.AuthenticationMiddleware())
		{
			messages.GET("/:id", messageHandler.ReadMessage)
			messages.GET("", messageHandler.ReadMessages)
			messages.POST("", messageHandler.CreateMessage)
			messages.PUT("/:id", messageHandler.UpdateMessage)
			messages.DELETE("/:id", messageHandler.DeleteMessage)
		}

		// User routes
		users := v1.Group("/customer")
		users.Use(http.AuthenticationMiddleware())
		{
			users.GET("/:id", customerHandler.ReadUser)
			users.GET("", customerHandler.ReadUsers)
			users.POST("", customerHandler.CreateUser)
			users.PUT("", customerHandler.UpdateUser)
			users.DELETE("", customerHandler.DeleteUser)
		}

		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", loginHandler.LoginAccount)
		}

		// Webhook routes
		v1.POST("/membership/webhooks", customerHandler.UpdateMembershipStatus)

		// CMS routes (authenticated)
		cms := v1.Group("/cms")
		cms.Use(http.AuthenticationMiddleware())
		{
			categories := cms.Group("/categories")
			{
				categories.POST("", blogHandler.CreateCategory)
				categories.GET("", blogHandler.ListCategories)
				categories.GET("/cursor", blogHandler.ListCategoriesCursor)
				categories.GET("/:id", blogHandler.GetCategory)
				categories.PUT("/:id", blogHandler.UpdateCategory)
				categories.DELETE("/:id", blogHandler.DeleteCategory)
			}

			posts := cms.Group("/posts")
			{
				posts.POST("", blogHandler.CreatePost)
				posts.GET("", blogHandler.ListPosts)
				posts.GET("/cursor", blogHandler.ListPostsCursor)
				posts.GET("/:id", blogHandler.GetPost)
				posts.PUT("/:id", blogHandler.UpdatePost)
				posts.DELETE("/:id", blogHandler.DeletePost)
				posts.POST("/:id/publish", blogHandler.PublishPost)
				posts.GET("/search", blogHandler.SearchPosts)
			}

			tags := cms.Group("/tags")
			{
				tags.POST("", tagHandler.CreateTag)
				tags.GET("", tagHandler.ListTags)
			}

			tasks := cms.Group("/tasks")
			{
				tasks.POST("", taskHandler.CreateTask)
				tasks.GET("", taskHandler.ListTasks)
				tasks.GET("/cursor", taskHandler.ListTasksCursor)
				tasks.GET("/statistics", taskHandler.GetTaskStatistics)
				tasks.GET("/:id", taskHandler.GetTask)
				tasks.PUT("/:id", taskHandler.UpdateTask)
				tasks.DELETE("/:id", taskHandler.DeleteTask)
			}
		}

		// Public blog routes (no auth)
		blog := v1.Group("/blog")
		{
			blog.GET("/categories", blogHandler.ListCategories)
			blog.GET("/posts", blogHandler.ListPublishedPosts)
			blog.GET("/posts/:slug", blogHandler.GetPostBySlug)
			blog.GET("/tags", tagHandler.ListTags)
			blog.GET("/search", blogHandler.SearchBlogPostsPublic)
		}
		// Upload routes
		upload := v1.Group("/upload")
		upload.Use(http.AuthenticationMiddleware())
		{
			upload.POST("", uploadHandler.UploadFile)
		}
		// Upload routes for TinyEditor
		tinyEditor := v1.Group("/tiny-editor")
		tinyEditor.Use(http.RateLimitMiddleware(rateLimiter))
		{
			tinyEditor.POST("", uploadHandler.UploadFileTinyEditor)
			tinyEditor.DELETE("", uploadHandler.DeleteFileTinyEditor)
		}
	}
}

// setupV2Routes configures v2 API routes
func setupV2Routes(
	router *gin.Engine,
	loginHandler *handler.LoginHandler,
) {
	v2 := router.Group("/v2")
	{
		// Auth routes
		v2.POST("/login", loginHandler.LoginAccount)
	}
}

// startServers starts the HTTP servers
func startServers(router *gin.Engine, router2 *gin.Engine) {
	// Start main server
	loadConfig := config.LoadConfig()
	port := loadConfig.App.Port
	if len(port) != 4 && len(port) != 5 {
		panic("Port not accept")
	}
	if port == "" {
		port = "8080"
	}
	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	// Uncomment to run multiple servers concurrently
	// go func() {
	// 	if err := router.Run(":5000"); err != nil {
	// 		log.Fatalf("failed to run messages and users service: %v", err)
	// 	}
	// }()

	// if err := router2.Run(":4242"); err != nil {
	// 	log.Fatalf("failed to run payments service: %v", err)
	// }
}
