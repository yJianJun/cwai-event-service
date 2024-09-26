package app

import (
	"fmt"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"os"

	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/model"
	// "os/signal"
	// "syscall"
)

var (
	sigHandles []func() error
	configPath string
)

func NewServerCommand() *cobra.Command {
	glog.Info("run NewServerCommand.")
	opts := model.NewConfig()
	cmd := &cobra.Command{
		Use:           "daemon [OPTIONS]",
		Short:         "The daemon  server",
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

	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "../../conf/config.yaml", "config path")
	opts.BindFlags(cmd)

	return cmd
}

func runDaemon(opts *model.ServerConfig) error {
	glog.Info("run runDaemon.")

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
		klog.Infof("received signal: %s", sig)
		for _, handle := range sigHandles {
			if err := handle(); err != nil {
				klog.Errorf("failed to handle signal: %v", err)
			}
		}
		os.Exit(1)
	case err := <-errCh:
		return err

	}
	return nil
}
