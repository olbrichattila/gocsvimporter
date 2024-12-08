# Go CSV Importer: Multi-threaded Fast CSV to Database Tool

## Overview
`Go CSV Importer` is a powerful, multi-threaded command-line tool designed to import large CSV files into databases quickly and efficiently. It supports various database types, offers customizable import modes, and is optimized for performance.

### Key Features
- **Support for Large CSV Files**: Handles millions of rows with ease.
- **Multi-Database Compatibility**: Works with SQLite, MySQL, PostgreSQL, and Firebird.
- **Configurable Import Modes**:
  - Transactional or non-transactional imports.
  - Batch SQL inserts or row-by-row operations.
  - Single or multiple database connections.

---

## Installation

Install the latest version directly from Go:

```bash
go install github.com/olbrichattila/gocsvimporter/cmd/csvimporter@latest
```

---

## How It Works

1. **Analyze the CSV**: Determines data types and structures.
2. **Prepare the Database**: Automatically creates the necessary table, dropping it if it already exists.
3. **Import Data**: Optimizes the process based on database type and chosen parameters.

### Usage
Run the tool with:

```bash
csvimporter <csv_file> <table_name> [delimiter]
```

#### Parameters:
1. **CSV File**: Path to the CSV file.
2. **Table Name**: Target database table name.
3. **Delimiter** *(Optional)*: CSV delimiter (default: `,`).

#### Example:
```bash
csvimporter data.csv vehicles ";"
```

---

## Supported Databases

### SQLite
```env
DB_CONNECTION=sqlite
DB_DATABASE=./database/database.sqlite
```

### MySQL
```env
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=mydb
DB_USERNAME=myuser
DB_PASSWORD=mypassword
```

### PostgreSQL
```env
DB_CONNECTION=pgsql
DB_HOST=127.0.0.1
DB_PORT=5432
DB_DATABASE=postgres
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_SSLMODE=disable
```

### Firebird
```env
DB_CONNECTION=firebird
DB_HOST=127.0.0.1
DB_PORT=3050
DB_DATABASE=/path/to/database.fdb
DB_USERNAME=SYSDBA
DB_PASSWORD=masterkey
```

---

## Configuration

Create a `.env.csvimporter` file in the application's directory or export environment variables:

```env
BATCH_SIZE=500               # Rows per batch (default: 100)
MAX_CONNECTION_COUNT=25      # Maximum connections (default: 10)
BATCH_INSERT=on              # Enable/disable batch insert
MULTIPLE_CONNECTIONS=on      # Enable/disable multi-threading
TRANSACTIONAL=off            # Enable/disable transactions
```

*Note: Unsupported options for certain databases (e.g., SQLite) are ignored.*

---

## Performance

### Speed Test: 2 Million Rows
**System Configuration**:
- **OS**: Ubuntu Linux
- **Processor**: Intel® Core™ i7-3770S @ 3.10GHz
- **Storage**: SSD

| Database    | Duration   | Mode                     | Threads |
|-------------|------------|--------------------------|---------|
| **SQLite**  | 52 seconds | Transactional, Batch SQL | 1       |
| **MySQL**   | 65 seconds | Transactional, Multi-Threaded | 10      |
| **PostgreSQL** | 43 seconds | Transactional, Multi-Threaded | 10      |
| **Firebird**| 5m 42s     | Transactional, Multi-Threaded | 10      |

---

## Makefile Targets

### Test Imports
```bash
make vehicles
make customers
```

### Switch Environments
```bash
make switch-sqlite
make switch-mysql
make switch-pgsql
make switch-firebird
```

---

## Local Testing with Docker
A `docker-compose` setup is provided for testing:

```bash
cd docker
docker-compose up -d
```

---

## Roadmap

### Planned Improvements
- **Distributed Import**: Split CSV files across multiple instances (pods/servers) for faster parallel imports.
- **Enhanced Configuration**: Support more advanced database-specific settings.

---

Start importing your CSV files faster and more efficiently today!
