#!/bin/bash

# Script untuk menjalankan semua services secara lokal (tanpa Docker)

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Starting Bot Telegram Café (Local Mode)${NC}"

# Add Go bin to PATH if not already there
export PATH="$PATH:$(go env GOPATH)/bin"

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${RED}Error: .env file not found!${NC}"
    echo "Run: cp .env.example .env"
    exit 1
fi

# Check if .vars.json exists
if [ ! -f .vars.json ]; then
    echo -e "${RED}Error: .vars.json file not found!${NC}"
    echo "Run: cp .vars.json.example .vars.json"
    exit 1
fi

# Create data directory if not exists
mkdir -p data

# Load environment variables
export $(grep -v '^#' .env | xargs)

# Function to cleanup on exit
cleanup() {
    echo -e "\n${YELLOW}Stopping all services...${NC}"
    pkill -P $$ || true
    exit 0
}

trap cleanup SIGINT SIGTERM

# Start services in background
echo -e "${GREEN}Starting auth-service on port 8081...${NC}"
(cd services/auth-service && AUTH_SERVICE_PORT=8081 go run . 2>&1 | sed 's/^/[AUTH] /') &
AUTH_PID=$!

echo -e "${GREEN}Starting menu-service on port 8082...${NC}"
(cd services/menu-service && MENU_SERVICE_PORT=8082 go run . 2>&1 | sed 's/^/[MENU] /') &
MENU_PID=$!

echo -e "${GREEN}Starting promo-service on port 8083...${NC}"
(cd services/promo-service && PROMO_SERVICE_PORT=8083 go run . 2>&1 | sed 's/^/[PROMO] /') &
PROMO_PID=$!

echo -e "${GREEN}Starting info-service on port 8084...${NC}"
(cd services/info-service && INFO_SERVICE_PORT=8084 go run . 2>&1 | sed 's/^/[INFO] /') &
INFO_PID=$!

echo -e "${GREEN}Starting media-service on port 8085...${NC}"
(cd services/media-service && MEDIA_SERVICE_PORT=8085 go run . 2>&1 | sed 's/^/[MEDIA] /') &
MEDIA_PID=$!

# Wait for services to be ready
echo -e "${YELLOW}Waiting for services to be ready...${NC}"
sleep 3

# Start agent
echo -e "${GREEN}Starting Telegram bot agent...${NC}"
(cd agent && go run . 2>&1 | sed 's/^/[AGENT] /') &
AGENT_PID=$!

echo -e "${GREEN}✅ All services started!${NC}"
echo -e "${GREEN}Process IDs:${NC}"
echo "  AUTH:  $AUTH_PID"
echo "  MENU:  $MENU_PID"
echo "  PROMO: $PROMO_PID"
echo "  INFO:  $INFO_PID"
echo "  MEDIA: $MEDIA_PID"
echo "  AGENT: $AGENT_PID"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"

# Wait for all background processes
wait
