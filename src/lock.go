package importer

import (
	"strconv"
	"strings"
	"sync"
)

func newLocker(count int) *lockers {
	l := &lockers{}
	l.init(count)
	return l
}

type locker struct {
	mu     *sync.Mutex
	locked bool
}

func (l *locker) isLocked() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.locked
}

func (l *locker) lock() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.locked = true
}

func (l *locker) unLock() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.locked = false
}

func (l *locker) wait() {
	for {
		if !l.isLocked() {
			break
		}
	}
}

type lockers struct {
	locks []*locker
}

func (l *lockers) init(count int) {
	l.locks = make([]*locker, count)
	for i := range l.locks {
		l.locks[i] = &locker{
			mu: &sync.Mutex{},
		}
	}
}

func (l *lockers) getLockerById(id int) *locker {
	return l.locks[id]
}

func (l *lockers) waitAll() {
	for {
		allDone := true
		for i := range l.locks {
			if l.locks[i].isLocked() {
				allDone = false
			}
		}
		if allDone {
			break
		}
	}
}

func (l *lockers) getNextUnclockedId() int {
	for {
		for i := range l.locks {
			if !l.locks[i].isLocked() {
				return i
			}
		}
	}
}

func (l *lockers) getActiveThreadReport() string {
	var ids []string
	for i := range l.locks {
		if l.locks[i].isLocked() {
			ids = append(ids, strconv.Itoa(i))
		}
	}

	return strings.Join(ids, ", ")
}
