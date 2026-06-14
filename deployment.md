# Panduan Deployment GLK Web App

Dokumen ini berisi langkah-langkah untuk melakukan deployment aplikasi **GLK Web App** (Pelamar & Admin) ke server VPS menggunakan **Docker**.

## 1. Persiapan Server

Pastikan VPS Anda sudah terinstall:

- **Docker**
- **Docker Compose**
- **Git**

## 2. Struktur Project

Aplikasi ini terdiri dari dua layanan utama yang berjalan secara konkuren:

- **Web Pelamar**: Port `8081` (Default)
- **Web Admin**: Port `8082` (Default)
- **Database**: PostgreSQL
- **Cache**: Redis

## 3. Langkah-langkah Deployment

### A. Clone Repository

Masuk ke VPS Anda dan clone project ini:

```bash
git clone <url-repository-anda>
cd glk-web-app
```

### B. Konfigurasi Environment (`.env`)

Salin file `.env.example` menjadi `.env` dan sesuaikan nilainya:

```bash
cp .env.example .env
nano .env
```

Pastikan variabel berikut diisi dengan benar:

```env
# Database
DB_HOST=db
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=glk_db
DB_PORT=5432

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# SMTP (Untuk Magic Link)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=email-anda@gmail.com
SMTP_PASS=app-password-anda

# URL App (Sesuaikan dengan domain/IP VPS)
APP_URL=http://<ip-vps-anda>:8081
```

### C. Menjalankan Aplikasi dengan Docker Compose

Gunakan Docker Compose untuk membangun dan menjalankan seluruh kontainer:

```bash
# Build dan jalankan di background
docker-compose up -d --build
```

### D. Verifikasi Status

Cek apakah semua kontainer berjalan dengan baik:

```bash
docker-compose ps
```

Anda seharusnya melihat kontainer `app`, `db`, dan `redis` dalam status `Up`.

## 4. Akses Aplikasi

Setelah berhasil, aplikasi dapat diakses melalui:

- **Portal Pelamar**: `http://<IP-VPS>:8081`
- **Portal Admin**: `http://<IP-VPS>:8082/admin`

## 5. Perintah Penting Lainnya

### Melihat Log Aplikasi

Jika terjadi error, cek log kontainer:

```bash
docker-compose logs -f app
```

### Menghentikan Aplikasi

```bash
docker-compose down
```

### Update Aplikasi (Jika ada perubahan kode)

```bash
git pull origin main
docker-compose up -d --build
```

## 6. Keamanan (Opsional tapi Disarankan)

Jika Anda menggunakan Domain, disarankan untuk menggunakan **Nginx Reverse Proxy** dan **SSL (Certbot)** di depan kontainer Docker untuk keamanan (HTTPS) pada port 80/443.

---

**Gurulesku**
© 2026 GLK Web App Deployment Guide.
