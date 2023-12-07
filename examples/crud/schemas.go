package main

import (
	"time"

	"github.com/gofrs/uuid"
)

type Error struct {
	Message string `json:"message,omitempty"`
}

type Todo struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description,omitempty"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Completed   bool      `json:"completed,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
}

type NewTodo struct {
	Description string   `json:"description"`
	Tags        []string `json:"tags,omitempty"`
}

type UpdateTodo struct {
	Description string   `json:"description,omitempty"`
	Completed   bool     `json:"completed,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type GetTodosQuery struct {
	Search        string   `query:"search"`
	ShowCompleted bool     `query:"completed"`
	Tags          []string `query:"tags"`
}
