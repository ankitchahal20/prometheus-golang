package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

// NewCounterVec() method is used to create a new counter metric
var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func init() {
	// to expose the created metric in the HTTP handler we must register the metric to Prometheus using the register() method.
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}

func PostPrintHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Reached in PostPrintHandler")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	s := "success"
	w.Write([]byte(s))
}

func main() {
	router := mux.NewRouter()
	router.Use(prometheusMiddleware)

	// Prometheus endpoint
	router.Path("/prometheus").Handler(promhttp.Handler())

	// Serving static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	//router.HandleFunc("/", PostPrintHandler).Methods("POST")
	fmt.Println("Serving requests on port 9000")
	err := http.ListenAndServe(":9000", router)
	log.Fatal(err)
}
