# 🤖 Bot Telegram Café/Resto

Bot Telegram berbasis **microservices** untuk mengelola café/resto dengan fitur lengkap untuk admin dan pelanggan.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat\u0026logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## ✨ Fitur Utama

### 👥 Untuk Pelanggan
- 📋 Lihat menu lengkap dengan kategori
- 💰 Cek harga dan deskripsi produk
- 🎉 Info promo dan diskon terkini
- ℹ️ Informasi café (alamat, jam buka, kontak)

### 👨‍💼 Untuk Admin
- ➕ **CRUD Menu** - Kelola menu dan kategori
- 🎁 **CRUD Promo** - Buat dan kelola promo
- ℹ️ **CRUD Info Café** - Update informasi café
- 👥 **Multi-Admin** - Dukungan multiple admin
- 🔐 **Autentikasi Aman** - Admin via Telegram ID

## 🚀 Quick Start

### 1. Setup Project

```bash
# Clone repository
git clone https://github.com/alrescha79-cmd/bot-cafe.git
cd bot-cafe

# Initialize (creates .env and .vars.json from examples)
make init

# Configure
nano .env          # Set TELEGRAM_BOT_TOKEN
nano .vars.json    # Set admin Telegram IDs
```

### 2. Pilih Development Mode

**Option A: Lokal dengan Hot Reload (Recommended)**
```bash
make dev-local-hot
```

**Option B: Docker dengan Hot Reload**
```bash
make dev
```

**Option C: Lokal Tanpa Hot Reload**
```bash
make dev-local
```

### 3. Mulai Gunakan Bot

Buka Telegram, cari bot Anda, dan ketik `/start`

✅ Selesai! Bot siap digunakan.

## 📚 Dokumentasi Lengkap

- **[Development Setup](docs/guides/development-setup.md)** - Panduan setup development lengkap
- **[VPS Deployment](docs/guides/vps-deployment.md)** - Deploy ke VPS production
- **[Troubleshooting](docs/guides/troubleshooting.md)** - Solusi masalah umum
- **[Architecture](docs/reference/architecture.md)** - Arsitektur microservices
- **[API Reference](docs/reference/api.md)** - Dokumentasi API
- **[Makefile Commands](docs/reference/makefile-commands.md)** - Referensi command
- **[Admin Workflows](docs/examples/admin-workflows.md)** - Contoh penggunaan admin

## 🏗️ Arsitektur

Aplikasi ini menggunakan **5 microservices**:

| Service | Port | Fungsi |
|---------|------|--------|
| auth-service | 8081 | Autentikasi admin |
| menu-service | 8082 | Manajemen menu & kategori |
| promo-service | 8083 | Manajemen promo |
| info-service | 8084 | Informasi café |
| media-service | 8085 | Upload media (optional) |

Plus **1 agent** (Telegram Bot) yang berkomunikasi dengan semua services.

Setiap service memiliki database SQLite sendiri untuk independensi dan scalability.

## 🛠️ Tech Stack

- **Backend**: Go 1.21+
- **Database**: SQLite3 (per service)
- **Bot Framework**: go-telegram-bot-api
- **Containerization**: Docker & Docker Compose
- **Hot Reload**: Air v1.49.0

## 🔧 Commands Penting

```bash
# Setup & Dependencies
make init           # Setup awal project
make deps           # Install Go dependencies

# Development
make dev-local-hot  # Run lokal dengan hot reload (tanpa Docker)
make dev            # Run dengan Docker + hot reload
make dev-local      # Run lokal tanpa hot reload
make stop           # Stop semua services

# Build & Clean
make build          # Build semua services
make clean          # Bersihkan build artifacts

# Deployment
make docker-build   # Build Docker images
make docker-up      # Start containers
make docker-down    # Stop containers
make docker-logs    # Lihat logs
```

Lihat [Makefile Commands Reference](docs/reference/makefile-commands.md) untuk command lengkap.

## 🔥 Hot Reload

Aplikasi mendukung **automatic hot reload** saat development:

- Edit file `.go` → Auto rebuild → Service restart
- **Lokal**: Menggunakan [Air](https://github.com/cosmtrek/air) v1.49.0
- **Docker**: Built-in volume mounting

No manual restart needed! 🎉

## 🗄️ Database

Setiap microservice memiliki database terpisah:

- `data/auth.db` - Admin & sessions
- `data/menu.db` - Menu & kategori
- `data/promo.db` - Promo & diskon
- `data/info.db` - Info café
- `data/media.db` - Media files

## 🔐 Security

- ✅ Admin authentication via Telegram ID
- ✅ Session-based admin verification
- ✅ Input sanitization (SQL injection prevention)
- ✅ Environment variables untuk secrets
- ✅ `.vars.json` tidak masuk git (`.gitignore`)

## 📊 Project Structure

```
bot-cafe/
├── agent/              # Telegram Bot
├── services/           # Microservices
│   ├── auth-service/
│   ├── menu-service/
│   ├── promo-service/
│   ├── info-service/
│   └── media-service/
├── shared/             # Shared utilities
├── deployments/        # Docker configs
├── docs/               # Documentation
└── Makefile           # Build commands
```

## 🐛 Troubleshooting

**Bot tidak respond?**
```bash
# Check logs
make docker-logs

# Verify bot token
cat .env | grep TELEGRAM_BOT_TOKEN
```

**Access denied untuk admin?**
```bash
# Verify admin ID di .vars.json
cat .vars.json

# Restart agent
docker restart cafe-bot-agent
# atau untuk lokal: make stop && make dev-local-hot
```

**Hot reload tidak bekerja?**

Lihat [Troubleshooting Guide](docs/guides/troubleshooting.md#hot-reload-issues) untuk solusi lengkap.

## 🚀 Roadmap

- [ ] Broadcast message ke users
- [ ] Upload foto menu via bot
- [ ] Order tracking system
- [ ] Payment gateway integration
- [ ] Analytics dashboard
- [ ] Multi-language support

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

- **Anggun Caksono** - [Github/alrescha79-cmd](https://github.com/alrescha79-cmd)

Built with ❤️ using microservices architecture best practices.

---

**Happy Coding!** 🚀

For questions or issues, please open an issue on GitHub.
