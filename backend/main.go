package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cursor-react-go/backend/handlers"
	"github.com/cursor-react-go/backend/repository"
	"github.com/dgraph-io/badger/v3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func initDB() (*badger.DB, error) {
	// Create a data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, err
	}

	// Open the Badger database with minimal logging
	opts := badger.DefaultOptions("data/todos.db").
		WithLogger(nil)

	return badger.Open(opts)
}

func main() {
	// Initialize BadgerDB
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repository and handler
	todoRepo := repository.NewBadgerTodoRepository(db)
	todoHandler := handlers.NewTodoHandler(todoRepo)

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
	}))

	// Routes
	api := e.Group("/api")
	{
		todos := api.Group("/todos")
		todos.GET("", todoHandler.GetTodos)
		todos.POST("", todoHandler.CreateTodo)
		todos.PUT("/:id", todoHandler.UpdateTodo)
		todos.DELETE("/:id", todoHandler.DeleteTodo)
	}

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
