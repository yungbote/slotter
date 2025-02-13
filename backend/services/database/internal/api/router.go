package api

import (
  "github.com/gin-gonic/gin"
  "github.com/yungbote/slotter/backend/services/database/internal/handlers"
)

func NewRouter() *gin.Engine {
  r := gin.Default()

  r.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "OK"})
  })

  api := r.Group("/v1")
  {
    api.POST("/users", handlers.CreateUser)
    api.GET("/users", handlers.GetAllUsers)
    api.GET("/users/:id", handlers.GetUserByID)
    api.PUT("/users/:id", handlers.UpdateUser)
    api.DELETE("/users/:id", handlers.DeleteUser)

    api.POST("/companies", handlers.CreateCompany)
    api.GET("/companies", handlers.GetAllCompanies)
    api.GET("/companies/:id", handlers.GetCompanyByID)
    api.PUT("/companies/:id", handlers.UpdateCompany)
    api.DELETE("/companies/:id", handlers.DeleteCompany)

    api.POST("/transactionrecords", handlers.CreateTransactionRecord)
    api.GET("/transactionrecords", handlers.GetAllTransactionRecords)
    api.GET("/transactionrecords/:id", handlers.GetTransactionRecordsByID)
    api.PUT("/transactionrecords/:id", handlers.UpdateTransactionRecord)
    api.DELETE("/transactionrecords/:id", handlers.DeleteTransactionRecord)

    api.POST("/warehouses", handlers.CreateWarehouse)
    api.GET("/warehouses", handlers.GetAllWarehouses)
    api.GET("/warehouses/:id", handlers.GetWarehouseByID)
    api.PUT("/warehouses/:id", handlers.UpdateWarehouse)
    api.DELETE("/warehouses/:id", handlers.DeleteWarehouse)

    api.POST("/roles", handlers.CreateRole)
    api.GET("/roles", handlers.GetAllRoles)
    api.GET("/roles/:id", handlers.GetRoleByID)
    api.PUT("/roles/:id", handlers.UpdateRole)
    api.DELETE("/roles/:id", handlers.DeleteRole)

  }
  return r
}
