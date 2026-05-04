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

// LoginRequest структура запроса для входа
type LoginRequest struct {
	Username string `json:"username" example:"ivanov"`
	Password string `json:"password" example:"password123"`
}

// RegisterRequest структура запроса для регистрации
type RegisterRequest struct {
	Username string `json:"username" example:"ivanov"`
	Password string `json:"password" example:"password123"`
	Role     string `json:"role" example:"worker"`
}

// LoginResponse структура ответа при успешном входе
type LoginResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

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

// Register godoc
// @Summary Регистрация пользователя
// @Description Создание нового пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} model.User
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
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

// Login godoc
// @Summary Авторизация пользователя
// @Description Вход в систему, возвращает JWT-токен
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Учетные данные"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
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
// CreateTemplate godoc
// @Summary Создание шаблона
// @Description Создание нового шаблона процесса с задачами
// @Tags Templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTemplateRequest true "Данные шаблона"
// @Success 201 {object} model.Template
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /templates [post]
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

// ListTemplates godoc
// @Summary Получение списка шаблонов
// @Description Возвращает все шаблоны процессов
// @Tags Templates
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Template
// @Failure 500 {object} ErrorResponse
// @Router /templates [get]
func (h *Handler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.service.GetAllTemplates(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, templates)
}

// GetTemplate godoc
// @Summary Получение шаблона по ID
// @Description Возвращает шаблон с задачами
// @Tags Templates
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID шаблона"
// @Success 200 {object} model.Template
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /templates/{id} [get]
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
// CreateProcess godoc
// @Summary Создание процесса из шаблона
// @Description Создает экземпляр процесса на основе шаблона
// @Tags Processes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProcessRequest true "Данные процесса"
// @Success 201 {object} model.Process
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /processes [post]
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

// CreateEmptyProcess - создание процесса без шаблона
func (h *Handler) CreateEmptyProcess(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if req.Title == "" {
		h.respondError(w, http.StatusBadRequest, "title is required")
		return
	}

	process, err := h.service.CreateEmptyProcess(r.Context(), req.Title)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, process)
}

// GetProcess godoc
// @Summary Получение процесса по ID
// @Description Возвращает процесс со всеми задачами и зависимостями
// @Tags Processes
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID процесса"
// @Success 200 {object} model.Process
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /processes/{id} [get]
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
// UpdateTaskStatus godoc
// @Summary Обновление статуса задачи
// @Description Изменяет статус задачи (pending/in_progress/done)
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи"
// @Param request body UpdateStatusRequest true "Новый статус"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id}/status [patch]
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

type CreateTemplateRequest struct {
	Name        string               `json:"name" example:"Процесс найма"`
	Description string               `json:"description" example:"Стандартный процесс найма сотрудника"`
	Tasks       []model.TaskTemplate `json:"tasks"`
}

// CreateProcessRequest структура для создания процесса
type CreateProcessRequest struct {
	TemplateID int    `json:"template_id" example:"1"`
	Title      string `json:"title" example:"Найм Иванова"`
}

// UpdateStatusRequest структура для обновления статуса
type UpdateStatusRequest struct {
	Status string `json:"status" example:"done" enums:"pending,in_progress,done"`
}

// ListProcesses - список всех процессов
func (h *Handler) ListProcesses(w http.ResponseWriter, r *http.Request) {
	processes, err := h.service.GetAllProcesses(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, processes)
}

// АДМИНКА

// ListUsers - список пользователей
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, users)
}

// UpdateUser - изменить пользователя
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/admin/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req struct {
		Username string `json:"username"`
		Role     string `json:"role"`
		Password string `json:"password,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := h.service.UpdateUser(r.Context(), id, req.Username, req.Role, req.Password); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, SuccessResponse{Message: "user updated"})
}

// DeleteUser - удалить пользователя
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/admin/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, SuccessResponse{Message: "user deleted"})
}

// DeleteTemplate - удалить шаблон
func (h *Handler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/admin/templates/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid template id")
		return
	}

	if err := h.service.DeleteTemplate(r.Context(), id); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, SuccessResponse{Message: "template deleted"})
}

// DeleteProcess - удалить процесс
func (h *Handler) DeleteProcess(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/admin/processes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid process id")
		return
	}

	if err := h.service.DeleteProcess(r.Context(), id); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, SuccessResponse{Message: "process deleted"})
}
