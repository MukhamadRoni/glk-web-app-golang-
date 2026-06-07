# Project Plan: Web Admin & Web Pelamar (Golang MVC)

## 1. System Overview & Hardware Constraints

We are building a web application consisting of two parts: a **Web Admin** and a **Web Pelamar** (Job Applicant Web).

- **Database:** PostgreSQL
- **Language/Framework:** Golang (using Fiber Framework for speed and standard html/template engine for views).
- **Architecture Pattern:** MVC (Model-View-Controller) styled similarly to CodeIgniter 3 for clear file separation.
- **Infrastructure Constraint:** Deploying on a VPS with 8 GB RAM, 4 Cores, but strictly **limited to 30 GB Storage**. Docker images MUST be optimized using multi-stage builds to keep image size under 50MB.

---

## 2. Directory Structure to Generate

Please follow this exact structural pattern when generating code:

proyek-pelamar/
├── config/ # Database connection & ENV loaders
│ └── database.go
├── models/ # GORM Structs & DB Queries
│ ├── pelamar.go
│ └── admin.go
├── controllers/ # Route Handlers & Business Logic
│ ├── auth.go
│ ├── admin.go
│ └── pelamar.go
├── views/ # HTML Templates (Go html/template)
│ ├── layouts/
│ ├── admin/
│ └── pelamar/
├── static/ # Static Assets (CSS, JS)
├── Dockerfile # Multi-stage optimized build
├── docker-compose.yml # Go App + PostgreSQL
├── go.mod
└── main.go # Application entrypoint & routing

---

## 3. Tech Stack Requirements

- **Backend:** Go (Golang) 1.22+
- **Web Framework:** Gofiber (://github.com) or Gin Gonic.
- **ORM:** GORM (gorm.io/gorm) with PostgreSQL driver.
- **Template Engine:** Go standard `html/template` or Fiber Views.
- **Database:** PostgreSQL 16-alpine (Strictly alpine version to save disk space).

---

## 4. Phase-by-Phase Execution Plan

### Phase 1: Infrastructure & Base Setup

1. Create `go.mod` and install required dependencies (`gofiber`, `gorm`, `postgres driver`, `godotenv`).
2. Write an optimized **multi-stage `Dockerfile`** that compiles the Go binary and discards development caches, targeting a final size of < 50MB.
3. Write `docker-compose.yml` defining the Go application and a PostgreSQL service. Ensure database volumes are mapped correctly.
4. Write `config/database.go` to handle safe connections and automatic database migration (AutoMigrate) via GORM.

### Phase 2: Database Schema & Models

1. Create `models/pelamar.go` with columns: ID, Name, Email, Password (hashed), CV_URL, Status, CreatedAt.
2. Create `models/admin.go` with columns: ID, Username, Password (hashed), CreatedAt.
3. Include standard CRUD functions inside the model files using GORM methods.

### Phase 3: Authentication & Core Controllers

1. Implement JWT or Session-based authentication inside `controllers/auth.go`.
2. Implement routing logic in `main.go` to separate endpoints:
   - **Public/Applicant Routes:** `/login`, `/register`, `/apply`, `/dashboard`
   - **Admin Routes (Protected):** `/admin/login`, `/admin/dashboard`, `/admin/applicants`

### Phase 4: Frontend Views (HTML Template)

1. Set up a base layout in `views/layouts/base.html` using a CDN for Tailwind CSS (no local npm node_modules allowed to save storage space).
2. Create simple forms for Login, Register, Profile, and Admin Dashboard data tables.

---

## 5. Immediate First Instruction for the AI

"Based on the plan.md above, let's start with **Phase 1**. Please generate the `go.mod` configuration, the optimized `Dockerfile`, the `docker-compose.yml`, and the `config/database.go` file to establish our environment."
