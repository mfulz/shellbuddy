package main

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Config struct {
	DBPath            string
	Timezone          string
	MaxEntries        int
	VimMode           bool
	IgnoreFromHistory []string
}

type appRuntime struct {
	db                *DB
	location          *time.Location
	add               bool
	etype             []EntryType
	search            string
	maxEntries        int
	vimMode           bool
	shellConfig       ShellConfig
	ignoreFromHistory []string
}

func initRuntime(a *args) (*appRuntime, error) {
	ret := new(appRuntime)
	ret.etype = nil
	config := Config{}
	config.Timezone = "Europe/Berlin"
	config.VimMode = false
	config.MaxEntries = 10

	cfile, err := ioutil.ReadFile(a.config)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(cfile, &config)
	if err != nil {
		return nil, err
	}

	ret.db, err = openDB(config.DBPath)
	if err != nil {
		return nil, err
	}

	ret.location, err = time.LoadLocation(config.Timezone)
	if err != nil {
		return nil, err
	}

	ret.add = a.add

	if a.etype != "" {
		etypesl := strings.Split(a.etype, ",")
		if len(etypesl) > 0 {
			ret.etype = make([]EntryType, len(etypesl))
			for i, e := range etypesl {
				etype, err := StringToEntryType(e)
				if err != nil {
					return nil, err
				}
				ret.etype[i] = etype
			}
		}
	}

	ret.search = a.search
	ret.maxEntries = config.MaxEntries
	ret.vimMode = config.VimMode

	if a.shellType == "" {
		ret.shellConfig.shellType, err = getShellType()
		if err != nil {
			return nil, err
		}
	} else {
		ret.shellConfig.shellType, err = ShellToType(a.shellType)
		if err != nil {
			return nil, err
		}
	}

	ret.shellConfig.historyFile = a.historyFile
	if ret.shellConfig.historyFile == "" {
		ret.shellConfig.historyFile, err = getShellHistoryFile()
		if err != nil {
			return nil, err
		}
	}

	ret.ignoreFromHistory = config.IgnoreFromHistory

	return ret, nil
}

func (r *appRuntime) GetEntries() ([]Entry, error) {
	if len(r.etype) > 0 {
		return r.db.GetAllEntriesByType(r.search, r.etype)
	} else {
		return r.db.GetAllEntries(r.search)
	}

	return nil, fmt.Errorf("Something went wrong")
}

func (r *appRuntime) AddEntries() error {
	p, err := os.Getwd()
	if err != nil {
		return err
	}

	err = r.db.AddPath(p)
	if err != nil {
		return err
	}

	history, err := ioutil.ReadFile(r.shellConfig.historyFile)
	if err != nil {
		return err
	}

	historyl := strings.Split(string(history), "\n")
	if len(historyl) > 2 {
		// remove last entry if empty
		if historyl[len(historyl)-1] == "" {
			historyl = historyl[0 : len(historyl)-1]
		}

		switch r.shellConfig.shellType {
		case ZSH:
			for i, c := range historyl {
				csplit := strings.SplitN(c, ";", 2)
				if len(csplit) == 2 {
					historyl[i] = csplit[1]
				}
			}
		case BASH:
			break
		default:
			return fmt.Errorf("Unknown shell: %s", r.shellConfig.shellType)
		}

		var cmd string
		index := len(historyl) - 1
		for {
			ignore := false
			cmd = historyl[index]
			for _, i := range r.ignoreFromHistory {
				if i == strings.SplitN(cmd, " ", 2)[0] {
					ignore = true
					break
				}
			}
			if ignore {
				if index == 0 {
					return nil
				}
				index--
			} else {
				break
			}
		}
		err = r.db.AddCmd(cmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *appRuntime) Run() ([]Entry, error) {
	defer r.db.Close()
	if r.add {
		return []Entry{}, r.AddEntries()
	}

	ret, err := r.GetEntries()
	if err != nil {
		return nil, err
	}

	for _, e := range ret {
		*e.Timestamp = e.Timestamp.In(r.location)
	}

	return ret, nil
}
