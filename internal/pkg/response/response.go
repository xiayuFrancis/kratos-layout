package response

import (
	commonv1 "kratosdemo/api/common/v1"

	"github.com/bytedance/sonic"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// 响应状态码定义
const (
	// 成功
	CodeSuccess = 0
	// 通用错误
	CodeError = 1
	// 参数错误
	CodeInvalidParam = 400
	// 未授权
	CodeUnauthorized = 401
	// 禁止访问
	CodeForbidden = 403
	// 资源不存在
	CodeNotFound = 404
	// 服务器内部错误
	CodeServerError = 500
)

// Success 返回成功响应
func Success(data proto.Message) *commonv1.Response {
	// 创建空的 Struct
	st, _ := structpb.NewStruct(map[string]interface{}{})
	resp := &commonv1.Response{
		Code: CodeSuccess,
		Msg:  "success",
		Ext:  st,
	}

	if data != nil {
		// 使用 protojson 将 proto 消息转换为 JSON
		m := protojson.MarshalOptions{
			EmitUnpopulated: true,
			UseProtoNames:   true,
		}
		jsonBytes, err := m.Marshal(data)
		if err == nil {
			// 解析 JSON 字节到结构体
			var jsonMap map[string]interface{}
			if err := sonic.Unmarshal(jsonBytes, &jsonMap); err == nil {
				st, err := structpb.NewStruct(jsonMap)
				if err == nil {
					resp.Ext = st
				}
			}
		}
	}

	return resp
}

// Error 返回错误响应
func Error(code int32, msg string) *commonv1.Response {
	// 创建空的 Struct
	st, _ := structpb.NewStruct(map[string]interface{}{})
	return &commonv1.Response{
		Code: code,
		Msg:  msg,
		Ext:  st, // 空 JSON 对象
	}
}

// ErrorWithData 返回带数据的错误响应
func ErrorWithData(code int32, msg string, data proto.Message) *commonv1.Response {
	// 创建空的 Struct
	st, _ := structpb.NewStruct(map[string]interface{}{})
	resp := &commonv1.Response{
		Code: code,
		Msg:  msg,
		Ext:  st, // 默认为空 JSON 对象
	}

	if data != nil {
		// 使用 protojson 将 proto 消息转换为 JSON
		m := protojson.MarshalOptions{
			EmitUnpopulated: true,
			UseProtoNames:   true,
		}
		jsonBytes, err := m.Marshal(data)
		if err == nil {
			// 解析 JSON 字节到结构体
			var jsonMap map[string]interface{}
			if err := sonic.Unmarshal(jsonBytes, &jsonMap); err == nil {
				st, err := structpb.NewStruct(jsonMap)
				if err == nil {
					resp.Ext = st
				}
			}
		}
	}

	return resp
}
