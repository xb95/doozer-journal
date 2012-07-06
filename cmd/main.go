// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/doozer-journal

package main

import (
	"flag"
	"fmt"
	"github.com/soundcloud/doozer"
	"os"
	"text/template"
)

type Command struct {
	Run       func(cmd *Command, args []string)
	Desc      string
	Name      string
	UsageLine string
	Conn      *doozer.Conn
}

func (cmd *Command) Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s\n\n", cmd.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", cmd.Desc)

	os.Exit(1)
}

var (
	file string
	uri  string
)

var commands = []*Command{
	cmdJournal,
	cmdRestore,
	cmdSnapshot,
}

func init() {
	flag.Usage = usage
	flag.StringVar(&file, "file", "./doozerd.log", "location of the journal file")
	flag.StringVar(&uri, "uri", "doozer:?ca=localhost:8046", "doozerd cluster uri")
	flag.Parse()
}

func main() {
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
	}

	for _, cmd := range commands {
		if cmd.Name == args[0] && cmd.Run != nil {
			conn, err := doozer.DialUri(uri, "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error connecting to %s: %s\n", uri, err.Error())
				os.Exit(1)
			}

			cmd.Conn = conn
			cmd.Run(cmd, args)

			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command %#q\n\n", flag.Args()[0])
	flag.Usage()
}

func usage() {
	t := template.New("top")
	template.Must(t.Parse(usageTmpl))
	data := struct {
		Commands []*Command
		Globals  map[string]string
	}{
		commands,
		map[string]string{"file": file, "uri": uri},
	}

	if err := t.Execute(os.Stderr, data); err != nil {
		panic(err)
	}

	os.Exit(1)
}

var usageTmpl = `Usage: doozer-journal [globals] command

Globals:
  -file   location of backup file ({{.Globals.file}})
  -uri    doozerd cluster URI     ({{.Globals.uri}})

Commands:{{range .Commands}}
  {{.Name | printf "%-10s"}} {{.Desc}}{{end}}
`
