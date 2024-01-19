package exporter

import (
	"log"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/kubernetes-sigs/nfs-ganesha-server-and-external-provisioner/pkg/dbus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	nfsV41RequestedDesc = prometheus.NewDesc(
		"ganesha_exports_nfs_v41_requested_bytes_total",
		"Number of requested bytes for NFSv4.1 operations",
		[]string{"direction", "exportid", "path"}, nil,
	)
	nfsV41TransferedDesc = prometheus.NewDesc(
		"ganesha_exports_nfs_v41_transfered_bytes_total",
		"Number of transfered bytes for NFSv4.1 operations",
		[]string{"direction", "exportid", "path"}, nil,
	)
	nfsV41OperationsDesc = prometheus.NewDesc(
		"ganesha_exports_nfs_v41_operations_total",
		"Number of operations for NFSv4.1",
		[]string{"direction", "exportid", "path"}, nil,
	)
	nfsV41ErrorsDesc = prometheus.NewDesc(
		"ganesha_exports_nfs_v41_operations_errors_total",
		"Number of operations in error for NFSv4.1",
		[]string{"direction", "exportid", "path"}, nil,
	)
	nfsV41LatencyDesc = prometheus.NewDesc(
		"ganesha_exports_nfs_v41_operations_latency_seconds_total",
		"Cumulative time consumed by operations for NFSv4.1",
		[]string{"direction", "exportid", "path"}, nil,
	)
	nfsV41QueueWaitDesc = prometheus.NewDesc(
		"ganesha_exports_nfs_v41_operations_queue_wait_seconds_total",
		"Cumulative time spent in rpc wait queue for NFSv4.1",
		[]string{"direction", "exportid", "path"}, nil,
	)
)

// ExportsCollector Collector for ganesha exports
type ExportsCollector struct {
	exportMgr              dbus.ExportMgr
	nfsv40, nfsv41, nfsv42 bool
}

var t = true

// NewExportsCollector creates a new collector
func NewExportsCollector(v40, v41, v42 *bool) ExportsCollector {
	return ExportsCollector{
		exportMgr: dbus.NewExportMgr(),
		nfsv40:    *v40,
		nfsv41:    *v41,
		nfsv42:    *v42,
	}
}

func (ic ExportsCollector) GetExportMgr() dbus.ExportMgr {
	return ic.exportMgr
}

// Describe prometheus description
func (ic ExportsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(ic, ch)
}

// Collect do the actual job
func (ic ExportsCollector) Collect(ch chan<- prometheus.Metric) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}

		log.Println("get stats panic")
		spew.Dump(recovered)
	}()
	_, exports := ic.exportMgr.ShowExports()
	for _, export := range exports {
		exportid := strconv.FormatUint(uint64(export.ExportID), 10)
		path := export.Path
		if ic.nfsv41 {
			stats := dbus.BasicStats{}
			if export.NFSv41 {
				stats = ic.exportMgr.GetNFSv41IO(export.ExportID)
			}
			ch <- prometheus.MustNewConstMetric(
				nfsV41RequestedDesc,
				prometheus.CounterValue,
				float64(stats.Read.Requested),
				"read", exportid, path)
			ch <- prometheus.MustNewConstMetric(
				nfsV41TransferedDesc,
				prometheus.CounterValue,
				float64(stats.Read.Transfered),
				"read", exportid, path)
			ch <- prometheus.MustNewConstMetric(
				nfsV41OperationsDesc,
				prometheus.CounterValue,
				float64(stats.Read.Total),
				"read", exportid, path)
			ch <- prometheus.MustNewConstMetric(
				nfsV41ErrorsDesc,
				prometheus.CounterValue,
				float64(stats.Read.Errors),
				"read", exportid, path)
			ch <- prometheus.MustNewConstMetric(
				nfsV41LatencyDesc,
				prometheus.CounterValue,
				float64(stats.Read.Latency)/1e9,
				"read", exportid, path)
			ch <- prometheus.MustNewConstMetric(
				nfsV41QueueWaitDesc,
				prometheus.CounterValue,
				float64(stats.Read.QueueWait)/1e9,
				"read", exportid, path)
			ch <- prometheus.MustNewConstMetric(
				nfsV41RequestedDesc,
				prometheus.CounterValue,
				float64(stats.Write.Requested),
				"write", exportid, path)
			ch <- prometheus.MustNewConstMetric(
				nfsV41TransferedDesc,
				prometheus.CounterValue,
				float64(stats.Write.Transfered),
				"write", exportid, path)
			ch <- prometheus.MustNewConstMetric(
				nfsV41OperationsDesc,
				prometheus.CounterValue,
				float64(stats.Write.Total),
				"write", exportid, path)
			ch <- prometheus.MustNewConstMetric(
				nfsV41ErrorsDesc,
				prometheus.CounterValue,
				float64(stats.Write.Errors),
				"write", exportid, path)
			ch <- prometheus.MustNewConstMetric(
				nfsV41LatencyDesc,
				prometheus.CounterValue,
				float64(stats.Write.Latency)/1e9,
				"write", exportid, path)
			ch <- prometheus.MustNewConstMetric(
				nfsV41QueueWaitDesc,
				prometheus.CounterValue,
				float64(stats.Write.QueueWait)/1e9,
				"write", exportid, path)
		}
	}
}
