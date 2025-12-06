# Containerization and Database Migration Guide

This guide details the steps to containerize the Go PDF Converter application and migrate the database from SQLite to PostgreSQL.

## 1. Docker Configuration

### Dockerfile

Create a `Dockerfile` in the project root:

```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git for fetching dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./main"]
```

### docker-compose.yml

Create a `docker-compose.yml` in the project root:

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=pdfconverter
    depends_on:
      - db

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=pdfconverter
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

volumes:
  postgres_data:
```

## 2. Go Dependencies

Add the PostgreSQL driver to your project:

```bash
go get github.com/lib/pq
```

## 3. Code Changes

### database.go

Update `database.go` to use PostgreSQL. Note the changes in SQL placeholders (`$1` instead of `?`) and types (`SERIAL`, `BYTEA`).

```go
package main

import (
 "database/sql"
 "log"
 "time"

 _ "github.com/lib/pq" // PostgreSQL driver
)

type FileRecord struct {
 ID           int
 OriginalName string
 PDFData      []byte
 UploadedAt   time.Time
}

type Database struct {
 db *sql.DB
}

func NewDatabase(connStr string) (*Database, error) {
 // Open connection to Postgres
 db, err := sql.Open("postgres", connStr)
 if err != nil {
  return nil, err
 }

 // Verify connection
 if err := db.Ping(); err != nil {
  return nil, err
 }

 // Create table if not exists
 createTableSQL := `CREATE TABLE IF NOT EXISTS files (
  id SERIAL PRIMARY KEY,
  original_name TEXT NOT NULL,
  pdf_data BYTEA NOT NULL,
  uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
 );`

 _, err = db.Exec(createTableSQL)
 if err != nil {
  return nil, err
 }

 return &Database{db: db}, nil
}

func (d *Database) SaveFile(originalName string, pdfData []byte) (int64, error) {
 var id int64
 // Postgres uses $1, $2 placeholders and RETURNING to get the ID
 err := d.db.QueryRow(
  "INSERT INTO files (original_name, pdf_data) VALUES ($1, $2) RETURNING id",
  originalName, pdfData,
 ).Scan(&id)
 
 if err != nil {
  return 0, err
 }
 return id, nil
}

func (d *Database) GetFile(id int64) (*FileRecord, error) {
 var record FileRecord
 // Postgres uses $1 placeholder
 err := d.db.QueryRow(
  "SELECT id, original_name, pdf_data, uploaded_at FROM files WHERE id = $1",
  id,
 ).Scan(&record.ID, &record.OriginalName, &record.PDFData, &record.UploadedAt)

 if err != nil {
  return nil, err
 }
 return &record, nil
}

func (d *Database) GetFiles(ids []int64) ([]*FileRecord, error) {
 records := make([]*FileRecord, 0, len(ids))

 for _, id := range ids {
  record, err := d.GetFile(id)
  if err != nil {
   log.Printf("Error fetching file ID %d: %v", id, err)
   continue
  }
  records = append(records, record)
 }

 return records, nil
}

func (d *Database) Close() error {
 return d.db.Close()
}
```

### main.go

Update `main.go` to configure the database connection using environment variables.

```go
package main

import (
 "fmt"
 "log"
 "net/http"
 "os"
 "os/signal"
 "syscall"
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
 connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
  dbHost, dbPort, dbUser, dbPassword, dbName)

 // Initialize database
 db, err := NewDatabase(connStr)
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
```

## 4. Running the Application

1. **Clean up dependencies**:

    ```bash
    go mod tidy
    ```

2. **Start with Docker Compose**:

    ```bash
    docker-compose up --build
    ```

3. **Access the application**:
    Open [http://localhost:8080](http://localhost:8080) in your browser.
