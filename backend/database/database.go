package database

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"carwash-app/backend/models"
)

var DB *gorm.DB

func Initialize() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot get home dir: %w", err)
	}

	dataDir := filepath.Join(homeDir, ".sparklewash")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("cannot create data dir: %w", err)
	}

	dbPath := filepath.Join(dataDir, "sparklewash.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("cannot open database: %w", err)
	}

	// Enable WAL mode for better concurrent access
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA foreign_keys=ON")

	// Auto migrate
	err = db.AutoMigrate(
		&models.User{},
		&models.Shift{},
		&models.WashPackage{},
		&models.MenuItem{},
		&models.Discount{},
		&models.Transaction{},
		&models.TransactionItem{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	DB = db

	// Seed default data if empty
	seedDefaults()

	return nil
}

func seedDefaults() {
	// Seed superadmin if no users exist
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		DB.Create(&models.User{
			Username: "admin",
			Password: string(hash),
			FullName: "Super Admin",
			Role:     "superadmin",
			IsActive: true,
		})

		hash2, _ := bcrypt.GenerateFromPassword([]byte("kasir123"), bcrypt.DefaultCost)
		DB.Create(&models.User{
			Username: "kasir1",
			Password: string(hash2),
			FullName: "Kasir Satu",
			Role:     "kasir",
			IsActive: true,
		})
	}

	// Seed packages if empty
	var pkgCount int64
	DB.Model(&models.WashPackage{}).Count(&pkgCount)
	if pkgCount == 0 {
		packages := []models.WashPackage{
			{Name: "Cuci Kilat", Description: "Cuci exterior cepat, vakum interior", Price: 35000, Duration: 20, Category: "cuci", Icon: "⚡", IsActive: true},
			{Name: "Cuci Standar", Description: "Cuci exterior + interior + semir ban", Price: 55000, Duration: 35, Category: "cuci", Icon: "🚿", IsActive: true},
			{Name: "Cuci Premium", Description: "Full wash + wax + poles body + vakum detail", Price: 95000, Duration: 60, Category: "cuci", Icon: "✨", IsActive: true},
			{Name: "Cuci Salon", Description: "Detailing full + engine wash + parfum", Price: 150000, Duration: 90, Category: "detailing", Icon: "💎", IsActive: true},
			{Name: "Nano Coating", Description: "Perlindungan cat nano ceramic", Price: 350000, Duration: 120, Category: "coating", Icon: "🛡️", IsActive: true},
			{Name: "Cuci Motor", Description: "Cuci motor + semir body", Price: 20000, Duration: 15, Category: "cuci", Icon: "🏍️", IsActive: true},
		}
		DB.Create(&packages)
	}

	// Seed menu items if empty
	var menuCount int64
	DB.Model(&models.MenuItem{}).Count(&menuCount)
	if menuCount == 0 {
		items := []models.MenuItem{
			{Name: "Nasi Goreng Spesial", Price: 25000, Category: "makanan", Icon: "🍛", Stock: 50, IsActive: true},
			{Name: "Mie Goreng", Price: 22000, Category: "makanan", Icon: "🍜", Stock: 50, IsActive: true},
			{Name: "Ayam Geprek", Price: 28000, Category: "makanan", Icon: "🍗", Stock: 30, IsActive: true},
			{Name: "Nasi Padang", Price: 30000, Category: "makanan", Icon: "🍚", Stock: 30, IsActive: true},
			{Name: "Bakso Urat", Price: 20000, Category: "makanan", Icon: "🍲", Stock: 40, IsActive: true},
			{Name: "Sate Ayam (10 tusuk)", Price: 25000, Category: "makanan", Icon: "🥘", Stock: 30, IsActive: true},
			{Name: "Es Teh Manis", Price: 8000, Category: "minuman", Icon: "🧊", Stock: 100, IsActive: true},
			{Name: "Es Jeruk", Price: 10000, Category: "minuman", Icon: "🍊", Stock: 80, IsActive: true},
			{Name: "Kopi Susu", Price: 15000, Category: "minuman", Icon: "☕", Stock: 80, IsActive: true},
			{Name: "Jus Alpukat", Price: 18000, Category: "minuman", Icon: "🥤", Stock: 50, IsActive: true},
			{Name: "Air Mineral", Price: 5000, Category: "minuman", Icon: "💧", Stock: 200, IsActive: true},
			{Name: "Kentang Goreng", Price: 15000, Category: "snack", Icon: "🍟", Stock: 60, IsActive: true},
			{Name: "Pisang Goreng", Price: 12000, Category: "snack", Icon: "🍌", Stock: 40, IsActive: true},
			{Name: "Roti Bakar", Price: 15000, Category: "snack", Icon: "🍞", Stock: 40, IsActive: true},
		}
		DB.Create(&items)
	}

	// Seed discount if empty
	var discCount int64
	DB.Model(&models.Discount{}).Count(&discCount)
	if discCount == 0 {
		discounts := []models.Discount{
			{Name: "Diskon Member 10%", Type: "percentage", Value: 10, MinOrder: 50000, IsActive: true, ValidFrom: "2026-01-01", ValidUntil: "2026-12-31"},
			{Name: "Potongan Rp 10.000", Type: "fixed", Value: 10000, MinOrder: 100000, IsActive: true, ValidFrom: "2026-01-01", ValidUntil: "2026-12-31"},
			{Name: "Promo Weekend 15%", Type: "percentage", Value: 15, MinOrder: 75000, IsActive: true, ValidFrom: "2026-01-01", ValidUntil: "2026-12-31"},
		}
		DB.Create(&discounts)
	}
}
