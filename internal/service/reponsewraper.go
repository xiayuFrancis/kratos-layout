package service

import (
	"blog/api/common"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
)

// ResponseWrapper 用于包装响应
type ResponseWrapper struct{}

// Success 成功响应
func (r *ResponseWrapper) Success(data interface{}) (*common.Response, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &common.Response{
		Code: 0,
		Msg:  "success",
		Data: dataBytes,
	}, nil
}

// Error 错误响应
func (r *ResponseWrapper) Error(err error) *common.Response {
	if e, ok := err.(*errors.Error); ok {
		return &common.Response{
			Code: int32(e.Code),
			Msg:  e.Reason,
		}
	}
	// 默认错误码 500
	return &common.Response{
		Code: 500,
		Msg:  err.Error(),
	}
}

func NewResponseWrapper() *ResponseWrapper {
	return &ResponseWrapper{}
}
