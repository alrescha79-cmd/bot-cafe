#!/bin/bash

# Script untuk menjalankan semua services dengan hot reload menggunakan Air (tanpa Docker)

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}🔥 Starting Bot Telegram Café with Hot Reload (Local Mode)${NC}"

# Add Go bin to PATH if not already there
export PATH="$PATH:$(go env GOPATH)/bin"

# Check if air is installed
if ! command -v air &> /dev/null; then
    echo -e "${YELLOW}Air not found. Installing air v1.49.0 (compatible with Go 1.21)...${NC}"
    if go install github.com/cosmtrek/air@v1.49.0; then
        echo -e "${GREEN}Air installed successfully!${NC}"
    else
        echo -e "${RED}Failed to install Air.${NC}"
        echo -e "${YELLOW}Falling back to run without hot reload...${NC}"
        exec "$(dirname "$0")/run-local.sh"
        exit $?
    fi
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

# Create data and tmp directories with subdirectories for each service
echo -e "${BLUE}Creating necessary directories...${NC}"
mkdir -p data
mkdir -p tmp/auth-service
mkdir -p tmp/menu-service
mkdir -p tmp/promo-service
mkdir -p tmp/info-service
mkdir -p tmp/media-service
mkdir -p tmp/agent

# Load environment variables
export $(grep -v '^#' .env | xargs)

# Create .air.toml for each service if not exists
create_air_config() {
    local service=$1
    local port=$2
    local dir=$3
    
    # Determine relative path based on directory depth
    if [[ "$dir" == "agent" ]]; then
        local tmp_prefix="../tmp"
    else
        local tmp_prefix="../../tmp"
    fi
    
    cat > "$dir/.air.toml" << EOF
root = "."
testdata_dir = "testdata"
tmp_dir = "$tmp_prefix/$service"

[build]
  args_bin = []
  bin = "$tmp_prefix/$service/main"
  cmd = "go build -o $tmp_prefix/$service/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "$tmp_prefix/$service/build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
EOF
}

# Create air configs
echo -e "${BLUE}Creating Air configurations...${NC}"
create_air_config "auth-service" "8081" "services/auth-service"
create_air_config "menu-service" "8082" "services/menu-service"
create_air_config "promo-service" "8083" "services/promo-service"
create_air_config "info-service" "8084" "services/info-service"
create_air_config "media-service" "8085" "services/media-service"
create_air_config "agent" "" "agent"

# Function to cleanup on exit
cleanup() {
    echo -e "\n${YELLOW}Stopping all services...${NC}"
    pkill -P $$ || true
    rm -f services/*/.air.toml agent/.air.toml
    exit 0
}

trap cleanup SIGINT SIGTERM

# Start services with Air in background
echo -e "${GREEN}Starting auth-service on port 8081 with hot reload...${NC}"
(cd services/auth-service && AUTH_SERVICE_PORT=8081 air 2>&1 | sed 's/^/[AUTH] /') &

echo -e "${GREEN}Starting menu-service on port 8082 with hot reload...${NC}"
(cd services/menu-service && MENU_SERVICE_PORT=8082 air 2>&1 | sed 's/^/[MENU] /') &

echo -e "${GREEN}Starting promo-service on port 8083 with hot reload...${NC}"
(cd services/promo-service && PROMO_SERVICE_PORT=8083 air 2>&1 | sed 's/^/[PROMO] /') &

echo -e "${GREEN}Starting info-service on port 8084 with hot reload...${NC}"
(cd services/info-service && INFO_SERVICE_PORT=8084 air 2>&1 | sed 's/^/[INFO] /') &

echo -e "${GREEN}Starting media-service on port 8085 with hot reload...${NC}"
(cd services/media-service && MEDIA_SERVICE_PORT=8085 air 2>&1 | sed 's/^/[MEDIA] /') &

# Wait for services to be ready
echo -e "${YELLOW}Waiting for services to be ready...${NC}"
sleep 3

# Start agent with Air
echo -e "${GREEN}Starting Telegram bot agent with hot reload...${NC}"
(cd agent && air 2>&1 | sed 's/^/[AGENT] /') &

echo -e "${GREEN}✅ All services started with HOT RELOAD!${NC}"
echo -e "${BLUE}Edit any .go file and it will auto-reload!${NC}"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"

# Wait for all background processes
wait
