package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	c "random-cache/internal/config"
	h "random-cache/internal/handler"
	w "random-cache/internal/worker"

	"github.com/buaazp/fasthttprouter"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func initRouter(router *fasthttprouter.Router, fastHTTPHandler *h.FastHTTPHandler, cfg *c.Config) {
	router.GET(fmt.Sprintf("/%s", cfg.EndPointStr), fastHTTPHandler.LastTwoElemsFromCache)
}

func wrapHandler(r *fasthttprouter.Router) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		r.Handler(ctx)
	}
}

func gracefullShutdown(stopProgramChan chan os.Signal, log *logrus.Logger,
	fastHTTPServer *fasthttp.Server, worker *w.Worker) {
	s := <-stopProgramChan

	log.Infof("Caught signal %v: terminating", s)

	err := fastHTTPServer.Shutdown()
	if err != nil {
		log.Errorf("Error when try fastHTTPServer.Shutdown[gracefullShutdown], err: %s", err)
	}

	worker.Close()

	stopProgramChan <- s
}

func main() {
	stopProgramChan := make(chan os.Signal, 1)
	signal.Notify(stopProgramChan, syscall.SIGINT, syscall.SIGTERM)

	// need validate cfg
	cfg := &c.Config{
		BindAddress:                         "localhost",
		BindPort:                            8080,
		LenStringForAddToCache:              6,
		FrequencyAddToCacheSec:              2,
		HowMuchLastElemsFromCacheNeedReturn: 2,
		EndPointStr:                         "last_two_elems_from_cache",
	}

	// default level: info
	logger := logrus.New()

	worker := w.New(cfg, logger)

	// start timer for add in cache
	// начнёт добавлять уже до непосредственного старта сервера, по идее надо подумать как
	// отловить событие старта сервера и уже потом запускать таймер
	go worker.Init()

	handler := h.New(worker, logger)

	router := fasthttprouter.New()

	initRouter(router, handler, cfg)

	wrappedFastHTTPHandler := wrapHandler(router)

	fastHTTPServer := &fasthttp.Server{
		Handler:          wrappedFastHTTPHandler,
		DisableKeepalive: true,
	}

	go gracefullShutdown(stopProgramChan, logger, fastHTTPServer, worker)

	logger.Info("Start web service")

	err := fastHTTPServer.ListenAndServe(fmt.Sprintf("%s:%d", cfg.BindAddress, cfg.BindPort))
	if err != nil {
		logger.Fatalf("error when try fastHTTPServer.ListenAndServe[main], %s", err)
	}

	<-stopProgramChan
}
