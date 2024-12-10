package glog_clean

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	program  = filepath.Base(os.Args[0])
	host     = "unknownhost"
	userName = "unknownuser"
)

func init() {
	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		userName = current.Username
	}

	// Sanitize userName since it may contain filepath separators on Windows.
	userName = strings.Replace(userName, `\`, "_", -1)
}

// shortHostname returns its argument, truncating at the first period.
// For instance, given "www.google.com" it returns "www".
func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

// logPrefix: 获取glog日志文件前缀
func logPrefix() string {
	return fmt.Sprintf("%s.%s.%s.log", program, host, userName)
}

// getTagName: 获取日志文件名称中的tag，即日志等级，例如INFO,ERROR
func getTagName(logName string) string {
	// 参考：https://github.com/golang/glog/blob/master/glog_file.go#L83
	// ${program}.${host}.${userName}.log.${tag}.${time}.${pid}
	// 示例：vc-scheduler.volcano-scheduler-5b7b6bcd7c-slqnp.root.log.WARNING.20210303-183058.1
	parts := strings.Split(logName, ".")
	if len(parts) >= 5 {
		return parts[4]
	}
	return ""
}
