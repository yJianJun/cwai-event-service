package app

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/config"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/router"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/util"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/validatorx"
	"github.com/golang/glog"
	"net/http"
)

type Daemon struct {
	Config *config.ServerConfig
	Server *http.Server
}

func NewDaemon(cfg *config.ServerConfig) *Daemon {
	address := cfg.App.Host + ":" + cfg.App.Port
	routers := router.InitRoute()
	util.InitElasticSearch(cfg)
	// 参数校验器初始化、如错误提示中文转译、注册自定义校验器等
	validatorx.Init()
	srv := &http.Server{
		Addr:    address,
		Handler: routers,
	}
	return &Daemon{
		Config: cfg,
		Server: srv,
	}
}

// run the main operator
func (da *Daemon) Run() error {
	glog.Info("Run of Daemon")

	startError := make(chan error)
	go func(errCh chan error) {
		if err := da.Server.ListenAndServe(); err != nil {
			glog.Errorf("failed to start http server: %v", err)
			startError <- err
		}
	}(startError)
	err := <-startError
	if err != nil {
		return err
	}
	return nil
}
