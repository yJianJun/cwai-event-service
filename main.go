// the server binary is responsible for receiving and handing request
package main

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/app"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/model"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/golang/glog"
)

//	@title			CTCCL事件监听
//	@version		1.0
//	@description	监听CTCCL上报事件服务

//	@contact.name	yejianjun
//	@contact.email	yejianjun@ideal.sh.cn

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath		/ctccl
func main() {
	time.Local = time.FixedZone("CST", 8*3600)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	defer glog.Flush()
	var err error
	model.DB, err = model.ConnectDatabase()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	command := app.NewServerCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
