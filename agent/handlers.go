package main

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/son/bot-cafe/shared"
)

func handleMessage(msg *tgbotapi.Message) {
	userID := msg.From.ID
	_ = msg.From.UserName // username tersimpan tapi tidak digunakan untuk saat ini

	// Check for commands
	if msg.IsCommand() {
		handleCommand(msg)
		return
	}

	// Check if user is in dialog state
	if state, exists := userStates[userID]; exists {
		handleDialogState(msg, state)
		return
	}

	// Default response
	text := "Gunakan menu atau ketik /start untuk memulai."
	sendMessage(msg.Chat.ID, text, nil)
}

func handleCommand(msg *tgbotapi.Message) {
	userID := msg.From.ID
	username := msg.From.UserName
	command := msg.Command()

	switch command {
	case "start":
		handleStartCommand(msg)
	case "menu":
		showUserMenu(msg.Chat.ID)
	case "promo":
		showPromos(msg.Chat.ID, true)
	case "info":
		showCafeInfo(msg.Chat.ID)
	case "admin":
		if isAdmin(userID, username) {
			showAdminMenu(msg.Chat.ID)
		} else {
			sendMessage(msg.Chat.ID, "⚠️ Anda tidak memiliki akses admin.", nil)
		}
	case "cancel":
		delete(userStates, userID)
		delete(userTempData, userID)
		sendMessage(msg.Chat.ID, "❌ Operasi dibatalkan.", nil)
	default:
		sendMessage(msg.Chat.ID, "Perintah tidak dikenal. Gunakan /start untuk melihat menu.", nil)
	}
}

func handleStartCommand(msg *tgbotapi.Message) {
	userID := msg.From.ID
	username := msg.From.UserName

	shared.LogInfo("[START] User %d (@%s) executed /start", userID, username)

	// Admins go directly to admin menu
	isAdminUser := isAdmin(userID, username)
	shared.LogInfo("[START] isAdmin check result: %v for user %d (@%s)", isAdminUser, userID, username)

	if isAdminUser {
		welcomeText := "👋 Selamat datang di Bot Café!\n\n"
		welcomeText += "Anda login sebagai *Admin*."

		msg2 := tgbotapi.NewMessage(msg.Chat.ID, welcomeText)
		msg2.ParseMode = "Markdown"
		bot.Send(msg2)

		// Show admin menu directly
		shared.LogInfo("[START] Showing ADMIN menu for user %d (@%s)", userID, username)
		showAdminMenu(msg.Chat.ID)
		return
	}

	// Regular users see the standard welcome menu
	shared.LogInfo("[START] Showing USER menu for user %d (@%s)", userID, username)
	welcomeText := "👋 Selamat datang di Bot Café!\n\n"
	welcomeText += "Pilih menu di bawah ini:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Lihat Menu", "show_user_menu"),
			tgbotapi.NewInlineKeyboardButtonData("🎉 Lihat Promo", "show_promo"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ℹ️ Info Café", "show_info"),
		),
	)

	msg2 := tgbotapi.NewMessage(msg.Chat.ID, welcomeText)
	msg2.ParseMode = "Markdown"
	msg2.ReplyMarkup = keyboard
	bot.Send(msg2)
}

func handleCallback(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID
	username := callback.From.UserName
	data := callback.Data

	// Answer callback to remove loading state
	bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	parts := strings.Split(data, ":")
	action := parts[0]

	switch action {
	case "menu_category":
		if len(parts) > 1 {
			category := parts[1]
			showMenuByCategory(callback.Message.Chat.ID, category)
		}
	case "menu_detail":
		if len(parts) > 1 {
			menuID, _ := strconv.Atoi(parts[1])
			showMenuDetail(callback.Message.Chat.ID, menuID)
		}
	case "admin_menu":
		if !isAdmin(userID, username) {
			sendMessage(callback.Message.Chat.ID, "⚠️ Akses ditolak.", nil)
			return
		}
		showAdminMenuManagement(callback.Message.Chat.ID)
	case "admin_promo":
		if !isAdmin(userID, username) {
			sendMessage(callback.Message.Chat.ID, "⚠️ Akses ditolak.", nil)
			return
		}
		showAdminPromoManagement(callback.Message.Chat.ID)
	case "admin_info":
		if !isAdmin(userID, username) {
			sendMessage(callback.Message.Chat.ID, "⚠️ Akses ditolak.", nil)
			return
		}
		showAdminInfoManagement(callback.Message.Chat.ID)
	case "admin_category":
		if !isAdmin(userID, username) {
			sendMessage(callback.Message.Chat.ID, "⚠️ Akses ditolak.", nil)
			return
		}
		showAdminCategoryManagement(callback.Message.Chat.ID)

	// Menu CRUD Operations
	case "menu_create":
		if !isAdmin(userID, username) {
			return
		}
		startAddMenuDialog(callback.Message.Chat.ID, userID)
	case "menu_read_all":
		if !isAdmin(userID, username) {
			return
		}
		showMenuList(callback.Message.Chat.ID, "view")
	case "menu_update_list":
		if !isAdmin(userID, username) {
			return
		}
		showMenuList(callback.Message.Chat.ID, "update")
	case "menu_delete_list":
		if !isAdmin(userID, username) {
			return
		}
		showMenuList(callback.Message.Chat.ID, "delete")

	// Promo CRUD Operations
	case "promo_create":
		if !isAdmin(userID, username) {
			return
		}
		startAddPromoDialog(callback.Message.Chat.ID, userID)
	case "promo_read_all":
		if !isAdmin(userID, username) {
			return
		}
		showPromoList(callback.Message.Chat.ID, "view")
	case "promo_update_list":
		if !isAdmin(userID, username) {
			return
		}
		showPromoList(callback.Message.Chat.ID, "update")
	case "promo_delete_list":
		if !isAdmin(userID, username) {
			return
		}
		showPromoList(callback.Message.Chat.ID, "delete")

	// Category CRUD Operations
	case "category_create":
		if !isAdmin(userID, username) {
			return
		}
		startAddCategoryDialog(callback.Message.Chat.ID, userID)
	case "category_read_all":
		if !isAdmin(userID, username) {
			return
		}
		showCategoryList(callback.Message.Chat.ID, "view")
	case "category_delete_list":
		if !isAdmin(userID, username) {
			return
		}
		showCategoryList(callback.Message.Chat.ID, "delete")

	// Info CRUD Operations
	case "info_read":
		if !isAdmin(userID, username) {
			return
		}
		showCafeInfoDetail(callback.Message.Chat.ID)
	case "info_update":
		if !isAdmin(userID, username) {
			return
		}
		sendMessage(callback.Message.Chat.ID, "⚠️ Fitur edit info café akan segera ditambahkan.", nil)

	// Legacy handlers (keep for compatibility)
	case "add_menu":
		if !isAdmin(userID, username) {
			return
		}
		startAddMenuDialog(callback.Message.Chat.ID, userID)
	case "edit_menu":
		if !isAdmin(userID, username) {
			return
		}
		if len(parts) > 1 {
			menuID, _ := strconv.Atoi(parts[1])
			startEditMenuDialog(callback.Message.Chat.ID, userID, menuID)
		}
	case "confirm_delete_menu":
		if !isAdmin(userID, username) {
			return
		}
		if len(parts) > 1 {
			menuID, _ := strconv.Atoi(parts[1])
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("✅ Ya, Hapus", fmt.Sprintf("delete_menu:%d", menuID)),
					tgbotapi.NewInlineKeyboardButtonData("❌ Batal", "menu_delete_list"),
				),
			)
			sendMessage(callback.Message.Chat.ID, "⚠️ Apakah Anda yakin ingin menghapus menu ini?", keyboard)
		}
	case "delete_menu":
		if !isAdmin(userID, username) {
			return
		}
		if len(parts) > 1 {
			menuID, _ := strconv.Atoi(parts[1])
			deleteMenu(callback.Message.Chat.ID, menuID)
		}
	case "edit_promo":
		if !isAdmin(userID, username) {
			return
		}
		if len(parts) > 1 {
			promoID, _ := strconv.Atoi(parts[1])
			sendMessage(callback.Message.Chat.ID, fmt.Sprintf("⚠️ Fitur edit promo akan segera ditambahkan. (Promo ID: %d)", promoID), nil)
		}
	case "confirm_delete_promo":
		if !isAdmin(userID, username) {
			return
		}
		if len(parts) > 1 {
			promoID, _ := strconv.Atoi(parts[1])
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("✅ Ya, Hapus", fmt.Sprintf("delete_promo:%d", promoID)),
					tgbotapi.NewInlineKeyboardButtonData("❌ Batal", "promo_delete_list"),
				),
			)
			sendMessage(callback.Message.Chat.ID, "⚠️ Apakah Anda yakin ingin menghapus promo ini?", keyboard)
		}
	case "delete_promo":
		if !isAdmin(userID, username) {
			return
		}
		if len(parts) > 1 {
			promoID, _ := strconv.Atoi(parts[1])
			deletePromo(callback.Message.Chat.ID, promoID)
		}
	case "confirm_delete_category":
		if !isAdmin(userID, username) {
			return
		}
		if len(parts) > 1 {
			categoryID, _ := strconv.Atoi(parts[1])
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("✅ Ya, Hapus", fmt.Sprintf("delete_category:%d", categoryID)),
					tgbotapi.NewInlineKeyboardButtonData("❌ Batal", "category_delete_list"),
				),
			)
			sendMessage(callback.Message.Chat.ID, "⚠️ Apakah Anda yakin ingin menghapus kategori ini?\n\n*Perhatian:* Menu dengan kategori ini mungkin terpengaruh.", keyboard)
		}
	case "delete_category":
		if !isAdmin(userID, username) {
			return
		}
		if len(parts) > 1 {
			categoryID, _ := strconv.Atoi(parts[1])
			deleteCategory(callback.Message.Chat.ID, categoryID)
		}
	case "show_user_menu":
		showUserMenu(callback.Message.Chat.ID)
	case "show_promo":
		showPromos(callback.Message.Chat.ID, true)
	case "show_info":
		showCafeInfo(callback.Message.Chat.ID)
	case "show_admin_panel":
		if !isAdmin(userID, username) {
			sendMessage(callback.Message.Chat.ID, "⚠️ Akses ditolak.", nil)
			return
		}
		showAdminMenu(callback.Message.Chat.ID)
	case "back":
		if len(parts) > 1 {
			target := parts[1]
			switch target {
			case "admin":
				showAdminMenu(callback.Message.Chat.ID)
			case "user":
				showUserMenu(callback.Message.Chat.ID)
			case "start":
				handleStartCommand(&tgbotapi.Message{Chat: &tgbotapi.Chat{ID: callback.Message.Chat.ID}, From: callback.From})
			}
		}
	}
}

func sendMessage(chatID int64, text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	bot.Send(msg)
}

func sendMessageWithoutMarkdown(chatID int64, text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	bot.Send(msg)
}
