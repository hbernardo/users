package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/hbernardo/users/go-src/infra"
	"github.com/hbernardo/users/go-src/lib"
	"github.com/hbernardo/users/go-src/srv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	usersDataFilePath = "data/users.json"
)

type serviceConfig struct {
	ServerPort int `env:"PORT,required"`

	HealthCheckPort    int    `env:"HEALTH_CHECK_PORT,required"`
	LivenessProbePath  string `env:"LIVENESS_PROBE_PATH,required"`
	ReadinessProbePath string `env:"READINESS_PROBE_PATH,required"`

	RateLimitMaxFrequency   int           `env:"RATE_LIMIT_MAX_FREQUENCY,required"`
	RateLimitBurstSize      int           `env:"RATE_LIMIT_BURST_SIZE,required"`
	RateLimitMemoryDuration time.Duration `env:"RATE_LIMIT_MEMORY_DURATION,required"`

	CORSAllowOrigin  string   `env:"CORS_ALLOW_ORIGIN,required"`
	CORSAllowMethods []string `env:"CORS_ALLOW_METHODS,required"`
	CORSAllowHeaders []string `env:"CORS_ALLOW_HEADERS,required"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"error"`
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

	// Reading users data file
	usersData, err := readUsersDataJSONFile()
	if err != nil {
		return err
	}

	// Default HTTP Server
	httpSrv := srv.NewHTTPServer(config.ServerPort,
		srv.NewUsersHandler(
			lib.NewUsersService(
				infra.NewUsersRepo(usersData),
			),
		),
		srv.CORSMiddleware(
			config.CORSAllowOrigin,
			config.CORSAllowMethods,
			config.CORSAllowHeaders,
		),
		srv.RateLimiterMiddleware(
			config.RateLimitMaxFrequency,
			config.RateLimitBurstSize,
			config.RateLimitMemoryDuration,
		),
		srv.PanicRecoveryMiddleware,
	)
	defer httpSrv.Close(ctx)
	httpSrv.ListenAndServe()

	// Health Server (for liveness and readiness probes)
	healthSrv := srv.NewHTTPServer(config.HealthCheckPort,
		srv.NewHealthHandler(config.LivenessProbePath, config.ReadinessProbePath),
	)
	defer healthSrv.Close(ctx)
	healthSrv.ListenAndServe()

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

func readUsersDataJSONFile() ([]lib.User, error) {
	jsonFile, err := os.Open(usersDataFilePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var usersData []lib.User

	err = json.Unmarshal(jsonBytes, &usersData)
	if err != nil {
		return nil, err
	}

	return usersData, nil
}

func waitSignal() os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	return <-sig
}
