package main

import (
	"bytes"
	"time"

	// "crypto/rand"
	"fmt"
	"net/http"

	"github.com/go-petname/petname"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)
var (
	reg = prometheus.NewRegistry()
)


func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/bar", barHandler)
	http.ListenAndServe(":8080", nil)
}

func registerCounter(name, help string, tagList map[string]string) {
	labels := []string{}
	values := []string{}
	
	for label, value :=range tagList{
		labels = append(labels, label)
		values = append(values, value)
	}
	// counterName := fmt.Sprintf("%s_%s_%s",name, "a","b")
	counterName := fmt.Sprintf("%s_%s_%s",name, petname.Generate(3,"_"), time.Now().Format("20060102150405"))
	counter := prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: counterName,
		Help: help,
	},
	labels,    
	)
	counter.WithLabelValues(values...).Inc()
	// counter.
	prometheus.Register(counter)
	// c := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
	// 	Name: counterName,
	// 	Help: help,
	// }, labels)
	// 	c.WithLabelValues(values...).Inc()
	// prometheus.NewGaugeVec(prometheus.GaugeOpts{}, []string{})
}





func barHandler(w http.ResponseWriter, r *http.Request) {
	t := map[string]string{
		"name":"foo",
	}
	// if rand.Int()%2 == 0 {
	// 	delete(t, "name")
	// 	t["key"]="bar"
	// }

	defer func () {
		registerCounter("bar_counter","http bar count",t)
	}()

	
	fmt.Println("barHandler status:", t)

	bs := bytes.NewBufferString("success")
	w.Write(bs.Bytes())
}





