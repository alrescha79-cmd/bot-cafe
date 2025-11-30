# 🚀 Quick Reference - Development Modes

## Cara Menjalankan Bot (Pilih Salah Satu)

### 1️⃣ TERCEPAT: Lokal Tanpa Docker (Recommended untuk Dev)
```bash
make dev-local
```
**Kelebihan:**
- ✅ Satu perintah, semua jalan
- ✅ Tidak perlu Docker
- ✅ Cepat startup (~3 detik)
- ✅ Log semua services dalam 1 terminal
- ✅ Tekan Ctrl+C untuk stop semua

**Kapan Digunakan:**
- Development sehari-hari
- Testing cepat
- Debugging
- Tidak mau install Docker

---

### 2️⃣ BEST: Lokal + Hot Reload (Tanpa Docker)
```bash
make dev-local-hot
```
**Kelebihan:**
- ✅ Satu perintah, semua jalan
- ✅ Auto reload saat edit file .go
- ✅ Tidak perlu Docker
- ✅ Development super cepat
- ✅ Log semua services dalam 1 terminal

**Kekurangan:**
- ⚠️ Perlu install Air (otomatis diinstall)
- ⚠️ Agak lebih lambat startup (~5 detik)

**Kapan Digunakan:**
- Development intensif
- Sering edit code
- Tidak mau restart manual
- Tidak mau pakai Docker

---

### 3️⃣ PRODUCTION-LIKE: Docker + Hot Reload
```bash
make dev
```
**Kelebihan:**
- ✅ Environment seperti production
- ✅ Auto reload saat edit file .go
- ✅ Isolated environment
- ✅ Easy deployment later

**Kekurangan:**
- ⚠️ Perlu Docker & Docker Compose
- ⚠️ Startup lebih lama (~10 detik)

**Kapan Digunakan:**
- Testing production environment
- Collaborative development
- CI/CD testing
- Final testing sebelum deploy

---

### 4️⃣ MANUAL: Build & Run Sendiri
```bash
make build
./bin/auth-service &
./bin/menu-service &
./bin/promo-service &
./bin/info-service &
./bin/media-service &
./bin/agent
```
**Kelebihan:**
- ✅ Full control
- ✅ Bisa debug per service
- ✅ Flexible

**Kekurangan:**
- ⚠️ Perlu banyak terminal/commands
- ⚠️ Manual restart tiap service
- ⚠️ Ribet

**Kapan Digunakan:**
- Debug spesifik service
- Development 1 service saja
- Troubleshooting

---

## 📊 Comparison Table

| Mode | Command | Docker | Hot Reload | Startup | Complexity | Use Case |
|------|---------|--------|------------|---------|------------|----------|
| **Lokal** | `make dev-local` | ❌ | ❌ | ⚡ Fast (3s) | 😊 Simple | Daily dev |
| **Lokal + Hot** | `make dev-local-hot` | ❌ | ✅ | ⚡ Medium (5s) | 😊 Simple | Intensive dev |
| **Docker + Hot** | `make dev` | ✅ | ✅ | 🐌 Slow (10s) | 😐 Medium | Production-like |
| **Manual** | Multiple commands | ❌ | ❌ | ⚡ Fast | 😫 Complex | Debugging |

---

## 🎯 Rekomendasi Berdasarkan Situasi

### Pertama Kali Setup
```bash
make init              # Setup config files
nano .env             # Tambah bot token
nano .vars.json       # Tambah admin ID
make deps             # Install dependencies
make dev-local        # Start bot!
```

### Development Sehari-hari
```bash
make dev-local-hot    # Start dengan hot reload
# Edit code, auto reload!
# Ctrl+C untuk stop
```

### Testing Feature Baru
```bash
make dev-local        # Quick start
# Test di Telegram
make stop            # Stop
```

### Sebelum Commit/Push
```bash
make dev              # Test di Docker
# Pastikan works di isolated env
make docker-down      # Stop
```

### Production Deployment
```bash
make docker-build     # Build images
make docker-up        # Deploy
```

---

## 🛑 Cara Stop Services

### Lokal (dev-local atau dev-local-hot)
```bash
# Di terminal yang running bot:
Ctrl+C

# Atau di terminal lain:
make stop
```

### Docker
```bash
make docker-down
```

---

## 🐛 Troubleshooting Quick Fix

### Port sudah dipakai
```bash
# Cek process yang pakai port
lsof -i :8081
lsof -i :8082
# dst...

# Kill process
kill -9 <PID>

# Atau kill semua sekaligus
make stop
```

### Air tidak terinstall (untuk dev-local-hot)
```bash
go install github.com/cosmtrek/air@latest
```

### Dependencies error
```bash
make deps
go mod tidy
```

### Database corrupt
```bash
make clean
# Start ulang
```

---

## 💡 Pro Tips

1. **Default recommendation**: Gunakan `make dev-local-hot` untuk development
2. **Quick test**: Gunakan `make dev-local` 
3. **Production test**: Gunakan `make dev` sebelum deploy
4. **One terminal**: Semua mode kecuali manual hanya butuh 1 terminal!
5. **Stop cepat**: Cukup `Ctrl+C` untuk stop semua services
6. **Check status**: Jika gagal start, cek port dengan `lsof -i :8081-8085`

---

## 📝 Summary

**Untuk 90% kasus, gunakan:**
```bash
make dev-local-hot
```

**Satu perintah, semua jalan, auto reload, tanpa Docker!** 🚀
