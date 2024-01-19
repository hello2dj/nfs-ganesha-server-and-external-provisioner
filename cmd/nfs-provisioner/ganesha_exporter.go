package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/kubernetes-sigs/nfs-ganesha-server-and-external-provisioner/pkg/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func runExporter() {
	ec := exporter.NewExportsCollector(nfsv40, nfsv41, nfsv42)
	cc := exporter.NewClientsCollector(nfsv40, nfsv41, nfsv42)

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)
	if *exporterCollector {
		reg.MustRegister(ec)
	}
	if *clientCollector {
		reg.MustRegister(cc)
	}
	http.Handle(*metricsPath, recoverHandler(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))

	log.Println("Listening on", *metricslistenAddress)
	log.Fatalln(http.ListenAndServe(*metricslistenAddress, nil))
}

func recoverHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		// Defer recovery func
		defer func() {
			// Recover potential panic, if this returns nil, everything
			// went as expected and there are no errors to be handled
			recovered := recover()
			if recovered == nil {
				return
			}

			// We expect panic to return a real error,
			// but we can also handle all other types
			var err error
			switch r := recovered.(type) {
			case error:
				err = r
			default:
				err = errors.New("unknown panic reason")
			}

			// Marshal an error response
			jsonBody, err := json.Marshal(map[string]interface{}{
				"error": err.Error(),

				// This adds the current stack information
				// so we can trace back which code panic'd
				"stack": string(debug.Stack()),
			})
			if err != nil {
				return
			}

			// Send JSON error response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(jsonBody)
		}()

		// Invoke wrapped handler func
		next.ServeHTTP(w, request)
	}
}
