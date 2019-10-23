package handler

import (
	"encoding/json"
	"random-cache/internal/worker"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type FastHTTPHandler struct {
	worker *worker.Worker
	logger *logrus.Logger
}

func New(worker *worker.Worker, logger *logrus.Logger) *FastHTTPHandler {
	return &FastHTTPHandler{
		worker: worker,
		logger: logger,
	}
}

type OkResponseWithResult struct {
	Result worker.LastTwoElemsFromCache `json:"result"`
}

func (f *FastHTTPHandler) sendError(text string, code int, ctx *fasthttp.RequestCtx) {
	ctx.Error(text, code)
}

func (f *FastHTTPHandler) sendOkResponseWithResult(ctx *fasthttp.RequestCtx, targetStruct interface{}) error {
	ctx.SetContentType("application/json")
	enc := json.NewEncoder(ctx.Response.BodyWriter())
	return enc.Encode(targetStruct)
}

func (t *FastHTTPHandler) LastTwoElemsFromCache(ctx *fasthttp.RequestCtx) {
	lastTwoElemsFromCache, err := t.worker.LastTwoElemsFromCache()

	if err != nil {
		t.sendError(err.Error(), 500, ctx)
		return
	}

	err = t.sendOkResponseWithResult(ctx, lastTwoElemsFromCache)

	// удалить если ожидается высокая нагрузка, так как ошибка может быть вызвана и закрытием сокета клиента
	if err != nil {
		t.logger.Errorf("Error when try t.sendOkResponseWithResult[(t *FastHTTPHandler) "+
			"LastTwoElemsFromCache], err: %s", err)
		return
	}
}
