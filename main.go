package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"fmt"
	"os/exec"
	"runtime"
	"time"
)

func esbackup(timestamp string) {
		out, err := exec.Command("/elasticsearch_backup.sh", timestamp ).Output()
		if err != nil {
			elasticearchbackup.With(prometheus.Labels{"env_label": "app_env"}).Set(0)
			fmt.Println("Command Failed")
			fmt.Println(string(out[:]))
			fmt.Printf("%s", err)
		}
		elasticearchbackup.With(prometheus.Labels{"env_label": "app_env"}).Set(1)
		fmt.Println("Command Successfully Executed")
		output := string(out[:])
		fmt.Println(output)
}

func runbackup() {
	for {
	currentTime := time.Now()
	timestamp := currentTime.Format("2006-01-02_15-04")
	fmt.Println(timestamp)
	cassbackup(timestamp)
	esbackup(timestamp)

    //Sleep for 100 ms
	time.Sleep(60 * time.Minute)
	}
}

var (
    cassandrbackup = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cassandra_backup_status",
			Help: "Cassandra backup job status",
		},
		[]string{
			// Which cluster name where backup is executed
			"app_env",
		})
)

var (
    elasticearchbackup = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "elasticsearch_backup_status",
			Help: "elasticsearch backup job status",
		},
		[]string{
			// Which cluster name where backup is executed
			"app_env",
		})
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(cassandrbackup)
	prometheus.MustRegister(elasticearchbackup)

}

func main() {

	if runtime.GOOS == "windows" {
		fmt.Println("Can't Execute this on a windows machine")
	} else {
		go runbackup()
	}

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":3112", nil)
}
