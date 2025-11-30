# 🚀 Deployment Guide

## Quick Start (5 menit)

### 1. Prerequisites

```bash
# Cek Docker
docker --version
docker-compose --version

# Cek Go (optional, untuk dev tanpa Docker)
go version
```

### 2. Get Bot Token

1. Buka Telegram, cari [@BotFather](https://t.me/botfather)
2. Kirim `/newbot`
3. Ikuti instruksi
4. Copy token yang diberikan

### 3. Setup Project

```bash
# Clone repository
git clone <your-repo-url>
cd bot-cafe

# Initialize
make init

# Edit .env - masukkan bot token Anda
nano .env
# Set: TELEGRAM_BOT_TOKEN=your_token_here

# Edit .vars.json - masukkan Telegram ID Anda
nano .vars.json
# Set admin_telegram_ids dengan ID Anda
```

**Cara mendapatkan Telegram ID:**
1. Cari bot [@userinfobot](https://t.me/userinfobot)
2. Kirim `/start`
3. Copy ID yang diberikan

### 4. Run!

```bash
# Start dengan hot reload
make dev

# Atau step by step
make docker-up
make docker-logs
```

### 5. Test Bot

1. Buka Telegram
2. Cari bot Anda (nama yang Anda set di BotFather)
3. Kirim `/start`
4. Enjoy! 🎉

---

## Development Workflow

### Daily Development

```bash
# Start services
make dev

# Edit code (auto reload!)
# services/menu-service/handlers.go
# agent/handlers.go

# Check logs
make docker-logs

# Restart if needed
make docker-restart

# Stop when done
make docker-down
```

### Testing Changes

```bash
# Test specific service
curl -X POST http://localhost:8082 \
  -H "Content-Type: application/json" \
  -d '{"action":"list","payload":{}}'

# Check health
curl http://localhost:8082/health

# View database
docker exec -it cafe-menu-service sh
sqlite3 /data/menu.db
> SELECT * FROM menus;
> .quit
```

### Debugging

```bash
# View logs for specific service
docker logs -f cafe-bot-agent
docker logs -f cafe-menu-service

# Enter container
docker exec -it cafe-bot-agent sh

# Check environment
docker exec cafe-bot-agent env

# Rebuild specific service
docker-compose -f deployments/docker-compose.yml up -d --build menu-service
```

---

## Production Deployment

### Option 1: Docker Compose (Simple)

**Untuk VPS atau server sederhana**

```bash
# 1. Setup server
ssh user@your-server

# 2. Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# 3. Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 4. Clone repository
git clone <your-repo-url>
cd bot-cafe

# 5. Setup
make init
nano .env          # Set bot token
nano .vars.json    # Set admin IDs

# 6. Run
docker-compose -f deployments/docker-compose.yml up -d

# 7. Monitor
docker-compose -f deployments/docker-compose.yml logs -f

# 8. Auto-restart on reboot
# Create systemd service
sudo nano /etc/systemd/system/bot-cafe.service
```

**systemd service file:**
```ini
[Unit]
Description=Bot Cafe Telegram
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/user/bot-cafe
ExecStart=/usr/local/bin/docker-compose -f deployments/docker-compose.yml up -d
ExecStop=/usr/local/bin/docker-compose -f deployments/docker-compose.yml down
User=user

[Install]
WantedBy=multi-user.target
```

```bash
# Enable service
sudo systemctl enable bot-cafe
sudo systemctl start bot-cafe

# Check status
sudo systemctl status bot-cafe
```

---

### Option 2: Kubernetes (Scalable)

**Untuk production dengan high availability**

#### Prerequisites
```bash
# kubectl installed
kubectl version

# Cluster ready (GKE, EKS, AKS, atau self-hosted)
kubectl get nodes
```

#### 1. Create Namespace
```bash
kubectl create namespace bot-cafe
```

#### 2. Create Secrets
```bash
# Bot token
kubectl create secret generic bot-secrets \
  --from-literal=telegram-bot-token=YOUR_TOKEN \
  -n bot-cafe

# Admin config
kubectl create configmap admin-config \
  --from-file=.vars.json \
  -n bot-cafe
```

#### 3. Create Persistent Volumes
```yaml
# k8s/pv.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: auth-data
  namespace: bot-cafe
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
---
# Repeat for menu-data, promo-data, info-data, media-data
```

```bash
kubectl apply -f k8s/pv.yaml
```

#### 4. Deploy Services
```yaml
# k8s/menu-service.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: menu-service
  namespace: bot-cafe
spec:
  replicas: 2
  selector:
    matchLabels:
      app: menu-service
  template:
    metadata:
      labels:
        app: menu-service
    spec:
      containers:
      - name: menu-service
        image: your-registry/menu-service:latest
        ports:
        - containerPort: 8082
        env:
        - name: MENU_SERVICE_PORT
          value: "8082"
        - name: MENU_DB_PATH
          value: "/data/menu.db"
        volumeMounts:
        - name: data
          mountPath: /data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: menu-data
---
apiVersion: v1
kind: Service
metadata:
  name: menu-service
  namespace: bot-cafe
spec:
  selector:
    app: menu-service
  ports:
  - port: 8082
    targetPort: 8082
```

```bash
# Deploy all services
kubectl apply -f k8s/auth-service.yaml
kubectl apply -f k8s/menu-service.yaml
kubectl apply -f k8s/promo-service.yaml
kubectl apply -f k8s/info-service.yaml
kubectl apply -f k8s/media-service.yaml
kubectl apply -f k8s/agent.yaml

# Check status
kubectl get pods -n bot-cafe
kubectl get services -n bot-cafe
```

#### 5. Monitor
```bash
# Logs
kubectl logs -f deployment/agent -n bot-cafe

# Scale
kubectl scale deployment menu-service --replicas=3 -n bot-cafe

# Update
kubectl set image deployment/menu-service menu-service=your-registry/menu-service:v2 -n bot-cafe
```

---

## CI/CD Pipeline

### GitHub Actions

**.github/workflows/deploy.yml**
```yaml
name: Build and Deploy

on:
  push:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    - name: Run tests
      run: go test ./...

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Build Docker images
      run: |
        docker build -t your-registry/auth-service:${{ github.sha }} -f deployments/Dockerfile.service --build-arg SERVICE_NAME=auth-service .
        docker build -t your-registry/menu-service:${{ github.sha }} -f deployments/Dockerfile.service --build-arg SERVICE_NAME=menu-service .
        # ... other services
        docker build -t your-registry/agent:${{ github.sha }} -f deployments/Dockerfile.agent .
    
    - name: Push to registry
      run: |
        echo ${{ secrets.DOCKER_PASSWORD }} | docker login -u ${{ secrets.DOCKER_USERNAME }} --password-stdin
        docker push your-registry/auth-service:${{ github.sha }}
        docker push your-registry/menu-service:${{ github.sha }}
        # ... other services
        docker push your-registry/agent:${{ github.sha }}

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
    - name: Deploy to server
      uses: appleboy/ssh-action@master
      with:
        host: ${{ secrets.SERVER_HOST }}
        username: ${{ secrets.SERVER_USER }}
        key: ${{ secrets.SSH_KEY }}
        script: |
          cd /home/user/bot-cafe
          git pull
          docker-compose -f deployments/docker-compose.yml pull
          docker-compose -f deployments/docker-compose.yml up -d
```

---

## Monitoring

### Prometheus + Grafana

```yaml
# docker-compose.monitoring.yml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - cafe-network

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    networks:
      - cafe-network

networks:
  cafe-network:
    external: true
```

**prometheus.yml:**
```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'bot-services'
    static_configs:
      - targets:
        - 'auth-service:8081'
        - 'menu-service:8082'
        - 'promo-service:8083'
        - 'info-service:8084'
        - 'media-service:8085'
```

```bash
# Start monitoring
docker-compose -f docker-compose.monitoring.yml up -d

# Access Grafana
# http://localhost:3000
# Default: admin/admin
```

---

## Backup & Recovery

### Backup Database

```bash
# Backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/bot-cafe/$DATE"

mkdir -p $BACKUP_DIR

docker cp cafe-auth-service:/data/auth.db $BACKUP_DIR/
docker cp cafe-menu-service:/data/menu.db $BACKUP_DIR/
docker cp cafe-promo-service:/data/promo.db $BACKUP_DIR/
docker cp cafe-info-service:/data/info.db $BACKUP_DIR/
docker cp cafe-media-service:/data/media.db $BACKUP_DIR/

tar -czf $BACKUP_DIR.tar.gz $BACKUP_DIR
rm -rf $BACKUP_DIR

echo "Backup completed: $BACKUP_DIR.tar.gz"
```

### Restore Database

```bash
#!/bin/bash
BACKUP_FILE=$1

if [ -z "$BACKUP_FILE" ]; then
  echo "Usage: ./restore.sh backup_file.tar.gz"
  exit 1
fi

tar -xzf $BACKUP_FILE
BACKUP_DIR=$(basename $BACKUP_FILE .tar.gz)

docker cp $BACKUP_DIR/auth.db cafe-auth-service:/data/
docker cp $BACKUP_DIR/menu.db cafe-menu-service:/data/
docker cp $BACKUP_DIR/promo.db cafe-promo-service:/data/
docker cp $BACKUP_DIR/info.db cafe-info-service:/data/
docker cp $BACKUP_DIR/media.db cafe-media-service:/data/

docker-compose -f deployments/docker-compose.yml restart

echo "Restore completed"
```

### Automated Backup (Cron)

```bash
# Add to crontab
crontab -e

# Daily backup at 2 AM
0 2 * * * /home/user/bot-cafe/scripts/backup.sh

# Keep backups for 30 days
0 3 * * * find /backup/bot-cafe -type f -mtime +30 -delete
```

---

## Security Checklist

- [ ] Bot token disimpan di environment variable (bukan hardcode)
- [ ] `.vars.json` tidak di-commit ke git
- [ ] File `.env` tidak di-commit ke git
- [ ] Telegram ID admin sudah benar
- [ ] Database volumes di-backup secara berkala
- [ ] HTTPS untuk production (jika ada webhook)
- [ ] Firewall dikonfigurasi (hanya expose port yang perlu)
- [ ] Docker images dari source terpercaya
- [ ] Update dependencies secara berkala
- [ ] Monitor logs untuk aktivitas mencurigakan

---

## Troubleshooting

### Bot tidak merespons
```bash
# Check logs
docker logs cafe-bot-agent

# Check token
docker exec cafe-bot-agent env | grep TELEGRAM

# Test connectivity
docker exec cafe-bot-agent ping api.telegram.org
```

### Service tidak bisa connect
```bash
# Check network
docker network inspect bot-cafe_cafe-network

# Check service status
docker-compose -f deployments/docker-compose.yml ps

# Restart services
docker-compose -f deployments/docker-compose.yml restart
```

### Database error
```bash
# Check permissions
docker exec cafe-menu-service ls -la /data

# Check database
docker exec cafe-menu-service sqlite3 /data/menu.db ".tables"

# Recreate database
docker volume rm bot-cafe_menu-data
docker-compose -f deployments/docker-compose.yml up -d menu-service
```

### Out of memory
```bash
# Check resources
docker stats

# Increase limits in docker-compose.yml
services:
  menu-service:
    deploy:
      resources:
        limits:
          memory: 512M
```

---

## Performance Optimization

1. **Database Indexing**: Sudah ada di schema
2. **Connection Pooling**: Default Go HTTP client
3. **Caching**: Add Redis untuk frequent queries
4. **Load Balancing**: Kubernetes service load balancer
5. **CDN**: Untuk media files
6. **Database Optimization**: VACUUM, ANALYZE
7. **Monitoring**: Track slow queries

---

## Scaling Guide

### Horizontal Scaling (Kubernetes)
```bash
# Scale services
kubectl scale deployment menu-service --replicas=5 -n bot-cafe
kubectl scale deployment promo-service --replicas=3 -n bot-cafe

# Auto-scaling
kubectl autoscale deployment menu-service --min=2 --max=10 --cpu-percent=70 -n bot-cafe
```

### Vertical Scaling (Docker Compose)
```yaml
services:
  menu-service:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 1G
        reservations:
          cpus: '1'
          memory: 512M
```

---

Selamat! Bot Anda siap production! 🚀

Untuk pertanyaan lebih lanjut, buka issue di GitHub atau hubungi maintainer.
