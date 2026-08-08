package main

import (
	"os/exec"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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
	labelPrefix := cases.Title(language.Und).String(p.MetricKeyPrefix())

	graph := func(label, unit string, metrics ...mp.Metrics) mp.Graphs {
		return mp.Graphs{
			Label:   labelPrefix + ": " + label,
			Unit:    unit,
			Metrics: metrics,
		}
	}
	diffMetric := func(name, label string) mp.Metrics {
		return mp.Metrics{Name: name, Label: label, Diff: true}
	}
	stackedDiffMetric := func(name, label string) mp.Metrics {
		return mp.Metrics{Name: name, Label: label, Stacked: true, Diff: true}
	}

	return map[string]mp.Graphs{
		"dnsupdate": graph("Dynamic DNS Update", "integer",
			diffMetric("dnsupdate-answers", "Answers"),
			diffMetric("dnsupdate-changes", "Changes"),
			diffMetric("dnsupdate-queries", "Queries"),
			diffMetric("dnsupdate-refused", "Refused"),
		),
		"notifications": graph("DNS Notifications", "integer",
			diffMetric("incoming-notifications", "Incoming"),
		),
		"packetcache": graph("Packet Cache", "integer",
			stackedDiffMetric("packetcache-hit", "Hits"),
			stackedDiffMetric("packetcache-miss", "Misses"),
		),
		"query-cache": graph("Query Cache", "integer",
			stackedDiffMetric("query-cache-hit", "Hits"),
			stackedDiffMetric("query-cache-miss", "Misses"),
		),
		"cache-size": graph("Cache Sizes", "integer",
			mp.Metrics{Name: "packetcache-size", Label: "Packet cache"},
			mp.Metrics{Name: "key-cache-size", Label: "Key cache"},
			mp.Metrics{Name: "signature-cache-size", Label: "Signature cache"},
			mp.Metrics{Name: "meta-cache-size", Label: "Metadata cache"},
		),
		"fails": graph("Failed packets", "integer",
			diffMetric("servfail-packets", "SERVFAIL packets"),
			diffMetric("corrupt-packets", "Corrupt packets"),
			diffMetric("timedout-packets", "Timedout packets"),
			diffMetric("overload-drops", "Dropped because backends overload"),
		),
		"backend": graph("Backend", "integer",
			diffMetric("backend-queries", "Backend queries"),
		),
		"tcp-connection": graph("Connections", "integer",
			mp.Metrics{Name: "open-tcp-connections", Label: "TCP Connections"},
			mp.Metrics{Name: "fd-usage", Label: "FD usage"},
		),
		"signatures": graph("DNSSEC Signatures", "integer",
			diffMetric("signatures", "Signatures created"),
		),
		"latency": graph("Latency (microseconds)", "integer",
			mp.Metrics{Name: "latency", Label: "Latency"},
		),
		"qsize": graph("Queue Size", "integer",
			mp.Metrics{Name: "qsize-q", Label: "Queue size"},
		),
		"answers": graph("Answers", "integer",
			diffMetric("tcp-answers", "TCP"),
			diffMetric("udp-answers", "UDP"),
			stackedDiffMetric("tcp4-answers", "TCP4"),
			stackedDiffMetric("udp4-answers", "UDP4"),
			stackedDiffMetric("tcp6-answers", "TCP6"),
			stackedDiffMetric("udp6-answers", "UDP6"),
		),
		"queries": graph("Queries", "integer",
			diffMetric("tcp-queries", "TCP"),
			diffMetric("udp-queries", "UDP"),
			stackedDiffMetric("tcp4-queries", "TCP4"),
			stackedDiffMetric("udp4-queries", "UDP4"),
			stackedDiffMetric("tcp6-queries", "TCP6"),
			stackedDiffMetric("udp6-queries", "UDP6"),
			diffMetric("udp-do-queries", "UDP DO queries"),
		),
		"answer-bytes": graph("Answer Bytes", "integer",
			diffMetric("tcp-answers-bytes", "TCP"),
			diffMetric("udp-answers-bytes", "UDP"),
			stackedDiffMetric("tcp4-answers-bytes", "TCP4"),
			stackedDiffMetric("udp4-answers-bytes", "UDP4"),
			stackedDiffMetric("tcp6-answers-bytes", "TCP6"),
			stackedDiffMetric("udp6-answers-bytes", "UDP6"),
		),
		"memory": graph("Memory Usage", "bytes",
			mp.Metrics{Name: "real-memory-usage", Label: "Usage"},
		),
		"cpu": graph("CPU Usage (milliseconds)", "integer",
			diffMetric("user-msec", "User"),
			diffMetric("sys-msec", "System"),
		),
	}
}

func (u *Plugin) FetchMetrics() (map[string]float64, error) {
	buf, err := exec.Command(u.CommandPath, "show", "*").Output()
	if err != nil {
		return nil, err
	}
	return u.ParseMetrics(buf), nil
}

func (u *Plugin) ParseMetrics(buf []byte) map[string]float64 {
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
	return result
}

func (u *Plugin) Run() {
	plugin := mp.NewMackerelPlugin(u)
	plugin.Run()
}
