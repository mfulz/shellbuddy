package main

import (
	"flag"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type ShellType int

const (
	ZSH  ShellType = 0
	BASH ShellType = 1
)

func (s ShellType) String() string {
	switch s {
	case ZSH:
		return "zsh"
	case BASH:
		return "bash"
	default:
		return "unknown"
	}
}

func ShellToType(t string) (ShellType, error) {
	if t == "zsh" {
		return ZSH, nil
	}

	if t == "bash" {
		return BASH, nil
	}

	return -1, fmt.Errorf("Unknown shell type: %s", t)
}

type args struct {
	config string
	add    bool
	path   bool
	cmd    bool
	search string
	init   bool
}

type Config struct {
	DBPath      string
	Timezone    string
	MaxEntries  int
	VimMode     bool
	ShellType   string
	HistoryFile string
}

type appRuntime struct {
	db          *DB
	location    *time.Location
	add         bool
	path        bool
	cmd         bool
	search      string
	maxEntries  int
	vimMode     bool
	shellType   ShellType
	historyFile string
}

func initRuntime(a *args) (*appRuntime, error) {
	ret := new(appRuntime)
	config := Config{}
	config.VimMode = false
	config.MaxEntries = 10
	config.ShellType = "bash"
	config.HistoryFile = ".bash_history"

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
	ret.cmd = a.cmd
	ret.path = a.path
	ret.search = a.search
	ret.maxEntries = config.MaxEntries
	ret.vimMode = config.VimMode
	ret.shellType, err = ShellToType(config.ShellType)
	if err != nil {
		return ret, err
	}
	ret.historyFile = config.HistoryFile

	return ret, nil
}

func (r *appRuntime) GetEntries() ([]Entry, error) {
	if r.path && r.cmd {
		return r.db.GetAllEntries(r.search)
	}

	if r.path {
		return r.db.GetPathes(r.search)
	}

	if r.cmd {
		return r.db.GetCmds(r.search)
	}

	return nil, fmt.Errorf("Something went wrong")
}

func (r *appRuntime) AddEntries() error {
	if r.path {
		p, err := os.Getwd()
		if err != nil {
			return err
		}

		err = r.db.AddPath(p)
		if err != nil {
			return err
		}
	}

	if r.cmd {
		history, err := ioutil.ReadFile(r.historyFile)
		if err != nil {
			return err
		}

		historyl := strings.Split(string(history), "\n")
		if len(historyl) > 2 {
			var cmd string
			switch r.shellType {
			case ZSH:
				cmds := historyl[len(historyl)-2]
				cmd = strings.SplitAfterN(cmds, ";", 2)[1]
				if cmd == "h" {
					cmds = historyl[len(historyl)-3]
					cmd = strings.SplitAfterN(cmds, ";", 2)[1]
				}
			case BASH:
				cmds := historyl[len(historyl)-2]
				cmd = cmds
			default:
				return fmt.Errorf("Unknown shell: %s", r.shellType)
			}
			err = r.db.AddCmd(cmd)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *appRuntime) Run() ([]Entry, error) {
	if r.add {
		return []Entry{}, r.AddEntries()
	}

	return r.GetEntries()
}

func parseArgs() (*args, error) {
	ret := new(args)

	h, isset := os.LookupEnv("HOME")
	if isset {
		h = path.Join(h, ".shellbuddy", "config")
	} else {
		h = "./config"
	}

	flag.StringVar(&ret.config, "config", h, "Configuration file to use. Defaults to ~/.shellbuddy/config")
	flag.BoolVar(&ret.add, "add", false, "If you want to add / update an entry")
	flag.BoolVar(&ret.cmd, "cmd", false, "If you want to add / update commands")
	flag.BoolVar(&ret.path, "path", false, "If you want to add / update pathes")
	flag.StringVar(&ret.search, "search", "", "Select entries by string")
	flag.BoolVar(&ret.init, "init", false, "If you want to create a default config file")

	flag.Parse()

	if ret.init {
		err := writeConfig(ret)
		if err != nil {
			return nil, err
		} else {
			os.Exit(0)
		}
	}

	if !ret.cmd && !ret.path {
		return nil, fmt.Errorf("Neither path nor cmd selected")
	}

	return ret, nil
}

func main() {
	a, err := parseArgs()
	if err != nil {
		panic(err)
	}

	r, err := initRuntime(a)
	if err != nil {
		panic(err)
	}

	e, err := r.Run()
	if err != nil {
		panic(err)
	}

	if len(e) == 0 {
		return
	}

	prompt := promptui.Select{
		Label:        "ShellBuddy",
		Items:        e,
		Size:         r.maxEntries,
		HideSelected: true,
		IsVimMode:    r.vimMode,
	}

	i, _, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "%v", e[i].text)
}
