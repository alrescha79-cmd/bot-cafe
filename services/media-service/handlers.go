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
		response = h.createMedia(req.Payload)
	case "read":
		response = h.getMedia(req.Payload)
	case "list":
		response = h.listMedia(req.Payload)
	case "delete":
		response = h.deleteMedia(req.Payload)
	default:
		sendErrorResponse(w, shared.NewError(shared.ErrCodeInvalidInput, "Unknown action", nil))
		return
	}

	sendResponse(w, response)
}

// createMedia creates a new media record
func (h *Handler) createMedia(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	fileName, _ := data["file_name"].(string)
	fileURL, _ := data["file_url"].(string)
	fileType, _ := data["file_type"].(string)
	entityID, _ := data["entity_id"].(float64)
	entityType, _ := data["entity_type"].(string)

	// Validate inputs
	if err := shared.ValidateNotEmpty(fileName, "Nama file"); err != nil {
		return errorResponse(err.(*shared.AppError))
	}
	if err := shared.ValidateNotEmpty(fileURL, "URL file"); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	media := &Media{
		FileName:   fileName,
		FileURL:    fileURL,
		FileType:   fileType,
		EntityID:   int(entityID),
		EntityType: entityType,
	}

	result, err := h.repo.CreateMedia(media)
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"media": result,
	})
}

// getMedia gets a media by ID
func (h *Handler) getMedia(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	id, ok := data["id"].(float64)
	if !ok {
		return errorResponse(shared.NewInvalidInputError("ID diperlukan"))
	}

	media, err := h.repo.GetMediaByID(int(id))
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"media": media,
	})
}

// listMedia lists media by entity
func (h *Handler) listMedia(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	entityID, _ := data["entity_id"].(float64)
	entityType, _ := data["entity_type"].(string)

	medias, err := h.repo.ListMediaByEntity(int(entityID), entityType)
	if err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"medias": medias,
	})
}

// deleteMedia deletes a media
func (h *Handler) deleteMedia(payload interface{}) *shared.Response {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return errorResponse(shared.NewInvalidInputError("Invalid payload"))
	}

	id, ok := data["id"].(float64)
	if !ok {
		return errorResponse(shared.NewInvalidInputError("ID diperlukan"))
	}

	if err := h.repo.DeleteMedia(int(id)); err != nil {
		return errorResponse(err.(*shared.AppError))
	}

	return successResponse(map[string]interface{}{
		"message": "Media berhasil dihapus",
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
