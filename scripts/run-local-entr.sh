#!/bin/bash

# Script untuk menjalankan dengan hot reload menggunakan entr (lebih simple)

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}🔥 Starting Bot Telegram Café with Hot Reload (entr mode)${NC}"

# Check if entr is installed
if ! command -v entr &> /dev/null; then
    echo -e "${YELLOW}entr not found. Please install it:${NC}"
    echo "  - Ubuntu/Debian: sudo apt install entr"
    echo "  - macOS: brew install entr"
    echo "  - Arch: sudo pacman -S entr"
    echo ""
    echo -e "${YELLOW}Or use without hot reload: make dev-local${NC}"
    exit 1
fi

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

# Create data directory
mkdir -p data tmp

# Load environment variables
export $(grep -v '^#' .env | xargs)

# Function to cleanup on exit
cleanup() {
    echo -e "\n${YELLOW}Stopping all services...${NC}"
    pkill -P $$ || true
    exit 0
}

trap cleanup SIGINT SIGTERM

# Function to run service with entr
run_with_entr() {
    local service=$1
    local port=$2
    local prefix=$3
    
    while true; do
        echo -e "${GREEN}Starting $service on port $port...${NC}"
        (
            cd "services/$service"
            find . -name "*.go" | entr -r sh -c "PORT=$port go run . 2>&1 | sed 's/^/[$prefix] /'"
        ) &
        wait $!
        echo -e "${YELLOW}$service crashed, restarting...${NC}"
        sleep 1
    done
}

# Start services with entr
echo -e "${BLUE}Starting services with hot reload...${NC}"

run_with_entr "auth-service" "8081" "AUTH" &
sleep 1

run_with_entr "menu-service" "8082" "MENU" &
sleep 1

run_with_entr "promo-service" "8083" "PROMO" &
sleep 1

run_with_entr "info-service" "8084" "INFO" &
sleep 1

run_with_entr "media-service" "8085" "MEDIA" &
sleep 2

# Start agent with entr
echo -e "${GREEN}Starting Telegram bot agent...${NC}"
(
    cd agent
    while true; do
        find . -name "*.go" | entr -r sh -c "go run . 2>&1 | sed 's/^/[AGENT] /'"
        echo -e "${YELLOW}Agent crashed, restarting...${NC}"
        sleep 1
    done
) &

echo -e "${GREEN}✅ All services started with HOT RELOAD (entr)!${NC}"
echo -e "${BLUE}Edit any .go file and it will auto-reload!${NC}"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"

# Wait for all background processes
wait
