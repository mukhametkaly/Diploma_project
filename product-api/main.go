package main

import (
	"flag"
	"fmt"
	"github.com/djumanoff/amqp"
	"github.com/mukhametkaly/Diploma/product-api/src/catalog"
	"github.com/mukhametkaly/Diploma/product-api/src/config"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/sirupsen/logrus"
)

var service catalog.Service

func main() {
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address only port :8080")
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = level.NewFilter(logger, level.AllowAll())
	logger = &serializedLogger{Logger: logger}
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	logr := &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.Level(5),
		Formatter: &logrus.TextFormatter{
			FullTimestamp: true,
		},
	}
	catalog.Loger = logr

	err := config.GetConfigs()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	service = catalog.NewService()
	service = catalog.NewLoggingService(log.With(logger, "component", "catalog"), service)
	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/v1/product/", catalog.MakeHandler(service, httpLogger))
	http.Handle("/v1/product/", accessControl(mux))
	http.HandleFunc("/v1/check", config.Healthchecks)
	errs := make(chan error, 2)

	var consumerCfg amqp.ConsumerConfig
	consumerCfg.PrefetchCount = 100

	sessionrabbit, err := catalog.GetRabbitSession()
	if err != nil {
		level.Error(logger).Log("Error creating session to rabbit:", fmt.Sprintf("%v", err))
		return
	}

	serverRabbit, err := (*sessionrabbit).Consumer(consumerCfg)
	if err != nil {
		level.Error(logger).Log("Error creating consumer to rabbit:", fmt.Sprintf("%v", err))
		return
	}

	if err := serverRabbit.Queue(amqp.Queue{
		Name:    "event.catalog",
		Durable: true,
		AutoAck: false,
		NoWait:  true,
		Args:    nil,
		Bindings: []amqp.QueueBinding{
			{
				RoutingKey: "event.catalog.update.nomenclatures",
				Exchange:   "X:routing.topic",
				NoWait:     true,
				Args:       nil,
				Handler: func(message amqp.Message) *amqp.Message {
					return catalog.UpdateOrderAndDeliveryStatusEvent(message, service)
				},
			},
		},
	}); err != nil {
		panic("err in queue (name = FileUploaderEvent)")
	}

	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening catalog-api V1.0.0")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization,X-Owner,darvis-dialog-id")

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

type serializedLogger struct {
	mtx sync.Mutex
	log.Logger
}

func (l *serializedLogger) Log(keyvals ...interface{}) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.Logger.Log(keyvals...)
}

func healthchecks(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
