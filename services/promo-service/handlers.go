package main

import (
	"encoding/json"
	"net/http"
	"time"

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
		response = h.createPromo(req.Payload)
	case "read":
		response = h.getPromo(req.Payload)
	case "update":
		response = h.updatePromo(req.Payload)
	case "delete":
		response = h.deletePromo(req.Payload)
	case "list":
		response = h.listPromos(req.Payload)
	default:
		sendErrorResponse(w, shared.NewError(shared.ErrCodeInvalidInput, "Unknown action", nil))
		return
	}

	sendResponse(w, response)
}

// createPromo creates a new promo
func (h *Handler) createPromo(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	title, _ := data["title"].(string)
	description, _ := data["description"].(string)
	discountRaw := data["discount"]
	discountType, _ := data["discount_type"].(string)
	startDateStr, _ := data["start_date"].(string)
	endDateStr, _ := data["end_date"].(string)
	isActive := true
	if val, ok := data["is_active"].(bool); ok {
		isActive = val
	}

	// Validate inputs
	if err := shared.ValidateNotEmpty(title, "Judul promo"); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	discount, err := shared.ValidatePrice(discountRaw)
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	if discountType != "percentage" && discountType != "amount" {
		return errorResponse(shared.NewInvalidInputError("Tipe diskon harus 'percentage' atau 'amount'"))
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return errorResponse(shared.NewInvalidInputError("Format tanggal mulai tidak valid (YYYY-MM-DD)"))
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return errorResponse(shared.NewInvalidInputError("Format tanggal akhir tidak valid (YYYY-MM-DD)"))
	}

	if endDate.Before(startDate) {
		return errorResponse(shared.NewInvalidInputError("Tanggal akhir harus setelah tanggal mulai"))
	}

	promo := &Promo{
		Title:        shared.SanitizeInput(title),
		Description:  shared.SanitizeInput(description),
		Discount:     discount,
		DiscountType: discountType,
		StartDate:    startDate,
		EndDate:      endDate,
		IsActive:     isActive,
	}

	result, err := h.repo.CreatePromo(promo)
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"promo": result,
	})
}

// getPromo gets a promo by ID
func (h *Handler) getPromo(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	id, ok := data["id"].(float64)
	if !ok {
		return errorResponse(shared.NewInvalidInputError("ID diperlukan"))
	}

	promo, err := h.repo.GetPromoByID(int(id))
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"promo": promo,
	})
}

// updatePromo updates a promo
func (h *Handler) updatePromo(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	id, ok := data["id"].(float64)
	if !ok {
		return errorResponse(shared.NewInvalidInputError("ID diperlukan"))
	}

	// Get existing promo
	promo, err := h.repo.GetPromoByID(int(id))
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	// Update fields if provided
	if title, ok := data["title"].(string); ok && title != "" {
		promo.Title = shared.SanitizeInput(title)
	}
	if desc, ok := data["description"].(string); ok {
		promo.Description = shared.SanitizeInput(desc)
	}
	if discountRaw, ok := data["discount"]; ok {
		discount, err := shared.ValidatePrice(discountRaw)
		if err != nil {
			return errorResponse(err.(*shared.AppError))
		}
		promo.Discount = discount
	}
	if discountType, ok := data["discount_type"].(string); ok && discountType != "" {
		if discountType != "percentage" && discountType != "amount" {
			return errorResponse(shared.NewInvalidInputError("Tipe diskon harus 'percentage' atau 'amount'"))
		}
		promo.DiscountType = discountType
	}
	if startDateStr, ok := data["start_date"].(string); ok && startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return errorResponse(shared.NewInvalidInputError("Format tanggal mulai tidak valid"))
		}
		promo.StartDate = startDate
	}
	if endDateStr, ok := data["end_date"].(string); ok && endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return errorResponse(shared.NewInvalidInputError("Format tanggal akhir tidak valid"))
		}
		promo.EndDate = endDate
	}
	if isActive, ok := data["is_active"].(bool); ok {
		promo.IsActive = isActive
	}

	if err := h.repo.UpdatePromo(promo); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"promo": promo,
	})
}

// deletePromo deletes a promo
func (h *Handler) deletePromo(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	id, ok := data["id"].(float64)
	if !ok {
		return errorResponse(shared.NewInvalidInputError("ID diperlukan"))
	}

	if err := h.repo.DeletePromo(int(id)); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"message": "Promo berhasil dihapus",
	})
}

// listPromos lists promos
func (h *Handler) listPromos(payload interface{}) *shared.Response {
	activeOnly := false

	if data, ok := payload.(map[string]interface{}); ok {
		if active, ok := data["active_only"].(bool); ok {
			activeOnly = active
		}
	}

	promos, err := h.repo.ListPromos(activeOnly)
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"promos": promos,
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
