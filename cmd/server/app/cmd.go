package app

import (
	"context"
	goflag "flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"os"

	"work.ctyun.cn/git/cwai/cwai-event-service/pkg/config"
	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

var (
	sigHandles []func() error
	configPath string
)

func NewServerCommand() *cobra.Command {
	logger.Info(context.TODO(), "run NewEventServiceCommand.")
	opts := config.NewConfig()
	cmd := &cobra.Command{
		Use:           "daemon [OPTIONS]",
		Short:         "The daemon server",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			if err := opts.ReadYAML(configPath); err != nil {
				fmt.Fprintf(os.Stderr, "[WARN] read config file failed: %s\n", err)
			}
			if err := opts.Parse(); err != nil {
				fmt.Fprintf(os.Stderr, "[WARN] parse config failed: %s\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "config: %+v\n", opts)
			if err := runDaemon(opts); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(2)
			}
		},
	}

	// 将config加入命令行flag，配置通过flag库读取的flag
	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "config path")
	if err := opts.BindFlags(cmd); err != nil {
		logger.Errorf(context.TODO(), "failed to bind flags, err: %s", err.Error())
		return cmd
	}

	// parse golog flags
	cmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)
	if err := goflag.CommandLine.Parse([]string{}); err != nil {
		logger.Errorf(context.TODO(), "failed to parse params, err: %s", err.Error())
		return cmd
	}

	return cmd
}

func runDaemon(opts *config.ServerConfig) error {
	logger.Info(context.TODO(), "run runDaemon.")

	var (
		errCh    = make(chan error, 1)
		signalCh = make(chan os.Signal, 1)
	)
	daemon := NewDaemon(opts)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		errCh <- daemon.Run()
	}()
	//shutdown gracefully
	sigHandles = append(sigHandles, daemon.Shutdown)
	select {

	case sig := <-signalCh:
		logger.Infof(context.TODO(), "received signal: %s", sig)
		for _, handle := range sigHandles {
			if err := handle(); err != nil {
				logger.Errorf(context.TODO(), "failed to handle signal: %v", err)
			}
		}
		os.Exit(1)
	case err := <-errCh:
		return err

	}
	return nil
}
