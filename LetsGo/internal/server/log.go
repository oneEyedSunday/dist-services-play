package server

import (
	"fmt"
	"sync"
)

type Log struct {
	mu sync.Mutex
	records []Record
}

type Record struct {
	Value []byte
	Offset uint64
}

func NewLog() *Log {
	return &Log{}
}

func (l* Log) Append(record Record) (uint64, error) {
	l.acquireAndReleaseLock()

	record.Offset = uint64(len(l.records))
	l.records = append(l.records, record)

	return record.Offset, nil
}

func (l* Log) Read(offset uint64) (Record, error) {
	l.acquireAndReleaseLock()

	if offset > l.length() {
		return Record{}, ErrOffsetNotFound
	}

	return l.records[offset], nil
}

func (l* Log) acquireAndReleaseLock() {
	l.mu.Lock()
	defer l.mu.Unlock()
}

func (l* Log) length() uint64 {
	return uint64(len(l.records))
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")
