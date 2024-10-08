# CSV importer

## It imports a CSV file to database.
### Large CSV files are supported, optimised for speed

### Install as a command line
```
go install github.com/olbrichattila/gocsvimporter/cmd/csvimporter@latest
```

How it works:

- Analyses the CSV file and determine the file types
- Creates the table, drops if already exists
- Import data

Import modes are different per database type. I've tried to find the best settings for each on them

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

Create a file called ```.env.csvimporter``` next to the application or export the variables like ```export DB_CONNECTION=sqlite```

Usage:

```
csvimporter source.csv desttablename ";"
```

where ";" is csv separator, and this parameter is optional, if not set then the delimiter is a comma ","

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

## Optional parameters in .env.csvimporter

- The BATCH_SIZE, is how many rows are sent to the database engine per insertSQL (Firebird does not souport it and will be ignored)
- The MAX_CONNECTION_COUNT is how many connection (max) should be established to the database at the same time (SQLite does not support that)

The default bath size is 100, the default max connection count is 10.

If the values are incorrectly set, (not a number), then it will fall back to default values

- The BATCH_INSERT can be set to ON or OFF, This will overwrite the default configuration / database type (in Firebird it's ignored)
- The MULTIPLE_CONNECTIONS can be set to ON or OFF, This will overwrite the default configuration / database type (in sqLite it's ignored)
- The TRANSACTIONAL can be set to ON or OFF, This will overwrite the default configuration / database type

```
BATCH_SIZE=500
MAX_CONNECTION_COUNT=25
BATCH_INSERT=on
MULTIPLE_CONNECTIONS=on
TRANSACTIONAL=off
```

## Speed analyses with 2 million rows:
12 Columns: (Index, Customer Id, First Name, Last Name, Company, City, Country, Phone 1, Phone 2, Email, Subscription Date, Website)

Running on: Ubuntu Linux
Dell Inc. OptiPlex 7010
Intel® Core™ i7-3770S CPU @ 3.10GHz × 8
SSD

### Sqlite
```
Analyzing CSV...
Found 12 fields
Row count:2000000

Running in transactional mode
Running in batch insert mode
1 Connection opened
1 Transaction started
Importing: 100% Active threads: [ ] 
Done
1 transactions committed
1 connections closed

Full Analysis time: 0 minutes 15 seconds
Full duration time: 0 minutes 36 seconds
Total: 0 minutes 52 seconds
```

### MySql
```
Analyzing CSV...
Found 12 fields
Row count:2000000

Running in transactional mode
Running in multiple threads mode
Running in batch insert mode
10 Connection opened
10 Transaction started
Importing: 100% Active threads: [OOO OOOOOO] 
Done
10 transactions committed
10 connections closed

Full Analysis time: 0 minutes 15 seconds
Full duration time: 0 minutes 50 seconds
Total: 1 minutes 5 seconds
```

### PostgesQl
```
Analyzing CSV...
Found 12 fields
Row count:2000000

Running in transactional mode
Running in multiple threads mode
Running in batch insert mode
10 Connection opened
10 Transaction started
Importing: 100% Active threads: [OOO OOOOOO] 
Done
10 transactions committed
10 connections closed

Full Analysis time: 0 minutes 15 seconds
Full duration time: 0 minutes 28 seconds
Total: 0 minutes 43 seconds
```

### Firebird
```
Analyzing CSV...
Found 12 fields
Row count:2000000

Running in transactional mode
Running in multiple threads mode

10 Connection opened
10 Transaction started
Importing: 100% Active threads: [OOOOOOOOOO] 
Done
10 transactions committed
10 connections closed

Full Analysis time: 0 minutes 16 seconds
Full duration time: 5 minutes 26 seconds
Total: 5 minutes 42 seconds
```

## Make targets
```
## Some test imports
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

## Next steps, possible improvements
- Make it distributable, when analyzing the file, record some file pointer number and split up the importer by distribution the application between (pods, servers, virtual machines). Each one of them will open the CSV in a readonly / file shared mode and will start pushing to the database engine from it's designated file pointer until it designated target file pointer. In theory this could work in one machine to process the file in go routines per block. The distributed importers would work the same way


