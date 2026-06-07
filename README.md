# 🚀 GLK Web App (Mega Project)

![Golang](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber-20232A?style=for-the-badge&logo=gofiber&logoColor=00ADD8)
![GORM](https://img.shields.io/badge/GORM-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white)
![Bootstrap](https://img.shields.io/badge/Bootstrap-563D7C?style=for-the-badge&logo=bootstrap&logoColor=white)

**GLK Web App** adalah sistem berskala besar (*Mega Project*) berbasis **Golang Fiber** yang dirancang khusus untuk memfasilitasi kebutuhan manajemen internal dan publik secara terintegrasi. Sistem ini dibagi menjadi beberapa portal multi-tenant yang dapat berjalan secara konkuren (seperti Portal Admin dan Portal Pelamar).

---

## ✨ Fitur Utama

- **🏗️ Multi-Portal Architecture**
  Aplikasi menjalankan dua server terpisah dalam satu instansi *binary*:
  - **Portal Pelamar** (Port `8081`) - Antarmuka publik.
  - **Portal Admin** (Port `8082`) - *Backend management dashboard*.
  
- **🔐 Dynamic RBAC (Role-Based Access Control)**
  - Otorisasi halaman dan menu secara terpusat.
  - Opsi perizinan spesifik hingga ke *sub-menu* yang disimpan di dalam *Database* (Tidak lagi bergantung pada *hardcoded JSON*).
  - Integrasi session lintas modular untuk perlindungan rute secara dinamis.

- **🧩 Mega Project Modular Tree**
  Semua menu dan sistem disiapkan untuk dapat beradaptasi dengan konsep *Mega Project*, dipisahkan atas modul-modul independen:
  - 💼 **Recruitment**
  - 🪪 **Admkar (Administrasi Karyawan)**
  - 💰 **Payroll**
  
- **⚡ Advanced Rendering & UI**
  - Menggunakan **Go HTML/Template** dengan fitur *Hot-Reload* (`Reload: true`).
  - Dilengkapi antarmuka responsif berbasis Bootstrap 5.
  - **Grid.js** digunakan untuk manajemen *Advance Tables* (Search, Sort, Pagination).

---

## 🛠️ Tech Stack

*   **Backend:** Go (Golang)
*   **Web Framework:** [Fiber v2](https://gofiber.io/)
*   **ORM:** [GORM](https://gorm.io/)
*   **Database:** MySQL
*   **Frontend UI:** Bootstrap 5, BoxIcons, Waves Effect, Grid.js
*   **Environment Management:** Godotenv

---

## 🚀 Instalasi & Menjalankan Aplikasi

Ikuti panduan berikut untuk menjalankan aplikasi di *environment* lokal Anda:

### 1. Persiapan Database
1. Pastikan layanan MySQL Anda telah berjalan.
2. Buat database kosong bernama `glk_db` (atau sesuai konfigurasi di `.env`).

### 2. Konfigurasi Environment (`.env`)
Buat file `.env` di *root* proyek (atau gandakan dari `.env.example` jika ada) dan sesuaikan konfigurasi koneksi databasenya:
```env
DB_DSN=user:password@tcp(127.0.0.1:3306)/glk_db?charset=utf8mb4&parseTime=True&loc=Local
PELAMAR_PORT=8081
ADMIN_PORT=8082
```

### 3. Instalasi Dependensi
Jalankan perintah berikut untuk mengunduh seluruh *package* Golang yang dibutuhkan:
```bash
go mod tidy
```

### 4. Jalankan Aplikasi
Jalankan file utama `main.go`. Aplikasi otomatis akan melakukan **Auto-Migrate** tabel dan melakukan **Seeding** (pembuatan *Super Admin* dan stuktur Menu/Modul) ke database.
```bash
go run main.go
```
*Tunggu hingga terminal menampilkan informasi bahwa server berhasil berjalan di port `8081` dan `8082`.*

---

## 🔑 Default Akses (Seeder)
Saat aplikasi pertama kali terhubung dengan database kosong, ia akan membuatkan akun default:

*   **URL Portal Admin:** `http://localhost:8082/admin/login`
*   **Username:** `admin`
*   **Password:** `admin`
*   **Role:** `Super Admin` (Memiliki izin penuh ke seluruh modul)

---

## 📂 Struktur Direktori Utama

```text
📁 glk-web-app
├── 📁 config/           # Konfigurasi Database dan Environment
├── 📁 controllers/      # Routing logic (Admin & Pelamar)
├── 📁 models/           # Definisi GORM schema & Entities (Admin, Master, Menu, dll)
├── 📁 static/           # Asset statis (CSS, JS, Images, Grid.js)
├── 📁 views/            # File HTML Templates (Base Layout, Admin, Auth, dll)
├── 📄 .env              # Environment Variables
├── 📄 main.go           # Entry point aplikasi (Fiber init, Auto-Migrate, Seeder)
└── 📄 README.md         # Dokumentasi Proyek
```

---

## 🎨 Modifikasi Tema
Warna utama (*primary color*) dikendalikan via *CSS variables* yang dapat Anda temukan pada tag `<style>` di `views/layouts/base.html`. 
Setelan saat ini adalah: `#F87242` (Oranye). Anda tidak perlu merekompilasi CSS bawaan (*app.min.css*) untuk mengubah warnanya.

---
*Developed with ❤️ for Guru Lesku Mega Project.*
