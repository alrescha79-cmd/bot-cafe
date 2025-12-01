# Development Setup Guide

Complete guide untuk setup development environment bot café.

## Prerequisites

### Required
- **Go 1.21+** - [Download](https://go.dev/dl/)
- **Git** - Version control
- **Telegram Bot Token** - Dari [@BotFather](https://t.me/botfather)

### Optional (pilih salah satu)
- **Docker & Docker Compose** - Untuk development dengan container
- **Air v1.49.0** - Untuk hot reload lokal tanpa Docker

## 🚀 Initial Setup

### 1. Clone Repository

```bash
git clone https://github.com/alrescha79-cmd/bot-cafe.git
cd bot-cafe
```

### 2. Initialize Project

```bash
make init
```

Perintah ini akan:
- Copy `.env.example` → `.env`
- Copy `.vars.json.example` → `.vars.json`
- Create `data/` directory
- Download Go dependencies

### 3. Configure Environment

#### `.env` File

Edit `.env` dan set bot token Anda:

```env
TELEGRAM_BOT_TOKEN=your_bot_token_from_botfather
```

**Cara mendapat bot token:**
1. Chat dengan [@BotFather](https://t.me/botfather)
2. Ketik `/newbot`
3. Ikuti instruksi (nama bot, username)
4. Copy token yang diberikan

#### `.vars.json` File

Edit `.vars.json` dan set admin Telegram IDs:

```json
{
  "admin_telegram_ids": [
    "123456789"
  ],
  "admin_usernames": [
    "your_telegram_username"
  ]
}
```

**Cara mendapat Telegram ID Anda:**
1. Chat dengan [@userinfobot](https://t.me/userinfobot)
2. Bot akan kirim ID Anda
3. Copy angka tersebut ke `.vars.json`

> ⚠️ **Penting**: File `.vars.json` sudah ada di `.gitignore` dan tidak akan di-commit.

## 🎯 Development Modes

Ada 3 cara menjalankan bot untuk development:

### Mode 1: Lokal dengan Hot Reload (Recommended)

**Best for**: Development harian, fast iteration

```bash
make dev-local-hot
```

**Features:**
- ✅ Hot reload otomatis saat edit code
- ✅ Tidak perlu Docker
- ✅ Fast rebuild (hanya rebuild yang berubah)
- ✅ Native performance

**Installation (otomatis):**

Script akan otomatis install Air v1.49.0 jika belum ada. Atau install manual:

```bash
make install-air
# atau
go install github.com/cosmtrek/air@v1.49.0
```

**How it works:**
- Air watch semua file `.go` di setiap service
- Saat file berubah → Auto rebuild → Auto restart
- Juga watch direktori `shared/` untuk shared code

**Logs:**
```
[AUTH] watching .
[AUTH] watching ../../tmp/../shared
[MENU] building...
[MENU] running...
✅ All services started with HOT RELOAD!
```

**Stop:**
```bash
Ctrl+C
# atau di terminal lain:
make stop
```

### Mode 2: Docker dengan Hot Reload

**Best for**: Production-like environment, consistent environment

```bash
make dev
```

**Features:**
- ✅ Hot reload via volume mounting
- ✅ Isolated environment (seperti production)
- ✅ Tidak perlu install Go lokal
- ✅ Consistent di semua OS

**How it works:**
- Code di-mount ke container sebagai volume
- Air watch perubahan file dari dalam container
- Auto rebuild dan restart service

**Useful commands:**
```bash
# Lihat logs semua services
make docker-logs

# Lihat logs specific service
docker logs -f cafe-bot-agent
docker logs -f cafe-menu-service

# Restart containers
docker restart cafe-bot-agent

# Stop all
make docker-down
```

### Mode 3: Lokal Tanpa Hot Reload

**Best for**: Quick testing, debugging specific issues

```bash
make dev-local
```

**Features:**
- ⚡ Paling cepat untuk start
- ✅ Semua services running di satu terminal
- ✅ Simple, no dependencies

**How it works:**
- Build semua services
- Running semua dalam background processes
- Semua logs muncul di satu terminal

**Stop:**
```bash
Ctrl+C      # Di terminal yang running
make stop   # Dari terminal lain
```

## 📊 Development Workflow

### Typical Workflow (Hot Reload)

1. **Start services:**
   ```bash
   make dev-local-hot
   ```

2. **Edit code** - Misalnya edit `agent/handlers.go`

3. **Save file** - Air auto-detect dan rebuild
   ```
   [AGENT] building...
   [AGENT] running...
   ```

4. **Test di Telegram** - Changes langsung apply!

5. **Repeat** - Edit, save, test. No restart needed.

### Working with Shared Code

Ketika edit file di `shared/` (e.g., `shared/logger.go`):

```
[AUTH] building...
[MENU] building...
[AGENT] building...
```

Semua service yang import `shared` akan rebuild otomatis!

### Adding New Features

```bash
# 1. Create feature branch
git checkout -b feature/add-payment

# 2. Edit code dengan hot reload running
make dev-local-hot

# 3. Test di Telegram
# Bot auto-reload setiap kali save

# 4. Commit changes
git add .
git commit -m "Add payment feature"
```

## 🐛 Debugging

### Check Service Health

```bash
# Health check endpoints
curl http://localhost:8081/health  # auth-service
curl http://localhost:8082/health  # menu-service
curl http://localhost:8083/health  # promo-service
curl http://localhost:8084/health  # info-service
curl http://localhost:8085/health  # media-service
```

### Test API Manually

```bash
# List all menus
curl -X POST http://localhost:8082 \
  -H "Content-Type: application/json" \
  -d '{"action":"list","payload":{}}'

# Get cafe info
curl -X POST http://localhost:8084 \
  -H "Content-Type: application/json" \
  -d '{"action":"read","payload":{}}'
```

### View Logs

**Hot Reload Lokal:**
- Logs muncul langsung di terminal
- Setiap service punya prefix: `[AUTH]`, `[MENU]`, dll

**Docker:**
```bash
# All logs
make docker-logs

# Specific service
docker logs -f cafe-auth-service
docker logs -f cafe-bot-agent
```

### Debug dengan VS Code

Add `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Agent",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/agent",
      "env": {
        "TELEGRAM_BOT_TOKEN": "your_token_here"
      }
    }
  ]
}
```

## 🧪 Testing

```bash
# Run all tests
make test

# Test specific service
cd services/menu-service && go test -v

# Test with coverage
go test -cover ./...

# Test specific function
go test -v -run TestCreateMenu
```

## 🔧 Troubleshooting

### Port Already in Use

```bash
# Check what's using ports
lsof -i :8081-8085

# Kill stuck processes
make stop

# Or kill specific port
kill -9 $(lsof -t -i:8081)
```

### Air Command Not Found

```bash
# Check if Air installed
which air

# Add Go bin to PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Make it permanent (add to ~/.zshrc or ~/.bashrc)
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc

# Install Air
make install-air
```

### Hot Reload Not Working

| Issue | Solution |
|-------|----------|
| File changes not detected | Check Air is watching correct directories |
| Air crashes on rebuild | Check Go syntax errors in logs |
| Shared code changes ignored | Air config includes `shared/` directory (fixed in v1.49.0) |

See [Troubleshooting Guide](troubleshooting.md#hot-reload-issues) for more solutions.

### Database Issues

```bash
# Reset all databases
make clean
make dev-local-hot

# Check database location
ls -la data/
# Should show: auth.db, menu.db, promo.db, info.db, media.db
```

### Admin Access Denied

1. **Verify Telegram ID di `.vars.json`**
   ```bash
   cat .vars.json
   ```

2. **Check agent logs**
   ```bash
   # Look for authentication logs
   grep "AUTH" logs
   ```

3. **Restart agent**
   ```bash
   make stop
   make dev-local-hot
   ```

## 🚀 Next Steps

- **Add new features** - Edit code, test dengan hot reload
- **Read [API Documentation](../reference/api.md)** - Understand service APIs
- **Check [Admin Workflows](../examples/admin-workflows.md)** - Learn admin features
- **Ready for production?** - See [VPS Deployment Guide](vps-deployment.md)

## 📝 Tips & Best Practices

### 1. Always Use Hot Reload
```bash
# Instead of manual restart:
make dev-local-hot

# NOT:
make build && ./bin/agent
```

### 2. Test in Telegram Frequently
- Keep Telegram open while developing
- Test each change immediately
- Use `/cancel` to reset dialog states

### 3. Check Logs Often
- Watch logs untuk errors
- Understanding log flow helps debugging
- Each service logs dengan prefix

### 4. Commit Often
```bash
git add .
git commit -m "Descriptive message"
```

### 5. Keep Dependencies Updated
```bash
make deps
go mod tidy
```

---

**Happy Hacking!** 🎉

Jika ada masalah, check [Troubleshooting Guide](troubleshooting.md) atau buka issue di GitHub.
