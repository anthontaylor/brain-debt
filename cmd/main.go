package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegermetrics "github.com/uber/jaeger-lib/metrics"

	brain "github.com/anthontaylor/brain-debt"
	"github.com/anthontaylor/brain-debt/inmem"
	"github.com/anthontaylor/brain-debt/topic"
	//	"github.com/anthontaylor/brain-debt/tracing"
	"github.com/anthontaylor/brain-debt/user"
)

const (
	defaultPort = "8080"
)

func main() {
	var (
		addr       = envString("PORT", defaultPort)
		httpAddr   = flag.String("http.addr", ":"+addr, "HTTP listen address")
		jaegerAddr = flag.String("jaeger", "jaeger:5775", "Jaeger host:port")
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	{
		if *jaegerAddr != "" {
			transport, err := jaeger.NewUDPTransport(*jaegerAddr, 0)
			if err != nil {
				level.Error(logger).Log("err", err)
				os.Exit(1)
			}
			cfg := jaegerconfig.Configuration{
				Sampler: &jaegerconfig.SamplerConfig{
					Type:  jaeger.SamplerTypeConst,
					Param: 1.0,
				},
			}
			closer, err := cfg.InitGlobalTracer(
				"brain_debt",
				jaegerconfig.Logger(logAdapter{logger}),
				jaegerconfig.Metrics(jaegermetrics.NullFactory),
				jaegerconfig.Reporter(jaeger.NewRemoteReporter(transport)),
			)
			if err != nil {
				level.Error(logger).Log("err", err)
				os.Exit(1)
			}
			defer closer.Close()
			level.Info(logger).Log("tracing", "enabled", "jaeger", *jaegerAddr)
		} else {
			level.Info(logger).Log("tracing", "disabled")
		}
	}

	var (
		users  brain.UserRepository
		topics brain.TopicRepository
	)

	users = inmem.NewUserRepository()
	topics = inmem.NewTopicRepository()

	fieldKeys := []string{"method"}

	var us user.Service
	us = user.NewService(users)
	us = user.NewLoggingService(log.With(logger, "component", "user"), us)
	us = user.NewTracingService(us)
	us = user.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "user_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "user_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		us,
	)

	var tp topic.Service
	tp = topic.NewService(topics)
	tp = topic.NewLoggingService(log.With(logger, "component", "topic"), tp)
	tp = topic.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "topic_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "topic_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		tp,
	)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/user/", user.MakeHandler(us, httpLogger))
	mux.Handle("/topics/", topic.MakeHandler(tp, httpLogger))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
