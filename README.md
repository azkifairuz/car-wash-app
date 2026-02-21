# SparkleWash v2 - Sistem Cuci Mobil & Cafe

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
