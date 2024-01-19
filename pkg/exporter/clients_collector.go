package exporter

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/kubernetes-sigs/nfs-ganesha-server-and-external-provisioner/pkg/dbus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	clientsNfsV41RequestedDesc = prometheus.NewDesc(
		"ganesha_clients_nfs_v41_requested_bytes_total",
		"Number of requested bytes for NFSv4.1 operations",
		[]string{"direction", "clientip"}, nil,
	)
	clientsNfsV41TransferedDesc = prometheus.NewDesc(
		"ganesha_clients_nfs_v41_transfered_bytes_total",
		"Number of transfered bytes for NFSv4.1 operations",
		[]string{"direction", "clientip"}, nil,
	)
	clientsNfsV41OperationsDesc = prometheus.NewDesc(
		"ganesha_clients_nfs_v41_operations_total",
		"Number of operations for NFSv4.1",
		[]string{"direction", "clientip"}, nil,
	)
	clientsNfsV41ErrorsDesc = prometheus.NewDesc(
		"ganesha_clients_nfs_v41_operations_errors_total",
		"Number of operations in error for NFSv4.1",
		[]string{"direction", "clientip"}, nil,
	)
	clientsNfsV41LatencyDesc = prometheus.NewDesc(
		"ganesha_clients_nfs_v41_operations_latency_seconds_total",
		"Cumulative time consumed by operations for NFSv4.1",
		[]string{"direction", "clientip"}, nil,
	)
	clientsNfsV41QueueWaitDesc = prometheus.NewDesc(
		"ganesha_clients_nfs_v41_operations_queue_wait_seconds_total",
		"Cumulative time spent in rpc wait queue for NFSv4.1",
		[]string{"direction", "clientip"}, nil,
	)
)

// ClientsCollector Collector for ganesha clients
type ClientsCollector struct {
	clientMgr              dbus.ClientMgr
	nfsv40, nfsv41, nfsv42 bool
}

// NewClientsCollector creates a new collector
func NewClientsCollector(v40, v41, v42 *bool) ClientsCollector {
	return ClientsCollector{
		clientMgr: dbus.NewClientMgr(),
		nfsv40:    *v40,
		nfsv41:    *v41,
		nfsv42:    *v42,
	}
}

// Describe prometheus description
func (ic ClientsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(ic, ch)
}

// Collect do the actual job
func (ic ClientsCollector) Collect(ch chan<- prometheus.Metric) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}

		log.Println("get stats panic")
		spew.Dump(recovered)
	}()
	_, clients := ic.clientMgr.ShowClients()
	for _, client := range clients {
		clientip := client.Client
		if ic.nfsv41 {
			stats := dbus.BasicStats{}
			if client.NFSv41 {
				stats = ic.clientMgr.GetNFSv41IO(client.Client)
			}
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41RequestedDesc,
				prometheus.CounterValue,
				float64(stats.Read.Requested),
				"read", clientip)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41TransferedDesc,
				prometheus.CounterValue,
				float64(stats.Read.Transfered),
				"read", clientip)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41OperationsDesc,
				prometheus.CounterValue,
				float64(stats.Read.Total),
				"read", clientip)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41ErrorsDesc,
				prometheus.CounterValue,
				float64(stats.Read.Errors),
				"read", clientip)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41LatencyDesc,
				prometheus.CounterValue,
				float64(stats.Read.Latency)/1e9,
				"read", clientip)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41QueueWaitDesc,
				prometheus.CounterValue,
				float64(stats.Read.QueueWait)/1e9,
				"read", clientip)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41RequestedDesc,
				prometheus.CounterValue,
				float64(stats.Write.Requested),
				"write", clientip)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41TransferedDesc,
				prometheus.CounterValue,
				float64(stats.Write.Transfered),
				"write", clientip)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41OperationsDesc,
				prometheus.CounterValue,
				float64(stats.Write.Total),
				"write", clientip)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41ErrorsDesc,
				prometheus.CounterValue,
				float64(stats.Write.Errors),
				"write", clientip)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41LatencyDesc,
				prometheus.CounterValue,
				float64(stats.Write.Latency)/1e9,
				"write", clientip)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41QueueWaitDesc,
				prometheus.CounterValue,
				float64(stats.Write.QueueWait)/1e9,
				"write", clientip)
		}
	}
}
