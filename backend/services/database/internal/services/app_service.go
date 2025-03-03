package services

import (
  "context"
  "errors"
  "fmt"
  "strings"

  "github.com/google/uuid"
  "golang.org/x/crypto/bcrypt"

  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/repos"
  "github.com/yungbote/slotter/backend/services/database/internal/events"
  "github.com/yungbote/slotter/backend/services/database/internal/auth"
  "github.com/yungbote/slotter/backend/services/database/internal/services/avatar"
  "github.com/yungbote/slotter/backend/services/database/internal/services/s3"
)

type ParserService interface {
  ParseFile(ctx context.Context, fileName string, fileData []byte, transactionFileID, companyID, warehouseID uuid.UUID) (int, error)
}

type AppSvc interface {
  //Auth
  RegisterUserLocal(ctx context.Context, email, password, firstName, lastName string, createCompanyName string) (*models.User, string, error)
  LoginUserLocal(ctx context.Context, email, password string) (*models.User, string, error)
  LoginWithGoogle(ctx context.Context, code, state string) (*models.User, string, error)
  LogoutUser()

  //Company
  CreateCompany(ctx context.Context, name string, generateAvatar bool) (*models.Company, error)
  GetCompanyByID(ctx context.Context, companyID uuid.UUID) (*models.Company, error)
  UpdateCompanyAvatar(ctx context.Context, userID uuid.UUID, newAvatar ) error
  UpdateCompanyName(ctx context.Context, userID uuid.UUID, newName string)

  //Warehouse
  CreateWarehouse(ctx context.Context, userID uuid.UUID, createWarehouseName string) error
  GetWarehouseByID(ctx context.Context, warehouseID uuid.UUID) (*models.Warehouse, error)
  UpdateWarehouseName(ctx context.Context, userID uuid.UUID, newWarehouseName string) error
  DeleteWarehouse(ctx context.Context, userID, warehouseID uuid.UUID) error
  ListWarehouses(ctx context.Context, userID uuid.UUID, f repos.WarehouseFilter) ([]*models.Warehouse, error)

  //Location
  CreateLocation(ctx context.Context, userID, warehouseID uuid.UUID, locationPath, locationNamePath string) error
  GetLocationByID(ctx context.Context, locationID uuid.UUID) (*models.Location, error)
  GetLocationByPath(ctx context.Context, userID, warehouseID uuid.UUID, locationPath string) (*models.Location, error)
  DeleteLocation(ctx context.Context, userID, locationID uuid.UUID) error
  ListLocations(ctx context.Context, userID uuid.UUID, f repos.LocationFilter) ([]*models.Location, error)

  //TransactionFile
  UploadTransactionFile()
  UpdateTransactionFileName(ctx context.Context, userID, fileID uuid.UUID) error
  DeleteTransactionFile(ctx context.Context, userID, fileID uuid.UUID) error
  ListTransactionFiles(ctx context.Context, userID uuid.UUID, f repos.TransactionFileFilter) ([]*models.TransactionFile, error)

  //TransactionRecord
  CreateTransactionRecord(ctx context.Context, userID, warehouseID uuid.UUID, transactionType, orderName, description string, transactionQ, completedQ int64, completedDate time.Time, locationPath, locationNamePath, itemName string) error
  GetTransactionRecordByID(ctx context.Context, recordID uuid.UUID) (*models.TransactionRecord, error)
  UpdateTransactionRecordOrderName(ctx context.Context, recordID uuid.UUID, newOrderName string) error
  UpdateTransactionRecordDescription(ctx context.Context, recordID uuid.UUID, newDescription string) error
  UpdateTransactionRecordTransactionQuantity(ctx context.Context, recordID uuid.UUID, newTQuantity int64) error
  UpdateTransactionRecordCompletedQuantity(ctx context.Context, recordID uuid.UUID, newCQuantity int64) error
  UpdateTransactionRecordCompletedDate(ctx context.Context, recordID uuid.UUID, newDate time.Time) error
  ListTransactionRecords(ctx context.Context, userID uuid.UUID, f repos.TransactionRecordFilter) ([]*models.TransactionRecord, error)

  //User
  UpdateUserAvatar(ctx context.Context, userID uuid.UUID, newAvatar )
  UpdateUserFirstName(ctx context.Context, userID uuid.UUID, newFirst string) error
  UpdateUserLastName(ctx context.Context, userID uuid.UUID, newLast string) error
  UpdateUserEmail(ctx context.Context, userID uuid.UUID, newEmail string) error
  UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPass string) error
  DeleteUser(ctx context.Context, userID uuid.UUID) error
  ListUsers(ctx context.Context, userID uuid.UUID, f repos.UserFilter) ([]*models.User, error)

  //Item
  ListItems(ctx context.Context, userID uuid.UUID, f repos.ItemFilter) ([]*models.Item, error)

  //Utility
  generateUserAvatar(ctx context.Context, firstName, lastName string) (string, error)
  generateCompanyAvatar(ctx contex.Context, companyName string) (string, error)
}

type appSvc struct {
  csvc      CSvc
  usvc      USvc
  wsvc      WSvc
  lsvc      LSvc
  tfsvc     TFSvc
  trsvc     TRSvc
  isvc      ISvc
  
  avatarsvc avatar.AvatarService
  s3svc     s3.S3Service

  tokensvc  auth.TokenService
  oauthsvc  auth.OAuthService

  pub       events.PubSubPublisher
  uact      repos.UserActionRepo

  parsersvc ParserService
}

func NewAppSvc(csvc: CSvc, usvc: USvc, wsvc: WSvc, lsvc: LSvc, tfsvc: TFSvc, trsvc: TRSvc, isvc: ISvc, avatarsvc: avatar.AvatarService, s3svc: s3.S3Service, tokensvc: auth.TokenService, oauthsvc: auth.OAuthService, pub: events.PubSubPublisher, uact repos.UserActionRepo, parsersvc ParserService) AppSvc {
  return &appSvc{csvc: csvc, usvc: usvc, wsvc: wsvc, lsvc: lsvc, tfsvc: tfsvc, trsvc: trsvc, isvc: isvc, avatarsvc: avatarsvc, s3svc: s3svc, tokensvc: tokensvc, oauthsvc: oauthsvc, pub: pub, uact: uact, parsersvc: parsersvc}
}

func (s *appSvc) RegisterUserLocal(ctx context.Context, email, password, firstName, lastName string, createCompanyName string, companyID uuid.UUID) (*models.User, string, error) {
  email = strings.TrimSpace(email)
  if email == "" || password == "" || firstName == "" || lastName == "" {
    return nil, "", fmt.Errorf("missing required fields for registration")
  }
  existing, err := s.usvc.GetUserByEmail(email)
  if err == nil && existing != nil {
    return nil, "", fmt.Errorf("email already in use")
  }
  var ncID *uuid.UUID
  if createCompanyName != "" {
    avatarURL, err := s.generateCompanyAvatar(ctx, createCompanyName)
    if err != nil {
      return nil, "", fmt.Errorf("Failed to generate company avatar: %w", err)
    }
    newCo := models.Company{
      Name:       createCompanyName,
      AvatarURL:  avatarURL,
    }
    createdCo, err := s.csvc.CreateCompany(newCo)
    if err != nil {
      return nil, "", fmt.Errorf("failed to create new company: %w", err)
    }
    ncID = &createdCo.IDkcreatedCo.ID
    _ = s.pub.PublishCompanyEvent(createdCo.ID, "COMPANY_CREATED", map[string]interface{}{"company_name": createCompanyName})
  }
  hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    return nil, "", fmt.Errorf("failed to hash password: %w", err)
  }
  userAvatar, err := s.generateUserAvatar(ctx, firstName, lastName)
  if err != nil {
    userAvatar = ""
  }
  u := models.User{Email: email, Password: string(hashed), FirstName: firstName, LastName: lastName, AvatarURL: userAvatar}
  if ncID != nil {
    u.CompanyID = ncID
  }
  createdUser, err := s.usvc.CreateUser(u)
  if err != nil {
    return nil, "", fmt.Errorf("failed to create user: %w", err)
  }
  token, err := s.tokensvc.GenerateToken(createdUser.ID, createdUser.Email)
  if err != nil {
    return nil, "", fmt.Errorf("failed to generate JWT: %w", err)
  }
  _ = s.pub.PublishCompanyEvent(createdUser.CompanyID, "USER_REGISTERED", map[string]interface{}{"user_id": createdUser.ID, "email": createdUser.Email, "company_id": createdUser.CompanyID})
  return createdUser, token, nil
}

func (s *appSvc) LoginUserLocal(ctx context.Context, email, password string) (*models.User, string, error) {
  if email == "" || password == "" {
    return nil, "", fmt.Errorf("missing email or password")
  }
  user, err := s.usvc.GetUserByEmail(email)
  if err != nil || user == nil {
    return nil, "", errors.New("Invalid credentials")
  }
  if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
    return nil, "", errors.New("Invalid credentials")
  }
  token, err := s.tokenService.GenerateToken(user.ID, user.Email)
  if err != nil {
    return nil, "", err
  }
  _ = s.pub.PublishCompanyEvent(user.CompanyID, "USER_LOGIN", map[string]interface{}{"user_id": user.ID, "email": user.Email})
  _ = s.pub.PublishUserEvent(user.ID, "USER_LOGIN", map[string]interface{}{"user_id": user.ID, "email": user.Email})
  
  return user, token, nil
}

func (s *appSvc) LoginWithGoogle(ctx context.Context, code, state string) (*models.User, string, error) {
  user, err := s.oauthsvc.HandleGoogleCallback(code, state)
  if err != nil {
    return nil, "", err
  }
  token, err := s.tokensvc.GenerateToken(user.ID, user.Email)
  if err != nil {
    return nil, "", err
  }
  _ = s.pub.PublishCompanyEvent(user.CompanyID, "USER_LOGIN", map[string]interface{}{"user_id": user.ID, "email": user.Email})
  _ = s.pub.PublishUserEvent(user.ID, "USER_LOGIN", map[string]interface{}{"user_id": user.ID, "email": user.Email})
  
  return user, token, nil
}

func (s *appSvc) CreateCompany(ctx context.Context, name string, generateAvatar bool) (*models.Company, error) {
  if name == "" {
    return nil, fmt.Errorf("company name is required")
  }
  var avatarURL string
  if generateAvatar {
    var err error
    avatarURL, err = s.generateCompanyAvatar(ctx, name)
    if err != nil {
      avatarURL = ""
    }
  }
  newCo := models.Company{Name: name, AvatarURL: avatarURL}
  created, err := s.csvc.CreateCompany(newCo)
  if err != nil {
    return nil, err
  }
  _ = s.pub.PublishCompanyEvent(created.ID, "COMPANY_CREATED", map[string]interface{}{"company_name": created.Name, "company_id": created.ID})
  
  return created, nil
}

func (s *appSvc) GetCompanyByID(ctx context.Context, companyID uuid.UUID) (*models.Company, error) {
  return s.csvc.GetCompanyByID(companyID)
}

func (s *appSvc) UpdateCompanyAvatar(ctx context.Context, userID uuid.UUID, newAvatar []byte) (string, error) {
  user, err := s.usvc.GetUserByID(userID)
  if err != nil {
    return "", fmt.Errorf("user not found: %w", err)
  }
  co, err := s.csvc.GetCompanyByID(user.CompanyID)
  if err != nil {
    return "", fmt.Errorf("company not found: %w", err)
  }
  if len(newAvatar) == 0 {
    return "", fmt.Errorf("avatar data is empty")
  }
  //Upload new avatar to s3
  fileName := fmt.Sprintf("company-avatar-%s.png", co.ID)
  contentType := "image/png"
  url, err := s.s3svc.UploadImage(ctx, fileName, newAvatar, contentType)
  if err != nil {
    return "", fmt.Errorf("failed to upload new company avatar: %w", err)
  }
  if err := s.csvc.UpdateCompanyAvatarURL(co.ID, url); err != nil {
    return "", fmt.Errorf("failed to update company avatar URL: %w", err)
  }
  _ = s.pub.PublishCompanyEvent(co.ID, "COMPANY_AVATAR_UPDATED", map[string]interface{}{"company_id": co.ID, "updated_by": userID, "avatar_url": url})
  return url, nil
}

func (s *appSvc) UpdateCompanyName(ctx context.Context, userID uuid.UUID, newName string) error {
  if newName == "" {
    return fmt.Errorf("new company name is required")
  }
  user, err := s.usvc.GetUserByID(userID)
  if err != nil {
    return fmt.Errorf("user not found: %w", err)
  }
  co, err := s.csvc.GetCompanyByID(user.CompanyID)
  if err != nil || co == nil {
    return fmt.Errorf("company not found: %w", err)
  }
  if err := s.csvc.UpdateCompanyName(co.ID, newName); err != nil {
    return fmt.Errorf("failed to update company name: %w", err)
  }
  _ = s.pub.PublishCompanyEvent(co.ID, "COMPANY_NAME_UPDATED", map[string]interface{}{"company_id": co.ID, "updated_by": userID, "new_name": newName})
  return nil
}

func (s *appSvc) CreateWarehouse(ctx context.Context, userID uuid.UUID, warehouseName string) (*models.Warehouse, error) {
  if warehouseName == "" {
    return nil, fmt.Errorf("warehouse name is required")
  }
  user, err := s.usvc.GetUserByID(userID)
  if err != nil || user == nil {
    return nil, fmt.Errorf("user not found or invalid")
  }
  if user.CompanyID == nil {
    return nil, fmt.Errorf("cannot create warehouse: user has no associated company id")
  }
  warehouse := models.Warehouse{
    Name:             warehouseName,
    CompanyID:        user.CompanyID,
  }
  created, err := s.wsvc.CreateWarehouse(warehouse)
  if err != nil {
    return nil, fmt.Errorf("failed to create warehouse: %w", err)
  }
  if created == nil {
    return nil, fmt.Errorf("create warehouse returned nil")
  }
  _ = s.pub.PublishCompanyEvent(*user.CompanyID, "WAREHOUSE_CREATED", map[string]interface{}{"warehouse_id": created.ID, "warehouse_name": created.Name, "created_by": userID})
  return created, nil
}

func (s *appSvc) UpdateWarehouseName(ctx context.Context, userID, warehouseID uuid.UUID) error {
  if newName == "" {
    return fmt.Errorf("new warehouse name is required")
  }
  wh, err := s.GetWarehouseByID(ctx, userID, warehouseID)
  if err != nil {
    return err
  }
  err = s.wsvc.UpdateWarehouseName(wh.ID, newName)
  if err != nil {
    return fmt.Errorf("failed to update warehouse name: %w", err)
  }
  _ = s.pub.PublishCompanyEvent(*wh.CompanyID, "WAREHOUSE_NAME_UPDATED", map[string]interface{}{"warehouse_id": wh.ID, "new_name": newName, "updated_by": userID})
  return nil
}

func (s *appSvc) DeleteWarehouse(ctx context.Context, userID, warehouseID uuid.UUID) error {
  wh, err := s.GetWarehouseByID(ctx, userID, warehouseID)
  if err != nil {
    return err
  }
  err = s.wsvc.DeleteWarehouse(wh.ID)
  if err != nil {
    return fmt.Errorf("failed to delete warehouse: %w", err)
  }
  _ = s.pub.PublishCompanyEvent(*wh.CompanyID, "WAREHOUSE_DELETED", map[string]interface{}{"warehouse_id": wh.ID, "deleted_by": userID})
  return nil
}

func (s *appSvc) ListWarehouses(ctx context.Context, userID uuid.UUID, f repos.WarehouseFilter) ([]*models.Warehouse, error) {
  user, err := s.usvc.GetUserByID(userID)
  if err != nil || user == nil {
    return nil, fmt.Errorf("unauthorized user")
  }
  if user.CompanyID == nil {
    return nil, fmt.Errorf("user has no company")
  }
  f.CompanyID = *user.CompanyID
  return s.wsvc.ListWarehouses(f)
}

func (s *appSvc) CreateLocation(ctx context.Context, userID, warehouseID uuid.UUID, locationPath, locationNamePath string) (*models.Location, error) {
  wh, err := s.GetWarehouseByID(ctx, userID, warehouseID)
  if err != nil {
    return nil, err
  }
  if locationPath == "" {
    return nil, fmt.Errorf("location path is required")
  }
  loc := models.Location{
    WarehouseID:      &wh.ID,
    LocationPath:     locationPath,
    LocationNamePath: locationNamePath,
  }
  created, err := s.lsvc.CreateLocation(loc)
  if err != nil {
    return nil, fmt.Errorf("failed to create location: %w", err)
  }
  _ = s.pub.PublishCompanyEvent(*wh.CompanyID, "LOCATION_CREATED", map[string]interface{}{"warehouse_id": wh.ID, "location_id": created.ID, "created_by": userID})
  return created, nil
}

func (s *appSvc) GetLocationByID(ctx context.Context, userID, locationID uuid.UUID) (*models.Location, error) {
  if locationID == uuid.Nil {
    return nil, fmt.Errorf("invalid locationID")
  }
  user, err := s.usvc.GetUserByID(userID)
  if err != nil || user == nil {
    return nil, fmt.Errorf("unauthorized user")
  }
  loc, err := s.lsvc.GetLocationByID(locationID)
  if err != nil {
    return nil, err
  }
  if loc.WarehouseID == nil {
    return nil, fmt.Errorf("location missing warehouse reference")
  }
  wh, err := s.wsvc.GetWarehouseByID(*loc.WarehouseID)
  if err != nil {
    return nil, fmt.Errorf("failed to get warehouse for location: %w", err)
  }
  if wh.CompanyID == nil || *wh.CompanyID != *user.CompanyID {
    return nil, fmt.Errorf("location is not in user's company warehouse")
  }
  return loc, nil
}

func (s *appSvc) GetLocationByPath(ctx context.Context, userID, warehouseID uuid.UUID, locationPath string) (*models.Location, error) {
  if warehouseID == uuid.Nil || locationPath == "" {
    return nil, fmt.Errorf("invalid warehouseID or locationPath")
  }
  _, err := s.GetWarehouseByID(ctx, userID, warehouseID)
  if err != nil {
    return nil, err
  }
  loc, err := s.lsvc.GetLocationByPath(*user.CompanyID, *(&warehouseID), locationPath)
  if err != nil {
    return nil, err
  }
  return loc, nil
}

func (s *appSvc) DeleteLocation(ctx context.Context, userID, locationID uuid.UUID) error {
  loc, err := s.GetLocationByID(ctx, userID, locationID)
  if err != nil {
    return err
  }
  wh, err := s.wsvc.GetWarehouseByID(*loc.WarehouseID)
  if err != nil {
    return fmt.Errorf("failed to get warehouse: %w", err)
  }
  err = s.lsvc.DeleteLocation(loc.ID)
  if err != nil {
    return fmt.Errorf("failed to delete location: %w", err)
  }
  _ = s.pub.PublishCompanyEvent(*wh.CompanyID, "LOCATION_DELETED", map[string]interface{}{"location_id": loc.ID, "deleted_by": userID, "warehouse_id": wh.ID})
  return nil
}

func (s *appSvc) ListLocations(ctx context.Context, userID uuid.UUID, f repos.LocationFilter) ([]*models.Location, error) {
  user, err := s.usvc.GetUserByID(userID)
  if err != nil || user == nil {
    return nil, fmt.Errorf("unauthorized user")
  }
  if user.CompanyID == nil {
    return nil, fmt.Errorf("user has no company")
  }
  f.CompanyID = *user.CompanyID
  return s.lsvc.ListLocations(f)
}

func (s *appSvc) UploadTransactionFile(ctx context.Context, userID, warehouseID uuid.UUID, fileName string, data []byte) (*models.TransactionFile, error) {
  user, err := s.usvc.GetUserByID(userID)
  if err != nil || user == nil {
    return nil, fmt.Errorf("user not found: %w", err)
  }
  if user.CompanyID == nil {
    return nil, fmt.Errorf("user has no associated company")
  }
  if warehouseID == uuid.Nil {
    return nil, fmt.Errorf("invalid warehouseID")
  }
  url, err := s.s3svc.UploadFile(ctx, fileName, data)
  if err != nil {
    return nil, fmt.Errorf("failed to upload file to s3: %w", err)
  }
  fileExt := s.extractExt(fileName)
  tf := models.TransactionFile{
    FileName:         fileName,
    FileExtension:    fileExt,
    FilePathURL:      url,
    WarehouseID:      &warehouseID,
    CompanyID:        user.CompanyID,
  }
  createdFile, err := s.tfsvc.CreateTransactionFile(tf)
  if err != nil {
    return nil, fmt.Errorf("failed to create transaction file record: %w", err)
  }
  recordsCreated, err := s.parserSvc.ParseFile(ctx, fileName, data, createdFile.ID, *user.CompanyID, warehouseID)
  if err != nil {
    return nil, fmt.Errorf("failed to parse transaction file: %w", err)
  }
  _ = s.pub.PublishCompanyEvent(*user.CompanyID, "TRANSACTION_FILE_UPLOADED", map[string]interface{}{"transaction_file_id": createdFile.ID, "records_created": recordsCreated, "uploaded_by": userID, "file_path_url": url})
  return createdFile, nil
}

func (s *appSvc) UpdateTransactionFileName(ctx context.Context, userID, fileID uuid.UUID, newName string) error {
  if newName == "" {
    return fmt.Errorf("newName is empty")
  }
  user, err := s.usvc.GetUserByID(userID)
  if err != nil {
    return err
  }
  tf, err := s.tfsvc.GetTransactionFileByID(fileID)
  if err != nil {
    return err
  }
  if tf.CompanyID == nil || user.CompanyID == nil || *tf.CompanyID != *user.CompanyID {
    return fmt.Errorf("transaction file does not belong to user's company")
  }
  if err := s.tfsvc.UpdateTransactionFileName(fileID, newName); err != nil {
    return err
  }
  _ = s.pub.PublishCompanyEvent(*tf.CompanyID, "TRANSACTION_FILE_RENAMED", map[string]interface{}{"file_id": fileID, "new_file_name": newName, "updated_by": userID})
  return nil
}

func (s *appSvc) DeleteTransactionFile(ctx context.Context, userID, fileID uuid.UUID) error {
  user, err := s.usvc.GetUserByID(userID)
  if err != nil {
    return err
  }
  tf, err := s.tfsvc.GetTransactionFileByID(fileID)
  if err != nil {
    return err
  }
  if tf.CompanyID == nil || user.CompanyID == nil || *tf.CompanyID != *user.CompanyID {
    return fmt.Errorf("transaction file does not belong to user's company")
  }
  if err := s.tfsvc.DeleteTransactionFile(fileID); err != nil {
    return fmt.Errorf("failed to delete transaction file: %w", err)
  }
  _ = s.pub.PublishCompanyEvent(*tf.CompanyID, "TRANSACTION_FILE_DELETED", map[string]interface{}{"file_id": fileID, "deleted_by": userID})
  return nil
}

func (s *appSvc) ListTransactionFiles(ctx context.Context, userID uuid.UUID, f repos.TransactionFileFilter) ([]*models.TransactionFile, error) {
  user, err := s.usvc.GetUserByID(userID)
  if err != nil || user == nil {
    return nil, fmt.Errorf("unauthorized user")
  }
  if user.CompanyID == nil {
    return nil, fmt.Errorf("user has no company")
  }
  f.CompanyID = *user.CompanyID
  return s.tfsvc.ListTransactionFiles(f)
}

func (s *appSvc) extractExt(name string) string {
  idx := -1
  for i := len(name) - 1; i >= 0; i-- {
    if name[i] == '.' {
      idx = i
      break
    }
  }
  if idx == -1 || idx == len(name)-1 {
    return ""
  }
  return name[idx:]
}

func (s *appSvc) CreateTransactionRecord(ctx context.Context, userID, warehouseID uuid.UUID, transactionType, orderName, description string, transactionQ, completedQ int, completedDate time.Time, locationPath, locationNamePath, itemName string) (*models.TransactionRecord, error) {
  user, err := s.usvc.GetUserByID(userID)
  if err != nil || user == nil || user.CompanyID == nil {
    return nil, fmt.Errorf("invalid user or user has no company")
  }
  wh, err := s.GetWarehouseByID(ctx, userID, warehouseID)
  if err != nil {
    return nil, fmt.Errorf("warehouse invalid: %w", err)
  }
  loc, err := s.lsvc.GetLocationByPath(*wh.WarehouseID, locationPath)
  if err != nil || loc == nil {
    newLoc := models.Location{
      WarehouseID:      &wh.ID,
      LocationPath:     locationPath,
      LocationNamePath: locationNamePath,
    }
    loc, err := s.lsvc.CreateLocation(newLoc)
    if err != nil {
      return nil, fmt.Errorf("could not create location: %w", err)
    }
  }
  item, err := s.isvc.GetItemByNameAndCompanyID(*user.CompanyID, itemName)
  if err != nil || item == nil {
    i := models.Item{
      CompanyID:      user.CompanyID,
      Name:           itemName,
    }
    item, err := s.isvc.CreateItem(i)
    if err != nil {
      return nil, fmt.Errorf("could not create item: %w", err)
    }
    _ = s.wsvc.LinkToItem(wh.ID, item.ID)
  }
  _ = s.lsvc.LinkToItem(loc.ID, item.ID)

  rec := models.TransactionRecord{
    CompanyID:            user.CompanyID,
    WarehouseID:          &wh.ID,
    LocationID:           &loc.ID,
    ItemID:               &item.ID,
    TransactionType:      transactionType,
    OrderName:            orderName,
    Description:          description,
    TransactionQuantity:  transactionQ,
    CompletedDate:        completedDate,
    CompletedQuantity:    completedQ,
  }
  createdRec, err := s.trsvc.CreateTransactionRecord(rec)
  if err != nil {
    return nil, fmt.Errorf("failed to create transaction record: %w", err)
  }
  _ = s.pub.PublishCompanyEvent(*user.CompanyID, "TRANSACTION_RECORD_CREATED", map[string]interface{}{"record_id": createdRec.ID, "warehouse_id": wh.ID, "location_id": loc.ID, "item_id": item.ID, "user_id": userID})
  return createdRec, nil
}

func (s *appSvc) GetTransactionRecordByID(ctx context.Context, userID, recordID uuid.UUID) (*models.TransactionRecord, error) {
  user, err := s.usvc.GetUserByID(userID)
  if err != nil {
    return nil, fmt.Errorf("unauthorized user")
  }
  rec, err := s.trsvc.GetTransactionRecordByID(recordID)
  if err != nil {
    return nil, err
  }
  if rec.CompanyID == nil || user.CompanyID == nil || *rec.CompanyID != *user.CompanyID {
    return nil, fmt.Errorf("record is not in user's company")
  }
  return rec, nil
}

func (s *appSvc) UpdateTransactionRecordOrderName(ctx context.Context, userID, recordID uuid.UUID, newOrderName string) error {
  rec, err := s.GetTransactionRecordByID(ctx, userID, recordID)
  if err != nil {
    return err
  }
  if err := s.trsvc.UpdateTransactionRecordOrderName(recordID, newOrderName); err != nil {
    return err
  }
  _ = s.pub.PublishCompanyEvent(*rec.CompanyID, "TRANSACTION_RECORD_UPDATED", map[string]interface{}{"record_id": recordID, "updated_by": userID, "order_name": newOrderName})
  return nil
}

func (s *appSvc) UpdateTransactionRecordDescription(ctx context.Context, userID, recordID uuid.UUID, newDescription string) error {
  rec, err := s.GetTransactionRecordByID(ctx, userID, recordID)
  if err != nil {
    return err
  }
  if err := s.trsvc.UpdateTransactionRecordDescription(recordID, newDescription); err != nil {
    return err
  }
  _ = s.pub.PublishCompanyEvent(*rec.CompanyID, "TRANSACTION_RECORD_UPDATED", map[string]interface{}{"record_id": recordID, "updated_by": userID, "description": newDescription})
  return nil
}

func (s *appSvc) UpdateTransactionRecordTransactionQuantity(ctx context.Context, userID, recordID uuid.UUID, newTQuantity int64) error {
  rec, err := s.GetTransactionRecordByID(ctx, userID, recordID)
  if err != nil {
    return err
  }
  if err := s.trsvc.UpdateTransactionRecordTransactionQuantity(recordID, newTQuantity); err != nil {
    return err
  }
  _ = s.pub.PublishCompanyEvent(*rec.CompanyID, "TRANSACTION_RECORD_UPDATED", map[string]interface{}{"record_id": recordID, "updated_by": userID, "transaction_quantity": newTQuantity})
  return nil
}

func (s *appSvc) UpdateTransactionRecordCompletedQuantity(ctx context.Context, userID, recordID uuid.UUID, newCQuantity int64) error {
  rec, err := s.GetTransactionRecordByID(ctx, userID, recordID)
  if err != nil {
    return err
  }
  if err := s.trsvc.UpdateTransactionRecordCompletedQuantity(recordID, newCQuantity); err != nil {
    return err
  }
  _ = s.pub.PublishCompanyEvent(*rec.CompanyID, "TRANSACTION_RECORD_UPDATED", map[string]interface{}{"record_id": recordID, "updated_by": userID, "completed_quantity": newCQuantity})
  return nil
}

func (s *appSvc) UpdateTransactionRecordCompletedDate(ctx context.Context, userID, recordID uuid.UUID, newDate time.Time) error {
  rec, err := s.GetTransactionRecordByID(ctx, userID, recordID)
  if err != nil {
    return err
  }
  if err := s.trsvc.UpdateTransactionRecordCompletedDate(recordID, newDate); err != nil {
    return err
  }
  _ = s.pub.PublishCompanyEvent(*rec.CompanyID, "TRANSACTION_RECORD_UPDATED", map[string]interface{}{"record_id": recordID, "updated_by": userID, "completed_date": newDate})
  return nil
}

func (s *appSvc) UpdateTransactionRecordTransactionType(ctx context.Context, userID, recordID uuid.UUID, newType string) error {
  if newType == "" {
    return fmt.Errorf("transaction type is required")
  }
  rec, err := s.GetTransactionRecordByID(ctx, userID, recordID)
  if err != nil {
    return err
  }
  if err := s.trsvc.UpdateTransactionRecordTransactionType(recordID, newType); err != nil {
    return err
  }
  _ = s.pub.PublishCompanyEvent(*rec.CompanyID, "TRANSACTION_RECORD_UPDATED", map[string]interface{}{"record_id": recordID, "updated_by": userID, "transaction_type": newType})
  return nil
}

func (s *appSvc) ListTransactionRecords(ctx context.Context, userID uuid.UUID, f repos.TransactionRecordFilter) ([]*models.TransactionRecord, error) {
  user, err := s.usvc.GetUserByID(userID)
  if err != nil || user == nil {
    return nil, fmt.Errorf("unauthorized user")
  }
  if user.CompanyID == nil {
    return nil, fmt.Errorf("user has no company")
  }
  f.CompanyID = *user.CompanyID
  return s.trsvc.ListTransactionRecords(f)
}

func (s *appSvc) UpdateUserAvatar(ctx context.Context, userID uuid.UUID, newAvatar []byte) (string, error) {
  user, err := s.usvc.GetUserByID(userID)
  if err != nil {
    return "", fmt.Errorf("user not found: %w", err)
  }
  if len(newAvatar) == 0 {
    return "", fmt.Errorf("avatar image is empty")
  }
  fileName := fmt.Sprintf("user-avatar-%s.png", userID)
  contentType := "image/png"
  url, err := s.s3svc.UploadImage(ctx, fileName, newAvatar, contentType)
  if err != nil {
    return "", fmt.Errorf("failed to upload user avatar: %w", err)
  }
  if err := s.usvc.UpdateUserAvatarURL(userID, url); err != nil {
    return "", fmt.Errorf("failed to update user avatar URL: %w", err)
  }
  if user.CompanyID != nil {
    _ = s.pub.PublishCompanyEvent(*user.CompanyID, "USER_AVATAR_UPDATED", map[string]interface{}{"user_id": userID, "avatar_url": url})
  }
  _ = s.pub.PublishUserEvent(userID, "USER_AVATAR_UPDATED", map[string]interface{}{"user_id": userID, "avatar_url": url})
  return nil
}

func (s *appSvc) UpdateUserFirstName(ctx context.Context, userID uuid.UUID, newFirst string) error {
  if newFirst == "" {
    return fmt.Errorf("first name is required")
  }
  err := s.usvc.UpdateUserFirstName(userID, newFirst)
  if err != nil {
    return fmt.Errorf("failed to update user first name: %w", err)
  }
  user, err := s.usvc.GetUserByID(userID)
  if err != nil {
    return err
  }
  _ = s.pub.PublishCompanyEvent(*user.CompanyID, "USER_FIRST_NAME_UPDATED", map[string]interface{}{"user_id": userID, "first_name": newFirst})
  _ = s.pub.PublishUserEvent(userID, "USER_FIRST_NAME_UPDATED", map[string]interface{}{"user_id": userID, "first_name": newFirst})
  return nil
}

func (s *appSvc) UpdateUserLastName(ctx context.Context, userID uuid.UUID, newLast string) error {
  if newLast == "" {
    return fmt.Errorf("last name is required")
  }
  err := s.usvc.UpdateUserLastName(userID, newLast)
  if err != nil {
    return fmt.Errorf("failed to update user last name: %w", err)
  }
  user, err := s.usvc.GetUserByID(userID)
  if err != nil {
    return err
  }
  _ = s.pub.PublishCompanyEvent(*user.CompanyID, "USER_LAST_NAME_UPDATED", map[string]interface{}{"user_id": userID, "last_name": newLast})
  _ = s.pub.PublishUserEvent(userID, "USER_LAST_NAME_UPDATED", map[string]interface{}{"user_id": userID, "last_name": newLast})
  return nil
}

func (s *appSvc) UpdateUserEmail(ctx context.Context, userID uuid.UUID, newEmail string) error {
  if newEmail == "" {
    return fmt.Errorf("new email is empty")
  }
  returned, err := s.usvc.GetUserByEmail(newEmail)
  if returned != nil {
    return fmt.Errorf("email already in use")
  }
  err := s.usvc.UpdateUserEmail(userID, newEmail)
  if err != nil {
    return fmt.Errorf("failed to update user email: %w", err)
  }
  return nil
}

func (s *appSvc) UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPass string) error {
  if newPass == "" {
    return fmt.Errorf("new password is empty")
  }
  hashed, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
  if err != nil {
    return fmt.Errorf("failed to hash new password: %w", err)
  }
  err = s.usvc.UpdateUserPassword(targetUserID, string(hashed))
  if err != nil {
    return fmt.Errorf("failed to update user password: %w", err)
  }
  _ = s.pub.PublishUserEvent(userID, "USER_PASSWORD_UPDATED", map[string]interface{}{"user_id": userID, "password": newPass})
  return nil
}

func (s *appSvc) DeleteUser(ctx context.Context, userID, targetUserID uuid.UUID) error {
  if userID != targetUserID {
    return fmt.Errorf("permission denied to delete another user")
  }
  user, err := s.usvc.GetUserByID(targetUserID)
  if err != nil || user == nil {
    return fmt.Errorf("failed to find user: %w", err)
  }
  if err := s.usvc.DeleteUser(targetUserID); err != nil {
    return fmt.Errorf("failed to delete user: %w", err)
  }
  _ = s.pub.PublishCompanyEvent(*user.CompanyID, "USER_DELETED", map[string]interface{}{"user_id": user.ID})
  return nil
}

func (s *appSvc) ListUsers(ctx context.Context, userID uuid.UUID, f repos.UserFilter) ([]*models.User, error) {
  user, err := s.usvc.GetUserByID(userID)
  if err != nil {
    return nil, fmt.Errorf("unauthorized user")
  }
  if user.CompanyID == nil {
    return nil, fmt.Errorf("this user has no company")
  }
  f.CompanyID = *user.CompanyID
  return s.usvc.ListUsers(f)
}

func (s *appSvc) ListItems(ctx context.Context, userID uuid.UUID, f repos.ItemFilter) ([]*models.Item, error) {
  user, err := s.usvc.GetUserByID(userID)
  if err != nil || user == nil {
    return nil, fmt.Errorf("invalid user")
  }
  if user.CompanyID == nil {
    return nil, fmt.Errorf("user has no company")
  }
  f.CompanyID = *user.CompanyID
  return s.isvc.ListItems(f)
}

func (s *appSvc) generateUserAvatar(ctx context.Context, firstName, lastName string) (string, error) {
  seed := fmt.Sprintf("%s-%s", strings.ToLower(firstName), strings.ToLower(lastName))
  return s.avatarsvc.GenerateAndUploadAvatar(ctx, "adventurer", seed)
}

func (s *appSvc) generateCompanyAvatar(ctx context.Context, companyName string) (string, error) {
  seed := strings.ToLower(strings.ReplaceAll(companyName, " ", "-"))
  return s.avatarsvc.GenerateAndUploadAvatar(ctx, "bottts", seed)
}

