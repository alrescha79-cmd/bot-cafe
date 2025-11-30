# 🏗️ Architecture Documentation

## Overview

Bot Telegram Café/Resto dibangun menggunakan **Microservices Architecture** dengan 5 service independen dan 1 agent sebagai interface ke Telegram.

## System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                      Telegram API                            │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                  Telegram Bot Agent                          │
│                  (Main Orchestrator)                         │
│  - Receives messages from users                              │
│  - Routes requests to appropriate microservices              │
│  - Formats responses for users                               │
│  - Manages user states & dialogs                             │
└──┬────────┬────────┬────────┬────────┬─────────────────────┘
   │        │        │        │        │
   │        │        │        │        │
   ▼        ▼        ▼        ▼        ▼
┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐
│Auth  │ │Menu  │ │Promo │ │Info  │ │Media │
│Service│ │Service│ │Service│ │Service│ │Service│
│:8081 │ │:8082 │ │:8083 │ │:8084 │ │:8085 │
└──┬───┘ └──┬───┘ └──┬───┘ └──┬───┘ └──┬───┘
   │        │        │        │        │
   ▼        ▼        ▼        ▼        ▼
┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐
│auth  │ │menu  │ │promo │ │info  │ │media │
│.db   │ │.db   │ │.db   │ │.db   │ │.db   │
└──────┘ └──────┘ └──────┘ └──────┘ └──────┘
```

## Microservices Details

### 1. Auth Service (Port 8081)
**Responsibility:** Manajemen autentikasi dan otorisasi admin

**Database:** `auth.db`
- Table: `admins` - Data admin
- Table: `sessions` - Session tokens

**API Actions:**
- `verify` - Verifikasi apakah user adalah admin
- `login` - Create session token
- `logout` - Hapus session
- `list` - List semua admin
- `register` - Register admin baru
- `update_status` - Update status admin

**Key Features:**
- Token-based authentication
- Session management dengan expiry
- Admin verification dari `.vars.json`
- Auto cleanup expired sessions

---

### 2. Menu Service (Port 8082)
**Responsibility:** Manajemen menu makanan/minuman dan kategori

**Database:** `menu.db`
- Table: `menus` - Data menu items
- Table: `categories` - Kategori menu

**API Actions:**
- `create` - Tambah menu baru
- `read` - Baca detail menu
- `update` - Update menu
- `delete` - Hapus menu
- `list` - List menus dengan filter
- `list_categories` - List kategori
- `create_category` - Tambah kategori
- `delete_category` - Hapus kategori

**Key Features:**
- Filter by category
- Filter by availability
- Price validation
- Category management
- Photo URL support

---

### 3. Promo Service (Port 8083)
**Responsibility:** Manajemen promo dan diskon

**Database:** `promo.db`
- Table: `promos` - Data promosi

**API Actions:**
- `create` - Tambah promo baru
- `read` - Baca detail promo
- `update` - Update promo
- `delete` - Hapus promo
- `list` - List promos

**Key Features:**
- Percentage atau amount discount
- Date range validation
- Active/inactive status
- Auto-filter by date

---

### 4. Info Service (Port 8084)
**Responsibility:** Manajemen informasi café

**Database:** `info.db`
- Table: `cafe_info` - Informasi café (single row)

**API Actions:**
- `read` - Baca info café
- `update` - Update info café

**Key Features:**
- Single café info (singleton pattern)
- Operating hours
- Contact information
- Address & description

---

### 5. Media Service (Port 8085)
**Responsibility:** Manajemen file media (foto menu/promo)

**Database:** `media.db`
- Table: `media` - Metadata file media

**API Actions:**
- `create` - Upload/link media
- `read` - Baca detail media
- `list` - List media by entity
- `delete` - Hapus media

**Key Features:**
- Entity linkage (menu/promo)
- File metadata storage
- URL-based storage

---

### 6. Telegram Bot Agent
**Responsibility:** Interface dengan Telegram dan orchestration

**Components:**
- `main.go` - Initialization & config loading
- `handlers.go` - Message & callback handlers
- `menu_user.go` - User menu functions
- `menu_admin.go` - Admin menu functions

**Key Features:**
- User state management
- Dialog flow untuk CRUD
- Keyboard navigation
- Admin verification
- HTTP client untuk microservices
- Message formatting

## Communication Pattern

### Request Flow

```
User → Telegram → Agent → Microservice → Database
                          ↓
User ← Telegram ← Agent ← Microservice ← Database
```

### Standard Request Format
```json
{
  "action": "create|read|update|delete|list",
  "payload": {
    "key": "value"
  }
}
```

### Standard Response Format
```json
{
  "success": true/false,
  "data": {...},
  "error": {
    "code": "ERR_CODE",
    "message": "Error message"
  }
}
```

## Data Flow Examples

### Example 1: User melihat menu
```
1. User: Klik "Lihat Menu"
2. Agent: Detect callback "menu_category:Coffee"
3. Agent → Menu Service: POST {"action":"list", "payload":{"category":"Coffee"}}
4. Menu Service → DB: SELECT * FROM menus WHERE category='Coffee'
5. DB → Menu Service: Return menu list
6. Menu Service → Agent: {"success":true, "data":{"menus":[...]}}
7. Agent: Format response dengan emoji & harga
8. Agent → Telegram: Send formatted message
9. Telegram → User: Display menu
```

### Example 2: Admin tambah menu
```
1. Admin: Klik "Tambah Menu"
2. Agent: Set state = "add_menu_name"
3. Admin: Kirim "Cappuccino"
4. Agent: Save to tempData, state = "add_menu_price"
5. Admin: Kirim "25000"
6. Agent: Validate price, state = "add_menu_category"
7. Admin: Kirim "Coffee"
8. Agent: state = "add_menu_description"
9. Admin: Kirim "Kopi susu"
10. Agent → Menu Service: POST {"action":"create", "payload":{...}}
11. Menu Service: Validate & save to DB
12. Menu Service → Agent: {"success":true, "data":{"menu":{...}}}
13. Agent: Clear state & tempData
14. Agent → Telegram: "✅ Menu berhasil ditambahkan!"
```

## Database Schema

### auth.db
```sql
CREATE TABLE admins (
  id INTEGER PRIMARY KEY,
  telegram_id TEXT UNIQUE,
  username TEXT,
  is_active BOOLEAN,
  created_at DATETIME
);

CREATE TABLE sessions (
  id INTEGER PRIMARY KEY,
  admin_id INTEGER,
  token TEXT UNIQUE,
  expires_at DATETIME,
  created_at DATETIME,
  FOREIGN KEY (admin_id) REFERENCES admins(id)
);
```

### menu.db
```sql
CREATE TABLE categories (
  id INTEGER PRIMARY KEY,
  name TEXT UNIQUE,
  created_at DATETIME
);

CREATE TABLE menus (
  id INTEGER PRIMARY KEY,
  name TEXT,
  description TEXT,
  price INTEGER,
  category TEXT,
  photo_url TEXT,
  is_available BOOLEAN,
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (category) REFERENCES categories(name)
);
```

### promo.db
```sql
CREATE TABLE promos (
  id INTEGER PRIMARY KEY,
  title TEXT,
  description TEXT,
  discount INTEGER,
  discount_type TEXT,
  start_date DATETIME,
  end_date DATETIME,
  is_active BOOLEAN,
  created_at DATETIME,
  updated_at DATETIME
);
```

### info.db
```sql
CREATE TABLE cafe_info (
  id INTEGER PRIMARY KEY,
  name TEXT,
  address TEXT,
  phone TEXT,
  email TEXT,
  opening_hour TEXT,
  closing_hour TEXT,
  description TEXT,
  updated_at DATETIME
);
```

### media.db
```sql
CREATE TABLE media (
  id INTEGER PRIMARY KEY,
  file_name TEXT,
  file_url TEXT,
  file_type TEXT,
  entity_id INTEGER,
  entity_type TEXT,
  created_at DATETIME
);
```

## Shared Package

### `shared/database.go`
- `InitDB()` - Initialize SQLite connection
- `ExecuteSchema()` - Run SQL schema

### `shared/http_client.go`
- `NewHTTPClient()` - Create HTTP client
- `Post()` - Send POST request
- `Get()` - Send GET request
- Request/Response structs

### `shared/errors.go`
- Standard error codes
- Error constructors
- `AppError` struct

### `shared/logger.go`
- `LogInfo()` - Info logging
- `LogError()` - Error logging

### `shared/utils.go`
- `SanitizeInput()` - SQL injection prevention
- `ValidatePrice()` - Price validation
- `ValidateNotEmpty()` - Required field validation
- `ValidatePhotoURL()` - URL validation
- `FormatPrice()` - Format currency

## Security Considerations

1. **Input Sanitization**
   - All inputs sanitized in `shared/utils.go`
   - Remove SQL injection patterns
   - Trim whitespace

2. **Admin Authentication**
   - Bootstrap dari `.vars.json`
   - Session token verification
   - Token expiry (24 hours)

3. **Database**
   - Prepared statements (SQL injection safe)
   - Separate DB per service
   - Volume-mounted in Docker

4. **Configuration**
   - `.vars.json` in `.gitignore`
   - Environment variables for secrets
   - No hardcoded credentials

## Scalability

### Horizontal Scaling
- Each microservice dapat di-scale independently
- Stateless services (kecuali database)
- Load balancer dapat ditambahkan

### Vertical Scaling
- Increase resources per container
- Database optimization
- Caching layer (Redis) dapat ditambahkan

### Future Enhancements
- Message queue (RabbitMQ/Kafka) untuk async operations
- Centralized logging (ELK stack)
- Monitoring (Prometheus/Grafana)
- API Gateway
- Service mesh (Istio)

## Development Workflow

### Hot Reload with Air
```
1. Edit file .go
2. Air detects change
3. Recompile service
4. Restart process
5. Ready in < 2 seconds
```

### Docker Development
```
docker-compose up → All services start
Code change → Auto rebuild & restart
docker-compose logs -f → Monitor all services
```

### Local Development
```
make deps → Download dependencies
make build → Build all services
Run services individually for debugging
```

## Testing Strategy

### Unit Tests
- Test individual functions
- Mock external dependencies
- Coverage target: 80%

### Integration Tests
- Test service endpoints
- Test database operations
- Test service communication

### E2E Tests
- Test full user flows
- Test admin operations
- Test error handling

## Monitoring & Debugging

### Health Checks
```bash
# All services expose /health endpoint
curl http://localhost:8081/health
```

### Logs
```bash
# Docker logs
docker logs -f cafe-bot-agent

# Service logs
tail -f services/menu-service/build-errors.log
```

### Database Inspection
```bash
# Connect to DB
sqlite3 data/menu.db

# Query data
SELECT * FROM menus;
```

## Deployment Considerations

### Docker Compose (Development)
- Hot reload enabled
- Volume mounts for code
- Shared network

### Kubernetes (Production)
- Deployment per service
- Services for internal communication
- ConfigMaps for configuration
- Secrets for sensitive data
- Persistent volumes for databases

### CI/CD Pipeline
```
1. Git push
2. Run tests
3. Build Docker images
4. Push to registry
5. Deploy to cluster
6. Run smoke tests
```

## Design Principles

1. **Single Responsibility**: Each service has one clear purpose
2. **Separation of Concerns**: Agent doesn't handle business logic
3. **Fail Fast**: Validate early, return errors quickly
4. **Idempotency**: Operations can be repeated safely
5. **Graceful Degradation**: Service failures don't crash system
6. **Security by Design**: Input validation, sanitization, authentication
7. **Developer Experience**: Hot reload, clear logs, good documentation

---

This architecture provides:
- ✅ **Modularity**: Easy to add/remove/modify services
- ✅ **Scalability**: Scale services independently
- ✅ **Maintainability**: Clear separation of concerns
- ✅ **Testability**: Easy to test individual components
- ✅ **Developer Friendly**: Hot reload, clear structure
- ✅ **Production Ready**: Docker, monitoring, error handling
