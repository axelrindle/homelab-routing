package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "embed"

	"github.com/axelrindle/traefik-configuration-provider/app"
	"github.com/axelrindle/traefik-configuration-provider/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:embed banner.txt
var banner string

var (
	Version        = "dev"
	CommitHash     = "unknown"
	BuildTimestamp = "unknown"
)

var (
	configFile  string
	showVersion bool
	healthcheck bool
)

func BuildVersion() string {
	return fmt.Sprintf("%s-%s (%s)", Version, CommitHash, BuildTimestamp)
}

func makeLogger(c *config.Config) (*zap.Logger, error) {
	if c.Environment == "production" {
		return zap.NewProduction()
	} else {
		conf := zap.NewDevelopmentConfig()
		conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return conf.Build()
	}
}

func main() {
	println(banner)

	flag.StringVar(&configFile, "config", "config.yml", "path to the config file")
	flag.BoolVar(&showVersion, "version", false, "show program version")
	flag.BoolVar(&healthcheck, "healthcheck", false, "run a healthcheck")
	flag.Parse()
	if showVersion {
		println(BuildVersion())
		return
	}

	config := &config.Config{}
	config.Load(configFile)

	logger, err := makeLogger(config)
	if err != nil {
		log.Fatal(err)
	}

	if healthcheck {
		addressParts := strings.Split(config.Server.Address, ":")
		if len(addressParts) != 2 {
			logger.Fatal("invalid server address", zap.String("address", config.Server.Address))
		}
		port := addressParts[1]

		res, err := http.Get("http://127.0.0.1:" + port + "/status")
		if err != nil {
			log.Fatal("healthcheck failed", zap.Error(err))
		}

		if res.StatusCode != 200 {
			log.Fatal("healthcheck failed", zap.Int("status", res.StatusCode))
		}

		return
	}

	logger.Info("starting program", zap.String("mode", config.Environment))

	defer logger.Sync()

	app := &app.App{
		Config: config,
		Logger: logger.With(zap.String("component", "app")),
	}
	app.Init()

	go app.Boot()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down …")

	app.Shutdown()
}
