package testutils

import (
	"encoding/json"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
)

type EchoBuilder interface {
	Request(httpMethod string, path string, body any) RequestBuilder
}

type RequestBuilder interface {
	AddParams(params map[string]string) RequestBuilder
	GetContextAndResponseRecorder() (echo.Context, *httptest.ResponseRecorder)
}

type builder struct {
	e   *echo.Echo
	ctx echo.Context
	rec *httptest.ResponseRecorder
}

func NewEchoBuilder() EchoBuilder {
	return &builder{
		e: echo.New(),
	}
}

func (b *builder) Request(httpMethod string, path string, body any) RequestBuilder {
	bytes, err := json.Marshal(body)
	if err != nil {
		panic("failed to marshal body: " + err.Error())
	}

	b.rec = httptest.NewRecorder()
	req := httptest.NewRequest(httpMethod, path, strings.NewReader(string(bytes)))
	b.ctx = b.e.NewContext(req, b.rec)

	return b
}

func (b *builder) AddParams(params map[string]string) RequestBuilder {
	if b.ctx == nil {
		panic("context is nil, please call AddRequest before adding params")
	}

	for k, v := range params {
		b.ctx.SetParamNames(k)
		b.ctx.SetParamValues(v)
	}

	return b
}

func (b *builder) GetContextAndResponseRecorder() (echo.Context, *httptest.ResponseRecorder) {
	if b.ctx == nil || b.rec == nil {
		panic("context and/or response recorder are nil, please call AddRequest before getting context")
	}

	return b.ctx, b.rec
}
