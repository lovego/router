package goa

import (
	"log"
	"reflect"

	"github.com/lovego/errs"
)

func convertHandler(h interface{}, path string) func(*Context) {
	if handler, ok := h.(func(*Context)); ok {
		return handler
	}

	val := reflect.ValueOf(h)

	typ := val.Type()
	if typ.Kind() != reflect.Func {
		log.Panic("handler must be a func.")
	}
	if typ.NumIn() != 2 {
		log.Panic("handler func must have exactly two parameters.")
	}
	if typ.NumOut() != 0 {
		log.Panic("handler func must have no return values.")
	}

	reqConvertFunc, hasCtx := newReqConvertFunc(typ.In(0), path)
	respTyp, respWriteFunc := newRespWriteFunc(typ.In(1), hasCtx)

	return func(ctx *Context) {
		req, err := reqConvertFunc(ctx)
		if err != nil {
			ctx.Data(nil, errs.New("args-err", err.Error()))
			return
		}
		resp := reflect.New(respTyp)
		val.Call([]reflect.Value{req, resp})
		if respWriteFunc != nil {
			respWriteFunc(ctx, resp.Elem())
		}
	}
}

// handler example
func handlerExample(req *struct {
	Title string
	Desc  string
	Param struct {
		Id int64
	}
	Query struct {
		Id   int64
		Page int64
	}
	Header struct {
		Cookie string
	}
	Body struct {
		Id   int64
		Name string
	}
	Session struct {
		UserId int64
	}
	Ctx *Context
}, resp *struct {
	Error error
	Data  struct {
		Id   int64
		Name string
	}
	Header struct {
		SetCookie string
	}
}) {
	// resp.Body, resp.Error = users.Get(req.Params.Id)
}
