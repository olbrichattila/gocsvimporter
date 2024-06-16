# CSV importer

## It imports a CSV file to database.
### Large CSV files are supported, optimised for speed

How it works:

- Analyses the CSV file and determine the file types
- Creates the table, drops if already exists
- Import data

Import modes are diffent per database type. I've tried to find the best settings for each on them

The import can run
- with/without transaction
- batched-SQL/Insert SQL per row
- One connection/multiple connections

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

## Optional parameters in .env

- The batch size, is how many rows are sent to the database engine per insertSQL (Firebird does not souport it and will be ignored)
- The max connection count is how many connection (max) should be established to the database at the same time (SQLite does not support that)
The defualt bath size is 100, the default max connection count is 10.
If the values are incorrectly set, (not a number), then it will fall back to default values

```
BATCH_SIZE=500
MAX_CONNECTION_COUNT=25
```

## Speed analisys with 2 million rows:
12 Columns: (Index, Customer Id, First Name, Last Name, Company, City, Country, Phone 1, Phone 2, Email, Subscription Date, Website)

Running on: Ubuntu Linux
Dell Inc. OptiPlex 7010
Intel® Core™ i7-3770S CPU @ 3.10GHz × 8
SSD

### Sqlite
```
Analising CSV...
Found 12 fields
Row count:2000000

Running in transactional mode
Running in batch insert mode
1 Connection opened
1 Transaction started
Importing: 100% Active threads: [ ] 
Done
1 transactions commtted
1 connections closed

Full Analysis time: 0 minutes 15 seconds
Full duration time: 0 minutes 36 seconds
Total: 0 minutes 52 seconds
```

### MySql
```
Analising CSV...
Found 12 fields
Row count:2000000

Running in transactional mode
Running in multiple threads mode
Running in batch insert mode
10 Connection opened
10 Transaction started
Importing: 100% Active threads: [OOO OOOOOO] 
Done
10 transactions commtted
10 connections closed

Full Analysis time: 0 minutes 15 seconds
Full duration time: 0 minutes 50 seconds
Total: 1 minutes 5 seconds
```

### PostgesQl
```
Analising CSV...
Found 12 fields
Row count:2000000

Running in transactional mode
Running in multiple threads mode
Running in batch insert mode
10 Connection opened
10 Transaction started
Importing: 100% Active threads: [OOO OOOOOO] 
Done
10 transactions commtted
10 connections closed

Full Analysis time: 0 minutes 15 seconds
Full duration time: 0 minutes 28 seconds
Total: 0 minutes 43 seconds
```

### Firebird
```
Analising CSV...
Found 12 fields
Row count:2000000

Running in transactional mode
Running in multiple threads mode

10 Connection opened
10 Transaction started
Importing: 100% Active threads: [OOOOOOOOOO] 
Done
10 transactions commtted
10 connections closed

Full Analysis time: 0 minutes 16 seconds
Full duration time: 5 minutes 26 seconds
Total: 5 minutes 42 seconds
```

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

## TODO: Notes
- Code cleanup
    
