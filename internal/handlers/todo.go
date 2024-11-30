package handlers

import (
	"net/http"
	"strconv"
	"time"

	"todo/internal/db"
	"todo/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new todo
// @Description Create a new todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body models.Todo true "Todo object"
// @Success 201 {object} models.Todo
// @Failure 400 {object} models.ErrorResponse
// @Router /todos [post]
func CreateTodo(ctx *gin.Context) {
	var newTodo models.Todo
	if err := ctx.ShouldBindJSON(&newTodo); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	query := `
		INSERT INTO todos (title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	timestamp := time.Now()
	err := db.DB.QueryRow(query,
		newTodo.Title,
		newTodo.Description,
		newTodo.Completed,
		timestamp,
		timestamp,
	).Scan(&newTodo.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	newTodo.CreatedAt = timestamp
	newTodo.UpdatedAt = timestamp

	ctx.JSON(http.StatusCreated, newTodo)
}

// @Summary Get all todos
// @Description Get a list of all todos
// @Tags todos
// @Produce json
// @Success 200 {array} models.Todo
// @Failure 500 {object} models.ErrorResponse
// @Router /todos [get]
func GetTodos(ctx *gin.Context) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at 
		FROM todos`

	rows, err := db.DB.Query(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	defer rows.Close()

	var todoList []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		); err != nil {
			ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		todoList = append(todoList, todo)
	}

	ctx.JSON(http.StatusOK, todoList)
}

// @Summary Get a todo
// @Description Get a todo by its ID
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /todos/{id} [get]
func GetTodo(ctx *gin.Context) {
	todoID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
		return
	}

	var todo models.Todo
	query := `
		SELECT id, title, description, completed, created_at, updated_at 
		FROM todos WHERE id = $1`

	err = db.DB.QueryRow(query, todoID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Todo not found foo"})
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

// @Summary Update a todo
// @Description Update a todo by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body models.Todo true "Todo object"
// @Success 200 {object} models.Todo
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /todos/{id} [put]
func UpdateTodo(ctx *gin.Context) {
	todoID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
		return
	}

	var updatedTodo models.Todo
	if err := ctx.ShouldBindJSON(&updatedTodo); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	query := `
		UPDATE todos 
		SET title = $1, description = $2, completed = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, created_at`

	timestamp := time.Now()
	err = db.DB.QueryRow(query,
		updatedTodo.Title,
		updatedTodo.Description,
		updatedTodo.Completed,
		timestamp,
		todoID,
	).Scan(&updatedTodo.ID, &updatedTodo.CreatedAt)

	if err != nil {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Todo not found"})
		return
	}

	updatedTodo.UpdatedAt = timestamp

	ctx.JSON(http.StatusOK, updatedTodo)
}

// @Summary Delete a todo
// @Description Delete a todo by its ID
// @Tags todos
// @Param id path int true "Todo ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /todos/{id} [delete]
func DeleteTodo(ctx *gin.Context) {
	todoID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid ID format"})
		return
	}

	query := `DELETE FROM todos WHERE id = $1`
	result, err := db.DB.Exec(query, todoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Todo not found"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
