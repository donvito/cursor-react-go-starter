package repository

import (
	"encoding/json"
	"strings"

	"github.com/dgraph-io/badger/v3"
)

// BadgerTodoRepository implements TodoRepository using BadgerDB
type BadgerTodoRepository struct {
	db *badger.DB
}

// NewBadgerTodoRepository creates a new BadgerDB-backed todo repository
func NewBadgerTodoRepository(db *badger.DB) TodoRepository {
	return &BadgerTodoRepository{db: db}
}

// validate checks if a todo item is valid
func (t Todo) validate() error {
	if strings.TrimSpace(t.ID) == "" {
		return &ErrInvalidInput{Message: "id cannot be empty"}
	}
	if strings.TrimSpace(t.Title) == "" {
		return &ErrInvalidInput{Message: "title cannot be empty"}
	}
	return nil
}

func (r *BadgerTodoRepository) GetAll() ([]Todo, error) {
	var todos []Todo
	err := r.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			var todo Todo
			err := it.Item().Value(func(val []byte) error {
				return json.Unmarshal(val, &todo)
			})
			if err != nil {
				return &ErrDatabase{Op: "GetAll.Unmarshal", Err: err}
			}
			todos = append(todos, todo)
		}
		return nil
	})
	if err != nil {
		return nil, &ErrDatabase{Op: "GetAll", Err: err}
	}
	return todos, nil
}

func (r *BadgerTodoRepository) GetByID(id string) (Todo, error) {
	if strings.TrimSpace(id) == "" {
		return Todo{}, &ErrInvalidInput{Message: "id cannot be empty"}
	}

	var todo Todo
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err == badger.ErrKeyNotFound {
			return &ErrNotFound{ID: id}
		}
		if err != nil {
			return &ErrDatabase{Op: "GetByID", Err: err}
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &todo)
		})
	})
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func (r *BadgerTodoRepository) Create(todo Todo) error {
	if err := todo.validate(); err != nil {
		return err
	}

	value, err := json.Marshal(todo)
	if err != nil {
		return &ErrDatabase{Op: "Create.Marshal", Err: err}
	}

	return r.db.Update(func(txn *badger.Txn) error {
		// Check if todo already exists
		_, err := txn.Get([]byte(todo.ID))
		if err == nil {
			return &ErrInvalidInput{Message: "todo with this ID already exists"}
		}
		if err != badger.ErrKeyNotFound {
			return &ErrDatabase{Op: "Create.Check", Err: err}
		}

		return txn.Set([]byte(todo.ID), value)
	})
}

func (r *BadgerTodoRepository) Update(todo Todo) error {
	if err := todo.validate(); err != nil {
		return err
	}

	value, err := json.Marshal(todo)
	if err != nil {
		return &ErrDatabase{Op: "Update.Marshal", Err: err}
	}

	return r.db.Update(func(txn *badger.Txn) error {
		// Check if todo exists
		_, err := txn.Get([]byte(todo.ID))
		if err == badger.ErrKeyNotFound {
			return &ErrNotFound{ID: todo.ID}
		}
		if err != nil {
			return &ErrDatabase{Op: "Update.Check", Err: err}
		}

		return txn.Set([]byte(todo.ID), value)
	})
}

func (r *BadgerTodoRepository) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return &ErrInvalidInput{Message: "id cannot be empty"}
	}

	return r.db.Update(func(txn *badger.Txn) error {
		// Check if todo exists
		_, err := txn.Get([]byte(id))
		if err == badger.ErrKeyNotFound {
			return &ErrNotFound{ID: id}
		}
		if err != nil {
			return &ErrDatabase{Op: "Delete.Check", Err: err}
		}

		return txn.Delete([]byte(id))
	})
}
