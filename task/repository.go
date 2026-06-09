package task

import (
	"database/sql"
)

// Repository manages the database persistence for the Task entity.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(title string) (*Task, error) {
	var t Task
	query := "INSERT INTO tasks (title) VALUES ($1) RETURNING id, finished, created_at"
	err := r.db.QueryRow(query, title).Scan(&t.ID, &t.Finished, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	t.Title = title
	return &t, nil
}

func (r *Repository) List() ([]Task, error) {
	query := "SELECT id, title, finished, created_at FROM tasks ORDER BY id ASC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Finished, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *Repository) Update(id int, title string, finished bool) (int64, error) {
	query := "UPDATE tasks SET title = $1, finished = $2 WHERE id = $3"
	result, err := r.db.Exec(query, title, finished, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *Repository) Delete(id int) (int64, error) {
	query := "DELETE FROM tasks WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
