package main

import "fmt"

func displayHelp() {
	fmt.Println(`
Usage: csvimporter <csv_file> <table_name> [delimiter]
  Where delimiter is optional, by default it is ,
  In case of semicolon, please use quotes ";"

Optional parameters:
	-sh=headerFileName (only save header info into a JSON file after file analysis)
	-lh=headerFileName (skip file analysis and uses the pre-saved header for faster regular imports)

Database settings.
Create a file called .env.csvimporter

Examples:
## Sqlite
--------------------------
DB_CONNECTION=sqlite
DB_DATABASE=./database/database.sqlite
--------------------------

## MySql
--------------------------
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=migrator
DB_USERNAME=migrator
DB_PASSWORD=H8E7kU8Y


## Postgresql
--------------------------
DB_CONNECTION=pgsql
DB_HOST=127.0.0.1
DB_PORT=5432
DB_DATABASE=postgres
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_SSLMODE=disable
--------------------------

## Firebird SQL
--------------------------
DB_CONNECTION=firebird
DB_HOST=127.0.0.1
DB_PORT=3050
DB_DATABASE=/firebird/data/employee.fdb
DB_USERNAME=SYSDBA
DB_PASSWORD=masterkey
--------------------------

## Optional parameters in .env.csvimporter

- The BATCH_SIZE, is how many rows are sent to the database engine per insertSQL (Firebird does not souport it and will be ignored)
- The MAX_CONNECTION_COUNT is how many connection (max) should be established to the database at the same time (SQLite does not support that)

The default bath size is 100, the default max connection count is 10.

If the values are incorrectly set, (not a number), then it will fall back to default values

- The BATCH_INSERT can be set to ON or OFF, This will overwrite the default configuration / database type (in Firebird it's ignored)
- The MULTIPLE_CONNECTIONS can be set to ON or OFF, This will overwrite the default configuration / database type (in sqLite it's ignored)
- The TRANSACTIONAL can be set to ON or OFF, This will overwrite the default configuration / database type

--------------------------
BATCH_SIZE=500
MAX_CONNECTION_COUNT=25
BATCH_INSERT=on
MULTIPLE_CONNECTIONS=on
TRANSACTIONAL=off
--------------------------`)
}
