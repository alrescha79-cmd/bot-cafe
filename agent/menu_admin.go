package main

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/son/bot-cafe/shared"
)

// ADMIN MENU FUNCTIONS

func showAdminMenu(chatID int64) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Kelola Menu", "admin_menu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎉 Kelola Promo", "admin_promo"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ℹ️ Kelola Info Café", "admin_info"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📁 Kelola Kategori", "admin_category"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Kembali ke Menu Utama", "back:start"),
		),
	)

	sendMessage(chatID, "👨‍💼 *Panel Admin*\n\nPilih menu:", keyboard)
}

func showAdminMenuManagement(chatID int64) {
	text := "📋 *Manajemen Menu*\n\n"
	text += "Pilih operasi yang ingin dilakukan:\n\n"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Tambah Menu Baru", "menu_create"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 Lihat Semua Menu", "menu_read_all"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✏️ Edit Menu", "menu_update_list"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Hapus Menu", "menu_delete_list"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali ke Panel Admin", "back:admin"),
		),
	)

	sendMessage(chatID, text, keyboard)
}

func showAdminPromoManagement(chatID int64) {
	text := "🎉 *Manajemen Promo*\n\n"
	text += "Pilih operasi yang ingin dilakukan:\n\n"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Tambah Promo Baru", "promo_create"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 Lihat Semua Promo", "promo_read_all"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✏️ Edit Promo", "promo_update_list"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Hapus Promo", "promo_delete_list"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali ke Panel Admin", "back:admin"),
		),
	)

	sendMessage(chatID, text, keyboard)
}

func showAdminInfoManagement(chatID int64) {
	text := "ℹ️ *Manajemen Info Café*\n\n"
	text += "Pilih operasi yang ingin dilakukan:\n\n"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 Lihat Info Café", "info_read"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✏️ Edit Info Café", "info_update"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali ke Panel Admin", "back:admin"),
		),
	)

	sendMessage(chatID, text, keyboard)
}

func showAdminCategoryManagement(chatID int64) {
	text := "📁 *Manajemen Kategori Menu*\n\n"
	text += "Pilih operasi yang ingin dilakukan:\n\n"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Tambah Kategori Baru", "category_create"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 Lihat Semua Kategori", "category_read_all"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Hapus Kategori", "category_delete_list"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali ke Panel Admin", "back:admin"),
		),
	)

	sendMessage(chatID, text, keyboard)
}

// READ Operations - List Views

func showMenuList(chatID int64, forOperation string) {
	resp, err := httpClient.Post(menuServiceURL, shared.Request{
		Action:  "list",
		Payload: map[string]interface{}{},
	})

	if err != nil || !resp.Success {
		sendMessage(chatID, "⚠️ Gagal memuat menu.", nil)
		return
	}

	menusData, ok := resp.Data.(map[string]interface{})["menus"].([]interface{})
	text := "📋 *Daftar Menu*\n\n"
	var keyboard [][]tgbotapi.InlineKeyboardButton

	if ok && len(menusData) > 0 {
		for _, item := range menusData {
			menu := item.(map[string]interface{})
			name := menu["name"].(string)
			price := int(menu["price"].(float64))
			id := int(menu["id"].(float64))
			available := menu["is_available"].(bool)
			category := menu["category"].(string)

			status := "✅"
			if !available {
				status = "❌"
			}

			text += fmt.Sprintf("%s *%s*\n", status, name)
			text += fmt.Sprintf("   💰 %s | 📁 %s\n\n", shared.FormatPrice(price), category)

			if forOperation == "update" {
				keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("✏️ Edit "+name, fmt.Sprintf("edit_menu:%d", id)),
				))
			} else if forOperation == "delete" {
				keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🗑️ Hapus "+name, fmt.Sprintf("confirm_delete_menu:%d", id)),
				))
			}
		}
	} else {
		text += "Belum ada menu.\n"
	}

	keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "admin_menu"),
	))

	sendMessage(chatID, text, tgbotapi.NewInlineKeyboardMarkup(keyboard...))
}

func showPromoList(chatID int64, forOperation string) {
	resp, err := httpClient.Post(promoServiceURL, shared.Request{
		Action:  "list",
		Payload: map[string]interface{}{"active_only": false},
	})

	if err != nil || !resp.Success {
		sendMessage(chatID, "⚠️ Gagal memuat promo.", nil)
		return
	}

	promosData, ok := resp.Data.(map[string]interface{})["promos"].([]interface{})
	text := "🎉 *Daftar Promo*\n\n"
	var keyboard [][]tgbotapi.InlineKeyboardButton

	if ok && len(promosData) > 0 {
		for _, item := range promosData {
			promo := item.(map[string]interface{})
			title := promo["title"].(string)
			id := int(promo["id"].(float64))
			discount := int(promo["discount"].(float64))
			discountType := promo["discount_type"].(string)
			isActive := promo["is_active"].(bool)

			status := "✅"
			if !isActive {
				status = "❌"
			}

			text += fmt.Sprintf("%s *%s*\n", status, title)
			if discountType == "percentage" {
				text += fmt.Sprintf("   🎁 Diskon %d%%\n\n", discount)
			} else {
				text += fmt.Sprintf("   🎁 Diskon %s\n\n", shared.FormatPrice(discount))
			}

			if forOperation == "update" {
				keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("✏️ Edit "+title, fmt.Sprintf("edit_promo:%d", id)),
				))
			} else if forOperation == "delete" {
				keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🗑️ Hapus "+title, fmt.Sprintf("confirm_delete_promo:%d", id)),
				))
			}
		}
	} else {
		text += "Belum ada promo.\n"
	}

	keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "admin_promo"),
	))

	sendMessage(chatID, text, tgbotapi.NewInlineKeyboardMarkup(keyboard...))
}

func showCategoryList(chatID int64, forOperation string) {
	resp, err := httpClient.Post(menuServiceURL, shared.Request{
		Action: "list_categories",
	})

	if err != nil || !resp.Success {
		sendMessage(chatID, "⚠️ Gagal memuat kategori.", nil)
		return
	}

	categoriesData, ok := resp.Data.(map[string]interface{})["categories"].([]interface{})
	text := "📁 *Daftar Kategori*\n\n"
	var keyboard [][]tgbotapi.InlineKeyboardButton

	if ok && len(categoriesData) > 0 {
		for _, item := range categoriesData {
			category := item.(map[string]interface{})
			name := category["name"].(string)
			id := int(category["id"].(float64))

			text += fmt.Sprintf("• *%s*\n", name)

			if forOperation == "delete" {
				keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🗑️ Hapus "+name, fmt.Sprintf("confirm_delete_category:%d", id)),
				))
			}
		}
	} else {
		text += "Belum ada kategori.\n"
	}

	keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "admin_category"),
	))

	sendMessage(chatID, text, tgbotapi.NewInlineKeyboardMarkup(keyboard...))
}

func showCafeInfoDetail(chatID int64) {
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

	text := "ℹ️ *Informasi Café*\n\n"
	text += fmt.Sprintf("*Nama:* %s\n", name)
	text += fmt.Sprintf("*Alamat:* %s\n", address)
	text += fmt.Sprintf("*Telepon:* %s\n", phone)
	text += fmt.Sprintf("*Jam Operasional:* %s - %s\n", openingHour, closingHour)
	if description != "" {
		text += fmt.Sprintf("*Deskripsi:* %s\n", description)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "admin_info"),
		),
	)

	sendMessage(chatID, text, keyboard)
}

// Dialog state handlers

func handleDialogState(msg *tgbotapi.Message, state string) {
	userID := msg.From.ID

	switch state {
	case "add_menu_name":
		handleAddMenuName(msg, userID)
	case "add_menu_price":
		handleAddMenuPrice(msg, userID)
	case "add_menu_category":
		handleAddMenuCategory(msg, userID)
	case "add_menu_description":
		handleAddMenuDescription(msg, userID)
	case "add_promo_title":
		handleAddPromoTitle(msg, userID)
	case "add_promo_description":
		handleAddPromoDescription(msg, userID)
	case "add_promo_discount":
		handleAddPromoDiscount(msg, userID)
	case "add_promo_discount_type":
		handleAddPromoDiscountType(msg, userID)
	case "add_promo_start_date":
		handleAddPromoStartDate(msg, userID)
	case "add_promo_end_date":
		handleAddPromoEndDate(msg, userID)
	case "add_category_name":
		handleAddCategoryName(msg, userID)
	default:
		delete(userStates, userID)
		sendMessage(msg.Chat.ID, "State tidak dikenal. Gunakan /cancel untuk membatalkan.", nil)
	}
}

func startAddMenuDialog(chatID int64, userID int64) {
	userStates[userID] = "add_menu_name"
	userTempData[userID] = make(map[string]interface{})
	sendMessage(chatID, "➕ *Tambah Menu Baru*\n\nMasukkan nama menu:\n\n(Ketik /cancel untuk membatalkan)", nil)
}

func handleAddMenuName(msg *tgbotapi.Message, userID int64) {
	name := strings.TrimSpace(msg.Text)
	if name == "" {
		sendMessage(msg.Chat.ID, "⚠️ Nama menu tidak boleh kosong. Coba lagi:", nil)
		return
	}

	userTempData[userID]["name"] = name
	userStates[userID] = "add_menu_price"
	sendMessage(msg.Chat.ID, "Masukkan harga menu (angka saja, tanpa Rp):", nil)
}

func handleAddMenuPrice(msg *tgbotapi.Message, userID int64) {
	price, err := strconv.Atoi(strings.TrimSpace(msg.Text))
	if err != nil || price < 0 {
		sendMessage(msg.Chat.ID, "⚠️ Harga tidak valid. Masukkan angka positif:", nil)
		return
	}

	userTempData[userID]["price"] = price
	userStates[userID] = "add_menu_category"

	// Get categories
	resp, err := httpClient.Post(menuServiceURL, shared.Request{
		Action: "list_categories",
	})

	text := "Pilih kategori:\n\n"
	if err == nil && resp.Success {
		categoriesData, ok := resp.Data.(map[string]interface{})["categories"].([]interface{})
		if ok {
			for _, item := range categoriesData {
				category := item.(map[string]interface{})
				name := category["name"].(string)
				text += fmt.Sprintf("• %s\n", name)
			}
		}
	}
	text += "\nKetik nama kategori:"

	sendMessage(msg.Chat.ID, text, nil)
}

func handleAddMenuCategory(msg *tgbotapi.Message, userID int64) {
	category := strings.TrimSpace(msg.Text)
	if category == "" {
		sendMessage(msg.Chat.ID, "⚠️ Kategori tidak boleh kosong. Coba lagi:", nil)
		return
	}

	userTempData[userID]["category"] = category
	userStates[userID] = "add_menu_description"
	sendMessage(msg.Chat.ID, "Masukkan deskripsi menu (atau ketik - untuk skip):", nil)
}

func handleAddMenuDescription(msg *tgbotapi.Message, userID int64) {
	description := strings.TrimSpace(msg.Text)
	if description == "-" {
		description = ""
	}

	userTempData[userID]["description"] = description

	// Create menu
	data := userTempData[userID]
	resp, err := httpClient.Post(menuServiceURL, shared.Request{
		Action: "create",
		Payload: map[string]interface{}{
			"name":         data["name"],
			"price":        data["price"],
			"category":     data["category"],
			"description":  description,
			"is_available": true,
		},
	})

	delete(userStates, userID)
	delete(userTempData, userID)

	if err != nil || !resp.Success {
		sendMessage(msg.Chat.ID, "⚠️ Gagal menambahkan menu. Silakan coba lagi.", nil)
		return
	}

	menuData := resp.Data.(map[string]interface{})["menu"].(map[string]interface{})
	name := menuData["name"].(string)
	price := int(menuData["price"].(float64))

	text := fmt.Sprintf("✅ *Menu berhasil ditambahkan!*\n\n🍽️ %s\n💰 %s", name, shared.FormatPrice(price))
	sendMessage(msg.Chat.ID, text, nil)
}

func startEditMenuDialog(chatID int64, userID int64, menuID int) {
	sendMessage(chatID, "✏️ Fitur edit menu akan segera tersedia.", nil)
}

func deleteMenu(chatID int64, menuID int) {
	resp, err := httpClient.Post(menuServiceURL, shared.Request{
		Action: "delete",
		Payload: map[string]interface{}{
			"id": menuID,
		},
	})

	if err != nil || !resp.Success {
		sendMessage(chatID, "⚠️ Gagal menghapus menu.", nil)
		return
	}

	sendMessage(chatID, "✅ Menu berhasil dihapus!", nil)
	showAdminMenuManagement(chatID)
}

func startAddPromoDialog(chatID int64, userID int64) {
	userStates[userID] = "add_promo_title"
	userTempData[userID] = make(map[string]interface{})
	sendMessage(chatID, "➕ *Tambah Promo Baru*\n\nMasukkan judul promo:\n\n(Ketik /cancel untuk membatalkan)", nil)
}

func handleAddPromoTitle(msg *tgbotapi.Message, userID int64) {
	title := strings.TrimSpace(msg.Text)
	if title == "" {
		sendMessage(msg.Chat.ID, "⚠️ Judul promo tidak boleh kosong. Coba lagi:", nil)
		return
	}

	userTempData[userID]["title"] = title
	userStates[userID] = "add_promo_description"
	sendMessage(msg.Chat.ID, "Masukkan deskripsi promo (atau ketik - untuk skip):", nil)
}

func handleAddPromoDescription(msg *tgbotapi.Message, userID int64) {
	description := strings.TrimSpace(msg.Text)
	if description == "-" {
		description = ""
	}

	userTempData[userID]["description"] = description
	userStates[userID] = "add_promo_discount_type"
	sendMessage(msg.Chat.ID, "Pilih tipe diskon:\n\nKetik *percentage* atau *amount*:", nil)
}

func handleAddPromoDiscountType(msg *tgbotapi.Message, userID int64) {
	discountType := strings.ToLower(strings.TrimSpace(msg.Text))
	if discountType != "percentage" && discountType != "amount" {
		sendMessage(msg.Chat.ID, "⚠️ Tipe harus 'percentage' atau 'amount'. Coba lagi:", nil)
		return
	}

	userTempData[userID]["discount_type"] = discountType
	userStates[userID] = "add_promo_discount"

	if discountType == "percentage" {
		sendMessage(msg.Chat.ID, "Masukkan jumlah diskon (dalam %, angka saja):", nil)
	} else {
		sendMessage(msg.Chat.ID, "Masukkan jumlah diskon (dalam Rp, angka saja):", nil)
	}
}

func handleAddPromoDiscount(msg *tgbotapi.Message, userID int64) {
	discount, err := strconv.Atoi(strings.TrimSpace(msg.Text))
	if err != nil || discount < 0 {
		sendMessage(msg.Chat.ID, "⚠️ Diskon tidak valid. Masukkan angka positif:", nil)
		return
	}

	userTempData[userID]["discount"] = discount
	userStates[userID] = "add_promo_start_date"
	sendMessage(msg.Chat.ID, "Masukkan tanggal mulai (format: YYYY-MM-DD, contoh: 2025-01-01):", nil)
}

func handleAddPromoStartDate(msg *tgbotapi.Message, userID int64) {
	startDate := strings.TrimSpace(msg.Text)

	userTempData[userID]["start_date"] = startDate
	userStates[userID] = "add_promo_end_date"
	sendMessage(msg.Chat.ID, "Masukkan tanggal akhir (format: YYYY-MM-DD):", nil)
}

func handleAddPromoEndDate(msg *tgbotapi.Message, userID int64) {
	endDate := strings.TrimSpace(msg.Text)

	// Create promo
	data := userTempData[userID]
	resp, err := httpClient.Post(promoServiceURL, shared.Request{
		Action: "create",
		Payload: map[string]interface{}{
			"title":         data["title"],
			"description":   data["description"],
			"discount":      data["discount"],
			"discount_type": data["discount_type"],
			"start_date":    data["start_date"],
			"end_date":      endDate,
			"is_active":     true,
		},
	})

	delete(userStates, userID)
	delete(userTempData, userID)

	if err != nil || !resp.Success {
		errMsg := "⚠️ Gagal menambahkan promo. Periksa format tanggal (YYYY-MM-DD)."
		if resp != nil && resp.Error != nil {
			errMsg = "⚠️ " + resp.Error.Message
		}
		sendMessage(msg.Chat.ID, errMsg, nil)
		return
	}

	promoData := resp.Data.(map[string]interface{})["promo"].(map[string]interface{})
	title := promoData["title"].(string)

	text := fmt.Sprintf("✅ *Promo berhasil ditambahkan!*\n\n🎁 %s", title)
	sendMessage(msg.Chat.ID, text, nil)
}

func deletePromo(chatID int64, promoID int) {
	resp, err := httpClient.Post(promoServiceURL, shared.Request{
		Action: "delete",
		Payload: map[string]interface{}{
			"id": promoID,
		},
	})

	if err != nil || !resp.Success {
		sendMessage(chatID, "⚠️ Gagal menghapus promo.", nil)
		return
	}

	sendMessage(chatID, "✅ Promo berhasil dihapus.", nil)
	showAdminPromoManagement(chatID)
}

func deleteCategory(chatID int64, categoryID int) {
	resp, err := httpClient.Post(menuServiceURL, shared.Request{
		Action: "delete_category",
		Payload: map[string]interface{}{
			"id": categoryID,
		},
	})

	if err != nil || !resp.Success {
		errMsg := "⚠️ Gagal menghapus kategori."
		if resp.Error != nil {
			errMsg += "\n" + resp.Error.Message
		}
		sendMessage(chatID, errMsg, nil)
		return
	}

	sendMessage(chatID, "✅ Kategori berhasil dihapus.", nil)
	showAdminCategoryManagement(chatID)
}

func startAddCategoryDialog(chatID int64, userID int64) {
	userStates[userID] = "add_category_name"
	userTempData[userID] = make(map[string]interface{})
	sendMessage(chatID, "📁 *Tambah Kategori Baru*\n\nMasukkan nama kategori:\n\n_Ketik /cancel untuk membatalkan_", nil)
}

func handleAddCategoryName(msg *tgbotapi.Message, userID int64) {
	categoryName := strings.TrimSpace(msg.Text)
	if categoryName == "" {
		sendMessage(msg.Chat.ID, "⚠️ Nama kategori tidak boleh kosong. Coba lagi:", nil)
		return
	}

	// Create category
	resp, err := httpClient.Post(menuServiceURL, shared.Request{
		Action: "create_category",
		Payload: map[string]interface{}{
			"name": categoryName,
		},
	})

	delete(userStates, userID)
	delete(userTempData, userID)

	if err != nil || !resp.Success {
		errMsg := "⚠️ Gagal menambahkan kategori."
		if resp.Error != nil {
			errMsg += "\n" + resp.Error.Message
		}
		sendMessage(msg.Chat.ID, errMsg, nil)
		return
	}

	sendMessage(msg.Chat.ID, fmt.Sprintf("✅ Kategori *%s* berhasil ditambahkan!", categoryName), nil)
	showAdminCategoryManagement(msg.Chat.ID)
}
