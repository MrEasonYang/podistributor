package service

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_request_counter",
		Help: "The total number of http request in podistributor.",
	})
	succeedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_succeed_counter",
		Help: "The total number of succeed http request in podistributor.",
	})
	notFoundCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_not_found_counter",
		Help: "The total number of not found request (404) in podistributor.",
	})
	badRequestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_bad_request_counter",
		Help: "The total number of http bad request (400) in podistributor.",
	})
	internalErrorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_internal_error_counter",
		Help: "The total number of http internal error (500) in podistributor.",
	})
	forbiddenErrorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_forbidden_error_counter",
		Help: "The total number of http forbidden error (403) in podistributor.",
	})
	dbReqCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "db_req_counter",
		Help: "The total number of db request in podistributor.",
	})
	dbNotFoundCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "db_not_found_counter",
		Help: "The total number of records not found in db in podistributor.",
	})
	cachePutCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cache_put_counter",
		Help: "The total number of data putting into cache in podistributor.",
	})
	emptyCachePutCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "empty_cache_put_counter",
		Help: "The total number of empty data putting into cache in podistributor.",
	})
	cachePutFailureCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cache_put_failure_counter",
		Help: "The total number of cache putting action failures in podistributor.",
	})
	analysisHitCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "analysis_hit_counter",
		Help: "The total number of analysis url request in podistributor.",
	})
	analysisFailureCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "analysis_failure_counter",
		Help: "The total number of analysis failures in podistributor.",
	})
)