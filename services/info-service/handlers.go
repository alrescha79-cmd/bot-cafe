package main

import (
	"encoding/json"
	"net/http"

	"github.com/alrescha79-cmd/bot-cafe/shared"
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
	case "read":
		response = h.getCafeInfo()
	case "update":
		response = h.updateCafeInfo(req.Payload)
	default:
		sendErrorResponse(w, shared.NewError(shared.ErrCodeInvalidInput, "Unknown action", nil))
		return
	}

	sendResponse(w, response)
}

// getCafeInfo gets café information

func (h *Handler) getCafeInfo() *shared.Response {
	info, err := h.repo.GetCafeInfo()
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"info": info,
	})
}

// updateCafeInfo updates café information
func (h *Handler) updateCafeInfo(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	// Get existing info
	info, err := h.repo.GetCafeInfo()
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	// Update fields if provided
	if name, ok := data["name"].(string); ok && name != "" {
		info.Name = shared.SanitizeInput(name)
	}
	if address, ok := data["address"].(string); ok && address != "" {
		info.Address = shared.SanitizeInput(address)
	}
	if phone, ok := data["phone"].(string); ok && phone != "" {
		info.Phone = shared.SanitizeInput(phone)
	}
	if email, ok := data["email"].(string); ok {
		info.Email = shared.SanitizeInput(email)
	}
	if openingHour, ok := data["opening_hour"].(string); ok && openingHour != "" {
		info.OpeningHour = openingHour
	}
	if closingHour, ok := data["closing_hour"].(string); ok && closingHour != "" {
		info.ClosingHour = closingHour
	}
	if description, ok := data["description"].(string); ok {
		info.Description = shared.SanitizeInput(description)
	}

	if err := h.repo.UpdateCafeInfo(info); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"info": info,
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
