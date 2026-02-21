package backend

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"carwash-app/backend/database"
	"carwash-app/backend/models"
)

type App struct {
	ctx         context.Context
	currentUser *models.User
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	if err := database.Initialize(); err != nil {
		fmt.Println("Database init error:", err)
	}
}

func db() *gorm.DB {
	return database.DB
}

// ============================================================
// AUTH
// ============================================================

func (a *App) Login(username, password string) models.LoginResponse {
	var user models.User
	result := db().Where("username = ? AND is_active = ?", username, true).First(&user)
	if result.Error != nil {
		return models.LoginResponse{Success: false, Message: "Username tidak ditemukan"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return models.LoginResponse{Success: false, Message: "Password salah"}
	}

	a.currentUser = &user
	return models.LoginResponse{Success: true, Message: "Login berhasil", User: &user}
}

func (a *App) Logout() {
	a.currentUser = nil
}

func (a *App) GetCurrentUser() *models.User {
	return a.currentUser
}

// ============================================================
// USER MANAGEMENT (superadmin only)
// ============================================================

func (a *App) GetUsers() []models.User {
	var users []models.User
	db().Order("created_at DESC").Find(&users)
	return users
}

func (a *App) CreateUser(username, password, fullName, role string) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("gagal hash password")
	}

	user := models.User{
		Username: username,
		Password: string(hash),
		FullName: fullName,
		Role:     role,
		IsActive: true,
	}

	if result := db().Create(&user); result.Error != nil {
		return models.User{}, fmt.Errorf("username sudah dipakai")
	}
	return user, nil
}

func (a *App) UpdateUser(id uint, fullName, role string, isActive bool) (models.User, error) {
	var user models.User
	if result := db().First(&user, id); result.Error != nil {
		return models.User{}, fmt.Errorf("user tidak ditemukan")
	}
	user.FullName = fullName
	user.Role = role
	user.IsActive = isActive
	db().Save(&user)
	return user, nil
}

func (a *App) ResetPassword(id uint, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal hash password")
	}
	db().Model(&models.User{}).Where("id = ?", id).Update("password", string(hash))
	return nil
}

func (a *App) DeleteUser(id uint) error {
	return db().Delete(&models.User{}, id).Error
}

// ============================================================
// SHIFT MANAGEMENT
// ============================================================

func (a *App) GetShifts(dateFrom, dateTo string) []models.Shift {
	var shifts []models.Shift
	q := db().Preload("User").Order("date DESC, start_time ASC")
	if dateFrom != "" && dateTo != "" {
		q = q.Where("date BETWEEN ? AND ?", dateFrom, dateTo)
	} else if dateFrom != "" {
		q = q.Where("date >= ?", dateFrom)
	}
	q.Find(&shifts)
	return shifts
}

func (a *App) CreateShift(userId uint, date, startTime, endTime, notes string) (models.Shift, error) {
	shift := models.Shift{
		UserID:    userId,
		Date:      date,
		StartTime: startTime,
		EndTime:   endTime,
		Notes:     notes,
	}
	if result := db().Create(&shift); result.Error != nil {
		return models.Shift{}, result.Error
	}
	db().Preload("User").First(&shift, shift.ID)
	return shift, nil
}

func (a *App) UpdateShift(id uint, userId uint, date, startTime, endTime, notes string) (models.Shift, error) {
	var shift models.Shift
	if result := db().First(&shift, id); result.Error != nil {
		return models.Shift{}, fmt.Errorf("shift tidak ditemukan")
	}
	shift.UserID = userId
	shift.Date = date
	shift.StartTime = startTime
	shift.EndTime = endTime
	shift.Notes = notes
	db().Save(&shift)
	db().Preload("User").First(&shift, shift.ID)
	return shift, nil
}

func (a *App) DeleteShift(id uint) error {
	return db().Delete(&models.Shift{}, id).Error
}

// ============================================================
// PACKAGES CRUD
// ============================================================

func (a *App) GetPackages() []models.WashPackage {
	var pkgs []models.WashPackage
	db().Order("created_at ASC").Find(&pkgs)
	return pkgs
}

func (a *App) GetActivePackages() []models.WashPackage {
	var pkgs []models.WashPackage
	db().Where("is_active = ?", true).Order("price ASC").Find(&pkgs)
	return pkgs
}

func (a *App) CreatePackage(name, description string, price float64, duration int, category, icon string) (models.WashPackage, error) {
	pkg := models.WashPackage{
		Name: name, Description: description, Price: price,
		Duration: duration, Category: category, Icon: icon, IsActive: true,
	}
	if result := db().Create(&pkg); result.Error != nil {
		return models.WashPackage{}, result.Error
	}
	return pkg, nil
}

func (a *App) UpdatePackage(id uint, name, description string, price float64, duration int, category, icon string, isActive bool) (models.WashPackage, error) {
	var pkg models.WashPackage
	if result := db().First(&pkg, id); result.Error != nil {
		return models.WashPackage{}, fmt.Errorf("paket tidak ditemukan")
	}
	pkg.Name = name
	pkg.Description = description
	pkg.Price = price
	pkg.Duration = duration
	pkg.Category = category
	pkg.Icon = icon
	pkg.IsActive = isActive
	db().Save(&pkg)
	return pkg, nil
}

func (a *App) DeletePackage(id uint) error {
	return db().Delete(&models.WashPackage{}, id).Error
}

// ============================================================
// MENU ITEMS CRUD
// ============================================================

func (a *App) GetMenuItems() []models.MenuItem {
	var items []models.MenuItem
	db().Order("category ASC, name ASC").Find(&items)
	return items
}

func (a *App) GetActiveMenuItems() []models.MenuItem {
	var items []models.MenuItem
	db().Where("is_active = ? AND stock > 0", true).Order("category ASC, name ASC").Find(&items)
	return items
}

func (a *App) CreateMenuItem(name string, price float64, category, icon string, stock int) (models.MenuItem, error) {
	item := models.MenuItem{
		Name: name, Price: price, Category: category,
		Icon: icon, Stock: stock, IsActive: true,
	}
	if result := db().Create(&item); result.Error != nil {
		return models.MenuItem{}, result.Error
	}
	return item, nil
}

func (a *App) UpdateMenuItem(id uint, name string, price float64, category, icon string, stock int, isActive bool) (models.MenuItem, error) {
	var item models.MenuItem
	if result := db().First(&item, id); result.Error != nil {
		return models.MenuItem{}, fmt.Errorf("menu tidak ditemukan")
	}
	item.Name = name
	item.Price = price
	item.Category = category
	item.Icon = icon
	item.Stock = stock
	item.IsActive = isActive
	db().Save(&item)
	return item, nil
}

func (a *App) UpdateMenuStock(id uint, stock int) error {
	return db().Model(&models.MenuItem{}).Where("id = ?", id).Update("stock", stock).Error
}

func (a *App) DeleteMenuItem(id uint) error {
	return db().Delete(&models.MenuItem{}, id).Error
}

// ============================================================
// DISCOUNTS CRUD
// ============================================================

func (a *App) GetDiscounts() []models.Discount {
	var discounts []models.Discount
	db().Order("created_at DESC").Find(&discounts)
	return discounts
}

func (a *App) GetActiveDiscounts() []models.Discount {
	today := time.Now().Format("2006-01-02")
	var discounts []models.Discount
	db().Where("is_active = ? AND valid_from <= ? AND valid_until >= ?", true, today, today).Find(&discounts)
	return discounts
}

func (a *App) CreateDiscount(name, discType string, value, minOrder float64, validFrom, validUntil string) (models.Discount, error) {
	disc := models.Discount{
		Name: name, Type: discType, Value: value, MinOrder: minOrder,
		ValidFrom: validFrom, ValidUntil: validUntil, IsActive: true,
	}
	if result := db().Create(&disc); result.Error != nil {
		return models.Discount{}, result.Error
	}
	return disc, nil
}

func (a *App) UpdateDiscount(id uint, name, discType string, value, minOrder float64, validFrom, validUntil string, isActive bool) (models.Discount, error) {
	var disc models.Discount
	if result := db().First(&disc, id); result.Error != nil {
		return models.Discount{}, fmt.Errorf("diskon tidak ditemukan")
	}
	disc.Name = name
	disc.Type = discType
	disc.Value = value
	disc.MinOrder = minOrder
	disc.ValidFrom = validFrom
	disc.ValidUntil = validUntil
	disc.IsActive = isActive
	db().Save(&disc)
	return disc, nil
}

func (a *App) DeleteDiscount(id uint) error {
	return db().Delete(&models.Discount{}, id).Error
}

// ============================================================
// TRANSACTIONS
// ============================================================

type CartItem struct {
	ItemType string  `json:"itemType"` // "package" or "menu"
	ItemID   uint    `json:"itemId"`
	ItemName string  `json:"itemName"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func (a *App) CreateTransaction(customerName, plateNumber, carType, paymentMethod string, items []CartItem, discountId *uint) (models.Transaction, error) {
	if len(items) == 0 {
		return models.Transaction{}, fmt.Errorf("order harus memiliki minimal 1 item")
	}

	var subtotal float64
	var txItems []models.TransactionItem

	for _, ci := range items {
		itemSubtotal := ci.Price * float64(ci.Quantity)
		subtotal += itemSubtotal

		txItems = append(txItems, models.TransactionItem{
			ItemType: ci.ItemType,
			ItemID:   ci.ItemID,
			ItemName: ci.ItemName,
			Price:    ci.Price,
			Quantity: ci.Quantity,
			Subtotal: itemSubtotal,
		})

		// Decrease stock for menu items
		if ci.ItemType == "menu" {
			result := db().Model(&models.MenuItem{}).
				Where("id = ? AND stock >= ?", ci.ItemID, ci.Quantity).
				Update("stock", gorm.Expr("stock - ?", ci.Quantity))
			if result.RowsAffected == 0 {
				return models.Transaction{}, fmt.Errorf("stok %s tidak mencukupi", ci.ItemName)
			}
		}
	}

	// Calculate discount
	var discountAmount float64
	if discountId != nil && *discountId > 0 {
		var disc models.Discount
		if err := db().First(&disc, *discountId).Error; err == nil {
			if subtotal >= disc.MinOrder {
				if disc.Type == "percentage" {
					discountAmount = subtotal * disc.Value / 100
				} else {
					discountAmount = disc.Value
				}
			}
		}
	}

	total := subtotal - discountAmount
	if total < 0 {
		total = 0
	}

	cashierID := uint(0)
	if a.currentUser != nil {
		cashierID = a.currentUser.ID
	}

	now := time.Now()
	tx := models.Transaction{
		InvoiceNo:      fmt.Sprintf("INV-%s", now.Format("20060102150405")),
		CustomerName:   customerName,
		PlateNumber:    plateNumber,
		CarType:        carType,
		Items:          txItems,
		Subtotal:       subtotal,
		DiscountID:     discountId,
		DiscountAmount: discountAmount,
		Total:          total,
		PaymentMethod:  paymentMethod,
		Status:         "washing",
		CashierID:      cashierID,
	}

	if result := db().Create(&tx); result.Error != nil {
		return models.Transaction{}, result.Error
	}

	db().Preload("Items").Preload("Discount").Preload("Cashier").First(&tx, tx.ID)
	return tx, nil
}

func (a *App) GetTransactions(status, dateFrom, dateTo string) []models.Transaction {
	var txs []models.Transaction
	q := db().Preload("Items").Preload("Discount").Preload("Cashier").Order("created_at DESC")

	if status != "" && status != "all" {
		q = q.Where("status = ?", status)
	}
	if dateFrom != "" {
		q = q.Where("DATE(created_at) >= ?", dateFrom)
	}
	if dateTo != "" {
		q = q.Where("DATE(created_at) <= ?", dateTo)
	}

	q.Find(&txs)
	return txs
}

func (a *App) GetTodayTransactions() []models.Transaction {
	today := time.Now().Format("2006-01-02")
	return a.GetTransactions("", today, today)
}

func (a *App) UpdateTransactionStatus(id uint, status string) error {
	updates := map[string]interface{}{"status": status}
	if status == "done" || status == "paid" {
		now := time.Now()
		updates["completed_at"] = &now
	}
	return db().Model(&models.Transaction{}).Where("id = ?", id).Updates(updates).Error
}

func (a *App) CancelTransaction(id uint) error {
	var tx models.Transaction
	if err := db().Preload("Items").First(&tx, id).Error; err != nil {
		return fmt.Errorf("transaksi tidak ditemukan")
	}

	// Restore stock for menu items
	for _, item := range tx.Items {
		if item.ItemType == "menu" {
			db().Model(&models.MenuItem{}).Where("id = ?", item.ItemID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity))
		}
	}

	return db().Model(&tx).Update("status", "cancelled").Error
}

// ============================================================
// DASHBOARD
// ============================================================

func (a *App) GetDashboardStats(dateFrom, dateTo string) models.DashboardStats {
	today := time.Now().Format("2006-01-02")
	firstOfMonth := time.Now().Format("2006-01") + "-01"

	// If dates are provided, we use them for "today" stats but also for everything except monthlyRevenue
	// actually it's better to rename the fields or just use them as general stats.
	// We'll use dateFrom and dateTo for ALL stats except MonthlyRevenue which stays monthly.

	targetFrom := today
	targetTo := today
	if dateFrom != "" {
		targetFrom = dateFrom
	}
	if dateTo != "" {
		targetTo = dateTo
	}

	stats := models.DashboardStats{}

	// Revenue in period
	db().Model(&models.Transaction{}).
		Where("DATE(created_at) BETWEEN ? AND ? AND status IN ?", targetFrom, targetTo, []string{"done", "paid"}).
		Select("COALESCE(SUM(total), 0)").Scan(&stats.TodayRevenue)

	// Transactions in period
	db().Model(&models.Transaction{}).
		Where("DATE(created_at) BETWEEN ? AND ? AND status != ?", targetFrom, targetTo, "cancelled").
		Count((*int64)(&stats.TodayTransactions))

	// Active washes (not date bounded usually, but let's keep it as is)
	db().Model(&models.Transaction{}).
		Where("status = ?", "washing").
		Count((*int64)(&stats.ActiveWashes))

	// Food orders in period
	db().Model(&models.TransactionItem{}).
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Where("DATE(transactions.created_at) BETWEEN ? AND ? AND transaction_items.item_type = ? AND transactions.status != ?", targetFrom, targetTo, "menu", "cancelled").
		Select("COALESCE(SUM(transaction_items.quantity), 0)").Scan(&stats.FoodOrdersToday)

	// Monthly revenue (always current month for reference)
	db().Model(&models.Transaction{}).
		Where("DATE(created_at) >= ? AND status IN ?", firstOfMonth, []string{"done", "paid"}).
		Select("COALESCE(SUM(total), 0)").Scan(&stats.MonthlyRevenue)

	return stats
}

// ============================================================
// FILE EXPORT
// ============================================================

func (a *App) SaveExcelFile(base64Data string, defaultFilename string) (string, error) {
	filePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: defaultFilename,
		Title:           "Simpan Laporan Excel",
		Filters: []runtime.FileFilter{
			{DisplayName: "Excel Files (*.xlsx)", Pattern: "*.xlsx"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("gagal membuka dialog: %v", err)
	}
	if filePath == "" {
		return "", nil // user cancelled
	}

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("gagal decode data: %v", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("gagal menyimpan file: %v", err)
	}

	return filePath, nil
}
