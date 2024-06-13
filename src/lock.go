package importer

import (
	"sync"
)

func newLocker() *lock {
	return &lock{c: 0}
}

type lock struct {
	c int
	s sync.Mutex
}

func (l *lock) lock() {
	l.s.Lock()
	defer l.s.Unlock()
	l.c++
}

func (l *lock) unlock() {
	l.s.Lock()
	defer l.s.Unlock()
	l.c--
}

func (l *lock) count() int {
	l.s.Lock()
	defer l.s.Unlock()

	return l.c
}

func (l *lock) wait(count int) {
	for {
		if l.count() < count {
			break
		}
	}
}
