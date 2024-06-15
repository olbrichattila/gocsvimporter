package importer

import (
	"database/sql"
	"fmt"
	"math"
	"time"
)

const (
	numberOfThreads = 10
)

type batch = [][]any

func newImporter2(
	dBconfig dBConfiger,
	csv csvReader,
	sQLGen sQLGenerator,
	storer storager,
) importer2 {
	return &imp{
		dBconfig: dBconfig,
		csv:      csv,
		sQLGen:   sQLGen,
		storer:   storer,
	}
}

type importer2 interface {
	importCsv() (float64, float64, float64, error)
}

type threadConnection struct {
	db *sql.DB
	tx *sql.Tx
}

type threadConnections []*threadConnection

func (t *threadConnection) getExecutor() SQLExecutor {
	if t.tx != nil {
		return t.tx
	}

	return t.db
}

type imp struct {
	dBconfig    dBConfiger
	csv         csvReader
	sQLGen      sQLGenerator
	storer      storager
	connections threadConnections
	progress    int
	rowNr       int
	rowCount    int
}

func (i *imp) importCsv() (float64, float64, float64, error) {
	isTrasactional := i.dBconfig.needTransactions()
	if isTrasactional {
		fmt.Println("Running in transactional mode")
	}

	haveMultipleTheads := i.dBconfig.haveMultipleThreads()
	connectionCount := 1
	if haveMultipleTheads {
		fmt.Println("Running in multiple threads mode")
		connectionCount = numberOfThreads
	}

	haveBatchInsert := i.dBconfig.haveBatchInsert()
	if haveBatchInsert {
		fmt.Println("Running in batch insert mode")
	}

	startedAt := time.Now()
	i.csv.init()
	defer i.csv.close()
	headers := i.csv.header()
	fmt.Printf("Found %d fields\n", len(headers))
	err := i.createConnections(connectionCount, isTrasactional)
	if err != nil {
		return 0, 0, 0, err
	}

	defer i.closeConnections()

	i.dropAndCreateTable()
	insertSql := i.sQLGen.createInsertSQL()

	batchIndex := 0
	connectionId := 0
	var batchdata batch
	locks := newLocker(connectionCount)

	i.resetProgress()
	importStartedAt := time.Now()
	for i.csv.next() {
		i.showProgress(locks)

		if !haveBatchInsert {
			if haveMultipleTheads {
				connectionId = locks.getNextUnclockedId()
				go i.executeBatchInThread(locks.getLockerById(connectionId), isTrasactional, i.connections[connectionId], insertSql, i.csv.row()...)
			} else {
				i.executeBatchInThread(nil, isTrasactional, i.connections[0], insertSql, i.csv.row()...)
			}

			continue
		}

		batchdata = append(batchdata, i.csv.row())
		if batchIndex == batchSize {
			insertSql, bindingPars := i.sQLGen.createBatchInsertSQL(batchdata, true)

			if haveMultipleTheads {
				connectionId = locks.getNextUnclockedId()
				go i.executeBatchInThread(locks.getLockerById(connectionId), isTrasactional, i.connections[connectionId], insertSql, bindingPars...)
			} else {
				i.executeBatchInThread(nil, isTrasactional, i.connections[0], insertSql, bindingPars...)
			}

			batchIndex = 0
			batchdata = nil
		} else {
			batchIndex++
		}
	}

	if haveBatchInsert {
		connectionId = locks.getNextUnclockedId()
		insertSql, bindingPars := i.sQLGen.createBatchInsertSQL(batchdata, false)
		if haveMultipleTheads {
			go i.executeBatchInThread(locks.getLockerById(connectionId), isTrasactional, i.connections[connectionId], insertSql, bindingPars...)
		} else {
			go i.executeBatchInThread(nil, isTrasactional, i.connections[0], insertSql, bindingPars...)
		}
	}

	if haveMultipleTheads {
		locks.waitAll()
	}

	finishedAt := time.Now()

	importTime := finishedAt.Sub(importStartedAt).Seconds()
	pharseTime := importStartedAt.Sub(startedAt).Seconds()
	totalTime := importTime + pharseTime

	return pharseTime, importTime, totalTime, nil
}

func (i *imp) dropAndCreateTable() error {
	err := i.dropTable(i.connections[0].db)
	if err != nil {
		return err
	}

	return i.createTable(i.connections[0].db)
}

func (i *imp) dropTable(db *sql.DB) error {
	return i.storer.execute(db, i.sQLGen.getDropTableSQL())
}

func (i *imp) createTable(db *sql.DB) error {
	return i.storer.execute(db, i.sQLGen.ceateTableSQL(i.csv.header()))
}

func (i *imp) createConnections(count int, isTransactional bool) error {
	i.connections = make(threadConnections, count)

	for x := range i.connections {
		conn, err := i.dBconfig.getNewConnection()
		if err != nil {
			i.closeConnections()
			return err
		}
		fmt.Printf("%d Connection opened\n", x+1)
		var tx *sql.Tx
		if isTransactional {
			tx, err = conn.Begin()
			if err != nil {
				i.closeConnections()
				return err
			}
			fmt.Printf("%d Transaction started\n", x+1)
		}
		i.connections[x] = &threadConnection{db: conn, tx: tx}
	}

	return nil
}

func (i *imp) closeConnections() {
	for x := range i.connections {
		if i.connections[x] != nil {
			// We don't want to handle this error
			if i.connections[x].tx != nil {
				err := i.connections[x].tx.Commit()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("%d Transaction committed\n", x+1)
			}
			err := i.connections[x].db.Close()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%d Connection closed\n", x+1)
		}
	}
}

func (i *imp) resetProgress() {
	i.rowNr = 0
	i.progress = 0
	i.rowCount = i.csv.rowCount()
}

func (i *imp) showProgress(locks *lockers) {
	i.rowNr++
	percent := int(math.Ceil(float64(i.rowNr) / float64(i.rowCount) * 100))
	if percent != i.progress {
		i.progress = percent
		activeTreads := locks.getActiveThreadReport()
		fmt.Printf("\rImporting: %d%% %s          ", i.progress, activeTreads)
	}
}

func (i *imp) executeBatchInThread(
	l *locker,
	transactional bool,
	connection *threadConnection,
	insertSQL string,
	bindingPars ...any,
) {
	//fmt.Println(l.isLocked())
	var err error
	if l != nil {
		l.wait()
		l.lock()
	}
	defer func() {
		if l != nil {
			l.unLock()
		}
	}()

	err = i.storer.execute(connection.getExecutor(), insertSQL, bindingPars...)
	if err != nil {
		// log tx errors?
		fmt.Println(err)
		return
	}
}
