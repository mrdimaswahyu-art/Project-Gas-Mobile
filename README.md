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
