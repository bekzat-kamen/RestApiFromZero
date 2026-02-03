package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bekzat-kamen/RestApiFromZero/internal/models"
	"github.com/jmoiron/sqlx"
)

type TaskStore struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) *TaskStore {
	return &TaskStore{
		db: db,
	}
}

func (s *TaskStore) GetAll() ([]models.Task, error) {
	var tasks []models.Task

	query := `
		SELECT
			id,
			title,
			description,
			completed,
			created_at,
			updated_at
		FROM tasks
		ORDER BY created_at DESC
	`

	if err := s.db.Select(&tasks, query); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskStore) GetByID(id int) (*models.Task, error) {
	var task models.Task

	query := `
		SELECT
			id,
			title,
			description,
			completed,
			created_at,
			updated_at
		FROM tasks
		WHERE id = $1
	`

	err := s.db.Get(&task, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found with id %d", id)
		}
		return nil, err
	}

	return &task, nil
}
func (s *TaskStore) Create(input models.CreateTaskInput) (*models.Task, error) {
	var task models.Task
	now := time.Now()

	query := `
		INSERT INTO tasks (
			title,
			description,
			completed,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			title,
			description,
			completed,
			created_at,
			updated_at
	`

	err := s.db.
		QueryRowx(
			query,
			input.Title,
			input.Description,
			input.Completed,
			now,
			now,
		).
		StructScan(&task)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskStore) Update(id int, input models.UpdateTaskInput) (*models.Task, error) {
	existingTask, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		existingTask.Title = *input.Title
	}
	if input.Description != nil {
		existingTask.Description = *input.Description
	}
	if input.Completed != nil {
		existingTask.Completed = *input.Completed
	}

	existingTask.UpdatedAt = time.Now()

	query := `
		UPDATE tasks
		SET
			title = $1,
			description = $2,
			completed = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING
			id,
			title,
			description,
			completed,
			created_at,
			updated_at
	`

	var updatedTask models.Task
	err = s.db.
		QueryRowx(
			query,
			existingTask.Title,
			existingTask.Description,
			existingTask.Completed,
			existingTask.UpdatedAt,
			id,
		).
		StructScan(&updatedTask)

	if err != nil {
		return nil, err
	}

	return &updatedTask, nil
}

func (s *TaskStore) Delete(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found with id %d", id)
	}

	return nil
}
