package main

import (
	"log"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	casbinAdapter "github.com/vkhangstack/hexagonal-architecture/internal/adapters/casbin"
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
	noteService *services.NoteService,
	drawingService *services.DrawingService,
	timetableService *services.TimetableService,
	db *bun.DB,
) {
	// Initialize Casbin RBAC enforcer backed by PostgreSQL for per-user policies.
	authzAdapter, err := casbinAdapter.NewAuthorizationAdapterWithDB(db)
	if err != nil {
		log.Fatalf("failed to initialize authorization adapter: %v", err)
	}

	// Create routers
	router := gin.Default()
	// router2 := gin.Default()

	// Register profiling
	pprof.Register(router)
	// pprof.Register(router2)

	// Initialize handlers
	menuSvc := casbinAdapter.NewMenuServiceAdapter(authzAdapter)
	menuHandler := handler.NewMenuHandler(menuSvc)
	permissionHandler := handler.NewPermissionHandler(authzAdapter)
	messageHandler := handler.NewMessageHandler(*msgService)
	customerHandler := handler.NewUserHandler(*customerService, *firebaseService)
	loginHandler := handler.NewLoginHandler(*accountService)
	blogHandler := handler.NewBlogHandler(blogCategoryService, blogPostService, searchEngineService)
	tagHandler := handler.NewTagHandler(tagService)
	taskHandler := handler.NewTaskHandler(taskService, searchEngineService)
	uploadHandler := handler.NewUploadHandler(uploadService)
	noteHandler := handler.NewNoteHandler(*noteService)
	drawingHandler := handler.NewDrawingHandler(*drawingService)
	timetableHandler := handler.NewTimetableHandler(timetableService)

	// Setup route groups
	setupV1Routes(router, menuHandler, permissionHandler, messageHandler, customerHandler, loginHandler, blogHandler,
		tagHandler, taskHandler, uploadHandler, rateLimiter, noteHandler, drawingHandler, timetableHandler, authzAdapter)
	// setupV2Routes(router2, customerHandler)

	// Start servers
	startServers(router, nil)
}

// setupV1Routes configures v1 API routes
func setupV1Routes(
	router *gin.Engine,
	menuHandler *handler.MenuHandler,
	permissionHandler *handler.PermissionHandler,
	messageHandler *handler.MessageHandler,
	customerHandler *handler.UserHandler,
	loginHandler *handler.LoginHandler,
	blogHandler *handler.BlogHandler,
	tagHandler *handler.TagHandler,
	taskHandler *handler.TaskHandler,
	uploadHandler *handler.UploadHandler,
	rateLimiter *services.RateLimiter,
	noteHandler *handler.NoteHandler,
	drawingHandler *handler.DrawingHandler,
	timetableHandler *handler.TimetableHandler,
	authzAdapter *casbinAdapter.AuthorizationAdapter,
) {
	// Health check route
	router.GET("/health", http.TracingMiddleware(), handler.NewHealthHandler().HealthCheck)

	v1 := router.Group("/v1")
	{
		// Message routes
		messages := v1.Group("/messages")
		messages.Use(http.AuthenticationMiddleware(authzAdapter))
		messages.Use(http.AuthorizationMiddleware(authzAdapter, "messages"))
		{
			messages.GET("/:id", messageHandler.ReadMessage)
			messages.GET("", messageHandler.ReadMessages)
			messages.POST("", messageHandler.CreateMessage)
			messages.PUT("/:id", messageHandler.UpdateMessage)
			messages.DELETE("/:id", messageHandler.DeleteMessage)
		}

		// User routes
		users := v1.Group("/customer")
		users.Use(http.AuthenticationMiddleware(authzAdapter))
		users.Use(http.AuthorizationMiddleware(authzAdapter, "customer"))
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

		// Account routes (authenticated — no authz check, users can always see their own data)
		account := v1.Group("/account")
		account.Use(http.AuthenticationMiddleware(authzAdapter))
		{
			account.GET("/menu", menuHandler.GetMenu)
			account.GET("/permissions", permissionHandler.GetMyPermissions)
		}

		// Webhook routes
		v1.POST("/membership/webhooks", customerHandler.UpdateMembershipStatus)

		// CMS routes (authenticated + authorized per sub-group)
		cms := v1.Group("/cms")
		cms.Use(http.AuthenticationMiddleware(authzAdapter))
		{
			categories := cms.Group("/categories")
			categories.Use(http.AuthorizationMiddleware(authzAdapter, "cms/categories"))
			{
				categories.POST("", blogHandler.CreateCategory)
				categories.GET("", blogHandler.ListCategories)
				categories.GET("/cursor", blogHandler.ListCategoriesCursor)
				categories.GET("/:id", blogHandler.GetCategory)
				categories.PUT("/:id", blogHandler.UpdateCategory)
				categories.DELETE("/:id", blogHandler.DeleteCategory)
			}

			posts := cms.Group("/posts")
			posts.Use(http.AuthorizationMiddleware(authzAdapter, "cms/posts"))
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
			tags.Use(http.AuthorizationMiddleware(authzAdapter, "cms/tags"))
			{
				tags.POST("", tagHandler.CreateTag)
				tags.GET("", tagHandler.ListTags)
			}

			tasks := cms.Group("/tasks")
			tasks.Use(http.AuthorizationMiddleware(authzAdapter, "cms/tasks"))
			{
				tasks.POST("", taskHandler.CreateTask)
				tasks.GET("", taskHandler.ListTasks)
				tasks.GET("/cursor", taskHandler.ListTasksCursor)
				tasks.GET("/statistics", taskHandler.GetTaskStatistics)
				tasks.GET("/:id", taskHandler.GetTask)
				tasks.PUT("/:id", taskHandler.UpdateTask)
				tasks.DELETE("/:id", taskHandler.DeleteTask)
				tasks.GET("/search", taskHandler.SearchTasks)
			}

			notes := cms.Group("/notes")
			notes.Use(http.AuthorizationMiddleware(authzAdapter, "cms/notes"))
			{
				notes.POST("", noteHandler.CreateNote)
				notes.GET("", noteHandler.ListNotes)
				notes.GET("/cursor", noteHandler.ListNotesCursor)
				notes.GET("/:id", noteHandler.GetNote)
				notes.PUT("/:id", noteHandler.UpdateNote)
				notes.DELETE("/:id", noteHandler.DeleteNote)
			}

			drawings := cms.Group("/drawings")
			drawings.Use(http.AuthorizationMiddleware(authzAdapter, "cms/drawings"))
			{
				drawings.POST("", drawingHandler.CreateDrawing)
				drawings.GET("", drawingHandler.ListDrawings)
				drawings.GET("/cursor", drawingHandler.ListDrawingsCursor)
				drawings.GET("/:id", drawingHandler.GetDrawing)
				drawings.PUT("/:id", drawingHandler.UpdateDrawing)
				drawings.DELETE("/:id", drawingHandler.DeleteDrawing)
			}
			// Permission management — ROOT/ADMIN only (reuse customer resource guard)
			permissions := cms.Group("/permissions")
			permissions.Use(http.AuthorizationMiddleware(authzAdapter, "cms/menus"))
			{
				permissions.POST("/grant", permissionHandler.GrantPermission)
				permissions.POST("/revoke", permissionHandler.RevokePermission)
			}

			timetables := cms.Group("/timetables")
			timetables.Use(http.AuthorizationMiddleware(authzAdapter, "cms/timetables"))
			{
				timetables.POST("", timetableHandler.CreateTimetableEntry)
				timetables.GET("", timetableHandler.ListTimetableEntries)
				timetables.GET("/cursor", timetableHandler.ListTimetableEntriesCursor)
				timetables.GET("/:id", timetableHandler.GetTimetableEntry)
				timetables.PUT("/:id", timetableHandler.UpdateTimetableEntry)
				timetables.DELETE("/:id", timetableHandler.DeleteTimetableEntry)
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
		upload.Use(http.AuthenticationMiddleware(authzAdapter))
		upload.Use(http.AuthorizationMiddleware(authzAdapter, "cms/upload"))
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
