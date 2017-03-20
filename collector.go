package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"github.com/nimbus-cloud/shield_exporter/tasks"
)

type Exporter struct {
	backend string
	user 	string
	pass 	string

	failedBackups *prometheus.GaugeVec
}

func NewExporter(
	namespace string,
	backend	  string,
	user 	  string,
	pass 	  string,
) *Exporter {
	return &Exporter{
		backend: backend,
		user: user,
		pass: pass,
		failedBackups: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "failed_backup_tasks",
			Help:      "Failed backup tasks in the last 24 hours",
		},
			[]string{"type", "status", "job_name"},
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.failedBackups.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	e.failedBackups.Reset()

	if err := e.collect(); err != nil {
		return
	}

	e.failedBackups.Collect(ch)
}


func (e *Exporter) collect() error {

	tasks, err := tasks.GetFailedBackups(e.backend, e.user, e.pass)
	if err != nil {
		log.Fatalf("Error getting Shield tasks info: %s", err)
	}

	for _, task := range tasks {
		failed := 0.0
		if task.Status == "failed" {
			failed = 1.0
		}
		e.failedBackups.WithLabelValues(task.Op, task.Status, task.JobUUID).Set(failed)
	}

	return nil
}
