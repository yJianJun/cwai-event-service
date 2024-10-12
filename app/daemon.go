package app

import (
	"context"
	"net/http"
	"time"

	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/client"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/model"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/router"
	"github.com/golang/glog"
)

type Daemon struct {
	Config *model.ServerConfig
	Server *http.Server
}

func NewDaemon(cfg *model.ServerConfig) *Daemon {
	address := cfg.Host + ":" + cfg.Port
	routers := router.InitRoute()

	//init ccae client config
	client.NewClient(cfg)
	model.ConnectDatabase(cfg)
	model.InitES(cfg)

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

// shutdown gracefully
func (da *Daemon) Shutdown() error {
	glog.Info("Shutdown of Daemon")

	shutdownTimeout := da.Config.ShutTimeOut

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
	defer cancel()
	if err := da.Server.Shutdown(ctx); err != nil {
		glog.Fatal("Server Shutdown:", err)
		return err
	}

	// catching ctx.Done()
	select {
	case <-ctx.Done():
		glog.Infof("shutdown timeout")
	}

	glog.Infof("shutdown gracefully")
	return nil
}
