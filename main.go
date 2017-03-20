package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddress = flag.String("web.listen", ":9119", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.path", "/metrics", "Path under which to expose metrics.")
	namespace     = flag.String("namespace", "shield", "Namespace for the Shield failed backups metrics.")

	shieldBackend = flag.String("shield.Backend", "", "Shield Url to connect to.")
	shieldUser    = flag.String("shield.User", "shield", "User to connect to Shield instance.")
	shieldPass    = flag.String("shield.Pass", "", "Password to connect to Shield instance.")
)

func main() {
	flag.Parse()

	prometheus.MustRegister(NewExporter(
		*namespace,
		*shieldBackend,
		*shieldUser,
		*shieldPass))

	log.Printf("Starting Server: %s", *listenAddress)
	handler := promhttp.Handler()
	if *metricsPath == "" || *metricsPath == "/" {
		http.Handle(*metricsPath, handler)
	} else {
		http.Handle(*metricsPath, handler)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
			<head><title>Shield Exporter</title></head>
			<body>
			<h1>Shield Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
		})
	}

	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
