# Go PDF Converter

A web application built in Go for converting files (images, text documents, and DOCX) to PDF format.

## Features

-   **Single and Multiple File Conversion**: Upload one or more files and get a converted PDF or a ZIP archive of PDFs.
-   **Supported Formats**:
    -   Images: JPG, JPEG, PNG
    -   Documents: TXT, DOCX
-   **Database Integration**: Stores converted files in a PostgreSQL database.
-   **Dockerized**: Comes with a `docker-compose.yml` for easy setup and deployment.
-   **Modern UI**: A simple and clean user interface with drag-and-drop support.

## Project Structure

```
.
├── main.go               # Main application entry point
├── converter.go          # File to PDF conversion logic
├── handlers.go           # HTTP request handlers
├── database.go           # Database interaction logic
├── static/index.html     # Frontend HTML, CSS, and JS
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── Dockerfile            # Dockerfile for building the application
├── docker-compose.yml    # Docker Compose setup
├── traefik.yml           # Traefik configuration (optional)
└── README.md             # This file
```

## Prerequisites

-   [Docker](https://docs.docker.com/get-docker/)
-   [Docker Compose](https://docs.docker.com/compose/install/)

## How to Run

1.  **Clone the repository**:
    ```sh
    git clone <repository-url>
    cd go-pdf-converter
    ```

2.  **Run with Docker Compose**:
    ```sh
    docker-compose up --build
    ```
    This will start the Go application and a PostgreSQL database.

3.  **Access the application**:
    Open your web browser and go to `http://localhost:19080`.

## API Endpoints

### `POST /upload`

Uploads one or more files to be converted to PDF.

-   **Request:** `multipart/form-data` with one or more files in the `files` field.
-   **Success Response (200 OK)**:
    -   For a single file:
        ```json
        {
            "success": true,
            "message": "File uploaded and converted successfully",
            "file_id": 1
        }
        ```
    -   For multiple files:
        ```json
        {
            "success": true,
            "message": "2 files uploaded and converted successfully",
            "file_ids": [1, 2]
        }
        ```
-   **Error Response (4xx or 5xx)**:
    ```json
    {
        "success": false,
        "message": "Error message"
    }
    ```

### `GET /download`

Downloads a single converted PDF file.

-   **Query Parameter:** `id` (the ID of the file).
-   **Example:** `GET /download?id=1`

### `GET /download-zip`

Downloads multiple converted PDF files as a ZIP archive.

-   **Query Parameter:** `ids` (a comma-separated list of file IDs).
-   **Example:** `GET /download-zip?ids=1,2,3`

## Technical Details

-   **Backend**: Go
-   **Database**: PostgreSQL
-   **PDF Generation**: `gofpdf` for images and text, `libreoffice` for DOCX files.
-   **Containerization**: Docker and Docker Compose.
