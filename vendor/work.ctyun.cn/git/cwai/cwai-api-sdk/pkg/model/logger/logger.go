package logger

type GetLogParam struct {
	ResourceGroupID string `form:"resourceGroupID"` // 资源组名, 必填
	Namespace       string `form:"namespace"`       // 命名空间, 必填
	PodName         string `form:"podName"`         // pod名, 必填
	Container       string `form:"container"`       // 容器名, 必填
	TailLines       int64  `form:"tailLines"`       // 要显示的最新的日志条数, -1会显示所有的日志
	Timestamps      bool   `form:"timestamps"`      // 为true在日志中包含时间戳
	Previous        bool   `form:"previous"`        //如果为true, 输出pod中曾经运行过,但目前已终止的容器的日志
}
