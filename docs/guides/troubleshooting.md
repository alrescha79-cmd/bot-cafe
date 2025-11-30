# Troubleshooting Guide

Solusi untuk masalah umum saat development dan deployment.

## 🔥 Hot Reload Issues

### Air Installation Failed

**Problem:**
```
go: github.com/cosmtrek/air@latest requires go >= 1.25
```

**Solution:**

Install versi Air yang kompatibel dengan Go 1.21:

```bash
# Install Air v1.49.0
make install-air
# atau
go install github.com/cosmtrek/air@v1.49.0

# Verify
air -v
```

### Air Command Not Found

**Problem:**
```
air: command not found
```

**Solution:**

```bash
# Check if Air is installed
which air

# Add to PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Make permanent (add to ~/.zshrc or ~/.bashrc)
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc

# If not installed
make install-air
```

### Hot Reload Not Detecting Changes

**Problem:** File changes tidak trigger rebuild

**Solutions:**

1. **Check Air is running:**
   ```bash
   # Should see output like:
   [AUTH] watching .
   [MENU] watching ../../tmp/../shared
   ```

2. **Verify shared directory watching:**
   - Fixed in latest version
   - Air now watches `shared/` directory
   - Changes to shared code trigger all service rebuilds

3. **Restart with clean state:**
   ```bash
   make stop
   rm -rf tmp/
   make dev-local-hot
   ```

### Alternative Hot Reload Methods

If Air doesn't work, try alternatives:

**Option 1: inotifywait (Linux)**
```bash
sudo apt install inotify-tools
make dev-local-watch
```

**Option 2: fswatch (macOS)**
```bash
brew install fswatch
make dev-local-watch
```

**Option 3: Docker hot reload**
```bash
make dev
```

**Option 4: No hot reload**
```bash
make dev-local
```

## 🔌 Port Issues

### Port Already in Use

**Problem:**
```
listen tcp :8081: bind: address already in use
```

**Solutions:**

```bash
# Method 1: Use make stop
make stop

# Method 2: Kill specific ports
lsof -i :8081 -i :8082 -i :8083 -i :8084 -i :8085
kill -9 <PID>

# Method 3: Kill all Go processes (caution!)
pkill go
pkill -f "auth-service|menu-service|promo-service"

# Method 4: Reboot (last resort)
sudo reboot
```

### Check Which Process Uses Port

```bash
# Linux/macOS
lsof -i :8081
netstat -tuln | grep 8081

# Find and kill
kill -9 $(lsof -t -i:8081)
```

## 🤖 Bot Issues

### Bot Not Responding

**Diagnosis:**

```bash
# Check agent logs
docker logs cafe-bot-agent
# or for local:
# Look for [AGENT] lines in terminal
```

**Common Causes & Solutions:**

| Cause | Solution |
|-------|----------|
| Invalid bot token | Check `.env`, verify token with @BotFather |
| Agent not running | `docker ps` atau check terminal output |
| Network issues | Check internet connection |
| Telegram API down | Wait and retry |

**Quick Fix:**
```bash
# Restart agent
docker restart cafe-bot-agent
# or for local:
make stop && make dev-local-hot
```

### Admin Access Denied

**Problem:** Admin commands tidak work, bot says "Access denied"

**Solutions:**

1. **Verify Telegram ID:**
   ```bash
   cat .vars.json
   # Should contain your Telegram ID
   ```

2. **Get your Telegram ID:**
   - Chat with @userinfobot
   - Copy the ID number
   - Add to `.vars.json`:
     ```json
     {
       "admin_telegram_ids": ["YOUR_ID_HERE"]
     }
     ```

3. **Restart service:**
   ```bash
   # Docker
   docker restart cafe-bot-agent
   
   # Local
   make stop
   make dev-local-hot
   ```

4. **Check logs for auth:**
   ```bash
   docker logs cafe-bot-agent | grep AUTH
   # Should show: "✅ MATCH! User XXXXX is admin"
   ```

### Bot Shows Old Data After Update

**Problem:** Changes to menu/promo not reflected

**Solutions:**

1. **Clear cache (if any):**
   ```bash
   # Restart bot
   docker restart cafe-bot-agent
   ```

2. **Check database:**
   ```bash
   ls -la data/
   # Should show updated timestamps
   ```

3. **Verify update was saved:**
   ```bash
   # Check logs for successful update
   docker logs cafe-menu-service | grep -i "update"
   ```

## 🗄️ Database Issues

### Database File Not Found

**Problem:**
```
Failed to initialize database: no such file or directory
```

**Solution:**

```bash
# Create data directory
mkdir -p data

# Restart services (will create databases)
make dev-local-hot
```

### Database Locked

**Problem:**
```
database is locked
```

**Solutions:**

```bash
# Stop all services
make stop

# Kill any processes holding the lock
lsof data/*.db
kill <PID>

# Restart
make dev-local-hot
```

### Reset Database

**⚠️ Warning: This deletes all data!**

```bash
# Stop services
make stop

# Delete databases
rm -rf data/*.db

# Restart (will create fresh databases)
make dev-local-hot
```

### Backup Before Reset

```bash
# Backup
tar -czf backup-$(date +%Y%m%d).tar.gz data/

# Then reset
rm -rf data/*.db
make dev-local-hot
```

## 🐳 Docker Issues

### Docker Daemon Not Running

**Problem:**
```
Cannot connect to the Docker daemon
```

**Solution:**

```bash
# Start Docker
sudo systemctl start docker

# Enable on boot
sudo systemctl enable docker

# Verify
docker ps
```

### Docker Compose Not Found

**Problem:**
```
docker-compose: command not found
```

**Solution:**

```bash
# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify
docker-compose --version
```

### Container Won't Start

**Diagnosis:**

```bash
# Check container status
docker ps -a

# Check logs
docker logs cafe-bot-agent

# Inspect container
docker inspect cafe-bot-agent
```

**Common Solutions:**

```bash
# Remove and recreate
docker-compose -f deployments/docker-compose.yml down
docker-compose -f deployments/docker-compose.yml up -d

# Rebuild images
docker-compose -f deployments/docker-compose.yml build --no-cache
docker-compose -f deployments/docker-compose.yml up -d
```

### Disk Space Full

**Problem:**
```
no space left on device
```

**Solution:**

```bash
# Check disk usage
df -h

# Clean Docker
docker system prune -a
docker volume prune

# Remove old images
docker image prune -a

# Remove old containers
docker container prune
```

## 💻 Development Issues

### Go Module Errors

**Problem:**
```
go: module not found
```

**Solution:**

```bash
# Download dependencies
make deps
# or
go mod download
go mod tidy

# Verify go.mod
cat go.mod
```

### Build Errors

**Problem:** Compilation errors

**Solutions:**

```bash
# Clean and rebuild
make clean
make build

# Check Go version
go version
# Should be 1.21+

# Update dependencies
go get -u ./...
go mod tidy
```

### Import Errors

**Problem:**
```
package X is not in GOROOT
```

**Solution:**

```bash
# Set Go environment
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Download missing packages
go get package-name
```

## 🌐 Network Issues

### Cannot Reach Telegram API

**Problem:** Bot can't connect to Telegram

**Diagnosis:**

```bash
# Test internet
ping google.com

# Test Telegram API
curl https://api.telegram.org

# Check DNS
nslookup api.telegram.org
```

**Solutions:**

1. **Check firewall:**
   ```bash
   sudo ufw status
   # Ensure outbound connections allowed
   ```

2. **Try different DNS:**
   ```bash
   # Google DNS
   echo "nameserver 8.8.8.8" | sudo tee /etc/resolv.conf
   ```

3. **Use proxy (if behind firewall):**
   - Configure HTTP_PROXY environment variable

### Service Can't Communicate

**Problem:** Agent can't reach menu-service, etc.

**Diagnosis:**

```bash
# Docker networking
docker network ls
docker network inspect bot-cafe_default

# Test service endpoints
curl http://localhost:8082/health
```

**Solution:**

```bash
# Restart Docker network
docker-compose -f deployments/docker-compose.yml down
docker-compose -f deployments/docker-compose.yml up -d
```

## 🔍 Debugging Tips

### Enable Verbose Logging

Modify code temporarily:

```go
// In any service main.go
func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    // ... rest of code
}
```

### Test API Endpoints Manually

```bash
# Test menu service
curl -X POST http://localhost:8082 \
  -H "Content-Type: application/json" \
  -d '{"action":"list","payload":{}}'

# Test auth service  
curl -X POST http://localhost:8081 \
  -H "Content-Type: application/json" \
  -d '{"action":"verify","payload":{"user_id":123456}}'
```

### Monitor Resource Usage

```bash
# CPU and Memory
top
htop

# Docker stats
docker stats

# Disk I/O
iotop
```

## 📞 Getting More Help

### Check Logs First

```bash
# Docker
docker-compose -f deployments/docker-compose.yml logs -f

# Local
# Look at terminal output with service prefixes
```

### Search Issues

1. Check existing issues on GitHub
2. Search error message on Google
3. Check Air documentation
4. Check Go Telegram Bot API docs

### Report an Issue

Include:
- OS and version
- Go version (`go version`)
- Docker version (`docker --version`)
- Complete error message
- Steps to reproduce
- Relevant logs

---

**Still Having Issues?**

Open an issue on GitHub with:
- Detailed description
- Error logs
- Steps taken
- System information

We're here to help! 🤝
