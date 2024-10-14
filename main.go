//	@title			CTCCL事件监听
// @version 1.0
//	@description	监听CTCCL上报事件服务
// @termsOfService http://example.com/terms/

//	@contact.name	yejianjun
// @contact.url http://www.example.com/support
//	@contact.email	yejianjun@ideal.sh.cn

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath		/ctccl

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

package main

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/app"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/golang/glog"
)

func main() {
	time.Local = time.FixedZone("CST", 8*3600)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	defer glog.Flush()
	var err error
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	command := app.NewServerCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
