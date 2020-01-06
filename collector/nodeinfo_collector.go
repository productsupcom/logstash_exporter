package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// NodeInfoCollector type
type NodeInfoCollector struct {
	endpoint string

	NodeInfos *prometheus.Desc
	OsInfos   *prometheus.Desc
	JvmInfos  *prometheus.Desc
}

// NewNodeInfoCollector function
func NewNodeInfoCollector(logstashEndpoint string) (Collector, error) {
	const subsystem = "info"

	return &NodeInfoCollector{
		endpoint: logstashEndpoint,

		NodeInfos: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "node"),
			"A metric with a constant '1' value labeled by Logstash version.",
			[]string{"version"},
			nil,
		),

		OsInfos: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "os"),
			"A metric with a constant '1' value labeled by name, arch, version and available_processors to the OS running Logstash.",
			[]string{"name", "arch", "version", "available_processors"},
			nil,
		),

		JvmInfos: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "jvm"),
			"A metric with a constant '1' value labeled by name, version and vendor of the JVM running Logstash.",
			[]string{"name", "version", "vendor"},
			nil,
		),
	}, nil
}

// Collect function implements nodestats_collector collector
func (c *NodeInfoCollector) Collect(ch chan<- prometheus.Metric) error {
	stats, err := NodeInfo(c.endpoint)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		c.NodeInfos,
		prometheus.CounterValue,
		float64(1),
		stats.Version,
	)

	ch <- prometheus.MustNewConstMetric(
		c.OsInfos,
		prometheus.CounterValue,
		float64(1),
		stats.Os.Name,
		stats.Os.Arch,
		stats.Os.Version,
		strconv.Itoa(stats.Os.AvailableProcessors),
	)

	ch <- prometheus.MustNewConstMetric(
		c.JvmInfos,
		prometheus.CounterValue,
		float64(1),
		stats.Jvm.VMName,
		stats.Jvm.VMVersion,
		stats.Jvm.VMVendor,
	)

	return nil
}
