# Belajar Format JSON — Garuda Auto Service

File ini dibuat untuk membantu memahami format JSON yang dipakai aplikasi,
tanpa harus langsung menghafal semuanya. Buka file
[`contoh-format-json.json`](./contoh-format-json.json) di samping file ini,
lalu baca penjelasannya pelan-pelan.

> Semua data di file contoh adalah data fiktif. Jangan menyimpan data asli
> konsumen di dalam file contoh yang diunggah ke Git.

## 1. JSON itu apa?

JSON (*JavaScript Object Notation*) adalah cara menulis data menggunakan
teks biasa. Hampir semua REST API memakai JSON karena mudah dibaca manusia
dan mudah diproses komputer.

Contoh paling sederhana:

```json
{
  "nama": "Budi Santoso",
  "umur": 35,
  "aktif": true
}
```

Aturan yang paling sering membuat error:

1. Nama field dan teks **wajib** memakai tanda kutip dua (`"`), bukan kutip satu.
2. Angka **tidak** memakai tanda kutip: `45000`, bukan `"45000"`.
3. `true`, `false`, dan `null` ditulis tanpa tanda kutip.
4. Antar-field dipisahkan koma, tetapi **field terakhir tidak boleh** diberi koma.
5. JSON **tidak mendukung komentar** seperti `//` atau `#`.
6. Satu file JSON hanya boleh memiliki satu nilai paling luar (biasanya objek `{}` atau array `[]`).

## 2. Dua bentuk data yang sering dipakai

### Objek `{ }`

Dipakai untuk satu data yang memiliki nama-nama field.

```json
{
  "name": "Budi Santoso",
  "car_brand": "Toyota"
}
```

### Array `[ ]`

Dipakai untuk daftar banyak data.

```json
[
  { "name": "Budi Santoso" },
  { "name": "Siti Aminah" }
]
```

Di aplikasi ini, satu konsumen adalah **objek**, sedangkan daftar konsumen
adalah **array of object** (array yang berisi objek-objek).

## 3. Alur data di aplikasi

| Kebutuhan | Method | Endpoint | Bagian di file contoh |
|---|---|---|---|
| Login dan mendapat token | `POST` | `/api/auth/login` | Belum dicontohkan di file JSON |
| Menambah konsumen | `POST` | `/api/customers` | `1_tambah_konsumen` |
| Menambah riwayat servis | `POST` | `/api/customers/:id/services` | `2_tambah_riwayat_servis` |
| Melihat daftar konsumen | `GET` | `/api/customers` | `3_ambil_daftar_konsumen` |
| Melihat diagram bulanan | `GET` | `/api/stats/monthly?months=12` | `4_statistik_bulanan` |

Semua endpoint di atas memerlukan header:

```json
{
  "Authorization": "Bearer <token-login-anda>",
  "Content-Type": "application/json"
}
```

`<token-login-anda>` didapat dari respons login. Ganti bagian tersebut dengan
token asli saat mencoba lewat aplikasi lain (misalnya Postman).

## 4. Arti setiap field konsumen

| Field | Tipe | Wajib? | Arti |
|---|---|---|---|
| `id` | angka | Hanya di respons | Nomor unik konsumen dari database |
| `name` | teks | Ya | Nama konsumen |
| `phone` | teks | Tidak | Nomor telepon / WhatsApp |
| `address` | teks | Tidak | Alamat konsumen |
| `plate_number` | teks | Ya | Plat nomor kendaraan, contoh `B 1234 XYZ` |
| `car_brand` | teks | Ya | Merk mobil, contoh `Toyota` |
| `car_type` | teks | Ya | Tipe mobil, contoh `Avanza 1.3 G` |
| `services` | array | Hanya di respons | Daftar riwayat servis kendaraan |
| `created_at` | teks tanggal | Hanya di respons | Waktu data dibuat |
| `updated_at` | teks tanggal | Hanya di respons | Waktu data terakhir diubah |

Saat **menambah** konsumen, cukup kirim field yang wajib dan field opsional
yang sudah diketahui. Jangan mengirim `id`, `created_at`, atau `updated_at`
karena itu dibuat otomatis oleh server.

## 5. Arti setiap field riwayat servis

| Field | Tipe | Wajib? | Arti |
|---|---|---|---|
| `id` | angka | Hanya di respons | Nomor unik catatan servis |
| `customer_id` | angka | Hanya di respons | ID konsumen pemilik kendaraan |
| `service_date` | teks `YYYY-MM-DD` | Ya | Tanggal servis dikerjakan |
| `mileage` | angka | Tidak | Kilometer kendaraan saat servis |
| `description` | teks | Ya | Rincian pengerjaan servis |
| `next_service_date` | teks `YYYY-MM-DD` | Tidak | Perkiraan tanggal servis berikutnya |
| `next_mileage` | angka | Tidak | Perkiraan kilometer servis berikutnya |
| `notes` | teks | Tidak | Catatan tambahan atau keluhan |
| `created_at` | teks tanggal | Hanya di respons | Waktu catatan dibuat |
| `updated_at` | teks tanggal | Hanya di respons | Waktu catatan terakhir diubah |

Perhatikan: `service_date` dan `next_service_date` adalah **teks**, bukan
angka, karena formatnya harus `YYYY-MM-DD`. Contoh yang benar:
`"2026-10-28"`.

## 6. Arti data diagram bulanan

Respons `/api/stats/monthly` berbentuk:

```json
{
  "start_month": "2026-08",
  "end_month": "2026-10",
  "months": [
    {
      "month": "2026-10",
      "label": "Okt",
      "year": 2026,
      "new_customers": 3,
      "services": 5
    }
  ],
  "totals": {
    "new_customers": 13,
    "services": 21
  }
}
```

| Field | Arti |
|---|---|
| `start_month` | Bulan pertama pada diagram, format `YYYY-MM` |
| `end_month` | Bulan terakhir pada diagram (bulan berjalan) |
| `months` | Daftar data per bulan, urut dari bulan terlama |
| `month` | Kunci bulan, format `YYYY-MM` |
| `label` | Nama bulan singkat untuk ditampilkan di diagram |
| `year` | Tahun dari bulan tersebut |
| `new_customers` | Jumlah konsumen baru yang dibuat pada bulan itu |
| `services` | Jumlah riwayat servis dengan tanggal pada bulan itu |
| `totals` | Jumlah seluruh konsumen baru dan servis pada rentang tersebut |

Parameter `?months=12` berarti ambil 12 bulan terakhir. Nilai yang diizinkan
adalah 1 sampai 24. Jika tidak diisi, server memakai 12 bulan.

## 7. Latihan bertahap

Lakukan satu langkah per hari agar tidak kewalahan.

### Hari 1 — Membaca objek

Buka bagian `1_tambah_konsumen` pada file JSON. Tunjukkan mana yang disebut
nama field dan mana yang disebut nilai. Perhatikan bahwa `"name"` memakai
tanda kutip, sedangkan angka tidak.

### Hari 2 — Membedakan request dan response

Di bagian yang sama, bandingkan `request_body` dengan `response_201`.
Apa saja field yang kita kirim? Apa saja field tambahan yang muncul dari
server?

### Hari 3 — Mengenal array

Buka bagian `3_ambil_daftar_konsumen`. Perhatikan tanda `[` dan `]` yang
mengelilingi data konsumen. Itu artinya respons berisi **banyak** konsumen.
Di dalamnya ada `services` yang juga berupa array.

### Hari 4 — Mencoba lewat PowerShell

Login terlebih dahulu untuk mendapatkan token:

```powershell
$login = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/auth/login" `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"username":"staf","password":"password123"}'

$token = $login.token
```

Lalu kirim data konsumen:

```powershell
$body = @{
  name         = "Budi Santoso"
  phone        = "081234567890"
  address      = "Jl. Merdeka No. 10"
  plate_number = "B 1234 XYZ"
  car_brand    = "Toyota"
  car_type     = "Avanza 1.3 G"
} | ConvertTo-Json

Invoke-RestMethod `
  -Uri "http://localhost:8080/api/customers" `
  -Method Post `
  -Headers @{ Authorization = "Bearer $token" } `
  -ContentType "application/json" `
  -Body $body
```

### Hari 5 — Melihat diagram dari API

```powershell
Invoke-RestMethod `
  -Uri "http://localhost:8080/api/stats/monthly?months=12" `
  -Headers @{ Authorization = "Bearer $token" }
```

Bandingkan hasilnya dengan diagram yang tampil di dashboard.

## 8. Kesalahan JSON yang paling sering terjadi

| Kesalahan | Contoh salah | Contoh benar |
|---|---|---|
| Kutip satu | `{'name': 'Budi'}` | `{"name": "Budi"}` |
| Koma di akhir | `{"name": "Budi",}` | `{"name": "Budi"}` |
| Komentar | `// data konsumen` | Hapus komentarnya |
| Angka diberi kutip | `"mileage": "45000"` | `"mileage": 45000` |
| Lupa kutip pada teks | `"name": Budi` | `"name": "Budi"` |
| Tanggal salah format | `"2026-10-28"` ditulis `28-10-2026` | `"2026-10-28"` |

Jika muncul pesan error dari aplikasi, baca field `error` pada respons.
Contohnya:

```json
{
  "error": "Nama, plat nomor, merk mobil, dan tipe mobil wajib diisi"
}
```

Itu berarti ada field wajib yang belum dikirim atau namanya salah ketik.
