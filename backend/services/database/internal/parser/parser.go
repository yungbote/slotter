package parsing

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"

	"github.com/yungbote/slotter/backend/services/database/internal/models"
	"github.com/yungbote/slotter/backend/services/database/internal/services"
)

type ParserService interface {
	// ParseFile reads a file from disk or memory buffer, extracts location/item info,
	// creates them if needed, links them, and creates transaction records referencing that fileID.
	// Returns how many transaction records were created, or error if any step fails.
	ParseFile(filePath string, transactionFileID, companyID, warehouseID uuid.UUID) (int, error)
}

type parserService struct {
	lsvc  services.LSvc
	isvc  services.ISvc
	trsvc services.TRSvc
	tfsvc services.TFSvc
	wsvc  services.WSvc // optional if you want to link items to the warehouse
}

// Ensure we only treat certain columns as known transaction columns, and the rest as location columns.
var knownTransactionCols = map[string]bool{
	"id":                   true,
	"transaction type":     true,
	"order number":         true,
	"item number":          true,
	"description":          true,
	"transaction quantity": true,
	"completed date":       true,
	"completed by":         true,
	"completed quantity":   true,
}

// For caching location + item references
type locationCache struct {
	LocationPath     string
	LocationNamePath string
	ID               uuid.UUID
}

type itemCache struct {
	Name string
	ID   uuid.UUID
}

type transactionRow struct {
	TransactionType     string
	OrderName           string
	Description         string
	TransactionQuantity int
	CompletedQuantity   int
	CompletedDate       *time.Time
	LocationPathKey     string
	ItemNameKey         string
}

// NewParserService is the constructor. Pass in your existing Svc interfaces.
func NewParserService(
	lsvc services.LSvc,
	isvc services.ISvc,
	trsvc services.TRSvc,
	tfsvc services.TFSvc,
	wsvc services.WSvc,
) ParserService {
	return &parserService{
		lsvc:  lsvc,
		isvc:  isvc,
		trsvc: trsvc,
		tfsvc: tfsvc,
		wsvc:  wsvc,
	}
}

// ParseFile is the main entry point. It guesses file type from ext, calls parseCSV or parseXLSX, etc.
func (p *parserService) ParseFile(
	filePath string,
	transactionFileID, companyID, warehouseID uuid.UUID,
) (int, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".csv" {
		return p.parseCSVAndFlush(filePath, transactionFileID, companyID, warehouseID)
	} else if ext == ".xlsx" {
		return p.parseXLSXAndFlush(filePath, transactionFileID, companyID, warehouseID)
	} else if ext == ".xls" {
		return 0, fmt.Errorf(".xls parsing not implemented")
	}
	return 0, fmt.Errorf("unsupported file extension: %s", ext)
}

// parseCSVAndFlush handles CSV reading line-by-line, then calls flushToDB.
func (p *parserService) parseCSVAndFlush(
	filePath string,
	transactionFileID, companyID, warehouseID uuid.UUID,
) (int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	var header []string
	lineIndex := 0

	// Temporary in-memory caches
	locationMap := make(map[string]*locationCache)
	itemMap := make(map[string]*itemCache)
	var rows []transactionRow

	for {
		line, errRead := reader.ReadString('\n')
		if errRead != nil && errRead != io.EOF {
			return 0, fmt.Errorf("csv read error: %w", errRead)
		}
		if len(line) == 0 && errRead == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		cols := strings.Split(line, ",")
		if lineIndex == 0 {
			// This is the header row
			for i, c := range cols {
				cols[i] = strings.ToLower(strings.TrimSpace(c))
			}
			header = cols
		} else {
			// Body row
			if len(cols) != len(header) {
				// skip or handle error
				continue
			}
			rowMap := make(map[string]string)
			for i, h := range header {
				rowMap[h] = strings.TrimSpace(cols[i])
			}
			locCols := determineLocationCols(header)
			handleRow(rowMap, locCols, locationMap, itemMap, &rows)
		}
		lineIndex++
		if errRead == io.EOF {
			break
		}
	}

	return flushToDB(p, transactionFileID, companyID, warehouseID, locationMap, itemMap, rows)
}

// parseXLSXAndFlush handles XLSX reading with excelize, then calls flushToDB.
func (p *parserService) parseXLSXAndFlush(
	filePath string,
	transactionFileID, companyID, warehouseID uuid.UUID,
) (int, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("cannot open XLSX: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return 0, errors.New("xlsx has no sheets")
	}

	rows, err := f.Rows(sheetName)
	if err != nil {
		return 0, fmt.Errorf("failed to get rows: %w", err)
	}

	locationMap := make(map[string]*locationCache)
	itemMap := make(map[string]*itemCache)
	var txRows []transactionRow

	var header []string
	rowIndex := 0

	for rows.Next() {
		rowCells, errRow := rows.Columns()
		if errRow != nil {
			return 0, fmt.Errorf("xlsx row read error: %w", errRow)
		}

		if rowIndex == 0 {
			// This is header
			header = make([]string, len(rowCells))
			for i, cell := range rowCells {
				header[i] = strings.ToLower(strings.TrimSpace(cell))
			}
		} else {
			rowMap := make(map[string]string)
			for i, h := range header {
				val := ""
				if i < len(rowCells) {
					val = strings.TrimSpace(rowCells[i])
				}
				rowMap[h] = val
			}
			locCols := determineLocationCols(header)
			handleRow(rowMap, locCols, locationMap, itemMap, &txRows)
		}
		rowIndex++
	}

	return flushToDB(p, transactionFileID, companyID, warehouseID, locationMap, itemMap, txRows)
}

// determineLocationCols picks columns not in the knownTransactionCols set.
func determineLocationCols(header []string) []string {
	var locCols []string
	for _, h := range header {
		if !knownTransactionCols[h] {
			locCols = append(locCols, h)
		}
	}
	return locCols
}

// handleRow extracts the transaction row from rowMap and populates locationMap/itemMap as needed.
func handleRow(
	rowMap map[string]string,
	locCols []string,
	locationMap map[string]*locationCache,
	itemMap map[string]*itemCache,
	txRows *[]transactionRow,
) {
	locPath, locNamePath := buildLocationPath(rowMap, locCols)

	// Known columns
	itemName := strings.TrimSpace(rowMap["item number"])
	tranType := strings.TrimSpace(rowMap["transaction type"])
	orderName := strings.TrimSpace(rowMap["order number"])
	desc := strings.TrimSpace(rowMap["description"])

	qtyStr := strings.TrimSpace(rowMap["transaction quantity"])
	if qtyStr == "" {
		qtyStr = "0"
	}
	qty := parseInt(qtyStr)

	compQtyStr := strings.TrimSpace(rowMap["completed quantity"])
	if compQtyStr == "" {
		compQtyStr = "0"
	}
	compQty := parseInt(compQtyStr)

	dateStr := strings.TrimSpace(rowMap["completed date"])
	dateVal := parseDate(dateStr)

	// Cache location
	if _, exists := locationMap[locPath]; !exists {
		locationMap[locPath] = &locationCache{
			LocationPath:     locPath,
			LocationNamePath: locNamePath,
			ID:               uuid.Nil,
		}
	}

	// Cache item
	if _, exists := itemMap[itemName]; !exists {
		itemMap[itemName] = &itemCache{
			Name: itemName,
			ID:   uuid.Nil,
		}
	}

	// Accumulate row
	*txRows = append(*txRows, transactionRow{
		TransactionType:     tranType,
		OrderName:           orderName,
		Description:         desc,
		TransactionQuantity: qty,
		CompletedQuantity:   compQty,
		CompletedDate:       dateVal,
		LocationPathKey:     locPath,
		ItemNameKey:         itemName,
	})
}

// buildLocationPath forms slash-delimited path plus a pipe-delimited name path from leftover columns.
func buildLocationPath(rowMap map[string]string, locCols []string) (string, string) {
	var pathParts []string
	var nameParts []string
	for _, c := range locCols {
		val := strings.TrimSpace(rowMap[c])
		if val != "" {
			pathParts = append(pathParts, val)
			nameParts = append(nameParts, fmt.Sprintf("%s=%s", strings.Title(c), val))
		}
	}
	return strings.Join(pathParts, "/"), strings.Join(nameParts, "|")
}

// parseInt tries to parse an int from string
func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

// parseDate tries "2006-01-02" layout
func parseDate(s string) *time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return &t
	}
	return nil
}

// flushToDB handles creation of location/item records and transactionRecords in DB
func flushToDB(
	p *parserService,
	transactionFileID, companyID, warehouseID uuid.UUID,
	locationMap map[string]*locationCache,
	itemMap map[string]*itemCache,
	rows []transactionRow,
) (int, error) {
	// 1) Resolve/create location records
	for _, loc := range locationMap {
		if loc.ID != uuid.Nil {
			continue
		}
		// Check if location already exists
		existing, err := p.lsvc.GetLocationByPath(companyID, warehouseID, loc.LocationPath)
		if err == nil && existing != nil {
			loc.ID = existing.ID
		} else {
			// create
			lModel := models.Location{
				WarehouseID:      &warehouseID,
				LocationPath:     loc.LocationPath,
				LocationNamePath: loc.LocationNamePath,
			}
			created, errCreate := p.lsvc.CreateLocation(lModel)
			if errCreate != nil {
				return 0, fmt.Errorf("failed to create location '%s': %w", loc.LocationPath, errCreate)
			}
			loc.ID = created.ID
		}
	}

	// 2) Resolve/create item records
	for _, itm := range itemMap {
		if itm.ID != uuid.Nil {
			continue
		}
		// Attempt to find existing
		existing, err := p.isvc.GetByItemNameAndCompanyID(companyID, itm.Name)
		if err == nil && existing != nil {
			itm.ID = existing.ID
		} else {
			// Create new item
			newItem := models.Item{
				Name:      itm.Name,
				CompanyID: &companyID,
			}
			created, errCreate := p.isvc.CreateItem(newItem)
			if errCreate != nil {
				return 0, fmt.Errorf("failed to create item '%s': %w", itm.Name, errCreate)
			}
			itm.ID = created.ID

			// Optionally link item to warehouse
			if errLink := p.wsvc.LinkToItem(warehouseID, itm.ID); errLink != nil {
				// not fatal, but might indicate a mismatch
				return 0, fmt.Errorf("failed to link item '%s' to warehouse '%s': %w", itm.Name, warehouseID, errLink)
			}
		}
	}

	// 3) Link location<->item for each transaction row
	linkedSet := make(map[string]bool)
	for _, row := range rows {
		loc := locationMap[row.LocationPathKey]
		itm := itemMap[row.ItemNameKey]
		if loc.ID == uuid.Nil || itm.ID == uuid.Nil {
			continue
		}
		key := fmt.Sprintf("%s_%s", loc.ID.String(), itm.ID.String())
		if !linkedSet[key] {
			if err := p.lsvc.LinkToItem(loc.ID, itm.ID); err != nil {
				return 0, fmt.Errorf("failed to link location '%s' with item '%s': %w", loc.LocationPath, itm.Name, err)
			}
			linkedSet[key] = true
		}
	}

	// 4) Create transaction records
	createdCount := 0
	for _, row := range rows {
		loc := locationMap[row.LocationPathKey]
		itm := itemMap[row.ItemNameKey]
		if loc.ID == uuid.Nil || itm.ID == uuid.Nil {
			continue
		}

		rec := models.TransactionRecord{
			CompanyID:          &companyID,
			WarehouseID:        &warehouseID,
			LocationID:         &loc.ID,
			TransactionFileID:  &transactionFileID,
			ItemID:             &itm.ID,
			TransactionType:    row.TransactionType,
			OrderName:          row.OrderName,
			Description:        row.Description,
			TransactionQuantity: row.TransactionQuantity,
			CompletedQuantity:   row.CompletedQuantity,
		}
		if row.CompletedDate != nil {
			rec.CompletedDate = *row.CompletedDate
		}

		_, err := p.trsvc.CreateTransactionRecord(rec)
		if err != nil {
			return createdCount, fmt.Errorf("failed to create transaction record: %w", err)
		}
		createdCount++
	}

	// 5) Link transaction file <-> items & locations
	for _, loc := range locationMap {
		if loc.ID != uuid.Nil {
			if err := p.tfsvc.LinkToLocation(transactionFileID, loc.ID); err != nil {
				return createdCount, fmt.Errorf("failed to link tf->location: %w", err)
			}
		}
	}
	for _, itm := range itemMap {
		if itm.ID != uuid.Nil {
			if err := p.tfsvc.LinkToItem(transactionFileID, itm.ID); err != nil {
				return createdCount, fmt.Errorf("failed to link tf->item: %w", err)
			}
		}
	}

	return createdCount, nil
}

