package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTodo(t *testing.T) {
	todo := Todo{Title: "foo title", Description: "foo description"}
	assert.Equal(t, "foo title", todo.Title)
	assert.Equal(t, "foo description", todo.Description)
	assert.Equal(t, "foo title", todo.Title)
	assert.Equal(t, "foo description", todo.Description)
	assert.Equal(t, false, todo.Completed)
	assert.Zero(t, todo.CreatedAt)
	assert.Zero(t, todo.UpdatedAt)
}
