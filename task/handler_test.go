package task_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo_api/task"
	"todo_api/worker"
)

type mockRepository struct {
	createFunc func(title string) (*task.Task, error)
	listFunc   func() ([]task.Task, error)
	updateFunc func(id int, title string, finished bool) (int64, error)
	deleteFunc func(id int) (int64, error)
}

func (m *mockRepository) Create(title string) (*task.Task, error) {
	return m.createFunc(title)
}

func (m *mockRepository) List() ([]task.Task, error) {
	return m.listFunc()
}

func (m *mockRepository) Update(id int, title string, finished bool) (int64, error) {
	return m.updateFunc(id, title, finished)
}

func (m *mockRepository) Delete(id int) (int64, error) {
	return m.deleteFunc(id)
}

func TestCreateHandler(t *testing.T) {
	repo := &mockRepository{
		createFunc: func(title string) (*task.Task, error) {
			return &task.Task{ID: 1, Title: title}, nil
		},
	}
	w := worker.NewWorker(nil)
	w.Start()
	defer w.Stop()

	h := task.NewHandler(repo, w)

	payload := []byte(`{"title": "Test unit"}`)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(payload))
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}
}

func TestListHandler(t *testing.T) {
	repo := &mockRepository{
		listFunc: func() ([]task.Task, error) {
			return []task.Task{
				{ID: 1, Title: "Task 1", Finished: false},
				{ID: 2, Title: "Task 2", Finished: true},
			}, nil
		},
	}
	w := worker.NewWorker(nil)
	w.Start()
	defer w.Stop()

	h := task.NewHandler(repo, w)

	req := httptest.NewRequest("GET", "/tasks", nil)
	rec := httptest.NewRecorder()

	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var list []task.Task
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}

	if len(list) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(list))
	}
}

func TestUpdateHandler(t *testing.T) {
	repo := &mockRepository{
		updateFunc: func(id int, title string, finished bool) (int64, error) {
			if id == 1 {
				return 1, nil
			}
			return 0, nil
		},
	}
	w := worker.NewWorker(nil)
	w.Start()
	defer w.Stop()

	h := task.NewHandler(repo, w)

	payload := []byte(`{"title": "Updated Task", "finished": true}`)
	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(payload))
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	h.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestDeleteHandler(t *testing.T) {
	repo := &mockRepository{
		deleteFunc: func(id int) (int64, error) {
			if id == 1 {
				return 1, nil
			}
			return 0, nil
		},
	}
	w := worker.NewWorker(nil)
	w.Start()
	defer w.Stop()

	h := task.NewHandler(repo, w)

	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	h.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rec.Code)
	}
}
