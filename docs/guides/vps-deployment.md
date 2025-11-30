# VPS Deployment Guide

Complete guide untuk deploy bot café ke VPS production.

## Prerequisites

### VPS Requirements
- **OS**: Ubuntu 20.04+ / Debian 11+
- **RAM**: Minimum 1GB (Recommended 2GB)
- **Storage**: Minimum 10GB
- **Network**: Public IP dengan port 80/443 terbuka (optional untuk webhook)

### Required Software
- Docker & Docker Compose
- Git
- Telegram Bot Token

## 🚀 Quick Deployment

### 1. Prepare VPS

```bash
# Connect ke VPS
ssh root@your-vps-ip

# Update system
apt update && apt upgrade -y

# Install required packages
apt install -y git curl
```

### 2. Install Docker

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Verify installation
docker --version
docker-compose --version

# Start Docker
systemctl start docker
systemctl enable docker
```

### 3. Deploy Application

```bash
# Clone repository
cd /opt
git clone <your-repo-url> bot-cafe
cd bot-cafe

# Setup environment
cp .env.example .env
cp .vars.json.example .vars.json

# Edit configuration
nano .env
# Set: TELEGRAM_BOT_TOKEN=your_token

nano .vars.json
# Set admin Telegram IDs
```

### 4. Start Services

```bash
# Build and start
docker-compose -f deployments/docker-compose.yml up -d

# Check status
docker-compose -f deployments/docker-compose.yml ps

# Check logs
docker-compose -f deployments/docker-compose.yml logs -f
```

✅ Bot sekarang running di VPS!

## 🔧 Production Configuration

### Environment Variables

Edit `.env` untuk production:

```env
# Required
TELEGRAM_BOT_TOKEN=your_bot_token_here

# Service Ports (internal)
AUTH_SERVICE_PORT=8081
MENU_SERVICE_PORT=8082
PROMO_SERVICE_PORT=8083
INFO_SERVICE_PORT=8084
MEDIA_SERVICE_PORT=8085

# Database Paths
AUTH_DB_PATH=./data/auth.db
MENU_DB_PATH=./data/menu.db
PROMO_DB_PATH=./data/promo.db
INFO_DB_PATH=./data/info.db
MEDIA_DB_PATH=./data/media.db
```

### Admin Configuration

Edit `.vars.json`:

```json
{
  "admin_telegram_ids": [
    "123456789",
    "987654321"
  ],
  "admin_usernames": [
    "admin1",
    "admin2"
  ]
}
```

> ⚠️ **Security**: Jangan commit file ini! Sudah ada di `.gitignore`.

## 🔄 Running as System Service

### Create Systemd Service

```bash
# Create service file
sudo nano /etc/systemd/system/bot-cafe.service
```

Add configuration:

```ini
[Unit]
Description=Bot Cafe Telegram
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/bot-cafe
ExecStart=/usr/local/bin/docker-compose -f deployments/docker-compose.yml up -d
ExecStop=/usr/local/bin/docker-compose -f deployments/docker-compose.yml down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
```

### Enable and Start Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service (auto-start on boot)
sudo systemctl enable bot-cafe

# Start service
sudo systemctl start bot-cafe

# Check status
sudo systemctl status bot-cafe
```

### Service Commands

```bash
# Start
sudo systemctl start bot-cafe

# Stop
sudo systemctl stop bot-cafe

# Restart
sudo systemctl restart bot-cafe

# Status
sudo systemctl status bot-cafe

# Logs
journalctl -u bot-cafe -f
```

## 📊 Monitoring & Logs

### View Logs

```bash
# All services
cd /opt/bot-cafe
docker-compose -f deployments/docker-compose.yml logs -f

# Specific service
docker logs -f cafe-bot-agent
docker logs -f cafe-menu-service

# Last 100 lines
docker logs --tail 100 cafe-bot-agent
```

### Check Resource Usage

```bash
# Docker stats
docker stats

# Disk usage
df -h
du -sh /opt/bot-cafe/data

# Memory usage
free -h

# CPU load
top
```

### Health Checks

```bash
# Check containers
docker ps

# Test health endpoints (from inside VPS)
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
curl http://localhost:8084/health
curl http://localhost:8085/health

#Bot status
docker logs cafe-bot-agent | tail -20
```

## 🔄 Updates & Maintenance

### Update Application

```bash
cd /opt/bot-cafe

# Pull latest changes
git pull origin main

# Rebuild containers
docker-compose -f deployments/docker-compose.yml build

# Restart services
docker-compose -f deployments/docker-compose.yml up -d

# Or with systemd
sudo systemctl restart bot-cafe
```

### Backup Database

```bash
# Create backup directory
mkdir -p /opt/bot-cafe-backups

# Backup all databases
cd /opt/bot-cafe
tar -czf /opt/bot-cafe-backups/data-$(date +%Y%m%d-%H%M%S).tar.gz data/

# Keep only last 7 days backups
find /opt/bot-cafe-backups -name "data-*.tar.gz" -mtime +7 -delete
```

### Automated Backup (Cron)

```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * cd /opt/bot-cafe && tar -czf /opt/bot-cafe-backups/data-$(date +\%Y\%m\%d-\%H\%M\%S).tar.gz data/ && find /opt/bot-cafe-backups -name "data-*.tar.gz" -mtime +7 -delete
```

### Restore from Backup

```bash
cd /opt/bot-cafe

# Stop services
docker-compose -f deployments/docker-compose.yml down

# Restore
tar -xzf /opt/bot-cafe-backups/data-20250130-020000.tar.gz

# Start services
docker-compose -f deployments/docker-compose.yml up -d
```

## 🔐 Security Best Practices

### 1. Firewall Configuration

```bash
# Install UFW
apt install -y ufw

# Allow SSH (IMPORTANT - don't lock yourself out!)
ufw allow 22/tcp

# Enable firewall
ufw enable

# Check status
ufw status
```

> Bot tidak perlu port terbuka karena menggunakan long polling, bukan webhook.

### 2. Secure SSH

Edit `/etc/ssh/sshd_config`:

```bash
# Disable root login
PermitRootLogin no

# Disable password authentication (use SSH keys)
PasswordAuthentication no

# Change default port (optional)
Port 2222
```

Restart SSH:
```bash
systemctl restart sshd
```

### 3. Regular Updates

```bash
# Auto-update security patches
apt install -y unattended-upgrades
dpkg-reconfigure --priority=low unattended-upgrades
```

### 4. File Permissions

```bash
cd /opt/bot-cafe

# Secure configuration files
chmod 600 .env .vars.json

# Database directory
chmod 755 data/
chmod 644 data/*.db
```

## 🐛 Troubleshooting

### Bot Not Starting

```bash
# Check logs
docker logs cafe-bot-agent

# Common issues:
# 1. Invalid bot token
cat .env | grep TELEGRAM_BOT_TOKEN

# 2. Services not running
docker ps

# 3. Port conflicts
netstat -tuln | grep 808

# Restart everything
docker-compose -f deployments/docker-compose.yml restart
```

### Database Issues

```bash
# Check database files
ls -lah data/

# Permissions
chmod 644 data/*.db

# Reset databases (WARNING: deletes all data!)
rm -rf data/*.db
docker-compose -f deployments/docker-compose.yml restart
```

### Memory Issues

```bash
# Check memory
free -h

# Restart services to free memory
docker-compose -f deployments/docker-compose.yml restart

# Or increase VPS RAM
```

### Disk Full

```bash
# Check disk space
df -h

# Clean Docker
docker system prune -a

# Clean old logs
docker logs cafe-bot-agent --tail 1000 > /tmp/bot.log
# Then check /tmp/bot.log and decide what to clean
```

## 📈 Scaling

### Horizontal Scaling (Multiple Instances)

Untuk handle traffic tinggi, scale services:

```bash
# Scale specific service
docker-compose -f deployments/docker-compose.yml up -d --scale menu-service=3

# Check running instances
docker ps | grep menu-service
```

### Load Balancer

Untuk production besar, tambahkan nginx:

```nginx
upstream menu_service {
    server localhost:8082;
    server localhost:8083;
    server localhost:8084;
}

server {
    listen 80;
    location / {
        proxy_pass http://menu_service;
    }
}
```

## 🚀 Advanced: Multiple Environments

### Production + Staging

```bash
# Production
/opt/bot-cafe-prod

# Staging
/opt/bot-cafe-staging
```

Gunakan different `.env` untuk masing-masing:
- Production: live bot token
- Staging: test bot token

## 🎯 Monitoring dengan Webhook (Optional)

Jika ingin gunakan webhook instead of long polling:

### 1. Setup SSL

```bash
# Install certbot
apt install -y certbot

# Get SSL certificate
certbot certonly --standalone -d yourdomain.com
```

### 2. Configure Nginx

```nginx
server {
    listen 443 ssl;
    server_name yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    location /webhook {
        proxy_pass http://localhost:8080;
    }
}
```

### 3. Update Bot Code

Modify agent untuk webhook mode (not covered here, requires code changes).

## 📝 Checklist Deployment

- [ ] VPS setup (Ubuntu, 2GB RAM)
- [ ] Docker installed
- [ ] Repository cloned
- [ ] `.env` configured dengan bot token
- [ ] `.vars.json` configured dengan admin IDs
- [ ] Services running (`docker ps`)
- [ ] Bot responding di Telegram
- [ ] Systemd service enabled
- [ ] Backup cron configured
- [ ] Firewall configured
- [ ] Monitoring setup

## 🆘 Getting Help

### Check Logs First

```bash
# Bot logs
docker logs cafe-bot-agent --tail 50

# All services
docker-compose -f deployments/docker-compose.yml logs --tail 50
```

### Common Commands

```bash
# Restart everything
sudo systemctl restart bot-cafe

# View running containers
docker ps

# Exec into container
docker exec -it cafe-bot-agent sh

# Clean restart
docker-compose -f deployments/docker-compose.yml down
docker-compose -f deployments/docker-compose.yml up -d
```

---

**Deployment Success!** 🎉

Bot sekarang running di production. Monitor regularly dan backup data secara teratur.

Questions? Check [Troubleshooting Guide](troubleshooting.md) atau open GitHub issue.
