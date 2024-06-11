# CSV importer

## It imports a CSV file to database.
### Large CSV files are supported, optimised for speed

How it works:

- Analyses the CSV file and determine the file types
- Creates the table, drops if already exists
- Import data

The import happening in one transaction, and in multiple SQL baches (except FirebirdSQL which does not support it)

Currently the following database types are supported:

- SqLite
- MySql
- Postgresql
- FirebirdSql

Usage: ```go run . <params>``` or ```./importcs <params>```

Example:
```
go run . data.csv vehicles ";"
./importcs data.csv vehicles ";"
```

Parameters
1. CSV file name
2. Table name to create
3. The CSV delimiter, (in quotes). This parameter is optional, if not set then ","

### Database settings.

Create a file called ```.env``` next to the application or export the variables like ```export DB_CONNECTION=sqlite```

Examples:
## Sqlite
```
DB_CONNECTION=sqlite
DB_DATABASE=./database/database.sqlite
```

## MySql
```
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=migrator
DB_USERNAME=migrator
DB_PASSWORD=H8E7kU8Y
```

## Postgresql
```
DB_CONNECTION=pgsql
DB_HOST=127.0.0.1
DB_PORT=5432
DB_DATABASE=postgres
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_SSLMODE=disable
```

## Firebird SQL
```
DB_CONNECTION=firebird
DB_HOST=127.0.0.1
DB_PORT=3050
DB_DATABASE=/firebird/data/employee.fdb
DB_USERNAME=SYSDBA
DB_PASSWORD=masterkey
```

## Speed analisys with 2 million rows:
Running on: Ubuntu Linux
Dell Inc. OptiPlex 7010
Intel® Core™ i7-3770S CPU @ 3.10GHz × 8
SSD

### Sqlite
```1 minutes 50 seconds```

### MySql
```2 minutes 38 seconds```

### PostgesQl
```2 minutes 51 seconds```

### Firebird
``` ```

## Make targets
```
## Some test inports
make vehicles:
make customers:
	
## Switch between test environments
make switch-sqlite:
make switch-mysql:
make switch-pgsql:
make switch-firebird:
```

## Test locally:
There is a docker folder:
```
cd docker
docker-compose up -d
```

## What is next
- full code cleanup
- tests
- paralell SQL to inmpove speed



