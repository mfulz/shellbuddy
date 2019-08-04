package main

import (
	"flag"
	"fmt"
	"os"
)

type args struct {
	config string
	add    bool
	path   bool
	cmd    bool
}

func parseArgs() (*args, error) {
	ret := new(args)

	flag.StringVar(&ret.config, "config", "~/.shellbuddy/config", "Configuration file to use. Defaults to ~/.shellbuddy/config")
	flag.BoolVar(&ret.add, "add", false, "If you want to add / update an entry")
	flag.BoolVar(&ret.cmd, "cmd", false, "If you want to add / update commands")
	flag.BoolVar(&ret.path, "path", false, "If you want to add / update pathes")

	flag.Parse()

	if !ret.cmd && !ret.path {
		return nil, fmt.Errorf("Neither path nor cmd selected")
	}

	return ret, nil
}

func main() {
	db, err := openDB("/home/mfulz/test.db")
	if err != nil {
		panic(err)
	}

	d, _ := os.Getwd()
	err = db.AddPath(d)
	if err != nil {
		panic(err)
	}

	p, _ := db.GetPathes()
	fmt.Println(p)
}
