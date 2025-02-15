package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE
	);`

	if _, err := db.Exec(createTable); err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}

	return db
}

func TestSQLiteTodoRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteTodoRepository(db)

	tests := []struct {
		name    string
		todo    Todo
		wantErr bool
		errType error
	}{
		{
			name: "valid todo",
			todo: Todo{
				ID:        "1",
				Title:     "Test Todo",
				Completed: false,
			},
			wantErr: false,
		},
		{
			name: "empty id",
			todo: Todo{
				ID:        "",
				Title:     "Test Todo",
				Completed: false,
			},
			wantErr: true,
			errType: &ErrInvalidInput{},
		},
		{
			name: "empty title",
			todo: Todo{
				ID:        "2",
				Title:     "",
				Completed: false,
			},
			wantErr: true,
			errType: &ErrInvalidInput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.todo)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errType != nil {
				if _, ok := err.(interface{ Error() string }); !ok {
					t.Errorf("Create() error = %v, want error type %T", err, tt.errType)
				}
			}
		})
	}
}

func TestSQLiteTodoRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteTodoRepository(db)

	// Create a test todo
	testTodo := Todo{
		ID:        "1",
		Title:     "Test Todo",
		Completed: false,
	}
	if err := repo.Create(testTodo); err != nil {
		t.Fatalf("failed to create test todo: %v", err)
	}

	tests := []struct {
		name    string
		id      string
		want    Todo
		wantErr bool
		errType error
	}{
		{
			name: "existing todo",
			id:   "1",
			want: testTodo,
		},
		{
			name:    "non-existing todo",
			id:      "999",
			wantErr: true,
			errType: &ErrNotFound{},
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: true,
			errType: &ErrInvalidInput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GetByID() got = %v, want %v", got, tt.want)
			}
			if tt.wantErr && tt.errType != nil {
				if _, ok := err.(interface{ Error() string }); !ok {
					t.Errorf("GetByID() error = %v, want error type %T", err, tt.errType)
				}
			}
		})
	}
}

func TestSQLiteTodoRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteTodoRepository(db)

	// Create a test todo
	testTodo := Todo{
		ID:        "1",
		Title:     "Test Todo",
		Completed: false,
	}
	if err := repo.Create(testTodo); err != nil {
		t.Fatalf("failed to create test todo: %v", err)
	}

	tests := []struct {
		name    string
		todo    Todo
		wantErr bool
		errType error
	}{
		{
			name: "valid update",
			todo: Todo{
				ID:        "1",
				Title:     "Updated Todo",
				Completed: true,
			},
			wantErr: false,
		},
		{
			name: "non-existing todo",
			todo: Todo{
				ID:        "999",
				Title:     "Non-existing Todo",
				Completed: false,
			},
			wantErr: true,
			errType: &ErrNotFound{},
		},
		{
			name: "empty title",
			todo: Todo{
				ID:        "1",
				Title:     "",
				Completed: true,
			},
			wantErr: true,
			errType: &ErrInvalidInput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(tt.todo)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errType != nil {
				if _, ok := err.(interface{ Error() string }); !ok {
					t.Errorf("Update() error = %v, want error type %T", err, tt.errType)
				}
			}
		})
	}
}

func TestSQLiteTodoRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteTodoRepository(db)

	// Create a test todo
	testTodo := Todo{
		ID:        "1",
		Title:     "Test Todo",
		Completed: false,
	}
	if err := repo.Create(testTodo); err != nil {
		t.Fatalf("failed to create test todo: %v", err)
	}

	tests := []struct {
		name    string
		id      string
		wantErr bool
		errType error
	}{
		{
			name:    "existing todo",
			id:      "1",
			wantErr: false,
		},
		{
			name:    "non-existing todo",
			id:      "999",
			wantErr: true,
			errType: &ErrNotFound{},
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: true,
			errType: &ErrInvalidInput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errType != nil {
				if _, ok := err.(interface{ Error() string }); !ok {
					t.Errorf("Delete() error = %v, want error type %T", err, tt.errType)
				}
			}
		})
	}
}
