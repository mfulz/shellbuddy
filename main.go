package main

import (
	"flag"
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"path"
	"strings"
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
	stdinPrefix string
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
	flag.StringVar(&ret.stdinPrefix, "stdinpre", "", "Add this text to the input buffer before the output")
	flag.BoolVar(&ret.stdin, "stdin", false, "Write directly to the shell's input buffer via ioctl")

	flag.Parse()

	if ret.stdinPrefix != "" && !ret.stdin {
		ret.stdin = true
	}

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

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Text | cyan }} ({{ .Etype | red }})",
		Inactive: "  {{ .Text | cyan }} ({{ .Etype | red }})",
		Selected: "\U0001F336 {{ .Text | red | cyan }}",
		Details: `
--------- Entry ----------
{{ "Text:" | faint }}	{{ .Text }}
{{ "Type:" | faint }}	{{ .Etype }}
{{ "Prio:" | faint }}	{{ .Prio }}
{{ "Last called:" | faint }}	{{ .Timestamp }}`,
	}

	searcher := func(input string, index int) bool {
		entry := e[index]
		name := strings.Replace(strings.ToLower(entry.Text), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:        "ShellBuddy",
		Items:        e,
		Templates:    templates,
		Searcher:     searcher,
		Size:         r.maxEntries,
		HideSelected: true,
		IsVimMode:    r.vimMode,
	}

	i, _, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	if a.stdin {
		var stdinPrefix string
		if a.stdinPrefix != "" {
			stdinPrefix = a.stdinPrefix + " "
		}
		err = WriteToShellStdin(stdinPrefix + e[i].Text)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "%v", e[i].Text)
	}
}
