package main

import (
	"fmt"
	"os"
	"path"
)

func writeConfig(a *args) error {
	configPath := path.Dir(a.config)

	_, err := os.Stat(a.config)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		return fmt.Errorf("Config file already existing")
	}

	_, err = os.Stat(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err = os.MkdirAll(configPath, os.FileMode(0755))
		if err != nil {
			return err
		}
	}

	cf, err := os.Create(a.config)
	if err != nil {
		return err
	}
	defer cf.Close()

	h, isset := os.LookupEnv("HOME")
	if !isset {
		return fmt.Errorf("Cannot retrieve home directory")
	}

	line := "# DBPath set this to the location of the sqlite database file\n"
	if _, err := cf.WriteString(line); err != nil {
		return err
	}
	line = "DBPath = \"" + path.Join(h, ".shellbuddy", "shellbuddy.db") + "\"\n\n"
	if _, err := cf.WriteString(line); err != nil {
		return err
	}

	line = "# Timezone set this to your desired timezone\n"
	if _, err := cf.WriteString(line); err != nil {
		return err
	}
	line = "Timezone = \"Europe/Berlin\"\n\n"
	if _, err := cf.WriteString(line); err != nil {
		return err
	}

	line = "# MaxEntries set this to the maximum Entries that should be displayed in one selection of the prompt\n"
	if _, err := cf.WriteString(line); err != nil {
		return err
	}
	line = "MaxEntries = 5\n\n"
	if _, err := cf.WriteString(line); err != nil {
		return err
	}

	line = "# VimMode set this to true to enable vim style keybindings\n"
	if _, err := cf.WriteString(line); err != nil {
		return err
	}
	line = "VimMode = false\n\n"
	if _, err := cf.WriteString(line); err != nil {
		return err
	}

	line = "# IgnoreFromHistory set this to a list for commands that should be ignored from the history file. Normally the functions to use shellbuddy\n"
	if _, err := cf.WriteString(line); err != nil {
		return err
	}
	line = "IgnoreFromHistory = [\"h\", \"c\"]\n\n"
	if _, err := cf.WriteString(line); err != nil {
		return err
	}

	return nil
}
