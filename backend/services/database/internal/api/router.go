package api

import (
  "github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
  r := gin.Default()

  r.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "OK"})
  })

  api := r.Group("/v1")
  {
    api.POST("/users", CreateUser)
    api.GET("/users", GetAllUsers)
    api.GET("/users/:id", GetUserByID)
    api.PUT("/users/:id", UpdateUser)
    api.DELETE("/users/:id", DeleteUser)

    api.POST("/companies", CreateCompany)
    api.GET("/companies", GetAllCompanies)
    api.GET("/companies/:id", GetCompanyByID)
    api.PUT("/companies/:id", UpdateCompany)
    api.DELETE("/companies/:id", DeleteCompany)

    api.POST("/transactionrecords", CreateTransactionRecord)
    api.GET("/transactionrecords", GetAllTransactionRecords)
    api.GET("/transactionrecords/:id", GetTransactionRecordsByID)
    api.PUT("/transactionrecords/:id", UpdateTransactionRecord)
    api.DELETE("/transactionrecords/:id", DeleteTransactionRecord)

    api.POST("/warehouses", CreateWarehouse)
    api.GET("/warehouses", GetAllWarehouses)
    api.GET("/warehouses/:id", GetWarehouseByID)
    api.PUT("/warehouses/:id", UpdateWarehouse)
    api.DELETE("/warehouses/:id", DeleteWarehouse)

    api.POST("/roles", CreateRole)
    api.GET("/roles", GetAllRoles)
    api.GET("/roles/:id", GetRoleByID)
    api.PUT("/roles/:id", UpdateRole)
    api.DELETE("/roles/:id", DeleteRole)

  }
  return r
}
