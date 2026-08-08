package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var version string
var commit string

const (
	OK = iota
	WARNING
	CRITICAL
	UNKNOWN
)

type Opt struct {
	Version     bool   `short:"v" long:"version" description:"Show version"`
	Prefix      string `long:"prefix" default:"pdns" description:"Metric key prefix"`
	CommandPath string `long:"control-command" default:"/usr/bin/pdns_control" description:"Path to pdns_control command"`
}

type Plugin struct {
	Prefix      string
	CommandPath string
}

func (p *Plugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "pdns"
	}
	return p.Prefix
}

func (p *Plugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := cases.Title(language.Und).String(p.Prefix)
	return map[string]mp.Graphs{
		"dnsupdate": {
			Label: labelPrefix + ": Dynamic DNS Update",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "dnsupdate-answers", Label: "Answers", Diff: true},
				{Name: "dnsupdate-changes", Label: "Changes", Diff: true},
				{Name: "dnsupdate-queries", Label: "Queries", Diff: true},
				{Name: "dnsupdate-refused", Label: "Refused", Diff: true},
			},
		},
		"notifications": {
			Label: labelPrefix + ": DNS Notifications",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "incoming-notifications", Label: "Incoming", Diff: true},
			},
		},
		"packetcache": {
			Label: labelPrefix + ": Packet Cache",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "packetcache-hit", Label: "Hits", Stacked: true, Diff: true},
				{Name: "packetcache-miss", Label: "Misses", Stacked: true, Diff: true},
			},
		},
		"query-cache": {
			Label: labelPrefix + ": Query Cache",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "query-cache-hit", Label: "Hits", Stacked: true, Diff: true},
				{Name: "query-cache-miss", Label: "Misses", Stacked: true, Diff: true},
			},
		},
		"cache-size": {
			Label: labelPrefix + ": Cache Sizes",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "packetcache-size", Label: "Packet cache"},
				{Name: "key-cache-size", Label: "Key cache"},
				{Name: "signature-cache-size", Label: "Signature cache"},
				{Name: "meta-cache-size", Label: "Metadata cache"},
			},
		},
		"fails": {
			Label: labelPrefix + ": Failed packets",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "servfail-packets", Label: "SERVFAIL packets", Diff: true},
				{Name: "corrupt-packets", Label: "Corrupt packets", Diff: true},
				{Name: "timedout-packets", Label: "Timedout packets", Diff: true},
				{Name: "overload-drops", Label: "Dropped because backends overload", Diff: true},
			},
		},
		"backend": {
			Label: labelPrefix + ": Backend",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "backend-queries", Label: "Backend queries", Diff: true},
			},
		},
		"tcp-connection": {
			Label: labelPrefix + ": Connections",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "open-tcp-connections", Label: "TCP Connections"},
				{Name: "fd-usage", Label: "FD usage"},
			},
		},
		"signatures": {
			Label: labelPrefix + ": DNSSEC Signatures",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "signatures", Label: "Signatures created", Diff: true},
			},
		},
		"latency": {
			Label: labelPrefix + ": Latency (microseconds)",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "latency", Label: "Latency"},
			},
		},
		"qsize": {
			Label: labelPrefix + ": Queue Size",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "qsize-q", Label: "Queue size"},
			},
		},
		"answers": {
			Label: labelPrefix + ": Answers",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "tcp-answers", Label: "TCP", Diff: true},
				{Name: "udp-answers", Label: "UDP", Diff: true},
				{Name: "tcp4-answers", Label: "TCP4", Stacked: true, Diff: true},
				{Name: "udp4-answers", Label: "UDP4", Stacked: true, Diff: true},
				{Name: "tcp6-answers", Label: "TCP6", Stacked: true, Diff: true},
				{Name: "udp6-answers", Label: "UDP6", Stacked: true, Diff: true},
			},
		},
		"queries": {
			Label: labelPrefix + ": Queries",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "tcp-queries", Label: "TCP", Diff: true},
				{Name: "udp-queries", Label: "UDP", Diff: true},
				{Name: "tcp4-queries", Label: "TCP4", Stacked: true, Diff: true},
				{Name: "udp4-queries", Label: "UDP4", Stacked: true, Diff: true},
				{Name: "tcp6-queries", Label: "TCP6", Stacked: true, Diff: true},
				{Name: "udp6-queries", Label: "UDP6", Stacked: true, Diff: true},
				{Name: "udp-do-queries", Label: "UDP DO queries", Diff: true},
			},
		},
		"answer-bytes": {
			Label: labelPrefix + ": Answer Bytes",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "tcp-answers-bytes", Label: "TCP", Diff: true},
				{Name: "udp-answers-bytes", Label: "UDP", Diff: true},
				{Name: "tcp4-answers-bytes", Label: "TCP4", Stacked: true, Diff: true},
				{Name: "udp4-answers-bytes", Label: "UDP4", Stacked: true, Diff: true},
				{Name: "tcp6-answers-bytes", Label: "TCP6", Stacked: true, Diff: true},
				{Name: "udp6-answers-bytes", Label: "UDP6", Stacked: true, Diff: true},
			},
		},
		"memory": {
			Label: labelPrefix + ": Memory Usage",
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "real-memory-usage", Label: "Usage"},
			},
		},
		"cpu": {
			Label: labelPrefix + ": CPU Usage (milliseconds)",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "user-msec", Label: "User", Diff: true},
				{Name: "sys-msec", Label: "System", Diff: true},
			},
		},
	}
}

func (u *Plugin) FetchMetrics() (map[string]float64, error) {
	buf, err := exec.Command(u.CommandPath, "show", "*").Output()
	if err != nil {
		return nil, err
	}
	result := map[string]float64{}
	for b := range strings.SplitSeq(string(buf), ",") {
		kv := strings.SplitN(b, "=", 2)
		if len(kv) != 2 {
			continue
		}
		f, err := strconv.ParseFloat(kv[1], 64)
		if err != nil {
			continue
		}
		result[kv[0]] = f
	}
	return result, nil
}

func (u *Plugin) Run() {
	plugin := mp.NewMackerelPlugin(u)
	plugin.Run()
}

func main() {
	opt := &Opt{}
	psr := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	_, err := psr.Parse()
	if opt.Version {
		if commit == "" {
			commit = "dev"
		}
		fmt.Printf(
			"%s-%s\n%s/%s, %s, %s\n",
			filepath.Base(os.Args[0]),
			version,
			runtime.GOOS,
			runtime.GOARCH,
			runtime.Version(),
			commit)
		os.Exit(OK)
	} else if flags.WroteHelp(err) {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(OK)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(UNKNOWN)
	}

	u := &Plugin{
		Prefix:      opt.Prefix,
		CommandPath: opt.CommandPath,
	}
	u.Run()
}
