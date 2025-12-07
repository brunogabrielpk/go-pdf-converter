# PDF Converter Web Application

A web application built in Go for converting files (images and text documents) to PDF format.

## Features

- **Single File Upload**: Upload one file at a time and download the converted PDF
- **Multiple File Upload**: Upload multiple files and download them as a ZIP archive
- **Supported Formats**:
  - Images: JPG, JPEG, PNG
  - Documents: TXT
- **PostgreSQL Storage**: All converted PDFs are stored in a PostgreSQL database
- **Modern Web Interface**: Drag-and-drop file upload with real-time feedback
- **Containerized**: Fully containerized with Docker and Docker Compose

## Project Structure

```
pdf-converter/
├── main.go           # HTTP server and routing
├── database.go       # PostgreSQL database operations
├── converter.go      # PDF conversion logic
├── handlers.go       # HTTP request handlers
├── static/
│   └── index.html    # Frontend interface
├── go.mod            # Go module file
├── Dockerfile        # Docker build instructions
└── docker-compose.yml # Docker Compose configuration
```

## Installation

### Prerequisites

- **Docker** and **Docker Compose** (Recommended)
- OR **Go** (version 1.25 or later) and a running **PostgreSQL** instance

## Running the Application

### Using Docker Compose (Recommended)

1. Navigate to the project directory:

```bash
cd pdf-converter
```

2. Start the application and database:

```bash
docker-compose up --build
```

3. Open your browser and navigate to:

```
http://localhost:8080
```

### Running Locally

1. Ensure you have a PostgreSQL database running.

2. Set the necessary environment variables (or rely on defaults for localhost):
   - `DB_HOST`: Database host (default: localhost)
   - `DB_PORT`: Database port (default: 5432)
   - `DB_USER`: Database user (default: postgres)
   - `DB_PASSWORD`: Database password (default: postgres)
   - `DB_NAME`: Database name (default: pdfconverter)

3. Run the application:

```bash
go run .
```

## Usage

### Single File Conversion

1. Click the upload area or drag and drop a file
2. Select a supported file (JPG, JPEG, PNG, or TXT)
3. Click "Convert to PDF"
4. Download the converted PDF file

### Multiple File Conversion

1. Click the upload area or drag and drop multiple files
2. Select multiple supported files
3. Click "Convert to PDF"
4. Download all converted PDFs as a ZIP archive

## API Endpoints

### POST /upload

Upload one or more files for conversion to PDF.

**Request**: multipart/form-data with "files" field

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

- `id`: File ID returned from upload

### GET /download-zip?ids={file_ids}

Download multiple PDFs as a ZIP archive.

**Parameters**:

- `ids`: Comma-separated list of file IDs (e.g., "1,2,3")

## Technical Details

### Dependencies

- `github.com/lib/pq`: PostgreSQL driver
- `github.com/jung-kurt/gofpdf`: PDF generation library

### Image to PDF Conversion

- Images are decoded and embedded in PDF pages
- Page size automatically adjusts to image dimensions
- Maximum size limited to A4 to prevent oversized PDFs

### Text to PDF Conversion

- Text files are rendered with Courier font
- Automatic line wrapping for long lines
- Standard A4 page format

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

- **DB_HOST**: Database host address
- **DB_PORT**: Database port
- **DB_USER**: Database username
- **DB_PASSWORD**: Database password
- **DB_NAME**: Database name

## Limitations

- Maximum upload size: 32 MB per request
- Supported file formats: JPG, JPEG, PNG, TXT
- Text files use fixed font (Courier)
- Images are resized to fit A4 if larger

## License

This project is open source and available for educational purposes.
