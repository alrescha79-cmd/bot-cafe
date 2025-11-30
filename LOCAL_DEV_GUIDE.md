# 🎯 Local Development Setup - Complete!

## ✅ What's New

Sekarang Anda bisa menjalankan **seluruh bot dengan SATU PERINTAH** tanpa Docker!

---

## 🚀 Cara Menggunakan

### Option 1: Tanpa Hot Reload (Tercepat)
```bash
make dev-local
```
atau
```bash
make run-local
```

**Apa yang terjadi:**
- ✅ Semua 5 microservices start bersamaan
- ✅ Bot Telegram agent start otomatis
- ✅ Semua logs dalam 1 terminal (dengan prefix `[AUTH]`, `[MENU]`, dll)
- ✅ Tekan `Ctrl+C` untuk stop semua sekaligus

**Startup time:** ~3 detik

---

### Option 2: Dengan Hot Reload (Best for Development)
```bash
make dev-local-hot
```

**Apa yang terjadi:**
- ✅ Semua services start dengan Air (hot reload tool)
- ✅ Edit file `.go` apa saja → Auto rebuild & restart service
- ✅ Tidak perlu restart manual
- ✅ Tidak perlu Docker
- ✅ Air akan diinstall otomatis jika belum ada

**Startup time:** ~5 detik
**Reload time:** ~2 detik per service

---

## 📁 Files Created

### 1. `scripts/run-local.sh`
Script untuk menjalankan semua services secara bersamaan tanpa Docker.

**Features:**
- Auto-check `.env` dan `.vars.json`
- Start semua 5 microservices di background
- Start bot agent
- Colored output dengan prefix per service
- Cleanup otomatis saat Ctrl+C

### 2. `scripts/run-local-hot.sh`
Script untuk menjalankan dengan Air (hot reload) tanpa Docker.

**Features:**
- Auto-install Air jika belum ada
- Generate `.air.toml` config untuk tiap service
- Hot reload untuk semua services
- Auto cleanup config saat stop

### 3. Updated `Makefile`
Tambahan commands:
- `make run-local` - Run tanpa hot reload
- `make dev-local` - Alias untuk run-local
- `make dev-local-hot` - Run dengan hot reload
- `make stop` - Stop semua services

### 4. `QUICK_START.md`
Dokumentasi lengkap dengan:
- Comparison table semua mode
- Rekomendasi berdasarkan use case
- Troubleshooting guide
- Pro tips

---

## 🎮 Example Usage

### Pertama Kali
```bash
# Setup (hanya sekali)
make init
nano .env              # Tambah bot token
nano .vars.json        # Tambah admin ID
make deps              # Install dependencies

# Run!
make dev-local-hot
```

### Development Sehari-hari
```bash
# Start bot
make dev-local-hot

# Edit code di services/menu-service/handlers.go
# Save → Auto reload!

# Test di Telegram

# Ctrl+C untuk stop
```

---

## 📊 Comparison: Before vs After

### BEFORE (Ribet)
```bash
# Terminal 1
cd services/auth-service && go run .

# Terminal 2
cd services/menu-service && go run .

# Terminal 3
cd services/promo-service && go run .

# Terminal 4
cd services/info-service && go run .

# Terminal 5
cd services/media-service && go run .

# Terminal 6
cd agent && go run .

# Need 6 terminals! 😫
```

### AFTER (Easy!)
```bash
make dev-local-hot
# Done! One command, one terminal! 🎉
```

---

## 🔥 Hot Reload Demo

```bash
# Start bot
make dev-local-hot

# Edit file
nano services/menu-service/handlers.go
# Change something, save

# Output (otomatis):
[MENU] Detected file change...
[MENU] Rebuilding...
[MENU] Build successful!
[MENU] Restarting...
[MENU] Ready to serve!

# Ready in ~2 seconds!
```

---

## 🛑 Stop Services

### Lokal
```bash
# Method 1: Ctrl+C di terminal yang running
Ctrl+C

# Method 2: Di terminal lain
make stop
```

### Docker (kalau pakai)
```bash
make docker-down
```

---

## 🎯 When to Use What?

| Situation | Command | Why |
|-----------|---------|-----|
| Daily development | `make dev-local-hot` | Fast + auto reload |
| Quick test | `make dev-local` | Fastest startup |
| Production-like test | `make dev` | Docker environment |
| Debug specific service | Manual run | Full control |

---

## 💡 Pro Tips

1. **Default**: Gunakan `make dev-local-hot`
2. **Logs**: Semua service punya prefix (`[AUTH]`, `[MENU]`, dll) untuk mudah dibaca
3. **Port conflict**: Jalankan `make stop` sebelum start ulang
4. **Dependencies**: Jalankan `make deps` kalau ada error package
5. **Clean start**: `make clean` untuk hapus DB dan mulai fresh

---

## 🐛 Troubleshooting

### "Port already in use"
```bash
make stop
# atau
lsof -i :8081
kill -9 <PID>
```

### "Air not found"
```bash
go install github.com/cosmtrek/air@latest
# Atau biarkan script install otomatis
```

### "Cannot connect to service"
```bash
# Tunggu 3-5 detik setelah start
# Services butuh waktu init
```

### Services tidak auto-reload
```bash
# Pastikan edit file .go
# Air hanya watch file .go, bukan .json/.env
# Untuk .env changes, restart manual (Ctrl+C, start ulang)
```

---

## 📚 Documentation

Lihat dokumentasi lengkap di:
- `QUICK_START.md` - Quick reference semua mode
- `README.md` - Complete guide
- `DEPLOYMENT.md` - Production deployment

---

## 🎉 Summary

Sekarang Anda punya **3 cara mudah** untuk development:

1. **`make dev-local`** → Fast, no Docker, no hot reload
2. **`make dev-local-hot`** → Fast, no Docker, **WITH hot reload** ⭐
3. **`make dev`** → Docker, with hot reload

**Recommendation: `make dev-local-hot`** untuk development sehari-hari! 🚀

---

**Happy Coding!** 💻✨
