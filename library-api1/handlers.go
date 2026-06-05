package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

// Структура для JSON-лога
type LogEntry struct {
	Timestamp string                 `json:"timestamp"` // Время события
	Level     string                 `json:"level"`     // Уровень (INFO, ERROR)
	Message   string                 `json:"message"`   // Что произошло
	Fields    map[string]interface{} `json:"fields"`    // Дополнительные данные (ID, названия и т.д.)
}

// Вспомогательная функция для отправки лога в консоль в формате JSON
func logJSON(level, message string, fields map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		Fields:    fields,
	}
	// Кодируем структуру в JSON и сразу выводим в стандартный поток вывода (консоль)
	json.NewEncoder(os.Stdout).Encode(entry)
}

// 1. ОБРАБОТЧИК КНИГ (/books)
func handleBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Добавление книги
		if r.Method == http.MethodPost {
			var b Book
			if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			createBook(db, &b)

			// JSON Лог
			logJSON("INFO", "book_created", map[string]interface{}{
				"book_id": b.ID,
				"title":   b.Title,
				"author":  b.Author,
			})

			json.NewEncoder(w).Encode(b)
			return
		}

		// Просмотр всех книг
		if r.Method == http.MethodGet {
			books, err := getBooks(db) // Принимаем и книги, и ошибку
			if err != nil {
				// Логируем ошибку в нашем новом JSON-формате
				logJSON("ERROR", "get_books_failed", map[string]interface{}{
					"error": err.Error(),
				})
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(books)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// 2. ОБРАБОТЧИК ПОЛЬЗОВАТЕЛЕЙ (/users)
func handleUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodPost {
			var u User
			if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			createUser(db, &u)

			// JSON Лог
			logJSON("INFO", "user_registered", map[string]interface{}{
				"user_id": u.ID,
				"name":    u.Name,
				"email":   u.Email,
			})

			json.NewEncoder(w).Encode(u)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// 3. ОБРАБОТЧИК ВЫДАЧИ КНИГ (/issues)
func handleIssues(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodPost {
			var req struct {
				UserID string `json:"user_id"`
				BookID string `json:"book_id"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err := issueBook(db, req.UserID, req.BookID)
			if err != nil {
				logJSON("ERROR", "book_issue_failed", map[string]interface{}{
					"user_id": req.UserID,
					"book_id": req.BookID,
					"error":   err.Error(),
				})
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// JSON Лог успешной выдачи
			logJSON("INFO", "book_issued", map[string]interface{}{
				"user_id": req.UserID,
				"book_id": req.BookID,
			})

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// 4. ОБРАБОТЧИК ВОЗВРАТА КНИГ (/returns)
func handleReturns(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodPost {
			var req struct {
				BookID string `json:"book_id"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err := returnBook(db, req.BookID)
			if err != nil {
				logJSON("ERROR", "book_return_failed", map[string]interface{}{
					"book_id": req.BookID,
					"error":   err.Error(),
				})
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// JSON Лог успешного возврата
			logJSON("INFO", "book_returned", map[string]interface{}{
				"book_id": req.BookID,
			})

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
