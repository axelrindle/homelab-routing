package main

import (
	"flag"
	"fmt"
	"log"

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

var showVersion bool

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

	config := &config.Config{}
	config.Load()

	logger, err := makeLogger(config)
	if err != nil {
		log.Fatal(err)
	}

	flag.BoolVar(&showVersion, "version", false, "show program version")
	flag.Parse()
	if showVersion {
		println(BuildVersion())
		return
	}

	logger.Info("starting program", zap.String("mode", config.Environment))

	defer logger.Sync()

	app := &app.App{
		Config: config,
		Logger: logger,
	}
	app.Init()
}
