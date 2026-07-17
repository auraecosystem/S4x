package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os_signal"
	"syscall"
	"time"
)

// User represents a sample data payload
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// 1. Create a secure request multiplexer (router)
	mux := http.NewServeMux()

	// 2. Define routes
	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("GET /api/user", handleGetUser)

	// 3. Configure server timeouts for production safety
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 4. Run the server in a separate goroutine to allow graceful shutdown
	go func() {
		log.Printf("Server starting on HTTP://localhost%s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 5. Channel to catch OS termination signals
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-shutdownChan
	log.Println("Shutdown signal received. Closing connections...")

	// 6. Context with a timeout to force shutdown if active requests hang
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped.")
}

// handleHome writes a plain text response
func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to the home endpoint!")
}

// handleGetUser parses query parameters and returns a JSON object
func handleGetUser(w http.ResponseWriter, r *http.Request) {
	// Parse query parameter (e.g., /api/user?name=Alice)
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}

	user := User{
		ID:   101,
		Name: name,
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
