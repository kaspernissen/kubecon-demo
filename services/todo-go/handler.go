package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func AllTodos(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "AllTodos")
	attr := attribute.KeyValue{Key: "code.function", Value: attribute.StringValue("AllTodos")}
	span.SetAttributes(attr)

	todos, err := getAllTodos(ctx)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to fetch todos")
		span.RecordError(err)
		span.End()
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		Logger.Error().Err(err).Msg("Failed to fetch todos")
		return
	}
	span.End()
	Logger.Info().Msg("AllTodos: request received")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	todo, err := getTodo(r.Context(), id)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		Logger.Error().Err(err).Msgf("Todo not found: %s", id)
		return
	}

	Logger.Info().Msg("GetTodo: request received")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		Logger.Error().Err(err).Msg("Invalid JSON payload")
		return
	}

	t, err := createTodo(r.Context(), todo)
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		Logger.Error().Err(err).Msg("Failed to create todo")
		return
	}

	Logger.Info().Msg("CreateTodo: request received")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := deleteTodo(r.Context(), id)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		Logger.Error().Err(err).Msgf("Failed to delete todo: %s", id)
		return
	}

	Logger.Info().Msg("DeleteTodo: request received")
	w.WriteHeader(http.StatusNoContent)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		Logger.Error().Err(err).Msg("Invalid JSON payload")
		return
	}

	t, err := updateTodo(r.Context(), updatedTodo.Name, id)
	if err != nil {
		if err.Error() == "todo not found" {
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		}
		Logger.Error().Err(err).Msgf("Failed to update todo: %s", id)
		return
	}

	Logger.Info().Msgf("Updated todo: %s -> %s", id, updatedTodo.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}
