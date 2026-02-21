package models

import (
	"time"
)

// ============================================================
// User & Auth
// ============================================================

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"` // hashed, never sent to frontend
	FullName  string    `json:"fullName" gorm:"not null"`
	Role      string    `json:"role" gorm:"not null;default:kasir"` // "superadmin" or "kasir"
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt" ts_type:"string"`
	UpdatedAt time.Time `json:"updatedAt" ts_type:"string"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}

// ============================================================
// Shift Management
// ============================================================

type Shift struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userId" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Date      string    `json:"date" gorm:"not null"`      // "2026-02-21"
	StartTime string    `json:"startTime" gorm:"not null"` // "08:00"
	EndTime   string    `json:"endTime" gorm:"not null"`   // "16:00"
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt" ts_type:"string"`
	UpdatedAt time.Time `json:"updatedAt" ts_type:"string"`
}

// ============================================================
// Car Wash Packages
// ============================================================

type WashPackage struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null"`
	Duration    int       `json:"duration"` // minutes
	Category    string    `json:"category"` // "cuci", "detailing", "coating"
	Icon        string    `json:"icon"`
	IsActive    bool      `json:"isActive" gorm:"default:true"`
	CreatedAt   time.Time `json:"createdAt" ts_type:"string"`
	UpdatedAt   time.Time `json:"updatedAt" ts_type:"string"`
}

// ============================================================
// Food & Beverage Menu
// ============================================================

type MenuItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	Category  string    `json:"category"` // "makanan", "minuman", "snack"
	Icon      string    `json:"icon"`
	Stock     int       `json:"stock" gorm:"default:0"`
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt" ts_type:"string"`
	UpdatedAt time.Time `json:"updatedAt" ts_type:"string"`
}

// ============================================================
// Discounts
// ============================================================

type Discount struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"not null"`
	Type       string    `json:"type" gorm:"not null"`      // "percentage" or "fixed"
	Value      float64   `json:"value" gorm:"not null"`     // e.g. 10 for 10% or 5000 for Rp5000
	MinOrder   float64   `json:"minOrder" gorm:"default:0"` // minimum order to apply
	IsActive   bool      `json:"isActive" gorm:"default:true"`
	ValidFrom  string    `json:"validFrom"`  // "2026-01-01"
	ValidUntil string    `json:"validUntil"` // "2026-12-31"
	CreatedAt  time.Time `json:"createdAt" ts_type:"string"`
	UpdatedAt  time.Time `json:"updatedAt" ts_type:"string"`
}

// ============================================================
// Transactions
// ============================================================

type Transaction struct {
	ID             uint              `json:"id" gorm:"primaryKey"`
	InvoiceNo      string            `json:"invoiceNo" gorm:"uniqueIndex;not null"`
	CustomerName   string            `json:"customerName"`
	PlateNumber    string            `json:"plateNumber"`
	CarType        string            `json:"carType"`
	Items          []TransactionItem `json:"items" gorm:"foreignKey:TransactionID"`
	Subtotal       float64           `json:"subtotal"`
	DiscountID     *uint             `json:"discountId"`
	Discount       *Discount         `json:"discount,omitempty" gorm:"foreignKey:DiscountID"`
	DiscountAmount float64           `json:"discountAmount"`
	Total          float64           `json:"total"`
	PaymentMethod  string            `json:"paymentMethod"`
	Status         string            `json:"status" gorm:"default:washing"` // "washing", "done", "paid", "cancelled"
	CashierID      uint              `json:"cashierId"`
	Cashier        User              `json:"cashier" gorm:"foreignKey:CashierID"`
	CreatedAt      time.Time         `json:"createdAt" ts_type:"string"`
	UpdatedAt      time.Time         `json:"updatedAt" ts_type:"string"`
	CompletedAt    *time.Time        `json:"completedAt" ts_type:"string"`
}

type TransactionItem struct {
	ID            uint    `json:"id" gorm:"primaryKey"`
	TransactionID uint    `json:"transactionId" gorm:"not null"`
	ItemType      string  `json:"itemType" gorm:"not null"` // "package" or "menu"
	ItemID        uint    `json:"itemId" gorm:"not null"`
	ItemName      string  `json:"itemName" gorm:"not null"`
	Price         float64 `json:"price" gorm:"not null"`
	Quantity      int     `json:"quantity" gorm:"not null;default:1"`
	Subtotal      float64 `json:"subtotal" gorm:"not null"`
}

// ============================================================
// Dashboard Stats
// ============================================================

type DashboardStats struct {
	TodayRevenue      float64 `json:"todayRevenue"`
	TodayTransactions int64   `json:"todayTransactions"`
	ActiveWashes      int64   `json:"activeWashes"`
	FoodOrdersToday   int     `json:"foodOrdersToday"`
	MonthlyRevenue    float64 `json:"monthlyRevenue"`
}
