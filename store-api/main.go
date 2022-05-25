package main

import (
	"flag"
	"fmt"
	"github.com/djumanoff/amqp"
	"github.com/mukhametkaly/Diploma_project/auth-api/src/auth"
	"github.com/mukhametkaly/Diploma_project/auth-api/src/config"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/sirupsen/logrus"
)

var authService auth.Service

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
	auth.Loger = logr

	err := config.GetConfigs()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	authService = auth.NewService()
	authService = auth.NewLoggingService(log.With(logger, "component", "merchant"), authService)
	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/v1/auth/", auth.MakeHandler(authService, httpLogger))
	http.Handle("/v1/auth/", accessControl(mux))
	http.HandleFunc("/v1/auth/check", config.Healthchecks)
	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening merchant-api V1.0.0")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		srv, err := auth.Server()
		if err != nil {
			panic(fmt.Errorf("Fatal error connect Rabbit: %s \n", err))
		}
		if err := (*srv).Endpoint("request.auth.#", func(message amqp.Message) *amqp.Message {
			return auth.MakeAuthServiceRabbitMQ(authService, message)
		}); err != nil {
			fmt.Println("err = ", err)
		}
		logr.Debug("auth in auth-api auth RabbitMQ server started")
		errs <- (*srv).Start()
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
