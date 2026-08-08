package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sampleOutput = `backend-queries=5370225,corrupt-packets=28341,deferred-cache-inserts=2328,deferred-cache-lookup=661,deferred-packetcache-inserts=10031595,deferred-packetcache-lookup=379441,dnsupdate-answers=0,dnsupdate-changes=0,dnsupdate-queries=0,dnsupdate-refused=0,incoming-notifications=0,noerror-packets=292629,nxdomain-packets=110701,overload-drops=0,packetcache-hit=0,packetcache-miss=146921962,packetcache-size=0,query-cache-hit=31768384,query-cache-miss=5369564,query-cache-size=14,rd-queries=144656754,recursing-answers=6387867,recursing-questions=0,recursion-unanswered=0,security-status=3,servfail-packets=10286,signatures=0,tcp-answers=1608086,tcp-answers-bytes=95156373,tcp-cookie-queries=0,tcp-queries=1608148,tcp4-answers=1608086,tcp4-answers-bytes=95156373,tcp4-queries=1608148,tcp6-answers=0,tcp6-answers-bytes=0,tcp6-queries=0,timedout-packets=0,udp-answers=145713583,udp-answers-bytes=5379923337,udp-cookie-queries=0,udp-do-queries=837540,udp-queries=145713585,udp4-answers=139325716,udp4-answers-bytes=5379923337,udp4-queries=145713585,udp6-answers=0,udp6-answers-bytes=0,udp6-queries=0,unauth-packets=133545839,zone-cache-hit=13633551,zone-cache-miss=541012914,zone-cache-size=51,backend-latency=0,cache-latency=2,cpu-iowait=2567089,cpu-steal=151944,fd-usage=46,key-cache-size=386,latency=32,meta-cache-size=1021,open-tcp-connections=0,qsize-q=0,real-memory-usage=23306240,receive-latency=78,ring-logmessages-capacity=10000,ring-logmessages-size=0,ring-noerror-queries-capacity=10000,ring-noerror-queries-size=0,ring-nxdomain-queries-capacity=10000,ring-nxdomain-queries-size=0,ring-queries-capacity=10000,ring-queries-size=0,ring-remotes-capacity=10000,ring-remotes-corrupt-capacity=10000,ring-remotes-corrupt-size=0,ring-remotes-size=0,ring-remotes-unauth-capacity=10000,ring-remotes-unauth-size=0,ring-servfail-queries-capacity=10000,ring-servfail-queries-size=0,ring-unauth-queries-capacity=10000,ring-unauth-queries-size=0,send-latency=17,signature-cache-size=0,sys-msec=5563305,udp-in-csum-errors=1581,udp-in-errors=1581,udp-noport-errors=318562,udp-recvbuf-errors=0,udp-sndbuf-errors=0,udp6-in-csum-errors=0,udp6-in-errors=0,udp6-noport-errors=0,udp6-recvbuf-errors=0,udp6-sndbuf-errors=0,uptime=31900540,user-msec=6273487,xfr-queue=0,`

var expectedMetrics = map[string]float64{
	"backend-queries":              5370225,
	"corrupt-packets":              28341,
	"deferred-cache-inserts":       2328,
	"deferred-cache-lookup":        661,
	"deferred-packetcache-inserts": 10031595,
	"deferred-packetcache-lookup":  379441,
	"backend-latency":              0,
}

func TestParseMetrics(t *testing.T) {
	p := &Plugin{
		Prefix:      "pdns",
		CommandPath: "/usr/bin/pdns_control",
	}

	metrics := p.ParseMetrics([]byte(sampleOutput))

	require.Greater(t, len(metrics), 0, "ParseMetrics returned empty metrics")

	for name, expectedValue := range expectedMetrics {
		value, ok := metrics[name]
		assert.True(t, ok, "Metric %s not found in parsed metrics", name)
		assert.Equal(t, expectedValue, value, "Metric %s has unexpected value", name)
	}

}
