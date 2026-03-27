package backend

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"carwash-app/backend/database"
	"carwash-app/backend/models"
)

// ============================================================
// PRINTER CONFIGURATION (stored in SQLite)
// ============================================================

type PrinterConfig struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Mode        string `json:"mode" gorm:"not null;default:browser"` // "browser", "network", "usb"
	NetworkIP   string `json:"networkIp"`                            // e.g. "192.168.1.100"
	NetworkPort string `json:"networkPort" gorm:"default:9100"`      // e.g. "9100"
	USBName     string `json:"usbName"`                              // Windows shared printer name
	PaperWidth  string `json:"paperWidth" gorm:"default:80"`         // "58" or "80" mm
	StoreName   string `json:"storeName" gorm:"default:SparkleWash"`
	StoreAddr   string `json:"storeAddr"`
	StorePhone  string `json:"storePhone"`
	FooterText  string `json:"footerText" gorm:"default:Terima Kasih!"`
	AutoPrint   bool   `json:"autoPrint" gorm:"default:false"` // auto print after order
}

// ============================================================
// RECEIPT DATA STRUCT
// ============================================================

type ReceiptData struct {
	InvoiceNo      string        `json:"invoiceNo"`
	Date           string        `json:"date"`
	CustomerName   string        `json:"customerName"`
	PlateNumber    string        `json:"plateNumber"`
	CarType        string        `json:"carType"`
	Items          []ReceiptItem `json:"items"`
	Subtotal       float64       `json:"subtotal"`
	DiscountAmount float64       `json:"discountAmount"`
	Total          float64       `json:"total"`
	PaymentMethod  string        `json:"paymentMethod"`
	CashierName    string        `json:"cashierName"`
}

type ReceiptItem struct {
	Name     string  `json:"name"`
	Qty      int     `json:"qty"`
	Price    float64 `json:"price"`
	Subtotal float64 `json:"subtotal"`
}

// ============================================================
// PRINTER SERVICE - Methods bound to Wails
// ============================================================

// GetPrinterConfig returns current printer config from DB
func (a *App) GetPrinterConfig() PrinterConfig {
	var cfg PrinterConfig
	result := database.DB.First(&cfg)
	if result.Error != nil {
		// Return defaults if not found
		cfg = PrinterConfig{
			Mode:       "browser",
			PaperWidth: "80",
			StoreName:  "SparkleWash",
			StoreAddr:  "",
			StorePhone: "",
			FooterText: "Terima Kasih!",
			AutoPrint:  false,
		}
	}
	return cfg
}

// SavePrinterConfig saves printer settings
func (a *App) SavePrinterConfig(mode, networkIP, networkPort, usbName, paperWidth, storeName, storeAddr, storePhone, footerText string, autoPrint bool) (PrinterConfig, error) {
	var cfg PrinterConfig
	database.DB.First(&cfg)

	cfg.Mode = mode
	cfg.NetworkIP = networkIP
	cfg.NetworkPort = networkPort
	cfg.USBName = usbName
	cfg.PaperWidth = paperWidth
	cfg.StoreName = storeName
	cfg.StoreAddr = storeAddr
	cfg.StorePhone = storePhone
	cfg.FooterText = footerText
	cfg.AutoPrint = autoPrint

	if cfg.ID == 0 {
		database.DB.Create(&cfg)
	} else {
		database.DB.Save(&cfg)
	}

	return cfg, nil
}

// TestPrinter sends a test print
func (a *App) TestPrinter() error {
	cfg := a.GetPrinterConfig()

	testData := ReceiptData{
		InvoiceNo:    "TEST-001",
		Date:         time.Now().Format("02/01/2006 15:04"),
		CustomerName: "Test Customer",
		PlateNumber:  "B 1234 XX",
		CarType:      "Test Car",
		Items: []ReceiptItem{
			{Name: "Cuci Standar", Qty: 1, Price: 55000, Subtotal: 55000},
			{Name: "Es Teh Manis", Qty: 2, Price: 8000, Subtotal: 16000},
		},
		Subtotal:       71000,
		DiscountAmount: 0,
		Total:          71000,
		PaymentMethod:  "Cash",
		CashierName:    "Admin",
	}

	return a.doPrint(cfg, testData)
}

// PrintReceipt prints a receipt from transaction ID
func (a *App) PrintReceipt(txID uint) error {
	var tx models.Transaction
	if err := database.DB.Preload("Items").Preload("Cashier").First(&tx, txID).Error; err != nil {
		return fmt.Errorf("transaksi tidak ditemukan")
	}

	cfg := a.GetPrinterConfig()

	var items []ReceiptItem
	for _, item := range tx.Items {
		items = append(items, ReceiptItem{
			Name:     item.ItemName,
			Qty:      item.Quantity,
			Price:    item.Price,
			Subtotal: item.Subtotal,
		})
	}

	data := ReceiptData{
		InvoiceNo:      tx.InvoiceNo,
		Date:           tx.CreatedAt.Format("02/01/2006 15:04"),
		CustomerName:   tx.CustomerName,
		PlateNumber:    tx.PlateNumber,
		CarType:        tx.CarType,
		Items:          items,
		Subtotal:       tx.Subtotal,
		DiscountAmount: tx.DiscountAmount,
		Total:          tx.Total,
		PaymentMethod:  tx.PaymentMethod,
		CashierName:    tx.Cashier.FullName,
	}

	return a.doPrint(cfg, data)
}

// PrintReceiptDirect prints from frontend data (for immediate print after order)
func (a *App) PrintReceiptDirect(jsonData string) error {
	var data ReceiptData
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return fmt.Errorf("data struk tidak valid: %v", err)
	}
	cfg := a.GetPrinterConfig()
	return a.doPrint(cfg, data)
}

// ============================================================
// INTERNAL PRINT DISPATCHER
// ============================================================

func (a *App) doPrint(cfg PrinterConfig, data ReceiptData) error {
	switch cfg.Mode {
	case "network":
		return a.printNetwork(cfg, data)
	case "usb":
		return a.printUSB(cfg, data)
	default:
		return a.printBrowser(cfg, data)
	}
}

// ============================================================
// MODE 1: BROWSER PRINT (Universal - works everywhere)
// ============================================================

func (a *App) printBrowser(cfg PrinterConfig, data ReceiptData) error {
	html := a.generateReceiptHTML(cfg, data)

	tmpFile := filepath.Join(os.TempDir(), "sparklewash_receipt.html")
	if err := os.WriteFile(tmpFile, []byte(html), 0644); err != nil {
		return fmt.Errorf("gagal menulis file struk: %v", err)
	}

	// Try opening in default browser (cross-platform)
	// Windows
	if err := exec.Command("cmd", "/c", "start", "", tmpFile).Start(); err == nil {
		return nil
	}
	// macOS
	if err := exec.Command("open", tmpFile).Start(); err == nil {
		return nil
	}
	// Linux
	if err := exec.Command("xdg-open", tmpFile).Start(); err == nil {
		return nil
	}

	return fmt.Errorf("tidak bisa membuka browser untuk cetak")
}

// ============================================================
// MODE 2: NETWORK TCP (ESC/POS Direct - Fast, no dialog)
// ============================================================

func (a *App) printNetwork(cfg PrinterConfig, data ReceiptData) error {
	addr := cfg.NetworkIP + ":" + cfg.NetworkPort
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("gagal konek ke printer %s: %v", addr, err)
	}
	defer conn.Close()

	escpos := a.generateESCPOS(cfg, data)
	_, err = conn.Write(escpos)
	if err != nil {
		return fmt.Errorf("gagal kirim data ke printer: %v", err)
	}

	return nil
}

// ============================================================
// MODE 3: USB / Windows Shared Printer (ESC/POS via file handle)
// ============================================================

func (a *App) printUSB(cfg PrinterConfig, data ReceiptData) error {
	printerPath := cfg.USBName

	// On Windows, shared printers are accessed like: \\localhost\PrinterName
	if !strings.HasPrefix(printerPath, `\\`) && !strings.HasPrefix(printerPath, "/dev/") {
		printerPath = `\\localhost\` + printerPath
	}

	f, err := os.OpenFile(printerPath, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("gagal buka printer '%s': %v. Pastikan printer sudah di-share di Windows.", cfg.USBName, err)
	}
	defer f.Close()

	escpos := a.generateESCPOS(cfg, data)
	_, err = f.Write(escpos)
	if err != nil {
		return fmt.Errorf("gagal kirim data ke printer: %v", err)
	}

	return nil
}

// ============================================================
// ESC/POS COMMAND BUILDER (for network & USB modes)
// ============================================================

func (a *App) generateESCPOS(cfg PrinterConfig, data ReceiptData) []byte {
	var b []byte

	// Determine column width based on paper
	cols := 48 // 80mm
	if cfg.PaperWidth == "58" {
		cols = 32
	}

	// Initialize printer
	b = append(b, 0x1B, 0x40) // ESC @ (init)

	// Center align
	b = append(b, 0x1B, 0x61, 0x01) // ESC a 1

	// Bold on + Double height for store name
	b = append(b, 0x1B, 0x45, 0x01) // ESC E 1 (bold)
	b = append(b, 0x1D, 0x21, 0x11) // GS ! 0x11 (double width+height)
	b = append(b, []byte(cfg.StoreName+"\n")...)
	b = append(b, 0x1D, 0x21, 0x00) // GS ! 0x00 (normal size)
	b = append(b, 0x1B, 0x45, 0x00) // ESC E 0 (bold off)

	// Store info
	if cfg.StoreAddr != "" {
		b = append(b, []byte(cfg.StoreAddr+"\n")...)
	}
	if cfg.StorePhone != "" {
		b = append(b, []byte(cfg.StorePhone+"\n")...)
	}

	// Separator
	b = append(b, []byte(strings.Repeat("-", cols)+"\n")...)

	// Left align
	b = append(b, 0x1B, 0x61, 0x00) // ESC a 0

	// Transaction info
	b = append(b, []byte(fmt.Sprintf("No   : %s\n", data.InvoiceNo))...)
	b = append(b, []byte(fmt.Sprintf("Tgl  : %s\n", data.Date))...)
	b = append(b, []byte(fmt.Sprintf("Nama : %s\n", data.CustomerName))...)
	b = append(b, []byte(fmt.Sprintf("Plat : %s\n", data.PlateNumber))...)
	if data.CarType != "" && data.CarType != "-" {
		b = append(b, []byte(fmt.Sprintf("Mobil: %s\n", data.CarType))...)
	}
	b = append(b, []byte(fmt.Sprintf("Kasir: %s\n", data.CashierName))...)

	// Separator
	b = append(b, []byte(strings.Repeat("-", cols)+"\n")...)

	// Items
	for _, item := range data.Items {
		nameStr := item.Name
		if len(nameStr) > cols-2 {
			nameStr = nameStr[:cols-2]
		}
		b = append(b, []byte(nameStr+"\n")...)

		qtyPrice := fmt.Sprintf("  %dx %s", item.Qty, formatRp(item.Price))
		subStr := formatRp(item.Subtotal)
		spaces := cols - len(qtyPrice) - len(subStr)
		if spaces < 1 {
			spaces = 1
		}
		b = append(b, []byte(qtyPrice+strings.Repeat(" ", spaces)+subStr+"\n")...)
	}

	// Separator
	b = append(b, []byte(strings.Repeat("-", cols)+"\n")...)

	// Totals
	printRow := func(label string, amount string) {
		spaces := cols - len(label) - len(amount)
		if spaces < 1 {
			spaces = 1
		}
		b = append(b, []byte(label+strings.Repeat(" ", spaces)+amount+"\n")...)
	}

	printRow("Subtotal", formatRp(data.Subtotal))

	if data.DiscountAmount > 0 {
		printRow("Diskon", "-"+formatRp(data.DiscountAmount))
	}

	// Bold for total
	b = append(b, 0x1B, 0x45, 0x01)
	b = append(b, 0x1D, 0x21, 0x01) // double height
	printRow("TOTAL", formatRp(data.Total))
	b = append(b, 0x1D, 0x21, 0x00)
	b = append(b, 0x1B, 0x45, 0x00)

	printRow("Bayar", strings.ToUpper(data.PaymentMethod))

	// Separator
	b = append(b, []byte(strings.Repeat("-", cols)+"\n")...)

	// Footer - center
	b = append(b, 0x1B, 0x61, 0x01)
	b = append(b, []byte(cfg.FooterText+"\n")...)
	b = append(b, []byte("\n")...)

	// Feed and cut
	b = append(b, []byte("\n\n\n")...)
	b = append(b, 0x1D, 0x56, 0x00) // GS V 0 (full cut)

	return b
}

// ============================================================
// HTML RECEIPT GENERATOR (for browser mode)
// ============================================================

func (a *App) generateReceiptHTML(cfg PrinterConfig, data ReceiptData) string {
	width := "76mm"
	pageSize := "80mm auto"
	if cfg.PaperWidth == "58" {
		width = "54mm"
		pageSize = "58mm auto"
	}

	var itemsHTML string
	for _, item := range data.Items {
		itemsHTML += fmt.Sprintf(`<tr>
			<td>%s</td>
			<td class="c">%d</td>
			<td class="r">%s</td>
			<td class="r">%s</td>
		</tr>`, item.Name, item.Qty, formatRp(item.Price), formatRp(item.Subtotal))
	}

	discountRow := ""
	if data.DiscountAmount > 0 {
		discountRow = fmt.Sprintf(`<tr><td colspan="3">Diskon</td><td class="r">-%s</td></tr>`, formatRp(data.DiscountAmount))
	}

	carInfo := ""
	if data.CarType != "" && data.CarType != "-" {
		carInfo = fmt.Sprintf("<div>Mobil : %s</div>", data.CarType)
	}

	storeAddrHTML := ""
	if cfg.StoreAddr != "" {
		storeAddrHTML = fmt.Sprintf(`<div class="center small">%s</div>`, cfg.StoreAddr)
	}
	storePhoneHTML := ""
	if cfg.StorePhone != "" {
		storePhoneHTML = fmt.Sprintf(`<div class="center small">%s</div>`, cfg.StorePhone)
	}

	// Use backtick-safe approach: no literal </script> in Go string
	scriptTag := "<" + "script>window.onload=function(){window.print();window.onafterprint=function(){window.close()}}</" + "script>"

	html := fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>Struk %s</title>
<style>
@page{size:%s;margin:0}
@media print{body{margin:0}html{background:#fff}}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:'Courier New',monospace;font-size:12px;width:%s;margin:0 auto;padding:3mm 2mm;background:#fff;color:#000}
.center{text-align:center}
.small{font-size:10px}
.bold{font-weight:bold}
.sep{border-top:1px dashed #000;margin:5px 0}
table{width:100%%;border-collapse:collapse}
td{padding:1px 0;font-size:11px;vertical-align:top}
.c{text-align:center}
.r{text-align:right}
.total-row td{font-size:14px;font-weight:bold;padding-top:4px}
.info div{font-size:11px;line-height:1.4}
</style></head>
<body>
<div class="center bold" style="font-size:16px">%s</div>
%s%s
<div class="sep"></div>
<div class="info">
<div>No   : %s</div>
<div>Tgl  : %s</div>
<div>Nama : %s</div>
<div>Plat : %s</div>
%s
<div>Kasir: %s</div>
</div>
<div class="sep"></div>
<table>
<tr class="bold"><td>Item</td><td class="c">Qty</td><td class="r">Harga</td><td class="r">Sub</td></tr>
%s
</table>
<div class="sep"></div>
<table>
<tr><td colspan="3">Subtotal</td><td class="r">%s</td></tr>
%s
<tr class="total-row"><td colspan="3">TOTAL</td><td class="r">%s</td></tr>
<tr><td colspan="3">Bayar</td><td class="r">%s</td></tr>
</table>
<div class="sep"></div>
<div class="center" style="margin-top:6px;font-size:11px">%s</div>
<div class="center small" style="margin-top:8px;color:#999">SparkleWash POS</div>
%s
</body></html>`,
		data.InvoiceNo,
		pageSize, width,
		cfg.StoreName,
		storeAddrHTML, storePhoneHTML,
		data.InvoiceNo, data.Date,
		data.CustomerName, data.PlateNumber,
		carInfo, data.CashierName,
		itemsHTML,
		formatRp(data.Subtotal),
		discountRow,
		formatRp(data.Total),
		strings.ToUpper(data.PaymentMethod),
		cfg.FooterText,
		scriptTag,
	)

	return html
}

// ============================================================
// HELPER
// ============================================================

func formatRp(amount float64) string {
	// Format as "Rp 55.000"
	str := fmt.Sprintf("%.0f", amount)
	n := len(str)
	if n <= 3 {
		return "Rp " + str
	}

	var result []byte
	for i, ch := range str {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, byte(ch))
	}
	return "Rp " + string(result)
}