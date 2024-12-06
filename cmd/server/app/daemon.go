package app

import (
	"context"
	"net/http"
	"time"

	"work.ctyun.cn/git/cwai/cwai-event-service/pkg/config"
	"work.ctyun.cn/git/cwai/cwai-event-service/pkg/router"
	"work.ctyun.cn/git/cwai/cwai-event-service/pkg/utils"
	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

type Daemon struct {
	Config *config.ServerConfig
	Server *http.Server
}

func NewDaemon(cfg *config.ServerConfig) *Daemon {
	address := cfg.App.Host + ":" + cfg.App.Port
	routers := router.InitRoute()

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
	logger.Infof(context.TODO(), "Run of Daemon, daemon:%+v %+v", da.Config, da.Server)

	// logger init
	logger.Init(&logger.Config{
		ZapConfig: logger.ZapConfig{
			Name:        da.Config.LoggerInfo.Name,
			Level:       logger.LogLevel(da.Config.LoggerInfo.Level),
			TraceLevel:  logger.LogLevel(da.Config.LoggerInfo.TraceLevel),
			LogToStdOut: true,
		},
		LumberjackConfig: logger.LumberjackConfig{
			LogToDir:     da.Config.LoggerInfo.LogToDir,
			MaxSizeInMiB: da.Config.LoggerInfo.MaxSizeInMiB,
			MaxAgeInDays: da.Config.LoggerInfo.MaxAgeInDays,
		},
		FieldKeys: []string{"reason", "ID", "who", "what", "error", "config"},
	})

	logger.Debug(context.TODO(), "Logger test: Debug.")
	logger.Info(context.TODO(), "Logger test: Info.")
	logger.Error(context.TODO(), "Logger test: error.")

	err := utils.InitElasticSearch()
	if err != nil {
		logger.Errorf(context.TODO(), "Init elasticsearch failed: %v\n", err)
		return err
	}

	startError := make(chan error)
	go func(errCh chan error) {
		if err := da.Server.ListenAndServe(); err != nil {
			logger.Errorf(context.TODO(), "failed to start http server: %v", err)
			startError <- err
		}
	}(startError)
	err = <-startError
	if err != nil {
		return err
	}
	return nil
}

// shutdown gracefully
func (da *Daemon) Shutdown() error {
	logger.Info(context.TODO(), "Shutdown of Daemon")

	shutdownTimeout := da.Config.App.ShutTimeOut

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
	defer cancel()
	if err := da.Server.Shutdown(ctx); err != nil {
		logger.Fatal(context.TODO(), "Server Shutdown: "+err.Error())
		return err
	}

	// catching ctx.Done()
	select {
	case <-ctx.Done():
		logger.Infof(context.TODO(), "shutdown timeout")
	}

	logger.Infof(context.TODO(), "shutdown gracefully")
	return nil
}
