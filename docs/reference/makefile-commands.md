# Makefile Commands Reference

Complete reference untuk semua Makefile commands.

## 📦 Setup & Dependencies

### `make help`
Tampilkan semua available commands dengan deskripsi.

```bash
make help
```

### `make init` / `make setup`
Setup awal project (first-time setup).

```bash
make init
```

**What it does:**
- Copy `.env.example` → `.env`
- Copy `.vars.json.example` → `.vars.json`
- Create `data/` directory
- Download Go dependencies

**After running:**
- Edit `.env` - Set bot token
- Edit `.vars.json` - Set admin IDs

### `make deps`
Install/update Go dependencies.

```bash
make deps
```

**What it does:**
- `go mod download` - Download dependencies
- `go mod tidy` - Clean up go.mod

## 🚀 Development Commands

### `make dev-local-hot`
Run lokal dengan hot reload (recommended untuk development).

```bash
make dev-local-hot
```

**Features:**
- ✅ Hot reload via Air v1.49.0
- ✅ Auto rebuild on file changes
- ✅ Watches `shared/` directory
- ✅ Fast iteration

**Requirements:**
- Air v1.49.0 (auto-installed if missing)

**Stop:** `Ctrl+C` atau `make stop`

### `make dev-local` / `make run-local` / `make run` / `make dev`
Run lokal tanpa hot reload.

```bash
make dev-local
```

**Features:**
- ⚡ Fast startup
- ✅ No dependencies
- ✅ All services in one terminal

**Stop:** `Ctrl+C` atau `make stop`

### `make dev-local-watch`
Run lokal dengan inotifywait/fswatch hot reload.

```bash
make dev-local-watch
```

**Requirements:**
- Linux: `inotify-tools`
- macOS: `fswatch`

### `make stop`
Stop semua running services (local mode).

```bash
make stop
```

**What it does:**
- Kill all service processes
- Clean up background jobs

## 🐳 Docker Commands

### `make dev`
Start development environment dengan Docker hot reload.

```bash
make dev
```

**What it does:**
- Create `.vars.json` if not exists
- Start containers with `docker-compose up -d`
- Follow logs

**Equivalent to:**
```bash
docker-compose -f deployments/docker-compose.yml up -d
docker-compose -f deployments/docker-compose.yml logs -f
```

### `make docker-build`
Build Docker images.

```bash
make docker-build
```

**When to use:**
- After changing Dockerfile
- After adding new dependencies

### `make docker-up`
Start Docker containers (detached mode).

```bash
make docker-up
```

**What it does:**
- Create `.vars.json` if missing
- Start all containers in background

### `make docker-down`
Stop and remove Docker containers.

```bash
make docker-down
```

### `make docker-logs`
View logs dari all services.

```bash
make docker-logs
```

**Options:**
```bash
# Follow logs (continuous)
make docker-logs

# Last 100 lines
docker-compose -f deployments/docker-compose.yml logs --tail 100

# Specific service
docker logs cafe-bot-agent
```

### `make docker-restart`
Restart Docker containers.

```bash
make docker-restart
```

## 🔨 Build Commands

### `make build`
Build semua services (binaries).

```bash
make build
```

**Output:**
```
bin/
├── auth-service
├── menu-service
├── promo-service
├── info-service
├── media-service
└── agent
```

**What it does:**
- Build each service dari source
- Output ke `bin/` directory

### `make clean`
Bersihkan build artifacts dan databases.

```bash
make clean
```

**⚠️ Warning:** Menghapus semua data!

**What it deletes:**
- `bin/` - All binaries
- `data/` - All databases
- `tmp/` - Temporary files
- All `*.db` files

## 🧪 Testing Commands

### `make test`
Run all tests.

```bash
make test
```

**Equivalent to:**
```bash
go test ./...
```


## 🛠️ Utility Commands

### `make install-air`
Install Air v1.49.0 untuk hot reload.

```bash
make install-air
```

**What it does:**
- Install `github.com/cosmtrek/air@v1.49.0`
- Add to `$GOPATH/bin`

**Verify:**
```bash
air -v
# Output: v1.49.0
```

## 📊 Command Cheat Sheet

### Daily Development

```bash
# Start development
make dev-local-hot

# Edit code → Auto reload!

# Stop
Ctrl+C
```

### First Time Setup

```bash
# 1. Initialize
make init

# 2. Configure
nano .env
nano .vars.json

# 3. Run
make dev-local-hot
```

### With Docker

```bash
# Start
make dev

# View logs
make docker-logs

# Restart
make docker-restart

# Stop
make docker-down
```

### Building

```bash
# Build all
make build

# Run manually
./bin/agent
```

### Clean Slate

```bash
# Stop everything
make stop
make docker-down

# Clean all
make clean

# Fresh start
make init
make dev-local-hot
```

## 🔗 Command Dependencies

```
make init
  ├─ cp .env.example .env
  ├─ cp .vars.json.example .vars.json
  ├─ mkdir data
  └─ go mod download

make build
  ├─ cd services/*/
  └─ go build

make dev-local-hot
  ├─ Check air installed
  ├─ Create air configs
  └─ Run with air

make dev
  ├─ Create .vars.json (if missing)
  └─ docker-compose up -d

make clean
  ├─ rm -rf bin/
  ├─ rm -rf data/
  ├─ rm -rf tmp/
  └─ find -delete *.db
```

## 💡 Tips

### Combine Commands

```bash
# Clean then start fresh
make clean && make dev-local-hot

# Build then run
make build && ./bin/agent

# Stop Docker then start local
make docker-down && make dev-local-hot
```

### Check Before Running

```bash
# Check ports are free
lsof -i :8081-8085

# Check Docker status
docker ps

# Check Go version
go version
```

### Logs & Debugging

```bash
# Follow logs
make docker-logs

# Specific service
docker logs -f cafe-bot-agent

# Local logs
# Watch terminal output with service prefixes
```

---

**Pro Tip:** Use `make help` to quickly see all available commands!

For more details, check the [Development Setup Guide](../guides/development-setup.md).
