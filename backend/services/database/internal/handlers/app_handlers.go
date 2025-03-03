package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	// Internal
	"github.com/yungbote/slotter/backend/services/database/internal/repos"
	"github.com/yungbote/slotter/backend/services/database/internal/services"
)

// AppHandler holds a reference to your core AppSvc.
type AppHandler struct {
	appSvc services.AppSvc
}

// NewAppHandler constructs an AppHandler with the given AppSvc.
func NewAppHandler(appSvc services.AppSvc) *AppHandler {
	return &AppHandler{appSvc: appSvc}
}

// RegisterRoutes attaches all routes to the provided router group.
func (h *AppHandler) RegisterRoutes(rg *gin.RouterGroup) {

	// AUTH
	rg.POST("/register", h.RegisterUserLocal)
	rg.POST("/login", h.LoginUserLocal)
	rg.POST("/logout", h.LogoutUser)
	rg.POST("/login/google", h.LoginWithGoogle) // Could be GET with a redirect, depending on your OAuth flow

	// COMPANY
	rg.POST("/company", h.CreateCompany)
	rg.GET("/company/:company_id", h.GetCompanyByID)
	rg.PUT("/company/:company_id/name", h.UpdateCompanyName)
	rg.PUT("/company/:company_id/avatar", h.UpdateCompanyAvatar)

	// WAREHOUSE
	rg.POST("/warehouse", h.CreateWarehouse)
	rg.GET("/warehouse/:warehouse_id", h.GetWarehouseByID)
	rg.PUT("/warehouse/:warehouse_id/name", h.UpdateWarehouseName)
	rg.DELETE("/warehouse/:warehouse_id", h.DeleteWarehouse)
	rg.GET("/warehouses", h.ListWarehouses)

	// LOCATION
	rg.POST("/warehouse/:warehouse_id/location", h.CreateLocation)
	rg.GET("/location/:location_id", h.GetLocationByID)
	rg.DELETE("/location/:location_id", h.DeleteLocation)
	rg.GET("/locations", h.ListLocations)

	// TRANSACTION FILE
	rg.POST("/warehouse/:warehouse_id/transaction-file/upload", h.UploadTransactionFile)
	rg.PUT("/transaction-file/:file_id/name", h.UpdateTransactionFileName)
	rg.DELETE("/transaction-file/:file_id", h.DeleteTransactionFile)
	rg.GET("/transaction-files", h.ListTransactionFiles)

	// TRANSACTION RECORD
	rg.POST("/warehouse/:warehouse_id/transaction-record", h.CreateTransactionRecord)
	rg.GET("/transaction-record/:record_id", h.GetTransactionRecordByID)
	rg.PUT("/transaction-record/:record_id/order-name", h.UpdateTransactionRecordOrderName)
	rg.PUT("/transaction-record/:record_id/description", h.UpdateTransactionRecordDescription)
	rg.PUT("/transaction-record/:record_id/transaction-quantity", h.UpdateTransactionRecordTransactionQuantity)
	rg.PUT("/transaction-record/:record_id/completed-quantity", h.UpdateTransactionRecordCompletedQuantity)
	rg.PUT("/transaction-record/:record_id/completed-date", h.UpdateTransactionRecordCompletedDate)
	rg.PUT("/transaction-record/:record_id/transaction-type", h.UpdateTransactionRecordTransactionType)
	rg.GET("/transaction-records", h.ListTransactionRecords)

	// USER
	rg.PUT("/user/:user_id/avatar", h.UpdateUserAvatar)
	rg.PUT("/user/:user_id/first-name", h.UpdateUserFirstName)
	rg.PUT("/user/:user_id/last-name", h.UpdateUserLastName)
	rg.PUT("/user/:user_id/email", h.UpdateUserEmail)
	rg.PUT("/user/:user_id/password", h.UpdateUserPassword)
	rg.DELETE("/user/:user_id", h.DeleteUser)
	rg.GET("/users", h.ListUsers)

	// ITEM
	rg.GET("/items", h.ListItems)
}

// ---------------------------------------------------------------------------
// AUTH Handlers
// ---------------------------------------------------------------------------

// RegisterUserLocal handles POST /register
func (h *AppHandler) RegisterUserLocal(c *gin.Context) {
	type reqBody struct {
		Email             string `json:"email"`
		Password          string `json:"password"`
		FirstName         string `json:"first_name"`
		LastName          string `json:"last_name"`
		CreateCompanyName string `json:"create_company_name"`
	}

	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, accessToken, refreshToken, err := h.appSvc.RegisterUserLocal(
		c.Request.Context(),
		body.Email,
		body.Password,
		body.FirstName,
		body.LastName,
		body.CreateCompanyName,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user, "access_token": accessToken, "refresh_token": refreshToken})
}

// LoginUserLocal handles POST /login
func (h *AppHandler) LoginUserLocal(c *gin.Context) {
	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	user, accessToken, refreshToken, err := h.appSvc.LoginUserLocal(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user, "access_token": accessToken, "refresh_token": refreshToken})
}

// LogoutUser handles POST /logout
func (h *AppHandler) LogoutUser(c *gin.Context) {
	// Typically, you'd remove an auth token from a store or add it to a blacklist.
	// If you have a "LogoutUser" method, call it. We also retrieve userID from context if needed.
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)
	err := h.appSvc.LogoutUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalSeverError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AppHandler) RefreshToken(c *gin.Context) {
	var body struct {
		UserID				string			`json:"user_id"`
		RefreshToken	string			`json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	uid, err := uuid.Parse(body.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	newAccessToken, newRefreshToken, err := h.appSvc.RefreshTokens(c.Request.Context(), uid, body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken, "refresh_token": newRefreshToken})
}

// LoginWithGoogle handles POST /login/google (or GET with code & state).
func (h *AppHandler) LoginWithGoogle(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	user, token, err := h.appSvc.LoginWithGoogle(c.Request.Context(), code, state)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

// ---------------------------------------------------------------------------
// COMPANY Handlers
// ---------------------------------------------------------------------------

// CreateCompany handles POST /company
func (h *AppHandler) CreateCompany(c *gin.Context) {
	// Some apps let only logged-in users create a company.
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	_ = userIDVal.(uuid.UUID) // just for reference; your AppSvc might or might not require it

	type reqBody struct {
		Name           string `json:"name"`
		GenerateAvatar bool   `json:"generate_avatar"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	company, err := h.appSvc.CreateCompany(c.Request.Context(), body.Name, body.GenerateAvatar)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, company)
}

// GetCompanyByID handles GET /company/:company_id
func (h *AppHandler) GetCompanyByID(c *gin.Context) {
	companyIDStr := c.Param("company_id")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company_id"})
		return
	}
	company, err := h.appSvc.GetCompanyByID(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, company)
}

// UpdateCompanyName handles PUT /company/:company_id/name
func (h *AppHandler) UpdateCompanyName(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	companyIDStr := c.Param("company_id")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company_id"})
		return
	}

	type reqBody struct {
		NewName string `json:"new_name"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.appSvc.UpdateCompanyName(c.Request.Context(), userID, companyID, body.NewName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "company name updated"})
}

// UpdateCompanyAvatar handles PUT /company/:company_id/avatar
func (h *AppHandler) UpdateCompanyAvatar(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	companyIDStr := c.Param("company_id")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company_id"})
		return
	}

	newAvatar, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read avatar data"})
		return
	}

	url, err := h.appSvc.UpdateCompanyAvatar(c.Request.Context(), userID, companyID, newAvatar)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatar_url": url})
}

// ---------------------------------------------------------------------------
// WAREHOUSE Handlers
// ---------------------------------------------------------------------------

// CreateWarehouse handles POST /warehouse
func (h *AppHandler) CreateWarehouse(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	type reqBody struct {
		Name string `json:"name"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	warehouse, err := h.appSvc.CreateWarehouse(c.Request.Context(), userID, body.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, warehouse)
}

// GetWarehouseByID handles GET /warehouse/:warehouse_id
func (h *AppHandler) GetWarehouseByID(c *gin.Context) {
	// user ID might or might not be needed
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	warehouseIDStr := c.Param("warehouse_id")
	warehouseID, err := uuid.Parse(warehouseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse_id"})
		return
	}

	warehouse, err := h.appSvc.GetWarehouseByID(c.Request.Context(), warehouseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	// Could check ownership if needed. 
	c.JSON(http.StatusOK, warehouse)
}

// UpdateWarehouseName handles PUT /warehouse/:warehouse_id/name
func (h *AppHandler) UpdateWarehouseName(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	warehouseIDStr := c.Param("warehouse_id")
	warehouseID, err := uuid.Parse(warehouseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse_id"})
		return
	}

	type reqBody struct {
		NewName string `json:"new_name"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.appSvc.UpdateWarehouseName(c.Request.Context(), userID, warehouseID, body.NewName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "warehouse name updated"})
}

// DeleteWarehouse handles DELETE /warehouse/:warehouse_id
func (h *AppHandler) DeleteWarehouse(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	warehouseIDStr := c.Param("warehouse_id")
	warehouseID, err := uuid.Parse(warehouseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse_id"})
		return
	}

	if err := h.appSvc.DeleteWarehouse(c.Request.Context(), userID, warehouseID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "warehouse deleted"})
}

// ListWarehouses handles GET /warehouses
func (h *AppHandler) ListWarehouses(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var f repos.WarehouseFilter
	// Optionally parse query params for filtering, e.g. start_date, end_date
	// For brevity, skipping. You can do something like:
	// if s := c.Query("start_date"); s != "" { parse time, set f.StartDate }
	warehouses, err := h.appSvc.ListWarehouses(c.Request.Context(), userID, f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, warehouses)
}

// ---------------------------------------------------------------------------
// LOCATION Handlers
// ---------------------------------------------------------------------------

// CreateLocation handles POST /warehouse/:warehouse_id/location
func (h *AppHandler) CreateLocation(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	warehouseIDStr := c.Param("warehouse_id")
	warehouseID, err := uuid.Parse(warehouseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse_id"})
		return
	}

	type reqBody struct {
		LocationPath     string `json:"location_path"`
		LocationNamePath string `json:"location_name_path"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	location, err := h.appSvc.CreateLocation(c.Request.Context(), userID, warehouseID, body.LocationPath, body.LocationNamePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, location)
}

// GetLocationByID handles GET /location/:location_id
func (h *AppHandler) GetLocationByID(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	locationIDStr := c.Param("location_id")
	locationID, err := uuid.Parse(locationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location_id"})
		return
	}

	loc, err := h.appSvc.GetLocationByID(c.Request.Context(), userID, locationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loc)
}

// DeleteLocation handles DELETE /location/:location_id
func (h *AppHandler) DeleteLocation(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	locationIDStr := c.Param("location_id")
	locationID, err := uuid.Parse(locationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location_id"})
		return
	}

	if err := h.appSvc.DeleteLocation(c.Request.Context(), userID, locationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "location deleted"})
}

// ListLocations handles GET /locations
func (h *AppHandler) ListLocations(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var f repos.LocationFilter
	// Optionally parse query params for filtering
	locations, err := h.appSvc.ListLocations(c.Request.Context(), userID, f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, locations)
}

// ---------------------------------------------------------------------------
// TRANSACTION FILE Handlers
// ---------------------------------------------------------------------------

// UploadTransactionFile handles POST /warehouse/:warehouse_id/transaction-file/upload
func (h *AppHandler) UploadTransactionFile(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	warehouseIDStr := c.Param("warehouse_id")
	warehouseID, err := uuid.Parse(warehouseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse_id"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
		return
	}

	fileName := fileHeader.Filename
	fileData, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer fileData.Close()

	buf := make([]byte, fileHeader.Size)
	_, err = fileData.Read(buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file data"})
		return
	}

	tf, err := h.appSvc.UploadTransactionFile(c.Request.Context(), userID, warehouseID, fileName, buf)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tf)
}

// UpdateTransactionFileName handles PUT /transaction-file/:file_id/name
func (h *AppHandler) UpdateTransactionFileName(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	fileIDStr := c.Param("file_id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file_id"})
		return
	}

	type reqBody struct {
		NewName string `json:"new_name"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.appSvc.UpdateTransactionFileName(c.Request.Context(), userID, fileID, body.NewName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transaction file name updated"})
}

// DeleteTransactionFile handles DELETE /transaction-file/:file_id
func (h *AppHandler) DeleteTransactionFile(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	fileIDStr := c.Param("file_id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file_id"})
		return
	}

	err = h.appSvc.DeleteTransactionFile(c.Request.Context(), userID, fileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transaction file deleted"})
}

// ListTransactionFiles handles GET /transaction-files
func (h *AppHandler) ListTransactionFiles(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var f repos.TransactionFileFilter
	// parse any query params if needed
	tfiles, err := h.appSvc.ListTransactionFiles(c.Request.Context(), userID, f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tfiles)
}

// ---------------------------------------------------------------------------
// TRANSACTION RECORD Handlers
// ---------------------------------------------------------------------------

// CreateTransactionRecord handles POST /warehouse/:warehouse_id/transaction-record
func (h *AppHandler) CreateTransactionRecord(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	warehouseIDStr := c.Param("warehouse_id")
	warehouseID, err := uuid.Parse(warehouseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse_id"})
		return
	}

	type reqBody struct {
		TransactionType  string `json:"transaction_type"`
		OrderName        string `json:"order_name"`
		Description      string `json:"description"`
		TransactionQ     int64  `json:"transaction_quantity"`
		CompletedQ       int64  `json:"completed_quantity"`
		CompletedDateStr string `json:"completed_date"` // parse to time
		LocationPath     string `json:"location_path"`
		LocationNamePath string `json:"location_name_path"`
		ItemName         string `json:"item_name"`
	}

	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var completedDate time.Time
	if body.CompletedDateStr != "" {
		cd, parseErr := time.Parse("2006-01-02", body.CompletedDateStr)
		if parseErr == nil {
			completedDate = cd
		}
	}

	err = h.appSvc.CreateTransactionRecord(
		c.Request.Context(),
		userID,
		warehouseID,
		body.TransactionType,
		body.OrderName,
		body.Description,
		body.TransactionQ,
		body.CompletedQ,
		completedDate,
		body.LocationPath,
		body.LocationNamePath,
		body.ItemName,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transaction record created"})
}

// GetTransactionRecordByID handles GET /transaction-record/:record_id
func (h *AppHandler) GetTransactionRecordByID(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	recordIDStr := c.Param("record_id")
	recordID, err := uuid.Parse(recordIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record_id"})
		return
	}

	rec, err := h.appSvc.GetTransactionRecordByID(c.Request.Context(), userID, recordID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rec)
}

// UpdateTransactionRecordOrderName handles PUT /transaction-record/:record_id/order-name
func (h *AppHandler) UpdateTransactionRecordOrderName(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	recordIDStr := c.Param("record_id")
	recordID, err := uuid.Parse(recordIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record_id"})
		return
	}

	type reqBody struct {
		NewOrderName string `json:"new_order_name"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.appSvc.UpdateTransactionRecordOrderName(c.Request.Context(), userID, recordID, body.NewOrderName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transaction record order_name updated"})
}

// UpdateTransactionRecordDescription handles PUT /transaction-record/:record_id/description
func (h *AppHandler) UpdateTransactionRecordDescription(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	recordIDStr := c.Param("record_id")
	recordID, err := uuid.Parse(recordIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record_id"})
		return
	}

	type reqBody struct {
		NewDescription string `json:"new_description"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.appSvc.UpdateTransactionRecordDescription(c.Request.Context(), userID, recordID, body.NewDescription)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transaction record description updated"})
}

// UpdateTransactionRecordTransactionQuantity handles PUT /transaction-record/:record_id/transaction-quantity
func (h *AppHandler) UpdateTransactionRecordTransactionQuantity(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	recordIDStr := c.Param("record_id")
	recordID, err := uuid.Parse(recordIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record_id"})
		return
	}

	type reqBody struct {
		NewTQuantity int64 `json:"new_transaction_quantity"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.appSvc.UpdateTransactionRecordTransactionQuantity(c.Request.Context(), userID, recordID, body.NewTQuantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transaction record transaction_quantity updated"})
}

// UpdateTransactionRecordCompletedQuantity handles PUT /transaction-record/:record_id/completed-quantity
func (h *AppHandler) UpdateTransactionRecordCompletedQuantity(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	recordIDStr := c.Param("record_id")
	recordID, err := uuid.Parse(recordIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record_id"})
		return
	}

	type reqBody struct {
		NewCQuantity int64 `json:"new_completed_quantity"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.appSvc.UpdateTransactionRecordCompletedQuantity(c.Request.Context(), userID, recordID, body.NewCQuantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transaction record completed_quantity updated"})
}

// UpdateTransactionRecordCompletedDate handles PUT /transaction-record/:record_id/completed-date
func (h *AppHandler) UpdateTransactionRecordCompletedDate(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	recordIDStr := c.Param("record_id")
	recordID, err := uuid.Parse(recordIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record_id"})
		return
	}

	type reqBody struct {
		NewDateStr string `json:"new_completed_date"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var newDate time.Time
	if body.NewDateStr != "" {
		parsed, parseErr := time.Parse("2006-01-02", body.NewDateStr)
		if parseErr == nil {
			newDate = parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (YYYY-MM-DD)"})
			return
		}
	}

	err = h.appSvc.UpdateTransactionRecordCompletedDate(c.Request.Context(), userID, recordID, newDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transaction record completed_date updated"})
}

// UpdateTransactionRecordTransactionType handles PUT /transaction-record/:record_id/transaction-type
func (h *AppHandler) UpdateTransactionRecordTransactionType(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	recordIDStr := c.Param("record_id")
	recordID, err := uuid.Parse(recordIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record_id"})
		return
	}

	type reqBody struct {
		NewType string `json:"new_transaction_type"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.appSvc.UpdateTransactionRecordTransactionType(c.Request.Context(), userID, recordID, body.NewType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transaction record transaction_type updated"})
}

// ListTransactionRecords handles GET /transaction-records
func (h *AppHandler) ListTransactionRecords(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var f repos.TransactionRecordFilter
	// parse query params if needed (e.g. f.StartDate, f.EndDate)
	recs, err := h.appSvc.ListTransactionRecords(c.Request.Context(), userID, f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recs)
}

// ---------------------------------------------------------------------------
// USER Handlers
// ---------------------------------------------------------------------------

// UpdateUserAvatar handles PUT /user/:user_id/avatar
func (h *AppHandler) UpdateUserAvatar(c *gin.Context) {
	// "user_id" in the route is the user being updated, but the context user might be the same or an admin
	requestingUserIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	requestingUserID := requestingUserIDVal.(uuid.UUID)

	paramUserIDStr := c.Param("user_id")
	targetUserID, err := uuid.Parse(paramUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	newAvatar, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read avatar"})
		return
	}

	url, err := h.appSvc.UpdateUserAvatar(c.Request.Context(), targetUserID, newAvatar)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"avatar_url": url, "updated_for_user_id": targetUserID, "requested_by_user_id": requestingUserID})
}

// UpdateUserFirstName handles PUT /user/:user_id/first-name
func (h *AppHandler) UpdateUserFirstName(c *gin.Context) {
	requestingUserIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	_ = requestingUserIDVal.(uuid.UUID)

	paramUserIDStr := c.Param("user_id")
	targetUserID, err := uuid.Parse(paramUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	type reqBody struct {
		NewFirst string `json:"new_first_name"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.appSvc.UpdateUserFirstName(c.Request.Context(), targetUserID, body.NewFirst)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user first_name updated"})
}

// UpdateUserLastName handles PUT /user/:user_id/last-name
func (h *AppHandler) UpdateUserLastName(c *gin.Context) {
	requestingUserIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	_ = requestingUserIDVal.(uuid.UUID)

	paramUserIDStr := c.Param("user_id")
	targetUserID, err := uuid.Parse(paramUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	type reqBody struct {
		NewLast string `json:"new_last_name"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.appSvc.UpdateUserLastName(c.Request.Context(), targetUserID, body.NewLast)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user last_name updated"})
}

// UpdateUserEmail handles PUT /user/:user_id/email
func (h *AppHandler) UpdateUserEmail(c *gin.Context) {
	requestingUserIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	_ = requestingUserIDVal.(uuid.UUID)

	paramUserIDStr := c.Param("user_id")
	targetUserID, err := uuid.Parse(paramUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	type reqBody struct {
		NewEmail string `json:"new_email"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.appSvc.UpdateUserEmail(c.Request.Context(), targetUserID, body.NewEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user email updated"})
}

// UpdateUserPassword handles PUT /user/:user_id/password
func (h *AppHandler) UpdateUserPassword(c *gin.Context) {
	requestingUserIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	requestingUserID := requestingUserIDVal.(uuid.UUID)

	paramUserIDStr := c.Param("user_id")
	targetUserID, err := uuid.Parse(paramUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	type reqBody struct {
		NewPassword string `json:"new_password"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Usually you'd check requestingUserID == targetUserID or if an admin.
	err = h.appSvc.UpdateUserPassword(c.Request.Context(), requestingUserID, targetUserID, body.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user password updated"})
}

// DeleteUser handles DELETE /user/:user_id
func (h *AppHandler) DeleteUser(c *gin.Context) {
	requestingUserIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	requestingUserID := requestingUserIDVal.(uuid.UUID)

	paramUserIDStr := c.Param("user_id")
	targetUserID, err := uuid.Parse(paramUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	err = h.appSvc.DeleteUser(c.Request.Context(), requestingUserID, targetUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// ListUsers handles GET /users
func (h *AppHandler) ListUsers(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var f repos.UserFilter
	// parse optional query string e.g. f.EmailLike, f.Role, etc.
	users, err := h.appSvc.ListUsers(c.Request.Context(), userID, f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// ---------------------------------------------------------------------------
// ITEM Handler
// ---------------------------------------------------------------------------

// ListItems handles GET /items
func (h *AppHandler) ListItems(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var f repos.ItemFilter
	// parse query params if needed
	items, err := h.appSvc.ListItems(c.Request.Context(), userID, f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

