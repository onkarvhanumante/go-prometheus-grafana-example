package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-petname/petname"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

func main() {
	fmt.Println("some job .....")
	sleepDuration, e := strconv.Atoi(os.Args[1])
	fmt.Println("sleeping for ", sleepDuration, " seconds", e)
	start := time.Now()
	time.Sleep(time.Second * time.Duration(sleepDuration))
	elapsed := time.Since(start)
	fmt.Println("completed job in ", elapsed, " seconds")
	pushGauge("someFooJob", "someFooHelp", map[string]string{
		"dev":"yes",
	}, elapsed.Seconds())
}


func pushGauge(name, help string, tagList map[string]string, val float64) {
	labels := []string{}
	values := []string{}
	
	for label, value :=range tagList{
		labels = append(labels, label)
		values = append(values, value)
	}

	gaugeName := fmt.Sprintf("%s_%s_%s",name, petname.Generate(3,"_"), time.Now().Format("20060102150405"))
	completionTime := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: gaugeName,
			Help: help,
		},
		labels,
	)
	completionTime.WithLabelValues(values...).Set(val)

	registry := prometheus.NewRegistry()
	registry.MustRegister(completionTime)

	pusher := push.New("http://localhost:9091", "db_backup").Gatherer(registry)
	if err := pusher.Add(); err != nil {
		fmt.Println("Could not push to Pushgateway:", err)
	} else {
		fmt.Println("Success")
	}
}
