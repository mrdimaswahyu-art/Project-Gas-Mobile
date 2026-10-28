# Garuda Auto Service — Web App

Aplikasi web manajemen pelanggan, data kendaraan, dan riwayat servis untuk bengkel **Garuda Auto Service**. Dibangun menggunakan Go (Gin), PostgreSQL (GORM), dan Frontend (HTML/CSS/JavaScript murni).

## Fitur Utama

1. **Autentikasi Staf:**
   - Registrasi & Login staf bengkel dengan hash password (bcrypt) dan JWT Bearer token.
2. **Pencarian Real-Time:**
   - Cari data berdasarkan Plat Nomor, Tipe Mobil, atau Nama Pelanggan.
3. **Filter Merk Mobil:**
   - Filter dropdown dinamis berdasarkan merk mobil yang ada di bengkel.
4. **Manajemen Konsumen & Kendaraan:**
   - Tambah, lihat detail, edit, dan hapus (CRUD) data konsumen serta kendaraannya.
5. **Manajemen Riwayat Servis:**
   - Catat tanggal servis, kilometer (KM), rincian pengerjaan servis, dan jadwal follow-up servis berikutnya (tanggal & KM).
   - Penghapusan otomatis (*cascade delete*) jika data konsumen dihapus.
6. **Diagram Konsumen Bulanan:**
   - Menampilkan jumlah konsumen baru dan riwayat servis per bulan selama 12 bulan terakhir.

## Belajar Format JSON

Bagi yang ingin mempelajari bentuk data API secara perlahan:

- [`docs/contoh-format-json.json`](docs/contoh-format-json.json) — contoh request dan response JSON yang valid.
- [`docs/belajar-json.md`](docs/belajar-json.md) — penjelasan field, aturan JSON, dan latihan bertahap.

Endpoint statistik bulanan: `GET /api/stats/monthly?months=12`
(nilai `months` yang diizinkan 1–24).

## Cara Menjalankan Aplikasi

1. Pastikan variabel lingkungan sudah diatur di file `.env` (salin dari `.env.example`):
   ```bash
   APP_PORT=8080
   JWT_SECRET=rahasia-kunci-jwt-anda
   DATABASE_URL=postgresql://username:password@localhost:5432/neondb?sslmode=require
   ```

2. Jalankan aplikasi:
   ```bash
   go run main.go
   ```

3. Buka browser di:
   - Dashboard: `http://localhost:8080/`
   - Login: `http://localhost:8080/login`
   - Registrasi: `http://localhost:8080/register`

## Menjalankan Pengujian (Testing)

```bash
# Pengujian koneksi & migrasi database
go test -v ./config

# Pengujian alur autentikasi dan CRUD data servis
go test -v ./handlers
```
