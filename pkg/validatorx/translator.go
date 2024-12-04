package validatorx

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 绑定请求体中的json至form结构体，并拷贝至另一结构体
func ShouldBindJSON[T any](g *gin.Context, data T) {
	if err := g.ShouldBindJSON(data); err != nil {
		// 统一recover处理
		panic(ConvBindValidationError(data, err))
	}
}

// 转译参数校验错误，并将参数校验错误为业务异常错误（统一recover处理）
func ConvBindValidationError(data any, err error) error {
	if e, ok := err.(validator.ValidationErrors); ok {
		// 调用validatorx.Translate2Str方法进行校验错误转译
		panic(common.CommonError{Code: 403, Msg: Translate2Str(data, e)})
	}
	return err
}
