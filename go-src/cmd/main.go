package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	"github.com/hbernardo/users/go-src/srv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type serviceConfig struct {
	HealthCheckPort int    `env:"HEALTH_CHECK_PORT,required"`
	ServerPort      int    `env:"PORT,required"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"error"`
}

func getServiceConfig() (*serviceConfig, error) {
	config := &serviceConfig{}
	err := env.Parse(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// CLI commands
var (
	rootCmd = &cobra.Command{
		Use:   "insights-platform-api",
		Short: "Insights platform backend service",
	}
	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Start the server",
		RunE:  runHTTP,
	}
)

func init() {
	rootCmd.AddCommand(httpCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func runHTTP(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	config, err := getServiceConfig()
	if err != nil {
		return err
	}

	err = configureLog(config.LogLevel)
	if err != nil {
		return err
	}

	// TODO: test, should be refactored
	go func() {
		http.HandleFunc("/health/live", func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(w, "live\n")
		})
		http.HandleFunc("/health/ready", func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(w, "ready\n")
		})

		http.ListenAndServe(fmt.Sprintf(":%d", config.HealthCheckPort), nil)
	}()

	httpSrv := srv.NewHTTPServer(config.ServerPort)
	defer httpSrv.Close(ctx)

	httpSrv.ListenAndServe()

	sig := waitSignal() // blocking until signal
	log.WithFields(log.Fields{
		"signal": sig.String(),
	}).Debug("received signal, exiting...")

	return nil
}

func configureLog(logLevel string) error {
	lv, err := log.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	log.SetLevel(lv)

	log.SetFormatter(&log.JSONFormatter{})

	return nil
}

func waitSignal() os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	return <-sig
}
