# 🔧 Troubleshooting Guide - Hot Reload

## ❌ Error: Air Installation Failed

### Problem
```
go: github.com/cosmtrek/air@latest requires go >= 1.25
go: version constraints conflict
```

### Root Cause
Air telah berpindah dari `github.com/cosmtrek/air` ke `github.com/air-verse/air` dan memerlukan Go 1.25+, sedangkan project ini menggunakan Go 1.21.

### Solutions

#### Solution 1: Gunakan Air v1.49.0 (Recommended) ✅
```bash
# Install versi kompatibel
go install github.com/cosmtrek/air@v1.49.0

# Atau gunakan make command
make install-air

# Pastikan PATH sudah diset
export PATH="$PATH:$(go env GOPATH)/bin"

# Jalankan
make dev-local-hot
```

#### Solution 2: Gunakan Alternative Hot Reload (inotifywait/fswatch)
```bash
# Linux (Ubuntu/Debian)
sudo apt install inotify-tools

# macOS
brew install fswatch

# Jalankan
make dev-local-watch
```

#### Solution 3: Run Tanpa Hot Reload
```bash
# Paling simple, no dependencies
make dev-local
```

#### Solution 4: Gunakan Docker (sudah built-in hot reload)
```bash
make dev
```

---

## ❌ Error: Air Command Not Found

### Problem
```
air: command not found
```

### Solution
```bash
# Check apakah Air terinstall
ls ~/go/bin/air

# Jika ada, tambahkan ke PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Atau tambahkan ke .bashrc/.zshrc
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc

# Jika tidak ada, install
make install-air
```

---

## ❌ Error: Port Already in Use

### Problem
```
bind: address already in use
```

### Solution
```bash
# Method 1: Gunakan make stop
make stop

# Method 2: Kill manual
lsof -i :8081 -i :8082 -i :8083 -i :8084 -i :8085
kill -9 <PID>

# Method 3: Kill all Go processes (hati-hati!)
pkill go
```

---

## ❌ Error: inotifywait/fswatch Not Found

### Problem
```
inotifywait: command not found
```

### Solution Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install inotify-tools
make dev-local-watch
```

### Solution macOS
```bash
brew install fswatch
make dev-local-watch
```

### Or Use Air Instead
```bash
make install-air
make dev-local-hot
```

---

## 🔍 Comparison: Hot Reload Tools

| Tool | Install | Pros | Cons | Command |
|------|---------|------|------|---------|
| **Air** | `go install` | Fast, Go-specific | Needs Go 1.21+ compatible version | `make dev-local-hot` |
| **inotifywait** | `apt install` | Native Linux, stable | Linux only | `make dev-local-watch` |
| **fswatch** | `brew install` | Native macOS, stable | macOS only | `make dev-local-watch` |
| **Docker** | Built-in | No extra install, isolated | Slower startup | `make dev` |
| **None** | N/A | No dependencies | Manual restart needed | `make dev-local` |

---

## 📋 Quick Decision Tree

```
Want hot reload?
├─ Yes
│  ├─ Have Go 1.21?
│  │  ├─ Yes → make install-air && make dev-local-hot ✅
│  │  └─ No → Upgrade Go or use alternative
│  │
│  ├─ Linux → sudo apt install inotify-tools && make dev-local-watch
│  ├─ macOS → brew install fswatch && make dev-local-watch
│  └─ Any OS with Docker → make dev
│
└─ No → make dev-local (simplest!)
```

---

## 🎯 Recommended Setup by OS

### Linux (Ubuntu/Debian)
```bash
# Option A: Air (if Go 1.21+)
make install-air
make dev-local-hot

# Option B: inotifywait (native)
sudo apt install inotify-tools
make dev-local-watch
```

### macOS
```bash
# Option A: Air (if Go 1.21+)
make install-air
make dev-local-hot

# Option B: fswatch (native)
brew install fswatch
make dev-local-watch
```

### Windows
```bash
# Use Docker (easiest)
make dev

# Or WSL2 + Linux instructions
```

---

## 💡 Pro Tips

1. **Air version locked**: Project menggunakan Air v1.49.0 yang kompatibel dengan Go 1.21
2. **PATH important**: Pastikan `$(go env GOPATH)/bin` ada di PATH
3. **Docker fallback**: Jika semua gagal, gunakan `make dev` (Docker)
4. **No hot reload**: `make dev-local` tetap bisa digunakan kapan saja
5. **Check installation**: `which air` atau `air -v` untuk verifikasi

---

## 🔬 Debug Commands

```bash
# Check Go version
go version

# Check Air installation
which air
air -v
ls ~/go/bin/air

# Check GOPATH
go env GOPATH

# Check PATH
echo $PATH | grep go

# Test Air config
cd services/auth-service
air -v

# Check running processes
ps aux | grep go
lsof -i :8081-8085

# Check available ports
netstat -tuln | grep 808
```

---

## 📝 Summary

**Best practices:**
1. **First choice**: `make dev-local-hot` (if Air installed)
2. **Alternative**: `make dev-local-watch` (if inotifywait/fswatch available)
3. **Fallback**: `make dev-local` (always works, no dependencies)
4. **Production-like**: `make dev` (Docker with hot reload)

**Key point**: Semua mode tetap berjalan dengan baik, hot reload hanya untuk convenience!
