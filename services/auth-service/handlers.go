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
	case "verify":
		response = h.verifyAdmin(req.Payload)
	case "login":
		response = h.login(req.Payload)
	case "logout":
		response = h.logout(req.Payload)
	case "list":
		response = h.listAdmins()
	case "register":
		response = h.registerAdmin(req.Payload)
	case "update_status":
		response = h.updateStatus(req.Payload)
	default:
		sendErrorResponse(w, shared.NewError(shared.ErrCodeInvalidInput, "Unknown action", nil))
		return
	}

	sendResponse(w, response)
}

// verifyAdmin verifies if user is admin
func (h *Handler) verifyAdmin(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	telegramID, ok := data["telegram_id"].(string)
	if !ok || telegramID == "" {
		return errorResponse(shared.NewInvalidInputError("telegram_id is required"))
	}

	admin, err := h.repo.GetAdminByTelegramID(telegramID)
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	if !admin.IsActive {
		return errorResponse(shared.NewUnauthorizedError())
	}

	return successResponse(map[string]interface{}{
		"is_admin": true,
		"admin":    admin,
	})
}

// login creates a session for admin
func (h *Handler) login(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	telegramID, ok := data["telegram_id"].(string)
	if !ok || telegramID == "" {
		return errorResponse(shared.NewInvalidInputError("telegram_id is required"))
	}

	admin, err := h.repo.GetAdminByTelegramID(telegramID)
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	if !admin.IsActive {
		return errorResponse(shared.NewUnauthorizedError())
	}

	session, err := h.repo.CreateSession(admin.ID)
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"admin":   admin,
		"session": session,
	})
}

// logout removes session
func (h *Handler) logout(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	token, ok := data["token"].(string)
	if !ok || token == "" {
		return errorResponse(shared.NewInvalidInputError("token is required"))
	}

	if err := h.repo.DeleteSession(token); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"message": "Logout berhasil",
	})
}

// listAdmins lists all admins
func (h *Handler) listAdmins() *shared.Response {
	admins, err := h.repo.ListAdmins()
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"admins": admins,
	})
}

// registerAdmin registers a new admin
func (h *Handler) registerAdmin(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	telegramID, _ := data["telegram_id"].(string)
	username, _ := data["username"].(string)

	if telegramID == "" || username == "" {
		return errorResponse(shared.NewInvalidInputError("telegram_id dan username diperlukan"))
	}

	admin, err := h.repo.CreateAdmin(telegramID, username)
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"admin": admin,
	})
}

// updateStatus updates admin status
func (h *Handler) updateStatus(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	telegramID, _ := data["telegram_id"].(string)
	isActive, _ := data["is_active"].(bool)

	if telegramID == "" {
		return errorResponse(shared.NewInvalidInputError("telegram_id diperlukan"))
	}

	if err := h.repo.UpdateAdminStatus(telegramID, isActive); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"message": "Status admin berhasil diperbarui",
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
