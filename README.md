# 🤖 Bot Telegram Café/Resto

Bot Telegram berbasis **arsitektur microservices** untuk mengelola café/resto dengan fitur lengkap untuk admin dan user.

## 📋 Fitur

### User Features
- 📋 Melihat daftar menu (makanan & minuman)
- 💰 Melihat harga menu
- 🎉 Melihat promo aktif
- ℹ️ Melihat info café (alamat, jam buka, kontak)
- 📱 Navigasi mudah via keyboard Telegram

### Admin Features
- ➕ CRUD Menu (Create, Read, Update, Delete)
- 🎁 CRUD Promo
- ℹ️ CRUD Info Café
- 📁 CRUD Kategori Menu
- 👨‍💼 Multi-admin support
- 🔐 Autentikasi admin via `.vars.json`

## 🏗️ Arsitektur

Aplikasi ini dibangun dengan **5 microservices**:

1. **auth-service** (Port 8081) - Manajemen autentikasi admin
2. **menu-service** (Port 8082) - Manajemen menu dan kategori
3. **promo-service** (Port 8083) - Manajemen promo
4. **info-service** (Port 8084) - Manajemen info café
5. **media-service** (Port 8085) - Manajemen media/foto (optional)

Plus **1 agent** (Bot Telegram) yang berkomunikasi dengan semua microservices.

## 📁 Struktur Proyek

```
/bot-cafe
├── /agent                    # Telegram Bot Agent
│   ├── main.go
│   ├── handlers.go
│   ├── menu_user.go
│   └── menu_admin.go
├── /services
│   ├── /auth-service
│   ├── /menu-service
│   ├── /promo-service
│   ├── /info-service
│   └── /media-service
├── /shared                   # Shared utilities
│   ├── database.go
│   ├── http_client.go
│   ├── errors.go
│   ├── logger.go
│   └── utils.go
├── /deployments
│   ├── docker-compose.yml
│   ├── Dockerfile.service
│   └── Dockerfile.agent
├── .env.example
├── .vars.json.example
├── Makefile
└── README.md
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (untuk development tanpa Docker)
- Telegram Bot Token (dari [@BotFather](https://t.me/botfather))

### 1. Clone & Setup

```bash
# Clone repository
git clone <repo-url>
cd bot-cafe

# Initialize project
make init

# Edit konfigurasi
nano .env          # Set TELEGRAM_BOT_TOKEN
nano .vars.json    # Set admin Telegram IDs
```

### 2. Konfigurasi

#### `.env`
```env
TELEGRAM_BOT_TOKEN=your_bot_token_here
```

#### `.vars.json`
```json
{
  "admin_telegram_ids": [
    "123456789"
  ],
  "admin_usernames": [
    "your_username"
  ]
}
```

> **Penting**: File `.vars.json` tidak akan di-commit ke git (ada di `.gitignore`)

### 3. Pilih Mode Development

#### Option A: Dengan Docker (Recommended - Hot Reload)

```bash
# Start semua services dengan hot reload
make dev
```

Bot akan otomatis reload saat ada perubahan code di dalam Docker!

#### Option B: Lokal Tanpa Docker (Satu Perintah)

```bash
# Install dependencies (pertama kali)
make deps

# Jalankan semua services sekaligus
make run-local
# atau
make dev-local
```

Semua services akan berjalan bersamaan dalam satu terminal. Tekan `Ctrl+C` untuk stop semua.

#### Option C: Lokal Dengan Hot Reload (Tanpa Docker)

```bash
# Install Air jika belum (akan otomatis jika belum ada)
go install github.com/cosmtrek/air@latest

# Jalankan dengan hot reload lokal
make dev-local-hot
```

Bot akan otomatis reload saat ada perubahan code, **tanpa Docker**!

#### Option D: Manual (Build & Run)

```bash
# Build semua services
make build

# Jalankan (butuh terminal terpisah untuk setiap service)
./bin/auth-service &
./bin/menu-service &
./bin/promo-service &
./bin/info-service &
./bin/media-service &
./bin/agent
```

## 🎮 Cara Menggunakan Bot

### Sebagai User

1. Start bot: `/start`
2. Pilih menu:
   - 📋 Lihat Menu
   - 🎉 Lihat Promo
   - ℹ️ Info Café

### Sebagai Admin

1. Pastikan Telegram ID Anda ada di `.vars.json`
2. Start bot: `/start`
3. Klik **👨‍💼 Panel Admin**
4. Pilih menu admin:
   - Kelola Menu
   - Kelola Promo
   - Kelola Info Café
   - Kelola Kategori

## 🔧 Commands (Makefile)

```bash
# Setup & Dependencies
make help           # Lihat semua commands
make init           # Setup awal project
make deps           # Install Go dependencies

# Development - Lokal (Tanpa Docker)
make run-local      # Jalankan semua services lokal (satu perintah)
make dev-local      # Alias untuk run-local
make dev-local-hot  # Jalankan lokal dengan hot reload (tanpa Docker)
make stop           # Stop semua services lokal

# Development - Docker (Dengan Hot Reload)
make dev            # Start dev environment dengan Docker hot reload
make docker-up      # Start containers
make docker-down    # Stop containers
make docker-logs    # Lihat logs
make docker-build   # Build Docker images

# Build & Test
make build          # Build semua services
make test           # Run tests
make clean          # Bersihkan build artifacts
```

### 🌟 Rekomendasi Command

- **Pertama kali**: `make init` → edit `.env` & `.vars.json`
- **Development cepat tanpa Docker**: `make dev-local` (satu perintah, semua jalan!)
- **Development dengan hot reload tanpa Docker**: `make dev-local-hot`
- **Development dengan Docker**: `make dev`
- **Stop services lokal**: `make stop` atau `Ctrl+C`

## 📡 API Contract

Semua microservices menggunakan format request/response yang sama:

### Request Format
```json
{
  "action": "create|read|update|delete|list",
  "payload": {
    "key": "value"
  }
}
```

### Response Format (Success)
```json
{
  "success": true,
  "data": {
    "result": "data"
  }
}
```

### Response Format (Error)
```json
{
  "success": false,
  "error": {
    "code": "ERR_CODE",
    "message": "User-friendly error message"
  }
}
```

## 🗄️ Database

Setiap microservice memiliki database SQLite sendiri:

- `auth.db` - Data admin & sessions
- `menu.db` - Data menu & kategori
- `promo.db` - Data promo
- `info.db` - Data info café
- `media.db` - Data media/foto

Database disimpan di volume Docker `/data`

## 🔐 Security

- Admin ID disimpan di `.vars.json` (tidak masuk git)
- Session token untuk admin
- Input sanitization untuk mencegah SQL injection
- Validasi input untuk semua operasi CRUD

## 🐳 Hot Reload (Development)

Hot reload menggunakan [Air](https://github.com/cosmtrek/air):

- **Otomatis reload** saat file `.go` berubah
- **Tidak perlu restart** Docker container
- **Cepat** - hanya recompile service yang berubah

Edit code → Save → Auto reload! 🔥

## 📝 Contoh Alur CRUD Menu (Admin)

1. Admin: klik "Kelola Menu" → "Tambah Menu"
2. Bot: "Masukkan nama menu"
3. Admin: "Cappuccino"
4. Bot: "Masukkan harga"
5. Admin: "25000"
6. Bot: "Pilih kategori"
7. Admin: "Coffee"
8. Bot: "Masukkan deskripsi (- untuk skip)"
9. Admin: "Kopi susu premium"
10. Bot: "✅ Menu berhasil ditambahkan!"

## 🧪 Testing

```bash
# Run all tests
make test

# Test specific service
cd services/menu-service && go test

# Test with coverage
go test -cover ./...
```

## 📊 Monitoring & Logs

```bash
# Lihat logs semua services
make docker-logs

# Logs service tertentu
docker logs cafe-auth-service
docker logs cafe-menu-service
docker logs cafe-bot-agent

# Logs realtime
docker logs -f cafe-bot-agent
```

## 🛠️ Development Tips

### Menambah Fitur Baru

1. Buat endpoint di microservice terkait
2. Update handler di agent
3. Tambahkan UI/keyboard di bot
4. Test perubahan (auto reload!)

### Debugging

```bash
# Check service health
curl http://localhost:8081/health
curl http://localhost:8082/health

# Test API manually
curl -X POST http://localhost:8082 \
  -H "Content-Type: application/json" \
  -d '{"action":"list","payload":{}}'
```

## 🚧 Troubleshooting

### Bot tidak respond
- Cek `TELEGRAM_BOT_TOKEN` di `.env`
- Cek logs: `make docker-logs`

### "Access denied" untuk admin
- Cek Telegram ID di `.vars.json`
- Restart agent: `docker restart cafe-bot-agent`

### Hot reload tidak bekerja
- Cek volume mounts di `docker-compose.yml`
- Restart containers: `make docker-restart`

## 🗺️ Roadmap

- [ ] Broadcast message ke semua users
- [ ] Upload foto menu via bot
- [ ] Laporan penjualan
- [ ] Integrasi payment gateway
- [ ] Multi-language support

## 📄 License

MIT License

## 👤 Author

Built with ❤️ following microservices architecture best practices

## 🤝 Contributing

1. Fork the repo
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

---

**Happy Coding!** 🚀

Jika ada pertanyaan, silakan buka issue di GitHub.
