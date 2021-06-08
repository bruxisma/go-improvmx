package internal

import (
	"github.com/go-resty/resty/v2"
)

type Request resty.Request

func (request *Request) ToResty() *resty.Request {
	return (*resty.Request)(request)
}

func (r *Request) SetPathParameter(parameter, value string) *Request {
	return (*Request)(r.ToResty().SetPathParam(parameter, value))
}

func (r *Request) SetQueryParameter(parameter, value string) *Request {
	return (*Request)(r.ToResty().SetQueryParam(parameter, value))
}

func (r *Request) SetBody(body interface{}) *Request {
	return (*Request)(r.ToResty().SetBody(body))
}

func (request *Request) Get(url string) (interface{}, error) {
	response, error := request.ToResty().Get(url)
	return checkResponse(response, error)
}

func (request *Request) Post(url string) (interface{}, error) {
	response, error := request.ToResty().Post(url)
	return checkResponse(response, error)
}

func (request *Request) Put(url string) (interface{}, error) {
	response, error := request.ToResty().Put(url)
	return checkResponse(response, error)
}

func (request *Request) Delete(url string) (interface{}, error) {
	response, error := request.ToResty().Delete(url)
	return checkResponse(response, error)
}

func checkResponse(response *resty.Response, e error) (interface{}, error) {
	if e != nil {
		return nil, e
	}
	if response.IsError() {
		return nil, response.Error().(error)
	}
	return response.Result(), nil
}
