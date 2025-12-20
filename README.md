# Go PDF Converter

A web application built in Go for converting files (images, text documents, and DOCX) to PDF format.

## Features

- **Single File Upload**: Upload one file at a time and download the converted PDF.
- **Multiple File Upload**: Upload multiple files and download them as a ZIP archive.
- **Supported Formats**:
  - Images: JPG, JPEG, PNG
  - Documents: TXT, DOCX (requires LibreOffice)
- **PostgreSQL Storage**: All converted PDFs are stored in a PostgreSQL database.
- **Modern Web Interface**: Drag-and-drop file upload with real-time feedback.
- **Containerized**: Fully containerized with Docker and Docker Compose for easy setup.

## Project Structure

```
pdf-converter/
├── main.go               # HTTP server and routing
├── database.go           # PostgreSQL database operations
├── converter.go          # PDF conversion logic
├── handlers.go           # HTTP request handlers
├── static/
│   └── index.html        # Frontend interface
├── files.db              # SQLite database for local development
├── go.mod                # Go module file
├── Dockerfile            # Docker build instructions
└── docker-compose.yml    # Docker Compose configuration
```

## Installation

### Prerequisites

- **Docker** and **Docker Compose** (Recommended)
- **Go** (version 1.25.1 or later)
- **PostgreSQL** (if not using Docker)
- **LibreOffice** (for DOCX conversion)

## Getting Started

1.  Clone the repository:
    ```bash
    git clone https://github.com/brunogabrielpk/go-pdf-converter.git
    ```
2.  Navigate to the project directory:
    ```bash
    cd go-pdf-converter
    ```

## Running the Application

### Using Docker Compose (Recommended)

1.  Navigate to the project directory:
    ```bash
    cd go-pdf-converter
    ```
2.  Start the application and database:
    ```bash
    docker-compose up --build
    ```
3.  Open your browser and navigate to:
    ```
    http://localhost:40110
    ```

### Running Locally

1.  Ensure you have a PostgreSQL database running.
2.  Set the necessary environment variables (or rely on defaults for localhost):
    - `DB_HOST`: Database host (default: localhost)
    - `DB_PORT`: Database port (default: 5432)
    - `DB_USER`: Database user (default: postgres)
    - `DB_PASSWORD`: Database password (default: postgres)
    - `DB_NAME`: Database name (default: pdfconverter)
3.  Run the application:
    ```bash
    go run .
    ```

## Usage

### Single File Conversion

1.  Click the upload area or drag and drop a file.
2.  Select a supported file (JPG, JPEG, PNG, TXT, or DOCX).
3.  Click "Convert to PDF".
4.  Download the converted PDF file.

### Multiple File Conversion

1.  Click the upload area or drag and drop multiple files.
2.  Select multiple supported files.
3.  Click "Convert to PDF".
4.  Download all converted PDFs as a ZIP archive.

## API Endpoints

### POST /upload

Upload one or more files for conversion to PDF.

**Request**: `multipart/form-data` with a "files" field.

**Response**:
```json
{
  "success": true,
  "message": "File uploaded and converted successfully",
  "file_id": 1
}
```
For multiple files:
```json
{
  "success": true,
  "message": "3 files uploaded and converted successfully",
  "file_ids": [1, 2, 3]
}
```

### GET /download?id={file_id}

Download a single converted PDF file.

**Parameters**:
- `id`: The file ID returned from the upload.

### GET /download-zip?ids={file_ids}

Download multiple PDFs as a ZIP archive.

**Parameters**:
- `ids`: A comma-separated list of file IDs (e.g., "1,2,3").

## Technical Details

### Dependencies

- `github.com/lib/pq`: PostgreSQL driver.
- `github.com/jung-kurt/gofpdf`: PDF generation library for images and text.
- `libreoffice`: Required for converting DOCX files. See [DOCX_SUPPORT_GUIDE.md](DOCX_SUPPORT_GUIDE.md) for setup instructions.

### DOCX to PDF Conversion

The application uses LibreOffice in headless mode to convert DOCX files to PDF. This is a powerful feature that ensures high-fidelity conversions. For this to work, **LibreOffice must be installed on the server where the application is running**. The included Dockerfile does not have LibreOffice installed.

Please refer to the [DOCX_SUPPORT_GUIDE.md](DOCX_SUPPORT_GUIDE.md) for detailed instructions on how to set up your environment for DOCX conversion.

### Database Schema

```sql
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    original_name TEXT NOT NULL,
    pdf_data BYTEA NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Configuration

The application is configured via environment variables:

- `DB_HOST`: Database host address
- `DB_PORT`: Database port
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name

## Limitations

- **Maximum Upload Size**: 32 MB per request.
- **DOCX Conversion**: Requires LibreOffice to be installed on the server.
- **Image Resizing**: Images larger than A4 will be resized to fit.

## License

This project is open source and available for educational purposes.
