package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"todo_api/apierror"
)

// TaskRepository defines the database operations required for tasks.
type TaskRepository interface {
	Create(title string) (*Task, error)
	List() ([]Task, error)
	Update(id int, title string, finished bool) (int64, error)
	Delete(id int) (int64, error)
}

// Handler manages the HTTP transport for Task-related operations.
type Handler struct {
	repo TaskRepository
}

func NewHandler(repo TaskRepository) *Handler {
	return &Handler{repo: repo}
}

func validateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("title is required and cannot be empty")
	}
	return nil
}

type CreateInput struct {
	Title string `json:"title"`
}

func (i CreateInput) Validate() error {
	return validateTitle(i.Title)
}

type UpdateInput struct {
	Title    string `json:"title"`
	Finished bool   `json:"finished"`
}

func (i UpdateInput) Validate() error {
	return validateTitle(i.Title)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		apierror.Write(w, r, apierror.BadRequest("invalid JSON", err))
		return
	}

	if err := input.Validate(); err != nil {
		apierror.Write(w, r, apierror.BadRequest(err.Error(), nil))
		return
	}

	createdTask, err := h.repo.Create(input.Title)
	if err != nil {
		apierror.Write(w, r, apierror.Internal("Error inserting task in database", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.repo.List()
	if err != nil {
		apierror.Write(w, r, apierror.Internal("Failed fetching data", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apierror.Write(w, r, apierror.BadRequest("Invalid ID", err))
		return
	}

	var input UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.Write(w, r, apierror.BadRequest("Invalid JSON", err))
		return
	}

	if err := input.Validate(); err != nil {
		apierror.Write(w, r, apierror.BadRequest(err.Error(), nil))
		return
	}

	rowsAffected, err := h.repo.Update(id, input.Title, input.Finished)
	if err != nil {
		apierror.Write(w, r, apierror.Internal("Update failed", err))
		return
	}

	if rowsAffected == 0 {
		apierror.Write(w, r, apierror.NotFound("Task not found", nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task updated successfully"))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apierror.Write(w, r, apierror.BadRequest("Invalid ID", err))
		return
	}

	rowsAffected, err := h.repo.Delete(id)
	if err != nil {
		apierror.Write(w, r, apierror.Internal("Delete failed", err))
		return
	}

	if rowsAffected == 0 {
		apierror.Write(w, r, apierror.NotFound("Task not found", nil))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
