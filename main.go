package main

import (
	"flag"
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"path"
)

type args struct {
	config      string
	add         bool
	etype       string
	search      string
	init        bool
	shellType   string
	historyFile string
	stdin       bool
}

func parseArgs() (*args, error) {
	ret := new(args)

	h, err := getHomeDir()
	if err != nil {
		return nil, err
	}
	h = path.Join(h, ".shellbuddy", "config")

	flag.StringVar(&ret.config, "config", h, "Configuration file to use. Defaults to ~/.shellbuddy/config")
	flag.StringVar(&ret.shellType, "shell", "", "Specify the desired shell (\"bash\" or \"zsh\"). Normally this will be detected automatically.")
	flag.StringVar(&ret.historyFile, "history", "", "Specify the path to the shell's history file. Normally this will be detected automatically.")
	flag.BoolVar(&ret.add, "add", false, "Adding / Updating entries in the database")
	flag.StringVar(&ret.etype, "entries", "", "Select type of entries. Can be provided as comma seperated list (\"dirs,commands\"). If omitted all entries will be used")
	flag.StringVar(&ret.search, "search", "", "Select entries by search string")
	flag.BoolVar(&ret.init, "init", false, "Initialize configuration")
	flag.BoolVar(&ret.stdin, "stdin", false, "Write directly to the shell's input buffer via ioctl")

	flag.Parse()

	if ret.init {
		err := writeConfig(ret)
		if err != nil {
			return nil, err
		} else {
			os.Exit(0)
		}
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

	if a.stdin {
		err = WriteToShellStdin(e[i].text)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "%v", e[i].text)
	}
}
