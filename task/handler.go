package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
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
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := input.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdTask, err := h.repo.Create(input.Title)
	if err != nil {
		http.Error(w, "Error inserting task in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.repo.List()
	if err != nil {
		http.Error(w, "Failed fetching data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var input UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := input.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rowsAffected, err := h.repo.Update(id, input.Title, input.Finished)
	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task updated successfully"))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	rowsAffected, err := h.repo.Delete(id)
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
