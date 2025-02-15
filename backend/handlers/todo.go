package handlers

import (
	"net/http"

	"github.com/cursor-react-go/backend/repository"
	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	repo repository.TodoRepository
}

type TodoResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewTodoHandler(repo repository.TodoRepository) *TodoHandler {
	return &TodoHandler{repo: repo}
}

func toResponse(todo repository.Todo) TodoResponse {
	return TodoResponse{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: todo.Completed,
	}
}

func handleError(err error) (int, ErrorResponse) {
	switch e := err.(type) {
	case *repository.ErrNotFound:
		return http.StatusNotFound, ErrorResponse{Error: e.Error()}
	case *repository.ErrInvalidInput:
		return http.StatusBadRequest, ErrorResponse{Error: e.Error()}
	case *repository.ErrDatabase:
		return http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"}
	default:
		return http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"}
	}
}

func (h *TodoHandler) GetTodos(c echo.Context) error {
	todos, err := h.repo.GetAll()
	if err != nil {
		status, resp := handleError(err)
		return c.JSON(status, resp)
	}

	var response []TodoResponse
	for _, todo := range todos {
		response = append(response, toResponse(todo))
	}

	return c.JSON(http.StatusOK, response)
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	var req TodoResponse
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
	}

	todo := repository.Todo{
		ID:        req.ID,
		Title:     req.Title,
		Completed: req.Completed,
	}

	if err := h.repo.Create(todo); err != nil {
		status, resp := handleError(err)
		return c.JSON(status, resp)
	}

	return c.JSON(http.StatusCreated, req)
}

func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	id := c.Param("id")
	var req TodoResponse
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
	}

	todo := repository.Todo{
		ID:        id,
		Title:     req.Title,
		Completed: req.Completed,
	}

	if err := h.repo.Update(todo); err != nil {
		status, resp := handleError(err)
		return c.JSON(status, resp)
	}

	req.ID = id
	return c.JSON(http.StatusOK, req)
}

func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	id := c.Param("id")
	if err := h.repo.Delete(id); err != nil {
		status, resp := handleError(err)
		return c.JSON(status, resp)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Todo deleted successfully"})
}
