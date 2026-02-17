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

	"github.com/axelrindle/traefik-configuration-provider/internal/app"
	"github.com/axelrindle/traefik-configuration-provider/internal/config"
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
		conf := zap.NewProductionConfig()
		conf.Level = c.ZapLoggerLevel()
		conf.DisableCaller = true
		return conf.Build()
	} else {
		conf := zap.NewDevelopmentConfig()
		conf.Level = c.ZapLoggerLevel()
		conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return conf.Build()
	}
}

func main() {
	println(fmt.Sprintf(banner, BuildVersion()))

	flag.StringVar(&configFile, "config", "config.yml", "path to the config file")
	flag.BoolVar(&showVersion, "version", false, "show program version")
	flag.BoolVar(&healthcheck, "healthcheck", false, "run a healthcheck")
	flag.Parse()
	if showVersion {
		println(BuildVersion())
		return
	}

	cfg := config.Load(configFile)

	logger, err := makeLogger(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if healthcheck {
		addressParts := strings.Split(cfg.Server.Address, ":")
		if len(addressParts) != 2 {
			logger.Fatal("invalid server address", zap.String("address", cfg.Server.Address))
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

	logger.Info("starting program", zap.String("mode", cfg.Environment))

	defer logger.Sync()

	app := &app.App{
		Config: cfg,
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
