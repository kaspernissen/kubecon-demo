package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func getAllTodos(ctx context.Context) ([]Todo, error) {
	query := "SELECT id, name FROM todo"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.Id, &t.Name); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, rows.Err()
}

func getTodo(ctx context.Context, id string) (Todo, error) {
	query := "SELECT id, name FROM todo WHERE id = $1"
	var t Todo
	err := db.QueryRowContext(ctx, query, id).Scan(&t.Id, &t.Name)
	if err == sql.ErrNoRows {
		return Todo{}, err
	}
	return t, err
}

func updateTodo(ctx context.Context, name, id string) (Todo, error) {
	query := "UPDATE todo SET name = $1 WHERE id = $2 RETURNING id, name"
	var t Todo
	err := db.QueryRowContext(ctx, query, name, id).Scan(&t.Id, &t.Name)
	if err == sql.ErrNoRows {
		return Todo{}, fmt.Errorf("todo not found")
	}
	return t, err
}

func deleteTodo(ctx context.Context, id string) error {
	query := "DELETE FROM todo WHERE id = $1"
	_, err := db.ExecContext(ctx, query, id)
	return err
}

func createTodo(ctx context.Context, todo Todo) (Todo, error) {
	query := "INSERT INTO todo(name) VALUES ($1) RETURNING id"
	err := db.QueryRowContext(ctx, query, todo.Name).Scan(&todo.Id)
	return todo, err
}
