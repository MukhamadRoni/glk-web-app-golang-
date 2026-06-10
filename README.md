# GLK Web App (Golang)

Aplikasi Portal Rekrutmen Tutor LBB Gold Generation berbasis Golang (Fiber) dengan arsitektur microservices-lite (Dual Port) dan Database PostgreSQL. Aplikasi ini juga terintegrasi dengan Google Drive API via Apps Script dan Redis untuk manajemen _state_ tes akademik.

## 🌟 Ringkasan Fitur & Menu Saat Ini

### 👨‍💻 Portal Pelamar (Web Pelamar - Port 8081)
Aplikasi khusus untuk calon pelamar kerja.
1. **Login & Autentikasi**
   - Autentikasi menggunakan _Magic Link_ via Email (saat ini *bypass* langsung sukses untuk proses development).
2. **Dashboard Status Lamaran (`/dashboard`)**
   - Melihat histori aplikasi lamaran yang sudah dikirimkan beserta statusnya (Pending, Selesai Tes, Lulus, Ditolak, dll).
   - Menampilkan notifikasi "Lanjutkan Tes" jika pelamar memiliki Tes Akademik yang tertunda.
3. **Formulir Pendaftaran Guru (`/apply`)**
   - Pengisian biodata, domisili, jenjang yang dilamar, dan jadwal kesediaan mengajar.
   - Integrasi _Select2_ / _Choices.js_ yang terhubung ke _Master Data_ wilayah (Kota & Kecamatan) via API.
   - Fitur *Upload CV & Transkrip Nilai* langsung ke Google Drive via konektor Google Apps Script.
4. **Modul Tes Akademik (`/test/...`)**
   - **Intro (`/test/intro`)**: Halaman peraturan dan instruksi tes akademik.
   - **Ujian (`/test/soal`)**: Pengerjaan soal secara langsung. Menampilkan soal berdasarkan _Master Bank Soal_ (sesuai jenjang & mapel yang dilamar).
   - Fitur *Autosave* jawaban secara _real-time_ menggunakan **Redis**.
   - Fitur *Countdown Timer* cerdas yang otomatis mengumpulkan formulir jika waktu habis.

### 🛡️ Portal Admin (Web Admin - Port 8082)
Aplikasi _Back-Office_ untuk pengelolaan sistem dan Master Data oleh tim HR/Admin.
1. **Login Admin (`/admin/login`)**
   - Autentikasi standar menggunakan Username/Email dan Password.
2. **Dashboard HR (`/admin/dashboard`)**
   - Menampilkan metrik utama sistem (meskipun saat ini masih statis/_placeholder_).
3. **Data Pelamar (`/admin/applicants`)**
   - *Data Tables* yang menampilkan daftar semua pelamar yang mendaftar.
   - *Detail Pelamar*: Halaman profil detail untuk meninjau formulir pendaftaran, CV, Transkrip, jadwal luang, serta hasil tes pelamar.
4. **Modul Master Data (`/admin/recruitment/master/...`)**
   - **Master Wilayah (`/wilayah`)**: CRUD untuk tabel Kota dan Kecamatan di seluruh Indonesia.
   - **Master Jenis Pendidikan (`/jenis-pendidikan`)**: CRUD untuk tingkat/jenjang ajar (TK, SD, SMP, SMA, dsb).
   - **Master Mata Pelajaran (`/mapel`)**: CRUD untuk mengatur daftar bidang studi.
   - **Master Bank Soal (`/bank-soal`)**: Manajemen soal ujian (A, B, C) untuk tes akademik. Mendukung tipe soal _Multiple Choice_ dan pembatasan _Durasi Pengerjaan_ ujian.

---

## 🚀 Infrastruktur Pendukung
- **Framework**: Go Fiber v2
- **Database Utama**: PostgreSQL 16
- **Caching & Session Data**: Redis 7
- **UI Framework**: Bootstrap 5 + JQuery + MetisMenu (Template Skote)
- **Containerization**: Docker & Docker Compose

---

## 📋 To-Do List (Pengembangan Selanjutnya)

### Fitur Admin / HR
- [ ] **Sistem Penilaian (Scoring)**: Mengubah perhitungan nilai ujian dari *raw output* di log menjadi *Grade* atau presentase yang tersimpan di kolom *Database* khusus agar HRD dapat melihat nilai langsung di Dashboard.
- [ ] **Filter Data Pelamar**: Menambahkan filter canggih di halaman Data Pelamar (filter berdasarkan nilai tes, domisili, status).
- [ ] **Email Integration Sebenarnya**: Mengembalikan fungsi _SMTP Mailer_ pada Login Pelamar dan memberitahukan hasil tes/undangan *interview* secara otomatis via Email.
- [ ] **Manajemen Jadwal Wawancara**: Modul baru untuk HRD mengatur jadwal _interview_ offline/online bagi pelamar yang lulus tes akademik.
- [ ] **Role Management**: Pembuatan sistem Hak Akses/Role Admin agar tidak semua admin dapat mengedit *Master Bank Soal*.

### Fitur Pelamar
- [ ] **Profile Page**: Halaman dimana pelamar bisa mengubah _password_ atau mengatur ulang *biodata dasar* mereka sebelum melamar.
- [ ] **Lupa Akses Tes**: Pengecekan keamanan lebih ketat jika pelamar mencoba memintas (*bypass*) ujian atau *submit* berkali-kali ke server _Redis_.

### System & Architecture
- [ ] **Clean Code & Refactoring**: Memisahkan fungsi-fungsi di *controller* ke *layer Service / Repository* agar file *controller* tidak terlalu besar (terutama pada `ProcessApply`).
- [ ] **Production Ready**: Mematikan _Fiber Hot-Reload_ di environment `production` dan mempersiapkan *CI/CD pipeline*.
