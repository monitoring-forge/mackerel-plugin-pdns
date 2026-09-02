package main

import (
	"os"

	"github.com/monitoring-forge/flagrun"
)

var version string

type Opt struct {
	Version     bool   `short:"v" long:"version" description:"Show version"`
	Prefix      string `long:"prefix" default:"pdns" description:"Metric key prefix"`
	CommandPath string `long:"control-command" default:"/usr/bin/pdns_control" description:"Path to pdns_control command"`
}

func (o *Opt) Run(_ []string) {
	u := &Plugin{
		Prefix:      o.Prefix,
		CommandPath: o.CommandPath,
	}
	u.Run()
}

func main() {
	os.Exit(flagrun.Ship(&Opt{}, flagrun.Version(version)))
}
