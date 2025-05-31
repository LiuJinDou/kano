package dto

type Request struct {
}

// Response 通用返回结构
// swagger:model
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // 使用interface{}以支持不同类型的数据
}
