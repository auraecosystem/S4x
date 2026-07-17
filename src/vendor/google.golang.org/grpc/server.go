
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	ErrConfigNotFound = errors.New("config not found")
	ErrDatabaseOffline = errors.New("database offline")
)

type Config struct {
	Address   string
	DB        string
	ClientURL string
}

type DB struct{}

func (d *DB) Close() error {
	log.Println("Database closed")
	return nil
}

type Client struct {
	URL string
}

func (c *Client) DoSomething(ctx context.Context) error {
	log.Printf("Connected to client: %s\n", c.URL)
	return nil
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := run(ctx); err != nil {
		log.Printf("Fatal: %v", err)
		os.Exit(exitCode(err))
	}
}

func run(ctx context.Context) error {

	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	db, err := openDB(cfg.DB)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	client := NewClient(cfg.ClientURL)

	if err := client.DoSomething(ctx); err != nil {
		return fmt.Errorf("client error: %w", err)
	}

	return startHTTPServer(ctx, cfg)
}

func loadConfig() (*Config, error) {

	return &Config{
		Address:   ":8080",
		DB:        "localhost",
		ClientURL: "https://api.example.com",
	}, nil
}

func openDB(addr string) (*DB, error) {

	if addr == "" {
		return nil, ErrDatabaseOffline
	}

	log.Println("Database connected")

	return &DB{}, nil
}

func NewClient(url string) *Client {
	return &Client{
		URL: url,
	}
}

func startHTTPServer(ctx context.Context, cfg *Config) error {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("GET /api/user", handleGetUser)

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	errCh := make(chan error, 1)

	go func() {
		log.Printf("Server listening on http://localhost%s", cfg.Address)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}

		errCh <- nil
	}()

	select {

	case <-ctx.Done():

		log.Println("Shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			15*time.Second,
		)
		defer cancel()

		return server.Shutdown(shutdownCtx)

	case err := <-errCh:

		return err
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "Welcome to the home endpoint!")
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}

	user := User{
		ID:   101,
		Name: name,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
	}
}

func exitCode(err error) int {

	switch {

	case errors.Is(err, ErrConfigNotFound):
		return 2

	case errors.Is(err, ErrDatabaseOffline):
		return 3

	default:
		return 1
	}
}
