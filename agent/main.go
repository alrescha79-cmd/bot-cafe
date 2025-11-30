package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/son/bot-cafe/shared"
)

var (
	bot            *tgbotapi.BotAPI
	httpClient     *shared.HTTPClient
	adminIDs       []string
	adminUsernames []string

	// Service URLs
	authServiceURL  string
	menuServiceURL  string
	promoServiceURL string
	infoServiceURL  string
	mediaServiceURL string

	// User states untuk dialog CRUD
	userStates   = make(map[int64]string)
	userTempData = make(map[int64]map[string]interface{})
)

type VarsConfig struct {
	AdminTelegramIDs []string `json:"admin_telegram_ids"`
	AdminUsernames   []string `json:"admin_usernames"`
}

func main() {
	// Load environment variables
	godotenv.Load()

	// Load admin vars
	if err := loadAdminVars(); err != nil {
		log.Printf("Warning: Failed to load admin vars: %v", err)
	}

	// Get bot token
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	// Get service URLs
	authServiceURL = getEnv("AUTH_SERVICE_URL", "http://localhost:8081")
	menuServiceURL = getEnv("MENU_SERVICE_URL", "http://localhost:8082")
	promoServiceURL = getEnv("PROMO_SERVICE_URL", "http://localhost:8083")
	infoServiceURL = getEnv("INFO_SERVICE_URL", "http://localhost:8084")
	mediaServiceURL = getEnv("MEDIA_SERVICE_URL", "http://localhost:8085")

	// Initialize bot
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = false
	shared.LogInfo("Authorized on account %s", bot.Self.UserName)

	// Initialize HTTP client
	httpClient = shared.NewHTTPClient()

	// Configure updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Handle updates
	for update := range updates {
		if update.Message != nil {
			handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			handleCallback(update.CallbackQuery)
		}
	}
}

func loadAdminVars() error {
	varsFile := getEnv("ADMIN_VARS_FILE", ".vars.json")

	// Check if file exists
	if _, err := os.Stat(varsFile); os.IsNotExist(err) {
		log.Fatalf("ERROR: %s not found! Please create it with admin configuration.\nExample:\n{\n  \"admin_telegram_ids\": [\"YOUR_TELEGRAM_ID\"],\n  \"admin_usernames\": [\"your_username\"]\n}", varsFile)
	}

	// Read file
	data, err := os.ReadFile(varsFile)
	if err != nil {
		log.Fatalf("ERROR: Failed to read %s: %v", varsFile, err)
	}

	var config VarsConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("ERROR: Invalid JSON format in %s: %v", varsFile, err)
	}

	// Validate admin config
	if len(config.AdminTelegramIDs) == 0 && len(config.AdminUsernames) == 0 {
		log.Fatalf("ERROR: No admin configured in %s! Please add at least one admin_telegram_ids or admin_usernames", varsFile)
	}

	adminIDs = config.AdminTelegramIDs
	adminUsernames = config.AdminUsernames
	shared.LogInfo("Loaded %d admin IDs and %d admin usernames from %s", len(adminIDs), len(adminUsernames), varsFile)
	return nil
}

func isAdmin(userID int64, username string) bool {
	userIDStr := strconv.FormatInt(userID, 10)

	shared.LogInfo("[AUTH] Checking admin status for user %d (@%s)", userID, username)
	shared.LogInfo("[AUTH] Loaded admin IDs: %v", adminIDs)
	shared.LogInfo("[AUTH] Loaded admin usernames: %v", adminUsernames)

	// Check in vars
	for _, id := range adminIDs {
		shared.LogInfo("[AUTH] Comparing userID '%s' with admin ID '%s'", userIDStr, id)
		if id == userIDStr {
			shared.LogInfo("[AUTH] ✅ MATCH! User %d is admin (by ID)", userID)
			return true
		}
	}
	for _, un := range adminUsernames {
		shared.LogInfo("[AUTH] Comparing username '%s' with admin username '%s'", username, un)
		if un == username {
			shared.LogInfo("[AUTH] ✅ MATCH! User @%s is admin (by username)", username)
			return true
		}
	}

	// Verify with auth service
	shared.LogInfo("[AUTH] No match in vars, checking auth service...")
	resp, err := httpClient.Post(authServiceURL, shared.Request{
		Action: "verify",
		Payload: map[string]interface{}{
			"telegram_id": userIDStr,
		},
	})

	if err == nil && resp.Success {
		shared.LogInfo("[AUTH] ✅ User %d is admin (verified by auth service)", userID)
		return true
	}

	shared.LogInfo("[AUTH] ❌ User %d (@%s) is NOT admin", userID, username)
	return false
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
