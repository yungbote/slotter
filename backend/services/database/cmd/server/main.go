package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// GORM + Postgres
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// gin
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// redis
	"github.com/go-redis/redis/v8"

	// local imports
	"github.com/yungbote/slotter/backend/services/database/internal/events"
	"github.com/yungbote/slotter/backend/services/database/internal/handlers"
	"github.com/yungbote/slotter/backend/services/database/internal/middleware"
	"github.com/yungbote/slotter/backend/services/database/internal/models"
	"github.com/yungbote/slotter/backend/services/database/internal/parser"
	"github.com/yungbote/slotter/backend/services/database/internal/repos"
	"github.com/yungbote/slotter/backend/services/database/internal/server/websocket"
	"github.com/yungbote/slotter/backend/services/database/internal/services"
	"github.com/yungbote/slotter/backend/services/database/internal/services/avatar"
	"github.com/yungbote/slotter/backend/services/database/internal/services/s3"
)

func main() {
	// -------------------------------------------------------------------------
	// 1. Load environment variables
	// -------------------------------------------------------------------------
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; continuing with environment variables.")
	}

	// Required environment variables for the entire project might be:
	// - DB_DSN: Postgres DSN, e.g. "postgres://user:password@postgres:5432/dbname?sslmode=disable"
	// - REDIS_ADDR: e.g. "redis:6379"
	// - S3_BUCKET, AWS_REGION, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY (for S3)
	// - JWT_SECRET, JWT_ISSUER
	// - DICEBEAR_URL (for avatar generation) e.g. "https://api.dicebear.com"
	// - GIN_MODE: "release" or "debug"
	// (Some can be optional or have defaults)

	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		log.Fatal("Missing DB_DSN environment variable")
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("Missing REDIS_ADDR environment variable")
	}
	// S3 env checks happen in NewS3Service().

	// -------------------------------------------------------------------------
	// 2. Connect to Postgres with GORM
	// -------------------------------------------------------------------------
	gormLogger := logger.Default.LogMode(logger.Info)
	db, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}

	// Optional: auto-migrate your models:
	if err := db.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.Warehouse{},
		&models.Location{},
		&models.TransactionRecord{},
		&models.TransactionFile{},
		&models.Item{},
		&models.UserAction{},
	); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}

	// -------------------------------------------------------------------------
	// 3. Connect to Redis
	// -------------------------------------------------------------------------
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("Missing REDIS_ADDR")
	}
	// We'll use the standard redis client for your Pub/Sub:
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		// Password: "", // if needed
		// DB: 0,        // if needed
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to ping redis: %v", err)
	}
	log.Println("Connected to Redis successfully.")

	// -------------------------------------------------------------------------
	// 4. Initialize Repos
	// -------------------------------------------------------------------------
	userRepo := repos.NewURepo(db)
	companyRepo := repos.NewCRepo(db)
	warehouseRepo := repos.NewWRepo(db)
	locationRepo := repos.NewLRepo(db)
	transactionRecordRepo := repos.NewTRRepo(db)
	transactionFileRepo := repos.NewTFRepo(db)
	userActionRepo := repos.NewUARepo(db) // optional, but included for completeness
	itemRepo := repos.NewIRepo(db)

	// -------------------------------------------------------------------------
	// 5. Initialize Services
	// -------------------------------------------------------------------------
	// PubSub Publisher
	pub, err := events.NewRedisPublisher(redisAddr)
	if err != nil {
		log.Fatalf("failed to create redis publisher: %v", err)
	}

	// Token Service
	tokenSvc, err := services.NewTokenService()
	if err != nil {
		log.Fatalf("failed to init tokenService: %v", err)
	}

	refreshTokenSvc, err := services.NewRefreshTokenService(rdb)
	if err != nil {
		log.Fatalf("cannot init refreshTokenService: %v", err)
	}

	oauthSvc, err := services.NewOAuthService(..., userSvc)

	// S3 Service
	s3Svc, err := s3.NewS3Service()
	if err != nil {
		log.Fatalf("failed to init s3Service: %v", err)
	}
	avatarSvc := avatar.NewAvatarService(s3Svc)

	// Build core domain services
	userSvc := services.NewUSvc(userRepo /* pass more if needed... */)
	companySvc := services.NewCSvc(companyRepo)
	warehouseSvc := services.NewWSvc(warehouseRepo)
	locationSvc := services.NewLSvc(locationRepo)
	tfSvc := services.NewTFSvc(transactionFileRepo)
	trSvc := services.NewTRSvc(transactionRecordRepo)
	itemSvc := services.NewISvc(itemRepo)
	avatarSvc := avatar.NewAvatarService(s3Svc)

	// If you have an OAuth config for Google:
	oauthCfg := auth.OAuthConfig{
		GoogleClientID:  os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleSecret:    os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:     os.Getenv("GOOGLE_REDIRECT_URL"),
	}
	oauthService, err := services.NewOAuthService(oauthCfg, userSvc)
	if err != nil {
		log.Printf("Warning: failed to init OAuthService, Google login may not work: %v", err)
	}

	// Parser Service
	parserSvc := parser.NewParserService(locationSvc, itemSvc, trSvc, tfSvc, warehouseSvc)

	// Build the App Service
	appSvc := services.NewAppSvc(
		companySvc,
		userSvc,
		warehouseSvc,
		locationSvc,
		tfSvc,
		trSvc,
		itemSvc,
		avatarSvc,
		s3Svc,
		tokenSvc,
		refreshTokenSvc,
		oauthSvc,
		pub,
		userActionRepo,
		parserSvc,
	)


	// -------------------------------------------------------------------------
	// 6. Setup Gin + Middleware
	// -------------------------------------------------------------------------
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}
	gin.SetMode(ginMode)

	router := gin.Default()

	// If you have an AuthMiddleware that checks JWT tokens:
	authMW := middleware.AuthMiddleware(tokenSvc)

	// Create an API route group
	api := router.Group("/api")
	{
		// Optionally, apply the auth middleware to all routes except login/register,
		// or handle them individually.
		// e.g. `api.Use(authMW)` then define protected routes
		//
		// But in this example, we might do something like:
	}

	// -------------------------------------------------------------------------
	// 7. Register App Handlers
	// -------------------------------------------------------------------------
	appHandler := handlers.NewAppHandler(appSvc)
	// If we want some endpoints to be public and some protected:
	public := api.Group("/")
	{
		public.POST("/register", appHandler.RegisterUserLocal)
		public.POST("/login", appHandler.LoginUserLocal)
		public.POST("/login/google", appHandler.LoginWithGoogle)
	}
	protected := api.Group("/")
	protected.Use(authMW)
	{
		protected.POST("/logout", appHandler.LogoutUser)

		// company endpoints
		protected.POST("/company", appHandler.CreateCompany)
		protected.GET("/company/:company_id", appHandler.GetCompanyByID)
		protected.PUT("/company/:company_id/name", appHandler.UpdateCompanyName)
		protected.PUT("/company/:company_id/avatar", appHandler.UpdateCompanyAvatar)

		// warehouse endpoints
		protected.POST("/warehouse", appHandler.CreateWarehouse)
		protected.GET("/warehouse/:warehouse_id", appHandler.GetWarehouseByID)
		protected.PUT("/warehouse/:warehouse_id/name", appHandler.UpdateWarehouseName)
		protected.DELETE("/warehouse/:warehouse_id", appHandler.DeleteWarehouse)
		protected.GET("/warehouses", appHandler.ListWarehouses)

		// location endpoints
		protected.POST("/warehouse/:warehouse_id/location", appHandler.CreateLocation)
		protected.GET("/location/:location_id", appHandler.GetLocationByID)
		protected.DELETE("/location/:location_id", appHandler.DeleteLocation)
		protected.GET("/locations", appHandler.ListLocations)

		// transaction file endpoints
		protected.POST("/warehouse/:warehouse_id/transaction-file/upload", appHandler.UploadTransactionFile)
		protected.PUT("/transaction-file/:file_id/name", appHandler.UpdateTransactionFileName)
		protected.DELETE("/transaction-file/:file_id", appHandler.DeleteTransactionFile)
		protected.GET("/transaction-files", appHandler.ListTransactionFiles)

		// transaction record endpoints
		protected.POST("/warehouse/:warehouse_id/transaction-record", appHandler.CreateTransactionRecord)
		protected.GET("/transaction-record/:record_id", appHandler.GetTransactionRecordByID)
		protected.PUT("/transaction-record/:record_id/order-name", appHandler.UpdateTransactionRecordOrderName)
		protected.PUT("/transaction-record/:record_id/description", appHandler.UpdateTransactionRecordDescription)
		protected.PUT("/transaction-record/:record_id/transaction-quantity", appHandler.UpdateTransactionRecordTransactionQuantity)
		protected.PUT("/transaction-record/:record_id/completed-quantity", appHandler.UpdateTransactionRecordCompletedQuantity)
		protected.PUT("/transaction-record/:record_id/completed-date", appHandler.UpdateTransactionRecordCompletedDate)
		protected.PUT("/transaction-record/:record_id/transaction-type", appHandler.UpdateTransactionRecordTransactionType)
		protected.GET("/transaction-records", appHandler.ListTransactionRecords)

		// user endpoints
		protected.PUT("/user/:user_id/avatar", appHandler.UpdateUserAvatar)
		protected.PUT("/user/:user_id/first-name", appHandler.UpdateUserFirstName)
		protected.PUT("/user/:user_id/last-name", appHandler.UpdateUserLastName)
		protected.PUT("/user/:user_id/email", appHandler.UpdateUserEmail)
		protected.PUT("/user/:user_id/password", appHandler.UpdateUserPassword)
		protected.DELETE("/user/:user_id", appHandler.DeleteUser)
		protected.GET("/users", appHandler.ListUsers)

		// item endpoints
		protected.GET("/items", appHandler.ListItems)
	}

	// -------------------------------------------------------------------------
	// 8. WebSocket Setup (optional)
	// -------------------------------------------------------------------------
	subscriber := events.NewPubSubSubscriber(rdb)
	wsHandler := websocket.NewHandler(subscriber)
	router.GET("/ws", wsHandler.HandleWSConnection)

	// -------------------------------------------------------------------------
	// 9. Run HTTP Server
	// -------------------------------------------------------------------------
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	log.Printf("Server starting on port %s ...", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

