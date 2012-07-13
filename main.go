// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/doozer-journal

package main

import (
	"flag"
	"fmt"
	"github.com/soundcloud/doozer"
	"github.com/soundcloud/doozer-journal/journal"
	"github.com/soundcloud/logorithm"
	"os"
	"text/template"
)

const VERSION = "0.0.2"

type Command struct {
	Run       func(cmd *Command, args []string)
	Desc      string
	Name      string
	UsageLine string
	Conn      *doozer.Conn
	Rev       int64
}

func (cmd *Command) Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s\n\n", cmd.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", cmd.Desc)

	os.Exit(1)
}

var (
	debug  bool
	file   string
	log    *logorithm.L
	uri    string
	syslog bool
)

var commands = []*Command{
	cmdJournal,
	cmdRestore,
	cmdSnapshot,
}

func init() {
	flag.Usage = usage
	flag.BoolVar(&debug, "d", false, "debug output")
	flag.BoolVar(&syslog, "l", false, "syslog compliant logging")
	flag.StringVar(&file, "file", "./doozerd.log", "location of the journal file")
	flag.StringVar(&uri, "uri", "doozer:?ca=localhost:8046", "doozerd cluster uri")
	flag.Parse()

	if syslog {
		log = logorithm.New(os.Stdout, debug, "doozer-journal", VERSION, "journal", os.Getpid())
	}
}

func main() {
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	for _, cmd := range commands {
		if cmd.Name == args[0] && cmd.Run != nil {
			conn, err := doozer.DialUri(uri, "")
			if err != nil {
				exitWithError("Error connecting to %s: %s\n", uri, err.Error())
			}

			rev, err := conn.Rev()
			if err != nil {
				exitWithError("Unable to get revision: %s\n", err.Error())
			}

			cmd.Conn = conn
			cmd.Rev = rev
			cmd.Run(cmd, args)

			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command %#q\n\n", flag.Args()[0])
	flag.Usage()
	os.Exit(1)
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
}

var usageTmpl = `Usage: doozer-journal [globals] command

Globals:
  -file   location of backup file ({{.Globals.file}})
  -uri    doozerd cluster URI     ({{.Globals.uri}})

Commands:{{range .Commands}}
  {{.Name | printf "%-10s"}} {{.Desc}}{{end}}
`

func snapshot(conn *doozer.Conn, rev int64, j *journal.Journal) (err error) {
	err = doozer.Walk(conn, rev, "/", func(p string, i *doozer.FileInfo, e error) (err error) {
		if e != nil {
			return fmt.Errorf("Error walking tree: %s\n", e.Error())
		}

		if !i.IsDir {
			val, _, err := conn.Get(p, &rev)
			if err != nil {
				return fmt.Errorf("Error getting value for '%s': %s\n", p, err.Error())
			}

			e = j.Append(journal.NewEntry(i.Rev, journal.OpSet, p, val))
			if e != nil {
				return e
			}
		}

		return
	})

	return
}

func exitWithError(msg string, vargs ...interface{}) {
	if log != nil {
		log.Critical(msg, vargs...)
	} else {
		fmt.Fprintf(os.Stderr, msg, vargs...)
	}

	os.Exit(1)
}

func logInfo(msg string, vargs ...interface{}) {
	if log != nil {
		log.Info(msg, vargs...)
	} else {
		fmt.Fprintf(os.Stdout, msg, vargs...)
	}
}
