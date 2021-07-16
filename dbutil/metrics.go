package dbutil

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	successLabel = "SUCCESS"
	errorLabel   = "ERROR"
)

var (
	dbQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total number of database queries with query name and status",
		},
		[]string{"query", "status"},
	)
	dbQueryLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_latency_ms",
			Help:    "Observed latency of database queries measured in ms",
			Buckets: []float64{1, 5, 10, 20, 50, 100, 200, 500, 1000, 2000, 5000},
		},
		[]string{"query"},
	)
)

// RecordQuerySuccess records the execution of a named query as successfull
func RecordQuerySuccess(name string) {
	dbQueriesTotal.WithLabelValues(name, successLabel).Inc()
}

// RecordQuerySuccess records the execution of a named query as failed
func RecordQueryError(name string) {
	dbQueriesTotal.WithLabelValues(name, errorLabel).Inc()
}

func NewQueryTimer(name string) QueryTimer {
	return QueryTimer{
		QueryName: name,
		startTime: time.Now(),
	}
}

// QueryTimer struct to keep track of a queries execution time.
type QueryTimer struct {
	QueryName string
	startTime time.Time
}

func (t QueryTimer) Stop() {
	latencyMS := float64(time.Now().Sub(t.startTime).Milliseconds())
	dbQueryLatency.WithLabelValues(t.QueryName).Observe(latencyMS)
}
