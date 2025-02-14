package api

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "gorm.io/gorm"
  "github.com/yungbote/slotter/backend/services/database/internal/repositories"
  "github.com/yungbote/slotter/backend/services/database/internal/services"
  "github.com/yungbote/slotter/backend/services/database/internal/handlers"
)

func NewRouter(db *gorm.DB) *gin.Engine {
  r := gin.Default()

  r.GET("/health", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
  })

  //Repositories
  companyRepo := repositories.NewCompanyRepository(db)
  userRepo := repositories.NewUserRepository(db)
  warehouseRepo := repositories.NewWarehouseRepository(db)
  transactionrecordRepo := repositories.NewTransactionRecordRepository(db)

  //Services
  companyService := services.NewCompanyService(companyRepo)
  userService := services.NewUserService(userRepo)
  warehouseService := services.NewWarehouseService(warehouseRepo)
  transactionrecordService := services.NewTransactionRecordService(transactionrecordRepo)

  //Handlers
  companyHandler := handlers.NewCompanyHandler(companyService)
  userHandler := handlers.NewUserHandler(userService)
  warehouseHandler := handlers.NewWarehouseHandler(warehouseService)
  transactionrecordHandler := handlers.NewTransactionRecordHandler(transactionrecordService)

  api := r.Group("/v1")
  {
    //Company
    api.POST("/company", companyHandler.CreateCompany)
    api.GET("/company/:id", companyHandler.GetCompanyByID)
    api.PUT("/company/:id", companyHandler.UpdateCompany)
    api.DELETE("/company/:id", companyHandler.DeleteCompany)

    //Warehouse
    api.POST("/warehouse", warehouseHandler.CreateWarehouse)
    api.GET("/warehouse/:id", warehouseHandler.GetWarehouseByID)
    api.PUT("/warehouse/:id", warehouseHandler.UpdateWarehouse)
    api.DELETE("/warehouse/:id", warehouseHandler.DeleteWarehouse)

    //TransactionRecord
    api.POST("/transactionrecord", transactionrecordHandler.CreateTransactionRecord)
    api.GET("/transactionrecord/:id", transactionrecordHandler.GetTransactionRecordByID)
    api.PUT("/transactionrecord/:id", transactionrecordHandler.UpdateTransactionRecord)
    api.DELETE("/transactionrecor/:id", transactionrecordHandler.DeleteTransactionRecord)

    //User
    api.POST("/user", userHandler.CreateUser)
    api.GET("/user/:id", userHandler.GetUserByID)
    api.PUT("/user/:id", userHandler.UpdateUser)
    api.DELETE("/user/:id", userHandler.DeleteUser)
  }
  return r
}
