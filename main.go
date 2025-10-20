package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize database
	db, err := NewDatabase("./files.db")
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
	port := ":8080"
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
