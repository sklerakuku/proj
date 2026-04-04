package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/sklerakuku/tracker-web/internal/middleware"
	"github.com/sklerakuku/tracker-web/internal/model"
	"github.com/sklerakuku/tracker-web/internal/service"
	pkgjwt "github.com/sklerakuku/tracker-web/pkg/jwt"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, errMsg string) {
	h.respondJSON(w, status, ErrorResponse{Error: errMsg})
}

// Auth handlers
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	user, err := h.service.Register(r.Context(), req.Username, req.Password, req.Role)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	user, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		h.respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := pkgjwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// Template handlers
func (h *Handler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string               `json:"name"`
		Description string               `json:"description"`
		Tasks       []model.TaskTemplate `json:"tasks"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	template, err := h.service.CreateTemplate(r.Context(), req.Name, req.Description, req.Tasks)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, template)
}

func (h *Handler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.service.GetAllTemplates(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, templates)
}

func (h *Handler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/templates/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid template id")
		return
	}

	template, err := h.service.GetTemplate(r.Context(), id)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, template)
}

// Process handlers
func (h *Handler) CreateProcess(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TemplateID int    `json:"template_id"`
		Title      string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	process, err := h.service.CreateProcessFromTemplate(r.Context(), req.TemplateID, req.Title)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, process)
}

func (h *Handler) GetProcess(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/processes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid process id")
		return
	}

	process, err := h.service.GetProcess(r.Context(), id)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, process)
}

// Task handlers
func (h *Handler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	// extract task id from /tasks/{id}/status
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		h.respondError(w, http.StatusBadRequest, "invalid path")
		return
	}

	taskID, err := strconv.Atoi(parts[2])
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = h.service.UpdateTaskStatus(r.Context(), taskID, req.Status)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, SuccessResponse{Message: "task status updated"})
}

// Protected test endpoint
func (h *Handler) ProtectedTest(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "You are authenticated!",
		"user_id": claims.UserID,
		"role":    claims.Role,
	})
}
