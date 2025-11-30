package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/son/bot-cafe/shared"
)

// USER MENU FUNCTIONS

func showUserMenu(chatID int64) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("☕ Coffee", "menu_category:Coffee"),
			tgbotapi.NewInlineKeyboardButtonData("🍽️ Makanan", "menu_category:Makanan"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🥤 Minuman", "menu_category:Minuman"),
			tgbotapi.NewInlineKeyboardButtonData("🍪 Snack", "menu_category:Snack"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Kembali ke Menu Utama", "back:start"),
		),
	)

	sendMessage(chatID, "📋 *Menu Café*\n\nPilih kategori:", keyboard)
}

func showMenuByCategory(chatID int64, category string) {
	resp, err := httpClient.Post(menuServiceURL, shared.Request{
		Action: "list",
		Payload: map[string]interface{}{
			"category":       category,
			"available_only": true,
		},
	})

	if err != nil || !resp.Success {
		sendMessage(chatID, "⚠️ Gagal memuat menu.", nil)
		return
	}

	menusData, ok := resp.Data.(map[string]interface{})["menus"].([]interface{})
	if !ok || len(menusData) == 0 {
		sendMessage(chatID, fmt.Sprintf("Tidak ada menu dalam kategori *%s*", category), nil)
		return
	}

	text := fmt.Sprintf("📋 *Menu %s*\n\n", category)
	var keyboard [][]tgbotapi.InlineKeyboardButton

	for _, item := range menusData {
		menu := item.(map[string]interface{})
		name := menu["name"].(string)
		price := int(menu["price"].(float64))
		id := int(menu["id"].(float64))

		text += fmt.Sprintf("🍽️ *%s*\nHarga: %s\n\n", name, shared.FormatPrice(price))

		keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 "+name, fmt.Sprintf("menu_detail:%d", id)),
		))
	}

	keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "back:user"),
	))

	sendMessage(chatID, text, tgbotapi.NewInlineKeyboardMarkup(keyboard...))
}

func showMenuDetail(chatID int64, menuID int) {
	resp, err := httpClient.Post(menuServiceURL, shared.Request{
		Action: "read",
		Payload: map[string]interface{}{
			"id": menuID,
		},
	})

	if err != nil || !resp.Success {
		sendMessage(chatID, "⚠️ Menu tidak ditemukan.", nil)
		return
	}

	menuData := resp.Data.(map[string]interface{})["menu"].(map[string]interface{})
	name := menuData["name"].(string)
	description := menuData["description"].(string)
	price := int(menuData["price"].(float64))
	category := menuData["category"].(string)

	text := fmt.Sprintf("🍽️ *%s*\n\n", name)
	if description != "" {
		text += fmt.Sprintf("%s\n\n", description)
	}
	text += fmt.Sprintf("💰 Harga: %s\n", shared.FormatPrice(price))
	text += fmt.Sprintf("📁 Kategori: %s\n", category)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", fmt.Sprintf("menu_category:%s", category)),
		),
	)

	sendMessage(chatID, text, keyboard)
}

func showPromos(chatID int64, activeOnly bool) {
	resp, err := httpClient.Post(promoServiceURL, shared.Request{
		Action: "list",
		Payload: map[string]interface{}{
			"active_only": activeOnly,
		},
	})

	if err != nil || !resp.Success {
		sendMessage(chatID, "⚠️ Gagal memuat promo.", nil)
		return
	}

	promosData, ok := resp.Data.(map[string]interface{})["promos"].([]interface{})
	if !ok || len(promosData) == 0 {
		sendMessage(chatID, "Belum ada promo tersedia saat ini.", nil)
		return
	}

	text := "🎉 *Promo Tersedia*\n\n"

	for _, item := range promosData {
		promo := item.(map[string]interface{})
		title := promo["title"].(string)
		description := promo["description"].(string)
		discount := int(promo["discount"].(float64))
		discountType := promo["discount_type"].(string)

		text += fmt.Sprintf("🎁 *%s*\n", title)
		if description != "" {
			text += fmt.Sprintf("%s\n", description)
		}

		if discountType == "percentage" {
			text += fmt.Sprintf("Diskon: %d%%\n", discount)
		} else {
			text += fmt.Sprintf("Diskon: %s\n", shared.FormatPrice(discount))
		}
		text += "\n"
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Kembali ke Menu Utama", "back:start"),
		),
	)

	sendMessage(chatID, text, keyboard)
}

func showCafeInfo(chatID int64) {
	resp, err := httpClient.Post(infoServiceURL, shared.Request{
		Action: "read",
	})

	if err != nil || !resp.Success {
		sendMessage(chatID, "⚠️ Gagal memuat informasi café.", nil)
		return
	}

	infoData := resp.Data.(map[string]interface{})["info"].(map[string]interface{})
	name := infoData["name"].(string)
	address := infoData["address"].(string)
	phone := infoData["phone"].(string)
	openingHour := infoData["opening_hour"].(string)
	closingHour := infoData["closing_hour"].(string)
	description := infoData["description"].(string)

	text := fmt.Sprintf("ℹ️ *%s*\n\n", name)
	if description != "" {
		text += fmt.Sprintf("%s\n\n", description)
	}
	text += fmt.Sprintf("📍 Alamat: %s\n", address)
	text += fmt.Sprintf("📞 Telepon: %s\n", phone)
	text += fmt.Sprintf("🕐 Jam Buka: %s - %s\n", openingHour, closingHour)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Kembali ke Menu Utama", "back:start"),
		),
	)

	sendMessage(chatID, text, keyboard)
}
