package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Get database configuration from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Default to localhost if env vars are missing (for local testing)
	if dbHost == "" {
		dbHost = "localhost"
		dbPort = "5432"
		dbUser = "postgres"
		dbPassword = "postgres"
		dbName = "pdfconverter"
	}

	// Construct connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// Initialize database
	var db *Database
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {

		db, err = NewDatabase(connStr)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to databse (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(6 * time.Second) // wait 6 seconds before retrying
	}

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Println("Database initialized successfully")

	// Create handler
	handler := NewHandler(db)

	// Setup routes
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/upload", handler.HandleUpload)
	http.HandleFunc("/download", handler.HandleDownload)
	http.HandleFunc("/download-zip", handler.HandleDownloadZip)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Start server
	port := ":19080"
	log.Printf("Server starting on http://localhost%s", port)
	log.Println("Supported formats: JPG, JPEG, PNG, TXT")

	// Setup graceful shutdown
	go func() {
		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "./static/index.html")
}
