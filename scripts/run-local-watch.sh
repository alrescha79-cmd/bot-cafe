#!/bin/bash

# Script simple untuk hot reload tanpa dependencies eksternal
# Menggunakan inotifywait (Linux) atau fswatch (macOS)

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}🔥 Starting Bot Telegram Café with Hot Reload (Simple Mode)${NC}"

# Detect OS and check for file watcher
WATCHER=""
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    if command -v inotifywait &> /dev/null; then
        WATCHER="inotifywait"
    fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
    if command -v fswatch &> /dev/null; then
        WATCHER="fswatch"
    fi
fi

if [ -z "$WATCHER" ]; then
    echo -e "${YELLOW}No file watcher found. Install one:${NC}"
    echo "  - Linux: sudo apt install inotify-tools"
    echo "  - macOS: brew install fswatch"
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
    rm -f tmp/*.pid
    exit 0
}

trap cleanup SIGINT SIGTERM

# Function to run and watch service
watch_and_run() {
    local service=$1
    local port=$2
    local prefix=$3
    local dir="services/$service"
    
    while true; do
        # Start service
        (
            cd "$dir"
            ${port:+${port}=}go run . 2>&1 | sed "s/^/[$prefix] /"
        ) &
        local pid=$!
        echo $pid > "tmp/${service}.pid"
        
        # Watch for changes
        if [ "$WATCHER" = "inotifywait" ]; then
            inotifywait -q -r -e modify,create,delete --include '\.go$' "$dir"
        else
            fswatch -1 -r --include='\.go$' "$dir"
        fi
        
        # Restart on change
        echo -e "${YELLOW}[$prefix] File changed, reloading...${NC}"
        kill $pid 2>/dev/null || true
        sleep 1
    done
}

# Start services
echo -e "${BLUE}Starting services with file watcher ($WATCHER)...${NC}"

watch_and_run "auth-service" "AUTH_SERVICE_PORT=8081" "AUTH" &
sleep 0.5

watch_and_run "menu-service" "MENU_SERVICE_PORT=8082" "MENU" &
sleep 0.5

watch_and_run "promo-service" "PROMO_SERVICE_PORT=8083" "PROMO" &
sleep 0.5

watch_and_run "info-service" "INFO_SERVICE_PORT=8084" "INFO" &
sleep 0.5

watch_and_run "media-service" "MEDIA_SERVICE_PORT=8085" "MEDIA" &
sleep 1

# Start agent with watcher
echo -e "${GREEN}Starting Telegram bot agent...${NC}"
(
    while true; do
        (
            cd agent
            go run . 2>&1 | sed 's/^/[AGENT] /'
        ) &
        local pid=$!
        echo $pid > "tmp/agent.pid"
        
        if [ "$WATCHER" = "inotifywait" ]; then
            inotifywait -q -r -e modify,create,delete --include '\.go$' agent/
        else
            fswatch -1 -r --include='\.go$' agent/
        fi
        
        echo -e "${YELLOW}[AGENT] File changed, reloading...${NC}"
        kill $pid 2>/dev/null || true
        sleep 1
    done
) &

echo -e "${GREEN}✅ All services started with HOT RELOAD!${NC}"
echo -e "${BLUE}Edit any .go file and it will auto-reload!${NC}"
echo -e "${BLUE}Using file watcher: $WATCHER${NC}"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"

# Wait for all background processes
wait
