package main

import (
	"fmt"
	"time"
)

type EntryType int

const (
	PATH EntryType = 0
	CMD  EntryType = 1
)

type Entry struct {
	id        int64
	prio      int64
	text      string
	timestamp *time.Time
	etype     EntryType
}

func (e EntryType) String() string {
	switch e {
	case PATH:
		return "PATH"
	case CMD:
		return "COMMAND"
	default:
		return "UNKNOWN"
	}
}

func (e Entry) String() string {
	return fmt.Sprintf("%v (%v)", e.text, e.etype)
}
