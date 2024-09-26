// the server binary is responsible for receiving and handing request
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/cmd/server/app"

	"github.com/golang/glog"
)

func main() {
	time.Local = time.FixedZone("CST", 8*3600)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	defer glog.Flush()

	command := app.NewServerCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
