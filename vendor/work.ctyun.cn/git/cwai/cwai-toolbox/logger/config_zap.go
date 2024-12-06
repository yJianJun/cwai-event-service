package logger

type ZapConfig struct {
	// Name Root Logger 的名称
	Name string `json:"name" yaml:"name" mapstructure:"name"`
	// Level 进行显示和收集的最低日志等级
	// 当前可选 debug、info、warn、error、panic、fatal
	// 省略或填写为空，则默认为 debug 等级。
	Level LogLevel `json:"level" yaml:"level" mapstructure:"level"`
	// TraceLevel 触发 Stack Trace 输出的最低日志等级
	// 当前可选 debug、info、warn、error、panic、fatal
	// 省略或填写为空，则默认为 panic 等级。
	TraceLevel LogLevel `json:"traceLevel" yaml:"traceLevel" mapstructure:"traceLevel"`
	// LogToStdOut 是否输出到 stdout
	LogToStdOut bool `json:"logToStdOut" yaml:"logToStdOut" mapstructure:"logToStdOut"`
}
