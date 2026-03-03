package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

/*
	server.response -- [submodule]
One invocation format I envision is: server.[responder].Response(event). This is the
simplest extensible form I can think of. Obviously, we only need to implement
Response(event) for different backend frameworks, and then switch the responder
accordingly. At the same time, to ensure that different backends can easily implement it
in their own way, the responder does not need to be standardized; it is just a structure
that implements Response(event). This allows for constraints and ensures compatibility
when appropriate.
For simple HTTP events, only three event states are needed: success, failure, and error.
In fact, most responses can be defined in these three states.
A typical modern HTTP response format looks like this:
	content-type: application/json
	{
		code: <logic code>,
		msg: <logic message>
		data: <response data if success>
	}
JSON is the most common format used in RESTful applications, but this doesn't mean that only this format is used.
The response needs to dynamically change and transmit the actual value based on the given Content-Type.
*/

const (
	// Success(httpCode,msg,data)
	HTTPEVENT_SUCCESS uint8 = 0x00
	// Failure(httpCode,code, msg)
	HTTPEVENT_FAILURE uint8 = 0x01
	// Error(httpCode,msg)
	HTTPEVENT_ERROR uint8 = 0xf0

	// status text default
	DEFAULT_STATUS_TEXT = ""
)

type (
	// replace nil -> '{}'
	EmptyData struct{}
	// replace nil -> '[{}]'
	EmptyList []EmptyData
	Responder interface {
		ResponseHttp(re HttpEvent) error
	}
	EchoResponder struct {
		c echo.Context
	}
	HttpEvent struct {
		Type     uint8
		HttpCode int
		Code     int
		Msg      string
		Data     any
	}
)

func NewEchoResponder(c echo.Context) EchoResponder {
	return EchoResponder{
		c: c,
	}
}

func (rd *EchoResponder) Success(httpCode int, msg string, data any) error {
	return rd.c.JSON(httpCode, map[string]any{
		"code": 0,
		"msg":  msg,
		"data": data,
	})
}

func (rd *EchoResponder) Failure(httpCode, code int, msg string) error {
	return rd.c.JSON(httpCode, map[string]any{
		"code": code,
		"msg":  msg,
		"data": nil,
	})
}

func (rd *EchoResponder) Error(httpCode int, msg string) error {
	if msg == DEFAULT_STATUS_TEXT {
		msg = http.StatusText(httpCode)
	}
	return rd.c.JSON(httpCode, map[string]any{
		"code": -1,
		"msg":  msg,
		"data": nil,
	})
}

func (rd *EchoResponder) ResponseHttp(re HttpEvent) error {
	switch re.Type {
	case HTTPEVENT_SUCCESS:
		return rd.Success(re.HttpCode, re.Msg, re.Data)
	case HTTPEVENT_FAILURE:
		return rd.Failure(re.HttpCode, re.Code, re.Msg)
	case HTTPEVENT_ERROR:
		return rd.Error(re.Code, re.Msg)
	}
	return rd.c.JSON(re.HttpCode, re.Data)
}
