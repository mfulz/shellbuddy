package main

import (
	"fmt"
	"time"
)

type EntryType int

const (
	DIR     EntryType = 0
	COMMAND EntryType = 1
)

type Entry struct {
	id        int64
	Prio      int64
	Text      string
	Timestamp *time.Time
	Etype     EntryType
}

func (e EntryType) String() string {
	switch e {
	case DIR:
		return "DIR"
	case COMMAND:
		return "COMMAND"
	default:
		return "UNKNOWN"
	}
}

func StringToEntryType(entry string) (EntryType, error) {
	switch entry {
	case "dirs", "DIR":
		return DIR, nil
	case "commands", "COMMAND":
		return COMMAND, nil
	}

	return -1, fmt.Errorf("Unknown EntryType: %s", entry)
}

func (e Entry) String() string {
	return fmt.Sprintf("%v (%v)", e.Text, e.Etype)
}
