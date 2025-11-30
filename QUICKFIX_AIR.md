# ⚡ QUICK FIX - Air Installation Error

## Problem
```
make dev-local-hot
Error: go: github.com/cosmtrek/air@latest requires go >= 1.25
```

## ✅ Solution (Copy-Paste Ini)

```bash
# 1. Install Air versi yang kompatibel
go install github.com/cosmtrek/air@v1.49.0

# 2. Tambahkan ke PATH (pilih salah satu)
# Untuk session ini saja:
export PATH="$PATH:$(go env GOPATH)/bin"

# Atau permanent (tambahkan ke .bashrc/.zshrc):
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc

# 3. Verify installation
air -v

# 4. Jalankan bot
make dev-local-hot
```

## 🚀 Alternative Quick Solutions

### Tanpa Air (Paling Mudah)
```bash
make dev-local
# Semua jalan, tapi tanpa hot reload
```

### Dengan Docker (Production-like)
```bash
make dev
# Hot reload sudah built-in di Docker
```

### Dengan inotifywait (Linux)
```bash
sudo apt install inotify-tools
make dev-local-watch
```

### Dengan fswatch (macOS)
```bash
brew install fswatch
make dev-local-watch
```

---

## 📋 Makefile Commands Summary

```bash
make dev-local          # ⚡ Tercepat - tanpa hot reload
make dev-local-hot      # 🔥 Dengan hot reload (Air)
make dev-local-watch    # 👀 Dengan hot reload (inotifywait/fswatch)
make dev                # 🐳 Docker dengan hot reload
make install-air        # 📦 Install Air v1.49.0
make stop               # 🛑 Stop semua services
```

---

## 💡 Which One to Use?

- **Daily development**: `make dev-local-hot` (setelah install Air)
- **Quick test**: `make dev-local`
- **No Air installed**: `make dev-local` atau `make dev`
- **Production test**: `make dev`

---

**Done! Choose one and start coding! 🎉**
