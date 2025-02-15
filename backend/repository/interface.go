package repository

// Todo represents a todo item
type Todo struct {
	ID        string
	Title     string
	Completed bool
}

// TodoRepository defines the interface for todo storage operations
type TodoRepository interface {
	GetAll() ([]Todo, error)
	GetByID(id string) (Todo, error)
	Create(todo Todo) error
	Update(todo Todo) error
	Delete(id string) error
}
