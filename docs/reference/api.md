# API Documentation - Bot Telegram Café

## Auth Service (Port 8081)

### Endpoint: POST /

#### Actions

##### 1. Verify Admin
**Request:**
```json
{
  "action": "verify",
  "payload": {
    "telegram_id": "123456789"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "is_admin": true,
    "admin": {
      "id": 1,
      "telegram_id": "123456789",
      "username": "admin",
      "is_active": true,
      "created_at": "2025-01-01T00:00:00Z"
    }
  }
}
```

##### 2. Login (Create Session)
**Request:**
```json
{
  "action": "login",
  "payload": {
    "telegram_id": "123456789"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "admin": {...},
    "session": {
      "id": 1,
      "admin_id": 1,
      "token": "abc123...",
      "expires_at": "2025-01-02T00:00:00Z",
      "created_at": "2025-01-01T00:00:00Z"
    }
  }
}
```

##### 3. Logout
**Request:**
```json
{
  "action": "logout",
  "payload": {
    "token": "abc123..."
  }
}
```

##### 4. List Admins
**Request:**
```json
{
  "action": "list"
}
```

##### 5. Register Admin
**Request:**
```json
{
  "action": "register",
  "payload": {
    "telegram_id": "987654321",
    "username": "newadmin"
  }
}
```

---

## Menu Service (Port 8082)

### Endpoint: POST /

#### Actions

##### 1. Create Menu
**Request:**
```json
{
  "action": "create",
  "payload": {
    "name": "Cappuccino",
    "description": "Kopi susu premium",
    "price": 25000,
    "category": "Coffee",
    "photo_url": "https://example.com/photo.jpg",
    "is_available": true
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "menu": {
      "id": 1,
      "name": "Cappuccino",
      "description": "Kopi susu premium",
      "price": 25000,
      "category": "Coffee",
      "photo_url": "https://example.com/photo.jpg",
      "is_available": true,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  }
}
```

##### 2. Read Menu
**Request:**
```json
{
  "action": "read",
  "payload": {
    "id": 1
  }
}
```

##### 3. Update Menu
**Request:**
```json
{
  "action": "update",
  "payload": {
    "id": 1,
    "name": "Cappuccino Premium",
    "price": 30000
  }
}
```

##### 4. Delete Menu
**Request:**
```json
{
  "action": "delete",
  "payload": {
    "id": 1
  }
}
```

##### 5. List Menus
**Request:**
```json
{
  "action": "list",
  "payload": {
    "category": "Coffee",      // optional
    "available_only": true     // optional
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "menus": [
      {...},
      {...}
    ]
  }
}
```

##### 6. List Categories
**Request:**
```json
{
  "action": "list_categories"
}
```

##### 7. Create Category
**Request:**
```json
{
  "action": "create_category",
  "payload": {
    "name": "Dessert"
  }
}
```

##### 8. Delete Category
**Request:**
```json
{
  "action": "delete_category",
  "payload": {
    "name": "Dessert"
  }
}
```

---

## Promo Service (Port 8083)

### Endpoint: POST /

#### Actions

##### 1. Create Promo
**Request:**
```json
{
  "action": "create",
  "payload": {
    "title": "Diskon 50%",
    "description": "Diskon spesial bulan ini",
    "discount": 50,
    "discount_type": "percentage",  // or "amount"
    "start_date": "2025-01-01",
    "end_date": "2025-01-31",
    "is_active": true
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "promo": {
      "id": 1,
      "title": "Diskon 50%",
      "description": "Diskon spesial bulan ini",
      "discount": 50,
      "discount_type": "percentage",
      "start_date": "2025-01-01T00:00:00Z",
      "end_date": "2025-01-31T00:00:00Z",
      "is_active": true,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  }
}
```

##### 2. Read Promo
**Request:**
```json
{
  "action": "read",
  "payload": {
    "id": 1
  }
}
```

##### 3. Update Promo
**Request:**
```json
{
  "action": "update",
  "payload": {
    "id": 1,
    "title": "Diskon 60%",
    "discount": 60
  }
}
```

##### 4. Delete Promo
**Request:**
```json
{
  "action": "delete",
  "payload": {
    "id": 1
  }
}
```

##### 5. List Promos
**Request:**
```json
{
  "action": "list",
  "payload": {
    "active_only": true  // optional
  }
}
```

---

## Info Service (Port 8084)

### Endpoint: POST /

#### Actions

##### 1. Read Café Info
**Request:**
```json
{
  "action": "read"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "info": {
      "id": 1,
      "name": "Café Bot",
      "address": "Jl. Example No. 123",
      "phone": "081234567890",
      "email": "cafe@example.com",
      "opening_hour": "08:00",
      "closing_hour": "22:00",
      "description": "Selamat datang!",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  }
}
```

##### 2. Update Café Info
**Request:**
```json
{
  "action": "update",
  "payload": {
    "name": "Café Bot Premium",
    "address": "Jl. New Address",
    "phone": "081234567890",
    "email": "newcafe@example.com",
    "opening_hour": "07:00",
    "closing_hour": "23:00",
    "description": "Updated description"
  }
}
```

---

## Media Service (Port 8085)

### Endpoint: POST /

#### Actions

##### 1. Create Media
**Request:**
```json
{
  "action": "create",
  "payload": {
    "file_name": "cappuccino.jpg",
    "file_url": "https://example.com/cappuccino.jpg",
    "file_type": "image/jpeg",
    "entity_id": 1,
    "entity_type": "menu"
  }
}
```

##### 2. Read Media
**Request:**
```json
{
  "action": "read",
  "payload": {
    "id": 1
  }
}
```

##### 3. List Media by Entity
**Request:**
```json
{
  "action": "list",
  "payload": {
    "entity_id": 1,
    "entity_type": "menu"
  }
}
```

##### 4. Delete Media
**Request:**
```json
{
  "action": "delete",
  "payload": {
    "id": 1
  }
}
```

---

## Error Codes

| Code | Description |
|------|-------------|
| `ERR_INVALID_INPUT` | Input tidak valid atau tidak lengkap |
| `ERR_NOT_FOUND` | Resource tidak ditemukan |
| `ERR_UNAUTHORIZED` | Akses tidak diizinkan |
| `ERR_DATABASE` | Kesalahan database |
| `ERR_INTERNAL` | Kesalahan internal server |
| `ERR_DUPLICATE` | Data duplikat |
| `ERR_SERVICE` | Kesalahan microservice |

---

## Testing with cURL

### Test Menu Service
```bash
# List all menus
curl -X POST http://localhost:8082 \
  -H "Content-Type: application/json" \
  -d '{"action":"list","payload":{}}'

# Create menu
curl -X POST http://localhost:8082 \
  -H "Content-Type: application/json" \
  -d '{
    "action":"create",
    "payload":{
      "name":"Test Menu",
      "price":15000,
      "category":"Coffee",
      "is_available":true
    }
  }'
```

### Test Auth Service
```bash
# Verify admin
curl -X POST http://localhost:8081 \
  -H "Content-Type: application/json" \
  -d '{
    "action":"verify",
    "payload":{
      "telegram_id":"123456789"
    }
  }'
```

### Health Check
```bash
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
curl http://localhost:8084/health
curl http://localhost:8085/health
```
