package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	// Import the official pgx driver for PostgreSQL
	_ "://github.com"
)

// Global DB connection pool
var db *sql.DB

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// 1. Initialize Database Connection
	initDB()
	defer db.Close()

	mux := http.NewServeMux()

	// 2. Static Asset Routes (Serves index.html, CSS, JS from a "public" folder)
	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// 3. API Routes (Wrapped in CORS and Authentication Middleware)
	mux.Handle("GET /api/user", corsMiddleware(authMiddleware(http.HandlerFunc(handleGetUser))))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 4. Graceful Shutdown Routine
	go func() {
		log.Printf("Server running at http://localhost%s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	<-shutdownChan

	log.Println("Shutting down cleanly...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}
	log.Println("Server stopped.")
}

// initDB handles the secure PostgreSQL connection pool setup
func initDB() {
	// Use environment variables for connection strings in production
	connStr := "postgres://username:password@localhost:5432/dbname?sslmode=disable"
	var err error
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Database connection configuration failed: %v", err)
	}

	// Set connection pool limits
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		log.Fatalf("Could not reach database: %v", err)
	}
	log.Println("Database connection successfully established.")
}

// corsMiddleware injects global browser permission headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace * with your exact front-end domain in production
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight browser checks instantly
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// authMiddleware validates structural Bearer tokens
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing or malformed token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// NOTE: In production, pass 'tokenString' to a JWT validation library 
		// (e.g., ://github.com) to verify the signature.
		if tokenString != "valid-mock-token" { 
			http.Error(w, "Unauthorized: Invalid token signature", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleGetUser fetches data using safe parameterized SQL queries
func handleGetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "Bad Request: Missing user ID query parameter", http.StatusBadRequest)
		return
	}

	var user User
	// Parameterized queries ($1) completely protect against SQL Injection
	query := "SELECT id, name FROM users WHERE id = $1 LIMIT 1"
	err := db.QueryRowContext(r.Context(), query, userID).Scan(&user.ID, &user.Name)
	
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal Database Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
