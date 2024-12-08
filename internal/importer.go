package importer

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	database "github.com/olbrichattila/gocsvimporter/internal/db"
	"github.com/olbrichattila/gocsvimporter/internal/env"
	"github.com/olbrichattila/gocsvimporter/internal/storage"
)

const (
	defaultBatchSize       = 100
	defaultNumberOfThreads = 10
	progressRefreshCount   = 1000
)

type batch = [][]any

func newImporter(
	dBConfig database.DBConfiger,
	csv csvReader,
	sQLGen sQLGenerator,
	storer storage.Storager,
) importer {
	return &imp{
		dBConfig: dBConfig,
		csv:      csv,
		sQLGen:   sQLGen,
		storer:   storer,
	}
}

type importer interface {
	importCsv() (float64, float64, float64, error)
}

type imp struct {
	dBConfig    database.DBConfiger
	csv         csvReader
	sQLGen      sQLGenerator
	storer      storage.Storager
	connections threadConnections
	progress    int
	rowNr       int
	rowCount    int
}

func (i *imp) importCsv() (float64, float64, float64, error) {
	startedAt := time.Now()
	isTransactional, HaveMultipleThreads, HaveBatchInsert, connectionCount, err := i.init()
	if err != nil {
		return 0, 0, 0, err
	}
	defer i.csv.close()
	batchSize := i.getIntEnv(env.BatchSize, defaultBatchSize)
	if i.dBConfig.HaveBatchInsert() {
		fmt.Printf("Batch size is %d\n", batchSize)
	}

	err = i.createConnections(connectionCount, isTransactional)
	if err != nil {
		return 0, 0, 0, err
	}

	defer i.closeConnections()

	err = i.dropAndCreateTable()
	if err != nil {
		return 0, 0, 0, err
	}
	insertSQL := i.sQLGen.createInsertSQL()

	batchIndex := 0
	connectionID := 0
	var batchData batch
	locks := newLocker(connectionCount)

	i.resetProgress()
	importStartedAt := time.Now()
	for i.csv.next() {
		i.showProgress(locks)

		if !HaveBatchInsert {
			if HaveMultipleThreads {
				connectionID = locks.getNextUnlockedID()
				locks.getLockerByID(connectionID).lock()
				go i.executeBatchInThread(locks.getLockerByID(connectionID), i.connections[connectionID], insertSQL, i.csv.row()...)
				continue
			}

			i.executeBatch(i.connections[0], insertSQL, i.csv.row()...)
			continue
		}

		batchData = append(batchData, i.csv.row())
		if batchIndex == batchSize {
			insertSQL, bindingPars := i.sQLGen.createBatchInsertSQL(batchData, true)

			if HaveMultipleThreads {
				connectionID = locks.getNextUnlockedID()
				locks.getLockerByID(connectionID).lock()
				go i.executeBatchInThread(locks.getLockerByID(connectionID), i.connections[connectionID], insertSQL, bindingPars...)
			} else {
				i.executeBatch(i.connections[0], insertSQL, bindingPars...)
			}

			batchIndex = 0
			batchData = nil
			continue
		}

		batchIndex++
	}

	if HaveBatchInsert {
		insertSQL, bindingPars := i.sQLGen.createBatchInsertSQL(batchData, false)
		if HaveMultipleThreads {
			connectionID = locks.getNextUnlockedID()
			locks.getLockerByID(connectionID).lock()
			go i.executeBatchInThread(locks.getLockerByID(connectionID), i.connections[connectionID], insertSQL, bindingPars...)
		} else {
			i.executeBatch(i.connections[0], insertSQL, bindingPars...)
		}
	}

	if HaveMultipleThreads {
		locks.waitAll()
	}

	fmt.Printf("\nDone\n")
	finishedAt := time.Now()
	importTime := finishedAt.Sub(importStartedAt).Seconds()
	phraseTime := importStartedAt.Sub(startedAt).Seconds()
	totalTime := importTime + phraseTime

	return phraseTime, importTime, totalTime, nil
}

func (i *imp) init() (bool, bool, bool, int, error) {
	err := i.initCsv()
	if err != nil {
		return false, false, false, 0, err
	}

	isTransactional, HaveMultipleThreads, HaveBatchInsert, connectionCount := i.initConfig()
	return isTransactional, HaveMultipleThreads, HaveBatchInsert, connectionCount, nil
}

func (i *imp) initCsv() error {
	err := i.csv.init()
	if err != nil {
		return err
	}

	headers := i.csv.header()
	fmt.Printf("Found %d fields\nRow count:%d\n\n", len(headers), i.csv.rowCount())
	return nil
}

func (i *imp) initConfig() (bool, bool, bool, int) {
	isTransactional := i.dBConfig.NeedTransactions()
	if isTransactional {
		fmt.Println("Running in transactional mode")
	}

	HaveMultipleThreads := i.dBConfig.HaveMultipleThreads()
	connectionCount := 1
	if HaveMultipleThreads {
		fmt.Println("Running in multiple threads mode")
		connectionCount = i.getIntEnv(env.MaxConnectionCount, defaultNumberOfThreads)
	}

	HaveBatchInsert := i.dBConfig.HaveBatchInsert()
	if HaveBatchInsert {
		fmt.Println("Running in batch insert mode")
	}

	return isTransactional, HaveMultipleThreads, HaveBatchInsert, connectionCount
}

func (i *imp) getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return intValue
	}

	return defaultValue
}

func (i *imp) dropAndCreateTable() error {
	err := i.dropTable(i.connections[0].db)
	if err != nil {
		return err
	}

	return i.createTable(i.connections[0].db)
}

func (i *imp) dropTable(db *sql.DB) error {
	return i.storer.Execute(db, i.sQLGen.getDropTableSQL())
}

func (i *imp) createTable(db *sql.DB) error {
	return i.storer.Execute(db, i.sQLGen.cerateTableSQL(i.csv.header()))
}

func (i *imp) createConnections(count int, isTransactional bool) error {
	i.connections = make(threadConnections, count)
	for x := range i.connections {
		conn, err := i.dBConfig.GetNewConnection()
		if err != nil {
			i.closeConnections()
			return err
		}
		var tx *sql.Tx
		if isTransactional {
			tx, err = conn.Begin()
			if err != nil {
				i.closeConnections()
				return err
			}

		}
		i.connections[x] = &threadConnection{db: conn, tx: tx}
	}
	fmt.Printf("%d Connection opened\n", count)
	if isTransactional {
		fmt.Printf("%d Transaction started\n", count)
	}
	return nil
}

func (i *imp) closeConnections() {
	commitCount := 0
	closeCount := 0
	if i.connections == nil {
		return
	}
	for x := range i.connections {
		if i.connections[x].db != nil {
			if i.connections[x].tx != nil {
				err := i.connections[x].tx.Commit()
				if err != nil {
					fmt.Println("Commit error: " + err.Error())
				}
				i.connections[x].tx = nil
				commitCount++
			}
			err := i.connections[x].db.Close()
			if err != nil {
				fmt.Println(err)
			}
			i.connections[x].db = nil
			closeCount++
		}
	}
	fmt.Printf("%d transactions committed\n%d connections closed", commitCount, closeCount)
}

func (i *imp) resetProgress() {
	i.rowNr = 0
	i.progress = 0
	i.rowCount = i.csv.rowCount()
}

func (i *imp) showProgress(locks *lockers) {
	i.rowNr++
	if i.rowCount < 0 {
		i.showSimpleProgress(locks)
		return
	}

	i.showFullProgress(locks)
}

func (i *imp) showFullProgress(locks *lockers) {
	percent := int(math.Ceil(float64(i.rowNr) / float64(i.rowCount) * 100))
	if percent != i.progress {
		i.progress = percent
		fmt.Printf("\rImporting: %d%% Active threads: [%s] ", i.progress, locks.getActiveThreadReport())
	}
}

func (i *imp) showSimpleProgress(locks *lockers) {
	if i.rowNr%progressRefreshCount == 0 {
		fmt.Printf("\rImporting in progress... Active threads: [%s] ", locks.getActiveThreadReport())
	}
}

func (i *imp) executeBatchInThread(
	l *locker,
	connection *threadConnection,
	insertSQL string,
	bindingPars ...any,
) {
	defer l.unLock()
	err := i.storer.Execute(connection.getExecutor(), insertSQL, bindingPars...)
	if err != nil {
		fmt.Println("BatchThreadException: " + err.Error())
		return
	}
}

func (i *imp) executeBatch(
	connection *threadConnection,
	insertSQL string,
	bindingPars ...any,
) {
	err := i.storer.Execute(connection.getExecutor(), insertSQL, bindingPars...)
	if err != nil {
		fmt.Println("BatchException: " + err.Error())
		return
	}
}
