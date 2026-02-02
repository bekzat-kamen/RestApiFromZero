package models

import "time"

type Task struct {
	ID          int       `json:"id" db: "id"`
	Title       string    `json:"title" db: "title"`
	Description string    `json:"description" db: "description"`
	Completed   bool      `json:"completed" db: "completed"`
	CreatedAt   time.Time `json:"createdAt" db: "createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" db: "updatedAt"`
}

type CreateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type UpdateTaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}
