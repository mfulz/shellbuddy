package main

import (
	"fmt"
	"github.com/mitchellh/go-ps"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ShellType int

const (
	ZSH  ShellType = 0
	BASH ShellType = 1
)

type ShellConfig struct {
	shellType   ShellType
	historyFile string
}

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

func randomIdentifier() string {
	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, 20)
	for i := 0; i < 20; i++ {
		bytes[i] = byte(65 + rand.Intn(25))
	}

	return string(bytes)
}

func getHomeDir() (string, error) {
	ret, isset := os.LookupEnv("HOME")
	if isset {
		return ret, nil
	}

	return "", fmt.Errorf("Environment variable 'HOME' is not set. Cannot identify user home directory")
}

func getShellProcess() (string, error) {
	ppid := os.Getppid()
	parent, err := ps.FindProcess(ppid)
	if err != nil {
		return "", err
	}

	return parent.Executable(), nil
}

func retrieveShellCmd(cmd string) (string, error) {
	p, err := getShellProcess()
	if err != nil {
		return "", err
	}

	identifier := randomIdentifier()

	out, err := exec.Command(p, "-ic", "echo "+identifier+" && "+cmd).Output()
	if err != nil {
		return "", err
	}

	outl := strings.Split(string(out), "\n")
	if outl[len(outl)-1] == "" {
		outl = outl[:len(outl)-1]
	}
	var outs string
	start := false

	for i, l := range outl {
		if l == identifier {
			start = true
			continue
		}

		if !start {
			continue
		}

		if i < len(outl)-1 {
			outs = outs + l + "\n"
		} else {
			outs = outs + l
		}
	}

	return outs, nil
}

func getShellType() (ShellType, error) {
	out, err := retrieveShellCmd("echo $0")
	if err != nil {
		return -1, err
	}

	return ShellToType(out)
}

func getShellHistoryFile() (string, error) {
	out, err := retrieveShellCmd("echo $HISTFILE")
	if err != nil {
		return "", err
	}

	if out == "" {
		return out, fmt.Errorf("Unable to retrieve history file path")
	}

	return out, nil
}

func GetShellConfig() (ShellConfig, error) {
	ret := ShellConfig{}
	var err error

	ret.shellType, err = getShellType()
	if err != nil {
		return ret, err
	}

	ret.historyFile, err = getShellHistoryFile()
	if err != nil {
		return ret, err
	}

	return ret, nil
}
