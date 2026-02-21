# 🚗 SparkleWash v2 - Sistem Cuci Mobil & Cafe

Aplikasi desktop POS (Point of Sale) dengan **Wails v2** (Go + Vue 3) + **SQLite** database.

## ✨ Fitur

### Kasir
- **Dashboard** — Statistik pendapatan harian & bulanan, transaksi aktif
- **Order Baru** — Pilih paket cuci + makanan/minuman, pilih diskon, multi pembayaran
- **Antrian Aktif** — Monitor & update status cuci (Proses → Selesai → Lunas → Batal)
- **Riwayat** — Semua transaksi dengan filter status & detail modal

### Super Admin
- **CRUD Paket Cuci** — Tambah/edit/hapus paket, toggle aktif/nonaktif
- **CRUD Menu Makanan** — Kelola menu + stok, quick stock update
- **CRUD Diskon** — Persentase/nominal, min order, tanggal berlaku
- **Kelola User** — Buat akun kasir/admin, reset password, toggle aktif
- **Kelola Shift** — Atur jadwal shift kasir, filter per minggu

### Sistem
- **Login** — Autentikasi dengan bcrypt password hashing
- **Role-based access** — superadmin & kasir
- **SQLite database** — Persistent, file-based, zero config
- **Diskon** — Persentase atau nominal, dengan minimum order
- **Manajemen stok** — Auto decrease saat order, restore saat cancel

## 📦 Prasyarat

1. **Go** ≥ 1.21 — https://go.dev/dl/
2. **Node.js** ≥ 18 — https://nodejs.org/
3. **Wails CLI** v2:
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

## 🚀 Cara Menjalankan

```bash
cd carwash-app
cd frontend && npm install && cd ..
wails dev
```

### Build Production
```bash
wails build
# Output: build/bin/SparkleWash
```

## 🔑 Default Login

| Username | Password  | Role        |
|----------|-----------|-------------|
| admin    | admin123  | Super Admin |
| kasir1   | kasir123  | Kasir       |

## 📁 Struktur Project

```
carwash-app/
├── main.go
├── go.mod
├── wails.json
├── backend/
│   ├── app.go              # Semua handler (auth, CRUD, transaksi)
│   ├── database/
│   │   └── database.go     # SQLite init, migration, seeder
│   └── models/
│       └── models.go       # GORM models
└── frontend/
    ├── package.json
    ├── index.html
    └── src/
        ├── main.js          # Vue + Router + Auth guard
        ├── App.vue          # Layout + role-based sidebar
        ├── assets/style.css # Global dark theme
        ├── stores/auth.js   # Pinia auth + backend caller
        └── views/
            ├── Login.vue
            ├── Dashboard.vue
            ├── NewOrder.vue
            ├── ActiveOrders.vue
            ├── History.vue
            ├── ManagePackages.vue
            ├── ManageMenu.vue
            ├── ManageDiscounts.vue
            ├── ManageUsers.vue
            └── ManageShifts.vue
```

## 💾 Database

SQLite database disimpan di `~/.sparklewash/sparklewash.db`.
Auto migration — tabel dibuat otomatis saat pertama kali dijalankan.

## 🔮 Roadmap

- [ ] Laporan harian/bulanan (export PDF/Excel)
- [ ] Cetak struk (thermal printer)
- [ ] Member / loyalty program
- [ ] Notifikasi saat cuci selesai
- [ ] Multi-outlet support
