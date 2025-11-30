package main

import (
	"encoding/json"
	"net/http"

	"github.com/son/bot-cafe/shared"
)

// Handler handles HTTP requests
type Handler struct {
	repo *Repository
}

// NewHandler creates a new handler
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// HandleRequest handles all incoming requests
func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, shared.NewError(shared.ErrCodeInvalidInput, "Method not allowed", nil))
		return
	}

	var req shared.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, shared.NewError(shared.ErrCodeInvalidInput, "Invalid request format", err))
		return
	}

	var response *shared.Response

	switch req.Action {
	case "create":
		response = h.createMenu(req.Payload)
	case "read":
		response = h.getMenu(req.Payload)
	case "update":
		response = h.updateMenu(req.Payload)
	case "delete":
		response = h.deleteMenu(req.Payload)
	case "list":
		response = h.listMenus(req.Payload)
	case "list_categories":
		response = h.listCategories()
	case "create_category":
		response = h.createCategory(req.Payload)
	case "delete_category":
		response = h.deleteCategory(req.Payload)
	default:
		sendErrorResponse(w, shared.NewError(shared.ErrCodeInvalidInput, "Unknown action", nil))
		return
	}

	sendResponse(w, response)
}

// createMenu creates a new menu
func (h *Handler) createMenu(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	name, _ := data["name"].(string)
	description, _ := data["description"].(string)
	priceRaw := data["price"]
	category, _ := data["category"].(string)
	photoURL, _ := data["photo_url"].(string)
	isAvailable := true
	if val, ok := data["is_available"].(bool); ok {
		isAvailable = val
	}

	// Validate inputs
	if err := shared.ValidateNotEmpty(name, "Nama menu"); err != nil {
		return errorResponse(err.(*shared.AppError))
	}
	if err := shared.ValidateNotEmpty(category, "Kategori"); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	price, err := shared.ValidatePrice(priceRaw)
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	if err := shared.ValidatePhotoURL(photoURL); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	menu := &Menu{
		Name:        shared.SanitizeInput(name),
		Description: shared.SanitizeInput(description),
		Price:       price,
		Category:    category,
		PhotoURL:    photoURL,
		IsAvailable: isAvailable,
	}

	result, err := h.repo.CreateMenu(menu)
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"menu": result,
	})
}

// getMenu gets a menu by ID
func (h *Handler) getMenu(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	id, ok := data["id"].(float64)
	if !ok {
		return errorResponse(shared.NewInvalidInputError("ID diperlukan"))
	}

	menu, err := h.repo.GetMenuByID(int(id))
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"menu": menu,
	})
}

// updateMenu updates a menu
func (h *Handler) updateMenu(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	id, ok := data["id"].(float64)
	if !ok {
		return errorResponse(shared.NewInvalidInputError("ID diperlukan"))
	}

	// Get existing menu
	menu, err := h.repo.GetMenuByID(int(id))
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	// Update fields if provided
	if name, ok := data["name"].(string); ok && name != "" {
		menu.Name = shared.SanitizeInput(name)
	}
	if desc, ok := data["description"].(string); ok {
		menu.Description = shared.SanitizeInput(desc)
	}
	if priceRaw, ok := data["price"]; ok {
		price, err := shared.ValidatePrice(priceRaw)
		if err != nil {
			return errorResponse(err.(*shared.AppError))
		}
		menu.Price = price
	}
	if category, ok := data["category"].(string); ok && category != "" {
		menu.Category = category
	}
	if photoURL, ok := data["photo_url"].(string); ok {
		if err := shared.ValidatePhotoURL(photoURL); err != nil {
			return errorResponse(err.(*shared.AppError))
		}
		menu.PhotoURL = photoURL
	}
	if isAvailable, ok := data["is_available"].(bool); ok {
		menu.IsAvailable = isAvailable
	}

	if err := h.repo.UpdateMenu(menu); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"menu": menu,
	})
}

// deleteMenu deletes a menu
func (h *Handler) deleteMenu(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	id, ok := data["id"].(float64)
	if !ok {
		return errorResponse(shared.NewInvalidInputError("ID diperlukan"))
	}

	if err := h.repo.DeleteMenu(int(id)); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"message": "Menu berhasil dihapus",
	})
}

// listMenus lists menus with optional filters
func (h *Handler) listMenus(payload interface{}) *shared.Response {
	category := ""
	availableOnly := false

	if data, ok := payload.(map[string]interface{}); ok {
		if cat, ok := data["category"].(string); ok {
			category = cat
		}
		if avail, ok := data["available_only"].(bool); ok {
			availableOnly = avail
		}
	}

	menus, err := h.repo.ListMenus(category, availableOnly)
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"menus": menus,
	})
}

// listCategories lists all categories
func (h *Handler) listCategories() *shared.Response {
	categories, err := h.repo.ListCategories()
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"categories": categories,
	})
}

// createCategory creates a new category
func (h *Handler) createCategory(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	name, _ := data["name"].(string)
	if err := shared.ValidateNotEmpty(name, "Nama kategori"); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	category, err := h.repo.CreateCategory(shared.SanitizeInput(name))
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"category": category,
	})
}

// deleteCategory deletes a category
func (h *Handler) deleteCategory(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	name, _ := data["name"].(string)
	if err := shared.ValidateNotEmpty(name, "Nama kategori"); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	if err := h.repo.DeleteCategory(name); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"message": "Kategori berhasil dihapus",
	})
}

// Helper functions
func successResponse(data interface{}) *shared.Response {
	return &shared.Response{
		Success: true,
		Data:    data,
	}
}

func errorResponse(err *shared.AppError) *shared.Response {
	return &shared.Response{
		Success: false,
		Error: &shared.ErrorInfo{
			Code:    err.Code,
			Message: err.Message,
		},
	}
}

func sendResponse(w http.ResponseWriter, response *shared.Response) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, err *shared.AppError) {
	w.Header().Set("Content-Type", "application/json")
	response := errorResponse(err)
	json.NewEncoder(w).Encode(response)
}
