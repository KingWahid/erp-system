# ALISA — ERP Supplier Management System

> **ALISA** (Automated Lifecycle & Integrated Supplier Administration) adalah sistem ERP berbasis web untuk manajemen supplier end-to-end, mencakup onboarding, penilaian performa, manajemen invoice/outstanding, dan kontrol akses berbasis peran (RBAC).

---

## Daftar Isi

1. [Overview Arsitektur](#1-overview-arsitektur)
2. [Tech Stack — Backend](#2-tech-stack--backend)
3. [Tech Stack — Frontend](#3-tech-stack--frontend)
4. [Infrastruktur & DevOps](#4-infrastruktur--devops)
5. [Desain Database (ERD)](#5-desain-database-erd)
6. [Instalasi & Konfigurasi — Backend](#6-instalasi--konfigurasi--backend)
7. [Instalasi & Konfigurasi — Frontend](#7-instalasi--konfigurasi--frontend)
8. [Flow Frontend ke Backend](#8-flow-frontend-ke-backend)
9. [Sistem RBAC](#9-sistem-rbac)
10. [API Reference](#10-api-reference)
11. [Struktur Project](#11-struktur-project)

---

## 1. Overview Arsitektur

```
┌─────────────────────────────────────────────────────────────┐
│                        Browser / Client                      │
│              Next.js 16 (React 19 + TypeScript)             │
└─────────────────────┬───────────────────────────────────────┘
                      │ HTTP (port 80)
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                    Kong API Gateway (port 80)                │
│         Rate Limiting · CORS · Request Size Limit           │
│              X-API-Version header injection                  │
└─────────────────────┬───────────────────────────────────────┘
                      │ HTTP (port 8080)
                      ▼
┌─────────────────────────────────────────────────────────────┐
│               Go Backend — Echo v4 (port 8080)              │
│    JWT Auth Middleware → RBAC Permission Check → Handler    │
│         Clean Architecture: Handler → Usecase → Repository  │
└──────────────┬──────────────────────┬───────────────────────┘
               │                      │
               ▼                      ▼
┌──────────────────────┐   ┌──────────────────────────────────┐
│  PostgreSQL 16       │   │  Redis 7                         │
│  (Data Utama)        │   │  (Cache — Cache-Aside Pattern)   │
└──────────────────────┘   └──────────────────────────────────┘
```

**Pola desain:**
- Backend menggunakan **Clean Architecture** (Domain → Usecase → Repository → Handler)
- Dependency injection dengan **Uber Fx**
- API contract didefinisikan via **OpenAPI 3.0** dan kode di-generate menggunakan `oapi-codegen`
- Frontend menggunakan **Context API** untuk Auth & Permission, **TanStack Query** untuk server state
- Semua request melewati **Kong Gateway** sebelum sampai ke backend

---

## 2. Tech Stack — Backend

| Kategori | Teknologi | Versi | Keterangan |
|---|---|---|---|
| Bahasa | Go | 1.26 | Typed, compiled, concurrent |
| Web Framework | Echo | v4.15 | HTTP router + middleware |
| ORM | GORM | v1.31 | PostgreSQL driver pgx |
| Database | PostgreSQL | 16 | Data utama, relasional |
| Cache | Redis | 7 | Cache-Aside pattern via `go-redis/v9` |
| JWT | golang-jwt/jwt | v5.3 | HS256 signing, claims berisi roles + permissions |
| Config | Viper | v1.21 | Env vars + `.env` file |
| DI Container | Uber Fx | v1.24 | Dependency injection framework |
| Logging | Uber Zap | v1.28 | Structured JSON logging |
| API Codegen | oapi-codegen | — | Generate types + server interface dari OpenAPI |
| Validation | kin-openapi | v0.140 | OpenAPI request validation |
| API Gateway | Kong | 3.6 | DB-less mode, deklaratif via `kong.yml` |
| Container | Docker + Compose | — | Multi-stage build, non-root user |
| Password Hash | bcrypt | — | `golang.org/x/crypto` |

**Library tambahan:**
- `github.com/google/uuid` — UUID v4 generation
- `github.com/joho/godotenv` — Load `.env` lokal
- `github.com/spf13/viper` — Konfigurasi multi-source

---

## 3. Tech Stack — Frontend

| Kategori | Teknologi | Versi | Keterangan |
|---|---|---|---|
| Framework | Next.js | 16 (App Router) | SSR/CSR hybrid, Turbopack dev |
| Language | TypeScript | 5 | Type safety end-to-end |
| UI Library | Ant Design (antd) | 5.23 | Component library enterprise-grade |
| Styling | Tailwind CSS | 4.3 | Utility-first CSS |
| State — Server | TanStack Query | v5.64 | Data fetching, caching, sync |
| State — Client | React Context API | — | Auth + Permission context |
| HTTP Client | Axios | 1.7 | Interceptors untuk token injection & 401 handling |
| Forms | React Hook Form | 7.54 | Form state + validasi |
| Validasi Schema | Zod | 3.24 | Schema validation untuk form |
| Cookie | js-cookie | 3.0 | Simpan JWT token di cookie |
| Icons | @ant-design/icons | 5.6 | Icon set Ant Design |
| Date | Day.js | 1.11 | Format tanggal ringan |

---

## 4. Infrastruktur & DevOps

### Services Docker Compose

| Service | Image | Port | Keterangan |
|---|---|---|---|
| `postgres` | postgres:16-alpine | 5432 | Database utama |
| `redis` | redis:7-alpine | 6379 | Cache layer |
| `app` | Build lokal | 8080 | Go backend (non-root user) |
| `kong` | kong:3.6-ubuntu | 80, 8001 | API Gateway (DB-less) |

### Kong Gateway Configuration

Kong dikonfigurasi **DB-less** (deklaratif via `backend/kong/kong.yml`) dengan plugin:

| Plugin | Config |
|---|---|
| `rate-limiting` | 100 req/menit, 1000 req/jam |
| `cors` | Semua origin, method GET/POST/PUT/DELETE/OPTIONS |
| `request-size-limiting` | Max payload 10 MB |
| `response-transformer` | Tambah header `X-API-Version: v1` |

### Dockerfile — Multi-stage Build

```
Stage 1 (builder): golang:1.26-alpine
  └─ go mod download → go build → binary /app/bin/server

Stage 2 (runtime): alpine:3.19
  └─ Timezone Asia/Jakarta
  └─ Non-root user (appuser)
  └─ HEALTHCHECK via wget /health
  └─ EXPOSE 8080
```

---

## 5. Desain Database (ERD)

### Diagram Relasi

```
┌─────────────┐       ┌─────────────┐       ┌─────────────────┐
│    users    │───────│  user_roles │───────│     roles       │
├─────────────┤  M:M  ├─────────────┤  M:M  ├─────────────────┤
│ id (PK)     │       │ user_id(FK) │       │ id (PK)         │
│ name        │       │ role_id(FK) │       │ name (UNIQUE)   │
│ email       │       └─────────────┘       │ description     │
│ password    │                             │ is_system       │
│ is_active   │                             │ is_active       │
│ last_login  │                             └────────┬────────┘
└─────────────┘                                      │ M:M
                                             ┌───────┴────────┐
                                             │ role_permissions│
                                             ├────────────────┤
                                             │ role_id (FK)   │
                                             │ permission_id  │
                                             └───────┬────────┘
                                                     │
                                          ┌──────────┴──────────┐
                                          │     permissions      │
                                          ├─────────────────────┤
                                          │ id (PK)             │
                                          │ resource            │
                                          │ action              │
                                          │ endpoint_path       │
                                          │ endpoint_method     │
                                          │ description         │
                                          │ hide                │
                                          └─────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                          suppliers                               │
├──────────────────────────────────────────────────────────────────┤
│ id (PK) · code (UNIQUE) · supplier_no (UNIQUE) · name · alias   │
│ address · city · province · country · postal_code · phone       │
│ email · website · logo_url · notes                              │
│ status: draft|active|in_progress|blocked|inactive               │
│ stage: draft|in_review|in_assessment|active                     │
│ sla_hours · is_blocked · block_reason                           │
│ created_by · updated_by · created_at · updated_at · deleted_at  │
└────────┬─────────────────────────────────────────────────────────┘
         │ 1:N untuk semua tabel berikut
         ├──────────────────────────────────────────────────────────
         │
         ├──► supplier_addresses
         │    id · supplier_id(FK) · name · address · city
         │    province · country · postal_code · is_main
         │
         ├──► supplier_contacts
         │    id · supplier_id(FK) · name · position
         │    phone · mobile · email · is_primary
         │
         ├──► supplier_groups
         │    id · supplier_id(FK) · group_name · value · is_active
         │    (contoh: group_name="Industry", value="Manufacture")
         │
         ├──► supplier_materials
         │    id · supplier_id(FK) · material_group · material_id · is_active
         │
         ├──► supplier_performance_ratings
         │    id · supplier_id(FK) · price_rating(1-5)
         │    delivery_rating(1-5) · notes · reviewed_by · reviewed_at
         │
         ├──► supplier_stage_histories
         │    id · supplier_id(FK) · from_stage · to_stage
         │    notes · changed_by · elapsed_ms
         │
         └──► supplier_invoices
              id · supplier_id(FK) · invoice_number · project_name
              amount · currency · invoice_date · due_date · paid_date
              status: unpaid|partial|paid|overdue
              paid_amount · notes · created_by
```

### Supplier Workflow Stages

```
draft ──► in_review ──► in_assessment ──► active
  │                                          │
  └── (blocked kapan saja) ─────────────────►┘
```

### Roles & Permissions Default

| Role | Permissions |
|---|---|
| `admin` | Semua permissions |
| `manager` | supplier:*, material:*, rating:*, workflow:*, review:* (kecuali delete) |
| `viewer` | Semua action `read` saja |
| `supplier` | — (portal supplier, dikonfigurasi terpisah) |

**Permission key format:** `resource:action`
Contoh: `supplier:read`, `supplier:create`, `workflow:advance`, `review:approve`

---

## 6. Instalasi & Konfigurasi — Backend

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- `make` (opsional, tersedia di Git Bash / WSL di Windows)

### Cara 1 — Docker Compose (Recommended)

```bash
cd backend

# 1. Copy env file
cp .env.example .env  # atau buat .env baru (lihat contoh di bawah)

# 2. Jalankan semua services (postgres, redis, app, kong)
docker compose up --build -d

# 3. Cek status
docker compose ps
docker compose logs app -f
```

Setelah services berjalan:
- Backend API: `http://localhost:8080`
- Kong Gateway: `http://localhost:80`
- Kong Admin: `http://localhost:8001`

### Cara 2 — Lokal (tanpa Docker untuk app)

```bash
cd backend

# 1. Jalankan hanya infrastructure (postgres + redis)
docker compose up postgres redis -d

# 2. Copy dan edit .env
cp .env.example .env

# 3. Install dependencies
go mod download

# 4. Jalankan migrasi database
make migrate

# 5. (Opsional) Insert mock data
make seed-mock

# 6. Jalankan server
make run
# atau: go run ./cmd/main.go
```

### Environment Variables — Backend

Buat file `backend/.env`:

```env
# Application
APP_NAME=erp-system
APP_PORT=8080
APP_ENV=development        # development | production

# Database PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=erp_supplier
DB_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=ganti-dengan-secret-minimal-32-karakter-yang-aman
JWT_EXPIRATION_HOURS=24
```

> **Penting:** `JWT_SECRET` harus minimal 32 karakter dan tidak boleh dicommit ke repository.

### Make Commands

```bash
make generate    # Regenerate kode dari OpenAPI spec
make build       # Build binary ke bin/server
make run         # Jalankan server lokal
make migrate     # Jalankan migrasi database
make seed-mock   # Insert data mock untuk development
make drop        # DROP semua tabel (HANYA dev!)
make fresh       # drop + migrate + seed-mock
make up          # docker compose up --build -d
make down        # docker compose down
make docker-fresh # docker compose down -v + up --build
make tidy        # go mod tidy
make lint        # golangci-lint run
```

### Default Admin User

Setelah migrasi berhasil, user admin sudah tersedia:

```
Email    : admin@erp.local
Password : Admin@123
Role     : admin (full access)
```

---

## 7. Instalasi & Konfigurasi — Frontend

### Prerequisites

- Node.js 18+ (disarankan 20 LTS)
- npm atau yarn

### Setup

```bash
cd frontend

# 1. Install dependencies
npm install

# 2. Copy dan edit environment variables
cp .env.local.example .env.local
```

### Environment Variables — Frontend

Buat file `frontend/.env.local`:

```env
# URL backend (via Kong Gateway atau langsung)
NEXT_PUBLIC_API_URL=http://localhost:8000

# Nama cookie untuk JWT token
NEXT_PUBLIC_TOKEN_COOKIE=erp_token
```

> Gunakan `http://localhost:8000` (Kong) untuk lingkungan yang menyerupai production.
> Gunakan `http://localhost:8080` untuk bypass Kong dan akses langsung ke backend.

### Jalankan Development Server

```bash
npm run dev        # Development dengan Turbopack
npm run build      # Production build
npm run start      # Jalankan production build
npm run lint       # ESLint
npm run type-check # TypeScript type check
```

Frontend berjalan di: `http://localhost:3000`

---

## 8. Flow Frontend ke Backend

### 8.1 Authentication Flow

```
User input email+password
        │
        ▼
LoginPage (POST /auth/login)
        │
        ▼
AuthContext.login()
  └─► authService.login(payload)
        │
        ▼
axios apiClient
  └─► request interceptor: tidak ada token (login publik)
        │
        ▼
Kong Gateway (port 80)
  └─► rate-limit check → CORS headers
        │
        ▼
Go Backend: POST /auth/login
  └─► AuthHandler.Login()
  └─► AuthUsecase.Login()
        ├─► UserRepository.FindByEmail()
        ├─► bcrypt.CompareHashAndPassword()
        └─► JWTManager.Generate(userID, email, roles, permissions)
              └─► Claims: { user_id, email, roles[], permissions[] }
        │
        ▼
Response: { success: true, data: { token, user } }
        │
        ▼
AuthContext
  └─► Cookies.set("erp_token", token)  ← simpan di cookie
  └─► setUser(user)
  └─► setToken(token)
        │
        ▼
PermissionContext
  └─► Baca user.permissions dari AuthContext
  └─► Bangun Set<permission> untuk pengecekan cepat
        │
        ▼
Redirect ke /dashboard
```

### 8.2 Authenticated Request Flow

```
User aksi di halaman (misal: klik "Create Supplier")
        │
        ▼
PermissionContext.can("supplier:create")
  └─► Cek apakah permission ada di Set → jika tidak: tampilkan disabled/hidden
        │ (jika boleh)
        ▼
supplierService.createSupplier(payload)
  └─► apiClient.post("/suppliers", payload)
        │
        ▼
axios request interceptor
  └─► Baca Cookies.get("erp_token")
  └─► Set header: Authorization: Bearer <token>
        │
        ▼
Kong Gateway
  └─► Rate limiting, CORS
        │
        ▼
Go Backend: POST /suppliers
  └─► Router: operationMiddleware["CreateSupplier"] = [auth]
  └─► AuthMiddleware.Authenticate()
        ├─► Parse "Authorization: Bearer <token>"
        ├─► JWTManager.Validate(token)
        └─► Inject ke context: user_id, email, roles, permissions
  └─► SupplierHandler.CreateSupplier()
  └─► SupplierUsecase.Create()
  └─► SupplierRepository.Create()
        │
        ▼
Response: { success: true, data: { ...supplier } }
        │
        ▼
TanStack Query: invalidate cache "suppliers"
  └─► Refetch list otomatis
```

### 8.3 Session Restore on Page Reload

```
Browser reload
        │
        ▼
AuthProvider useEffect (mount)
  └─► Cookies.get("erp_token")
        ├─► Tidak ada → setIsLoading(false) → redirect /login
        └─► Ada token →
              authService.getProfile() → GET /auth/me
                ├─► Token valid → setUser(profile) → render app
                └─► Token expired → Cookies.remove() → redirect /login
```

### 8.4 Token Expired / 401 Handling

```
Response 401 dari backend
        │
        ▼
axios response interceptor (lib/axios.ts)
  └─► Cookies.remove("erp_token")
  └─► window.location.href = "/login"
```

---

## 9. Sistem RBAC

### Arsitektur RBAC

Sistem menggunakan **Role-Based Access Control (RBAC)** dengan model many-to-many:

```
User ──► [user_roles] ──► Role ──► [role_permissions] ──► Permission
```

Permission key format: **`resource:action`**

### Bagaimana Permission Dikemas ke JWT

Saat login berhasil:
1. Backend load user beserta semua roles dan permissions via GORM preload
2. `user.GetPermissions()` mengumpulkan semua permission dari semua roles (deduplication)
3. Permissions dimasukkan ke dalam JWT claims sebagai `[]string`
4. Frontend menerima token yang sudah berisi permissions — **tidak perlu request tambahan**

```go
// JWT Claims
type Claims struct {
    UserID string   `json:"user_id"`
    Email  string   `json:"email"`
    Roles  []string `json:"roles"`
    Perms  []string `json:"permissions"`  // e.g. ["supplier:read","supplier:create"]
    jwt.RegisteredClaims
}
```

### Permission Check — Backend

Middleware `RequirePermission` dan `RequireRole` tersedia untuk granular access control:

```go
// Contoh penggunaan di router
opMW := map[string][]echo.MiddlewareFunc{
    "DeleteSupplier": {
        auth,
        authMiddleware.RequirePermission("supplier:delete"),
    },
    "AdvanceSupplierStage": {
        auth,
        authMiddleware.RequirePermission("workflow:advance"),
    },
}
```

Middleware membaca permissions dari Echo context (yang sudah diisi oleh `Authenticate()`).

### Permission Check — Frontend

`PermissionContext` menyediakan helper functions:

```tsx
const { can, canAny, hasRole, isAdmin } = usePermission();

// Semua permission harus dipenuhi
can("supplier:create")                     // true/false

// Salah satu cukup
canAny("supplier:update", "supplier:delete")

// Cek role
hasRole("manager")
isAdmin  // shorthand untuk hasRole("admin")
```

**Admin bypass:** User dengan role `admin` otomatis lolos semua pengecekan permission di frontend.

### Daftar Permissions

| Resource | Action | Method | Endpoint | Deskripsi |
|---|---|---|---|---|
| supplier | read | GET | /suppliers | List dan detail supplier |
| supplier | create | POST | /suppliers | Buat supplier baru |
| supplier | update | PUT | /suppliers/:id | Update data supplier |
| supplier | delete | DELETE | /suppliers/:id | Soft delete supplier |
| supplier | block | POST | /suppliers/:id/block | Blokir/unblokir supplier |
| supplier | export | GET | /suppliers/export | Export data supplier |
| material | read | GET | /suppliers/:id/materials | Lihat material supplier |
| material | update | PUT | /suppliers/:id/materials | Update material list |
| rating | read | GET | /suppliers/:id/ratings | Lihat rating performa |
| rating | create | POST | /suppliers/:id/ratings | Tambah rating |
| workflow | advance | POST | /suppliers/:id/next-stage | Advance stage supplier |
| review | approve | POST | /reviews/:id/approve | Approve review (hidden) |
| review | reject | POST | /reviews/:id/reject | Reject review (hidden) |

---

## 10. API Reference

Base URL (via Kong): `http://localhost:80`
Base URL (direct): `http://localhost:8080`

### Authentication

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| POST | `/auth/register` | ❌ | Register user baru |
| POST | `/auth/login` | ❌ | Login, dapat JWT token |
| GET | `/auth/me` | ✅ | Get profile user saat ini |

### Supplier

| Method | Endpoint | Auth | Permission | Deskripsi |
|---|---|---|---|---|
| GET | `/suppliers/stats` | ✅ | — | Stats dashboard supplier |
| GET | `/suppliers` | ✅ | supplier:read | List supplier (pagination, search, filter) |
| POST | `/suppliers` | ✅ | supplier:create | Buat supplier baru |
| GET | `/suppliers/:id` | ✅ | supplier:read | Detail supplier |
| PUT | `/suppliers/:id` | ✅ | supplier:update | Update supplier |
| DELETE | `/suppliers/:id` | ✅ | supplier:delete | Soft delete supplier |
| POST | `/suppliers/:id/block` | ✅ | supplier:block | Blokir/unblokir supplier |
| POST | `/suppliers/:id/next-stage` | ✅ | workflow:advance | Advance ke stage berikutnya |

### Sub-resources Supplier

| Method | Endpoint | Deskripsi |
|---|---|---|
| GET/POST | `/suppliers/:id/addresses` | List / Tambah alamat |
| PUT/DELETE | `/suppliers/:id/addresses/:addrId` | Update / Hapus alamat |
| GET/POST | `/suppliers/:id/contacts` | List / Tambah kontak |
| PUT/DELETE | `/suppliers/:id/contacts/:contactId` | Update / Hapus kontak |
| GET/POST | `/suppliers/:id/groups` | List / Tambah grup |
| DELETE | `/suppliers/:id/groups/:groupId` | Hapus grup |
| GET/PUT | `/suppliers/:id/materials` | List / Update materials |
| GET/POST | `/suppliers/:id/ratings` | List / Tambah rating performa |
| GET | `/suppliers/:id/stage-history` | Riwayat perubahan stage |
| GET | `/suppliers/:id/outstandings` | List invoice/outstanding |

### Health Check

```
GET /health
→ { "status": "ok", "service": "erp-supplier-management" }
```

### Format Response

```json
// Success
{
  "success": true,
  "data": { ... },
  "meta": { "page": 1, "limit": 10, "total": 100, "total_pages": 10 }
}

// Error
{
  "success": false,
  "error": {
    "code": "ERR_SUPPLIER_NOT_FOUND",
    "message": "supplier not found",
    "message_id": "error_supplier_not_found"
  }
}
```

---

## 11. Struktur Project

```
erp-system/
├── backend/
│   ├── api/                          # OpenAPI spec + codegen config
│   │   ├── openapi.yaml              # API contract (sumber kebenaran)
│   │   ├── codegen.types.yaml        # Config generate types
│   │   └── codegen.server.yaml       # Config generate server interface
│   │
│   ├── cmd/
│   │   ├── main.go                   # Entrypoint
│   │   ├── fx/
│   │   │   ├── app.go                # Fx app setup
│   │   │   └── modules/              # Fx modules (DI)
│   │   │       ├── auth.go
│   │   │       ├── cache.go
│   │   │       ├── config.go
│   │   │       ├── database.go
│   │   │       ├── jwt.go
│   │   │       ├── logger.go
│   │   │       ├── server.go
│   │   │       └── supplier.go
│   │   ├── migrate/main.go           # CLI migration runner
│   │   ├── genhash/main.go           # Utility generate bcrypt hash
│   │   └── server/router.go          # Route registration
│   │
│   ├── internal/
│   │   ├── auth/
│   │   │   ├── domain/               # Entity (User, Role, Permission)
│   │   │   ├── handler/              # HTTP handler + adapter
│   │   │   ├── repository/           # User + Role repository (GORM)
│   │   │   └── usecase/              # Business logic auth
│   │   │
│   │   ├── supplier/
│   │   │   ├── domain/               # Entity supplier + sub-entities
│   │   │   ├── handler/              # HTTP handler + adapter
│   │   │   ├── repository/           # Supplier repo + cached repo
│   │   │   └── usecase/              # Business logic supplier + DTO
│   │   │
│   │   ├── workflow/                 # (In Progress)
│   │   │   ├── domain/
│   │   │   ├── repository/
│   │   │   └── usecase/
│   │   │
│   │   └── generated/                # Auto-generated dari OpenAPI
│   │       ├── types.gen.go
│   │       └── server.gen.go
│   │
│   ├── pkg/
│   │   ├── cache/                    # Cache-Aside abstraction + Redis client
│   │   ├── config/                   # Viper config loader
│   │   ├── database/                 # GORM PostgreSQL setup
│   │   ├── errors/                   # Custom error codes + helpers
│   │   ├── jwt/                      # JWT Manager (generate + validate)
│   │   ├── logger/                   # Zap logger setup
│   │   └── middleware/               # Auth middleware (Authenticate, RequirePermission, RequireRole)
│   │
│   ├── migrations/
│   │   ├── 001_init.sql              # Schema RBAC + Supplier + seed data
│   │   ├── 002_supplier_invoices.sql # Schema supplier invoices
│   │   └── seed/mock_data.sql        # Data mock untuk development
│   │
│   ├── kong/kong.yml                 # Kong declarative config
│   ├── Dockerfile                    # Multi-stage build
│   ├── docker-compose.yaml           # Services orchestration
│   ├── Makefile                      # Developer commands
│   ├── go.mod / go.sum
│   └── .env                          # Environment variables (jangan dicommit)
│
└── frontend/
    ├── src/
    │   ├── app/                      # Next.js App Router
    │   │   ├── (auth)/
    │   │   │   └── login/page.tsx    # Halaman login
    │   │   └── (dashboard)/
    │   │       ├── layout.tsx        # Layout sidebar + auth guard
    │   │       ├── dashboard/        # Halaman dashboard utama
    │   │       └── suppliers/        # Halaman manajemen supplier
    │   │           ├── page.tsx      # List supplier
    │   │           ├── [id]/         # Detail supplier
    │   │           └── configurations/ # Konfigurasi
    │   │
    │   ├── contexts/
    │   │   ├── AuthContext.tsx        # Auth state, login, logout
    │   │   └── PermissionContext.tsx  # can(), canAny(), hasRole()
    │   │
    │   ├── hooks/
    │   │   └── useAuth.ts            # Hook akses AuthContext
    │   │
    │   ├── lib/
    │   │   └── axios.ts              # Axios instance + interceptors
    │   │
    │   ├── services/
    │   │   ├── auth.service.ts       # Auth API calls
    │   │   └── supplier.service.ts   # Supplier API calls
    │   │
    │   └── types/
    │       └── api.ts                # TypeScript types (mirror dari OpenAPI)
    │
    ├── .env.local                    # Env variables (jangan dicommit)
    ├── .env.local.example
    ├── package.json
    ├── next.config.ts
    └── tsconfig.json
```

---

## Quick Start (All-in-One)

```bash
# Clone repository
git clone <repo-url>
cd erp-system

# 1. Setup backend .env
cd backend
cp .env.example .env   # edit JWT_SECRET minimal

# 2. Jalankan semua services
docker compose up --build -d

# 3. Setup frontend
cd ../frontend
cp .env.local.example .env.local
npm install
npm run dev
```

**Akses:**
- 🌐 Frontend: http://localhost:3000

**Login default:**
- Email: `admin@erp.local`
- Password: `Admin@123`

---

