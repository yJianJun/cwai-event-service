package app

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/config"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"os"
	// "os/signal"
	// "syscall"
)

var (
	sigHandles []func() error
	configPath string
)

func NewServerCommand() *cobra.Command {
	glog.Info("run NewServerCommand.")
	opts := config.NewConfig()
	cmd := &cobra.Command{
		Use:           "daemon [OPTIONS]",
		Short:         "The daemon  server",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
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

func runDaemon(opts *config.ServerConfig) error {
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

	select {

	case sig := <-signalCh:
		klog.Infof("received signal: %s", sig)
		os.Exit(1)
	case err := <-errCh:
		return err
	}
	return nil
}
