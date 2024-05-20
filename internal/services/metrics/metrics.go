package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	Name string
	Help string
}

var metricsInfo = []Metrics{
	{Name: "policies_purchased_today_gauge", Help: "Number of policies purchased today."},
	{Name: "authentications_today_gauge", Help: "Number of authentications today."},
}
var metrics = make(map[string]prometheus.Gauge, len(metricsInfo))

func Register() {
	for _, m := range metricsInfo {
		metrics[m.Name] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: m.Name,
				Help: m.Help,
			})
		prometheus.MustRegister(metrics[m.Name])
	}
}
