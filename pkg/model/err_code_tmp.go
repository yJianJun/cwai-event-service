package model

type ErrorCode string

/*
ErrorCode包含两部分：错误码和错误信息，以英文:分隔。错误码必须为三段式：服务.模块.错误。
错误信息尽可能满足：
1.错误信息是给用户看的，必须简洁扼要可读性强，尤其不要出现开发者才熟悉的短语或缩写等，要转化为用户熟悉的表达，例如，好的示例：资源组ID不存在；不好的示例：resourceGroupID不存在
2.错误信息要具体，不要一个笼统的ErrorCode到处用，理想情况是一个具体错误对应一个ErrorCode
3.错误信息风格要统一：例如请求参数错误，不要有的叫参数不合法，有的叫参数非法，有的叫请求字段错误
*/
const (
	EventUnAuthorized  ErrorCode = "Cwai.Event.UnAuthorized:没有访问权限"
	EventInvalidParam  ErrorCode = "Cwai.Event.InvalidParam:请求字段错误"
	EventForbidden     ErrorCode = "Cwai.Event.Forbidden:不可访问"
	EventInternalError ErrorCode = "Cwai.Event.InternalError:服务异常"
)
