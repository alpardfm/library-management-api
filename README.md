# 📚 **Library Management API**

## **Deskripsi Proyek**

**Library Management API** adalah sistem manajemen perpustakaan digital yang dibangun dengan **Go (Golang)** menggunakan arsitektur **Clean Architecture**. API ini menyediakan fungsi lengkap untuk mengelola perpustakaan modern termasuk manajemen buku, anggota, peminjaman, dan pengembalian.

## **🎯 Fitur Utama**

### **1. 🔐 Sistem Autentikasi & Otorisasi**

- **Registrasi pengguna** dengan validasi email
- **Login dengan JWT** (JSON Web Tokens)
- **Role-based access control** (RBAC):
    - **Admin**: Akses penuh ke semua fitur
    - **Librarian (Pustakawan)**: Kelola buku & peminjaman
    - **Member (Anggota)**: Pinjam & kembalikan buku
- **Token expiration** dengan konfigurasi waktu

### **2. 📖 Manajemen Buku**

- **CRUD lengkap** untuk data buku (Create, Read, Update, Delete)
- **Validasi ISBN** (International Standard Book Number)
- **Pencarian buku** dengan filter:
    - Judul buku
    - Penulis
    - ISBN
    - Genre
- **Status ketersediaan** buku (tersedia/tidak)
- **Manajemen stok** (total copy vs copy tersedia)

### **3. 👥 Manajemen Anggota**

- **Registrasi anggota** baru
- **Verifikasi status aktif/non-aktif**
- **Limit peminjaman** (maksimal buku per anggota)
- **Riwayat peminjaman** per anggota

### **4. 📅 Sistem Peminjaman**

- **Pinjam buku** dengan validasi:
    - Buku tersedia
    - Anggota aktif
    - Tidak melebihi limit
- **Atur tanggal jatuh tempo** (default: 14 hari)
- **Tracking status**:
    - Dipinjam (borrowed)
    - Dikembalikan (returned)
    - Terlambat (overdue)
- **Perhitungan denda** otomatis (Rp 1000/hari)

### **5. 📊 Laporan & Monitoring**

- **Daftar peminjaman aktif**
- **Buku yang terlambat** dikembalikan
- **Statistik penggunaan** perpustakaan
- **Pencatatan riwayat** lengkap

## **🏗️ Arsitektur Teknis**

### **Struktur Folder**
```
library-management-api/
├── 📁 cmd/api/           # Entry point aplikasi
├── 📁 internal/          # Kode internal aplikasi
│   ├── models/          # Struct database (User, Book, Borrow)
│   ├── repository/      # Layer akses database (GORM)
│   ├── service/         # Business logic
│   ├── handler/         # HTTP controllers (Gin)
│   ├── middleware/      # Auth, logging, recovery
│   └── dto/             # Data Transfer Objects
├── 📁 pkg/              # Package reusable
│   ├── database/        # Koneksi PostgreSQL
│   ├── auth/            # JWT authentication
│   └── utils/           # Helper functions
├── 📁 tests/            # Test suites lengkap
│   ├── unit/            # Unit tests
│   ├── integration/     # Integration tests
│   └── e2e/             # End-to-end tests
└── 📁 configs/          # Konfigurasi aplikasi

```

### **Teknologi Stack**

- **Bahasa**: Go 1.21+
- **Framework**: Gin Gonic (HTTP router)
- **Database**: PostgreSQL 15+
- **ORM**: GORM (Go ORM)
- **Authentication**: JWT (JSON Web Tokens)
- **Testing**: Testify, SQLMock
- **Container**: Docker & Docker Compose
- **Code Quality**: SonarQube, golangci-lint
- **CI/CD**: GitHub Actions

## **🔐 Keamanan**

### **Security Features**

- **Password hashing** dengan bcrypt
- **JWT token** dengan expiration
- **Role-based authorization**
- **Input validation** lengkap
- **SQL injection prevention** (GORM)
- **CORS configuration**
- **Rate limiting** (bisa diimplementasi)

### **Validasi Data**

- Validasi email format
- Validasi ISBN format
- Validasi tanggal peminjaman
- Validasi stok buku
- Custom validation rules

## **📡 API Endpoints**

### **Public Routes**

```
POST   /api/v1/auth/register    Registrasi anggota baru
POST   /api/v1/auth/login       Login dan dapatkan token
GET    /health                  Liveness check API
GET    /ready                   Readiness check API + DB ping

```

### **Protected Routes (Perlu Auth)**

```
# Books
GET    /api/v1/books            List semua buku (dengan pagination)
GET    /api/v1/books/:id        Detail buku spesifik
POST   /api/v1/books            Tambah buku baru (Admin/Librarian)
PUT    /api/v1/books/:id        Update buku (Admin/Librarian)
DELETE /api/v1/books/:id        Hapus buku (Admin/Librarian)

# Borrow
POST   /api/v1/borrow           Pinjam buku
POST   /api/v1/borrow/return    Kembalikan buku
GET    /api/v1/borrow/my-books  Riwayat peminjaman saya
GET    /api/v1/borrow/active    List peminjaman aktif (Admin/Librarian)
GET    /api/v1/borrow/overdue   List buku terlambat (Admin/Librarian)

```

### **Request/Response Examples**

### **Register User**

```
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "role": "member"
}

Response:
{
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "role": "member"
  }
}

```

### **Borrow Book**

```
POST /api/v1/borrow
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "book_id": 5
}

Response:
{
  "message": "Book borrowed successfully",
  "data": {
    "id": 100,
    "user_id": 1,
    "book_id": 5,
    "borrow_date": "2024-01-15T10:30:00Z",
    "due_date": "2024-01-29T10:30:00Z",
    "status": "borrowed"
  }
}

```

## **🗄️ Database Schema**

### **Tables Structure**

### **1. Users Table**

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'member',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

```

### **2. Books Table**

```sql
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    isbn VARCHAR(13) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    publisher VARCHAR(100),
    publication_year INTEGER,
    genre VARCHAR(50),
    total_copies INTEGER DEFAULT 1,
    available_copies INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

```

### **3. Borrow Records Table**

```sql
CREATE TABLE borrow_records (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
    borrow_date DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date DATE NOT NULL,
    return_date DATE,
    status VARCHAR(20) DEFAULT 'borrowed',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

```

## **🧪 Testing Strategy**

### **Test Coverage 100% Target**

```
✅ Unit Tests:    80% - Individual components
✅ Integration:   15% - API endpoints & DB integration
✅ E2E Tests:     5%  - Complete user flows
✅ Total:        100% - Full coverage

```

### **Testing Layers**

1. **Models**: Validasi data & business rules
2. **Repository**: Database operations dengan SQL mock
3. **Service**: Business logic dengan mock dependencies
4. **Handler**: HTTP requests/responses
5. **Middleware**: Auth, logging, error handling
6. **Integration**: API endpoints dengan database real

### **Test Tools**

- **Testify**: Assertions & mocking
- **SQLMock**: Mock database untuk unit tests
- **httptest**: HTTP testing
- **Docker**: Test database isolation

### **Testing Notes**

- Unit test `service` dan `repository` tetap memakai mock/`sqlmock`.
- Concurrency hardening untuk borrow flow dibuktikan lewat integration test Postgres nyata di [`tests/integration/borrow_concurrency_test.go`](github.com/alpardfm/library-management-api/tests/integration/borrow_concurrency_test.go).
- Constraint/index invariant yang dibuat di `AutoMigrate` memakai SQL Postgres-specific. Saat dialector bukan `postgres`, raw SQL invariant tersebut di-skip agar test environment non-Postgres tidak gagal palsu.

## **🚀 Deployment**

### **Local Development**

```bash
# 1. Clone repository
git clone <https://github.com/username/library-management-api.git>

# 2. Setup environment
cp .env.example .env
# Edit .env sesuai konfigurasi

# 3. Start services
docker-compose up -d

# 4. Run application
go run cmd/api/main.go

```

### **Docker Deployment**

```bash
# Build image
docker build -t library-api:latest .

# Run with Docker Compose
docker-compose -f docker-compose.prod.yml up -d

```

### **Kubernetes Deployment**

```bash
# Apply Kubernetes manifests
kubectl apply -f deployments/kubernetes/

```

## **📊 Monitoring & Logging**

### **Health Checks**

```
GET /health

Response:
{
  "status": "healthy",
  "app": "Library Management API",
  "version": "1.0.0",
  "env": "production"
}

GET /ready

Response:
{
  "status": "ready"
}

```

### **Logging Features**

- **Structured logging** dengan Zerolog
- **Request logging** standar (`request_id`, `method`, `path`, `status`, `latency_ms`, `user_id`, `role`)
- **Error logging** untuk panic/recovery dan 5xx request
- **Sensitive data tidak ikut dilog** (mis. token/header auth/password)

### **Metrics** (Opsional)

- API request count
- Database query performance
- Error rates
- Response time percentiles

## **🔧 Maintenance**

### **Database Migrations**

```bash
# Auto migrate on startup (development)
# Manual migration for production
go run cmd/migrate/main.go

```

### **Backup & Recovery**

- **Automatic backups** dengan pg_dump
- **Point-in-time recovery**
- **Data export** untuk reporting

### **Scaling Strategies**

- **Horizontal scaling** dengan load balancer
- **Database connection pooling**
- **Caching layer** (Redis opsional)
- **Message queue** untuk async tasks

## **🎯 Target Pengguna**

### **1. Perpustakaan Umum**

- Sekolah & Universitas
- Perpustakaan Kota/Kabupaten
- Perpustakaan Khusus (Rumah Sakit, Perusahaan)

### **2. Aplikasi Edukasi**

- Platform e-learning
- Sistem manajemen konten edukasi
- Aplikasi membaca digital

### **3. Bisnis**

- Manajemen inventaris buku
- Sistem rental buku
- Koleksi pribadi/organisasi

## **✨ Keunggulan**

### **1. Kode Berkualitas Tinggi**

- **100% test coverage** dengan SonarQube integration
- **Clean architecture** dengan separation of concerns
- **Zero code smells** & **zero duplication**
- **Comprehensive error handling**

### **2. Developer Experience**

- **Dokumentasi lengkap** dengan examples
- **Easy setup** dengan Docker
- **Comprehensive testing suite**
- **IDE friendly** dengan proper Go modules

### **3. Production Ready**

- **Graceful shutdown** handling
- **Health checks** & monitoring
- **Security best practices**
- **Scalable architecture**

### **4. Extensible**

- **Modular design** mudah ditambah fitur
- **Plugin architecture** untuk additional features
- **API versioning** support
- **Multi-database** support (MySQL, SQLite opsional)

## **📈 Roadmap**

### **Versi 1.0** (Current)

- ✅ Core features: Books, Users, Borrowing
- ✅ Authentication & Authorization
- ✅ Basic reporting

### **Versi 2.0** (Planned)

- 🔄 Notification system (email/SMS)
- 🔄 Reservation system (antrian buku)
- 🔄 Fine payment integration
- 🔄 Advanced analytics dashboard

### **Versi 3.0** (Future)

- 🔄 Mobile app support
- 🔄 QR Code book scanning
- 🔄 Machine learning recommendations
- 🔄 Multi-tenant support

## **👥 Kontribusi**

### **Cara Berkontribusi**

1. Fork repository
2. Create feature branch
3. Commit changes dengan meaningful messages
4. Push to branch
5. Create Pull Request

### **Coding Standards**

- **Go coding standards** dengan gofmt
- **Test coverage** minimal 90%
- **Documentation** untuk public APIs
- **Commit convention** (Conventional Commits)

## **📄 License**

**MIT License** - Bebas digunakan untuk keperluan komersial maupun non-komersial dengan attribution.

---

## **🎯 Kesimpulan**

**Library Management API** adalah solusi lengkap untuk manajemen perpustakaan digital modern dengan:

✅ **Fitur lengkap** untuk operasional perpustakaan

✅ **Arsitektur bersih** mudah maintain & scale

✅ **Keamanan terjamin** dengan JWT & RBAC

✅ **Test coverage 100%** kualitas kode terjamin

✅ **Production ready** dengan Docker & monitoring

✅ **Extensible** untuk kebutuhan masa depan

**Siap digunakan untuk:**

- 🏫 Perpustakaan sekolah/universitas
- 🏢 Perpustakaan perusahaan
- 📱 Aplikasi baca buku digital
- 🎓 Platform e-learning
- 📊 Sistem manajemen inventaris buku

**"From zero to production-ready library system in minutes!"** 🚀
